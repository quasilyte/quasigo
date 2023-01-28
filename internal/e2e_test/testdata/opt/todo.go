package opttest

// TODO: compile `s == ""` as `len(s) == 0`
//
//test:disasm_opt
// opttest.todoEmptyStringCompare code=17 frame=72 (3 slots: 1 params, 2 locals)
//   LoadStrConst temp1 = ""
//   StrEq temp0 = s temp1
//   JumpZero L0 temp0
//   ReturnOne
// L0:
//   LoadScalarConst temp0 = 2
//   ReturnScalar temp0
func todoEmptyStringCompare(s string) int {
	if s == "" {
		return 1
	}
	return 2
}

// TODO: optimize this to `ReturnScalar b`
//
//test:disasm_opt
// opttest.todoSimpleIfReturnBool code=6 frame=24 (1 slots: 1 params, 0 locals)
//   JumpZero L0 b
//   ReturnOne
// L0:
//   ReturnZero
func todoSimpleIfReturnBool(b bool) bool {
	if b {
		return true
	}
	return false
}

// TODO: x+0 -> x
//
//test:disasm_opt
// opttest.todoArith code=8 frame=72 (3 slots: 1 params, 2 locals)
//   Zero temp1
//   IntAdd64 temp0 = i temp1
//   ReturnScalar temp0
func todoArith(i int) int {
	return i + 0
}

// TODO: x+=1 -> x++
//
//test:disasm_opt
// opttest.todoInc code=9 frame=48 (2 slots: 1 params, 1 locals)
//   LoadScalarConst temp0 = 1
//   IntAdd64 i = i temp0
//   ReturnScalar i
func todoInc(i int) int {
	i += 1
	return i
}

// TODO: fuse into <= 0.
//
//test:disasm_opt
// opttest.todoFuseComparisons code=18 frame=72 (3 slots: 1 params, 2 locals)
//   Zero temp1
//   ScalarEq temp0 = i temp1
//   JumpNotZero L0 temp0
//   Zero temp1
//   IntLt temp0 = i temp1
// L0:
//   ReturnScalar temp0
func todoFuseComparisons(i int) bool {
	return i == 0 || i < 0
}

//test:disasm_opt
// opttest.todoInverseEq code=11 frame=96 (4 slots: 1 params, 3 locals)
//   Zero temp2
//   ScalarEq temp1 = i temp2
//   Not temp0 = temp1
//   ReturnScalar temp0
func todoInverseEq(i int) bool {
	return !(i == 0)
}
