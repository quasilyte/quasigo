// Package quasigo implements a Go subset compiler and interpreter.
//
// The implementation details are not part of the contract of this package.
package quasigo

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"github.com/quasilyte/quasigo/internal/qcompile"
	"github.com/quasilyte/quasigo/internal/qdisasm"
	"github.com/quasilyte/quasigo/internal/qruntime"
	"github.com/quasilyte/quasigo/qnative"
)

// TODO(quasilyte): document what is thread-safe and what not.
// TODO(quasilyte): add a readme.

// Env is used to hold both compilation and evaluation data.
type Env struct {
	data qruntime.Env
}

// EvalEnv is a goroutine-local handle for Env.
// To get one, use Env.GetEvalEnv() method.
type EvalEnv struct {
	data qruntime.EvalEnv
}

// NewEnv creates a new empty environment.
func NewEnv() *Env {
	env := &Env{}
	qruntime.InitEnv(&env.data)
	return env
}

// GetEvalEnv creates a new goroutine-local handle of env.
// Stack size is amount of bytes we allocate for all stack
// frames of this env.
func (env *Env) GetEvalEnv(stackSize int) *EvalEnv {
	ee := &EvalEnv{}
	qruntime.InitEvalEnv(&env.data, &ee.data, stackSize)
	return ee
}

// AddNativeMethod binds `$typeName.$methodName` symbol with f.
// A typeName should be fully qualified, like `github.com/user/pkgname.TypeName`.
// It method is defined on pointer type, the typeName should start with `*`,
// like `*github.com/user/pkgname.TypeName`.
func (env *Env) AddNativeMethod(typeName, methodName string, f func(ctx qnative.CallContext)) {
	env.data.AddNativeMethod(typeName, methodName, f)
}

// AddNativeFunc binds `$pkgPath.$funcName` symbol with f.
// A pkgPath should be a full package path in which funcName is defined.
//
// Generic funcs are bound per instantiation, the type args are
// a part of the funcName: `If[int]`, `Pair[string,[]uint8]`.
// Only the bound instantiations can be called.
// Generic stubs need a body to pass the type checking, but
// it's never compiled, so `var zero T; return zero` is enough.
// The type args spelling rules:
//   - predeclared aliases are interchangeable: byte/uint8, rune/int32, interface{}/any
//   - user-defined aliases should be spelled as the types they refer to
//   - named types are fully qualified: `*github.com/user/pkgname.TypeName`
//   - type args are separated by a comma without spaces
//
// The type args are taken from types.Info.Instances, so it should be
// populated (see CompileContext.Types).
func (env *Env) AddNativeFunc(pkgPath, funcName string, f func(ctx qnative.CallContext)) {
	env.data.AddNativeFunc(pkgPath, funcName, f)
}

// AddFunc binds f, it will be available under pkgPath and funcName
// associated with f.
func (env *Env) AddFunc(f Func) {
	dotIndex := strings.LastIndexByte(f.Name(), '.')
	pkgPath := f.Name()[:dotIndex]
	funcName := f.Name()[dotIndex+1:]
	env.data.AddFunc(pkgPath, funcName, f.data)
}

// GetFunc finds previously bound function searching for the `$pkgPath.$funcName` symbol.
func (env *Env) GetFunc(pkgPath, funcName string) Func {
	return Func{data: env.data.GetFunc(pkgPath, funcName)}
}

// CompileContext is used to provide necessary data to the compiler.
type CompileContext struct {
	// Env is shared environment that should be used for all functions
	// being compiled; then it should be used to execute these functions.
	Env *Env

	Package *types.Package
	Types   *types.Info
	Sizes   types.Sizes
	Fset    *token.FileSet

	Optimize bool
	Static   bool

	TestingContext any
}

