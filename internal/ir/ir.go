package ir

import (
	"github.com/quasilyte/quasigo/internal/bytecode"
)

//go:generate stringer -type=PseudoOp -trimprefix=Op
type PseudoOp byte

const (
	OpUnset PseudoOp = iota
	OpLabel
)

type InstArg uint16

func (a InstArg) ToSlot() Slot {
	return Slot{
		ID:   uint8(a >> 8),
		Kind: SlotKind(a),
	}
}

type Inst struct {
	Op     bytecode.Op
	Pseudo PseudoOp
	Arg0   InstArg
	Arg1   InstArg
	Arg2   InstArg
	Arg3   InstArg
}

func (inst *Inst) SetArg(i int, arg InstArg) {
	switch i {
	case 0:
		inst.Arg0 = arg
	case 1:
		inst.Arg1 = arg
	case 2:
		inst.Arg2 = arg
	default:
		inst.Arg3 = arg
	}
}

func (inst Inst) GetArg(i int) InstArg {
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

	StrConstants    []string
	ScalarConstants []uint64
}

func (fn *Func) SlotIndex(slot Slot) uint8 {
	switch slot.Kind {
	case SlotCallArg:
		return uint8(fn.NumFrameSlots) + slot.ID
	case SlotLocal:
		return uint8(fn.NumParams) + slot.ID
	case SlotTemp, SlotUniq:
		return uint8(fn.NumParams) + uint8(fn.NumLocals) + slot.ID
	default:
		return slot.ID
	}
}
