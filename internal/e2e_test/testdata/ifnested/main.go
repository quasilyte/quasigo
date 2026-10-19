package main

func nestedIf(a, b bool) int {
	if a {
		if b {
			return 1
		}
	} else {
		return 2
	}
	return 3
}

func nestedIfElse(a, b bool) int {
	x := 0
	if a {
		if b {
			return 1
		}
	} else if b {
		x = 20
	} else {
		return 2
	}
	return x + 3
}

func nestedSwitch(a bool, n int) int {
	if a {
		switch n {
		case 0:
			return 10
		case 1:
			if n > 5 {
				return 11
			}
		}
	} else {
		return 20
	}
	return 30
}

func nestedSwitchCase(n, m int) int {
	r := 0
	switch n {
	case 0:
		if m == 0 {
			return 1
		}
	case 1:
		r = 5
	default:
		r = 7
	}
	return r
}

func nestedFor(n int) int {
	if n > 0 {
		for i := 0; i < n; i++ {
			if i == 3 {
				return i
			}
		}
	} else {
		return -1
	}
	return 100
}

func voidNested(a, b bool) {
	if a {
		if b {
			println("ab")
			return
		}
	} else {
		println("not a")
		return
	}
	println("end")
}

func main() {
	for i := 0; i < 4; i++ {
		a := i%2 == 1
		b := i > 1
		println(nestedIf(a, b))
		println(nestedIfElse(a, b))
		voidNested(a, b)
	}
	for n := 0; n < 3; n++ {
		println(nestedSwitch(true, n))
		println(nestedSwitch(false, n))
		println(nestedSwitchCase(n, 0))
		println(nestedSwitchCase(n, 1))
	}
	println(nestedFor(2))
	println(nestedFor(5))
	println(nestedFor(0))
}
