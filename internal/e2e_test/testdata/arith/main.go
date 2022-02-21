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

//test:disasm_both
// main.intnegvar code=5 frame=48 (2 slots: 1 args, 0 locals, 1 temps)
//   IntNeg tmp0 = x
//   ReturnScalar tmp0
func intnegvar(x int) int {
	return -x
}

//test:disasm_both
// main.intnegconst code=5 frame=24 (1 slots: 0 args, 0 locals, 1 temps)
//   LoadScalarConst tmp0 = -50
//   ReturnScalar tmp0
func intnegconst() int {
	return -50
}

func byteexpr(x, y byte) {
	i := x + y
	println(i)
	i--
	println(i)
	println(x + y + x)
	println(y + x + y)
	println(5 + i)
	println(i + i)
	println(i - 5)
	println(5 - i)
	println(i * 3)
	println(x * 3 * y)
	println(x * y)
	i++
	println(i)
	i++
	println(i)
	println(-i)

	b := byte(210)
	println(b)
	println(b + b)
	println(b + b*x)
	println(b*b + x)
	println(b > x)
	println(b >= x)
	println(b < x)
	println(b <= x)
	println(b > y)
	println(b >= y)
	println(b < y)
	println(b <= y)
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
	println(-i)

	i2 := 329
	println(i2)
	println(i2 + i2)
	println(i2 + i2*x)
	println(i2*i2 + x)
	println(i2 > x)
	println(i2 >= x)
	println(i2 < x)
	println(i2 <= x)
	println(i2 > y)
	println(i2 >= y)
	println(i2 < y)
	println(i2 <= y)
}

func boolexpr(x, y bool) {
	println(x == y)
	println(x != y)
}

func main() {
	constexpr1()
	intexpr(14, 5)
	intexpr(-14, 5)
	intexpr(14, -5)
	byteexpr(10, 30)
	byteexpr(0, 0)
	byteexpr(255, 0)
	byteexpr(0, 255)
	boolexpr(false, false)
	boolexpr(false, true)
	boolexpr(true, false)
	boolexpr(true, true)
}
