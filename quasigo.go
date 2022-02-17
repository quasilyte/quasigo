// Package quasigo implements a Go subset compiler and interpreter.
//
// The implementation details are not part of the contract of this package.
package quasigo

import (
	"go/ast"
	"go/token"
	"go/types"

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
// It method is defined only on pointer type, the typeName should start with `*`.
func (env *Env) AddNativeMethod(typeName, methodName string, f func(qnative.CallContext)) {
	env.data.AddNativeFunc(typeName, methodName, f)
}

// AddNativeFunc binds `$pkgPath.$funcName` symbol with f.
// A pkgPath should be a full package path in which funcName is defined.
func (env *Env) AddNativeFunc(pkgPath, funcName string, f func(qnative.CallContext)) {
	env.data.AddNativeFunc(pkgPath, funcName, f)
}

// AddFunc binds `$pkgPath.$funcName` symbol with f.
func (env *Env) AddFunc(pkgPath, funcName string, f Func) {
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
	Fset    *token.FileSet
}

// Compile prepares an executable version of fn.
func Compile(ctx *CompileContext, fn *ast.FuncDecl) (Func, error) {
	internalCtx := qcompile.Context{
		Env:     &ctx.Env.data,
		Package: ctx.Package,
		Types:   ctx.Types,
		Fset:    ctx.Fset,
	}
	compiled, err := qcompile.Func(&internalCtx, fn)
	return Func{data: compiled}, err
}

// BindArgs prepares the arguments for the call.
// Bound args can be used many times if you don't need to change
// the call arguments.
//
// If BindArgs+Call is not convenient for you, consider using the
// simple wrapper that does this combination for you.
// Note, however, that reusing bound arguments, whether possible,
// if more efficient.
func (env *EvalEnv) BindArgs(args ...interface{}) {
	env.data.BindArgs(args...)
}

// Call invokes a given function.
// Before calling this function, be sure to bind arguments
// to the env using BindArgs.
func Call(env *EvalEnv, fn Func) CallResult {
	return CallResult{v: qruntime.Call(&env.data, fn.data)}
}

// CallResult is a return value of Call function.
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

func (fn Func) IsNil() bool { return fn.data == nil }
