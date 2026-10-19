package main

import "fmt"

func wrap(n int) string {
	return fmt.Sprintf("<%d>", n)
}

func wrap2(s string, n int) string {
	return fmt.Sprintf("[%s:%d]", s, n)
}

func main() {
	// Args that call funcs which use the variadic buffer themselves.
	println(fmt.Sprintf("%s %s", wrap(1), wrap(2)))
	println(fmt.Sprintf("%d %s %d", 5, wrap(6), 7))
	println(fmt.Sprintf("%s", wrap2("x", 3)))
	println(fmt.Sprintf(wrap(1)+"%s|%s", wrap2("a", 1), wrap2("b", 2)))

	// Float args should be formatted as floats.
	f := 2.5
	println(fmt.Sprintf("%v", f))
	println(fmt.Sprintf("%.2f|%v|%d", 1.25, f*2, 3))
	b := true
	by := byte(200)
	println(fmt.Sprintf("%v %v %v", b, by, "s"))
}
