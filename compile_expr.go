package quasigo

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"

	"github.com/quasilyte/go-ruleguard/ruleguard/goutil"
	"golang.org/x/tools/go/ast/astutil"
)

func (cl *compiler) compileTempExpr(e ast.Expr) int {
	if v, ok := e.(*ast.Ident); ok {
		if p, ok := cl.params[v.Name]; ok {
			return p.i
		}
		if l, ok := cl.locals[v.Name]; ok {
			return l.i
		}
	}
	tmp := cl.allocTmp()
	cl.compileExpr(tmp, e)
	return tmp
}

func (cl *compiler) compileRootTempExpr(e ast.Expr) int {
	slot := cl.compileTempExpr(e)
	cl.freeTmp()
	return slot
}

func (cl *compiler) compileRootExpr(dst int, e ast.Expr) {
	cl.compileExpr(dst, e)
	cl.freeTmp()
}

func (cl *compiler) compileExpr(dst int, e ast.Expr) {
	cv := cl.ctx.Types.Types[e].Value
	if cv != nil {
		cl.compileConstantValue(dst, e, cv)
		return
	}

	switch e := e.(type) {
	case *ast.ParenExpr:
		cl.compileExpr(dst, e.X)

	case *ast.Ident:
		cl.compileIdent(dst, e)

	case *ast.UnaryExpr:
		cl.compileUnaryExpr(dst, e)

	case *ast.BinaryExpr:
		cl.compileBinaryExpr(dst, e)

	case *ast.SliceExpr:
		cl.compileSliceExpr(dst, e)

	case *ast.SelectorExpr:
		cl.compileSelectorExpr(dst, e)

	case *ast.CallExpr:
		cl.compileCallExpr(dst, e)

	default:
		panic(cl.errorf(e, "can't compile %T yet", e))
	}
}

func (cl *compiler) compileUnaryExpr(dst int, e *ast.UnaryExpr) {
	switch e.Op {
	case token.NOT:
		cl.compileUnaryOp(dst, opNot, e.X)
	default:
		panic(cl.errorf(e, "can't compile unary %s yet", e.Op))
	}
}

func (cl *compiler) compileUnaryOp(dst int, op opcode, arg ast.Expr) {
	xslot := cl.compileTempExpr(arg)
	cl.emit8x2(op, dst, xslot)
}

func (cl *compiler) compileBinaryExpr(dst int, e *ast.BinaryExpr) {
	typ := cl.ctx.Types.TypeOf(e.X)

	switch e.Op {
	case token.LOR:
		cl.compileOr(dst, e)
	case token.LAND:
		cl.compileAnd(dst, e)

	case token.NEQ:
		switch {
		case identName(e.X) == "nil":
			cl.compileUnaryOp(dst, pickOp(typeIsInterface(cl.ctx.Types.TypeOf(e.Y)), opIsNotNilInterface, opIsNotNil), e.Y)
		case identName(e.Y) == "nil":
			cl.compileUnaryOp(dst, pickOp(typeIsInterface(typ), opIsNotNilInterface, opIsNotNil), e.X)

		case typeIsString(typ):
			cl.compileBinaryOp(dst, opStrNotEq, e)
		case typeIsInt(typ):
			cl.compileBinaryOp(dst, opIntNotEq, e)
		default:
			panic(cl.errorf(e, "!= is not implemented for %s operands", typ))
		}
	case token.EQL:
		switch {
		case identName(e.X) == "nil":
			cl.compileUnaryOp(dst, pickOp(typeIsInterface(cl.ctx.Types.TypeOf(e.Y)), opIsNilInterface, opIsNil), e.Y)
		case identName(e.Y) == "nil":
			cl.compileUnaryOp(dst, pickOp(typeIsInterface(typ), opIsNilInterface, opIsNil), e.X)

		case typeIsString(cl.ctx.Types.TypeOf(e.X)):
			cl.compileBinaryOp(dst, opStrEq, e)
		case typeIsInt(cl.ctx.Types.TypeOf(e.X)):
			cl.compileBinaryOp(dst, opIntEq, e)
		default:
			panic(cl.errorf(e, "== is not implemented for %s operands", typ))
		}

	case token.GTR:
		cl.compileIntBinaryOp(dst, e, opIntGt, typ)
	case token.GEQ:
		cl.compileIntBinaryOp(dst, e, opIntGtEq, typ)
	case token.LSS:
		cl.compileIntBinaryOp(dst, e, opIntLt, typ)
	case token.LEQ:
		cl.compileIntBinaryOp(dst, e, opIntLtEq, typ)

	case token.ADD:
		switch {
		case typeIsString(typ):
			cl.compileBinaryOp(dst, opConcat, e)
		case typeIsInt(typ):
			cl.compileBinaryOp(dst, opIntAdd, e)
		default:
			panic(cl.errorf(e, "+ is not implemented for %s operands", typ))
		}

	case token.SUB:
		cl.compileIntBinaryOp(dst, e, opIntSub, typ)
	case token.MUL:
		cl.compileIntBinaryOp(dst, e, opIntMul, typ)
	case token.QUO:
		cl.compileIntBinaryOp(dst, e, opIntDiv, typ)

	default:
		panic(cl.errorf(e, "can't compile binary %s yet", e.Op))
	}
}

