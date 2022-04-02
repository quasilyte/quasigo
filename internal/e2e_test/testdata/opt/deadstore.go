package opttest

//test:irdump
// block0 [0]:
//   Move temp0.v0 = i
//   Move temp1.v0 = temp0.v0
//   Move temp2.v0 = temp1.v0
//   Move temp3.v0 = temp2.v0
//   Move temp4.v0 = temp3.v0
//   ReturnScalar temp4.v0
//
//test:disasm_opt
// opttest.deadstore1 code=5 frame=48 (2 slots: 1 params, 1 locals)
//   Move temp0 = i
//   ReturnScalar temp0
func deadstore1(i int) int {
	x1 := i
	x2 := x1
	x3 := x2
	x4 := x3
	x5 := x4
	return x5
}

//test:disasm_opt
// opttest.deadstore2 code=5 frame=48 (2 slots: 1 params, 1 locals)
//   Len temp0 = s
//   ReturnScalar temp0
func deadstore2(s string) int {
	length := len(s)
	x1 := length
	x2 := x1
	return x2
}
