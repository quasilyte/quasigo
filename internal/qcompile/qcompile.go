package qcompile

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/quasilyte/quasigo/internal/qruntime"
)

type Context struct {
	Env *qruntime.Env

	Package *types.Package
	Types   *types.Info
	Fset    *token.FileSet
}

func Func(ctx *Context, fn *ast.FuncDecl) (*qruntime.Func, error) {
	return compile(ctx, fn)
}