func (cl *compiler) compileIntBinaryOp(dst int, e *ast.BinaryExpr, op opcode, typ types.Type) {
	switch {
	case typeIsInt(typ):
		cl.compileBinaryOp(dst, op, e)
	default:
		panic(cl.errorf(e, "%s is not implemented for %s operands", e.Op, typ))
	}
}

func (cl *compiler) compileBinaryOp(dst int, op opcode, e *ast.BinaryExpr) {
	xslot := cl.compileTempExpr(e.X)
	yslot := cl.compileTempExpr(e.Y)
	cl.emit8x3(op, dst, xslot, yslot)
}

func (cl *compiler) compileSliceExpr(dst int, slice *ast.SliceExpr) {
	if slice.Slice3 {
		panic(cl.errorf(slice, "can't compile 3-index slicing"))
	}

	// No need to do slicing, its no-op `s[:]`.
	if slice.Low == nil && slice.High == nil {
		cl.compileExpr(dst, slice.X)
		return
	}

	if !typeIsString(cl.ctx.Types.TypeOf(slice.X)) {
		panic(cl.errorf(slice.X, "can't compile slicing of something that is not a string"))
	}

	switch {
	case slice.Low == nil && slice.High != nil:
		strslot := cl.compileTempExpr(slice.X)
		toslot := cl.compileTempExpr(slice.High)
		cl.emit8x3(opStrSliceTo, dst, strslot, toslot)
	case slice.Low != nil && slice.High == nil:
		strslot := cl.compileTempExpr(slice.X)
		fromslot := cl.compileTempExpr(slice.Low)
		cl.emit8x3(opStrSliceFrom, dst, strslot, fromslot)
	default:
		strslot := cl.compileTempExpr(slice.X)
		fromslot := cl.compileTempExpr(slice.Low)
		toslot := cl.compileTempExpr(slice.High)
		cl.emit8x4(opStrSlice, dst, strslot, fromslot, toslot)
	}
}

func (cl *compiler) compileSelectorExpr(dst int, e *ast.SelectorExpr) {
	typ := cl.ctx.Types.TypeOf(e.X)
	key := funcKey{
		name:      e.Sel.String(),
		qualifier: typ.String(),
	}

	cl.compileCallArgs(nil, []ast.Expr{e.X}, nil)
	if cl.compileNativeCall(dst, key) {
		return
	}

	panic(cl.errorf(e, "can't compile %s field access", e.Sel))
}

func (cl *compiler) compileCallExpr(dst int, call *ast.CallExpr) {
	insideVariadic := cl.insideVariadic
	cl.compileCallExprImpl(dst, call)
	cl.insideVariadic = insideVariadic
}

