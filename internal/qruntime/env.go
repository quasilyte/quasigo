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
	vararg  []interface{}
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
	if numSlots < 4 {
		panic("stack size is too small")
	}
	slots := make([]Slot, numSlots)
	ee.nativeFuncs = env.NativeFuncs
	ee.userFuncs = env.UserFuncs
	ee.slots = slots
	ee.slotbase = &slots[0]
	ee.slotend = &slots[len(slots)-1]
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
	id := env.NameToFuncID[FuncKey{Qualifier: pkgPath, Name: funcName}]
	return env.UserFuncs[id]
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

func (env *EvalEnv) BindArgs(args ...interface{}) {
	for i, arg := range args {
		switch arg := arg.(type) {
		case int:
			env.slots[i].SetInt(arg)
		case bool:
			env.slots[i].SetBool(arg)
		case string:
			env.slots[i].SetString(arg)
		case []byte:
			env.slots[i].SetByteSlice(arg)
		default:
			env.slots[i].SetInterface(arg)
		}
	}
}
