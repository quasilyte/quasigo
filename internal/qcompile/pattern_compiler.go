package qcompile

import (
	"go/ast"

	"github.com/quasilyte/quasigo/internal/goutil"
)

// TODO:
// * switch with only default

// patternCompiler implements some adhoc optimizations based on the AST.
// These optimizations are not optional, they're the part of the normal compilation.
// This type is separated from the compiler to decompose the work a little bit.
type patternCompiler struct {
	cl *compiler
}

func (p *patternCompiler) CompileSliceExpr(dst int, slice *ast.SliceExpr) bool {
	// Try to recognize the no-op slicing.
	// Examples: s[:], s[0:], s[:len(s)], s[0:len(s)].
	// All of these can be simplified to just s.
	fromFirst := false
	toLast := false
	if slice.Low == nil {
		fromFirst = true
	} else {
		v, ok := intValueOf(p.cl.ctx.Types, slice.Low)
		fromFirst = ok && v == 0
	}
	if slice.High == nil {
		toLast = true
	} else {
		// Try to match `len(x)` where `x` is slice.X.
		asCall, ok := slice.High.(*ast.CallExpr)
		toLast = ok && goutil.MatchBuiltin(p.cl.ctx.Types, asCall.Fun, `len`) &&
			isSameExpr(asCall.Args[0], slice.X)
	}
	if fromFirst && toLast {
		p.cl.CompileExpr(dst, slice.X)
		return true
	}

	return false
}
