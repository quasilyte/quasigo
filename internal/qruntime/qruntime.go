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

type FuncDebugInfo struct {
	SlotNames []string
	NumLocals int
}
