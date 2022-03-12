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
		cl.compileBinarySearch(defaultBody, constCases)
	} else {
		cl.compileLinearSearch(stmt, defaultBody)
	}
}

func (cl switchCompiler) compileLinearSearch(stmt *ast.SwitchStmt, defaultBody []ast.Stmt) {
	labelEnd := cl.newLabel()

	for _, cc := range stmt.Body.List {
		cc := cc.(*ast.CaseClause)
		if cc.List == nil {
			continue // Default clause
		}
		labelNext := cl.newLabel()
		condtmp := cl.allocTmp()
		yslot := cl.allocTmp()
		if len(cc.List) == 1 {
			cl.CompileExpr(yslot, cc.List[0])
			cl.emit3(cl.opEq, condtmp, cl.tagslot, yslot)
			cl.emitCondJump(condtmp, bytecode.OpJumpZero, labelNext)
			cl.freeTmp()
		} else {
			// TODO: figure out a better way?
			labelMatched := cl.newLabel()
			for _, y := range cc.List {
				cl.CompileExpr(yslot, y)
				cl.emit3(cl.opEq, condtmp, cl.tagslot, yslot)
				cl.emitCondJump(condtmp, bytecode.OpJumpNotZero, labelMatched)
			}
			cl.emitJump(labelNext)
			cl.freeTmp()
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
		// if tag < firstCase { goto default }
		firstCase := cases[0]
		condtmp := cl.allocTmp()
		yslot := cl.allocTmp()
		cl.compileConstantValue(yslot, firstCase.expr, firstCase.value)
		cl.emit3(cl.opLt, condtmp, cl.tagslot, yslot)
		cl.emitCondJump(condtmp, bytecode.OpJumpNotZero, labelDefault)

		// if tag > lastCase { goto default }
		lastLase := cases[len(cases)-1]
		cl.compileConstantValue(yslot, lastLase.expr, lastLase.value)
		cl.emit3(cl.opGt, condtmp, cl.tagslot, yslot)
		cl.emitCondJump(condtmp, bytecode.OpJumpNotZero, labelDefault)

		cl.freeTmp()
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
			condtmp := cl.allocTmp()
			yslot := cl.allocTmp()
			for i := left; i <= right; i++ {
				c := &cases[i]
				c.label = cl.newLabel()
				cl.compileConstantValue(yslot, c.expr, c.value)
				cl.emit3(cl.opEq, condtmp, cl.tagslot, yslot)
				cl.emitCondJump(condtmp, bytecode.OpJumpNotZero, c.label)
				if i == right {
					cl.emitJump(labelDefault)
				}
			}
			cl.freeTmp()
			return
		}

		// Have many search cases, divide the search space by 2.

		mid := left + (right-left)/2
		c := &cases[mid]

		c.label = cl.newLabel()

		// perform `if tag == c.value { goto c.body }`
		condtmp := cl.allocTmp()
		yslot := cl.allocTmp()
		cl.compileConstantValue(yslot, c.expr, c.value)
		cl.emit3(cl.opEq, condtmp, cl.tagslot, yslot)
		cl.emitCondJump(condtmp, bytecode.OpJumpNotZero, c.label)

		// perform `if tag > c.value { goto labelGreater }`
		labelGreater := cl.newLabel()
		cl.emit3(cl.opGt, condtmp, cl.tagslot, yslot)
		cl.emitCondJump(condtmp, bytecode.OpJumpNotZero, labelGreater)
		cl.freeTmp()
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