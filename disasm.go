package quasigo

import (
	"fmt"
	"strings"

	"github.com/quasilyte/quasigo/internal/bytecode"
)

func disasm(env *Env, fn *Func) string {
	var out strings.Builder

	dbg, ok := env.debug.funcs[fn]
	if !ok {
		return "<unknown>\n"
	}

	numSlots := fn.frameSize / int(sizeofSlotValue)
	numLocals := dbg.numLocals
	numArgs := len(dbg.slotNames) - numLocals
	numTemps := numSlots - numArgs - numLocals
	fmt.Fprintf(&out, "%s code=%d frame=%d (%d slots: %d args, %d locals, %d temps)\n",
		fn.name, len(fn.code), fn.frameSize, numSlots, numArgs, numLocals, numTemps)

	slotName := func(index int) string {
		if index < len(dbg.slotNames) {
			return dbg.slotNames[index]
		}
		if index >= numSlots {
			return fmt.Sprintf("arg%d", index-numSlots)
		}
		return fmt.Sprintf("tmp%d", index-len(dbg.slotNames))
	}

	code := fn.code
	labels := map[int]string{}
	bytecode.Walk(code, func(pc int, op bytecode.Op) {
		if !op.IsJump() {
			return
		}
		offset := unpack16(fn.codeptr, pc+1)
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
				slot := int(code[pc+a.Offset])
				value = slotName(slot)
			case bytecode.ArgStrConst:
				index := int(code[pc+a.Offset])
				value = fmt.Sprintf("%q", fn.strConstants[index])
			case bytecode.ArgScalarConst:
				index := int(code[pc+a.Offset])
				value = fmt.Sprintf("%d", int64(fn.scalarConstants[index]))
			case bytecode.ArgOffset:
				offset := unpack16(fn.codeptr, pc+a.Offset)
				targetPC := pc + offset
				value = labels[targetPC]
			case bytecode.ArgFuncID:
				id := unpack16(fn.codeptr, pc+a.Offset)
				value = env.userFuncs[id].name + "()"
			case bytecode.ArgNativeFuncID:
				id := unpack16(fn.codeptr, pc+a.Offset)
				value = env.nativeFuncs[id].name + "()"
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
