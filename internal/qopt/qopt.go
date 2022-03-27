package qopt

import (
	"fmt"

	"github.com/quasilyte/quasigo/internal/bytecode"
	"github.com/quasilyte/quasigo/internal/ir"
)

// TODO:
//   Len temp0 = xs
//   Move arg0 = temp0
// =>
//   Len arg0 = xs

type Optimizer struct {
	fn *ir.Func

	idMap idMap
	idSet idSet

	varKillTable []varKillState
	blocksCache  []ir.Block
}

func (opt *Optimizer) PrepareFunc(fn *ir.Func) {
	opt.fn = fn

	numTemps := fn.NumFrameSlots - (opt.fn.NumParams + opt.fn.NumLocals)
	if cap(opt.varKillTable) >= numTemps {
		opt.varKillTable = opt.varKillTable[:numTemps]
	} else {
		opt.varKillTable = make([]varKillState, numTemps)
	}

	// Create basic blocks.
	fn.Blocks = createBlocks(fn, opt.blocksCache[:0])
	for i := range fn.Blocks {
		b := &fn.Blocks[i]
		if b.NumVarKill != 0 {
			opt.markUniq(b)
		}
	}
}

type varKillState struct {
	pos      uint16
	readPos  uint16
	readArg  uint8
	numReads uint8
	slotID   uint8
}

func (opt *Optimizer) markUniq(b *ir.Block) {
	if len(b.Code) == 0 {
		return
	}

	// table is preallocated, but we need to clear it.
	// Go compiler recognizes this loop and inserts memclearNoHeapPointers here.
	table := opt.varKillTable
	for i := range table {
		table[i] = varKillState{}
	}

	debugPrint := func(msg string) {
		if opt.fn.Name == "irtest.testopt" {
			print(msg)
		}
	}

	for i := len(b.Code) - 1; i >= 0; i-- {
		inst := b.Code[i]

		// Starting to track a variable.
		// If we were tracking some other variable before, it'll be overwritten.
		if inst.Pseudo == ir.OpVarKill {
			slotID := inst.Arg0.ToSlot().ID
			table[slotID] = varKillState{
				pos:    uint16(i),
				slotID: slotID,
			}
			debugPrint(fmt.Sprintf("block%d: track temp%d\n", i, slotID))
			continue
		}

		for argIndex, argInfo := range inst.Op.Args() {
			if argInfo.Kind != bytecode.ArgSlot {
				continue
			}
			slot := inst.GetArg(argIndex).ToSlot()
			if !slot.IsTemp() || table[slot.ID].pos == 0 {
				continue
			}
			if argInfo.IsWriteSlot() {
				info := table[slot.ID]
				switch info.numReads {
				case 0:
					// No reads for this var. Can safely remove it.
					b.Code[info.pos].Pseudo = ir.OpUnset
					b.Code[i].Op = bytecode.OpInvalid
					b.NumVarKill--
					debugPrint(fmt.Sprintf("block%d: clear temp%d\n", i, slot.ID))
				case 1:
					// Exactly one read for this var. Mark slot as uniq.
					b.Code[info.pos].Pseudo = ir.OpUnset
					b.Code[info.readPos].SetArgSlotKind(int(info.readArg), ir.SlotUniq)
					b.Code[i].SetArgSlotKind(0, ir.SlotUniq)
					b.NumVarKill--
					debugPrint(fmt.Sprintf("block%d: mark uniq temp%d\n", i, slot.ID))
				default:
					debugPrint(fmt.Sprintf("block%d: ignore temp%d\n", i, slot.ID))
				}
				table[slot.ID].pos = 0
			} else {
				table[slot.ID].readArg = uint8(argIndex)
				table[slot.ID].readPos = uint16(i)
				table[slot.ID].numReads++
				debugPrint(fmt.Sprintf("block%d: add read to temp%d\n", i, slot.ID))
			}
		}
	}
}