func (cl *compiler) compileCallExprImpl(dst int, call *ast.CallExpr) {
	if id, ok := astutil.Unparen(call.Fun).(*ast.Ident); ok {
		_, isBuiltin := cl.ctx.Types.ObjectOf(id).(*types.Builtin)
		if isBuiltin {
			cl.compileBuiltinCall(dst, id, call)
			return
		}
	}

	expr, fn := goutil.ResolveFunc(cl.ctx.Types, call.Fun)
	if fn == nil {
		panic(cl.errorf(call.Fun, "can't resolve the called function"))
	}

	// // TODO: just use Func.FullName as a key?
	key := funcKey{name: fn.Name()}
	sig := fn.Type().(*types.Signature)
	if sig.Recv() != nil {
		key.qualifier = sig.Recv().Type().String()
	} else {
		key.qualifier = fn.Pkg().Path()
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

	isMethod := expr != nil
	cl.compileCallArgs(expr, normalArgs, variadicArgs)

	if cl.compileNativeCall(dst, key) {
		if len(normalArgs) > maxNativeFuncArgs {
			panic(cl.errorf(call.Fun, "native funcs can't have more than %d args, got %d", maxNativeFuncArgs, len(normalArgs)))
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

func (cl *compiler) compileBuiltinCall(dst int, fn *ast.Ident, call *ast.CallExpr) {
	switch fn.Name {
	case `len`:
		s := call.Args[0]
		srcslot := cl.compileTempExpr(s)
		if !typeIsString(cl.ctx.Types.TypeOf(s)) {
			panic(cl.errorf(s, "can't compile len() with non-string argument yet"))
		}
		cl.emit8x2(opStrLen, dst, srcslot)

	case `println`:
		if len(call.Args) != 1 {
			panic(cl.errorf(call, "only 1-arg form of println() is supported"))
		}
		var funcName string
		argType := cl.ctx.Types.TypeOf(call.Args[0])
		switch {
		case typeIsInt(argType):
			funcName = "PrintInt"
		case typeIsString(argType):
			funcName = "PrintString"
		case typeIsBool(argType):
			funcName = "PrintBool"
		default:
			panic(cl.errorf(call.Args[0], "can't print %s type yet", argType.String()))
		}
		key := funcKey{qualifier: "builtin", name: funcName}
		cl.compileCallArgs(nil, call.Args, nil)
		if !cl.compileNativeCall(dst, key) {
			panic(cl.errorf(fn, "builtin.%s native func is not registered", funcName))
		}

	default:
		panic(cl.errorf(fn, "can't compile %s() builtin function call yet", fn))
	}
}

func (cl *compiler) compileCallVariadicArgs(args []ast.Expr) {
	cl.emit(opVariadicReset)
	tmpslot := cl.allocTmp()
	for _, arg := range args {
		cl.compileExpr(tmpslot, arg)
		argType := cl.ctx.Types.TypeOf(arg)
		switch {
		case typeIsBool(argType):
			cl.emit8(opPushVariadicBoolArg, tmpslot)
		case typeIsScalar(argType):
			cl.emit8(opPushVariadicScalarArg, tmpslot)
		case typeIsString(argType):
			cl.emit8(opPushVariadicStrArg, tmpslot)
		case typeIsInterface(argType):
			cl.emit8(opPushVariadicInterfaceArg, tmpslot)
		default:
			panic(cl.errorf(arg, "can't pass %s typed variadic arg", argType.String()))
		}
	}
}

func (cl *compiler) compileCallArgs(recv ast.Expr, args []ast.Expr, variadic []ast.Expr) {
	if len(args) == 1 {
		// Check that it's not a f(g()) call, where g() returns
		// a multi-value result; we can't compile that yet.
		if call, ok := args[0].(*ast.CallExpr); ok {
			results := cl.ctx.Types.TypeOf(call.Fun).(*types.Signature).Results()
			if results != nil && results.Len() > 1 {
				panic(cl.errorf(args[0], "can't pass tuple as a func argument"))
			}
		}
	}

	if recv != nil {
		allArgs := make([]ast.Expr, 0, len(args)+1)
		allArgs = append(allArgs, recv)
		allArgs = append(allArgs, args...)
		args = allArgs
	}

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
		tempSlots := make([]int, 0, 8)
		for _, arg := range args {
			tmpslot := cl.allocTmp()
			cl.compileExpr(tmpslot, arg)
			tempSlots = append(tempSlots, tmpslot)
		}
		if variadic != nil {
			cl.compileCallVariadicArgs(variadic)
		}
		// Move temporaries to args.
		for i, slot := range tempSlots {
			argslot := -(i + 1)
			arg := args[i]
			moveOp := cl.opMoveByType(arg, cl.ctx.Types.TypeOf(arg))
			cl.emit8x2(moveOp, argslot, slot)
		}
	} else {
		// Can move args directly to their slots.
		for i, arg := range args {
			argslot := -(i + 1)
			cl.compileExpr(argslot, arg)
		}
		if variadic != nil {
			cl.compileCallVariadicArgs(variadic)
		}
	}
}

func (cl *compiler) compileNativeCall(dst int, key funcKey) bool {
	funcID, ok := cl.ctx.Env.nameToNativeFuncID[key]
	if !ok {
		return false
	}

	// if variadic != 0 {
	// 	for _, arg := range variadicArgs {
	// 		cl.compileExpr(arg)
	// 		// int-typed values should appear in the interface{}-typed
	// 		// objects slice, so we get all variadic args placed in one place.
	// 		if typeIsInt(cl.ctx.Types.TypeOf(arg)) {
	// 			cl.emit(opConvIntToIface)
	// 		}
	// 	}
	// 	if len(variadicArgs) > 255 {
	// 		panic(cl.errorf(funcExpr, "too many variadic args"))
	// 	}
	// 	// Even if len(variadicArgs) is 0, we still need to overwrite
	// 	// the old variadicLen value, so the variadic func is not confused
	// 	// by some unrelated value.
	// 	cl.emit8(opSetVariadicLen, len(variadicArgs))
	// }

	cl.emitCall(opCallNative, dst, int(funcID))
	return true
}

func (cl *compiler) compileRecurCall(dst int) bool {
	cl.emit8(opCallRecur, dst)
	return true
}

func (cl *compiler) compileCall(dst int, key funcKey) bool {
	funcID, ok := cl.ctx.Env.nameToFuncID[key]
	if !ok {
		return false
	}

	cl.emitCall(opCall, dst, int(funcID))
	return true
}

func (cl *compiler) compileIdent(dst int, ident *ast.Ident) {
	tv := cl.ctx.Types.Types[ident]
	cv := tv.Value
	if cv != nil {
		cl.compileConstantValue(dst, ident, cv)
		return
	}

	if p, ok := cl.params[ident.String()]; ok {
		cl.emit8x2(cl.opMoveByType(ident, p.v.Type()), dst, p.i)
		return
	}
	if l, ok := cl.locals[ident.String()]; ok {
		cl.emit8x2(cl.opMoveByType(ident, l.v.Type()), dst, l.i)
		return
	}

	panic(cl.errorf(ident, "can't compile a %s (type %s) variable read", ident.String(), tv.Type))
}

func (cl *compiler) compileConstantValue(dst int, source ast.Expr, cv constant.Value) {
	switch cv.Kind() {
	case constant.Bool:
		v := constant.BoolVal(cv)
		id := cl.internBoolConstant(v)
		cl.emit8x2(opLoadScalarConst, dst, id)

	case constant.String:
		v := constant.StringVal(cv)
		id := cl.internStrConstant(v)
		cl.emit8x2(opLoadStrConst, dst, id)

	case constant.Int:
		v, exact := constant.Int64Val(cv)
		if !exact {
			panic(cl.errorf(source, "non-exact int value"))
		}
		id := cl.internIntConstant(int(v))
		cl.emit8x2(opLoadScalarConst, dst, id)

	case constant.Complex:
		panic(cl.errorf(source, "can't compile complex number constants yet"))

	case constant.Float:
		panic(cl.errorf(source, "can't compile float constants yet"))

	default:
		panic(cl.errorf(source, "unexpected constant %v", cv))
	}
}

func (cl *compiler) compileOr(dst int, e *ast.BinaryExpr) {
	labelEnd := cl.newLabel()
	cl.compileExpr(dst, e.X)
	cl.emitCondJump(dst, opJumpTrue, labelEnd)
	cl.compileExpr(dst, e.Y)
	cl.bindLabel(labelEnd)
}

func (cl *compiler) compileAnd(dst int, e *ast.BinaryExpr) {
	labelEnd := cl.newLabel()
	cl.compileExpr(dst, e.X)
	cl.emitCondJump(dst, opJumpFalse, labelEnd)
	cl.compileExpr(dst, e.Y)
	cl.bindLabel(labelEnd)
}
