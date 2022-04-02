package main

//test:disasm_both
// main.constexprFold code=5 frame=24 (1 slots: 0 params, 1 locals)
//   LoadScalarConst temp0 = 1138
//   ReturnScalar temp0
func constexprFold() int {
	return 40 + 549*2
}

//test:disasm_both
// main.returnZero code=1 frame=0 (0 slots: 0 params, 0 locals)
//   ReturnZero
func returnZero() int {
	return 0
}

//test:disasm_both
// main.returnZeroByte code=1 frame=0 (0 slots: 0 params, 0 locals)
//   ReturnZero
func returnZeroByte() byte {
	return 0
}

//test:disasm_both
// main.returnOne code=1 frame=0 (0 slots: 0 params, 0 locals)
//   ReturnOne
func returnOne() int {
	return 1
}

//test:disasm_both
// main.returnOneByte code=1 frame=0 (0 slots: 0 params, 0 locals)
//   ReturnOne
func returnOneByte() byte {
	return 1
}

//test:disasm_both
// main.returnFalse code=1 frame=0 (0 slots: 0 params, 0 locals)
//   ReturnZero
func returnFalse() bool {
	return false
}

//test:disasm_both
// main.returnTrue code=1 frame=0 (0 slots: 0 params, 0 locals)
//   ReturnOne
func returnTrue() bool {
	return true
}

func main() {
	println(returnZero())
	println(returnZeroByte())
	println(returnOne())
	println(returnOneByte())
	println(returnFalse())
	println(returnTrue())
}
