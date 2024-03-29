package qcompile

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"strconv"

	"github.com/quasilyte/quasigo/internal/goutil"
	"github.com/quasilyte/quasigo/internal/ir"
	"github.com/quasilyte/quasigo/internal/qruntime"

	"github.com/quasilyte/quasigo/internal/bytecode"
)

func (cl *funcCompiler) compileTempExpr(e ast.Expr) ir.Slot {
	if v, ok := e.(*ast.Ident); ok {
		if i := cl.scope.Lookup(v.Name); i != -1 {
			return ir.NewTempSlot(uint8(i))
		}
		if p, ok := cl.params[v.Name]; ok {
			return p.i
		}
	}
	temp := cl.allocTemp()
	cl.CompileExpr(temp, e)
	return temp
}

func (cl *funcCompiler) compileNestedTempExpr(e ast.Expr) ir.Slot {
	if v, ok := e.(*ast.Ident); ok {
		if i := cl.scope.Lookup(v.Name); i != -1 {
			return ir.NewTempSlot(uint8(i))
		}
		if p, ok := cl.params[v.Name]; ok {
			return p.i
		}
	}
	temp := cl.allocTemp()
	tempSeq := cl.tempSeq
	cl.CompileExpr(temp, e)
	tempAllocated := cl.tempSeq - tempSeq
	cl.killScopeVars(tempAllocated)
	return temp
}

func (cl *funcCompiler) compileRootExpr(dst ir.Slot, e ast.Expr) {
	cl.beginTempBlock()
	cl.CompileExpr(dst, e)
	cl.endTempBlock()
}

func (cl *funcCompiler) compileNestedRootExpr(dst ir.Slot, e ast.Expr) {
	if !cl.inTempBlock {
		cl.compileRootExpr(dst, e)
		return
	}
	tempSeq := cl.tempSeq
	cl.CompileExpr(dst, e)
	tempAllocated := cl.tempSeq - tempSeq
	cl.killScopeVars(tempAllocated)
}

func (cl *funcCompiler) CompileExpr(dst ir.Slot, e ast.Expr) {
	cv := cl.ctx.Types.Types[e].Value
	if cv != nil {
		cl.compileConstantValue(dst, e, cv)
		return
	}

	switch e := e.(type) {
	case *ast.ParenExpr:
		cl.CompileExpr(dst, e.X)

	case *ast.BasicLit:
		cl.compileBasicLit(dst, e)

	case *ast.Ident:
		cl.compileIdent(dst, e)

	case *ast.UnaryExpr:
		cl.compileUnaryExpr(dst, e)

	case *ast.BinaryExpr:
		cl.compileBinaryExpr(dst, e)

	case *ast.SliceExpr:
		cl.compileSliceExpr(dst, e)

	case *ast.IndexExpr:
		cl.compileIndexExpr(dst, e)

	case *ast.SelectorExpr:
		cl.compileSelectorExpr(dst, e)

	case *ast.CallExpr:
		cl.compileCallExpr(dst, e)

	default:
		panic(cl.errorf(e, "can't compile %T yet", e))
	}
}

func (cl *funcCompiler) compileUnaryExpr(dst ir.Slot, e *ast.UnaryExpr) {
	switch e.Op {
	case token.NOT:
		cl.compileUnaryOp(dst, bytecode.OpNot, e.X)

	case token.SUB:
		if typeIsFloat(cl.ctx.Types.TypeOf(e.X)) {
			cl.compileUnaryOp(dst, bytecode.OpFloatNeg, e.X)
		} else {
			cl.compileUnaryOp(dst, bytecode.OpIntNeg, e.X)
		}

	case token.XOR:
		cl.compileUnaryOp(dst, bytecode.OpIntBitwiseNot, e.X)

	default:
		panic(cl.errorf(e, "can't compile unary %s yet", e.Op))
	}
}

func (cl *funcCompiler) compileUnaryOp(dst ir.Slot, op bytecode.Op, arg ast.Expr) {
	xslot := cl.compileTempExpr(arg)
	cl.emit2(op, dst, xslot)
}

