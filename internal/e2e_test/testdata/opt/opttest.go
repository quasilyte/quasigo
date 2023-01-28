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

//test:disasm_opt
// opttest.binaryopTempReuse code=21 frame=144 (6 slots: 2 params, 4 locals)
//   IntAdd64 temp2 = x y
//   LoadScalarConst temp3 = 2
//   IntAdd64 temp1 = temp2 temp3
//   IntAdd64 temp2 = x y
//   IntMul64 temp0 = temp1 temp2
//   ReturnScalar temp0
func binaryopTempReuse(x, y int) int {
	return (x + y + 2) * (x + y)
}

//test:disasm_opt
// opttest.spectralNormEvalA code=39 frame=240 (10 slots: 2 params, 8 locals)
//   IntAdd64 temp4 = i j
//   IntAdd64 temp6 = i j
//   LoadScalarConst temp7 = 1
//   IntAdd64 temp5 = temp6 temp7
//   IntMul64 temp3 = temp4 temp5
//   LoadScalarConst temp4 = 2
//   IntDiv temp2 = temp3 temp4
//   IntAdd64 temp1 = temp2 i
//   LoadScalarConst temp2 = 1
//   IntAdd64 temp0 = temp1 temp2
//   ReturnScalar temp0
func spectralNormEvalA(i, j int) int { return ((i+j)*(i+j+1)/2 + i + 1) }

//test:disasm_opt
// opttest.spectralNormEvalA2 code=32 frame=216 (9 slots: 2 params, 7 locals)
//   LoadScalarConst temp0 = 1
//   IntAdd64 temp1 = i j
//   IntAdd64 temp6 = temp1 temp0
//   IntMul64 temp5 = temp1 temp6
//   LoadScalarConst temp6 = 2
//   IntDiv temp4 = temp5 temp6
//   IntAdd64 temp3 = temp4 i
//   IntAdd64 temp2 = temp3 temp0
//   ReturnScalar temp2
func spectralNormEvalA2(i, j int) int {
	one := 1
	ij := i + j
	return (ij*(ij+one)/2 + i + one)
}
