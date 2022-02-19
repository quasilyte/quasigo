package qcompile

import (
	"errors"
	"math"

	"github.com/quasilyte/quasigo/internal/bytecode"
	"github.com/quasilyte/quasigo/internal/ir"
)

type assembler struct {
	code []byte
}

func (asm *assembler) Assemble(fn *ir.Func) ([]byte, error) {
	asm.code = make([]byte, 0, len(fn.Code)*2)

	var labelOffsets [256]uint16
	encBuf := make([]byte, 0, 8)
	for _, inst := range fn.Code {
		if inst.Pseudo == ir.OpLabel {
			id := inst.Arg0
			if len(asm.code) > math.MaxUint16 {
				return nil, errors.New("label offset doesn't fit in uint16")
			}
			labelOffsets[id] = uint16(len(asm.code))
		}
		if inst.IsPseudo() || inst.Op == bytecode.OpInvalid {
			continue
		}

		encBuf = encBuf[:0]
		encBuf = append(encBuf, byte(inst.Op))
		inst.WalkArgs(func(arg bytecode.Argument, value int) {
			switch arg.Kind {
			case bytecode.ArgSlot, bytecode.ArgScalarConst, bytecode.ArgStrConst:
				encBuf = append(encBuf, byte(value))
			case bytecode.ArgOffset:
				// The actual jump targets are linked later.
				encBuf = append(encBuf, byte(value), 0)
			default:
				encBuf = append(encBuf, 0, 0)
				put16(encBuf, len(encBuf)-2, int(inst.Value))
			}
		})
		asm.code = append(asm.code, encBuf...)
	}

	asm.linkJumps(&labelOffsets)

	// This byte slice will never grow.
	// Cut down the excessive capacity.
	asm.code = asm.code[:len(asm.code):len(asm.code)]
	return asm.code, nil
}

func (asm *assembler) linkJumps(labelOffsets *[256]uint16) {
	bytecode.Walk(asm.code, func(pc int, op bytecode.Op) {
		if !op.IsJump() {
			return
		}
		labelID := asm.code[pc+1]
		targetPos := int(labelOffsets[labelID])
		jumpOffset := targetPos - pc
		patchPos := pc + 1
		put16(asm.code, patchPos, jumpOffset)
	})
}
