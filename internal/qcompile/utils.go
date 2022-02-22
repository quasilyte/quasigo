package qcompile

import (
	"encoding/binary"
	"go/ast"
	"go/constant"
	"go/types"

	"github.com/quasilyte/quasigo/internal/bytecode"
)

// isSameExpr is a simple form of expressions comparison operation.
// It's faster than astequal.Expr and it ignores some complicated expressions.
// We use it in the contexts where we don't need precise matching.
func isSameExpr(x, y ast.Expr) bool {
	switch x := x.(type) {
	case *ast.Ident:
		y, ok := y.(*ast.Ident)
		return ok && x.Name == y.Name
	case *ast.SelectorExpr:
		y, ok := y.(*ast.SelectorExpr)
		return ok && x.Sel.Name == y.Sel.Name && isSameExpr(x.X, y.X)
	default:
		return false
	}
}

func intValueOf(info *types.Info, e ast.Expr) (int64, bool) {
	cv := info.Types[e].Value
	if cv == nil || cv.Kind() != constant.Int {
		return 0, false
	}
	v, exact := constant.Int64Val(cv)
	if !exact {
		return 0, false
	}
	return v, true
}

func pickOp(cond bool, ifTrue, otherwise bytecode.Op) bytecode.Op {
	if cond {
		return ifTrue
	}
	return otherwise
}

func put16(code []byte, pos, value int) {
	binary.LittleEndian.PutUint16(code[pos:], uint16(value))
}

func typeIsScalar(typ types.Type) bool {
	basic, ok := typ.Underlying().(*types.Basic)
	if !ok {
		return false
	}
	switch basic.Kind() {
	case types.Int, types.UntypedInt, types.Bool, types.UntypedBool, types.Uint8:
		return true
	default:
		return false
	}
}

func typeIsBool(typ types.Type) bool {
	basic, ok := typ.Underlying().(*types.Basic)
	if !ok {
		return false
	}
	return basic.Info()&types.IsBoolean != 0
}

func typeIsByte(typ types.Type) bool {
	basic, ok := typ.Underlying().(*types.Basic)
	if !ok {
		return false
	}
	return basic.Kind() == types.Uint8
}

func typeIsInt(typ types.Type) bool {
	basic, ok := typ.Underlying().(*types.Basic)
	if !ok {
		return false
	}
	switch basic.Kind() {
	case types.Int, types.UntypedInt:
		return true
	default:
		return false
	}
}

func typeIsPointer(typ types.Type) bool {
	_, ok := typ.Underlying().(*types.Pointer)
	return ok
}

func typeIsInterface(typ types.Type) bool {
	_, ok := typ.Underlying().(*types.Interface)
	return ok
}

func typeIsString(typ types.Type) bool {
	basic, ok := typ.Underlying().(*types.Basic)
	if !ok {
		return false
	}
	return basic.Info()&types.IsString != 0
}

func identName(n ast.Expr) string {
	id, ok := n.(*ast.Ident)
	if ok {
		return id.Name
	}
	return ""
}
