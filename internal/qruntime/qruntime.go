package qruntime

type Func struct {
	StrConstants    []string
	ScalarConstants []uint64

	Codeptr *byte
	Code    []byte

	FrameSize  int
	FrameSlots byte

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

type FuncDebugInfo struct {
	SlotNames []string
	NumLocals int
}
