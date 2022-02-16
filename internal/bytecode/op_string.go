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
	_ = x[OpMoveScalar-3]
	_ = x[OpMoveStr-4]
	_ = x[OpMoveInterface-5]
	_ = x[OpMoveResult2-6]
	_ = x[OpNot-7]
	_ = x[OpIsNil-8]
	_ = x[OpIsNotNil-9]
	_ = x[OpIsNilInterface-10]
	_ = x[OpIsNotNilInterface-11]
	_ = x[OpStrLen-12]
	_ = x[OpStrSlice-13]
	_ = x[OpStrSliceFrom-14]
	_ = x[OpStrSliceTo-15]
	_ = x[OpConcat-16]
	_ = x[OpStrEq-17]
	_ = x[OpStrNotEq-18]
	_ = x[OpIntEq-19]
	_ = x[OpIntNotEq-20]
	_ = x[OpIntGt-21]
	_ = x[OpIntGtEq-22]
	_ = x[OpIntLt-23]
	_ = x[OpIntLtEq-24]
	_ = x[OpIntAdd-25]
	_ = x[OpIntSub-26]
	_ = x[OpIntMul-27]
	_ = x[OpIntDiv-28]
	_ = x[OpIntInc-29]
	_ = x[OpIntDec-30]
	_ = x[OpJump-31]
	_ = x[OpJumpFalse-32]
	_ = x[OpJumpTrue-33]
	_ = x[OpCall-34]
	_ = x[OpCallRecur-35]
	_ = x[OpCallNative-36]
	_ = x[OpCallVoidNative-37]
	_ = x[OpPushVariadicBoolArg-38]
	_ = x[OpPushVariadicScalarArg-39]
	_ = x[OpPushVariadicStrArg-40]
	_ = x[OpPushVariadicInterfaceArg-41]
	_ = x[OpVariadicReset-42]
	_ = x[OpReturnVoid-43]
	_ = x[OpReturnFalse-44]
	_ = x[OpReturnTrue-45]
	_ = x[OpReturnStr-46]
	_ = x[OpReturnScalar-47]
	_ = x[OpReturnInterface-48]
}

const _Op_name = "InvalidLoadScalarConstLoadStrConstMoveScalarMoveStrMoveInterfaceMoveResult2NotIsNilIsNotNilIsNilInterfaceIsNotNilInterfaceStrLenStrSliceStrSliceFromStrSliceToConcatStrEqStrNotEqIntEqIntNotEqIntGtIntGtEqIntLtIntLtEqIntAddIntSubIntMulIntDivIntIncIntDecJumpJumpFalseJumpTrueCallCallRecurCallNativeCallVoidNativePushVariadicBoolArgPushVariadicScalarArgPushVariadicStrArgPushVariadicInterfaceArgVariadicResetReturnVoidReturnFalseReturnTrueReturnStrReturnScalarReturnInterface"

var _Op_index = [...]uint16{0, 7, 22, 34, 44, 51, 64, 75, 78, 83, 91, 105, 122, 128, 136, 148, 158, 164, 169, 177, 182, 190, 195, 202, 207, 214, 220, 226, 232, 238, 244, 250, 254, 263, 271, 275, 284, 294, 308, 327, 348, 366, 390, 403, 413, 424, 434, 443, 455, 470}

func (i Op) String() string {
	if i >= Op(len(_Op_index)-1) {
		return "Op(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Op_name[_Op_index[i]:_Op_index[i+1]]
}