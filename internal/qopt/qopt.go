package qopt

import (
	"github.com/quasilyte/quasigo/internal/bytecode"
	"github.com/quasilyte/quasigo/internal/ir"
)

func Func(fn *ir.Func) {
	opt := optimizer{fn: fn}
	opt.Optimize()
}

type optimizer struct {
	fn *ir.Func

	idMap idMap
	idSet idSet
}

func (opt *optimizer) Optimize() {
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

func (opt *optimizer) filterBlock(block []ir.Inst) []ir.Inst {
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

func (opt *optimizer) walkBlocks(visit func([]ir.Inst)) {
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
		case bytecode.OpReturnFalse, bytecode.OpReturnTrue, bytecode.OpReturnVoid:
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

func (opt *optimizer) condInvert(block []ir.Inst) bool {
	// Not tmp0 = tmp1
	// JumpZero L0 tmp0
	// =>
	// JumpNotZero L0 tmp1

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

func (opt *optimizer) zeroComparisons(block []ir.Inst) bool {
	if len(block) < 3 {
		return false
	}

	// x != 0
	//
	// LoadScalarConst tmp1 = 0
	// ScalarNotEq tmp0 = x tmp1
	// JumpZero L0 tmp0
	// =>
	// JumpZero L0 x
	//
	// tmp0 must be uniq

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
	if block[len(block)-3].Op != bytecode.OpLoadScalarConst {
		return false
	}
	cmpDst := cmp.Arg0.ToSlot()
	if !cmpDst.IsUniq() || cmpDst.ID != jumpSlot.ID {
		return false
	}
	load := block[len(block)-3]
	if opt.fn.ScalarConstants[load.Arg1] != 0 {
		return false
	}
	if cmp.Arg2 != load.Arg0 {
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

func (opt *optimizer) injectConstants(block []ir.Inst) bool {
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
			if block[j].Op == bytecode.OpMoveStr {
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
			if block[j].Op == bytecode.OpMoveScalar {
				block[j].Op = bytecode.OpLoadScalarConst
				block[j].Arg1 = inst.Arg1
				block[i].Op = bytecode.OpInvalid
				changed = true
			}
			tracked.RemoveAt(key)
			storeHandled = true

		case bytecode.OpMoveScalar, bytecode.OpMoveStr:
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
