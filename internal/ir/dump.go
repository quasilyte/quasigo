package ir

import (
	"fmt"
	"strings"

	"github.com/quasilyte/quasigo/internal/bytecode"
)

func Dump(fn *Func) string {
	if fn.Blocks == nil {
		panic("dumping IR func without blocks")
	}

	env := fn.Env

	var out strings.Builder

	tempSeq := make(map[uint8]int)

	fmtSlot := func(slot Slot) string {
		switch slot.Kind {
		case SlotCallArg:
			return fmt.Sprintf("arg%d", slot.ID)
		case SlotParam:
			return fn.Debug.SlotNames[fn.SlotIndex(slot)]
		case SlotLocal:
			return fn.Debug.SlotNames[fn.SlotIndex(slot)]
		case SlotTemp:
			return fmt.Sprintf("temp%d", slot.ID)
		case SlotUniq:
			subscript := tempSeq[slot.ID]
			return fmt.Sprintf("temp%d.v%d", slot.ID, subscript-1)
		}
		return "?"
	}

	args := make([]string, 0, 4)
	for i, b := range fn.Blocks {
		if b.Label != 0 {
			fmt.Fprintf(&out, "block%d (L%d) [%d]:\n", i, b.Label-1, b.NumVarKill)
		} else {
			fmt.Fprintf(&out, "block%d [%d]:\n", i, b.NumVarKill)
		}

		for _, inst := range b.Code {
			if inst.Pseudo == OpVarKill {
				fmt.Fprintf(&out, "  VarKill temp%d\n", inst.Arg0.ToSlot().ID)
				continue
			}
			if inst.Op == bytecode.OpInvalid {
				continue
			}

			if inst.Op.HasDst() && inst.Arg0.ToSlot().IsUniq() {
				tempSeq[inst.Arg0.ToSlot().ID]++
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
					value = env.UserFuncs[arg].Name
				case bytecode.ArgNativeFuncID:
					value = env.NativeFuncs[arg].Name
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
	}

	return out.String()
}
