package ir

import (
	"fmt"
	"strings"

	"github.com/quasilyte/quasigo/internal/bytecode"
)

func Dump(fn *Func) string {
	var out strings.Builder

	tmpSeq := make(map[uint8]int)

	fmtSlot := func(slot Slot) string {
		switch slot.Kind {
		case SlotCallArg:
			return fmt.Sprintf("arg%d", slot.ID)
		case SlotParam:
			return fmt.Sprintf("p%d", slot.ID)
		case SlotLocal:
			return fmt.Sprintf("l%d", slot.ID)
		case SlotTemp:
			return fmt.Sprintf("tmp%d", slot.ID)
		case SlotUniq:
			subscript := tmpSeq[slot.ID]
			return fmt.Sprintf("tmp%d_%d", slot.ID, subscript-1)
		}
		return "?"
	}

	args := make([]string, 0, 4)
	for _, inst := range fn.Code {
		if inst.Pseudo == OpLabel {
			fmt.Fprintf(&out, "L%d:\n", inst.Arg0)
			continue
		}

		if inst.Op.HasDst() && inst.Arg0.ToSlot().IsUniq() {
			tmpSeq[inst.Arg0.ToSlot().ID]++
		}

		args = args[:0]
		for i, argInfo := range inst.Op.Args() {
			arg := inst.GetArg(i)
			var value string
			switch argInfo.Kind {
			case bytecode.ArgSlot:
				value = fmtSlot(arg.ToSlot())
			case bytecode.ArgStrConst:
				value = fmt.Sprintf("%q", fn.StrConstants[arg])
			case bytecode.ArgScalarConst:
				value = fmt.Sprintf("%d", int64(fn.ScalarConstants[arg]))
			case bytecode.ArgOffset:
				value = fmt.Sprintf("L%d", arg)
			case bytecode.ArgFuncID:
				value = fmt.Sprintf("func%d", arg)
			case bytecode.ArgNativeFuncID:
				value = fmt.Sprintf("nativefunc%d", arg)
			}
			if inst.Op.HasDst() && i == 0 && len(inst.Op.Args()) != 1 {
				args = append(args, value, "=")
			} else {
				args = append(args, value)
			}
		}

		out.WriteString("  ")
		out.WriteString(inst.Op.String())
		if len(args) != 0 {
			out.WriteByte(' ')
			out.WriteString(strings.Join(args, " "))
		}
		out.WriteByte('\n')
	}

	return out.String()
}
