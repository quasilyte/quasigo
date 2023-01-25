package ir

import (
	"github.com/quasilyte/quasigo/internal/bytecode"
	"github.com/quasilyte/quasigo/internal/qruntime"
)

//go:generate stringer -type=PseudoOp -trimprefix=Op
type PseudoOp byte

const (
	OpUnset PseudoOp = iota
	OpLabel
	OpVarKill
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

func (inst Inst) IsEmpty() bool {
	return inst.Op == bytecode.OpInvalid && inst.Pseudo == OpUnset
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

func (inst *Inst) SetArgSlotKind(i int, kind SlotKind) {
	slot := inst.GetArg(i).ToSlot()
	slot.Kind = kind
	inst.SetArg(i, slot.ToInstArg())
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
	Name string

	Code      []Inst
	Blocks    []Block
	NumParams int
	NumLocals int
	NumTemps  int

	StrConstants    []string
	ScalarConstants []uint64

	Debug qruntime.FuncDebugInfo
	Env   *qruntime.Env
}

func (fn *Func) NewScalarConstant(v uint64) int {
	for i := range fn.ScalarConstants {
		if fn.ScalarConstants[i] == v {
			return i
		}
	}
	fn.ScalarConstants = append(fn.ScalarConstants, v)
	return len(fn.ScalarConstants) - 1
}

func (fn *Func) NumFrameSlots() int {
	return fn.NumParams + fn.NumLocals + fn.NumTemps
}

func (fn *Func) SlotIndex(slot Slot) uint8 {
	switch slot.Kind {
	case SlotCallArg:
		return uint8(fn.NumFrameSlots()) + slot.ID
	case SlotTemp, SlotUniq:
		return uint8(fn.NumParams) + slot.ID
	default:
		return slot.ID
	}
}

type Block struct {
	Code       []Inst
	NumVarKill uint16
	Label      uint16
	Dirty      bool
}

func (b *Block) HasLabel() bool { return b.Label != 0 }

func (b *Block) LabelID() uint16 { return b.Label - 1 }
