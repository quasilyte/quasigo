package qcompile

import (
	"go/ast"
	"go/constant"
	"go/token"
	"sort"

	"github.com/quasilyte/quasigo/internal/bytecode"
	"github.com/quasilyte/quasigo/internal/ir"
)

type constCaseClause struct {
	value constant.Value
	expr  ast.Expr
	label label
	body  []ast.Stmt
}

type switchCompiler struct {
	*compiler
	tagslot ir.Slot
	opEq    bytecode.Op
	opLt    bytecode.Op
	opGt    bytecode.Op
}

func (cl *compiler) compileSwitchStmt(stmt *ast.SwitchStmt) {
	if stmt.Init != nil {
		panic(cl.errorf(stmt, "can't compile switch with init clause yet"))
	}
	if stmt.Tag == nil {
		panic(cl.errorf(stmt, "can't compile tagless switch yet"))
	}

	scl := switchCompiler{compiler: cl}
	tagType := cl.ctx.Types.TypeOf(stmt.Tag)
	switch {
	case typeIsInt(tagType) || typeIsByte(tagType):
		scl.opEq = bytecode.OpScalarEq
		scl.opLt = bytecode.OpIntLt
		scl.opGt = bytecode.OpIntGt
	case typeIsString(tagType):
		scl.opEq = bytecode.OpStrEq
		scl.opLt = bytecode.OpStrLt
		scl.opGt = bytecode.OpStrGt
	default:
		panic(cl.errorf(stmt.Tag, "can't compile switch over a value of type %s", tagType))
	}

	scl.tagslot = cl.allocAutoLocal()
	defer func() { cl.autoLocalSeq-- }()
	cl.CompileExpr(scl.tagslot, stmt.Tag)

	labelBreak := cl.newLabel()
	prevBreakTarget := cl.breakTarget
	cl.breakTarget = labelBreak
	scl.compile(stmt)
	cl.bindLabel(cl.breakTarget)
	cl.breakTarget = prevBreakTarget
}

func (cl switchCompiler) compile(stmt *ast.SwitchStmt) {
	collectConstCases := len(stmt.Body.List) >= 6

	var constCases []constCaseClause
	var defaultBody []ast.Stmt
	for _, cc := range stmt.Body.List {
		cc := cc.(*ast.CaseClause)
		if cc.List == nil { // Default case
			defaultBody = cc.Body
			continue
		}
		if !collectConstCases {
			continue
		}
		e := cc.List[0]
		cv := cl.ctx.Types.Types[e].Value
		if len(cc.List) != 1 || cv == nil {
			collectConstCases = false
			constCases = nil
			continue
		}
		constCases = append(constCases, constCaseClause{
			value: cv,
			expr:  e,
			body:  cc.Body,
		})
	}

	numCases := len(stmt.Body.List)
	if defaultBody != nil {
		numCases--
	}
	if len(constCases) == numCases {
		if cl.canCompileAsTable(constCases) {
			cl.compileTableSearch(defaultBody, constCases)
		} else {
			cl.compileBinarySearch(defaultBody, constCases)
		}
	} else {
		cl.compileLinearSearch(stmt, defaultBody)
	}
}

func (cl switchCompiler) canCompileAsTable(cases []constCaseClause) bool {
	if cl.opEq != bytecode.OpScalarEq || len(cases) < 8 {
		return false
	}
	first := cases[0].value
	if first.Kind() != constant.Int {
		return false
	}
	firstVal, exact := constant.Int64Val(first)
	if !exact || firstVal != 0 {
		return false
	}
	last := cases[len(cases)-1].value
	lastVal, exact := constant.Int64Val(last)
	if !exact || lastVal != int64(len(cases)-1) {
		return false
	}
	return true
}

func (cl switchCompiler) compileTableSearch(defaultBody []ast.Stmt, cases []constCaseClause) {
	labelEnd := cl.newLabel()
	labelDefault := cl.newLabel()

	// Like with a binary search, we need a value range check.
	{
		cl.beginTempBlock()

		// if tag < firstCase { goto default }
		firstCase := cases[0]
		condtemp := cl.allocTemp()
		yslot := cl.allocTemp()
		cl.compileConstantValue(yslot, firstCase.expr, firstCase.value)
		cl.emit3(cl.opLt, condtemp, cl.tagslot, yslot)
		cl.emitCondJump(condtemp, bytecode.OpJumpNotZero, labelDefault)

		// if tag > lastCase { goto default }
		lastLase := cases[len(cases)-1]
		cl.compileConstantValue(yslot, lastLase.expr, lastLase.value)
		cl.emit3(cl.opGt, condtemp, cl.tagslot, yslot)
		cl.emitCondJump(condtemp, bytecode.OpJumpNotZero, labelDefault)

		cl.endTempBlock()
	}

	cl.emit(ir.Inst{Op: bytecode.OpJumpTable, Arg0: cl.tagslot.ToInstArg()})
	for i := range cases {
		c := &cases[i]
		c.label = cl.newLabel()
		cl.emitJump(c.label)
	}

	for _, c := range cases {
		cl.bindLabel(c.label)
		cl.compileStmtList(c.body)
		if !cl.isUncondJump(cl.lastOp()) {
			cl.emitJump(labelEnd)
		}
	}

	cl.bindLabel(labelDefault)
	if defaultBody != nil {
		cl.compileStmtList(defaultBody)
	}
	cl.bindLabel(labelEnd)
}

