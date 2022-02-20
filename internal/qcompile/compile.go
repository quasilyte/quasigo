package qcompile

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"math"

	"github.com/quasilyte/quasigo/internal/bytecode"
	"github.com/quasilyte/quasigo/internal/ir"
	"github.com/quasilyte/quasigo/internal/qopt"
	"github.com/quasilyte/quasigo/internal/qruntime"
)

// TODO/ideas.
//
// - Add ret slots? Could use 1 less temp in some functions.

var voidType = &types.Tuple{}

const voidSlot = math.MaxInt

type compiler struct {
	ctx *Context

	fnName  *ast.Ident
	fnKey   qruntime.FuncKey
	fnType  *types.Signature
	retType types.Type

	insideVariadic bool

	locals map[string]frameSlotInfo
	tmpSeq int
	numTmp int

	strConstantsPool    map[string]int
	scalarConstantsPool map[uint64]int

	params map[string]frameSlotInfo

	code            []ir.Inst
	strConstants    []string
	scalarConstants []uint64

	breakTarget    label
	continueTarget label

	numLabels int
}

type frameSlotInfo struct {
	i int
	v *types.Var
}

type label struct {
	id int
}

type compileError string

func (e compileError) Error() string { return string(e) }

func (cl *compiler) compileFunc(fn *ast.FuncDecl) *qruntime.Func {
	switch cl.fnType.Results().Len() {
	case 0:
		cl.retType = voidType
	case 1:
		cl.retType = cl.fnType.Results().At(0).Type()
	default:
		panic(cl.errorf(fn.Name, "multi-result functions are not supported"))
	}

	if !cl.isSupportedType(cl.retType) {
		panic(cl.errorUnsupportedType(fn.Name, cl.retType, "function result"))
	}

	dbg := qruntime.FuncDebugInfo{
		SlotNames: make([]string, 0, len(cl.params)+len(cl.locals)),
	}
	cl.fnName = fn.Name
	cl.fnKey = qruntime.FuncKey{Qualifier: cl.ctx.Package.Path(), Name: fn.Name.String()}
	cl.params = make(map[string]frameSlotInfo, cl.fnType.Params().Len())
	for i := 0; i < cl.fnType.Params().Len(); i++ {
		p := cl.fnType.Params().At(i)
		paramName := p.Name()
		paramType := p.Type()
		if !cl.isSupportedType(paramType) {
			panic(cl.errorUnsupportedType(fn.Name, paramType, paramName+" param"))
		}
		cl.params[paramName] = frameSlotInfo{i: i, v: p}
		dbg.SlotNames = append(dbg.SlotNames, paramName)
	}

	cl.collectLocals(&dbg, fn.Body)
	dbg.NumLocals = len(cl.locals)

	cl.compileStmt(fn.Body)
	if cl.retType == voidType {
		cl.emitOp(bytecode.OpReturnVoid)
	}

	irFunc := ir.Func{
		Code:          cl.code,
		NumParams:     len(cl.params),
		NumLocals:     len(cl.locals),
		NumFrameSlots: len(cl.params) + len(cl.locals) + cl.numTmp,
	}

	if cl.ctx.Optimize {
		qopt.Func(&irFunc)
	}

	var asm assembler
	code, err := asm.Assemble(&irFunc)
	if err != nil {
		panic(cl.errorf(fn.Name, "unexpected result: %s", err))
	}
	if len(code) == 0 {
		panic(cl.errorf(fn.Name, "unexpected result: 0-sized bytecode"))
	}

	compiled := &qruntime.Func{
		Code:            code,
		Codeptr:         &code[0],
		StrConstants:    cl.strConstants,
		ScalarConstants: cl.scalarConstants,
		Name:            cl.fnKey.String(),
		FrameSize:       int(qruntime.SizeofSlot) * irFunc.NumFrameSlots,
		FrameSlots:      byte(irFunc.NumFrameSlots),
	}
	cl.ctx.Env.Debug.Funcs[compiled] = dbg

	// Now that we know the frame size, we need to fix the arguments passing offsets.
	// TODO: do this on IR level?
	bytecode.Walk(code, func(pc int, op bytecode.Op) {
		if !op.HasDst() {
			return
		}
		dstslot := int8(code[pc+1])
		if dstslot < 0 {
			actualIndex := -dstslot - 1
			code[pc+1] = byte(irFunc.NumFrameSlots + int(actualIndex))
		}
	})

	return compiled
}

