package quasigo

import "github.com/quasilyte/quasigo/internal/qruntime"

type funcKey struct {
	qualifier string
	name      string
}

func (k funcKey) String() string {
	if k.qualifier != "" {
		return k.qualifier + "." + k.name
	}
	return k.name
}

type nativeFunc struct {
	mappedFunc func(NativeCallContext)
	name       string // Needed for the readable disasm
	frameSize  int
}

func newEnv() *Env {
	return &Env{
		nameToNativeFuncID: make(map[funcKey]uint16),
		nameToFuncID:       make(map[funcKey]uint16),

		debug: newDebugInfo(),
	}
}

func (env *Env) addNativeFunc(key funcKey, f func(NativeCallContext)) {
	id := len(env.nativeFuncs)
	env.nativeFuncs = append(env.nativeFuncs, nativeFunc{
		mappedFunc: f,
		name:       key.String(),
		frameSize:  int(qruntime.SizeofSlot) * maxNativeFuncArgs,
	})
	env.nameToNativeFuncID[key] = uint16(id)
}

func (env *Env) addFunc(key funcKey, f *qruntime.Func) {
	id := len(env.userFuncs)
	env.userFuncs = append(env.userFuncs, f)
	env.nameToFuncID[key] = uint16(id)
}