func (cl switchCompiler) compileLinearSearch(stmt *ast.SwitchStmt, defaultBody []ast.Stmt) {
	labelEnd := cl.newLabel()

	for _, cc := range stmt.Body.List {
		cc := cc.(*ast.CaseClause)
		if cc.List == nil {
			continue // Default clause
		}
		labelNext := cl.newLabel()
		cl.beginTempBlock()
		condtemp := cl.allocTemp()
		yslot := cl.allocTemp()
		if len(cc.List) == 1 {
			cl.CompileExpr(yslot, cc.List[0])
			cl.emit3(cl.opEq, condtemp, cl.tagslot, yslot)
			cl.emitCondJump(condtemp, bytecode.OpJumpZero, labelNext)
			cl.endTempBlock()
		} else {
			// TODO: figure out a better way?
			labelMatched := cl.newLabel()
			for _, y := range cc.List {
				cl.CompileExpr(yslot, y)
				cl.emit3(cl.opEq, condtemp, cl.tagslot, yslot)
				cl.emitCondJump(condtemp, bytecode.OpJumpNotZero, labelMatched)
			}
			cl.emitJump(labelNext)
			cl.endTempBlock()
			cl.bindLabel(labelMatched)
		}
		cl.compileStmtList(cc.Body)
		if !cl.isUncondJump(cl.lastOp()) {
			cl.emitJump(labelEnd)
		}
		cl.bindLabel(labelNext)
	}

	if defaultBody != nil {
		cl.compileStmtList(defaultBody)
	}
	cl.bindLabel(labelEnd)
}

func (cl switchCompiler) compileBinarySearch(defaultBody []ast.Stmt, cases []constCaseClause) {
	const linearSearchLen = 4

	// Sort cases by their values in ascending order.
	sort.Slice(cases, func(i, j int) bool {
		x := cases[i]
		y := cases[j]
		return constant.Compare(x.value, token.LSS, y.value)
	})

	labelEnd := cl.newLabel()
	labelDefault := cl.newLabel()

	// Insert value bound checks. If tag values is out of bounds,
	// we don't need to perform the search.
	{
		cl.beginTempBlock()

		// if tag < firstCase { goto default }
		firstCase := cases[0]
		condtemp := cl.allocTemp()
		yslot := cl.allocTemp()
		cl.compileConstantValue(yslot, firstCase.expr, firstCase.value)
		cl.emit3(cl.opLt, condtemp, cl.tagslot, yslot)
		cl.emitCondJump(condtemp, bytecode.OpJumpNotZero, labelDefault)

		// if tag > lastCase { goto default }
		lastLase := cases[len(cases)-1]
		cl.compileConstantValue(yslot, lastLase.expr, lastLase.value)
		cl.emit3(cl.opGt, condtemp, cl.tagslot, yslot)
		cl.emitCondJump(condtemp, bytecode.OpJumpNotZero, labelDefault)

		cl.endTempBlock()
	}

	var walkTree func(left, right int)
	walkTree = func(left, right int) {
		if left > right {
			return
		}

		numCases := right - left
		if numCases <= linearSearchLen {
			// Have a few cases left in this path, should
			// switch to a linear search for this remainder.
			cl.beginTempBlock()
			condtemp := cl.allocTemp()
			yslot := cl.allocTemp()
			for i := left; i <= right; i++ {
				c := &cases[i]
				c.label = cl.newLabel()
				cl.compileConstantValue(yslot, c.expr, c.value)
				cl.emit3(cl.opEq, condtemp, cl.tagslot, yslot)
				cl.emitCondJump(condtemp, bytecode.OpJumpNotZero, c.label)
				if i == right {
					cl.emitJump(labelDefault)
				}
			}
			cl.endTempBlock()
			return
		}

		// Have many search cases, divide the search space by 2.

		mid := left + (right-left)/2
		c := &cases[mid]

		c.label = cl.newLabel()
		cl.beginTempBlock()

		// perform `if tag == c.value { goto c.body }`
		condtemp := cl.allocTemp()
		yslot := cl.allocTemp()
		cl.compileConstantValue(yslot, c.expr, c.value)
		cl.emit3(cl.opEq, condtemp, cl.tagslot, yslot)
		cl.emitCondJump(condtemp, bytecode.OpJumpNotZero, c.label)

		// perform `if tag > c.value { goto labelGreater }`
		labelGreater := cl.newLabel()
		cl.emit3(cl.opGt, condtemp, cl.tagslot, yslot)
		cl.emitCondJump(condtemp, bytecode.OpJumpNotZero, labelGreater)
		cl.endTempBlock()
		walkTree(left, mid-1)

		cl.bindLabel(labelGreater)
		walkTree(mid+1, right)
	}
	walkTree(0, len(cases)-1)

	for _, c := range cases {
		cl.bindLabel(c.label)
		cl.compileStmtList(c.body)
		if !cl.isUncondJump(cl.lastOp()) {
			cl.emitJump(labelEnd)
		}
	}

	cl.bindLabel(labelDefault)
	if defaultBody != nil {
		cl.compileStmtList(defaultBody)
	}
	cl.bindLabel(labelEnd)
}
