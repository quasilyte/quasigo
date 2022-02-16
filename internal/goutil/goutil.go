package goutil

import (
	"go/ast"
	"go/types"
)

// Unparen returns e with any enclosing parentheses stripped.
func Unparen(e ast.Expr) ast.Expr {
	for {
		p, ok := e.(*ast.ParenExpr)
		if !ok {
			return e
		}
		e = p.X
	}
}

func ResolveFunc(info *types.Info, callable ast.Expr) (ast.Expr, *types.Func) {
	switch callable := Unparen(callable).(type) {
	case *ast.Ident:
		sig, ok := info.ObjectOf(callable).(*types.Func)
		if !ok {
			return nil, nil
		}
		return nil, sig

	case *ast.SelectorExpr:
		sig, ok := info.ObjectOf(callable.Sel).(*types.Func)
		if !ok {
			return nil, nil
		}
		isMethod := sig.Type().(*types.Signature).Recv() != nil
		if _, ok := callable.X.(*ast.Ident); ok && !isMethod {
			return nil, sig
		}
		return callable.X, sig

	default:
		return nil, nil
	}
}
