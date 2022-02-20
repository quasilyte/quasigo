package main

const one int = 1

//test:disasm
// main.constexpr1 code=25 frame=48 (2 slots: 0 args, 1 locals, 1 temps)
//   LoadScalarConst x = 1
//   MoveScalar arg0 = x
//   CallVoidNative builtin.PrintInt()
//   LoadScalarConst arg0 = 1
//   CallVoidNative builtin.PrintInt()
//   LoadScalarConst tmp0 = 3
//   MoveScalar arg0 = tmp0
//   CallVoidNative builtin.PrintInt()
//   ReturnVoid
//
//test:disasm_opt
// main.constexpr1 code=22 frame=24 (1 slots: 0 args, 1 locals, 0 temps)
//   LoadScalarConst x = 1
//   MoveScalar arg0 = x
//   CallVoidNative builtin.PrintInt()
//   LoadScalarConst arg0 = 1
//   CallVoidNative builtin.PrintInt()
//   LoadScalarConst arg0 = 3
//   CallVoidNative builtin.PrintInt()
//   ReturnVoid
func constexpr1() {
	x := one
	println(x)
	println(one)
	println(one + 1 + one)
}

func main() {
	constexpr1()
}
