package qcompile

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/quasilyte/quasigo/internal/bytecode"
	"github.com/quasilyte/quasigo/internal/ir"
	"github.com/quasilyte/quasigo/internal/qopt"
	"github.com/quasilyte/quasigo/internal/qruntime"
)

var voidType = &types.Tuple{}

var voidSlot = ir.Slot{Kind: ir.SlotDiscard}

type funcCompiler struct {
	ctx *Context

	currentFunc ir.Func
	optimizer   qopt.Optimizer
	scope       scope

	fnName  *ast.Ident
	fnKey   qruntime.FuncKey
	fnType  *types.Signature
	retType types.Type

	insideVariadic bool

	hasCalls bool
	hasLoops bool

	tempSeq     int
	numTemp     int
	inTempBlock bool

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

func (cl *funcCompiler) buildIR(fn *ast.FuncDecl) *ir.Func {
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
		SlotNames: make([]string, 0, len(cl.params)),
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

	cl.compileStmt(fn.Body)
	if cl.retType == voidType {
		cl.emitOp(bytecode.OpReturnVoid)
	}

	cl.currentFunc = ir.Func{
		Name:            cl.fnKey.String(),
		Code:            cl.code,
		NumParams:       len(cl.params),
		NumTemps:        cl.numTemp,
		StrConstants:    cl.strConstants,
		ScalarConstants: cl.scalarConstants,
		Debug:           dbg,
		Env:             cl.ctx.Env,
	}

	return &cl.currentFunc
}

func (cl *funcCompiler) compileFunc(fn *ast.FuncDecl) *qruntime.Func {
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

	if len(cl.scope.vars) != 0 || len(cl.scope.depths) != 0 {
		panic("internal error: lexical scope is not empty")
	}

	compiled := &qruntime.Func{
		Code:            code,
		Codeptr:         &code[0],
		StrConstants:    irFunc.StrConstants,
		ScalarConstants: irFunc.ScalarConstants,
		Name:            irFunc.Name,
		FrameSize:       int(qruntime.SizeofSlot) * irFunc.NumFrameSlots(),
		FrameSlots:      byte(irFunc.NumFrameSlots()),
		NumParams:       byte(len(cl.params)),
		NumTemps:        byte(irFunc.NumTemps),
		CanInline:       cl.canInline(irFunc),
	}
	cl.ctx.Env.Debug.Funcs[compiled] = irFunc.Debug
	return compiled
}

func (cl *funcCompiler) canInline(fn *ir.Func) bool {
	return cl.ctx.Static &&
		!cl.hasCalls && !cl.hasLoops &&
		cl.numLabels <= 6 &&
		fn.NumFrameSlots() <= 16 &&
		len(fn.ScalarConstants) <= 8 &&
		len(fn.StrConstants) <= 8
}

func (cl *funcCompiler) emitVarKill(id int) {
	cl.emit(ir.Inst{
		Pseudo: ir.OpVarKill,
		Arg0:   ir.NewTempSlot(uint8(id)).ToInstArg(),
	})
}

func (cl *funcCompiler) newLabel() label {
	if cl.numLabels >= 255 {
		panic(cl.errorf(cl.fnName, "too many labels"))
	}
	l := label{id: cl.numLabels}
	cl.numLabels++
	return l
}

func (cl *funcCompiler) bindLabel(l label) {
	cl.emit(ir.Inst{
		Pseudo: ir.OpLabel,
		Arg0:   ir.InstArg(l.id),
	})
}

func (cl *funcCompiler) emit(inst ir.Inst) {
	cl.code = append(cl.code, inst)
}

func (cl *funcCompiler) emitOp(op bytecode.Op) {
	cl.code = append(cl.code, ir.Inst{Op: op})
}

func (cl *funcCompiler) emitJump(l label) {
	cl.emit(ir.Inst{Op: bytecode.OpJump, Arg0: ir.InstArg(l.id)})
}

func (cl *funcCompiler) emitCondJump(slot ir.Slot, op bytecode.Op, l label) {
	cl.emit(ir.Inst{
		Op:   op,
		Arg0: ir.InstArg(l.id),
		Arg1: slot.ToInstArg(),
	})
}

func (cl *funcCompiler) emit1(op bytecode.Op, a0 ir.Slot) {
	cl.emit(ir.Inst{Op: op, Arg0: a0.ToInstArg()})
}

func (cl *funcCompiler) emit2(op bytecode.Op, a0, a1 ir.Slot) {
	cl.emit(ir.Inst{Op: op, Arg0: a0.ToInstArg(), Arg1: a1.ToInstArg()})
}

func (cl *funcCompiler) emit3(op bytecode.Op, a0, a1, a2 ir.Slot) {
	cl.emit(ir.Inst{
		Op:   op,
		Arg0: a0.ToInstArg(),
		Arg1: a1.ToInstArg(),
		Arg2: a2.ToInstArg(),
	})
}

func (cl *funcCompiler) emit4(op bytecode.Op, a0, a1, a2, a3 ir.Slot) {
	cl.emit(ir.Inst{
		Op:   op,
		Arg0: a0.ToInstArg(),
		Arg1: a1.ToInstArg(),
		Arg2: a2.ToInstArg(),
		Arg3: a3.ToInstArg(),
	})
}

func (cl *funcCompiler) emitCall(op bytecode.Op, dst ir.Slot, funcid int) {
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

func (cl *funcCompiler) fatalf(format string, args ...interface{}) {
	loc := cl.ctx.Fset.Position(cl.fnName.Pos())
	panic(fmt.Sprintf("%s:%d: internal error: %s", loc.Filename, loc.Line, fmt.Sprintf(format, args...)))
}

func (cl *funcCompiler) errorUnsupportedType(e ast.Node, typ types.Type, where string) compileError {
	return cl.errorf(e, "%s type: %s is not supported, try something simpler", where, typ)
}

func (cl *funcCompiler) errorf(n ast.Node, format string, args ...interface{}) compileError {
	loc := cl.ctx.Fset.Position(n.Pos())
	message := fmt.Sprintf("%s:%d: %s", loc.Filename, loc.Line, fmt.Sprintf(format, args...))
	return compileError(message)
}

func (cl *funcCompiler) lastOp() bytecode.Op {
	for i := len(cl.code) - 1; i >= 0; i-- {
		if cl.code[i].Op == bytecode.OpInvalid {
			continue
		}
		return cl.code[i].Op
	}
	return bytecode.OpInvalid
}

func (cl *funcCompiler) isUncondJump(op bytecode.Op) bool {
	switch op {
	case bytecode.OpJump, bytecode.OpReturnZero, bytecode.OpReturnOne, bytecode.OpReturnStr, bytecode.OpReturnScalar:
		return true
	default:
		return false
	}
}

func (cl *funcCompiler) isSupportedType(typ types.Type) bool {
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

func (cl *funcCompiler) moveBool(dst ir.Slot, v bool) ir.Inst {
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

func (cl *funcCompiler) moveInt(dst ir.Slot, v int) ir.Inst {
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

func (cl *funcCompiler) internBoolConstant(v bool) int {
	if v {
		return cl.internScalarConstant(1)
	}
	return cl.internScalarConstant(0)
}

func (cl *funcCompiler) internIntConstant(v int) int {
	return cl.internScalarConstant(uint64(v))
}

func (cl *funcCompiler) internScalarConstant(v uint64) int {
	if id, ok := cl.scalarConstantsPool[v]; ok {
		return id
	}
	id := len(cl.scalarConstants)
	cl.scalarConstants = append(cl.scalarConstants, v)
	cl.scalarConstantsPool[v] = id
	return id
}

func (cl *funcCompiler) internStrConstant(s string) int {
	if id, ok := cl.strConstantsPool[s]; ok {
		return id
	}
	id := len(cl.strConstants)
	cl.strConstants = append(cl.strConstants, s)
	cl.strConstantsPool[s] = id
	return id
}

func (cl *funcCompiler) defineOrLookupVar(e ast.Expr, varname string, define bool) ir.Slot {
	if !define {
		return cl.getNamedSlot(e, varname)
	}

	// Can re-define a variable only if it's not from a current scope level.
	if i := cl.scope.LookupInCurrent(varname); i != -1 {
		return ir.NewTempSlot(uint8(i))
	}
	// Can re-define a parameter only if it's not a root function scope.
	if cl.scope.NumLevels() == 1 {
		if p, ok := cl.params[varname]; ok {
			return p.i
		}
	}

	slot := cl.allocTemp()
	cl.scope.PushVar(varname)
	return slot
}

func (cl *funcCompiler) getNamedSlot(v ast.Expr, varname string) ir.Slot {
	if i := cl.scope.Lookup(varname); i != -1 {
		return ir.NewTempSlot(uint8(i))
	}
	if p, ok := cl.params[varname]; ok {
		return p.i
	}
	panic(cl.errorf(v, "%s is not a writeable local variable", varname))
}

func (cl *funcCompiler) beginTempBlock() {
	if cl.inTempBlock {
		cl.fatalf("nested beginTempBlock call")
	}
	cl.inTempBlock = true
}

func (cl *funcCompiler) endTempBlock() {
	if !cl.inTempBlock {
		cl.fatalf("endTempBlock without beginTempBlock")
	}
	cl.killScopeVars(cl.tempSeq - cl.scope.NumLiveVars())
	cl.inTempBlock = false
}

func (cl *funcCompiler) killScopeVars(num int) {
	for i := 0; i < num; i++ {
		cl.tempSeq--
		cl.emitVarKill(cl.tempSeq)
	}
}

func (cl *funcCompiler) trackTemp(id int) {
	if cl.numTemp < id+1 {
		cl.numTemp = id + 1
	}
}

func (cl *funcCompiler) allocTemp() ir.Slot {
	id := cl.tempSeq
	cl.tempSeq++
	cl.trackTemp(id)
	return ir.NewTempSlot(uint8(id))
}

func (cl *funcCompiler) isSimpleExpr(e ast.Expr) bool {
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
