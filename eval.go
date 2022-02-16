package quasigo

import (
	"fmt"
)

//go:noinline
func panicStackOverflow(fn *Func) {
	panic(fmt.Sprintf("can't call %s func: stack overflow", fn.name))
}

func eval(env *EvalEnv, fn *Func, slotptr *slotValue) {
	pc := 0
	codeptr := fn.codeptr

	for {
		switch op := opcode(unpack8(codeptr, pc)); op {
		case opLoadStrConst:
			dstslot, constindex := unpack8x2(codeptr, pc+1)
			getslot(slotptr, dstslot).SetString(fn.strConstants[constindex])
			pc += 3
		case opLoadScalarConst:
			dstslot, constindex := unpack8x2(codeptr, pc+1)
			getslot(slotptr, dstslot).scalar = fn.scalarConstants[constindex]
			pc += 3

		case opMoveStr:
			dstslot, srcslot := unpack8x2(codeptr, pc+1)
			getslot(slotptr, dstslot).SetString(getslot(slotptr, srcslot).String())
			pc += 3
		case opMoveScalar:
			dstslot, srcslot := unpack8x2(codeptr, pc+1)
			getslot(slotptr, dstslot).scalar = getslot(slotptr, srcslot).scalar
			pc += 3
		case opMoveInterface:
			dstslot, srcslot := unpack8x2(codeptr, pc+1)
			getslot(slotptr, dstslot).MoveInterface(getslot(slotptr, srcslot))
			pc += 3
		case opMoveResult2:
			dstslot := unpack8(codeptr, pc+1)
			*getslot(slotptr, dstslot) = env.result2
			pc += 2

		case opStrLen:
			dstslot, srcslot := unpack8x2(codeptr, pc+1)
			getslot(slotptr, dstslot).scalar = uint64(len(getslot(slotptr, srcslot).String()))
			pc += 3
		case opStrSliceFrom:
			dstslot, strslot, fromslot := unpack8x3(codeptr, pc+1)
			str := getslot(slotptr, strslot).String()
			from := getslot(slotptr, fromslot).scalar
			getslot(slotptr, dstslot).SetString(str[from:])
			pc += 4
		case opStrSliceTo:
			dstslot, strslot, toslot := unpack8x3(codeptr, pc+1)
			str := getslot(slotptr, strslot).String()
			to := getslot(slotptr, toslot).scalar
			getslot(slotptr, dstslot).SetString(str[:to])
			pc += 4
		case opStrSlice:
			dstslot, strslot, fromslot, toslot := unpack8x4(codeptr, pc+1)
			str := getslot(slotptr, strslot).String()
			from := getslot(slotptr, fromslot).scalar
			to := getslot(slotptr, toslot).scalar
			getslot(slotptr, dstslot).SetString(str[from:to])
			pc += 5

		case opNot:
			dstslot, srcslot := unpack8x2(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(!getslot(slotptr, srcslot).Bool())
			pc += 3

		case opIsNil:
			dstslot, srcslot := unpack8x2(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(getslot(slotptr, srcslot).IsNil())
			pc += 3
		case opIsNotNil:
			dstslot, srcslot := unpack8x2(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(!getslot(slotptr, srcslot).IsNil())
			pc += 3
		case opIsNilInterface:
			dstslot, srcslot := unpack8x2(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(getslot(slotptr, srcslot).IsNilInterface())
			pc += 3
		case opIsNotNilInterface:
			dstslot, srcslot := unpack8x2(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(!getslot(slotptr, srcslot).IsNilInterface())
			pc += 3

		case opStrEq:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(getslot(slotptr, xslot).String() == getslot(slotptr, yslot).String())
			pc += 4
		case opStrNotEq:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(getslot(slotptr, xslot).String() != getslot(slotptr, yslot).String())
			pc += 4

		case opIntEq:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(getslot(slotptr, xslot).Int() == getslot(slotptr, yslot).Int())
			pc += 4
		case opIntNotEq:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(getslot(slotptr, xslot).Int() != getslot(slotptr, yslot).Int())
			pc += 4
		case opIntGt:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(getslot(slotptr, xslot).Int() > getslot(slotptr, yslot).Int())
			pc += 4
		case opIntGtEq:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(getslot(slotptr, xslot).Int() >= getslot(slotptr, yslot).Int())
			pc += 4
		case opIntLt:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(getslot(slotptr, xslot).Int() < getslot(slotptr, yslot).Int())
			pc += 4
		case opIntLtEq:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(getslot(slotptr, xslot).Int() <= getslot(slotptr, yslot).Int())
			pc += 4

		case opConcat:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).SetString(getslot(slotptr, xslot).String() + getslot(slotptr, yslot).String())
			pc += 4

		case opIntAdd:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).scalar = getslot(slotptr, xslot).scalar + getslot(slotptr, yslot).scalar
			pc += 4
		case opIntSub:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).scalar = getslot(slotptr, xslot).scalar - getslot(slotptr, yslot).scalar
			pc += 4
		case opIntMul:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).scalar = getslot(slotptr, xslot).scalar * getslot(slotptr, yslot).scalar
			pc += 4
		case opIntDiv:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).scalar = getslot(slotptr, xslot).scalar / getslot(slotptr, yslot).scalar
			pc += 4

		case opIntInc:
			dstslot := unpack8(codeptr, pc+1)
			getslot(slotptr, dstslot).scalar++
			pc += 2
		case opIntDec:
			dstslot := unpack8(codeptr, pc+1)
			getslot(slotptr, dstslot).scalar--
			pc += 2

		case opJump:
			offset := unpack16(codeptr, pc+1)
			pc += offset

		case opJumpFalse:
			srcslot := unpack8(codeptr, pc+3)
			if !getslot(slotptr, srcslot).Bool() {
				offset := unpack16(codeptr, pc+1)
				pc += offset
			} else {
				pc += 4
			}
		case opJumpTrue:
			srcslot := unpack8(codeptr, pc+3)
			if getslot(slotptr, srcslot).Bool() {
				offset := unpack16(codeptr, pc+1)
				pc += offset
			} else {
				pc += 4
			}

		case opCall:
			dstslot := unpack8(codeptr, pc+1)
			funcid := unpack16(codeptr, pc+2)
			callFunc := env.userFuncs[funcid]
			if !canAllocFrame(slotptr, env.slotend, callFunc.frameSize) {
				panicStackOverflow(fn)
			}
			eval(env, callFunc, nextFrameSlot(slotptr, fn.frameSize))
			*getslot(slotptr, dstslot) = env.result
			pc += 4
		case opCallRecur:
			dstslot := unpack8(codeptr, pc+1)
			if !canAllocFrame(slotptr, env.slotend, fn.frameSize) {
				panicStackOverflow(fn)
			}
			eval(env, fn, nextFrameSlot(slotptr, fn.frameSize))
			*getslot(slotptr, dstslot) = env.result
			pc += 2

		case opVariadicReset:
			env.vararg = env.vararg[:0]
			pc++
		case opPushVariadicBoolArg:
			srcslot := unpack8(codeptr, pc+1)
			env.vararg = append(env.vararg, getslot(slotptr, srcslot).Bool())
			pc += 2
		case opPushVariadicScalarArg:
			srcslot := unpack8(codeptr, pc+1)
			env.vararg = append(env.vararg, getslot(slotptr, srcslot).Int())
			pc += 2
		case opPushVariadicStrArg:
			srcslot := unpack8(codeptr, pc+1)
			env.vararg = append(env.vararg, getslot(slotptr, srcslot).String())
			pc += 2
		case opPushVariadicInterfaceArg:
			srcslot := unpack8(codeptr, pc+1)
			env.vararg = append(env.vararg, getslot(slotptr, srcslot).Interface())
			pc += 2

		case opCallNative:
			dstslot := unpack8(codeptr, pc+1)
			funcid := unpack16(codeptr, pc+2)
			callFunc := env.nativeFuncs[funcid]
			if !canAllocFrame(slotptr, env.slotend, callFunc.frameSize) {
				panicStackOverflow(fn)
			}
			callFunc.mappedFunc(NativeCallContext{
				env:     env,
				slotptr: nextFrameSlot(slotptr, fn.frameSize),
			})
			*getslot(slotptr, dstslot) = env.result
			pc += 4
		case opCallVoidNative:
			funcid := unpack16(codeptr, pc+1)
			callFunc := env.nativeFuncs[funcid]
			if !canAllocFrame(slotptr, env.slotend, callFunc.frameSize) {
				panicStackOverflow(fn)
			}
			callFunc.mappedFunc(NativeCallContext{
				slotptr: nextFrameSlot(slotptr, fn.frameSize),
			})
			pc += 3

		case opReturnVoid:
			return
		case opReturnTrue:
			env.result.SetBool(true)
			return
		case opReturnFalse:
			env.result.SetBool(false)
			return
		case opReturnStr:
			srcslot := unpack8(codeptr, pc+1)
			env.result.SetString(getslot(slotptr, srcslot).String())
			return
		case opReturnScalar:
			srcslot := unpack8(codeptr, pc+1)
			env.result.scalar = getslot(slotptr, srcslot).scalar
			return
		case opReturnInterface:
			srcslot := unpack8(codeptr, pc+1)
			env.result = *getslot(slotptr, srcslot)
			return

		default:
			panic(fmt.Sprintf("malformed bytecode: unexpected %s found", op))
		}
	}
}