func (cl *funcCompiler) compileBinaryExpr(dst ir.Slot, e *ast.BinaryExpr) {
	typ := cl.ctx.Types.TypeOf(e.X)

	switch e.Op {
	case token.LOR:
		cl.compileOr(dst, e)
	case token.LAND:
		cl.compileAnd(dst, e)

	case token.NEQ:
		switch {
		case identName(e.X) == "nil":
			cl.compileUnaryOp(dst, pickOp(typeIsInterface(cl.ctx.Types.TypeOf(e.Y)), bytecode.OpIsNotNilInterface, bytecode.OpIsNotNil), e.Y)
		case identName(e.Y) == "nil":
			cl.compileUnaryOp(dst, pickOp(typeIsInterface(typ), bytecode.OpIsNotNilInterface, bytecode.OpIsNotNil), e.X)

		case typeIsString(typ):
			cl.compileBinaryOp(dst, bytecode.OpStrNotEq, e)
		case typeIsScalar(typ):
			cl.compileBinaryOp(dst, bytecode.OpScalarNotEq, e)
		default:
			panic(cl.errorf(e, "!= is not implemented for %s bytecode.Operands", typ))
		}
	case token.EQL:
		switch {
		case identName(e.X) == "nil":
			cl.compileUnaryOp(dst, pickOp(typeIsInterface(cl.ctx.Types.TypeOf(e.Y)), bytecode.OpIsNilInterface, bytecode.OpIsNil), e.Y)
		case identName(e.Y) == "nil":
			cl.compileUnaryOp(dst, pickOp(typeIsInterface(typ), bytecode.OpIsNilInterface, bytecode.OpIsNil), e.X)

		case typeIsString(cl.ctx.Types.TypeOf(e.X)):
			cl.compileBinaryOp(dst, bytecode.OpStrEq, e)
		case typeIsScalar(cl.ctx.Types.TypeOf(e.X)):
			cl.compileBinaryOp(dst, bytecode.OpScalarEq, e)
		default:
			panic(cl.errorf(e, "== is not implemented for %s bytecode.Operands", typ))
		}

	case token.GTR:
		switch {
		case typeIsByte(typ) || typeIsInt(typ):
			cl.compileBinaryOp(dst, bytecode.OpIntGt, e)
		case typeIsFloat(typ):
			cl.compileBinaryOp(dst, bytecode.OpFloatGt, e)
		case typeIsString(typ):
			cl.compileBinaryOp(dst, bytecode.OpStrGt, e)
		default:
			panic(cl.errorf(e, "> is not implemented for %s bytecode.Operands", typ))
		}
	case token.LSS:
		switch {
		case typeIsByte(typ) || typeIsInt(typ):
			cl.compileBinaryOp(dst, bytecode.OpIntLt, e)
		case typeIsFloat(typ):
			cl.compileBinaryOp(dst, bytecode.OpFloatLt, e)
		case typeIsString(typ):
			cl.compileBinaryOp(dst, bytecode.OpStrLt, e)
		default:
			panic(cl.errorf(e, "< is not implemented for %s bytecode.Operands", typ))
		}
	case token.GEQ:
		cl.compileScalarBinaryOp(dst, e, bytecode.OpIntGtEq, typ)
	case token.LEQ:
		cl.compileScalarBinaryOp(dst, e, bytecode.OpIntLtEq, typ)

	case token.SHL:
		cl.compileBinaryOp(dst, bytecode.OpIntLshift, e)
	case token.SHR:
		cl.compileBinaryOp(dst, bytecode.OpIntRshift, e)

	case token.OR:
		cl.compileBinaryOp(dst, bytecode.OpIntOr, e)

	case token.ADD:
		switch {
		case typeIsString(typ):
			cl.compileBinaryOp(dst, bytecode.OpConcat, e)
		case typeIsByte(typ):
			cl.compileBinaryOp(dst, bytecode.OpIntAdd8, e)
		case typeIsInt(typ):
			cl.compileBinaryOp(dst, bytecode.OpIntAdd64, e)
		case typeIsFloat(typ):
			cl.compileBinaryOp(dst, bytecode.OpFloatAdd64, e)
		default:
			panic(cl.errorf(e, "+ is not implemented for %s bytecode.Operands", typ))
		}

	case token.SUB:
		switch {
		case typeIsByte(typ):
			cl.compileBinaryOp(dst, bytecode.OpIntSub8, e)
		case typeIsInt(typ):
			cl.compileBinaryOp(dst, bytecode.OpIntSub64, e)
		case typeIsFloat(typ):
			cl.compileBinaryOp(dst, bytecode.OpFloatSub64, e)
		default:
			panic(cl.errorf(e, "- is not implemented for %s bytecode.Operands", typ))
		}

	case token.XOR:
		cl.compileScalarBinaryOp(dst, e, bytecode.OpIntXor, typ)

	case token.MUL:
		switch {
		case typeIsByte(typ):
			cl.compileBinaryOp(dst, bytecode.OpIntMul8, e)
		case typeIsInt(typ):
			cl.compileBinaryOp(dst, bytecode.OpIntMul64, e)
		case typeIsFloat(typ):
			cl.compileBinaryOp(dst, bytecode.OpFloatMul64, e)
		default:
			panic(cl.errorf(e, "* is not implemented for %s bytecode.Operands", typ))
		}

	case token.QUO:
		op := pickOp(typeIsFloat(typ), bytecode.OpFloatDiv64, bytecode.OpIntDiv)
		cl.compileBinaryOp(dst, op, e)
	case token.REM:
		cl.compileIntBinaryOp(dst, e, bytecode.OpIntMod, typ)

	default:
		panic(cl.errorf(e, "can't compile binary %s yet", e.Op))
	}
}

