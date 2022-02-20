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

func intexpr(x, y int) {
	i := x + y
	println(x + y + x)
	println(y + x + y)
	println(5 + i)
	println(i + i)
	println(i - 5)
	println(5 - i)
	println(i * 3)
	println(x * 3 * y)
	println(x * y)
	println(x / y)
	println(y / x)
	println(i / 2)
	println((i * 3) / 10)
	i++
	println(i)
	i++
	println(i)
}

func main() {
	constexpr1()
	intexpr(14, 5)
	intexpr(-14, 5)
	intexpr(14, -5)
}
