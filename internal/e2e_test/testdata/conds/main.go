package main

//test:disasm_both
// main.notcond0 code=5 frame=48 (2 slots: 1 params, 1 locals)
//   Not temp0 = x
//   ReturnScalar temp0
func notcond0(x bool) bool {
	return !x
}

//test:disasm_both
// main.cond0 code=16 frame=96 (4 slots: 2 params, 2 locals)
//   Zero temp1
//   ScalarEq temp0 = x temp1
//   JumpNotZero L0 temp0
//   ScalarEq temp0 = y x
// L0:
//   ReturnScalar temp0
func cond0(x, y int) bool {
	return x == 0 || y == x
}

func cond1(x, y int) bool {
	return (x == 0 || x > 0) && (y < 5 || y >= 10)
}

func cond2(x, y int) bool {
	return (x != 0 || x < 0) || y < 5
}

func cond3(x, y int) bool {
	return x == 1 || x == 2 || y == 3 || y < 0
}

//test:disasm_both
// main.cond4 code=17 frame=96 (4 slots: 2 params, 2 locals)
//   LoadScalarConst temp1 = 2
//   ScalarEq temp0 = x temp1
//   JumpZero L0 temp0
//   ScalarEq temp0 = y x
// L0:
//   ReturnScalar temp0
func cond4(x, y int) bool {
	return x == 2 && y == x
}

func test0(x, y int) {
	println(cond0(x, y))
	println(cond0(y, x))
	println(cond0(x, x))
	println(cond0(y, y))
}

func test1(x, y int) {
	println(cond1(x, y))
	println(cond1(y, x))
	println(cond1(x, x))
	println(cond1(y, y))
}

func test2(x, y int) {
	println(cond2(x, y))
	println(cond2(y, x))
	println(cond2(x, x))
	println(cond2(y, y))
}

func test3(x, y int) {
	println(cond3(x, y))
	println(cond3(y, x))
	println(cond3(x, x))
	println(cond3(y, y))
}

func test4(x, y int) {
	println(cond4(x, y))
	println(cond4(y, x))
	println(cond4(x, x))
	println(cond4(y, y))
}

func testcond(x, y int) {
	test0(x, y)
	test1(x, y)
	test2(x, y)
	test3(x, y)
	test4(x, y)
}

func main() {
	testcond(-1, -1)
	testcond(-1, 0)
	testcond(1, 0)
	testcond(2, 0)
	testcond(2, 1)
	testcond(-2, 1)
	testcond(1031, 102)
	testcond(29, 10)
	testcond(-29, -10)
	testcond(-130, -130)
	testcond(0, -130)
	testcond(0, 130)
	testcond(10, 130)
	testcond(5, 10)

	println(notcond0(true))
	println(notcond0(false))
}