func (cl *funcCompiler) compileScalarBinaryOp(dst ir.Slot, e *ast.BinaryExpr, op bytecode.Op, typ types.Type) {
	if typeIsInt(typ) || typeIsByte(typ) {
		cl.compileBinaryOp(dst, op, e)
	} else {
		panic(cl.errorf(e, "%s is not implemented for %s bytecode.Operands", e.Op, typ))
	}
}

func (cl *funcCompiler) compileIntBinaryOp(dst ir.Slot, e *ast.BinaryExpr, op bytecode.Op, typ types.Type) {
	if typeIsInt(typ) {
		cl.compileBinaryOp(dst, op, e)
	} else {
		panic(cl.errorf(e, "%s is not implemented for %s bytecode.Operands", e.Op, typ))
	}
}

func (cl *funcCompiler) compileBinaryOp(dst ir.Slot, op bytecode.Op, e *ast.BinaryExpr) {
	xslot := cl.compileNestedTempExpr(e.X)
	yslot := cl.compileTempExpr(e.Y)
	cl.emit3(op, dst, xslot, yslot)
}

func (cl *funcCompiler) CompileSliceExpr(dst ir.Slot, isString bool, x, low, high ast.Expr) {
	switch {
	case low == nil && high != nil:
		strslot := cl.compileTempExpr(x)
		toslot := cl.compileTempExpr(high)
		cl.emit3(pickOp(isString, bytecode.OpStrSliceTo, bytecode.OpBytesSliceTo), dst, strslot, toslot)
	case low != nil && high == nil:
		strslot := cl.compileTempExpr(x)
		fromslot := cl.compileTempExpr(low)
		cl.emit3(pickOp(isString, bytecode.OpStrSliceFrom, bytecode.OpBytesSliceFrom), dst, strslot, fromslot)
	default:
		strslot := cl.compileTempExpr(x)
		fromslot := cl.compileTempExpr(low)
		toslot := cl.compileTempExpr(high)
		cl.emit4(pickOp(isString, bytecode.OpStrSlice, bytecode.OpBytesSlice), dst, strslot, fromslot, toslot)
	}
}

func (cl *funcCompiler) compileSliceExpr(dst ir.Slot, slice *ast.SliceExpr) {
	if slice.Slice3 {
		panic(cl.errorf(slice, "can't compile 3-index slicing"))
	}

	typ := cl.ctx.Types.TypeOf(slice.X)
	if !typeIsString(typ) && !typeIsByteSlice(typ) {
		panic(cl.errorf(slice.X, "can only slice strings and []byte"))
	}

	if cl.patternCompiler.CompileSliceExpr(dst, slice) {
		return
	}

	cl.CompileSliceExpr(dst, typeIsString(typ), slice.X, slice.Low, slice.High)
}

