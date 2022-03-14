package qruntime

type Func struct {
	StrConstants    []string
	ScalarConstants []uint64

	Codeptr *byte
	Code    []byte

	FrameSize  int
	FrameSlots byte
	NumParams  byte
	NumLocals  byte
	NumTemps   byte

	CanInline bool

	Name string
}

type FuncKey struct {
	Qualifier string
	Name      string
}

func (k FuncKey) String() string {
	if k.Qualifier != "" {
		return k.Qualifier + "." + k.Name
	}
	return k.Name
}

func Call(env *EvalEnv, fn *Func) Slot {
	eval(env, fn, env.slotbase)
	return env.result
}
