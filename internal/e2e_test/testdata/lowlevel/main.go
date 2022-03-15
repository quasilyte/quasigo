package main

//test:disasm_both
// main.returnZero code=1 frame=0 (0 slots: 0 args, 0 locals, 0 temps)
//   ReturnZero
func returnZero() int {
	return 0
}

//test:disasm_both
// main.returnZeroByte code=1 frame=0 (0 slots: 0 args, 0 locals, 0 temps)
//   ReturnZero
func returnZeroByte() byte {
	return 0
}

//test:disasm_both
// main.returnOne code=1 frame=0 (0 slots: 0 args, 0 locals, 0 temps)
//   ReturnOne
func returnOne() int {
	return 1
}

//test:disasm_both
// main.returnOneByte code=1 frame=0 (0 slots: 0 args, 0 locals, 0 temps)
//   ReturnOne
func returnOneByte() byte {
	return 1
}

//test:disasm_both
// main.returnFalse code=1 frame=0 (0 slots: 0 args, 0 locals, 0 temps)
//   ReturnZero
func returnFalse() bool {
	return false
}

//test:disasm_both
// main.returnTrue code=1 frame=0 (0 slots: 0 args, 0 locals, 0 temps)
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
