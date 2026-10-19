package main

//test:disasm_opt
// main.recurA code=19 frame=96 (4 slots: 1 params, 3 locals)
//   JumpNotZero L0 x
//   ReturnOne
// L0:
//   LoadScalarConst temp2 = 1
//   IntSub64 temp1 = x temp2
//   Move arg0 = temp1
//   CallRecur temp0
//   ReturnScalar temp0
func recurA(x int) bool {
	if x == 0 {
		return true
	}
	return recurA(x - 1)
}

func main() {
	println(recurA(39))
}