func (cl *compiler) collectLocals(dbg *qruntime.FuncDebugInfo, body *ast.BlockStmt) {
	ast.Inspect(body, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok || assign.Tok != token.DEFINE {
			return true
		}
		for _, lhs := range assign.Lhs {
			lhs, ok := lhs.(*ast.Ident)
			if !ok {
				continue
			}
			def, ok := cl.ctx.Types.Defs[lhs].(*types.Var)
			if !ok || def == nil {
				continue
			}
			if _, ok := cl.locals[lhs.String()]; ok {
				panic(cl.errorf(lhs, "%s variable shadowing is not allowed", lhs))
			}
			typ := cl.ctx.Types.TypeOf(lhs)
			if !cl.isSupportedType(typ) {
				panic(cl.errorUnsupportedType(lhs, typ, lhs.String()+" local variable"))
			}
			cl.locals[lhs.Name] = frameSlotInfo{
				i: len(cl.params) + len(cl.locals),
				v: def,
			}
			dbg.SlotNames = append(dbg.SlotNames, lhs.Name)
		}
		return true
	})
}

func (cl *compiler) newLabel() label {
	if cl.numLabels >= 255 {
		panic(cl.errorf(cl.fnName, "too many labels"))
	}
	l := label{id: cl.numLabels}
	cl.numLabels++
	return l
}

func (cl *compiler) bindLabel(l label) {
	cl.emit(ir.Inst{
		Pseudo: ir.OpLabel,
		Arg0:   uint8(l.id),
	})
}

func (cl *compiler) emit(inst ir.Inst) {
	cl.code = append(cl.code, inst)
}

func (cl *compiler) emitOp(op bytecode.Op) {
	cl.code = append(cl.code, ir.Inst{Op: op})
}

func (cl *compiler) emitJump(l label) {
	cl.emit(ir.Inst{Op: bytecode.OpJump, Arg0: uint8(l.id)})
}

func (cl *compiler) emitCondJump(slot int, op bytecode.Op, l label) {
	cl.emit(ir.Inst{
		Op:   op,
		Arg0: uint8(l.id),
		Arg1: uint8(slot),
	})
}

func (cl *compiler) emit1(op bytecode.Op, arg8 int) {
	cl.emit(ir.Inst{Op: op, Arg0: uint8(arg8)})
}

func (cl *compiler) emit2(op bytecode.Op, arg8a, arg8b int) {
	cl.emit(ir.Inst{Op: op, Arg0: uint8(arg8a), Arg1: uint8(arg8b)})
}

func (cl *compiler) emit3(op bytecode.Op, arg8a, arg8b, arg8c int) {
	cl.emit(ir.Inst{Op: op, Arg0: uint8(arg8a), Arg1: uint8(arg8b), Arg2: uint8(arg8c), Arg3: uint8(arg8c)})
}

func (cl *compiler) emit4(op bytecode.Op, arg8a, arg8b, arg8c, arg8d int) {
	cl.emit(ir.Inst{
		Op:   op,
		Arg0: uint8(arg8a),
		Arg1: uint8(arg8b),
		Arg2: uint8(arg8c),
		Arg3: uint8(arg8d),
	})
}

func (cl *compiler) emitCall(op bytecode.Op, dst int, funcid int) {
	if dst == voidSlot && op == bytecode.OpCallNative {
		cl.emit(ir.Inst{Op: bytecode.OpCallVoidNative, Value: uint16(funcid)})
		return
	}
	if dst == voidSlot && op == bytecode.OpCall {
		cl.emit(ir.Inst{Op: bytecode.OpCallVoid, Value: uint16(funcid)})
		return
	}
	cl.emit(ir.Inst{
		Op:    op,
		Value: uint16(funcid),
		Arg0:  uint8(dst),
	})
}

func (cl *compiler) errorUnsupportedType(e ast.Node, typ types.Type, where string) compileError {
	return cl.errorf(e, "%s type: %s is not supported, try something simpler", where, typ)
}