func (cl *funcCompiler) compileIndexExpr(dst ir.Slot, e *ast.IndexExpr) {
	typ := cl.ctx.Types.TypeOf(e.X)
	xslot := cl.compileTempExpr(e.X)
	indexslot := cl.compileTempExpr(e.Index)
	var op bytecode.Op
	switch {
	case typeIsString(typ):
		op = bytecode.OpStrIndex
	case typeIsSlice(typ):
		elemType := typ.Underlying().(*types.Slice).Elem()
		switch {
		case typeIsInt(elemType), typeIsFloat(elemType):
			op = bytecode.OpSliceIndexScalar64
		case typeIsBool(elemType), typeIsByte(elemType):
			op = bytecode.OpSliceIndexScalar8
		}
	}
	if op == bytecode.OpInvalid {
		panic(cl.errorf(e.X, "can't compile indexing of %s", typ))
	}
	cl.emit3(op, dst, xslot, indexslot)
}

func (cl *funcCompiler) compileSelectorExpr(dst ir.Slot, e *ast.SelectorExpr) {
	typ := cl.ctx.Types.TypeOf(e.X)
	key := qruntime.FuncKey{
		Name:      e.Sel.String(),
		Qualifier: typ.String(),
	}

	cl.compileCallArgs(nil, []ast.Expr{e.X}, nil)
	if cl.compileNativeCall(dst, key) {
		return
	}

	panic(cl.errorf(e, "can't compile %s field access", e.Sel))
}

func (cl *funcCompiler) compileCallExpr(dst ir.Slot, call *ast.CallExpr) {
	insideVariadic := cl.insideVariadic
	cl.compileCallExprImpl(dst, call)
	cl.insideVariadic = insideVariadic
}

func (cl *funcCompiler) compileIntConv(dst ir.Slot, call *ast.CallExpr) {
	x := call.Args[0]
	typ := cl.ctx.Types.TypeOf(x)
	if typeIsInt(typ) || typeIsByte(typ) {
		xslot := cl.compileTempExpr(x)
		cl.emit2(bytecode.OpMove, dst, xslot)
		return
	}
	panic(cl.errorf(call.Args[0], "can't convert %s to int", typ))
}

func (cl *funcCompiler) compileByteConv(dst ir.Slot, call *ast.CallExpr) {
	x := call.Args[0]
	typ := cl.ctx.Types.TypeOf(x)
	xslot := cl.compileTempExpr(x)
	switch {
	case typeIsByte(typ):
		cl.emit2(bytecode.OpMove, dst, xslot)
	case typeIsInt(typ):
		cl.emit2(bytecode.OpMove8, dst, xslot)
	default:
		panic(cl.errorf(call.Args[0], "can't convert %s to byte", typ))
	}
}

func (cl *funcCompiler) compileFloatConv(dst ir.Slot, call *ast.CallExpr) {
	x := call.Args[0]
	typ := cl.ctx.Types.TypeOf(x)
	xslot := cl.compileTempExpr(x)
	switch {
	case typeIsInt(typ):
		cl.emit2(bytecode.OpConvIntToFloat, dst, xslot)
	default:
		panic(cl.errorf(call.Args[0], "can't convert %s to float", typ))
	}
}

func (cl *funcCompiler) compileStringConv(dst ir.Slot, call *ast.CallExpr) {
	x := call.Args[0]
	typ := cl.ctx.Types.TypeOf(x)
	cl.compileCallArgs(nil, call.Args, nil)
	if typeIsByteSlice(typ) {
		key := qruntime.FuncKey{Qualifier: "builtin", Name: "bytesToString"}
		if !cl.compileNativeCall(dst, key) {
			panic(cl.errorf(call.Fun, "builtin.bytesToString native func is not registered"))
		}
	} else {
		panic(cl.errorf(call.Args[0], "can't convert %s to string", typ))
	}
}

