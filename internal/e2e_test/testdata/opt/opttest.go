package opttest

//test:irdump
// block0 [0]:
//   Len temp2.v0 = s
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
//test:disasm_opt
// opttest.zerocmp1 code=17 frame=96 (4 slots: 1 args, 0 locals, 3 temps)
//   Len temp2 = s
//   JumpZero L0 temp2
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
