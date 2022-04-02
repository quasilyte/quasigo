package opttest

//test:irdump
// block0 [0]:
//   Len temp2.v0 = s
//   Zero temp3.v0
//   ScalarEq temp1.v0 = temp2.v0 temp3.v0
//   Not temp0.v0 = temp1.v0
//   JumpZero L0 temp0.v0
// block1 [0]:
//   LoadStrConst temp0.v1 = "nonzero"
//   ReturnStr temp0.v1
// block2 (L0) [0]:
//   LoadStrConst temp0.v2 = "zero"
//   ReturnStr temp0.v2
//
// opttest.zerocmp1 code=17 frame=72 (3 slots: 1 params, 2 locals)
//   Len temp1 = s
//   JumpZero L0 temp1
//   LoadStrConst temp0 = "nonzero"
//   ReturnStr temp0
// L0:
//   LoadStrConst temp0 = "zero"
//   ReturnStr temp0
func zerocmp1(s string) string {
	if !(len(s) == 0) {
		return "nonzero"
	}
	return "zero"
}

//test:disasm_both
// opttest.nopStringSlice1 code=5 frame=48 (2 slots: 1 params, 1 locals)
//   Move temp0 = s
//   ReturnStr temp0
func nopStringSlice1(s string) string {
	return s[:]
}

//test:disasm_both
// opttest.nopStringSlice2 code=5 frame=48 (2 slots: 1 params, 1 locals)
//   Move temp0 = s
//   ReturnStr temp0
func nopStringSlice2(s string) string {
	return s[:][:][:]
}

//test:disasm_both
// opttest.nopStringSlice3 code=5 frame=48 (2 slots: 1 params, 1 locals)
//   Move temp0 = s
//   ReturnStr temp0
func nopStringSlice3(s string) string {
	return s[:len(s)]
}

//test:disasm_both
// opttest.nopStringSlice4 code=5 frame=48 (2 slots: 1 params, 1 locals)
//   Move temp0 = s
//   ReturnStr temp0
func nopStringSlice4(s string) string {
	return s[0:len(s)]
}

//test:disasm_both
// opttest.nopStringSlice5 code=5 frame=48 (2 slots: 1 params, 1 locals)
//   Move temp0 = s
//   ReturnStr temp0
func nopStringSlice5(s string) string {
	return s[0:]
}
