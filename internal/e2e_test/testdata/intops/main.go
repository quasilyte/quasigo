package main

//test:irdump
// block0 [0]:
//   Jump L1
// block1 (L2) [0]:
// block2 (L1) [0]:
//   Move temp0.v0 = b
//   IntMod b = a b
//   Move a = temp0.v0
// block3 [0]:
//   Zero temp1.v0
//   ScalarNotEq temp0.v1 = b temp1.v0
//   JumpNotZero L2 temp0.v1
// block4 (L0) [0]:
//   ReturnScalar a
//
//test:disasm
// main.gcd code=25 frame=96 (4 slots: 2 params, 2 locals)
//   Jump L0
// L1:
//   Move temp0 = b
//   IntMod b = a b
//   Move a = temp0
// L0:
//   Zero temp1
//   ScalarNotEq temp0 = b temp1
//   JumpNotZero L1 temp0
//   ReturnScalar a
//
//test:disasm_opt
// main.gcd code=19 frame=72 (3 slots: 2 params, 1 locals)
//   Jump L0
// L1:
//   Move temp0 = b
//   IntMod b = a b
//   Move a = temp0
// L0:
//   JumpNotZero L1 b
//   ReturnScalar a
func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

//test:disasm_both
// main.max code=12 frame=72 (3 slots: 2 params, 1 locals)
//   IntGt temp0 = a b
//   JumpZero L0 temp0
//   ReturnScalar a
// L0:
//   ReturnScalar b
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

//test:disasm_both
// main.sqrt code=101 frame=168 (7 slots: 1 params, 6 locals)
//   Zero temp1
//   ScalarEq temp0 = x temp1
//   JumpNotZero L0 temp0
//   LoadScalarConst temp2 = 1
//   ScalarEq temp0 = x temp2
// L0:
//   JumpZero L1 temp0
//   ReturnScalar x
// L1:
//   LoadScalarConst temp0 = 1
//   LoadScalarConst temp2 = 2
//   IntDiv temp1 = x temp2
//   Zero temp2
//   Jump L2
// L5:
//   IntAdd64 temp4 = temp0 temp1
//   LoadScalarConst temp5 = 2
//   IntDiv temp3 = temp4 temp5
//   IntMul64 temp4 = temp3 temp3
//   ScalarEq temp5 = temp4 x
//   JumpZero L3 temp5
//   ReturnScalar temp3
// L3:
//   IntLtEq temp5 = temp4 x
//   JumpZero L4 temp5
//   LoadScalarConst temp5 = 1
//   IntAdd64 temp0 = temp3 temp5
//   Move temp2 = temp3
//   Jump L2
// L4:
//   LoadScalarConst temp5 = 1
//   IntSub64 temp1 = temp3 temp5
// L2:
//   IntLtEq temp3 = temp0 temp1
//   JumpNotZero L5 temp3
//   ReturnScalar temp2
func sqrt(x int) int {
	if x == 0 || x == 1 {
		return x
	}
	start := 1
	end := x / 2
	result := 0
	for start <= end {
		mid := (start + end) / 2
		sqr := mid * mid
		if sqr == x {
			return mid
		}
		if sqr <= x {
			start = mid + 1
			result = mid
		} else {
			end = mid - 1
		}
	}
	return result
}

func evalA(i, j int) int { return ((i+j)*(i+j+1)/2 + i + 1) }

func testMax() {
	println(max(0, 0))
	println(max(4, 0))
	println(max(0, 3))
	println(max(0, -3))
	println(max(-4, -3))
}

func testGCD() {
	for a := -2; a < 20; a++ {
		for b := -6; b < 30; b += 3 {
			println(gcd(a, b))
		}
	}
}

func testSqrt() {
	println(sqrt(0))
	println(sqrt(1))
	println(sqrt(10))
	println(sqrt(15))
	println(sqrt(219))
	println(sqrt(2000))
	println(sqrt(36))
	println(sqrt(48))
	println(sqrt(81))
	println(sqrt(1024))
	println(sqrt(1025))
	println(sqrt(8321))
	println(sqrt(9999))
}

func testEvalA() {
	for i := -10; i < 25; i++ {
		for j := -10; j < 25; j++ {
			println(evalA(i, j))
		}
	}
}

func main() {
	testMax()
	testGCD()
	testSqrt()
	testEvalA()
}
