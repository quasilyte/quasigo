package main

import (
	"github.com/quasilyte/quasigo/internal/evaltest"
)

func pick(cond bool) string {
	return evaltest.If(cond, "yes", "no")
}

func main() {
	// Inferred type args.
	println(evaltest.If(true, 1, 2))
	println(evaltest.If(false, 1, 2))
	x := 10
	println(evaltest.If(x > 5, "big", "small"))
	println(pick(true))
	println(pick(false))

	// Explicit type args.
	println(evaltest.If[string](x > 50, "big", "small"))
	println(evaltest.If[int](x < 50, x, 0))

	// byte is an alias, the instantiation is If[[]uint8].
	bs := make([]byte, 0, 2)
	bs = append(bs, 'a')
	bs = append(bs, 'b')
	println(string(evaltest.If(len(bs) == 2, bs, nil)))
	println(len(evaltest.If(len(bs) == 3, bs, nil)))

	// Named types.
	p := evaltest.If(x == 10, evaltest.NewPoint(1, 2), nil)
	println(p.X)
	println(evaltest.If(x != 10, p, nil) == nil)
	println(evaltest.Len(evaltest.NewPoints(3)))

	// Several type args.
	println(evaltest.Describe("a", 1))
	println(evaltest.Describe("b", x+1))

	// Generic call results as generic call args.
	println(evaltest.If(evaltest.If(x > 0, false, true), 1, evaltest.If(x > 0, 3, 4)))
	println(evaltest.Describe(evaltest.If(x > 0, "pos", "neg"), evaltest.If(x > 0, x, -x)))

	// Discarded result.
	evaltest.If(true, 1, 2)
}
