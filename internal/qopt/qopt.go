package qopt

import (
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

	slotToIndexMap sliceMap[ir.Slot, uint16]
	idSet          idSet

	varKillTable []varKillState
	blocksCache  []ir.Block
}

func (opt *Optimizer) PrepareFunc(fn *ir.Func) {
	opt.fn = fn

	if cap(opt.varKillTable) >= fn.NumTemps {
		opt.varKillTable = opt.varKillTable[:fn.NumTemps]
	} else {
		opt.varKillTable = make([]varKillState, fn.NumTemps)
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

	numVarKills := b.NumVarKill
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
				case 1:
					// Exactly one read for this var. Mark slot as uniq.
					b.Code[info.pos].Pseudo = ir.OpUnset
					b.Code[info.readPos].SetArgSlotKind(int(info.readArg), ir.SlotUniq)
					b.Code[i].SetArgSlotKind(0, ir.SlotUniq)
					b.NumVarKill--
				default:
				}
				table[slot.ID].pos = 0
			} else {
				table[slot.ID].readArg = uint8(argIndex)
				table[slot.ID].readPos = uint16(i)
				table[slot.ID].numReads++
			}
		}
	}

	if numVarKills == b.NumVarKill {
		return // No changes
	}

	// Trim the block.
	for len(b.Code) > 0 && b.Code[len(b.Code)-1].IsEmpty() {
		b.Code = b.Code[:len(b.Code)-1]
	}
}

func (opt *Optimizer) OptimizePrepared() {
	numChanged := 0
	for numChanged < 5 {
		changed := false
		opt.walkBlocks(func(b *ir.Block) {
			if opt.injectConstants(b) {
				opt.filterBlock(b)
				b.Dirty = true
				changed = true
			}
			if opt.condInvert(b) {
				opt.filterBlock(b)
				b.Dirty = true
				changed = true
			}
			if opt.zeroComparisons(b) {
				opt.filterBlock(b)
				b.Dirty = true
				changed = true
			}
			if opt.removeDeadstores(b) {
				opt.filterBlock(b)
				b.Dirty = true
				changed = true
			}
		})
		if !changed {
			break
		}
		for i := range opt.fn.Blocks {
			b := &opt.fn.Blocks[i]
			if !b.Dirty {
				continue
			}
			if b.NumVarKill != 0 {
				opt.markUniq(b)
			}
		}
		numChanged++
	}

	changed := numChanged != 0
	if !changed {
		return
	}

	liveScalarConsts := make(map[uint64]int, len(opt.fn.ScalarConstants))
	newScalarConsts := make([]uint64, 0, len(opt.fn.ScalarConstants))
	internScalarConst := func(v uint64) int {
		if id, ok := liveScalarConsts[v]; ok {
			return id
		}
		id := len(newScalarConsts)
		newScalarConsts = append(newScalarConsts, v)
		liveScalarConsts[v] = id
		return id
	}

	// TODO: reuse this []bool?
	// TODO: use bitset?
	liveSlots := make([]bool, opt.fn.NumTemps)
	opt.walkBlocks(func(b *ir.Block) {
		usedSlots := opt.idSet
		usedSlots.Reset()
		for i := len(b.Code) - 1; i >= 0; i-- {
			inst := &b.Code[i]
			if inst.Op.HasDst() && inst.Arg0.ToSlot().IsUniq() {
				if !usedSlots.Contains(opt.fn.SlotIndex(inst.Arg0.ToSlot())) {
					inst.Op = bytecode.OpInvalid
					continue
				}
			}
			for i, argInfo := range inst.Op.Args() {
				arg := inst.GetArg(i)
				if argInfo.Kind == bytecode.ArgScalarConst {
					v := opt.fn.ScalarConstants[arg]
					inst.SetArg(i, ir.InstArg(internScalarConst(v)))
					continue
				}

				if argInfo.Kind != bytecode.ArgSlot {
					continue
				}
				slot := arg.ToSlot()
				if slot.IsCallArg() {
					continue
				}
				if slot.IsTemp() || slot.IsUniq() {
					liveSlots[slot.ID] = true
				}
				if argInfo.IsWriteSlot() {
					usedSlots.Remove(opt.fn.SlotIndex(slot))
				} else {
					usedSlots.Add(opt.fn.SlotIndex(slot))
				}
			}
		}
	})
	opt.fn.ScalarConstants = newScalarConsts

	newSlotIDs := make([]uint8, len(liveSlots))
	slotOffset := 0
	for id, isUsed := range liveSlots {
		if isUsed {
			newSlotIDs[id] = uint8(slotOffset)
			slotOffset++
		}
	}
	numTemps := slotOffset
	opt.walkBlocks(func(b *ir.Block) {
		for i := range b.Code {
			inst := &b.Code[i]
			for i, argInfo := range inst.Op.Args() {
				if argInfo.Kind != bytecode.ArgSlot {
					continue
				}
				slot := inst.GetArg(i).ToSlot()
				if !slot.IsTemp() && !slot.IsUniq() {
					continue
				}
				if !liveSlots[slot.ID] {
					continue
				}
				slot.ID = newSlotIDs[slot.ID]
				inst.SetArg(i, slot.ToInstArg())
			}
		}
	})

	opt.fn.NumTemps = numTemps
}

func (opt *Optimizer) filterBlock(b *ir.Block) {
	filtered := b.Code[:0]
	for _, inst := range b.Code {
		if inst.IsPseudo() || inst.Op != bytecode.OpInvalid {
			filtered = append(filtered, inst)
		}
	}
	tail := b.Code[len(filtered):]
	for i := range tail {
		tail[i].Op = bytecode.OpInvalid
	}
	b.Code = filtered
}

