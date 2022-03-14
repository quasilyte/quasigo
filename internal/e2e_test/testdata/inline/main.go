package main

//test:disasm
// main.fx3 code=9 frame=72 (3 slots: 1 args, 0 locals, 2 temps)
//   LoadScalarConst tmp1 = 1
//   IntAdd64 tmp0 = x tmp1
//   ReturnScalar tmp0
func fx3(x int) int { return x + 1 }

//test:disasm_opt
// main.fx2 code=32 frame=168 (7 slots: 0 args, 0 locals, 7 temps)
//   LoadScalarConst tmp2 = 1
//   LoadScalarConst tmp5 = 1
//   IntAdd64 tmp4 = tmp2 tmp5
//   Move tmp1 = tmp4
//   LoadScalarConst tmp3 = 2
//   LoadScalarConst tmp6 = 1
//   IntAdd64 tmp5 = tmp3 tmp6
//   Move tmp2 = tmp5
//   IntAdd64 tmp0 = tmp1 tmp2
//   ReturnScalar tmp0
func fx2() int { return fx3(1) + fx3(2) }

//test:disasm_opt
// main.fx1 code=35 frame=192 (8 slots: 0 args, 0 locals, 8 temps)
//   LoadScalarConst tmp3 = 1
//   LoadScalarConst tmp6 = 1
//   IntAdd64 tmp5 = tmp3 tmp6
//   Move tmp2 = tmp5
//   LoadScalarConst tmp4 = 2
//   LoadScalarConst tmp7 = 1
//   IntAdd64 tmp6 = tmp4 tmp7
//   Move tmp3 = tmp6
//   IntAdd64 tmp1 = tmp2 tmp3
//   Move tmp0 = tmp1
//   ReturnScalar tmp0
func fx1() int { return fx2() }

//test:disasm_opt
// main.fx code=38 frame=216 (9 slots: 0 args, 0 locals, 9 temps)
//   LoadScalarConst tmp4 = 1
//   LoadScalarConst tmp7 = 1
//   IntAdd64 tmp6 = tmp4 tmp7
//   Move tmp3 = tmp6
//   LoadScalarConst tmp5 = 2
//   LoadScalarConst tmp8 = 1
//   IntAdd64 tmp7 = tmp5 tmp8
//   Move tmp4 = tmp7
//   IntAdd64 tmp2 = tmp3 tmp4
//   Move tmp1 = tmp2
//   Move tmp0 = tmp1
//   ReturnScalar tmp0
func fx() int { return fx1() }

//test:disasm
// main.isDigit code=20 frame=96 (4 slots: 1 args, 0 locals, 3 temps)
//   LoadScalarConst tmp1 = 48
//   IntGtEq tmp0 = ch tmp1
//   JumpZero L0 tmp0
//   LoadScalarConst tmp2 = 57
//   IntLtEq tmp0 = ch tmp2
// L0:
//   ReturnScalar tmp0
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

//test:disasm
// main.isAlpha code=20 frame=96 (4 slots: 1 args, 0 locals, 3 temps)
//   LoadScalarConst tmp1 = 97
//   IntGtEq tmp0 = ch tmp1
//   JumpZero L0 tmp0
//   LoadScalarConst tmp2 = 122
//   IntLtEq tmp0 = ch tmp2
// L0:
//   ReturnScalar tmp0
func isAlpha(ch byte) bool {
	return ch >= 'a' && ch <= 'z'
}

//test:disasm_opt
// main.isAlphaNum code=54 frame=168 (7 slots: 1 args, 0 locals, 6 temps)
//   Move tmp1 = ch
//   LoadScalarConst tmp4 = 48
//   IntGtEq tmp3 = tmp1 tmp4
//   JumpZero L0 tmp3
//   LoadScalarConst tmp5 = 57
//   IntLtEq tmp3 = tmp1 tmp5
// L0:
//   Move tmp0 = tmp3
//   JumpNotZero L1 tmp0
//   Move tmp1 = ch
//   LoadScalarConst tmp4 = 97
//   IntGtEq tmp3 = tmp1 tmp4
//   JumpZero L2 tmp3
//   LoadScalarConst tmp5 = 122
//   IntLtEq tmp3 = tmp1 tmp5
// L2:
//   Move tmp0 = tmp3
// L1:
//   ReturnScalar tmp0
func isAlphaNum(ch byte) bool {
	return isDigit(ch) || isAlpha(ch)
}

