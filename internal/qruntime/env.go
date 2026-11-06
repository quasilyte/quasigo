package qruntime

const MaxNativeFuncArgs = 6

// Env is used to hold both compilation and evaluation data.
type Env struct {
	// TODO(quasilyte): store both native and user func ids in one map?

	NativeFuncs        []NativeFunc
	NameToNativeFuncID map[FuncKey]uint16

	UserFuncs    []*Func
	NameToFuncID map[FuncKey]uint16

	// debug contains all information that is only needed
	// for better debugging and compiled code introspection.
	// Right now it's always enabled, but we may allow stripping it later.
	Debug *DebugInfo
}

// EvalEnv is a goroutine-local handle for Env.
// To get one, use Env.GetEvalEnv() method.
type EvalEnv struct {
	nativeFuncs []NativeFunc
	userFuncs   []*Func

	slots    []Slot
	slotbase *Slot
	slotend  *Slot

	result  Slot
	result2 Slot
	vararg  []any
}

func InitEnv(env *Env) {
	env.NameToNativeFuncID = make(map[FuncKey]uint16)
	env.NameToFuncID = make(map[FuncKey]uint16)
	env.Debug = NewDebugInfo()

	env.AddNativeFunc("builtin", "makeSlice", nativeMakeSlice)
	env.AddNativeFunc("builtin", "append8", nativeAppend8)
	env.AddNativeFunc("builtin", "append64", nativeAppend64)
	env.AddNativeFunc("builtin", "bytesToString", nativeBytesToString)
}

func InitEvalEnv(env *Env, ee *EvalEnv, stackSize int) {
	numSlots := stackSize / int(SizeofSlot)
	if numSlots < 16+MaxNativeFuncArgs {
		panic("stack size is too small")
	}
	ee.nativeFuncs = env.NativeFuncs
	ee.userFuncs = env.UserFuncs

	// This covers both script slots plus some slots
	// available to call a native func at the end of the stack.
	slots := make([]Slot, numSlots)

	ee.slots = slots
	ee.slotbase = &slots[0]

	// Scripts can use slots up until the space reserved for
	// the pre-overflow native call. It can be treated as
	// there is some safe space for args beyond the normal stack space.
	numScriptSlots := numSlots - MaxNativeFuncArgs - 1
	ee.slotend = &slots[numScriptSlots]
}

type NativeFunc struct {
	mappedFunc func(NativeCallContext)
	Name       string // Needed for the readable disasm
	frameSize  int
}

func (env *Env) AddNativeMethod(typeName, methodName string, f func(NativeCallContext)) {
	env.addNativeFunc(FuncKey{Qualifier: typeName, Name: methodName}, f)
}

func (env *Env) AddNativeFunc(pkgPath, funcName string, f func(NativeCallContext)) {
	env.addNativeFunc(FuncKey{Qualifier: pkgPath, Name: funcName}, f)
}

func (env *Env) AddFunc(pkgPath, funcName string, f *Func) {
	env.addFunc(FuncKey{Qualifier: pkgPath, Name: funcName}, f)
}

func (env *Env) GetFunc(pkgPath, funcName string) *Func {
	id, ok := env.NameToFuncID[FuncKey{Qualifier: pkgPath, Name: funcName}]
	if ok {
		return env.UserFuncs[id]
	}
	return nil
}

func (env *Env) addNativeFunc(key FuncKey, f func(NativeCallContext)) {
	id := len(env.NativeFuncs)
	env.NativeFuncs = append(env.NativeFuncs, NativeFunc{
		mappedFunc: f,
		Name:       key.String(),
		frameSize:  int(SizeofSlot) * MaxNativeFuncArgs,
	})
	env.NameToNativeFuncID[key] = uint16(id)
}

func (env *Env) addFunc(key FuncKey, f *Func) {
	id := len(env.UserFuncs)
	env.UserFuncs = append(env.UserFuncs, f)
	env.NameToFuncID[key] = uint16(id)
}

// ArgsPusher binds the call arguments, one by one.
type ArgsPusher struct {
	slots []Slot
	n     int
}

func (env *EvalEnv) PrepareArgs() ArgsPusher {
	return ArgsPusher{slots: env.slots}
}

func (p *ArgsPusher) next() *Slot {
	slot := &p.slots[p.n]
	p.n++
	return slot
}

func (p *ArgsPusher) PushBool(v bool)     { p.next().SetBool(v) }
func (p *ArgsPusher) PushByte(v byte)     { p.next().SetByte(v) }
func (p *ArgsPusher) PushInt(v int)       { p.next().SetInt(v) }
func (p *ArgsPusher) PushFloat(v float64) { p.next().SetFloat(v) }
func (p *ArgsPusher) PushString(v string) { p.next().SetString(v) }
func (p *ArgsPusher) PushInterface(v any) { p.next().SetInterface(v) }

func (p *ArgsPusher) PushSlice[T any](v []T) { *p.next().Slice[T]() = v }
func (p *ArgsPusher) PushPtr[T any](v *T)    { p.next().SetPtr(v) }