func (opt *Optimizer) walkBlocks(visit func(b *ir.Block)) {
	for i := range opt.fn.Blocks {
		visit(&opt.fn.Blocks[i])
	}
}

func (opt *Optimizer) condInvert(b *ir.Block) bool {
	// Not temp0 = temp1
	// JumpZero L0 temp0
	// =>
	// JumpNotZero L0 temp1

	block := b.Code

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

func (opt *Optimizer) zeroComparisons(b *ir.Block) bool {
	block := b.Code

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

func (opt *Optimizer) removeDeadstores(b *ir.Block) bool {
	block := b.Code

	if len(block) > 255 {
		return false
	}

	changed := false
	movedUniqs := opt.slotToIndexMap
	movedUniqs.Reset()
	for i, inst := range block {
		switch inst.Op {
		case bytecode.OpMove:
			dst := inst.Arg0.ToSlot()
			src := inst.Arg1.ToSlot()
			if !dst.IsUniq() {
				break
			}
			srckey := movedUniqs.FindIndex(src)
			if srckey != -1 {
				srcpos := movedUniqs.GetValue(srckey)
				block[i].Arg1 = block[srcpos].Arg1
				block[srcpos].Op = bytecode.OpInvalid
				changed = true
			}
			movedUniqs.Add(dst, uint16(i))

		default:
			for argIndex, argInfo := range inst.Op.Args() {
				if argInfo.Kind != bytecode.ArgSlot || !argInfo.IsReadSlot() {
					continue
				}
				arg := inst.GetArg(argIndex).ToSlot()
				if !arg.IsUniq() {
					continue
				}
				srckey := movedUniqs.FindIndex(arg)
				if srckey != -1 {
					srcpos := movedUniqs.GetValue(srckey)
					block[i].SetArg(argIndex, block[srcpos].Arg1)
					block[srcpos].Op = bytecode.OpInvalid
					changed = true
				}
			}
		}

		for argIndex, argInfo := range inst.Op.Args() {
			if argInfo.Kind != bytecode.ArgSlot {
				continue
			}
			if !argInfo.IsReadSlot() {
				continue
			}
			arg := inst.GetArg(argIndex).ToSlot()
			if !arg.IsUniq() {
				continue
			}
			movedUniqs.Remove(arg)
		}
	}

	return changed
}

func (opt *Optimizer) injectConstants(b *ir.Block) bool {
	block := b.Code

	if len(block) > 255 {
		return false
	}

	getInt64Value := func(fn *ir.Func, inst ir.Inst) int64 {
		if inst.Op == bytecode.OpZero {
			return 0
		}
		return int64(fn.ScalarConstants[inst.Arg1])
	}

	changed := false
	constValues := opt.slotToIndexMap
	constValues.Reset()
	for i, inst := range block {
		switch inst.Op {
		case bytecode.OpLoadScalarConst, bytecode.OpLoadStrConst:
			dstslot := inst.Arg0.ToSlot()
			if !dstslot.IsUniq() {
				break
			}
			constValues.Add(dstslot, uint16(i))

		case bytecode.OpZero:
			dstslot := inst.Arg0.ToSlot()
			if !dstslot.IsUniq() {
				break
			}
			constValues.Add(dstslot, uint16(i))

		case bytecode.OpIntAdd64:
			xslot := inst.Arg1.ToSlot()
			yslot := inst.Arg2.ToSlot()
			if !xslot.IsUniq() || !yslot.IsUniq() {
				break
			}
			xkey := constValues.FindIndex(xslot)
			if xkey == -1 {
				break
			}
			ykey := constValues.FindIndex(yslot)
			if ykey == -1 {
				break
			}
			xpos := constValues.GetValue(xkey)
			ypos := constValues.GetValue(ykey)
			xload := block[xpos]
			yload := block[ypos]
			xvalue := getInt64Value(opt.fn, xload)
			yvalue := getInt64Value(opt.fn, yload)
			result := xvalue + yvalue
			if result == 0 {
				block[i].Op = bytecode.OpZero
			} else {
				block[i].Op = xload.Op
				block[i].Arg1 = ir.InstArg(opt.fn.NewScalarConstant(uint64(xvalue + yvalue)))
			}
			block[xpos].Op = bytecode.OpInvalid
			block[ypos].Op = bytecode.OpInvalid
			dstslot := inst.Arg0.ToSlot()
			if dstslot.IsUniq() {
				constValues.Add(dstslot, uint16(i))
			}
			changed = true

		case bytecode.OpMove:
			srcslot := inst.Arg1.ToSlot()
			key := constValues.FindIndex(srcslot)
			if key == -1 {
				break
			}
			j := constValues.GetValue(key)
			if block[j].Op == bytecode.OpZero {
				block[i].Op = bytecode.OpZero
				block[j].Op = bytecode.OpInvalid
			} else {
				block[i].Op = block[j].Op
				block[i].Arg1 = block[j].Arg1
				block[j].Op = bytecode.OpInvalid
			}
			dstslot := inst.Arg0.ToSlot()
			if dstslot.IsUniq() {
				constValues.Add(dstslot, uint16(i))
			}
			changed = true
		}

		for argIndex, argInfo := range inst.Op.Args() {
			if argInfo.Kind != bytecode.ArgSlot {
				continue
			}
			if !argInfo.IsReadSlot() {
				continue
			}
			arg := inst.GetArg(argIndex).ToSlot()
			if !arg.IsUniq() {
				continue
			}
			constValues.Remove(arg)
		}
	}

	return changed
}
