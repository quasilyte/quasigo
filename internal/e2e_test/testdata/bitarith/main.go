package main

func intMasks(x, y int) {
	println(x & y)
	println(x &^ y)
	println(x | y)
	println(x ^ y)
	z := x
	z &= y
	println(z)
	z = x
	z &^= y
	println(z)
}

func byteExpr(x, y byte) {
	println(x & y)
	println(x &^ y)
	z := x
	z &= 15
	println(z)
	z = x
	z &^= 15
	println(z)
}

func hasFlag(flags, flag int) bool {
	return flags&flag != 0
}

func main() {
	intMasks(12, 10)
	intMasks(-1, 255)
	intMasks(0, 7)
	byteExpr(255, 15)
	byteExpr(200, 100)

	println(hasFlag(5, 4))
	println(hasFlag(5, 2))
}
