package irtest

//test:irdump
// block0 [0]:
//   Move temp0 = x
//   JumpNotZero L0 temp0
// block1 [0]:
//   Move temp0 = y
// block2 (L0) [1]:
//   ReturnScalar temp0
//   VarKill temp0
func testOr1(x, y bool) bool {
	return x || y
}

//test:irdump
// block0 [0]:
//   JumpZero L0 x
// block1 [0]:
//   LoadScalarConst temp0.v0 = 10
//   ReturnScalar temp0.v0
// block2 (L0) [0]:
//   LoadScalarConst temp0.v1 = 20
//   ReturnScalar temp0.v1
func testIf1(x bool) int {
	if x {
		return 10
	}
	return 20
}