func (cl *funcCompiler) compileCallExprImpl(dst ir.Slot, call *ast.CallExpr) {
	calledExpr := goutil.Unparen(call.Fun)

	if id, ok := calledExpr.(*ast.Ident); ok {
		_, isBuiltin := cl.ctx.Types.ObjectOf(id).(*types.Builtin)
		if isBuiltin {
			cl.compileBuiltinCall(dst, id, call)
			return
		}
		switch id.Name {
		case "int":
			cl.compileIntConv(dst, call)
			return
		case "byte":
			cl.compileByteConv(dst, call)
			return
		case "float64":
			cl.compileFloatConv(dst, call)
			return
		case "string":
			cl.compileStringConv(dst, call)
			return
		}
	}

	expr, fn := goutil.ResolveFunc(cl.ctx.Types, calledExpr)
	if fn == nil {
		panic(cl.errorf(call.Fun, "can't resolve the called function"))
	}

	// TODO: just use Func.FullName as a key?
	key := qruntime.FuncKey{Name: fn.Name()}
	sig := fn.Type().(*types.Signature)
	if sig.Recv() != nil {
		key.Qualifier = sig.Recv().Type().String()
	} else {
		key.Qualifier = fn.Pkg().Path()
	}

	normalArgs := call.Args
	var variadicArgs []ast.Expr
	if sig.Variadic() {
		if cl.insideVariadic {
			panic(cl.errorf(call.Fun, "can't call %s: nested variadic calls are not supported", key))
		}
		cl.insideVariadic = true
		variadic := sig.Params().Len() - 1
		normalArgs = call.Args[:variadic]
		variadicArgs = call.Args[variadic:]
	}

	if !sig.Variadic() && cl.inlineCall(dst, expr, normalArgs, key) {
		return
	}

	isMethod := expr != nil
	cl.compileCallArgs(expr, normalArgs, variadicArgs)

	if cl.compileNativeCall(dst, key) {
		if len(normalArgs) > qruntime.MaxNativeFuncArgs {
			panic(cl.errorf(call.Fun, "native funcs can't have more than %d args, got %d", qruntime.MaxNativeFuncArgs, len(normalArgs)))
		}
		return
	}
	if key == cl.fnKey && !isMethod && !sig.Variadic() && cl.compileRecurCall(dst) {
		return
	}
	if !isMethod && !sig.Variadic() && cl.compileCall(dst, key) {
		return
	}

	panic(cl.errorf(call.Fun, "can't compile a call to %s func", key))
}

func (cl *funcCompiler) compileBuiltinCall(dst ir.Slot, fn *ast.Ident, call *ast.CallExpr) {
	switch fn.Name {
	case `make`:
		cl.compileMakeCall(dst, call)
	case `append`:
		cl.compileAppendCall(dst, call)

	case `len`:
		srcslot := cl.compileTempExpr(call.Args[0])
		cl.emit2(bytecode.OpMove, dst, srcslot)
	case `cap`:
		srcslot := cl.compileTempExpr(call.Args[0])
		cl.emit2(bytecode.OpCap, dst, srcslot)

	case `println`:
		if len(call.Args) != 1 {
			panic(cl.errorf(call, "only 1-arg form of println() is supported"))
		}
		var funcName string
		argType := cl.ctx.Types.TypeOf(call.Args[0])
		switch {
		case typeIsByte(argType):
			funcName = "PrintByte"
		case typeIsInt(argType):
			funcName = "PrintInt"
		case typeIsFloat(argType):
			funcName = "PrintFloat"
		case typeIsString(argType):
			funcName = "PrintString"
		case typeIsBool(argType):
			funcName = "PrintBool"
		default:
			panic(cl.errorf(call.Args[0], "can't print %s type yet", argType.String()))
		}
		key := qruntime.FuncKey{Qualifier: "builtin", Name: funcName}
		cl.compileCallArgs(nil, call.Args, nil)
		if !cl.compileNativeCall(dst, key) {
			panic(cl.errorf(fn, "builtin.%s native func is not registered", funcName))
		}

	default:
		panic(cl.errorf(fn, "can't compile %s() builtin function call yet", fn))
	}
}

func (cl *funcCompiler) compileAppendCall(dst ir.Slot, call *ast.CallExpr) {
	sliceType := cl.ctx.Types.TypeOf(call).Underlying().(*types.Slice)
	if !typeIsScalar(sliceType.Elem()) {
		panic(cl.errorf(call.Args[0], "can't append() to a slice with non-scalar elems yet"))
	}
	if len(call.Args) != 2 {
		panic(cl.errorf(call.Args[0], "can only compile the 2-arguments form of append()"))
	}
	var funcName string
	switch cl.ctx.Sizes.Sizeof(sliceType.Elem()) {
	case 1:
		funcName = "append8"
	case 8:
		funcName = "append64"
	default:
		panic(cl.errorf(call.Args[0], "can't append to a slice with elem type %s", sliceType.Elem()))
	}
	cl.compileCallArgs(nil, call.Args, nil)
	key := qruntime.FuncKey{Qualifier: "builtin", Name: funcName}
	if !cl.compileNativeCall(dst, key) {
		panic(cl.errorf(call.Fun, "builtin.%s native func is not registered", funcName))
	}
}

