package qcompile

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/quasilyte/quasigo/internal/ir"
	"github.com/quasilyte/quasigo/internal/qruntime"
)

type Context struct {
	Env *qruntime.Env

	Optimize bool

	Package *types.Package
	Types   *types.Info
	Sizes   types.Sizes
	Fset    *token.FileSet
}

type Compiler struct {
	instPool []ir.Inst
}

func NewCompiler() *Compiler {
	return &Compiler{
		instPool: make([]ir.Inst, 0, 128),
	}
}

func (c *Compiler) CompileFunc(ctx *Context, fn *ast.FuncDecl) (compiled *qruntime.Func, err error) {
	defer func() {
		if err != nil {
			return
		}
		rv := recover()
		if rv == nil {
			return
		}
		if compileErr, ok := rv.(compileError); ok {
			err = compileErr
			return
		}
		panic(rv) // not our panic
	}()

	p := patternCompiler{}
	cl := compiler{
		code: c.instPool[:0],
		ctx:  ctx,

		fnType: ctx.Types.ObjectOf(fn.Name).Type().(*types.Signature),

		strConstantsPool:    make(map[string]int),
		scalarConstantsPool: make(map[uint64]int),
		locals:              make(map[string]frameSlotInfo),

		patternCompiler: &p,
	}
	p.cl = &cl
	return cl.compileFunc(fn), nil
}
