package qopt

import (
	"github.com/quasilyte/quasigo/internal/bytecode"
	"github.com/quasilyte/quasigo/internal/ir"
)

func createBlocks(fn *ir.Func, out []ir.Block) []ir.Block {
	builder := blocksBuilder{fn: fn, out: out}
	return builder.Build()
}

type blocksBuilder struct {
	fn  *ir.Func
	out []ir.Block
}

func (builder *blocksBuilder) Build() []ir.Block {
	code := builder.fn.Code

	blockLabel := uint16(0)
	blockStart := 0
	numVarKill := 0
	seenBlockInst := false
	i := 0
	for {
		if i >= len(code) {
			break
		}
		inst := code[i]

		if inst.Pseudo == ir.OpVarKill {
			i++
			if seenBlockInst {
				numVarKill++
			} else {
				blockStart++
			}
			continue
		}

		if inst.Pseudo == ir.OpLabel {
			blockCode := code[blockStart:i]
			builder.pushBlock(ir.Block{
				Code:       blockCode,
				NumVarKill: uint16(numVarKill),
				Label:      blockLabel,
			})
			blockStart = i + 1
			numVarKill = 0
			seenBlockInst = false
			blockLabel = uint16(inst.Arg0) + 1
			i++
			continue
		}

		seenBlockInst = true

		switch inst.Op {
		case bytecode.OpJump, bytecode.OpJumpZero, bytecode.OpJumpNotZero:
			fallthrough
		case bytecode.OpReturnZero, bytecode.OpReturnOne, bytecode.OpReturnVoid:
			fallthrough
		case bytecode.OpReturnScalar, bytecode.OpReturnStr:
			i++
			for i < len(code) && code[i].Pseudo == ir.OpVarKill {
				i++
				numVarKill++
			}
			blockCode := code[blockStart:i]
			if len(blockCode) != 0 {
				builder.pushBlock(ir.Block{
					Code:       blockCode,
					NumVarKill: uint16(numVarKill),
					Label:      blockLabel,
				})
			}
			blockStart = i
			numVarKill = 0
			seenBlockInst = false
			blockLabel = 0
			continue
		}

		i++
	}

	if len(builder.out) == 0 {
		builder.pushBlock(ir.Block{
			Code:       code,
			NumVarKill: uint16(numVarKill),
		})
	}

	return builder.out
}

func (builder *blocksBuilder) pushBlock(b ir.Block) {
	if len(builder.out) != 0 {
		last := &builder.out[len(builder.out)-1]
		if len(last.Code) == 0 && (last.Label == 0 || b.Label == 0) {
			if b.Label != 0 {
				last.Label = b.Label
			}
			last.Code = b.Code
			last.NumVarKill = b.NumVarKill
			return
		}
	}
	builder.out = append(builder.out, b)
}
