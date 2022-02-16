package quasigo

import "github.com/quasilyte/quasigo/internal/qruntime"

type nativeFunc struct {
	mappedFunc func(NativeCallContext)
	name       string // Needed for the readable disasm
	frameSize  int
}

func newEnv() *Env {
	return &Env{
		nameToNativeFuncID: make(map[qruntime.FuncKey]uint16),
		nameToFuncID:       make(map[qruntime.FuncKey]uint16),

		debug: qruntime.NewDebugInfo(),
	}
}

func (env *Env) addNativeFunc(key qruntime.FuncKey, f func(NativeCallContext)) {
	id := len(env.nativeFuncs)
	env.nativeFuncs = append(env.nativeFuncs, nativeFunc{
		mappedFunc: f,
		name:       key.String(),
		frameSize:  int(qruntime.SizeofSlot) * maxNativeFuncArgs,
	})
	env.nameToNativeFuncID[key] = uint16(id)
}

func (env *Env) addFunc(key qruntime.FuncKey, f *qruntime.Func) {
	id := len(env.userFuncs)
	env.userFuncs = append(env.userFuncs, f)
	env.nameToFuncID[key] = uint16(id)
}
