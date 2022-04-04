package main

//test:disasm
// main.fx3 code=9 frame=72 (3 slots: 1 params, 2 locals)
//   LoadScalarConst temp1 = 1
//   IntAdd64 temp0 = x temp1
//   ReturnScalar temp0
func fx3(x int) int { return x + 1 }

//test:disasm_opt
// main.fx2 code=32 frame=168 (7 slots: 0 params, 7 locals)
//   LoadScalarConst temp2 = 1
//   LoadScalarConst temp5 = 1
//   IntAdd64 temp4 = temp2 temp5
//   Move temp1 = temp4
//   LoadScalarConst temp3 = 2
//   LoadScalarConst temp6 = 1
//   IntAdd64 temp5 = temp3 temp6
//   Move temp2 = temp5
//   IntAdd64 temp0 = temp1 temp2
//   ReturnScalar temp0
func fx2() int { return fx3(1) + fx3(2) }

//test:disasm_opt
// main.fx1 code=35 frame=192 (8 slots: 0 params, 8 locals)
//   LoadScalarConst temp3 = 1
//   LoadScalarConst temp6 = 1
//   IntAdd64 temp5 = temp3 temp6
//   Move temp2 = temp5
//   LoadScalarConst temp4 = 2
//   LoadScalarConst temp7 = 1
//   IntAdd64 temp6 = temp4 temp7
//   Move temp3 = temp6
//   IntAdd64 temp1 = temp2 temp3
//   Move temp0 = temp1
//   ReturnScalar temp0
func fx1() int { return fx2() }

//test:disasm_opt
// main.fx code=38 frame=216 (9 slots: 0 params, 9 locals)
//   LoadScalarConst temp4 = 1
//   LoadScalarConst temp7 = 1
//   IntAdd64 temp6 = temp4 temp7
//   Move temp3 = temp6
//   LoadScalarConst temp5 = 2
//   LoadScalarConst temp8 = 1
//   IntAdd64 temp7 = temp5 temp8
//   Move temp4 = temp7
//   IntAdd64 temp2 = temp3 temp4
//   Move temp1 = temp2
//   Move temp0 = temp1
//   ReturnScalar temp0
func fx() int { return fx1() }

//test:disasm
// main.isDigit code=20 frame=96 (4 slots: 1 params, 3 locals)
//   LoadScalarConst temp1 = 48
//   IntGtEq temp0 = ch temp1
//   JumpZero L0 temp0
//   LoadScalarConst temp2 = 57
//   IntLtEq temp0 = ch temp2
// L0:
//   ReturnScalar temp0
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

//test:disasm
// main.isAlpha code=20 frame=96 (4 slots: 1 params, 3 locals)
//   LoadScalarConst temp1 = 97
//   IntGtEq temp0 = ch temp1
//   JumpZero L0 temp0
//   LoadScalarConst temp2 = 122
//   IntLtEq temp0 = ch temp2
// L0:
//   ReturnScalar temp0
func isAlpha(ch byte) bool {
	return ch >= 'a' && ch <= 'z'
}

//test:disasm_opt
// main.isAlphaNum code=54 frame=168 (7 slots: 1 params, 6 locals)
//   Move temp1 = ch
//   LoadScalarConst temp4 = 48
//   IntGtEq temp3 = temp1 temp4
//   JumpZero L0 temp3
//   LoadScalarConst temp5 = 57
//   IntLtEq temp3 = temp1 temp5
// L0:
//   Move temp0 = temp3
//   JumpNotZero L1 temp0
//   Move temp1 = ch
//   LoadScalarConst temp4 = 97
//   IntGtEq temp3 = temp1 temp4
//   JumpZero L2 temp3
//   LoadScalarConst temp5 = 122
//   IntLtEq temp3 = temp1 temp5
// L2:
//   Move temp0 = temp3
// L1:
//   ReturnScalar temp0
func isAlphaNum(ch byte) bool {
	return isDigit(ch) || isAlpha(ch)
}

//test:disasm
// main.inlNot code=5 frame=48 (2 slots: 1 params, 1 locals)
//   Not temp0 = b
//   ReturnScalar temp0
func inlNot(b bool) bool {
	return !b
}

//test:disasm_opt
// main.testInlNot code=11 frame=120 (5 slots: 1 params, 4 locals)
//   Move temp1 = b
//   Not temp3 = temp1
//   Move temp0 = temp3
//   ReturnScalar temp0
func testInlNot(b bool) bool {
	return inlNot(b)
}

//test:disasm
// main.inlMultiStmt code=14 frame=48 (2 slots: 1 params, 1 locals)
//   JumpZero L0 b
//   LoadScalarConst temp0 = 10
//   ReturnScalar temp0
// L0:
//   LoadScalarConst temp0 = 20
//   ReturnScalar temp0
func inlMultiStmt(b bool) int {
	if b {
		return 10
	}
	return 20
}

