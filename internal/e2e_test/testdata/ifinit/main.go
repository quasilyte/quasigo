package main

import "strconv"

func double(x int) int {
	return x * 2
}

func divmod(x, y int) (int, int) {
	return x / y, x % y
}

func find(xs []int, want int) (int, bool) {
	for i := 0; i < len(xs); i++ {
		if xs[i] == want {
			return i, true
		}
	}
	return -1, false
}

func classify(x int) string {
	if y := double(x); y > 100 {
		return "huge"
	} else if z := y + x; z > 20 {
		return "big " + strconv.Itoa(y) + " " + strconv.Itoa(z)
	} else if z > 6 {
		return "medium " + strconv.Itoa(y) + " " + strconv.Itoa(z)
	} else {
		return "small " + strconv.Itoa(x) + " " + strconv.Itoa(y) + " " + strconv.Itoa(z)
	}
}

func parse(s string) {
	if n, err := strconv.Atoi(s); err == nil {
		println(n)
	} else {
		println(err.Error())
	}
}

func isEven(x int) bool {
	if r := x % 2; r == 0 {
		return true
	}
	return false
}

func main() {
	{
		xs := make([]int, 3)
		xs[0] = 4
		xs[1] = 5
		xs[2] = 6
		if i, ok := find(xs, 5); ok {
			println(i)
		}
		if i, ok := find(xs, 50); ok {
			println(i)
		} else {
			println(-100)
			println(i)
		}
		if _, ok := find(xs, 6); ok {
			println("found 6")
		}
	}

	{
		println(classify(1))
		println(classify(3))
		println(classify(9))
		println(classify(51))
	}

	parse("123")
	parse("bad")

	{
		for i := 0; i < 5; i++ {
			if q, r := divmod(17, i+1); r == 0 {
				println(q)
			} else if r > 2 {
				println(r)
			}
		}
	}

	{
		// Assignment to the existing variable.
		x := 0
		if x = double(21); x > 40 {
			println(x)
		}
		println(x)
	}

	{
		// The init variable shadows the outer one and is dead after the if.
		x := 10
		if x := x + 1; x > 10 {
			println(x)
			if x := x * 2; x > 20 {
				println(x)
			}
			println(x)
		}
		println(x)
		if x := 1; x > 5 {
			println("unreachable")
		}
		y := "after"
		println(y)
		println(x)
	}

	{
		// Init variable is visible in all branches.
		for i := 0; i < 4; i++ {
			if v := i * 3; v == 0 {
				println("zero")
			} else if v == 3 {
				println(v + 100)
			} else {
				println(v)
			}
		}
	}

	{
		i := 0
		if i++; i > 0 {
			println(i)
		}
		if println("init"); i == 1 {
			println("done")
		}
	}

	{
		sum := 0
		for i := 0; i < 10; i++ {
			if sq := i * i; sq%2 == 0 {
				sum += sq
			}
		}
		println(sum)
		println(isEven(4))
		println(isEven(7))
	}
}
