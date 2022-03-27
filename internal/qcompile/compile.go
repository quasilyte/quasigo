package qcompile

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/quasilyte/quasigo/internal/bytecode"
	"github.com/quasilyte/quasigo/internal/ir"
	"github.com/quasilyte/quasigo/internal/qopt"
	"github.com/quasilyte/quasigo/internal/qruntime"
)

var voidType = &types.Tuple{}

var voidSlot = ir.Slot{Kind: ir.SlotDiscard}

type compiler struct {
	ctx *Context

	currentFunc ir.Func
	optimizer   qopt.Optimizer

	fnName  *ast.Ident
	fnKey   qruntime.FuncKey
	fnType  *types.Signature
	retType types.Type

	insideVariadic bool

	hasCalls bool
	hasLoops bool

	locals       map[string]frameSlotInfo
	autoLocalSeq int
	numAutoLocal int
	tempSeq      int
	numTemp      int
	inTempBlock  bool

	strConstantsPool    map[string]int
	scalarConstantsPool map[uint64]int

	params map[string]frameSlotInfo

	code            []ir.Inst
	strConstants    []string
	scalarConstants []uint64

	breakTarget    label
	continueTarget label

	numLabels int

	patternCompiler *patternCompiler
}

type frameSlotInfo struct {
	i ir.Slot
	v *types.Var
}

type label struct {
	id int
}

type compileError string

func (e compileError) Error() string { return string(e) }

func (cl *compiler) buildIR(fn *ast.FuncDecl) *ir.Func {
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
	cl.hasCalls = false
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
		cl.params[paramName] = frameSlotInfo{
			i: ir.NewParamSlot(uint8(i)),
			v: p,
		}
		dbg.SlotNames = append(dbg.SlotNames, paramName)
	}

	cl.collectLocals(&dbg, fn.Body)

	cl.compileStmt(fn.Body)
	if cl.retType == voidType {
		cl.emitOp(bytecode.OpReturnVoid)
	}

	for i := 0; i < cl.numAutoLocal; i++ {
		dbg.SlotNames = append(dbg.SlotNames, fmt.Sprintf("auto%d", i))
	}
	cl.currentFunc = ir.Func{
		Name:            cl.fnKey.String(),
		Code:            cl.code,
		NumParams:       len(cl.params),
		NumLocals:       len(cl.locals) + cl.numAutoLocal,
		NumFrameSlots:   len(cl.params) + len(cl.locals) + cl.numAutoLocal + cl.numTemp,
		StrConstants:    cl.strConstants,
		ScalarConstants: cl.scalarConstants,
		Debug:           dbg,
		Env:             cl.ctx.Env,
	}

	return &cl.currentFunc
}

func (cl *compiler) compileFunc(fn *ast.FuncDecl) *qruntime.Func {
	irFunc := cl.buildIR(fn)

	if cl.ctx.Optimize {
		cl.optimizer.PrepareFunc(irFunc)
		if cl.ctx.TestingContext != nil {
			cl.ctx.TestingContext.FuncIR(irFunc)
		}
		cl.optimizer.OptimizePrepared()
	}

	var asm assembler
	code, err := asm.Assemble(irFunc)
	if err != nil {
		panic(cl.errorf(fn.Name, "unexpected result: %s", err))
	}
	if len(code) == 0 {
		panic(cl.errorf(fn.Name, "unexpected result: 0-sized bytecode"))
	}

	if cl.autoLocalSeq != 0 {
		panic("internal error: leaking auto locals?")
	}

	compiled := &qruntime.Func{
		Code:            code,
		Codeptr:         &code[0],
		StrConstants:    irFunc.StrConstants,
		ScalarConstants: irFunc.ScalarConstants,
		Name:            irFunc.Name,
		FrameSize:       int(qruntime.SizeofSlot) * irFunc.NumFrameSlots,
		FrameSlots:      byte(irFunc.NumFrameSlots),
		NumParams:       byte(len(cl.params)),
		NumLocals:       byte(len(cl.locals) + cl.numAutoLocal),
		CanInline:       cl.canInline(irFunc),
	}
	cl.ctx.Env.Debug.Funcs[compiled] = irFunc.Debug
	return compiled
}

func (cl *compiler) canInline(fn *ir.Func) bool {
	return cl.ctx.Static &&
		!cl.hasCalls && !cl.hasLoops &&
		cl.numLabels <= 6 &&
		fn.NumFrameSlots <= 16 &&
		len(fn.ScalarConstants) <= 8 &&
		len(fn.StrConstants) <= 8
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
			id := len(cl.locals)
			cl.locals[lhs.Name] = frameSlotInfo{
				i: ir.NewLocalSlot(uint8(id)),
				v: def,
			}
			dbg.SlotNames = append(dbg.SlotNames, lhs.Name)
		}
		return true
	})
}

func (cl *compiler) emitVarKill(id int) {
	cl.emit(ir.Inst{
		Pseudo: ir.OpVarKill,
		Arg0:   ir.NewTempSlot(uint8(id)).ToInstArg(),
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
		Arg0:   ir.InstArg(l.id),
	})
}

func (cl *compiler) emit(inst ir.Inst) {
	cl.code = append(cl.code, inst)
}

func (cl *compiler) emitOp(op bytecode.Op) {
	cl.code = append(cl.code, ir.Inst{Op: op})
}

func (cl *compiler) emitJump(l label) {
	cl.emit(ir.Inst{Op: bytecode.OpJump, Arg0: ir.InstArg(l.id)})
}

