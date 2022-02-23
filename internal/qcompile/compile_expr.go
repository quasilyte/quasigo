package qcompile

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"

	"github.com/quasilyte/quasigo/internal/goutil"
	"github.com/quasilyte/quasigo/internal/qruntime"

	"github.com/quasilyte/quasigo/internal/bytecode"
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
	cl.CompileExpr(tmp, e)
	return tmp
}

func (cl *compiler) compileRootTempExpr(e ast.Expr) int {
	slot := cl.compileTempExpr(e)
	cl.freeTmp()
	return slot
}

func (cl *compiler) compileRootExpr(dst int, e ast.Expr) {
	cl.CompileExpr(dst, e)
	cl.freeTmp()
}

func (cl *compiler) CompileExpr(dst int, e ast.Expr) {
	cv := cl.ctx.Types.Types[e].Value
	if cv != nil {
		cl.compileConstantValue(dst, e, cv)
		return
	}

	switch e := e.(type) {
	case *ast.ParenExpr:
		cl.CompileExpr(dst, e.X)

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

func (cl *compiler) compileUnaryExpr(dst int, e *ast.UnaryExpr) {
	switch e.Op {
	case token.NOT:
		cl.compileUnaryOp(dst, bytecode.OpNot, e.X)

	case token.SUB:
		cl.compileUnaryOp(dst, bytecode.OpIntNeg, e.X)

	default:
		panic(cl.errorf(e, "can't compile unary %s yet", e.Op))
	}
}

func (cl *compiler) compileUnaryOp(dst int, op bytecode.Op, arg ast.Expr) {
	xslot := cl.compileTempExpr(arg)
	cl.emit2(op, dst, xslot)
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
		cl.compileScalarBinaryOp(dst, e, bytecode.OpIntGt, typ)
	case token.GEQ:
		cl.compileScalarBinaryOp(dst, e, bytecode.OpIntGtEq, typ)
	case token.LSS:
		cl.compileScalarBinaryOp(dst, e, bytecode.OpIntLt, typ)
	case token.LEQ:
		cl.compileScalarBinaryOp(dst, e, bytecode.OpIntLtEq, typ)

	case token.ADD:
		switch {
		case typeIsString(typ):
			cl.compileBinaryOp(dst, bytecode.OpConcat, e)
		case typeIsInt(typ), typeIsByte(typ):
			cl.compileBinaryOp(dst, bytecode.OpIntAdd, e)
		default:
			panic(cl.errorf(e, "+ is not implemented for %s bytecode.Operands", typ))
		}

	case token.SUB:
		cl.compileScalarBinaryOp(dst, e, bytecode.OpIntSub, typ)
	case token.XOR:
		cl.compileScalarBinaryOp(dst, e, bytecode.OpIntXor, typ)
	case token.MUL:
		cl.compileScalarBinaryOp(dst, e, bytecode.OpIntMul, typ)
	case token.QUO:
		cl.compileIntBinaryOp(dst, e, bytecode.OpIntDiv, typ)

	default:
		panic(cl.errorf(e, "can't compile binary %s yet", e.Op))
	}
}

func (cl *compiler) compileScalarBinaryOp(dst int, e *ast.BinaryExpr, op bytecode.Op, typ types.Type) {
	if typeIsInt(typ) || typeIsByte(typ) {
		cl.compileBinaryOp(dst, op, e)
	} else {
		panic(cl.errorf(e, "%s is not implemented for %s bytecode.Operands", e.Op, typ))
	}
}

func (cl *compiler) compileIntBinaryOp(dst int, e *ast.BinaryExpr, op bytecode.Op, typ types.Type) {
	if typeIsInt(typ) {
		cl.compileBinaryOp(dst, op, e)
	} else {
		panic(cl.errorf(e, "%s is not implemented for %s bytecode.Operands", e.Op, typ))
	}
}

func (cl *compiler) compileBinaryOp(dst int, op bytecode.Op, e *ast.BinaryExpr) {
	xslot := cl.compileTempExpr(e.X)
	yslot := cl.compileTempExpr(e.Y)
	cl.emit3(op, dst, xslot, yslot)
}

func (cl *compiler) CompileSliceExpr(dst int, x, low, high ast.Expr) {
	switch {
	case low == nil && high != nil:
		strslot := cl.compileTempExpr(x)
		toslot := cl.compileTempExpr(high)
		cl.emit3(bytecode.OpStrSliceTo, dst, strslot, toslot)
	case low != nil && high == nil:
		strslot := cl.compileTempExpr(x)
		fromslot := cl.compileTempExpr(low)
		cl.emit3(bytecode.OpStrSliceFrom, dst, strslot, fromslot)
	default:
		strslot := cl.compileTempExpr(x)
		fromslot := cl.compileTempExpr(low)
		toslot := cl.compileTempExpr(high)
		cl.emit4(bytecode.OpStrSlice, dst, strslot, fromslot, toslot)
	}
}

func (cl *compiler) compileSliceExpr(dst int, slice *ast.SliceExpr) {
	if slice.Slice3 {
		panic(cl.errorf(slice, "can't compile 3-index slicing"))
	}

	if !typeIsString(cl.ctx.Types.TypeOf(slice.X)) {
		panic(cl.errorf(slice.X, "can't compile slicing of something that is not a string"))
	}

	if cl.patternCompiler.CompileSliceExpr(dst, slice) {
		return
	}
	cl.CompileSliceExpr(dst, slice.X, slice.Low, slice.High)
}

func (cl *compiler) compileIndexExpr(dst int, e *ast.IndexExpr) {
	if !typeIsString(cl.ctx.Types.TypeOf(e.X)) {
		panic(cl.errorf(e.X, "can't compile indexing of something that is not a string"))
	}
	strslot := cl.compileTempExpr(e.X)
	indexslot := cl.compileTempExpr(e.Index)
	cl.emit3(bytecode.OpStrIndex, dst, strslot, indexslot)
}

func (cl *compiler) compileSelectorExpr(dst int, e *ast.SelectorExpr) {
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

func (cl *compiler) compileCallExpr(dst int, call *ast.CallExpr) {
	insideVariadic := cl.insideVariadic
	cl.compileCallExprImpl(dst, call)
	cl.insideVariadic = insideVariadic
}

func (cl *compiler) compileIntConv(dst int, call *ast.CallExpr) {
	x := call.Args[0]
	typ := cl.ctx.Types.TypeOf(x)
	if typeIsInt(typ) || typeIsByte(typ) {
		xslot := cl.compileTempExpr(x)
		cl.emit2(cl.opMoveByType(x, typ), dst, xslot)
		return
	}
	panic(cl.errorf(call.Args[0], "can't convert %s to int", typ))
}

func (cl *compiler) compileCallExprImpl(dst int, call *ast.CallExpr) {
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
		}
	}

	expr, fn := goutil.ResolveFunc(cl.ctx.Types, calledExpr)
	if fn == nil {
		panic(cl.errorf(call.Fun, "can't resolve the called function"))
	}

	// // TODO: just use Func.FullName as a key?
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

func (cl *compiler) compileBuiltinCall(dst int, fn *ast.Ident, call *ast.CallExpr) {
	switch fn.Name {
	case `len`:
		s := call.Args[0]
		srcslot := cl.compileTempExpr(s)
		if !typeIsString(cl.ctx.Types.TypeOf(s)) {
			panic(cl.errorf(s, "can't compile len() with non-string argument yet"))
		}
		cl.emit2(bytecode.OpStrLen, dst, srcslot)

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

func (cl *compiler) compileCallVariadicArgs(args []ast.Expr) {
	cl.emitOp(bytecode.OpVariadicReset)
	tmpslot := cl.allocTmp()
	for _, arg := range args {
		cl.CompileExpr(tmpslot, arg)
		argType := cl.ctx.Types.TypeOf(arg)
		switch {
		case typeIsBool(argType):
			cl.emit1(bytecode.OpPushVariadicBoolArg, tmpslot)
		case typeIsScalar(argType):
			cl.emit1(bytecode.OpPushVariadicScalarArg, tmpslot)
		case typeIsString(argType):
			cl.emit1(bytecode.OpPushVariadicStrArg, tmpslot)
		case typeIsInterface(argType):
			cl.emit1(bytecode.OpPushVariadicInterfaceArg, tmpslot)
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
			cl.CompileExpr(tmpslot, arg)
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
			cl.emit2(moveOp, argslot, slot)
		}
	} else {
		// Can move args directly to their slots.
		for i, arg := range args {
			argslot := -(i + 1)
			cl.CompileExpr(argslot, arg)
		}
		if variadic != nil {
			cl.compileCallVariadicArgs(variadic)
		}
	}
}

func (cl *compiler) compileNativeCall(dst int, key qruntime.FuncKey) bool {
	funcID, ok := cl.ctx.Env.NameToNativeFuncID[key]
	if !ok {
		return false
	}

	cl.emitCall(bytecode.OpCallNative, dst, int(funcID))
	return true
}

func (cl *compiler) compileRecurCall(dst int) bool {
	cl.emit1(bytecode.OpCallRecur, dst)
	return true
}

func (cl *compiler) compileCall(dst int, key qruntime.FuncKey) bool {
	funcID, ok := cl.ctx.Env.NameToFuncID[key]
	if !ok {
		return false
	}

	cl.emitCall(bytecode.OpCall, dst, int(funcID))
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
		cl.emit2(cl.opMoveByType(ident, p.v.Type()), dst, p.i)
		return
	}
	if l, ok := cl.locals[ident.String()]; ok {
		cl.emit2(cl.opMoveByType(ident, l.v.Type()), dst, l.i)
		return
	}

	panic(cl.errorf(ident, "can't compile a %s (type %s) variable read", ident.String(), tv.Type))
}

func (cl *compiler) compileConstantValue(dst int, source ast.Expr, cv constant.Value) {
	switch cv.Kind() {
	case constant.Bool:
		v := constant.BoolVal(cv)
		id := cl.internBoolConstant(v)
		cl.emit2(bytecode.OpLoadScalarConst, dst, id)

	case constant.String:
		v := constant.StringVal(cv)
		id := cl.internStrConstant(v)
		cl.emit2(bytecode.OpLoadStrConst, dst, id)

	case constant.Int:
		v, exact := constant.Int64Val(cv)
		if !exact {
			panic(cl.errorf(source, "non-exact int value"))
		}
		id := cl.internIntConstant(int(v))
		cl.emit2(bytecode.OpLoadScalarConst, dst, id)

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
	cl.CompileExpr(dst, e.X)
	cl.emitCondJump(dst, bytecode.OpJumpNotZero, labelEnd)
	cl.CompileExpr(dst, e.Y)
	cl.bindLabel(labelEnd)
}

func (cl *compiler) compileAnd(dst int, e *ast.BinaryExpr) {
	labelEnd := cl.newLabel()
	cl.CompileExpr(dst, e.X)
	cl.emitCondJump(dst, bytecode.OpJumpZero, labelEnd)
	cl.CompileExpr(dst, e.Y)
	cl.bindLabel(labelEnd)
}
