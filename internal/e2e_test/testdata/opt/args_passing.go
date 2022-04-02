package opttest

func imul(x, y int) int {
	for {
		break
	}
	return x * y
}

func concat(x, y string) string {
	for {
		break
	}
	return x + y
}

//test:disasm_opt
// opttest.argsPassing1 code=32 frame=72 (3 slots: 0 params, 3 locals)
//   LoadScalarConst arg0 = 1
//   LoadScalarConst arg1 = 2
//   Call temp2 = opttest.imul()
//   Move arg0 = temp2
//   LoadScalarConst arg1 = 3
//   Call temp1 = opttest.imul()
//   Move arg0 = temp1
//   LoadScalarConst arg1 = 4
//   Call temp0 = opttest.imul()
//   ReturnScalar temp0
func argsPassing1() int {
	return imul(imul(imul(1, 2), 3), 4)
}

//test:disasm_opt
// opttest.argsPassing2 code=32 frame=72 (3 slots: 0 params, 3 locals)
//   LoadScalarConst arg0 = 3
//   LoadScalarConst arg1 = 4
//   Call temp2 = opttest.imul()
//   LoadScalarConst arg0 = 2
//   Move arg1 = temp2
//   Call temp1 = opttest.imul()
//   LoadScalarConst arg0 = 1
//   Move arg1 = temp1
//   Call temp0 = opttest.imul()
//   ReturnScalar temp0
func argsPassing2() int {
	x1 := 1
	x2 := 2
	x3 := 3
	x4 := 4
	return imul(x1, imul(x2, imul(x3, x4)))
}

//test:disasm_opt
// opttest.argsPassing3 code=32 frame=72 (3 slots: 0 params, 3 locals)
//   LoadStrConst arg0 = "1"
//   LoadStrConst arg1 = "2"
//   Call temp2 = opttest.concat()
//   Move arg0 = temp2
//   LoadStrConst arg1 = "3"
//   Call temp1 = opttest.concat()
//   Move arg0 = temp1
//   LoadStrConst arg1 = "4"
//   Call temp0 = opttest.concat()
//   ReturnStr temp0
func argsPassing3() string {
	return concat(concat(concat("1", "2"), "3"), "4")
}
