package main

//test:irdump
// block0 [0]:
//   LoadScalarConst temp0 = 10
//   LoadScalarConst temp1.v0 = 5
//   LoadScalarConst temp2.v0 = 10
//   IntAdd64 temp0 = temp0 temp2.v0
//   IntAdd64 temp0.v0 = temp0 temp1.v0
//   ReturnScalar temp0.v0
//
//test:disasm_both
// main.blockScopeTest1 code=19 frame=72 (3 slots: 0 params, 3 locals)
//   LoadScalarConst temp0 = 10
//   LoadScalarConst temp1 = 5
//   LoadScalarConst temp2 = 10
//   IntAdd64 temp0 = temp0 temp2
//   IntAdd64 temp0 = temp0 temp1
//   ReturnScalar temp0
func blockScopeTest1() int {
	result := 10
	x := 5
	{
		x := 10
		result += x
	}
	result += x
	return result
}

//test:irdump
// block0 [0]:
//   Zero temp0
//   LoadScalarConst temp1.v0 = 11
//   Zero temp2
//   LoadScalarConst temp2.v0 = 5
//   IntAdd64 temp0 = temp0 temp2.v0
//   IntAdd64 temp0.v0 = temp0 temp1.v0
//   ReturnScalar temp0.v0
//
//test:disasm_both
// main.blockScopeTest2 code=20 frame=72 (3 slots: 0 params, 3 locals)
//   Zero temp0
//   LoadScalarConst temp1 = 11
//   Zero temp2
//   LoadScalarConst temp2 = 5
//   IntAdd64 temp0 = temp0 temp2
//   IntAdd64 temp0 = temp0 temp1
//   ReturnScalar temp0
func blockScopeTest2() int {
	result := 0
	x := 11
	{
		x := 0
		{
			x = 5
		}
		result += x
	}
	result += x
	return result
}

//test:irdump
// block0 [0]:
//   Move temp0 = x
//   LoadScalarConst temp1.v0 = 10
//   IntAdd64 temp0 = temp0 temp1.v0
//   LoadScalarConst x = 5
//   IntAdd64 temp0.v0 = temp0 x
//   ReturnScalar temp0.v0
//
//test:disasm_both
// main.blockScopeTest3 code=19 frame=72 (3 slots: 1 params, 2 locals)
//   Move temp0 = x
//   LoadScalarConst temp1 = 10
//   IntAdd64 temp0 = temp0 temp1
//   LoadScalarConst x = 5
//   IntAdd64 temp0 = temp0 x
//   ReturnScalar temp0
func blockScopeTest3(x int) int {
	result := x
	{
		x := 10
		result += x
	}
	x = 5
	result += x
	return result
}

func testBlockScope() {
	println(blockScopeTest1())
	println(blockScopeTest2())
	println(blockScopeTest3(5))
}

//test:irdump
// block0 [0]:
//   Zero temp0
//   LoadScalarConst temp1 = 101
//   LoadScalarConst temp2 = 5
//   Jump L2
// block1 (L3) [0]:
// block2 (L1) [0]:
//   IntInc temp0
// block3 (L2) [0]:
//   IntInc temp2
// block4 [0]:
//   LoadScalarConst temp4.v0 = 10
//   IntLt temp3.v0 = temp2 temp4.v0
//   JumpNotZero L3 temp3.v0
// block5 (L0) [1]:
//   IntAdd64 temp0.v0 = temp0 temp1
//   ReturnScalar temp0.v0
//   VarKill temp1
//
//test:disasm_both
// main.forScopeTest1 code=32 frame=120 (5 slots: 0 params, 5 locals)
//   Zero temp0
//   LoadScalarConst temp1 = 101
//   LoadScalarConst temp2 = 5
//   Jump L0
// L1:
//   IntInc temp0
//   IntInc temp2
// L0:
//   LoadScalarConst temp4 = 10
//   IntLt temp3 = temp2 temp4
//   JumpNotZero L1 temp3
//   IntAdd64 temp0 = temp0 temp1
//   ReturnScalar temp0
func forScopeTest1() int {
	result := 0
	x := 101
	for x := 5; x < 10; x++ {
		result++
	}
	result += x
	return result
}

//test:irdump
// block0 [0]:
//   Zero temp0
//   LoadScalarConst temp1 = 101
//   LoadScalarConst temp2 = 5
//   Jump L2
// block1 (L3) [0]:
//   LoadScalarConst temp3 = 4
//   Jump L6
// block2 (L7) [0]:
// block3 (L5) [0]:
//   IntAdd64 temp0 = temp0 temp3
// block4 (L6) [0]:
//   IntDec temp3
// block5 [0]:
//   Zero temp5.v0
//   ScalarNotEq temp4.v0 = temp3 temp5.v0
//   JumpNotZero L7 temp4.v0
// block6 (L4) [0]:
// block7 (L1) [0]:
//   IntInc temp0
// block8 (L2) [0]:
//   IntInc temp2
// block9 [0]:
//   LoadScalarConst temp4.v1 = 10
//   IntLt temp3.v0 = temp2 temp4.v1
//   JumpNotZero L3 temp3.v0
// block10 (L0) [1]:
//   IntAdd64 temp0.v0 = temp0 temp1
//   ReturnScalar temp0.v0
//   VarKill temp1
//
//test:disasm_opt
// main.forScopeTest2 code=48 frame=120 (5 slots: 0 params, 5 locals)
//   Zero temp0
//   LoadScalarConst temp1 = 101
//   LoadScalarConst temp2 = 5
//   Jump L0
// L3:
//   LoadScalarConst temp3 = 4
//   Jump L1
// L2:
//   IntAdd64 temp0 = temp0 temp3
//   IntDec temp3
// L1:
//   JumpNotZero L2 temp3
//   IntInc temp0
//   IntInc temp2
// L0:
//   LoadScalarConst temp4 = 10
//   IntLt temp3 = temp2 temp4
//   JumpNotZero L3 temp3
//   IntAdd64 temp0 = temp0 temp1
//   ReturnScalar temp0
func forScopeTest2() int {
	result := 0
	x := 101
	for x := 5; x < 10; x++ {
		for x := 4; x != 0; x-- {
			result += x
		}
		result++
	}
	result += x
	return result
}

func testForScope() {
	println(forScopeTest1())
	println(forScopeTest2())
}

func main() {
	testBlockScope()
	testForScope()
}