func (cl *compiler) errorf(n ast.Node, format string, args ...interface{}) compileError {
	loc := cl.ctx.Fset.Position(n.Pos())
	message := fmt.Sprintf("%s:%d: %s", loc.Filename, loc.Line, fmt.Sprintf(format, args...))
	return compileError(message)
}

func (cl *compiler) lastOp() bytecode.Op {
	if len(cl.code) != 0 {
		return cl.code[len(cl.code)-1].Op
	}
	return bytecode.OpInvalid
}

func (cl *compiler) isUncondJump(op bytecode.Op) bool {
	switch op {
	case bytecode.OpJump, bytecode.OpReturnFalse, bytecode.OpReturnTrue, bytecode.OpReturnStr, bytecode.OpReturnScalar:
		return true
	default:
		return false
	}
}

func (cl *compiler) isSupportedType(typ types.Type) bool {
	if typ == voidType {
		return true
	}

	switch typ := typ.Underlying().(type) {
	case *types.Pointer:
		// 1. Pointers to structs are supported.
		_, isStruct := typ.Elem().Underlying().(*types.Struct)
		return isStruct

	case *types.Basic:
		// 2. Some of the basic types are supported.
		// TODO: support byte/uint8 and maybe float64.
		switch typ.Kind() {
		case types.Bool, types.Int, types.String:
			return true
		default:
			return false
		}

	case *types.Interface:
		// 3. Interfaces are supported.
		return true

	default:
		return false
	}
}

func (cl *compiler) opMoveByType(e ast.Expr, typ types.Type) bytecode.Op {
	switch {
	case typeIsScalar(typ):
		return bytecode.OpMoveScalar
	case typeIsString(typ):
		return bytecode.OpMoveStr
	case typeIsInterface(typ) || typeIsPointer(typ):
		return bytecode.OpMoveInterface
	default:
		panic(cl.errorf(e, "can't move %s typed value yet", typ.String()))
	}
}

func (cl *compiler) internBoolConstant(v bool) int {
	if v {
		return cl.internScalarConstant(1)
	}
	return cl.internScalarConstant(0)
}

func (cl *compiler) internIntConstant(v int) int {
	return cl.internScalarConstant(uint64(v))
}

func (cl *compiler) internScalarConstant(v uint64) int {
	if id, ok := cl.scalarConstantsPool[v]; ok {
		return id
	}
	id := len(cl.scalarConstants)
	cl.scalarConstants = append(cl.scalarConstants, v)
	cl.scalarConstantsPool[v] = id
	return id
}

func (cl *compiler) internStrConstant(s string) int {
	if id, ok := cl.strConstantsPool[s]; ok {
		return id
	}
	id := len(cl.strConstants)
	cl.strConstants = append(cl.strConstants, s)
	cl.strConstantsPool[s] = id
	return id
}

func (cl *compiler) isParamName(varname string) bool {
	_, ok := cl.params[varname]
	return ok
}

func (cl *compiler) getLocal(v ast.Expr, varname string) int {
	slot, ok := cl.locals[varname]
	if !ok {
		if cl.isParamName(varname) {
			panic(cl.errorf(v, "can't assign to %s, params are readonly", varname))
		}
		panic(cl.errorf(v, "%s is not a writeable local variable", varname))
	}
	return slot.i
}

func (cl *compiler) freeTmp() {
	cl.tmpSeq = 0
}

func (cl *compiler) allocTmp() int {
	index := cl.tmpSeq
	cl.tmpSeq++
	if cl.numTmp < cl.tmpSeq {
		cl.numTmp = cl.tmpSeq
	}
	return index + len(cl.params) + len(cl.locals)
}

func (cl *compiler) isSimpleExpr(e ast.Expr) bool {
	switch e := e.(type) {
	case *ast.ParenExpr:
		return cl.isSimpleExpr(e.X)
	case *ast.Ident, *ast.BasicLit:
		return true
	case *ast.IndexExpr:
		return cl.isSimpleExpr(e.X) && cl.isSimpleExpr(e.Index)
	default:
		return false
	}
}
