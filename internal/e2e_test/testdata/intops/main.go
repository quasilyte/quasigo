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

func main() {
	testMax()
	testGCD()
}
