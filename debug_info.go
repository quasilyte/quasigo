package quasigo

type debugInfo struct {
	funcs map[*Func]funcDebugInfo
}

type funcDebugInfo struct {
	slotNames []string
	numLocals int
}

func newDebugInfo() *debugInfo {
	return &debugInfo{
		funcs: make(map[*Func]funcDebugInfo),
	}
}
