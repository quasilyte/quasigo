package opttest

//test:irdump
// block0 [0]:
//   LoadScalarConst temp0.v0 = 10
//   Move temp1.v0 = temp0.v0
//   ReturnScalar temp1.v0
//
//test:disasm_opt
// opttest.constPropagationInt1 code=5 frame=24 (1 slots: 0 params, 1 locals)
//   LoadScalarConst temp0 = 10
//   ReturnScalar temp0
func constPropagationInt1() int {
	x := 10
	y := x
	return y
}

//test:irdump
// block0 [0]:
//   LoadScalarConst temp0.v0 = 20
//   Move temp1.v0 = temp0.v0
//   Move temp2.v0 = temp1.v0
//   ReturnScalar temp2.v0
//
//test:disasm_opt
// opttest.constPropagationInt2 code=5 frame=24 (1 slots: 0 params, 1 locals)
//   LoadScalarConst temp0 = 20
//   ReturnScalar temp0
func constPropagationInt2() int {
	x := 20
	y := x
	z := y
	return z
}

//test:irdump
// block0 [0]:
//   LoadScalarConst temp0.v0 = 30
//   Move temp1.v0 = temp0.v0
//   Move temp2.v0 = temp1.v0
//   Move temp3.v0 = temp2.v0
//   ReturnScalar temp3.v0
//
//test:disasm_opt
// opttest.constPropagationInt3 code=5 frame=24 (1 slots: 0 params, 1 locals)
//   LoadScalarConst temp0 = 30
//   ReturnScalar temp0
func constPropagationInt3() int {
	x := 30
	y := x
	z := y
	result := z
	return result
}

//test:irdump
// block0 [0]:
//   LoadScalarConst temp0.v0 = 40
//   LoadScalarConst temp1.v0 = 30
//   Move temp2.v0 = temp0.v0
//   Move temp3.v0 = temp1.v0
//   IntAdd64 temp4.v0 = temp2.v0 temp3.v0
//   ReturnScalar temp4.v0
//
//test:disasm_opt
// opttest.constPropagationInt4 code=5 frame=24 (1 slots: 0 params, 1 locals)
//   LoadScalarConst temp0 = 70
//   ReturnScalar temp0
//
//test:constants_opt
//   scalar constants: [70]
//   string constants: []
//
//test:constants
//   scalar constants: [40 30]
//   string constants: []
func constPropagationInt4() int {
	c1 := 40
	c2 := 30
	x := c1
	y := c2
	return x + y
}

//test:irdump
// block0 [0]:
//   LoadScalarConst temp0.v0 = 1
//   LoadScalarConst temp1.v0 = 2
//   Move temp2.v0 = temp0.v0
//   Move temp3.v0 = temp1.v0
//   IntAdd64 temp4.v0 = temp2.v0 temp3.v0
//   ReturnScalar temp4.v0
//
//test:disasm_opt
// opttest.constPropagationInt5 code=5 frame=24 (1 slots: 0 params, 1 locals)
//   LoadScalarConst temp0 = 3
//   ReturnScalar temp0
//
//test:constants_opt
//   scalar constants: [3]
//   string constants: []
//
//test:constants
//   scalar constants: [1 2]
//   string constants: []
func constPropagationInt5() int {
	c1 := 1
	c2 := len("go")
	x := c1
	y := c2
	result := x + y
	return result
}

//test:irdump
// block0 [0]:
//   LoadScalarConst temp0.v0 = 1
//   LoadScalarConst temp2.v0 = 1
//   IntAdd64 temp1.v0 = temp0.v0 temp2.v0
//   LoadScalarConst temp3.v0 = 1
//   IntAdd64 temp2.v1 = temp1.v0 temp3.v0
//   LoadScalarConst temp4.v0 = 1
//   IntAdd64 temp3.v1 = temp2.v1 temp4.v0
//   ReturnScalar temp3.v1
//
//test:disasm_opt
// opttest.constPropagationInt6 code=5 frame=24 (1 slots: 0 params, 1 locals)
//   LoadScalarConst temp0 = 4
//   ReturnScalar temp0
//
//test:constants_opt
//   scalar constants: [4]
//   string constants: []
//
//test:constants
//   scalar constants: [1]
//   string constants: []
func constPropagationInt6() int {
	x1 := 1
	x2 := x1 + 1
	x3 := x2 + 1
	return x3 + 1
}

//test:irdump
// block0 [0]:
//   Zero temp0.v0
//   Move temp1.v0 = temp0.v0
//   ReturnScalar temp1.v0
//
//test:disasm_opt
// opttest.constPropagationInt7 code=4 frame=24 (1 slots: 0 params, 1 locals)
//   Zero temp0
//   ReturnScalar temp0
//
//test:constants_opt
//   scalar constants: []
//   string constants: []
//
//test:constants
//   scalar constants: []
//   string constants: []
func constPropagationInt7() int {
	x1 := 0
	x2 := x1
	return x2
}

//test:irdump
// block0 [0]:
//   Zero temp0.v0
//   Zero temp1.v0
//   IntAdd64 temp2.v0 = temp0.v0 temp1.v0
//   Zero temp4.v0
//   IntAdd64 temp3.v0 = temp2.v0 temp4.v0
//   ReturnScalar temp3.v0
//
//test:disasm_opt
// opttest.constPropagationInt8 code=4 frame=24 (1 slots: 0 params, 1 locals)
//   Zero temp0
//   ReturnScalar temp0
func constPropagationInt8() int {
	x1 := 0
	x2 := 0
	x3 := x1 + x2
	return x3 + 0
}

//test:irdump
// block0 [0]:
//   LoadStrConst temp0.v0 = "str"
//   Move temp1.v0 = temp0.v0
//   Move temp2.v0 = temp1.v0
//   ReturnStr temp2.v0
//
//test:disasm_opt
// opttest.constPropagationStr1 code=5 frame=24 (1 slots: 0 params, 1 locals)
//   LoadStrConst temp0 = "str"
//   ReturnStr temp0
func constPropagationStr1() string {
	s1 := "str"
	s2 := s1
	s3 := s2
	return s3
}
