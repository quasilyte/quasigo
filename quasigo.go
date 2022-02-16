// Package quasigo implements a Go subset compiler and interpreter.
//
// The implementation details are not part of the contract of this package.
package quasigo

import (
	"go/ast"
	"go/token"
	"go/types"
)

// TODO(quasilyte): document what is thread-safe and what not.
// TODO(quasilyte): add a readme.

// Env is used to hold both compilation and evaluation data.
type Env struct {
	// TODO(quasilyte): store both native and user func ids in one map?

	nativeFuncs        []nativeFunc
	nameToNativeFuncID map[funcKey]uint16

	userFuncs    []*Func
	nameToFuncID map[funcKey]uint16

	// debug contains all information that is only needed
	// for better debugging and compiled code introspection.
	// Right now it's always enabled, but we may allow stripping it later.
	debug *debugInfo
}

// EvalEnv is a goroutine-local handle for Env.
// To get one, use Env.GetEvalEnv() method.
type EvalEnv struct {
	nativeFuncs []nativeFunc
	userFuncs   []*Func

	slots    []slotValue
	slotbase *slotValue
	slotend  *slotValue

	result  slotValue
	result2 slotValue
	vararg  []interface{}
}

type NativeCallContext struct {
	env     *EvalEnv
	slotptr *slotValue
}

func (ncc NativeCallContext) BoolArg(index int) bool {
	return getslot(ncc.slotptr, byte(index)).Bool()
}

func (ncc NativeCallContext) IntArg(index int) int {
	return getslot(ncc.slotptr, byte(index)).Int()
}

func (ncc NativeCallContext) StringArg(index int) string {
	return getslot(ncc.slotptr, byte(index)).String()
}

func (ncc NativeCallContext) InterfaceArg(index int) interface{} {
	return getslot(ncc.slotptr, byte(index)).Interface()
}

func (ncc NativeCallContext) VariadicArg() []interface{} {
	return ncc.env.vararg
}

func (ncc NativeCallContext) SetIntResult(v int)  { ncc.env.result.SetInt(v) }
func (ncc NativeCallContext) SetIntResult2(v int) { ncc.env.result2.SetInt(v) }

func (ncc NativeCallContext) SetBoolResult(v bool)  { ncc.env.result.SetBool(v) }
func (ncc NativeCallContext) SetBoolResult2(v bool) { ncc.env.result2.SetBool(v) }

func (ncc NativeCallContext) SetStringResult(v string)  { ncc.env.result.SetString(v) }
func (ncc NativeCallContext) SetStringResult2(v string) { ncc.env.result2.SetString(v) }

func (ncc NativeCallContext) SetInterfaceResult(v interface{})  { ncc.env.result.SetInterface(v) }
func (ncc NativeCallContext) SetInterfaceResult2(v interface{}) { ncc.env.result2.SetInterface(v) }

// NewEnv creates a new empty environment.
func NewEnv() *Env {
	return newEnv()
}

// GetEvalEnv creates a new goroutine-local handle of env.
// Stack size is amount of bytes we allocate for all stack
// frames of this env.
func (env *Env) GetEvalEnv(stackSize int) *EvalEnv {
	numSlots := stackSize / int(sizeofSlotValue)
	if numSlots < 4 {
		panic("stack size is too small")
	}
	slots := make([]slotValue, numSlots)
	return &EvalEnv{
		nativeFuncs: env.nativeFuncs,
		userFuncs:   env.userFuncs,
		slots:       slots,
		slotbase:    &slots[0],
		slotend:     &slots[len(slots)-1],
	}
}

// AddNativeMethod binds `$typeName.$methodName` symbol with f.
// A typeName should be fully qualified, like `github.com/user/pkgname.TypeName`.
// It method is defined only on pointer type, the typeName should start with `*`.
func (env *Env) AddNativeMethod(typeName, methodName string, f func(NativeCallContext)) {
	env.addNativeFunc(funcKey{qualifier: typeName, name: methodName}, f)
}

// AddNativeFunc binds `$pkgPath.$funcName` symbol with f.
// A pkgPath should be a full package path in which funcName is defined.
func (env *Env) AddNativeFunc(pkgPath, funcName string, f func(NativeCallContext)) {
	env.addNativeFunc(funcKey{qualifier: pkgPath, name: funcName}, f)
}

// AddFunc binds `$pkgPath.$funcName` symbol with f.
func (env *Env) AddFunc(pkgPath, funcName string, f *Func) {
	env.addFunc(funcKey{qualifier: pkgPath, name: funcName}, f)
}

// GetFunc finds previously bound function searching for the `$pkgPath.$funcName` symbol.
func (env *Env) GetFunc(pkgPath, funcName string) *Func {
	id := env.nameToFuncID[funcKey{qualifier: pkgPath, name: funcName}]
	return env.userFuncs[id]
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
func Compile(ctx *CompileContext, fn *ast.FuncDecl) (compiled *Func, err error) {
	return compile(ctx, fn)
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
	for i, arg := range args {
		switch arg := arg.(type) {
		case int:
			env.slots[i].SetInt(arg)
		case bool:
			env.slots[i].SetBool(arg)
		case string:
			env.slots[i].SetString(arg)
		default:
			env.slots[i].SetInterface(arg)
		}
	}
}

// Call invokes a given function.
// Before calling this function, be sure to bind arguments
// to the env using BindArgs.
func Call(env *EvalEnv, fn *Func) CallResult {
	eval(env, fn, env.slotbase)
	return CallResult{v: env.result}
}

// CallResult is a return value of Call function.
type CallResult struct {
	v slotValue
}

func (res CallResult) StringValue() string { return res.v.String() }

func (res CallResult) IntValue() int { return res.v.Int() }

func (res CallResult) BoolValue() bool { return res.v.Bool() }

// Disasm returns the compiled function disassembly text.
// This output is not guaranteed to be stable between versions
// and should be used only for debugging purposes.
func Disasm(env *Env, fn *Func) string {
	return disasm(env, fn)
}

// Func is a compiled function that is ready to be executed.
type Func struct {
	strConstants    []string
	scalarConstants []uint64

	codeptr *byte
	code    []byte

	frameSize  int
	frameSlots byte

	name string
}

// func (s *ValueStack) Arg(index uint) interface{} {
// 	i := s.base + index
// 	if i < uint(len(s.values)) {
// 		return s.values[i].Interface()
// 	}
// 	return nil
// }

// func (s *ValueStack) IntArg(index uint) int {
// 	i := s.base + index
// 	if i < uint(len(s.values)) {
// 		return s.values[i].Int()
// 	}
// 	return 0
// }

// // Push adds x to the stack.
// // Important: for int-typed values, use PushInt.
// func (s *ValueStack) Push(x interface{}) {
// 	s.values = append(s.values, slotValue{object: x})
// }

// // PushInt adds x to the stack.
// func (s *ValueStack) PushInt(x int) {
// 	s.values = append(s.values, slotValue{scalar: uint64(x)})
// }