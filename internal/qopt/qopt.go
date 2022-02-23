package qopt

import (
	"github.com/quasilyte/quasigo/internal/bytecode"
	"github.com/quasilyte/quasigo/internal/ir"
)

// TODO:
// * optimize expressions like `x + 0`
// * optimize `x += 1` to `x++`
// * optimize `Not x = y; JumpZero x` to `JumpNotZero x` if x is tmp value

func Func(fn *ir.Func) {
	opt := optimizer{fn: fn}
	opt.Optimize()
}

type optimizer struct {
	fn *ir.Func
}

func (opt *optimizer) Optimize() {
	scalarTempStores := idMap{make([]uint16, 0, 4)}
	strTempStores := idMap{make([]uint16, 0, 4)}
	changed := false
	for instIndex, inst := range opt.fn.Code {
		if inst.Pseudo == ir.OpLabel {
			scalarTempStores.Reset()
			strTempStores.Reset()
			continue
		}

		storeHandled := false
		switch inst.Op {
		case bytecode.OpLoadScalarConst:
			dst := int(inst.Arg0)
			if !opt.isTemp(dst) {
				continue
			}
			if i := scalarTempStores.Find(dst); i != -1 {
				scalarTempStores.UpdateValueAt(i, int(inst.Arg1))
			} else {
				scalarTempStores.Add(dst, int(inst.Arg1))
			}
			storeHandled = true
		case bytecode.OpMoveScalar:
			src := int(inst.Arg1)
			if i := scalarTempStores.Find(src); i != -1 {
				changed = true
				constindex := scalarTempStores.GetValue(i)
				opt.fn.Code[instIndex] = ir.Inst{
					Op:   bytecode.OpLoadScalarConst,
					Arg0: inst.Arg0,
					Arg1: constindex,
				}
			}

		case bytecode.OpLoadStrConst:
			dst := int(inst.Arg0)
			if !opt.isTemp(dst) {
				continue
			}
			if i := strTempStores.Find(dst); i != -1 {
				strTempStores.UpdateValueAt(i, int(inst.Arg1))
			} else {
				strTempStores.Add(dst, int(inst.Arg1))
			}
			storeHandled = true
		case bytecode.OpMoveStr:
			src := int(inst.Arg1)
			if i := strTempStores.Find(src); i != -1 {
				changed = true
				constindex := strTempStores.GetValue(i)
				opt.fn.Code[instIndex] = ir.Inst{
					Op:   bytecode.OpLoadStrConst,
					Arg0: inst.Arg0,
					Arg1: constindex,
				}
			}
		}

		if !storeHandled && inst.Op.HasDst() {
			dst := int(inst.Arg0)
			if !opt.isTemp(dst) {
				continue
			}
			if i := scalarTempStores.Find(dst); i != -1 {
				scalarTempStores.RemoveAt(i)
				continue
			}
			if i := strTempStores.Find(dst); i != -1 {
				strTempStores.RemoveAt(i)
				continue
			}
		}
	}

	if !changed {
		return
	}

	maxTempID := 0
	usedSlots := idSet{make([]uint8, 0, 8)}
	for i := len(opt.fn.Code) - 1; i >= 0; i-- {
		inst := &opt.fn.Code[i]
		if inst.Op.HasDst() && !opt.isCallArg(int(inst.Arg0)) && !usedSlots.Contains(int(inst.Arg0)) {
			inst.Op = bytecode.OpInvalid
			continue
		}
		inst.WalkArgs(func(arg bytecode.Argument, value int) {
			if arg.Kind != bytecode.ArgSlot {
				return
			}
			if opt.isCallArg(value) {
				return
			}
			if opt.isTemp(value) {
				if value > maxTempID {
					maxTempID = value
				}
			}
			if arg.IsWriteSlot() {
				usedSlots.Remove(value)
			} else {
				usedSlots.Add(value)
			}
		})
	}

	if maxTempID != 0 {
		opt.fn.NumFrameSlots = maxTempID + 1
	} else {
		opt.fn.NumFrameSlots = opt.fn.NumParams + opt.fn.NumLocals
	}
}

func (opt *optimizer) isCallArg(id int) bool {
	// TODO: something better?
	return id > 200
}

func (opt *optimizer) isTemp(id int) bool {
	return id >= (opt.fn.NumParams + opt.fn.NumLocals)
}
