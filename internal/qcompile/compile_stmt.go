package qcompile

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/quasilyte/quasigo/internal/bytecode"
)

func (cl *compiler) compileStmt(stmt ast.Stmt) {
	switch stmt := stmt.(type) {
	case *ast.BlockStmt:
		for i := range stmt.List {
			cl.compileStmt(stmt.List[i])
		}

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

func (cl *compiler) compileReturnStmt(ret *ast.ReturnStmt) {
	if cl.retType == voidType {
		cl.emit(bytecode.OpReturnVoid)
		return
	}

	if ret.Results == nil {
		panic(cl.errorf(ret, "'naked' return statements are not allowed"))
	}

	switch {
	case identName(ret.Results[0]) == "true":
		cl.emit(bytecode.OpReturnTrue)
	case identName(ret.Results[0]) == "false":
		cl.emit(bytecode.OpReturnFalse)
	default:
		typ := cl.ctx.Types.TypeOf(ret.Results[0])
		var op bytecode.Op
		switch {
		case typeIsScalar(typ):
			op = bytecode.OpReturnScalar
		case typeIsString(typ):
			op = bytecode.OpReturnStr
		case typeIsInterface(typ) || typeIsPointer(typ):
			op = bytecode.OpReturnInterface
		default:
			panic(cl.errorf(ret, "can't return %s typed value yet", typ.String()))
		}
		slot := cl.compileRootTempExpr(ret.Results[0])
		cl.emit8(op, slot)
	}
}

func (cl *compiler) compileIfStmt(stmt *ast.IfStmt) {
	if stmt.Else == nil {
		labelEnd := cl.newLabel()
		condslot := cl.compileRootTempExpr(stmt.Cond)
		cl.emitCondJump(condslot, bytecode.OpJumpFalse, labelEnd)
		cl.compileStmt(stmt.Body)
		cl.bindLabel(labelEnd)
		return
	}

	labelEnd := cl.newLabel()
	labelElse := cl.newLabel()
	condslot := cl.compileRootTempExpr(stmt.Cond)
	cl.emitCondJump(condslot, bytecode.OpJumpFalse, labelElse)
	cl.compileStmt(stmt.Body)
	if !cl.isUncondJump(cl.lastOp) {
		cl.emitJump(bytecode.OpJump, labelEnd)
	}
	cl.bindLabel(labelElse)
	cl.compileStmt(stmt.Else)
	cl.bindLabel(labelEnd)
}

func (cl *compiler) compileAssignStmt(assign *ast.AssignStmt) {
	if len(assign.Rhs) != 1 {
		panic(cl.errorf(assign, "only single right bytecode.Operand is allowed in assignments"))
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

	dst1 := assign.Lhs[0].(*ast.Ident)
	rhs := assign.Rhs[0]
	lhs1slot := cl.getLocal(dst1, dst1.Name)
	cl.compileRootExpr(lhs1slot, rhs)
	if len(assign.Lhs) == 2 {
		dst2 := assign.Lhs[1].(*ast.Ident)
		lhs2slot := cl.getLocal(dst2, dst2.Name)
		cl.emit8(bytecode.OpMoveResult2, lhs2slot)
	}
}

func (cl *compiler) compileIncDecStmt(stmt *ast.IncDecStmt) {
	varname, ok := stmt.X.(*ast.Ident)
	if !ok {
		panic(cl.errorf(stmt.X, "can assign only to simple variables"))
	}
	dst := cl.getLocal(varname, varname.String())
	if stmt.Tok == token.INC {
		cl.emit8(bytecode.OpIntInc, dst)
	} else {
		cl.emit8(bytecode.OpIntDec, dst)
	}
}

func (cl *compiler) compileBranchStmt(branch *ast.BranchStmt) {
	if branch.Label != nil {
		panic(cl.errorf(branch.Label, "can't compile %s with a label", branch.Tok))
	}

	switch branch.Tok {
	case token.BREAK:
		cl.emitJump(bytecode.OpJump, cl.breakTarget)
	default:
		panic(cl.errorf(branch, "can't compile %s yet", branch.Tok))
	}
}

func (cl *compiler) compileForStmt(stmt *ast.ForStmt) {
	labelBreak := cl.newLabel()
	labelContinue := cl.newLabel()
	prevBreakTarget := cl.breakTarget
	prevContinueTarget := cl.continueTarget
	cl.breakTarget = labelBreak
	cl.continueTarget = labelContinue

	switch {
	case stmt.Cond != nil && stmt.Init != nil && stmt.Post != nil:
		// Will be implemented later.
		panic(cl.errorf(stmt, "can't compile C-style for lobytecode.Ops yet"))

	case stmt.Cond != nil && stmt.Init == nil && stmt.Post == nil:
		// `for <cond> { ... }`
		labelBody := cl.newLabel()
		cl.emitJump(bytecode.OpJump, labelContinue)
		cl.bindLabel(labelBody)
		cl.compileStmt(stmt.Body)
		cl.bindLabel(labelContinue)
		condslot := cl.compileRootTempExpr(stmt.Cond)
		cl.emitCondJump(condslot, bytecode.OpJumpTrue, labelBody)
		cl.bindLabel(labelBreak)

	default:
		// `for { ... }`
		cl.bindLabel(labelContinue)
		cl.compileStmt(stmt.Body)
		cl.emitJump(bytecode.OpJump, labelContinue)
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