//test:disasm_opt
// main.testInlMultiStmt code=24 frame=120 (5 slots: 1 params, 4 locals)
//   Move temp1 = b
//   JumpZero L0 temp1
//   LoadScalarConst temp3 = 10
//   Move temp0 = temp3
//   Jump L1
// L0:
//   LoadScalarConst temp3 = 20
//   Move temp0 = temp3
// L1:
//   ReturnScalar temp0
func testInlMultiStmt(b bool) int {
	return inlMultiStmt(b)
}

//test:disasm
// main.inlLocal code=7 frame=24 (1 slots: 0 params, 1 locals)
//   LoadScalarConst temp0 = 10
//   IntInc temp0
//   ReturnScalar temp0
func inlLocal() int {
	loc := 10
	loc++
	return loc
}

//test:disasm_opt
// main.testInlLocal code=17 frame=72 (3 slots: 0 params, 3 locals)
//   LoadScalarConst temp2 = 10
//   IntInc temp2
//   Move temp1 = temp2
//   LoadScalarConst temp2 = 1
//   IntAdd64 temp0 = temp1 temp2
//   ReturnScalar temp0
func testInlLocal() int {
	return inlLocal() + 1
}

//test:disasm
// main.inlSwitch code=36 frame=96 (4 slots: 1 params, 3 locals)
//   Move temp0 = x
//   LoadScalarConst temp2 = 10
//   ScalarEq temp1 = temp0 temp2
//   JumpZero L0 temp1
//   ReturnOne
// L0:
//   LoadScalarConst temp2 = 20
//   ScalarEq temp1 = temp0 temp2
//   JumpZero L1 temp1
//   LoadScalarConst temp1 = -1
//   ReturnScalar temp1
// L1:
//   LoadScalarConst temp1 = 100
//   ReturnScalar temp1
func inlSwitch(x int) int {
	switch x {
	case 10:
		return 1
	case 20:
		return -1
	default:
		return 100
	}
}

//test:disasm_opt
// main.testInlSwitch code=51 frame=168 (7 slots: 1 params, 6 locals)
//   Move temp1 = x
//   Move temp3 = temp1
//   LoadScalarConst temp5 = 10
//   ScalarEq temp4 = temp3 temp5
//   JumpZero L0 temp4
//   LoadScalarConst temp0 = 1
//   Jump L1
// L0:
//   LoadScalarConst temp5 = 20
//   ScalarEq temp4 = temp3 temp5
//   JumpZero L2 temp4
//   LoadScalarConst temp4 = -1
//   Move temp0 = temp4
//   Jump L1
// L2:
//   LoadScalarConst temp4 = 100
//   Move temp0 = temp4
// L1:
//   ReturnScalar temp0
func testInlSwitch(x int) int {
	return inlSwitch(x)
}

func testChar(ch byte) {
	println(isAlpha(ch))
	println(isDigit(ch))
	println(isAlphaNum(ch))
}

func testBool(b bool) {
	println(testInlMultiStmt(b))
	println(testInlNot(b))
}

func intneg(x int) int {
	return -x
}

func noinline(x, y int) {
	for i := 0; i < 0; i++ {
	}
}

//test:irdump
// block0 [0]:
//   LoadScalarConst temp0.v0 = 3
//   LoadScalarConst temp1.v0 = -1
//   Move arg0 = temp0.v0
//   Move arg1 = temp1.v0
//   CallVoid main.noinline
//   Zero temp1
//   IntNeg temp3 = temp1
//   Move temp0 = temp3
// block1 (L0) [1]:
//   Move arg0 = temp0
//   CallVoidNative builtin.PrintInt
//   VarKill temp0
//   ReturnVoid
//
//test:disasm_opt
// main.testIntNeg code=24 frame=72 (3 slots: 0 params, 3 locals)
//   LoadScalarConst arg0 = 3
//   LoadScalarConst arg1 = -1
//   CallVoid main.noinline()
//   Zero temp1
//   IntNeg temp2 = temp1
//   Move temp0 = temp2
//   Move arg0 = temp0
//   CallVoidNative builtin.PrintInt()
//   ReturnVoid
func testIntNeg() {
	noinline(3, -1)
	println(intneg(0))
}

func main() {
	for i := -15; i < 40; i++ {
		println(testInlSwitch(i))
	}

	println(fx())
	println(testInlLocal())

	testChar('a')
	testChar('b')
	testChar('z')
	testChar('_')
	testChar('$')
	testChar(' ')
	testChar('0')
	testChar('3')
	testChar('9')
	testChar(0)

	testBool(true)
	testBool(false)

	testIntNeg()
}
