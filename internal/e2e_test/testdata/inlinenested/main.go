package main

func inc(x int) int {
	return x + 1
}

func add(x, y int) int {
	return x + y
}

func sub(x, y int) int {
	return x - y
}

func add3(x, y, z int) int {
	return x + y + z
}

func withTemp(x, y int) int {
	t := x * 2
	return t + y
}

func concat(a, b string) string {
	return a + b
}

func isPos(x int) bool {
	return x > 0
}

//test:disasm_opt
// main.nestedInlineDisasm code=25 frame=192 (8 slots: 0 params, 8 locals)
//   LoadScalarConst temp1 = 10
//   LoadScalarConst temp4 = 1
//   LoadScalarConst temp7 = 1
//   IntAdd64 temp6 = temp4 temp7
//   Move temp2 = temp6
//   IntAdd64 temp5 = temp1 temp2
//   Move temp0 = temp5
//   ReturnScalar temp0
func nestedInlineDisasm() int {
	return add(10, inc(1))
}

func main() {
	// Inlined calls as arguments of other inlined calls
	// should not share the temps slots.

	println(add(10, inc(1)))

	println(sub(inc(10), inc(1)))
	println(sub(inc(1), inc(10)))
	println(add(inc(1), add(2, inc(3))))
	println(add3(1, inc(1), add(inc(2), 3)))
	println(withTemp(3, withTemp(4, 5)))
	println(withTemp(withTemp(1, 2), withTemp(3, 4)))
	println(add(withTemp(1, 2), inc(inc(inc(0)))))
	println(concat("a", concat("b", concat("c", "d"))))
	println(concat(concat("a", "b"), concat("c", "d")))

	x := 5
	y := add(x, inc(x))
	z := sub(add(y, 1), inc(x))
	println(x)
	println(y)
	println(z)

	if isPos(sub(inc(1), 5)) {
		println("pos")
	} else {
		println("neg")
	}
	if isPos(add(inc(1), 5)) && isPos(sub(inc(0), inc(1))) {
		println("both")
	} else {
		println("not both")
	}
}