func (ctx *CompileContext) internal() qcompile.Context {
	internalCtx := qcompile.Context{
		Env:      &ctx.Env.data,
		Optimize: ctx.Optimize,
		Static:   ctx.Static,
		Package:  ctx.Package,
		Types:    ctx.Types,
		Sizes:    ctx.Sizes,
		Fset:     ctx.Fset,
	}
	if testingContext, ok := ctx.TestingContext.(qcompile.TestingContext); ok {
		internalCtx.TestingContext = testingContext
	}
	return internalCtx
}

// Compile prepares an executable version of fn.
//
// Calls to other funcs are resolved using ctx.Env, so a callee
// should be already bound (see Env.AddFunc) by the time its caller is compiled.
// Use CompileUnit if funcs call each other in an arbitrary order.
func Compile(ctx *CompileContext, fn *ast.FuncDecl) (Func, error) {
	internalCtx := ctx.internal()
	compiler := qcompile.NewCompiler()
	compiled, err := compiler.CompileFunc(&internalCtx, fn)
	return Func{data: compiled}, err
}

// CompilationUnit is a set of declarations that are compiled together.
// It can be used for an entire package compilation (or its selected portion).
type CompilationUnit struct {
	// Funcs are the funcs to compile. They can call each other (and themselves)
	// no matter what their order is.
	Funcs []*ast.FuncDecl
}

type UnitCompilationResult struct {
	Funcs []Func
}

// CompileUnit compiles all unit declarations at once.
// See CompilationUnit doc comment for more info.
func CompileUnit(ctx *CompileContext, unit *CompilationUnit) (UnitCompilationResult, error) {
	internalCtx := ctx.internal()
	compiler := qcompile.NewCompiler()
	result := UnitCompilationResult{}
	err := compiler.CompileUnit(&internalCtx, qcompile.UnitCompilationConfig{
		Funcs: unit.Funcs,
		OnFuncCompiled: func(compiled *qruntime.Func) {
			result.Funcs = append(result.Funcs, Func{data: compiled})
		},
	})
	return result, err
}

// ArgsPusher binds the call arguments, one by one.
// Every Push method writes the next argument, starting from the first one:
//
// args := evalEnv.PrepareArgs()
// args.PushString("hello")
// args.PushInt(10)
// args.PushSlice(xs)
// quasigo.Call(evalEnv, fn)
//
// The value types should match the called function params.
// When calling a function with N arguments, exactly N arguments
// should be pushed.
//
// Bound args can be used for many calls if you don't need to change them.
//
// ArgsPusher is efficient, zero alloc and reflection-free.
// You can build a more convenient args marshaling on top of it if you want to.
// You can either use a reflection approach and forward to it or
// use code generation that would insert the boilerplate for you.
type ArgsPusher = qruntime.ArgsPusher

// PrepareArgs returns a pusher that starts from the first call argument.
// See [ArgsPusher] doc comment for more info.
func (env *EvalEnv) PrepareArgs() ArgsPusher {
	return env.data.PrepareArgs()
}

// Call invokes a given function.
//
// Before calling this function, be sure to bind arguments
// to the env using PrepareArgs builder.
func Call(env *EvalEnv, fn Func) CallResult {
	return CallResult{v: qruntime.Call(&env.data, fn.data)}
}

// CallResult is a return value of a function call.
type CallResult struct {
	v qruntime.Slot
}

func (res CallResult) StringValue() string { return res.v.String() }

func (res CallResult) IntValue() int { return res.v.Int() }

func (res CallResult) BoolValue() bool { return res.v.Bool() }

// Disasm returns the compiled function disassembly text.
// This output is not guaranteed to be stable between versions
// and should be used only for debugging purposes.
func Disasm(env *Env, fn Func) string {
	return qdisasm.Func(&env.data, fn.data)
}

// Func is a compiled function that is ready to be executed.
type Func struct {
	data *qruntime.Func
}

func (fn Func) Name() string { return fn.data.Name }

func (fn Func) IsNil() bool { return fn.data == nil }

// TODO: remove this. Only needed for testing.
func (fn Func) ScalarConstants() []uint64 { return fn.data.ScalarConstants }

// TODO: remove this. Only needed for testing.
func (fn Func) StringConstants() []string { return fn.data.StrConstants }
