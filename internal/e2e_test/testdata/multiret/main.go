package main

import "strconv"

func ints(x, y int) (int, int) {
	return x, y
}

func swap(x, y int) (int, int) {
	return y, x
}

func consts() (int, int) {
	return 10, 20
}

func zeros() (int, int) {
	return 0, 0
}

func divmod(x, y int) (int, int) {
	return x / y, x % y
}

func add1(x int) int {
	return x + 1
}

// Results of different types are all returned as whole slots.
func strInt(s string, n int) (string, int) {
	return s + "!", n * 2
}

func boolFloat(b bool, f float64) (bool, float64) {
	return !b, f * 2.5
}

func floats(x float64) (float64, float64) {
	return x, 1
}

func bytesLen() ([]byte, int) {
	b := make([]byte, 3)
	b[0] = 7
	b[2] = 9
	return b, len(b)
}

func strBool(b bool) (string, bool) {
	if b {
		return "yes", true
	}
	return "no", false
}

func branch(x int) (int, string) {
	if x < 0 {
		return -1, "neg"
	}
	if x == 0 {
		return 0, "zero"
	}
	return 1, "pos"
}

func indexOf(xs []int, want int) (int, bool) {
	for i := 0; i < len(xs); i++ {
		if xs[i] == want {
			return i, true
		}
	}
	return -1, false
}

// Operands are evaluated left to right and only then returned.
func traceA() int {
	println("A")
	return 1
}

func traceB() int {
	println("B")
	return 2
}

func order() (int, int) {
	return traceA(), traceB()
}

func callsInReturn(x int) (int, int) {
	return add1(x), add1(add1(x))
}

// `return f()` where f() returns 2 values.
func forward(x, y int) (int, int) {
	return swap(x, y)
}

func forwardParen(x, y int) (int, int) {
	return (ints(x, y))
}

func fibPair(n int) (int, int) {
	if n == 0 {
		return 0, 1
	}
	a, b := fibPair(n - 1)
	return b, a + b
}

// The recursive call result is discarded, but n should be preserved.
func recDiscard(n int) (int, int) {
	if n == 0 {
		return 100, 200
	}
	_, b := recDiscard(n - 1)
	println(n)
	return n, b + 1
}

func main() {
	{
		a, b := ints(1, 2)
		println(a)
		println(b)
	}

	{
		a, b := swap(1, 2)
		println(a)
		println(b)
	}

	{
		a, b := consts()
		println(a)
		println(b)
		a, b = zeros()
		println(a)
		println(b)
	}

	{
		// Assignments to already defined variables.
		a, b := swap(1, 2)
		a, b = swap(a, b)
		println(a)
		println(b)
		b, a = swap(a, b)
		println(a)
		println(b)
		a, b = ints(a*10, b*10)
		println(a)
		println(b)
	}

	{
		// One new variable and one that is already defined.
		a, b := ints(1, 2)
		a, c := ints(3, 4)
		println(a)
		println(b)
		println(c)
	}

	{
		// Shadowing in the inner scope.
		a, b := ints(1, 2)
		{
			a, b := ints(10, 20)
			println(a)
			println(b)
		}
		println(a)
		println(b)
	}

	{
		q, r := divmod(17, 5)
		println(q)
		println(r)
		println(add1(q) + r)
	}

	{
		s, n := strInt("hi", 21)
		println(s)
		println(n)
	}

	{
		b, f := boolFloat(true, 2)
		println(b)
		println(f)
		b, f = boolFloat(b, f)
		println(b)
		println(f)
	}

	{
		f1, f2 := floats(2.5)
		println(f1)
		println(f2)
	}

	{
		bs, n := bytesLen()
		println(bs[0])
		println(bs[2])
		println(n)
	}

	{
		s1, b1 := strBool(true)
		s2, b2 := strBool(false)
		println(s1)
		println(b1)
		println(s2)
		println(b2)
	}

	{
		n1, s1 := branch(-10)
		n2, s2 := branch(0)
		n3, s3 := branch(10)
		println(n1)
		println(s1)
		println(n2)
		println(s2)
		println(n3)
		println(s3)
	}

	{
		xs := make([]int, 4)
		xs[0] = 5
		xs[1] = 6
		xs[2] = 7
		xs[3] = 8
		i, ok := indexOf(xs, 7)
		println(i)
		println(ok)
		i, ok = indexOf(xs, 100)
		println(i)
		println(ok)
	}

	{
		// Results are used inside a loop.
		sum := 0
		for i := 1; i <= 4; i++ {
			q, r := divmod(100, i)
			sum += q*10 + r
		}
		println(sum)
	}

	{
		a, b := order()
		println(a)
		println(b)
	}

	{
		a, b := callsInReturn(10)
		println(a)
		println(b)
	}

	{
		a, b := forward(1, 2)
		println(a)
		println(b)
		a, b = forwardParen(3, 4)
		println(a)
		println(b)
	}

	{
		a, b := fibPair(10)
		println(a)
		println(b)
	}

	{
		a, b := recDiscard(3)
		println(a)
		println(b)
	}

	{
		// Blank identifiers discard the result, but the call is still made.
		_, b := divmod(17, 5)
		println(b)
		a, _ := divmod(17, 5)
		println(a)
		_, _ = order()
		_, c := swap(1, 2)
		println(c)
	}

	{
		// Native funcs with 2 results.
		n, err := strconv.Atoi("12")
		println(n)
		println(err == nil)
		_, err2 := strconv.Atoi("zz")
		println(err2 != nil)
		n2, _ := strconv.Atoi("34")
		println(n2)
		_, _ = strconv.Atoi("56")
	}
}
