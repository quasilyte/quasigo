package opttest

//test:irdump
// block0 [0]:
//   Move temp2.v0 = s
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
//   Move temp1 = s
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

//test:disasm_opt
// opttest.nopStringSlice1 code=2 frame=24 (1 slots: 1 params, 0 locals)
//   ReturnStr s
func nopStringSlice1(s string) string {
	return s[:]
}

//test:disasm_opt
// opttest.nopStringSlice2 code=2 frame=24 (1 slots: 1 params, 0 locals)
//   ReturnStr s
func nopStringSlice2(s string) string {
	return s[:][:][:]
}

//test:disasm_opt
// opttest.nopStringSlice3 code=2 frame=24 (1 slots: 1 params, 0 locals)
//   ReturnStr s
func nopStringSlice3(s string) string {
	return s[:len(s)]
}

//test:disasm_opt
// opttest.nopStringSlice4 code=2 frame=24 (1 slots: 1 params, 0 locals)
//   ReturnStr s
func nopStringSlice4(s string) string {
	return s[0:len(s)]
}

//test:disasm_opt
// opttest.nopStringSlice5 code=2 frame=24 (1 slots: 1 params, 0 locals)
//   ReturnStr s
func nopStringSlice5(s string) string {
	return s[0:]
}

//test:disasm_opt
// opttest.addByteAsInt1 code=6 frame=72 (3 slots: 2 params, 1 locals)
//   IntAdd64 temp0 = b i
//   ReturnScalar temp0
func addByteAsInt1(b byte, i int) int {
	return int(b) + i
}

//test:disasm_opt
// opttest.addByteAsInt2 code=6 frame=72 (3 slots: 2 params, 1 locals)
//   IntAdd64 temp0 = i b
//   ReturnScalar temp0
func addByteAsInt2(b byte, i int) int {
	return i + int(b)
}
