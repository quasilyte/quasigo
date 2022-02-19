package main

func cond1(x, y int) bool {
	return (x == 0 || x > 0) && (y < 5 || y >= 10)
}

func cond2(x, y int) bool {
	return (x != 0 || x < 0) || y < 5
}

func cond3(x, y int) bool {
	return x == 1 || x == 2 || y == 3 || y < 0
}

func test1(x, y int) {
	println(cond1(x, y))
	println(cond1(y, x))
	println(cond1(x, x))
	println(cond1(y, y))
}

func test2(x, y int) {
	println(cond2(x, y))
	println(cond2(y, x))
	println(cond2(x, x))
	println(cond2(y, y))
}

func main() {
	test1(-1, -1)
	test1(-1, 0)
	test1(1, 0)
	test1(2, 0)
	test1(2, 1)
	test1(-2, 1)
	test1(1031, 102)
	test1(29, 10)
	test1(-29, -10)
	test1(-130, -130)
	test1(0, -130)
	test1(0, 130)
	test1(10, 130)
	test1(10, 10)
	test1(5, 10)
}
