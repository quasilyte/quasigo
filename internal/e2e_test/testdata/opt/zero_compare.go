package opttest

//test:disasm_opt
// opttest.zeroCompare1 code=14 frame=48 (2 slots: 1 params, 1 locals)
//   JumpZero L0 x
//   LoadStrConst temp0 = "a"
//   ReturnStr temp0
// L0:
//   LoadStrConst temp0 = "b"
//   ReturnStr temp0
func zeroCompare1(x int) string {
	if x != 0 {
		return "a"
	}
	return "b"
}

//test:disasm_opt
// opttest.zeroCompare2 code=17 frame=72 (3 slots: 1 params, 2 locals)
//   Len temp1 = s
//   JumpZero L0 temp1
//   LoadStrConst temp0 = "nonzero"
//   ReturnStr temp0
// L0:
//   LoadStrConst temp0 = "zero"
//   ReturnStr temp0
func zeroCompare2(s string) string {
	if len(s) != 0 {
		return "nonzero"
	}
	return "zero"
}

//test:disasm_opt
// opttest.zeroCompare3 code=17 frame=72 (3 slots: 1 params, 2 locals)
//   Len temp1 = s
//   JumpNotZero L0 temp1
//   LoadStrConst temp0 = "zero"
//   ReturnStr temp0
// L0:
//   LoadStrConst temp0 = "nonzero"
//   ReturnStr temp0
func zeroCompare3(s string) string {
	if len(s) == 0 {
		return "zero"
	}
	return "nonzero"
}

//test:disasm_opt
// opttest.zeroCompare4 code=17 frame=72 (3 slots: 1 params, 2 locals)
//   Len temp1 = s
//   JumpNotZero L0 temp1
//   LoadStrConst temp0 = "nonzero"
//   ReturnStr temp0
// L0:
//   LoadStrConst temp0 = "zero"
//   ReturnStr temp0
func zeroCompare4(s string) string {
	if !(len(s) != 0) {
		return "nonzero"
	}
	return "zero"
}