func (opt *Optimizer) OptimizePrepared() {
	changed := false
	opt.walkBlocks(func(block []ir.Inst) {
		if opt.injectConstants(block) {
			block = opt.filterBlock(block)
			changed = true
		}
		if opt.condInvert(block) {
			block = opt.filterBlock(block)
			changed = true
		}
		if opt.zeroComparisons(block) {
			changed = true
		}
	})

	if !changed {
		return
	}

	maxTempID := -1
	opt.walkBlocks(func(block []ir.Inst) {
		usedSlots := opt.idSet
		usedSlots.Reset()
		for i := len(block) - 1; i >= 0; i-- {
			inst := &block[i]
			if inst.Op.HasDst() && inst.Arg0.ToSlot().IsUniq() {
				if !usedSlots.Contains(opt.fn.SlotIndex(inst.Arg0.ToSlot())) {
					inst.Op = bytecode.OpInvalid
					continue
				}
			}
			for i, argInfo := range inst.Op.Args() {
				if argInfo.Kind != bytecode.ArgSlot {
					continue
				}
				slot := inst.GetArg(i).ToSlot()
				if slot.IsCallArg() {
					continue
				}
				if slot.IsTemp() || slot.IsUniq() {
					if int(slot.ID) > maxTempID {
						maxTempID = int(slot.ID)
					}
				}
				if argInfo.IsWriteSlot() {
					usedSlots.Remove(opt.fn.SlotIndex(slot))
				} else {
					usedSlots.Add(opt.fn.SlotIndex(slot))
				}
			}
		}
	})
	if maxTempID != -1 {
		opt.fn.NumFrameSlots = opt.fn.NumParams + opt.fn.NumLocals + maxTempID + 1
	} else {
		opt.fn.NumFrameSlots = opt.fn.NumParams + opt.fn.NumLocals
	}
}

func (opt *Optimizer) filterBlock(block []ir.Inst) []ir.Inst {
	filtered := block[:0]
	for _, inst := range block {
		if inst.IsPseudo() || inst.Op != bytecode.OpInvalid {
			filtered = append(filtered, inst)
		}
	}
	tail := block[len(filtered):]
	for i := range tail {
		tail[i].Op = bytecode.OpInvalid
	}
	for i := len(filtered); i < len(block); i++ {
		block[i].Op = bytecode.OpInvalid
	}
	return filtered
}

func (opt *Optimizer) walkBlocks(visit func([]ir.Inst)) {
	code := opt.fn.Code
	blockStart := 0
	for i, inst := range code {
		if inst.Pseudo == ir.OpLabel {
			block := code[blockStart:i]
			if len(block) != 0 {
				visit(block)
			}
			blockStart = i + 1
			continue
		}

		switch inst.Op {
		case bytecode.OpJump, bytecode.OpJumpZero, bytecode.OpJumpNotZero:
			fallthrough
		case bytecode.OpReturnZero, bytecode.OpReturnOne, bytecode.OpReturnVoid:
			fallthrough
		case bytecode.OpReturnScalar, bytecode.OpReturnStr:
			block := code[blockStart : i+1]
			if len(block) != 0 {
				visit(block)
			}
			blockStart = i + 1
		}
	}
}

func (opt *Optimizer) condInvert(block []ir.Inst) bool {
	// Not temp0 = temp1
	// JumpZero L0 temp0
	// =>
	// JumpNotZero L0 temp1

	if len(block) < 2 {
		return false
	}
	switch block[len(block)-1].Op {
	case bytecode.OpJumpZero, bytecode.OpJumpNotZero:
		// OK.
	default:
		return false
	}
	jump := block[len(block)-1]
	jumpSlot := jump.Arg1.ToSlot()
	if !jumpSlot.IsUniq() {
		return false
	}
	if block[len(block)-2].Op != bytecode.OpNot {
		return false
	}
	not := block[len(block)-2]
	if not.Arg0.ToSlot() != jumpSlot {
		return false
	}

	block[len(block)-1].Arg1 = not.Arg1
	switch jump.Op {
	case bytecode.OpJumpZero:
		block[len(block)-1].Op = bytecode.OpJumpNotZero
	case bytecode.OpJumpNotZero:
		block[len(block)-1].Op = bytecode.OpJumpZero
	}
	block[len(block)-2].Op = bytecode.OpInvalid

	return true
}