func (cl *funcCompiler) compileMakeCall(dst ir.Slot, call *ast.CallExpr) {
	sliceType, ok := cl.ctx.Types.TypeOf(call).Underlying().(*types.Slice)
	if !ok {
		panic(cl.errorf(call.Args[0], "can't make() a non-slice type yet"))
	}
	var funcName string
	if !typeIsScalar(sliceType.Elem()) {
		panic(cl.errorf(call.Args[0], "can't make() a slice with non-scalar elems yet"))
	}
	funcName = "makeSlice"

	elemSize := cl.ctx.Sizes.Sizeof(sliceType.Elem())

	var args []ast.Expr
	elemSizeArg := &ast.BasicLit{
		Kind:     token.INT,
		Value:    strconv.FormatInt(elemSize, 10),
		ValuePos: call.Args[1].Pos(),
	}
	if len(call.Args) == 2 {
		args = []ast.Expr{elemSizeArg, call.Args[1], call.Args[1]}
	} else {
		args = []ast.Expr{elemSizeArg, call.Args[1], call.Args[2]}

	}
	cl.compileCallArgs(nil, args, nil)
	key := qruntime.FuncKey{Qualifier: "builtin", Name: funcName}
	if !cl.compileNativeCall(dst, key) {
		panic(cl.errorf(call.Fun, "builtin.%s native func is not registered", funcName))
	}
}

func (cl *funcCompiler) compileCallVariadicArgs(args []ast.Expr) {
	cl.emitOp(bytecode.OpVariadicReset)
	tempslot := cl.allocTemp()
	for _, arg := range args {
		cl.CompileExpr(tempslot, arg)
		argType := cl.ctx.Types.TypeOf(arg)
		switch {
		case typeIsBool(argType):
			cl.emit1(bytecode.OpPushVariadicBoolArg, tempslot)
		case typeIsScalar(argType):
			cl.emit1(bytecode.OpPushVariadicScalarArg, tempslot)
		case typeIsString(argType):
			cl.emit1(bytecode.OpPushVariadicStrArg, tempslot)
		case typeIsInterface(argType):
			cl.emit1(bytecode.OpPushVariadicInterfaceArg, tempslot)
		default:
			panic(cl.errorf(arg, "can't pass %s typed variadic arg", argType.String()))
		}
	}
}

func (cl *funcCompiler) checkTupleArg(args []ast.Expr) {
	if len(args) != 1 {
		return
	}
	// Check that it's not a f(g()) call, where g() returns
	// a multi-value result; we can't compile that yet.
	if call, ok := args[0].(*ast.CallExpr); ok {
		sig, ok := cl.ctx.Types.TypeOf(call.Fun).(*types.Signature)
		if ok && sig.Results() != nil && sig.Results().Len() > 1 {
			panic(cl.errorf(args[0], "can't pass tuple as a func argument"))
		}
	}
}

func (cl *funcCompiler) prependRecv(recv ast.Expr, args []ast.Expr) []ast.Expr {
	if recv == nil {
		return args
	}
	allArgs := make([]ast.Expr, 0, len(args)+1)
	allArgs = append(allArgs, recv)
	allArgs = append(allArgs, args...)
	return allArgs
}

func (cl *funcCompiler) compileCallArgs(recv ast.Expr, args []ast.Expr, variadic []ast.Expr) {
	cl.checkTupleArg(args)
	args = cl.prependRecv(recv, args)

	// If all arguments are really simple and can be evaluated without
	// clobbering arg slots, we can evaluate arguments directly to their
	// slots. Otherwise we'll need to use temporaries.
	needTemporaries := false
	for _, arg := range args {
		if !cl.isSimpleExpr(arg) {
			needTemporaries = true
			break
		}
	}
	if !needTemporaries {
		for _, arg := range variadic {
			if !cl.isSimpleExpr(arg) {
				needTemporaries = true
				break
			}
		}
	}

	if needTemporaries {
		tempSlots := make([]ir.Slot, 0, 8)
		for _, arg := range args {
			tempslot := cl.allocTemp()
			cl.CompileExpr(tempslot, arg)
			tempSlots = append(tempSlots, tempslot)
		}
		if variadic != nil {
			cl.compileCallVariadicArgs(variadic)
		}
		// Move temporaries to args.
		for i, slot := range tempSlots {
			argslot := ir.NewCallArgSlot(uint8(i))
			cl.emit2(bytecode.OpMove, argslot, slot)
		}
	} else {
		// Can move args directly to their slots.
		for i, arg := range args {
			argslot := ir.NewCallArgSlot(uint8(i))
			cl.CompileExpr(argslot, arg)
		}
		if variadic != nil {
			cl.compileCallVariadicArgs(variadic)
		}
	}
}