func (cl *compiler) emitCondJump(slot ir.Slot, op bytecode.Op, l label) {
	cl.emit(ir.Inst{
		Op:   op,
		Arg0: ir.InstArg(l.id),
		Arg1: slot.ToInstArg(),
	})
}

func (cl *compiler) emit1(op bytecode.Op, a0 ir.Slot) {
	cl.emit(ir.Inst{Op: op, Arg0: a0.ToInstArg()})
}

func (cl *compiler) emit2(op bytecode.Op, a0, a1 ir.Slot) {
	cl.emit(ir.Inst{Op: op, Arg0: a0.ToInstArg(), Arg1: a1.ToInstArg()})
}

func (cl *compiler) emit3(op bytecode.Op, a0, a1, a2 ir.Slot) {
	cl.emit(ir.Inst{
		Op:   op,
		Arg0: a0.ToInstArg(),
		Arg1: a1.ToInstArg(),
		Arg2: a2.ToInstArg(),
	})
}

func (cl *compiler) emit4(op bytecode.Op, a0, a1, a2, a3 ir.Slot) {
	cl.emit(ir.Inst{
		Op:   op,
		Arg0: a0.ToInstArg(),
		Arg1: a1.ToInstArg(),
		Arg2: a2.ToInstArg(),
		Arg3: a3.ToInstArg(),
	})
}

func (cl *compiler) emitCall(op bytecode.Op, dst ir.Slot, funcid int) {
	if dst == voidSlot && op == bytecode.OpCallNative {
		cl.emit(ir.Inst{Op: bytecode.OpCallVoidNative, Arg0: ir.InstArg(funcid)})
		return
	}
	if dst == voidSlot && op == bytecode.OpCall {
		cl.emit(ir.Inst{Op: bytecode.OpCallVoid, Arg0: ir.InstArg(funcid)})
		return
	}
	cl.hasCalls = true
	cl.emit(ir.Inst{
		Op:   op,
		Arg0: dst.ToInstArg(),
		Arg1: ir.InstArg(funcid),
	})
}

func (cl *compiler) fatalf(format string, args ...interface{}) {
	loc := cl.ctx.Fset.Position(cl.fnName.Pos())
	panic(fmt.Sprintf("%s:%d: internal error: %s", loc.Filename, loc.Line, fmt.Sprintf(format, args...)))
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
	for i := len(cl.code) - 1; i >= 0; i-- {
		if cl.code[i].Op == bytecode.OpInvalid {
			continue
		}
		return cl.code[i].Op
	}
	return bytecode.OpInvalid
}

func (cl *compiler) isUncondJump(op bytecode.Op) bool {
	switch op {
	case bytecode.OpJump, bytecode.OpReturnZero, bytecode.OpReturnOne, bytecode.OpReturnStr, bytecode.OpReturnScalar:
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
		// TODO: float64.
		switch typ.Kind() {
		case types.Bool, types.Int, types.String, types.Uint8:
			return true
		default:
			return false
		}

	case *types.Interface:
		// 3. Interfaces are supported.
		return true

	case *types.Slice:
		// 4. Slices are supported as long as their elem type is supported.
		return cl.isSupportedType(typ.Elem())

	default:
		return false
	}
}

func (cl *compiler) moveBool(dst ir.Slot, v bool) ir.Inst {
	if v {
		id := cl.internBoolConstant(true)
		return ir.Inst{
			Op:   bytecode.OpLoadScalarConst,
			Arg0: dst.ToInstArg(),
			Arg1: ir.InstArg(id),
		}
	}
	return ir.Inst{Op: bytecode.OpZero, Arg0: dst.ToInstArg()}
}

func (cl *compiler) moveInt(dst ir.Slot, v int) ir.Inst {
	if v != 0 {
		id := cl.internIntConstant(v)
		return ir.Inst{
			Op:   bytecode.OpLoadScalarConst,
			Arg0: dst.ToInstArg(),
			Arg1: ir.InstArg(id),
		}
	}
	return ir.Inst{Op: bytecode.OpZero, Arg0: dst.ToInstArg()}
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

func (cl *compiler) getNamedSlot(v ast.Expr, varname string) ir.Slot {
	if p, ok := cl.params[varname]; ok {
		return p.i
	}
	if l, ok := cl.locals[varname]; ok {
		return l.i
	}
	panic(cl.errorf(v, "%s is not a writeable local variable", varname))
}

func (cl *compiler) beginTempBlock() {
	if cl.tempSeq != 0 {
		cl.fatalf("beginTempBlock with non-zero temp seq")
	}
	if cl.inTempBlock {
		cl.fatalf("nested beginTempBlock call")
	}
	cl.inTempBlock = true
}

func (cl *compiler) endTempBlock() {
	if !cl.inTempBlock {
		cl.fatalf("endTempBlock without beginTempBlock")
	}
	for cl.tempSeq > 0 {
		cl.tempSeq--
		cl.emitVarKill(cl.tempSeq)
	}
	cl.inTempBlock = false
}

func (cl *compiler) trackTemp(id int) {
	if cl.numTemp < id+1 {
		cl.numTemp = id + 1
	}
}

func (cl *compiler) allocTemp() ir.Slot {
	id := cl.tempSeq
	cl.tempSeq++
	cl.trackTemp(id)
	return ir.NewTempSlot(uint8(id))
}

func (cl *compiler) allocAutoLocal() ir.Slot {
	id := cl.autoLocalSeq + len(cl.locals)
	cl.autoLocalSeq++
	if cl.numAutoLocal < cl.autoLocalSeq {
		cl.numAutoLocal = cl.autoLocalSeq
	}
	return ir.NewLocalSlot(uint8(id))
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
