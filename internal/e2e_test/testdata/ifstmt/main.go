package main

//test:irdump
// block0 [0]:
//   JumpZero L0 cond
// block1 [0]:
//   ReturnScalar x
// block2 (L0) [0]:
//   ReturnScalar y
//
//test:disasm_both
// main.testIf1 code=8 frame=72 (3 slots: 3 params, 0 locals)
//   JumpZero L0 cond
//   ReturnScalar x
// L0:
//   ReturnScalar y
func testIf1(cond bool, x, y int) int {
	if cond {
		return x
	}
	return y
}

//test:irdump
// block0 [0]:
//   Zero temp0
//   JumpZero L1 cond
// block1 [0]:
//   Move temp0 = x
//   Jump L0
// block2 (L1) [0]:
//   Move temp0 = y
// block3 (L0) [1]:
//   ReturnScalar temp0
//   VarKill temp0
//
//test:disasm_both
// main.testIf2 code=17 frame=96 (4 slots: 3 params, 1 locals)
//   Zero temp0
//   JumpZero L0 cond
//   Move temp0 = x
//   Jump L1
// L0:
//   Move temp0 = y
// L1:
//   ReturnScalar temp0
func testIf2(cond bool, x, y int) int {
	result := 0
	if cond {
		result = x
	} else {
		result = y
	}
	return result
}

//test:irdump
// block0 [0]:
//   JumpZero L1 cond
// block1 [0]:
//   ReturnScalar x
// block2 (L1) [0]:
//   ReturnScalar y
// block3 [0]:
//
//test:disasm_both
// main.testIf3 code=8 frame=72 (3 slots: 3 params, 0 locals)
//   JumpZero L0 cond
//   ReturnScalar x
// L0:
//   ReturnScalar y
func testIf3(cond bool, x, y int) int {
	if cond {
		return x
	} else {
		return y
	}
}

//test:irdump
// block0 [0]:
//   Not temp0.v0 = cond
//   JumpZero L1 temp0.v0
// block1 [0]:
//   ReturnScalar x
// block2 (L1) [0]:
//   ReturnScalar y
// block3 [0]:
//
//test:disasm
// main.testIf4 code=11 frame=96 (4 slots: 3 params, 1 locals)
//   Not temp0 = cond
//   JumpZero L0 temp0
//   ReturnScalar x
// L0:
//   ReturnScalar y
//
//test:disasm_opt
// main.testIf4 code=8 frame=72 (3 slots: 3 params, 0 locals)
//   JumpNotZero L0 cond
//   ReturnScalar x
// L0:
//   ReturnScalar y
func testIf4(cond bool, x, y int) int {
	if !cond {
		return x
	} else {
		return y
	}
}

//test:irdump
// block0 [0]:
//   LoadScalarConst temp1.v0 = 2
//   ScalarNotEq temp0.v0 = x temp1.v0
//   JumpZero L1 temp0.v0
// block1 [0]:
//   LoadScalarConst temp0.v1 = 100
//   ReturnScalar temp0.v1
// block2 (L1) [0]:
//   JumpZero L2 cond
// block3 [0]:
//   IntAdd64 temp0.v2 = y y
//   ReturnScalar temp0.v2
// block4 (L2) [0]:
// block5 (L0) [0]:
//   LoadScalarConst temp0.v3 = 19
//   ReturnScalar temp0.v3
//
//test:disasm_opt
// main.testIf5 code=31 frame=120 (5 slots: 3 params, 2 locals)
//   LoadScalarConst temp1 = 2
//   ScalarNotEq temp0 = x temp1
//   JumpZero L0 temp0
//   LoadScalarConst temp0 = 100
//   ReturnScalar temp0
// L0:
//   JumpZero L1 cond
//   IntAdd64 temp0 = y y
//   ReturnScalar temp0
// L1:
//   LoadScalarConst temp0 = 19
//   ReturnScalar temp0
func testIf5(cond bool, x, y int) int {
	if x != 2 {
		return 100
	} else if cond {
		return y + y
	}
	return 19
}

func test2(cond bool, x, y int) {
	println(testIf1(cond, x, y))
	println(testIf2(cond, x, y))
	println(testIf3(cond, x, y))
	println(testIf4(cond, x, y))
	println(testIf5(cond, x, y))
}

func test(cond bool) {
	test2(cond, 10, 20)
	test2(cond, 20, 10)
	test2(cond, 0, 0)
	test2(cond, 2, 0)
	test2(cond, 0, 2)
}

func main() {
	test(true)
	test(false)
}