//test:disasm
// main.inlNot code=5 frame=48 (2 slots: 1 args, 0 locals, 1 temps)
//   Not tmp0 = b
//   ReturnScalar tmp0
func inlNot(b bool) bool {
	return !b
}

//test:disasm_opt
// main.testInlNot code=11 frame=120 (5 slots: 1 args, 0 locals, 4 temps)
//   Move tmp1 = b
//   Not tmp3 = tmp1
//   Move tmp0 = tmp3
//   ReturnScalar tmp0
func testInlNot(b bool) bool {
	return inlNot(b)
}

//test:disasm
// main.inlMultiStmt code=14 frame=48 (2 slots: 1 args, 0 locals, 1 temps)
//   JumpZero L0 b
//   LoadScalarConst tmp0 = 10
//   ReturnScalar tmp0
// L0:
//   LoadScalarConst tmp0 = 20
//   ReturnScalar tmp0
func inlMultiStmt(b bool) int {
	if b {
		return 10
	}
	return 20
}

//test:disasm_opt
// main.testInlMultiStmt code=24 frame=120 (5 slots: 1 args, 0 locals, 4 temps)
//   Move tmp1 = b
//   JumpZero L0 tmp1
//   LoadScalarConst tmp3 = 10
//   Move tmp0 = tmp3
//   Jump L1
// L0:
//   LoadScalarConst tmp3 = 20
//   Move tmp0 = tmp3
// L1:
//   ReturnScalar tmp0
func testInlMultiStmt(b bool) int {
	return inlMultiStmt(b)
}

//test:disasm
// main.inlLocal code=7 frame=24 (1 slots: 0 args, 1 locals, 0 temps)
//   LoadScalarConst loc = 10
//   IntInc loc
//   ReturnScalar loc
func inlLocal() int {
	loc := 10
	loc++
	return loc
}

//test:disasm_opt
// main.testInlLocal code=17 frame=72 (3 slots: 0 args, 0 locals, 3 temps)
//   LoadScalarConst tmp2 = 10
//   IntInc tmp2
//   Move tmp1 = tmp2
//   LoadScalarConst tmp2 = 1
//   IntAdd64 tmp0 = tmp1 tmp2
//   ReturnScalar tmp0
func testInlLocal() int {
	return inlLocal() + 1
}

//test:disasm
// main.inlSwitch code=40 frame=96 (4 slots: 1 args, 1 locals, 2 temps)
//   Move auto0 = x
//   LoadScalarConst tmp1 = 10
//   ScalarEq tmp0 = auto0 tmp1
//   JumpZero L0 tmp0
//   LoadScalarConst tmp0 = 1
//   ReturnScalar tmp0
// L0:
//   LoadScalarConst tmp1 = 20
//   ScalarEq tmp0 = auto0 tmp1
//   JumpZero L1 tmp0
//   LoadScalarConst tmp0 = -1
//   ReturnScalar tmp0
// L1:
//   LoadScalarConst tmp0 = 100
//   ReturnScalar tmp0
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
// main.testInlSwitch code=54 frame=192 (8 slots: 1 args, 0 locals, 7 temps)
//   Move tmp1 = x
//   Move tmp4 = tmp1
//   LoadScalarConst tmp6 = 10
//   ScalarEq tmp5 = tmp4 tmp6
//   JumpZero L0 tmp5
//   LoadScalarConst tmp5 = 1
//   Move tmp0 = tmp5
//   Jump L1
// L0:
//   LoadScalarConst tmp6 = 20
//   ScalarEq tmp5 = tmp4 tmp6
//   JumpZero L2 tmp5
//   LoadScalarConst tmp5 = -1
//   Move tmp0 = tmp5
//   Jump L1
// L2:
//   LoadScalarConst tmp5 = 100
//   Move tmp0 = tmp5
// L1:
//   ReturnScalar tmp0
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
}
