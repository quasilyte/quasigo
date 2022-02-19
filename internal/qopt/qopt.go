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
}

func (opt *optimizer) Optimize() {
	scalarTempStores := idMap{make([]uint16, 0, 4)}
	strTempStores := idMap{make([]uint16, 0, 4)}
	changed := false
	for instIndex, inst := range opt.fn.Code {
		if inst.Pseudo == ir.OpLabel {
			scalarTempStores.Reset()
			strTempStores.Reset()
			return
		}

		switch inst.Op {
		case bytecode.OpLoadScalarConst:
			dst := int(inst.Arg0)
			if !opt.isTemp(dst) {
				return
			}
			scalarTempStores.Add(dst, int(inst.Arg1))
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
				return
			}
			strTempStores.Add(dst, int(inst.Arg1))
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
