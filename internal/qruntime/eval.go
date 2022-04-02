package qruntime

import (
	"fmt"

	"github.com/quasilyte/quasigo/internal/bytecode"
)

//go:noinline
func panicStackOverflow(fn *Func) {
	panic(fmt.Sprintf("can't call %s func: stack overflow", fn.Name))
}

func eval(env *EvalEnv, fn *Func, slotptr *Slot) {
	pc := 0
	codeptr := fn.Codeptr

	for {
		switch op := bytecode.Op(unpack8(codeptr, pc)); op {
		case bytecode.OpLoadStrConst:
			dstslot, constindex := unpack8x2(codeptr, pc+1)
			getslot(slotptr, dstslot).SetString(fn.StrConstants[constindex])
			pc += 3
		case bytecode.OpLoadScalarConst:
			dstslot, constindex := unpack8x2(codeptr, pc+1)
			getslot(slotptr, dstslot).Scalar = fn.ScalarConstants[constindex]
			pc += 3

		case bytecode.OpZero:
			dstslot := unpack8(codeptr, pc+1)
			*getslot(slotptr, dstslot) = Slot{}
			pc += 2
		case bytecode.OpMove:
			dstslot, srcslot := unpack8x2(codeptr, pc+1)
			*getslot(slotptr, dstslot) = *getslot(slotptr, srcslot)
			pc += 3
		case bytecode.OpMove8:
			dstslot, srcslot := unpack8x2(codeptr, pc+1)
			getslot(slotptr, dstslot).SetByte(getslot(slotptr, srcslot).Byte())
			pc += 3
		case bytecode.OpMoveResult2:
			dstslot := unpack8(codeptr, pc+1)
			*getslot(slotptr, dstslot) = env.result2
			pc += 2

		case bytecode.OpLen:
			dstslot, srcslot := unpack8x2(codeptr, pc+1)
			getslot(slotptr, dstslot).Scalar = getslot(slotptr, srcslot).Scalar
			pc += 3
		case bytecode.OpCap:
			dstslot, srcslot := unpack8x2(codeptr, pc+1)
			getslot(slotptr, dstslot).Scalar = getslot(slotptr, srcslot).Scalar2
			pc += 3

		case bytecode.OpStrIndex:
			dstslot, strslot, indexslot := unpack8x3(codeptr, pc+1)
			str := getslot(slotptr, strslot).String()
			index := getslot(slotptr, indexslot).Int()
			getslot(slotptr, dstslot).SetByte(str[index])
			pc += 4
		case bytecode.OpSliceIndexScalar8:
			dstslot, sliceslot, indexslot := unpack8x3(codeptr, pc+1)
			slice := getslot(slotptr, sliceslot).ByteSlice()
			index := getslot(slotptr, indexslot).Int()
			getslot(slotptr, dstslot).SetByte(slice[index])
			pc += 4
		case bytecode.OpSliceIndexScalar64:
			dstslot, sliceslot, indexslot := unpack8x3(codeptr, pc+1)
			slice := getslot(slotptr, sliceslot).slice64()
			index := getslot(slotptr, indexslot).Int()
			getslot(slotptr, dstslot).Scalar = slice[index]
			pc += 4

		case bytecode.OpSliceSetScalar8:
			sliceslot, indexslot, valueslot := unpack8x3(codeptr, pc+1)
			slice := getslot(slotptr, sliceslot).ByteSlice()
			index := getslot(slotptr, indexslot).Int()
			slice[index] = byte(getslot(slotptr, valueslot).Scalar)
			pc += 4
		case bytecode.OpSliceSetScalar64:
			sliceslot, indexslot, valueslot := unpack8x3(codeptr, pc+1)
			slice := getslot(slotptr, sliceslot).slice64()
			index := getslot(slotptr, indexslot).Int()
			slice[index] = getslot(slotptr, valueslot).Scalar
			pc += 4

		case bytecode.OpStrSliceFrom:
			dstslot, strslot, fromslot := unpack8x3(codeptr, pc+1)
			str := getslot(slotptr, strslot).String()
			from := getslot(slotptr, fromslot).Scalar
			getslot(slotptr, dstslot).SetString(str[from:])
			pc += 4
		case bytecode.OpStrSliceTo:
			dstslot, strslot, toslot := unpack8x3(codeptr, pc+1)
			str := getslot(slotptr, strslot).String()
			to := getslot(slotptr, toslot).Scalar
			getslot(slotptr, dstslot).SetString(str[:to])
			pc += 4
		case bytecode.OpStrSlice:
			dstslot, strslot, fromslot, toslot := unpack8x4(codeptr, pc+1)
			str := getslot(slotptr, strslot).String()
			from := getslot(slotptr, fromslot).Scalar
			to := getslot(slotptr, toslot).Scalar
			getslot(slotptr, dstslot).SetString(str[from:to])
			pc += 5

		case bytecode.OpNot:
			dstslot, srcslot := unpack8x2(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(!getslot(slotptr, srcslot).Bool())
			pc += 3

		case bytecode.OpIsNil:
			dstslot, srcslot := unpack8x2(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(getslot(slotptr, srcslot).IsNil())
			pc += 3
		case bytecode.OpIsNotNil:
			dstslot, srcslot := unpack8x2(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(!getslot(slotptr, srcslot).IsNil())
			pc += 3
		case bytecode.OpIsNilInterface:
			dstslot, srcslot := unpack8x2(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(getslot(slotptr, srcslot).IsNilInterface())
			pc += 3
		case bytecode.OpIsNotNilInterface:
			dstslot, srcslot := unpack8x2(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(!getslot(slotptr, srcslot).IsNilInterface())
			pc += 3

		case bytecode.OpStrEq:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(getslot(slotptr, xslot).String() == getslot(slotptr, yslot).String())
			pc += 4
		case bytecode.OpStrNotEq:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(getslot(slotptr, xslot).String() != getslot(slotptr, yslot).String())
			pc += 4
		case bytecode.OpStrGt:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(getslot(slotptr, xslot).String() > getslot(slotptr, yslot).String())
			pc += 4
		case bytecode.OpStrLt:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(getslot(slotptr, xslot).String() < getslot(slotptr, yslot).String())
			pc += 4

		case bytecode.OpIntNeg:
			dstslot, xslot := unpack8x2(codeptr, pc+1)
			getslot(slotptr, dstslot).SetInt(-getslot(slotptr, xslot).Int())
			pc += 3

		case bytecode.OpScalarEq:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(getslot(slotptr, xslot).Scalar == getslot(slotptr, yslot).Scalar)
			pc += 4
		case bytecode.OpScalarNotEq:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(getslot(slotptr, xslot).Scalar != getslot(slotptr, yslot).Scalar)
			pc += 4
		case bytecode.OpIntGt:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(getslot(slotptr, xslot).Int() > getslot(slotptr, yslot).Int())
			pc += 4
		case bytecode.OpIntGtEq:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(getslot(slotptr, xslot).Int() >= getslot(slotptr, yslot).Int())
			pc += 4
		case bytecode.OpIntLt:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(getslot(slotptr, xslot).Int() < getslot(slotptr, yslot).Int())
			pc += 4
		case bytecode.OpIntLtEq:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).SetBool(getslot(slotptr, xslot).Int() <= getslot(slotptr, yslot).Int())
			pc += 4

		case bytecode.OpConcat:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).SetString(getslot(slotptr, xslot).String() + getslot(slotptr, yslot).String())
			pc += 4

		case bytecode.OpIntXor:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).SetInt(getslot(slotptr, xslot).Int() ^ getslot(slotptr, yslot).Int())
			pc += 4
		case bytecode.OpIntAdd8:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).SetByte(getslot(slotptr, xslot).Byte() + getslot(slotptr, yslot).Byte())
			pc += 4
		case bytecode.OpIntAdd64:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).Scalar = getslot(slotptr, xslot).Scalar + getslot(slotptr, yslot).Scalar
			pc += 4
		case bytecode.OpIntSub8:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).SetByte(getslot(slotptr, xslot).Byte() - getslot(slotptr, yslot).Byte())
			pc += 4
		case bytecode.OpIntSub64:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).Scalar = getslot(slotptr, xslot).Scalar - getslot(slotptr, yslot).Scalar
			pc += 4
		case bytecode.OpIntMul8:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).SetByte(getslot(slotptr, xslot).Byte() * getslot(slotptr, yslot).Byte())
			pc += 4
		case bytecode.OpIntMul64:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).Scalar = getslot(slotptr, xslot).Scalar * getslot(slotptr, yslot).Scalar
			pc += 4
		case bytecode.OpIntDiv:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).SetInt(getslot(slotptr, xslot).Int() / getslot(slotptr, yslot).Int())
			pc += 4
		case bytecode.OpIntMod:
			dstslot, xslot, yslot := unpack8x3(codeptr, pc+1)
			getslot(slotptr, dstslot).SetInt(getslot(slotptr, xslot).Int() % getslot(slotptr, yslot).Int())
			pc += 4

		case bytecode.OpIntInc:
			dstslot := unpack8(codeptr, pc+1)
			getslot(slotptr, dstslot).Scalar++
			pc += 2
		case bytecode.OpIntDec:
			dstslot := unpack8(codeptr, pc+1)
			getslot(slotptr, dstslot).Scalar--
			pc += 2

		case bytecode.OpJump:
			offset := unpack16(codeptr, pc+1)
			pc += offset

		case bytecode.OpJumpTable:
			slot := unpack8(codeptr, pc+1)
			pc += (getslot(slotptr, slot).Int() * 3) + 2

		case bytecode.OpJumpZero:
			srcslot := unpack8(codeptr, pc+3)
			if getslot(slotptr, srcslot).Scalar == 0 {
				offset := unpack16(codeptr, pc+1)
				pc += offset
			} else {
				pc += 4
			}
		case bytecode.OpJumpNotZero:
			srcslot := unpack8(codeptr, pc+3)
			if getslot(slotptr, srcslot).Scalar != 0 {
				offset := unpack16(codeptr, pc+1)
				pc += offset
			} else {
				pc += 4
			}

		case bytecode.OpCall:
			dstslot := unpack8(codeptr, pc+1)
			funcid := unpack16(codeptr, pc+2)
			callFunc := env.userFuncs[funcid]
			if !canAllocFrame(slotptr, env.slotend, callFunc.FrameSize) {
				panicStackOverflow(fn)
			}
			eval(env, callFunc, nextFrameSlot(slotptr, fn.FrameSize))
			*getslot(slotptr, dstslot) = env.result
			pc += 4
		case bytecode.OpCallRecur:
			dstslot := unpack8(codeptr, pc+1)
			if !canAllocFrame(slotptr, env.slotend, fn.FrameSize) {
				panicStackOverflow(fn)
			}
			eval(env, fn, nextFrameSlot(slotptr, fn.FrameSize))
			*getslot(slotptr, dstslot) = env.result
			pc += 2
		case bytecode.OpCallVoid:
			funcid := unpack16(codeptr, pc+1)
			callFunc := env.userFuncs[funcid]
			if !canAllocFrame(slotptr, env.slotend, callFunc.FrameSize) {
				panicStackOverflow(fn)
			}
			eval(env, callFunc, nextFrameSlot(slotptr, fn.FrameSize))
			pc += 3

		case bytecode.OpVariadicReset:
			env.vararg = env.vararg[:0]
			pc++
		case bytecode.OpPushVariadicBoolArg:
			srcslot := unpack8(codeptr, pc+1)
			env.vararg = append(env.vararg, getslot(slotptr, srcslot).Bool())
			pc += 2
		case bytecode.OpPushVariadicScalarArg:
			srcslot := unpack8(codeptr, pc+1)
			env.vararg = append(env.vararg, getslot(slotptr, srcslot).Int())
			pc += 2
		case bytecode.OpPushVariadicStrArg:
			srcslot := unpack8(codeptr, pc+1)
			env.vararg = append(env.vararg, getslot(slotptr, srcslot).String())
			pc += 2
		case bytecode.OpPushVariadicInterfaceArg:
			srcslot := unpack8(codeptr, pc+1)
			env.vararg = append(env.vararg, getslot(slotptr, srcslot).Interface())
			pc += 2

		case bytecode.OpCallNative:
			dstslot := unpack8(codeptr, pc+1)
			funcid := unpack16(codeptr, pc+2)
			callFunc := env.nativeFuncs[funcid]
			if !canAllocFrame(slotptr, env.slotend, callFunc.frameSize) {
				panicStackOverflow(fn)
			}
			callFunc.mappedFunc(NativeCallContext{
				env:     env,
				slotptr: nextFrameSlot(slotptr, fn.FrameSize),
			})
			*getslot(slotptr, dstslot) = env.result
			pc += 4
		case bytecode.OpCallVoidNative:
			funcid := unpack16(codeptr, pc+1)
			callFunc := env.nativeFuncs[funcid]
			if !canAllocFrame(slotptr, env.slotend, callFunc.frameSize) {
				panicStackOverflow(fn)
			}
			callFunc.mappedFunc(NativeCallContext{
				slotptr: nextFrameSlot(slotptr, fn.FrameSize),
			})
			pc += 3

		case bytecode.OpReturnVoid:
			return
		case bytecode.OpReturnZero:
			env.result.Scalar = 0
			return
		case bytecode.OpReturnOne:
			env.result.Scalar = 1
			return
		case bytecode.OpReturnStr:
			srcslot := unpack8(codeptr, pc+1)
			env.result.SetString(getslot(slotptr, srcslot).String())
			return
		case bytecode.OpReturnScalar:
			srcslot := unpack8(codeptr, pc+1)
			env.result.Scalar = getslot(slotptr, srcslot).Scalar
			return
		case bytecode.OpReturn:
			srcslot := unpack8(codeptr, pc+1)
			env.result = *getslot(slotptr, srcslot)
			return

		default:
			panic(fmt.Sprintf("malformed bytecode: unexpected %s found", op))
		}
	}
}
