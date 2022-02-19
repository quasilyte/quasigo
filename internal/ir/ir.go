package ir

import (
	"github.com/quasilyte/quasigo/internal/bytecode"
)

type PseudoOp byte

const (
	OpUnset PseudoOp = iota
	OpLabel
)

type Inst struct {
	Op     bytecode.Op
	Pseudo PseudoOp
	Value  uint16
	Arg0   uint8
	Arg1   uint8
	Arg2   uint8
	Arg3   uint8
}

func (inst Inst) WalkArgs(fn func(arg bytecode.Argument, value int)) {
	argIndex := 0
	for _, argInfo := range inst.Op.Args() {
		var v int
		switch argInfo.Kind {
		case bytecode.ArgSlot, bytecode.ArgScalarConst, bytecode.ArgStrConst:
			v = int(inst.ArgByIndex(argIndex))
			argIndex++
		case bytecode.ArgOffset:
			v = int(inst.Arg0)
			argIndex++
		default:
			v = int(inst.Value)
		}
		fn(argInfo, v)
	}
}

func (inst Inst) ArgByIndex(i int) uint8 {
	switch i {
	case 0:
		return inst.Arg0
	case 1:
		return inst.Arg1
	case 2:
		return inst.Arg2
	default:
		return inst.Arg3
	}
}

func (inst Inst) IsPseudo() bool {
	return inst.Pseudo != OpUnset
}

type Func struct {
	Code          []Inst
	NumParams     int
	NumLocals     int
	NumFrameSlots int
}
