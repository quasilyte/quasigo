package qcompile

import (
	"go/ast"
	"go/token"
	"go/types"
	"math"

	"github.com/quasilyte/quasigo/internal/ir"
	"github.com/quasilyte/quasigo/internal/qruntime"
)

type Context struct {
	Env *qruntime.Env

	Optimize bool
	Static   bool

	Package *types.Package
	Types   *types.Info
	Sizes   types.Sizes
	Fset    *token.FileSet

	TestingContext TestingContext

	// pending holds the funcs of the unit that is being compiled.
	// It is nil when compiling a standalone func.
	pending map[qruntime.FuncKey]*pendingFunc
}

// pendingFunc is a func of a compilation unit that has a reserved ID,
// but may not be compiled (and therefore bound to Env) yet.
type pendingFunc struct {
	id uint16
	fn *qruntime.Func
}

type TestingContext interface {
	FuncIR(fn *ir.Func)
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
	cl := funcCompiler{
		code: c.instPool[:0],
		ctx:  ctx,

		fnType: ctx.Types.ObjectOf(fn.Name).Type().(*types.Signature),

		strConstantsPool:    make(map[string]int),
		scalarConstantsPool: make(map[uint64]int),

		patternCompiler: &p,
	}
	p.cl = &cl
	return cl.compileFunc(fn), nil
}

type UnitCompilationConfig struct {
	Funcs []*ast.FuncDecl

	OnFuncCompiled func(compiled *qruntime.Func)
}

// CompileUnit compiles all funcs as one unit.
func (c *Compiler) CompileUnit(ctx *Context, config UnitCompilationConfig) error {
	base := len(ctx.Env.UserFuncs)
	if base+len(config.Funcs) > math.MaxUint16 {
		return compileError("too many funcs")
	}

	unitCtx := *ctx
	unitCtx.pending = make(map[qruntime.FuncKey]*pendingFunc, len(config.Funcs))
	entries := make([]pendingFunc, len(config.Funcs))

	// Pass 1: reserve the func IDs.
	for i, decl := range config.Funcs {
		key := qruntime.FuncKey{Qualifier: ctx.Package.Path(), Name: decl.Name.Name}
		if _, ok := unitCtx.pending[key]; ok {
			return cerrorf(ctx.Fset, decl.Name.Pos(), "%s is redeclared", decl.Name.Name)
		}
		entries[i] = pendingFunc{id: uint16(base + i)}
		unitCtx.pending[key] = &entries[i]
	}

	// Pass 2: compile all funcs.
	for i, decl := range config.Funcs {
		compiled, err := c.CompileFunc(&unitCtx, decl)
		if err != nil {
			return err
		}
		entries[i].fn = compiled
		config.OnFuncCompiled(compiled)
	}

	return nil
}

// lookupUserFunc finds a user func by its key, looking at the unit being
// compiled first. The returned fn is nil if the func is a part of the
// current unit and has not been compiled yet.
func (cl *funcCompiler) lookupUserFunc(key qruntime.FuncKey) (id uint16, fn *qruntime.Func, ok bool) {
	if p, found := cl.ctx.pending[key]; found {
		return p.id, p.fn, true
	}

	id, ok = cl.ctx.Env.NameToFuncID[key]
	if !ok {
		return 0, nil, false
	}

	return id, cl.ctx.Env.UserFuncs[id], true
}
