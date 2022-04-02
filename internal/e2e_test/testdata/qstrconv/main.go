package main

import (
	"strconv"
)

//test:disasm_both
// main.testAtoi code=25 frame=96 (4 slots: 1 params, 3 locals)
//   Move arg0 = s
//   CallNative temp0 = strconv.Atoi()
//   MoveResult2 temp1
//   Move arg0 = temp0
//   CallVoidNative builtin.PrintInt()
//   IsNilInterface temp2 = temp1
//   Move arg0 = temp2
//   CallVoidNative builtin.PrintBool()
//   ReturnVoid
func testAtoi(s string) {
	i, err := strconv.Atoi(s)
	println(i)
	println(err == nil)
}

func main() {
	s := "16"
	i := 11

	testAtoi(s)

	i2, err2 := strconv.Atoi("bad")
	println(i2)
	println(err2.Error())

	println(strconv.Itoa(140))
	println(strconv.Itoa(i) == s)

	i, err2 = strconv.Atoi("foo")
	println(i)
	println(err2.Error())

	i, err2 = strconv.Atoi("-349")
	println(i)
	if err2 == nil {
		println("err2 is nil")
	}
}
