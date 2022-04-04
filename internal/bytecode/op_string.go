// Code generated by "stringer -type=Op -trimprefix=Op"; DO NOT EDIT.

package bytecode

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OpInvalid-0]
	_ = x[OpLoadScalarConst-1]
	_ = x[OpLoadStrConst-2]
	_ = x[OpZero-3]
	_ = x[OpMove-4]
	_ = x[OpMove8-5]
	_ = x[OpMoveResult2-6]
	_ = x[OpNot-7]
	_ = x[OpIsNil-8]
	_ = x[OpIsNotNil-9]
	_ = x[OpIsNilInterface-10]
	_ = x[OpIsNotNilInterface-11]
	_ = x[OpLen-12]
	_ = x[OpCap-13]
	_ = x[OpStrSlice-14]
	_ = x[OpStrSliceFrom-15]
	_ = x[OpStrSliceTo-16]
	_ = x[OpStrIndex-17]
	_ = x[OpSliceIndexScalar8-18]
	_ = x[OpSliceIndexScalar64-19]
	_ = x[OpBytesSlice-20]
	_ = x[OpBytesSliceFrom-21]
	_ = x[OpBytesSliceTo-22]
	_ = x[OpSliceSetScalar8-23]
	_ = x[OpSliceSetScalar64-24]
	_ = x[OpConcat-25]
	_ = x[OpStrEq-26]
	_ = x[OpStrNotEq-27]
	_ = x[OpStrGt-28]
	_ = x[OpStrLt-29]
	_ = x[OpIntNeg-30]
	_ = x[OpIntBitwiseNot-31]
	_ = x[OpScalarEq-32]
	_ = x[OpScalarNotEq-33]
	_ = x[OpIntGt-34]
	_ = x[OpIntGtEq-35]
	_ = x[OpIntLt-36]
	_ = x[OpIntLtEq-37]
	_ = x[OpIntAdd8-38]
	_ = x[OpIntAdd64-39]
	_ = x[OpIntSub8-40]
	_ = x[OpIntSub64-41]
	_ = x[OpIntMul8-42]
	_ = x[OpIntMul64-43]
	_ = x[OpIntDiv-44]
	_ = x[OpIntMod-45]
	_ = x[OpIntXor-46]
	_ = x[OpIntOr-47]
	_ = x[OpIntLshift-48]
	_ = x[OpIntRshift-49]
	_ = x[OpIntInc-50]
	_ = x[OpIntDec-51]
	_ = x[OpJump-52]
	_ = x[OpJumpZero-53]
	_ = x[OpJumpNotZero-54]
	_ = x[OpJumpTable-55]
	_ = x[OpCall-56]
	_ = x[OpCallRecur-57]
	_ = x[OpCallVoid-58]
	_ = x[OpCallNative-59]
	_ = x[OpCallVoidNative-60]
	_ = x[OpPushVariadicBoolArg-61]
	_ = x[OpPushVariadicScalarArg-62]
	_ = x[OpPushVariadicStrArg-63]
	_ = x[OpPushVariadicInterfaceArg-64]
	_ = x[OpVariadicReset-65]
	_ = x[OpReturnVoid-66]
	_ = x[OpReturnZero-67]
	_ = x[OpReturnOne-68]
	_ = x[OpReturnStr-69]
	_ = x[OpReturnScalar-70]
	_ = x[OpReturn-71]
	_ = x[OpFloatAdd64-72]
	_ = x[OpFloatSub64-73]
	_ = x[OpFloatMul64-74]
	_ = x[OpFloatDiv64-75]
	_ = x[OpFloatGt-76]
	_ = x[OpFloatGtEq-77]
	_ = x[OpFloatLt-78]
	_ = x[OpFloatLtEq-79]
	_ = x[OpFloatNeg-80]
	_ = x[OpConvIntToFloat-81]
}

const _Op_name = "InvalidLoadScalarConstLoadStrConstZeroMoveMove8MoveResult2NotIsNilIsNotNilIsNilInterfaceIsNotNilInterfaceLenCapStrSliceStrSliceFromStrSliceToStrIndexSliceIndexScalar8SliceIndexScalar64BytesSliceBytesSliceFromBytesSliceToSliceSetScalar8SliceSetScalar64ConcatStrEqStrNotEqStrGtStrLtIntNegIntBitwiseNotScalarEqScalarNotEqIntGtIntGtEqIntLtIntLtEqIntAdd8IntAdd64IntSub8IntSub64IntMul8IntMul64IntDivIntModIntXorIntOrIntLshiftIntRshiftIntIncIntDecJumpJumpZeroJumpNotZeroJumpTableCallCallRecurCallVoidCallNativeCallVoidNativePushVariadicBoolArgPushVariadicScalarArgPushVariadicStrArgPushVariadicInterfaceArgVariadicResetReturnVoidReturnZeroReturnOneReturnStrReturnScalarReturnFloatAdd64FloatSub64FloatMul64FloatDiv64FloatGtFloatGtEqFloatLtFloatLtEqFloatNegConvIntToFloat"

var _Op_index = [...]uint16{0, 7, 22, 34, 38, 42, 47, 58, 61, 66, 74, 88, 105, 108, 111, 119, 131, 141, 149, 166, 184, 194, 208, 220, 235, 251, 257, 262, 270, 275, 280, 286, 299, 307, 318, 323, 330, 335, 342, 349, 357, 364, 372, 379, 387, 393, 399, 405, 410, 419, 428, 434, 440, 444, 452, 463, 472, 476, 485, 493, 503, 517, 536, 557, 575, 599, 612, 622, 632, 641, 650, 662, 668, 678, 688, 698, 708, 715, 724, 731, 740, 748, 762}

func (i Op) String() string {
	if i >= Op(len(_Op_index)-1) {
		return "Op(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Op_name[_Op_index[i]:_Op_index[i+1]]
}
