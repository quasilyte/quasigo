package main

const one int = 1

//test:disasm
// main.constexpr1 code=25 frame=48 (2 slots: 0 params, 2 locals)
//   LoadScalarConst temp0 = 1
//   Move arg0 = temp0
//   CallVoidNative builtin.PrintInt()
//   LoadScalarConst arg0 = 1
//   CallVoidNative builtin.PrintInt()
//   LoadScalarConst temp1 = 3
//   Move arg0 = temp1
//   CallVoidNative builtin.PrintInt()
//   ReturnVoid
//
//test:disasm_opt
// main.constexpr1 code=19 frame=0 (0 slots: 0 params, 0 locals)
//   LoadScalarConst arg0 = 1
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
// main.intnegvar code=5 frame=48 (2 slots: 1 params, 1 locals)
//   IntNeg temp0 = x
//   ReturnScalar temp0
func intnegvar(x int) int {
	return -x
}

//test:disasm_both
// main.intnegconst code=5 frame=24 (1 slots: 0 params, 1 locals)
//   LoadScalarConst temp0 = -50
//   ReturnScalar temp0
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

	minusTenInt := -10
	minusTen := byte(minusTenInt)
	println(minusTen > 10)
	println(minusTen == 10)

	b2 := byte(190)
	println(int(b + b2))
	println(int(b ^ b2))
	println(int(b - b2))
	println(int(b2 - b))
	println(int(b2 * b))
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

	i3 := 43848234
	println(i3 - i2)
	println(i3 - i2 - i3)
	println(i2 - i3 - i3)
	println(i2 + i3 - i3)
	println(i2 + i3 + i3)
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
