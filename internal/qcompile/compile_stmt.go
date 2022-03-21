package qcompile

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"

	"github.com/quasilyte/quasigo/internal/bytecode"
)

func (cl *compiler) compileStmt(stmt ast.Stmt) {
	switch stmt := stmt.(type) {
	case *ast.BlockStmt:
		cl.compileStmtList(stmt.List)

	case *ast.SwitchStmt:
		cl.compileSwitchStmt(stmt)

	case *ast.ReturnStmt:
		cl.compileReturnStmt(stmt)

	case *ast.IfStmt:
		cl.compileIfStmt(stmt)

	case *ast.AssignStmt:
		cl.compileAssignStmt(stmt)

	case *ast.IncDecStmt:
		cl.compileIncDecStmt(stmt)

	case *ast.BranchStmt:
		cl.compileBranchStmt(stmt)

	case *ast.ForStmt:
		cl.compileForStmt(stmt)

	case *ast.ExprStmt:
		cl.compileExprStmt(stmt)

	default:
		panic(cl.errorf(stmt, "can't compile %T yet", stmt))
	}
}

func (cl *compiler) compileStmtList(list []ast.Stmt) {
	for i := range list {
		cl.compileStmt(list[i])
	}
}

func (cl *compiler) compileReturnStmt(ret *ast.ReturnStmt) {
	if cl.retType == voidType {
		cl.emitOp(bytecode.OpReturnVoid)
		return
	}

	if ret.Results == nil {
		panic(cl.errorf(ret, "'naked' return statements are not allowed"))
	}

	cv := cl.ctx.Types.Types[ret.Results[0]].Value
	if cv != nil {
		// Return of a constant value.
		switch cv.Kind() {
		case constant.Bool:
			cl.emitOp(pickOp(constant.BoolVal(cv), bytecode.OpReturnOne, bytecode.OpReturnZero))
			return
		case constant.Int:
			v, exact := constant.Int64Val(cv)
			if exact && v == 0 {
				cl.emitOp(bytecode.OpReturnZero)
				return
			}
			if exact && v == 1 {
				cl.emitOp(bytecode.OpReturnOne)
				return
			}
		}
	}

	typ := cl.ctx.Types.TypeOf(ret.Results[0])
	var op bytecode.Op
	switch {
	case typeIsScalar(typ):
		op = bytecode.OpReturnScalar
	case typeIsString(typ):
		op = bytecode.OpReturnStr
	default:
		op = bytecode.OpReturn
	}
	slot := cl.compileRootTempExpr(ret.Results[0])
	cl.emit1(op, slot)
}

func (cl *compiler) compileIfStmt(stmt *ast.IfStmt) {
	if stmt.Else == nil {
		labelEnd := cl.newLabel()
		condslot := cl.compileRootTempExpr(stmt.Cond)
		cl.emitCondJump(condslot, bytecode.OpJumpZero, labelEnd)
		cl.compileStmt(stmt.Body)
		cl.bindLabel(labelEnd)
		return
	}

	labelEnd := cl.newLabel()
	labelElse := cl.newLabel()
	condslot := cl.compileRootTempExpr(stmt.Cond)
	cl.emitCondJump(condslot, bytecode.OpJumpZero, labelElse)
	cl.compileStmt(stmt.Body)
	if !cl.isUncondJump(cl.lastOp()) {
		cl.emitJump(labelEnd)
	}
	cl.bindLabel(labelElse)
	cl.compileStmt(stmt.Else)
	cl.bindLabel(labelEnd)
}

func (cl *compiler) compileAssignOp(op bytecode.Op, assign *ast.AssignStmt) {
	lhs := assign.Lhs[0].(*ast.Ident)
	rhs := assign.Rhs[0]
	dstslot := cl.getNamedSlot(lhs, lhs.Name)
	rhsslot := cl.compileRootTempExpr(rhs)
	cl.emit3(op, dstslot, dstslot, rhsslot)
}

func (cl *compiler) compileAssignIndex(e *ast.IndexExpr, assign *ast.AssignStmt) {
	if len(assign.Lhs) != 1 {
		panic(cl.errorf(assign, "only single lhs operand is allowed in index assignments"))
	}
	if assign.Tok != token.ASSIGN {
		panic(cl.errorf(assign, "only = index assignments are supported"))
	}
	typ := cl.ctx.Types.TypeOf(e.X)
	if !typeIsSlice(typ) {
		panic(cl.errorf(assign, "only slices support index assignments"))
	}
	elemType := typ.Underlying().(*types.Slice).Elem()
	var op bytecode.Op
	switch {
	case typeIsInt(elemType):
		op = bytecode.OpSliceSetScalar64
	case typeIsBool(elemType), typeIsByte(elemType):
		op = bytecode.OpSliceSetScalar8
	}
	cl.beginTempBlock()
	valueslot := cl.compileTempExpr(assign.Rhs[0])
	xslot := cl.compileTempExpr(e.X)
	indexslot := cl.compileTempExpr(e.Index)
	cl.emit3(op, xslot, indexslot, valueslot)
	cl.endTempBlock()
}

