package qdisasm

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/quasilyte/quasigo/internal/bytecode"
	"github.com/quasilyte/quasigo/internal/qruntime"
)

func Func(env *qruntime.Env, fn *qruntime.Func) string {
	var out strings.Builder

	dbg, ok := env.Debug.Funcs[fn]
	if !ok {
		return "<unknown>\n"
	}

	numSlots := fn.FrameSize / int(qruntime.SizeofSlot)
	numLocals := int(fn.NumLocals)
	numArgs := len(dbg.SlotNames) - numLocals
	numTemps := numSlots - numArgs - numLocals
	fmt.Fprintf(&out, "%s code=%d frame=%d (%d slots: %d args, %d locals, %d temps)\n",
		fn.Name, len(fn.Code), fn.FrameSize, numSlots, numArgs, numLocals, numTemps)

	slotName := func(index int) string {
		if index < len(dbg.SlotNames) {
			return dbg.SlotNames[index]
		}
		if index >= numSlots {
			return fmt.Sprintf("arg%d", index-numSlots)
		}
		return fmt.Sprintf("temp%d", index-len(dbg.SlotNames))
	}

	code := fn.Code
	labels := make(map[int]string)
	bytecode.Walk(code, func(pc int, op bytecode.Op) {
		if !op.IsJump() {
			return
		}
		if op == bytecode.OpJumpTable {
			return
		}
		offset := decode16(fn.Code, pc+1)
		targetPC := pc + offset
		if _, ok := labels[targetPC]; !ok {
			labels[targetPC] = fmt.Sprintf("L%d", len(labels))
		}
	})

	args := make([]string, 0, 4)
	bytecode.Walk(code, func(pc int, op bytecode.Op) {
		if l := labels[pc]; l != "" {
			fmt.Fprintf(&out, "%s:\n", l)
		}
		args = args[:0]

		for i, a := range op.Args() {
			var value string
			switch a.Kind {
			case bytecode.ArgSlot:
				slot := int(code[pc+int(a.Offset)])
				value = slotName(slot)
			case bytecode.ArgStrConst:
				index := int(code[pc+int(a.Offset)])
				value = fmt.Sprintf("%q", fn.StrConstants[index])
			case bytecode.ArgScalarConst:
				index := int(code[pc+int(a.Offset)])
				value = fmt.Sprintf("%d", int64(fn.ScalarConstants[index]))
			case bytecode.ArgOffset:
				offset := decode16(fn.Code, pc+int(a.Offset))
				targetPC := pc + offset
				value = labels[targetPC]
			case bytecode.ArgFuncID:
				id := decode16(fn.Code, pc+int(a.Offset))
				value = env.UserFuncs[id].Name + "()"
			case bytecode.ArgNativeFuncID:
				id := decode16(fn.Code, pc+int(a.Offset))
				value = env.NativeFuncs[id].Name + "()"
			}
			if op.HasDst() && i == 0 && len(op.Args()) != 1 {
				args = append(args, value, "=")
			} else {
				args = append(args, value)
			}
		}

		out.WriteString("  ")
		out.WriteString(op.String())
		if len(args) != 0 {
			out.WriteByte(' ')
			out.WriteString(strings.Join(args, " "))
		}
		out.WriteByte('\n')
	})

	return out.String()
}

func decode16(code []byte, pos int) int {
	return int(int16(binary.LittleEndian.Uint16(code[pos:])))
}
