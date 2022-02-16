package quasigo

import "github.com/quasilyte/quasigo/internal/qruntime"

type debugInfo struct {
	funcs map[*qruntime.Func]funcDebugInfo
}

type funcDebugInfo struct {
	slotNames []string
	numLocals int
}

func newDebugInfo() *debugInfo {
	return &debugInfo{
		funcs: make(map[*qruntime.Func]funcDebugInfo),
	}
}
