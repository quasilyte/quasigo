package main

import (
	"fmt"

	"github.com/quasilyte/quasigo/internal/evaltest"
)

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

	// Nested variadic calls.
	println(evaltest.Digits(1, 2, evaltest.Digits(3, 4, 5)))
	println(evaltest.Digits(evaltest.Digits(1, 2), 3, evaltest.Digits(4, evaltest.Digits(5, 6))))
	println(fmt.Sprintf("%s|%s", fmt.Sprintf("%d", 1), fmt.Sprintf("%d-%s", 2, fmt.Sprintf("%v", 3))))
	println(fmt.Sprintf("%q", fmt.Sprintf("str")))
	println(fmt.Sprintf(fmt.Sprintf("%%s:%s", "a"), "b"))
	println(fmt.Sprintf("%s %d", wrap(evaltest.Digits(1, 2)), evaltest.Digits(3, len(fmt.Sprintf("%d", 45)))))
	for i := 0; i < 3; i++ {
		println(evaltest.Digits(i, evaltest.Digits(i+1, i+2), i+3))
	}

	// A call without variadic args should not observe the previous call args.
	println(evaltest.Digits(1, 2, evaltest.Digits()))
	println(evaltest.Digits(7, 8))
	println(evaltest.Digits())
}