func (cl *compiler) compileAssignStmt(assign *ast.AssignStmt) {
	if len(assign.Rhs) != 1 {
		panic(cl.errorf(assign, "only single right operand is allowed in assignments"))
	}
	if indexing, ok := assign.Lhs[0].(*ast.IndexExpr); ok {
		cl.compileAssignIndex(indexing, assign)
		return
	}
	for _, lhs := range assign.Lhs {
		_, ok := lhs.(*ast.Ident)
		if !ok {
			panic(cl.errorf(lhs, "can assign only to simple variables"))
		}
	}
	if len(assign.Lhs) > 2 {
		panic(cl.errorf(assign, "at most 2 value results are supported"))
	}

	if len(assign.Lhs) == 1 {
		op := bytecode.OpInvalid
		typ := cl.ctx.Types.TypeOf(assign.Rhs[0])
		switch assign.Tok {
		case token.MUL_ASSIGN:
			op = pickOp(typeIsByte(typ), bytecode.OpIntMul8, bytecode.OpIntMul64)
		case token.XOR_ASSIGN:
			op = bytecode.OpIntXor
		case token.ADD_ASSIGN:
			switch {
			case typeIsString(typ):
				op = bytecode.OpConcat
			case typeIsByte(typ):
				op = bytecode.OpIntAdd8
			default:
				op = bytecode.OpIntAdd64
			}
		case token.SUB_ASSIGN:
			op = pickOp(typeIsByte(typ), bytecode.OpIntSub8, bytecode.OpIntSub64)
		}
		if op != bytecode.OpInvalid {
			cl.compileAssignOp(op, assign)
			return
		}
	}

	switch assign.Tok {
	case token.ASSIGN, token.DEFINE:
		// OK.
	default:
		panic(cl.errorf(assign, "can't compile %s assign op", assign.Tok))
	}

	dst1 := assign.Lhs[0].(*ast.Ident)
	rhs := assign.Rhs[0]
	lhs1slot := cl.getNamedSlot(dst1, dst1.Name)
	cl.compileRootExpr(lhs1slot, rhs)
	if len(assign.Lhs) == 2 {
		dst2 := assign.Lhs[1].(*ast.Ident)
		lhs2slot := cl.getNamedSlot(dst2, dst2.Name)
		cl.emit1(bytecode.OpMoveResult2, lhs2slot)
	}
}

func (cl *compiler) compileIncDecStmt(stmt *ast.IncDecStmt) {
	varname, ok := stmt.X.(*ast.Ident)
	if !ok {
		panic(cl.errorf(stmt.X, "can assign only to simple variables"))
	}
	dst := cl.getNamedSlot(varname, varname.String())
	if stmt.Tok == token.INC {
		cl.emit1(bytecode.OpIntInc, dst)
	} else {
		cl.emit1(bytecode.OpIntDec, dst)
	}
}

func (cl *compiler) compileBranchStmt(branch *ast.BranchStmt) {
	if branch.Label != nil {
		panic(cl.errorf(branch.Label, "can't compile %s with a label", branch.Tok))
	}

	switch branch.Tok {
	case token.BREAK:
		cl.emitJump(cl.breakTarget)
	case token.CONTINUE:
		cl.emitJump(cl.continueTarget)
	default:
		panic(cl.errorf(branch, "can't compile %s yet", branch.Tok))
	}
}

func (cl *compiler) compileForStmt(stmt *ast.ForStmt) {
	cl.hasLoops = true

	labelBreak := cl.newLabel()
	labelContinue := cl.newLabel()
	prevBreakTarget := cl.breakTarget
	prevContinueTarget := cl.continueTarget
	cl.breakTarget = labelBreak
	cl.continueTarget = labelContinue

	switch {
	case stmt.Cond != nil && stmt.Init == nil && stmt.Post == nil:
		// `for <cond> { ... }`
		labelBody := cl.newLabel()
		cl.emitJump(labelContinue)
		cl.bindLabel(labelBody)
		cl.compileStmt(stmt.Body)
		cl.bindLabel(labelContinue)
		condslot := cl.compileRootTempExpr(stmt.Cond)
		cl.emitCondJump(condslot, bytecode.OpJumpNotZero, labelBody)
		cl.bindLabel(labelBreak)

	case stmt.Cond == nil && stmt.Init == nil && stmt.Post == nil:
		// `for { ... }`
		cl.bindLabel(labelContinue)
		cl.compileStmt(stmt.Body)
		cl.emitJump(labelContinue)
		cl.bindLabel(labelBreak)

	default:
		// `for <init>; <cond>; <post> { ... }`
		labelStart := cl.newLabel()
		labelBody := cl.newLabel()
		if stmt.Init != nil {
			cl.compileStmt(stmt.Init)
		}
		if stmt.Cond != nil {
			cl.emitJump(labelStart)
		}
		cl.bindLabel(labelBody)
		cl.compileStmt(stmt.Body)
		cl.bindLabel(labelContinue)
		if stmt.Post != nil {
			cl.compileStmt(stmt.Post)
		}
		cl.bindLabel(labelStart)
		if stmt.Cond != nil {
			condslot := cl.compileRootTempExpr(stmt.Cond)
			cl.emitCondJump(condslot, bytecode.OpJumpNotZero, labelBody)
		} else {
			cl.emitJump(labelBody)
		}
		cl.bindLabel(labelBreak)
	}

	cl.breakTarget = prevBreakTarget
	cl.continueTarget = prevContinueTarget
}

func (cl *compiler) compileExprStmt(stmt *ast.ExprStmt) {
	if call, ok := stmt.X.(*ast.CallExpr); ok {
		sig := cl.ctx.Types.TypeOf(call.Fun).(*types.Signature)
		if sig.Results() != nil {
			panic(cl.errorf(call, "only void funcs can be used in stmt context"))
		}
		cl.compileRootExpr(voidSlot, call)
		return
	}

	panic(cl.errorf(stmt.X, "can't compile this expr stmt yet: %T", stmt.X))
}