func (cl *funcCompiler) compileNativeCall(dst ir.Slot, key qruntime.FuncKey) bool {
	funcID, ok := cl.ctx.Env.NameToNativeFuncID[key]
	if !ok {
		return false
	}

	cl.emitCall(bytecode.OpCallNative, dst, int(funcID))
	return true
}

func (cl *funcCompiler) compileRecurCall(dst ir.Slot) bool {
	cl.hasCalls = true
	cl.emit1(bytecode.OpCallRecur, dst)
	return true
}

func (cl *funcCompiler) compileCall(dst ir.Slot, key qruntime.FuncKey) bool {
	funcID, ok := cl.ctx.Env.NameToFuncID[key]
	if !ok {
		return false
	}

	cl.emitCall(bytecode.OpCall, dst, int(funcID))
	return true
}

func (cl *funcCompiler) compileBasicLit(dst ir.Slot, lit *ast.BasicLit) {
	switch lit.Kind {
	case token.INT:
		v, err := strconv.ParseInt(lit.Value, 0, 64)
		if err != nil {
			panic(cl.errorf(lit, "invalid int value literal"))
		}
		cl.compileConstantValue(dst, lit, constant.MakeInt64(v))

	default:
		panic(cl.errorf(lit, "unexpected basic lit %v", lit.Kind))
	}

}

func (cl *funcCompiler) compileIdent(dst ir.Slot, ident *ast.Ident) {
	tv := cl.ctx.Types.Types[ident]
	cv := tv.Value
	if cv != nil {
		cl.compileConstantValue(dst, ident, cv)
		return
	}

	if i := cl.scope.Lookup(ident.String()); i != -1 {
		cl.emit2(bytecode.OpMove, dst, ir.NewTempSlot(uint8(i)))
		return
	}
	if p, ok := cl.params[ident.String()]; ok {
		cl.emit2(bytecode.OpMove, dst, p.i)
		return
	}

	panic(cl.errorf(ident, "can't compile a %s (type %s) variable read", ident.String(), tv.Type))
}

func (cl *funcCompiler) compileConstantValue(dst ir.Slot, source ast.Expr, cv constant.Value) {
	switch cv.Kind() {
	case constant.Bool:
		v := constant.BoolVal(cv)
		cl.emit(cl.moveBool(dst, v))

	case constant.Int:
		v, exact := constant.Int64Val(cv)
		if !exact {
			panic(cl.errorf(source, "non-exact int value"))
		}
		cl.emit(cl.moveInt(dst, int(v)))

	case constant.Float:
		v, exact := constant.Float64Val(cv)
		if !exact {
			panic(cl.errorf(source, "non-exact float value"))
		}
		cl.emit(cl.moveFloat(dst, v))

	case constant.String:
		v := constant.StringVal(cv)
		id := cl.internStrConstant(v)
		cl.emit(ir.Inst{
			Op:   bytecode.OpLoadStrConst,
			Arg0: dst.ToInstArg(),
			Arg1: ir.InstArg(id),
		})

	case constant.Complex:
		panic(cl.errorf(source, "can't compile complex number constants yet"))

	default:
		panic(cl.errorf(source, "unexpected constant %v", cv))
	}
}

func (cl *funcCompiler) compileOr(dst ir.Slot, e *ast.BinaryExpr) {
	labelEnd := cl.newLabel()
	cl.compileNestedRootExpr(dst, e.X)
	cl.emitCondJump(dst, bytecode.OpJumpNotZero, labelEnd)
	cl.CompileExpr(dst, e.Y)
	cl.bindLabel(labelEnd)
}

func (cl *funcCompiler) compileAnd(dst ir.Slot, e *ast.BinaryExpr) {
	labelEnd := cl.newLabel()
	cl.compileNestedRootExpr(dst, e.X)
	cl.emitCondJump(dst, bytecode.OpJumpZero, labelEnd)
	cl.CompileExpr(dst, e.Y)
	cl.bindLabel(labelEnd)
}