func (opt *Optimizer) zeroComparisons(block []ir.Inst) bool {
	if len(block) < 3 {
		return false
	}

	// x != 0
	//
	// Zero temp1
	// ScalarNotEq temp0 = x temp1
	// JumpZero L0 temp0
	// =>
	// JumpZero L0 x
	//
	// temp0 must be uniq

	switch block[len(block)-1].Op {
	case bytecode.OpJumpZero, bytecode.OpJumpNotZero:
		// OK.
	default:
		return false
	}
	jump := block[len(block)-1]
	jumpSlot := jump.Arg1.ToSlot()
	if !jumpSlot.IsUniq() {
		return false
	}
	switch block[len(block)-2].Op {
	case bytecode.OpScalarEq, bytecode.OpScalarNotEq:
		// OK.
	default:
		return false
	}

	cmp := block[len(block)-2]
	if block[len(block)-3].Op != bytecode.OpZero {
		return false
	}
	cmpDst := cmp.Arg0.ToSlot()
	if !cmpDst.IsUniq() || cmpDst.ID != jumpSlot.ID {
		return false
	}
	zero := block[len(block)-3]
	if cmp.Arg2 != zero.Arg0 {
		return false
	}

	combinedOp := bytecode.OpInvalid
	switch {
	case cmp.Op == bytecode.OpScalarEq && jump.Op == bytecode.OpJumpZero:
		combinedOp = bytecode.OpJumpNotZero
	case cmp.Op == bytecode.OpScalarNotEq && jump.Op == bytecode.OpJumpZero:
		combinedOp = bytecode.OpJumpZero
	case cmp.Op == bytecode.OpScalarEq && jump.Op == bytecode.OpJumpNotZero:
		combinedOp = bytecode.OpJumpZero
	case cmp.Op == bytecode.OpScalarNotEq && jump.Op == bytecode.OpJumpNotZero:
		combinedOp = bytecode.OpJumpNotZero
	}

	if combinedOp != bytecode.OpInvalid {
		block[len(block)-1].Op = combinedOp
		block[len(block)-1].Arg1 = cmp.Arg1
		block[len(block)-2].Op = bytecode.OpInvalid
		block[len(block)-3].Op = bytecode.OpInvalid
		return true
	}

	return false
}

func (opt *Optimizer) injectConstants(block []ir.Inst) bool {
	if len(block) > 255 {
		return false
	}

	changed := false
	tracked := opt.idMap
	tracked.Reset()
	for i := len(block) - 1; i > 0; i-- {
		inst := block[i]
		storeHandled := false
		switch inst.Op {
		case bytecode.OpLoadStrConst:
			dstslot := inst.Arg0.ToSlot()
			if !dstslot.IsUniq() {
				continue
			}
			key := tracked.FindIndex(dstslot.ID)
			if key == -1 {
				continue
			}
			j := tracked.GetValue(key)
			if block[j].Op == bytecode.OpMove {
				block[j].Op = bytecode.OpLoadStrConst
				block[j].Arg1 = inst.Arg1
				block[i].Op = bytecode.OpInvalid
				changed = true
			}
			tracked.RemoveAt(key)
			storeHandled = true

		case bytecode.OpLoadScalarConst:
			dstslot := inst.Arg0.ToSlot()
			if !dstslot.IsUniq() {
				continue
			}
			key := tracked.FindIndex(dstslot.ID)
			if key == -1 {
				continue
			}
			j := tracked.GetValue(key)
			if block[j].Op == bytecode.OpMove {
				block[j].Op = bytecode.OpLoadScalarConst
				block[j].Arg1 = inst.Arg1
				block[i].Op = bytecode.OpInvalid
				changed = true
			}
			tracked.RemoveAt(key)
			storeHandled = true

		case bytecode.OpMove:
			dstslot := inst.Arg0.ToSlot()
			if dstslot.IsUniq() {
				break // handled below
			}
			if !dstslot.IsCallArg() {
				continue
			}
			srcslot := inst.Arg1.ToSlot()
			if !srcslot.IsUniq() {
				continue
			}
			tracked.Add(srcslot.ID, uint8(i))
		}

		if inst.Op.HasDst() && !storeHandled {
			dstslot := inst.Arg0.ToSlot()
			if dstslot.IsUniq() {
				tracked.Remove(dstslot.ID)
			}
		}
	}

	return changed
}
