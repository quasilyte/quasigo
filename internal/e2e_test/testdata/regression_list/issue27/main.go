package main

func intneg(x int) int {
	return -x
}

func noinline(x, y int) {
	for i := 0; i < 0; i++ {
	}
}

func testIntNeg() {
	noinline(3, -1)
	println(intneg(0))
}

func main() {
	testIntNeg()
}
