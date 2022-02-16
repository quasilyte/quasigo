package quasigo

import "github.com/quasilyte/quasigo/internal/qruntime"

type debugInfo struct {
	funcs map[*qruntime.Func]qruntime.FuncDebugInfo
}

func newDebugInfo() *debugInfo {
	return &debugInfo{
		funcs: make(map[*qruntime.Func]qruntime.FuncDebugInfo),
	}
}
