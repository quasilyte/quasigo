// Code generated by "stringer -type=SlotKind -trimprefix=Slot"; DO NOT EDIT.

package ir

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SlotInvalid-0]
	_ = x[SlotCallArg-1]
	_ = x[SlotParam-2]
	_ = x[SlotTemp-3]
	_ = x[SlotUniq-4]
	_ = x[SlotDiscard-5]
}

const _SlotKind_name = "InvalidCallArgParamTempUniqDiscard"

var _SlotKind_index = [...]uint8{0, 7, 14, 19, 23, 27, 34}

func (i SlotKind) String() string {
	if i >= SlotKind(len(_SlotKind_index)-1) {
		return "SlotKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SlotKind_name[_SlotKind_index[i]:_SlotKind_index[i+1]]
}
