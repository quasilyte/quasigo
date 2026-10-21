package qcompile

import (
	"go/ast"
	"go/token"
	"go/types"
	"math"

	"github.com/quasilyte/quasigo/internal/goutil"
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
	// The callees are compiled before their callers whenever possible,
	// so the callers can inline them.
	for _, i := range c.funcCompileOrder(&unitCtx, config.Funcs) {
		compiled, err := c.CompileFunc(&unitCtx, config.Funcs[i])
		if err != nil {
			return err
		}
		entries[i].fn = compiled
	}

	// The results are reported in the declaration order,
	// as it matches the reserved IDs.
	for i := range entries {
		config.OnFuncCompiled(entries[i].fn)
	}

	return nil
}

// compileOrder returns the indexes of funcs in the order they should be compiled.
//
// It's a post-order walk over the call graph, so a callee goes before its callers.
// Ties are resolved by the declaration order, so the order is stable.
// If there are cycles (recursion), they are broken at an arbitrary edge:
// funcs that call each other can't be inlined anyway.
func (c *Compiler) funcCompileOrder(ctx *Context, funcs []*ast.FuncDecl) []int {
	if !ctx.Optimize {
		order := make([]int, len(funcs))
		for i := range funcs {
			order[i] = i
		}
		return order
	}

	indexes := make(map[qruntime.FuncKey]int, len(funcs))
	for i, decl := range funcs {
		key := qruntime.FuncKey{Qualifier: ctx.Package.Path(), Name: decl.Name.Name}
		indexes[key] = i
	}

	callees := make([][]int, len(funcs))
	for i, decl := range funcs {
		ast.Inspect(decl.Body, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			_, fn := goutil.ResolveFunc(ctx.Types, call.Fun)
			if fn == nil || fn.Pkg() == nil || fn.Type().(*types.Signature).Recv() != nil {
				return true
			}
			key := qruntime.FuncKey{Qualifier: fn.Pkg().Path(), Name: fn.Name()}
			if j, ok := indexes[key]; ok {
				callees[i] = append(callees[i], j)
			}
			return true
		})
	}

	const (
		stateUnvisited = iota
		stateVisiting
		stateDone
	)
	states := make([]byte, len(funcs))
	order := make([]int, 0, len(funcs))
	var visit func(i int)
	visit = func(i int) {
		if states[i] != stateUnvisited {
			// Either done, or we're inside a cycle.
			return
		}
		states[i] = stateVisiting
		for _, j := range callees[i] {
			visit(j)
		}
		states[i] = stateDone
		order = append(order, i)
	}
	for i := range funcs {
		visit(i)
	}
	return order
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
