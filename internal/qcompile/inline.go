package qcompile

import (
	"encoding/binary"
	"go/ast"

	"github.com/quasilyte/quasigo/internal/bytecode"
	"github.com/quasilyte/quasigo/internal/ir"
	"github.com/quasilyte/quasigo/internal/qruntime"
)

type inliner struct {
	cl *funcCompiler

	labelRet label
	fn       *qruntime.Func
	irfn     ir.Func

	failed bool

	dst ir.Slot

	tempOffset int
}

var inlineableOps = [256]bool{
	bytecode.OpIntNeg:   true,
	bytecode.OpIntAdd8:  true,
	bytecode.OpIntAdd64: true,
	bytecode.OpIntSub8:  true,
	bytecode.OpIntSub64: true,
	bytecode.OpIntMul8:  true,
	bytecode.OpIntMul64: true,
	bytecode.OpIntXor:   true,
	bytecode.OpIntDiv:   true,
	bytecode.OpConcat:   true,

	bytecode.OpIntInc: true,
	bytecode.OpIntDec: true,

	bytecode.OpCap: true,

	bytecode.OpStrSlice:     true,
	bytecode.OpStrSliceFrom: true,
	bytecode.OpStrSliceTo:   true,

	bytecode.OpStrIndex:           true,
	bytecode.OpSliceIndexScalar8:  true,
	bytecode.OpSliceIndexScalar64: true,
	bytecode.OpSliceSetScalar8:    true,
	bytecode.OpSliceSetScalar64:   true,

	bytecode.OpNot:               true,
	bytecode.OpIsNil:             true,
	bytecode.OpIsNotNil:          true,
	bytecode.OpIsNilInterface:    true,
	bytecode.OpIsNotNilInterface: true,

	bytecode.OpScalarEq:    true,
	bytecode.OpScalarNotEq: true,
	bytecode.OpIntLtEq:     true,
	bytecode.OpIntLt:       true,
	bytecode.OpIntGtEq:     true,
	bytecode.OpIntGt:       true,
	bytecode.OpStrEq:       true,
	bytecode.OpStrNotEq:    true,
	bytecode.OpStrLt:       true,
	bytecode.OpStrGt:       true,

	bytecode.OpMove:            true,
	bytecode.OpMove8:           true,
	bytecode.OpLoadScalarConst: true,
	bytecode.OpLoadStrConst:    true,

	bytecode.OpReturn:       true,
	bytecode.OpReturnStr:    true,
	bytecode.OpReturnScalar: true,
	bytecode.OpReturnZero:   true,
	bytecode.OpReturnOne:    true,

	bytecode.OpJump:        true,
	bytecode.OpJumpZero:    true,
	bytecode.OpJumpNotZero: true,
}

func (cl *funcCompiler) inlineCall(dst ir.Slot, recv ast.Expr, args []ast.Expr, key qruntime.FuncKey) bool {
	if !cl.ctx.Optimize || !cl.ctx.Static {
		return false
	}

	funcID, ok := cl.ctx.Env.NameToFuncID[key]
	if !ok {
		return false
	}
	fn := cl.ctx.Env.UserFuncs[funcID]
	if !fn.CanInline || len(fn.Code) > 48 {
		return false
	}
	// If temps pressure is high and inlining candidate frame is big,
	// do not attempt any inlining.
	numFrameSlots := int(fn.NumParams) + int(fn.NumTemps)
	if cl.tempSeq+numFrameSlots > 100 {
		return false
	}

	numLabels := cl.numLabels
	numTemps := cl.numTemp
	tempSeq := cl.tempSeq
	inl := inliner{
		cl:       cl,
		fn:       fn,
		labelRet: cl.newLabel(),
		irfn: ir.Func{
			NumParams: int(fn.NumParams),
			NumTemps:  int(fn.NumTemps),
		},
		dst:        dst,
		tempOffset: cl.tempSeq,
	}
	inlined := inl.Inline(recv, args)
	if inlined {
		cl.bindLabel(inl.labelRet)
	} else {
		cl.numLabels = numLabels
		cl.numTemp = numTemps
		cl.tempSeq = tempSeq
	}
	return inlined
}

