package qruntime

type DebugInfo struct {
	Funcs map[*Func]FuncDebugInfo
}

func NewDebugInfo() *DebugInfo {
	return &DebugInfo{
		Funcs: make(map[*Func]FuncDebugInfo),
	}
}

type FuncDebugInfo struct {
	SlotNames []string
	NumLocals int
}