func (inl *inliner) Inline(recv ast.Expr, args []ast.Expr) bool {
	fn := inl.fn

	// This labels collection code is copied from the qdisasm package.
	// TODO: maybe we can re-use this code somehow?
	var labels map[int]label // Lazily allocated
	decode16 := func(code []byte, pos int) int {
		return int(int16(binary.LittleEndian.Uint16(code[pos:])))
	}
	bytecode.Walk(fn.Code, func(pc int, op bytecode.Op) {
		if inl.failed {
			return
		}
		if !inlineableOps[op] {
			inl.failed = true
			return
		}
		if !op.IsJump() {
			return
		}
		if op == bytecode.OpJumpTable {
			return
		}
		offset := decode16(fn.Code, pc+1)
		targetPC := pc + offset
		if labels == nil {
			labels = make(map[int]label, 2)
		}
		if _, ok := labels[targetPC]; !ok {
			labels[targetPC] = inl.cl.newLabel()
		}
	})

	inlined := make([]ir.Inst, 0, 8)
	bytecode.Walk(fn.Code, func(pc int, op bytecode.Op) {
		if inl.failed {
			return
		}

		if l, ok := labels[pc]; ok {
			inlined = append(inlined, ir.Inst{
				Pseudo: ir.OpLabel,
				Arg0:   ir.InstArg(l.id),
			})
		}

		// Some ops require ad-hoc handling.
		isReturnOp := false
		switch op {
		case bytecode.OpReturn:
			isReturnOp = true
		case bytecode.OpReturnScalar, bytecode.OpReturnStr:
			srcslot := inl.convertSlot(fn.Code[pc+1])
			inlined = append(inlined, ir.Inst{
				Op:   bytecode.OpMove,
				Arg0: inl.dst.ToInstArg(),
				Arg1: srcslot.ToInstArg(),
			})
			isReturnOp = true
		case bytecode.OpReturnZero:
			inlined = append(inlined, inl.cl.moveInt(inl.dst, 0))
			isReturnOp = true
		case bytecode.OpReturnOne:
			inlined = append(inlined, inl.cl.moveInt(inl.dst, 1))
			isReturnOp = true
		}
		if isReturnOp {
			isLast := (pc + op.Width()) == len(fn.Code)
			if !isLast {
				inlined = append(inlined, ir.Inst{
					Op:   bytecode.OpJump,
					Arg0: ir.InstArg(inl.labelRet.id),
				})
			}
			return
		}

		// Handle other generic ops by mapping the arguments.
		inst := ir.Inst{Op: op}
		for i, a := range op.Args() {
			switch a.Kind {
			case bytecode.ArgSlot:
				slot := fn.Code[pc+int(a.Offset)]
				inst.SetArg(i, inl.convertSlot(slot).ToInstArg())
			case bytecode.ArgScalarConst:
				constindex := fn.Code[pc+int(a.Offset)]
				inst.SetArg(i, inl.convertScalarConst(constindex))
			case bytecode.ArgStrConst:
				constindex := fn.Code[pc+int(a.Offset)]
				inst.SetArg(i, inl.convertStrConst(constindex))
			case bytecode.ArgOffset:
				offset := decode16(fn.Code, pc+int(a.Offset))
				targetPC := pc + offset
				l, ok := labels[targetPC]
				if !ok {
					inl.failed = true
					return
				}
				inst.SetArg(i, ir.InstArg(l.id))
			default:
				inl.failed = true
				return
			}
		}
		inlined = append(inlined, inst)
	})

	if inl.failed || len(inlined) == 0 {
		return false
	}

	inl.compileInlinedCallArgs(recv, args)
	inl.cl.code = append(inl.cl.code, inlined...)

	return true
}

func (inl *inliner) compileInlinedCallArgs(recv ast.Expr, args []ast.Expr) {
	inl.cl.checkTupleArg(args)
	args = inl.cl.prependRecv(recv, args)

	for i, arg := range args {
		id := int(inl.irfn.SlotIndex(ir.NewParamSlot(uint8(i)))) + inl.tempOffset
		inl.cl.trackTemp(id)
		argslot := ir.NewTempSlot(uint8(id))
		inl.cl.CompileExpr(argslot, arg)
	}
}

func (inl *inliner) convertStrConst(constindex byte) ir.InstArg {
	v := inl.fn.StrConstants[constindex]
	return ir.InstArg(inl.cl.internStrConstant(v))
}

func (inl *inliner) convertScalarConst(constindex byte) ir.InstArg {
	v := inl.fn.ScalarConstants[constindex]
	return ir.InstArg(inl.cl.internScalarConstant(v))
}

func (inl *inliner) convertSlot(id byte) ir.Slot {
	var orig ir.Slot
	switch {
	case id < inl.fn.NumParams:
		orig = ir.NewParamSlot(id)
	case id < inl.fn.FrameSlots:
		orig = ir.NewTempSlot(id)
	default:
		orig = ir.NewCallArgSlot(id)
	}
	newid := int(inl.irfn.SlotIndex(orig)) + inl.tempOffset
	inl.cl.trackTemp(newid)
	return ir.NewTempSlot(uint8(newid))
}
