package main

import (
	"fmt"

	"github.com/quasilyte/quasigo/internal/evaltest"
)

func newFoo(s string) *evaltest.Foo {
	return evaltest.NewFoo(s)
}

func testFooPair1(p *evaltest.FooPair) {
	println(p.First().String())
	println(p.Second().String())
	p.SetFirst(p.Second())
	println(p.First().String())
	println(p.Second().String())

	p2 := evaltest.NewFooPair(newFoo("a"), newFoo("b"))
	println(p2.First().String())
	println(p2.Second().String())
	println(p2.Get("first").String())
	println(p2.Get("second").String())
	p2.SetFirst(p2.Get("first"))
	println(p2.First().String())
	println(p2.Second().String())
	println(p2.Get("first").String())
	println(p2.Get("second").String())
	p2.SetFirst(p2.Get("second"))
	println(p2.First().String())
	println(p2.Second().String())
	println(p2.Get("first").String())
	println(p2.Get("second").String())
}

func testFooPair2() {
	p := evaltest.NewFooPair(newFoo("1"), newFoo("2"))
	println(p.First().String())
	p.SetFirstPrefix(p.Get("first").String())
	println(p.First().String())
	p.SetFirstPrefix(p.Get("second").String())
	println(p.First().String())
}

func testFooPair3() {
	p := evaltest.NewFooPair(newFoo("1"), newFoo("2"))
	println(p.First().String())
	p.SetFirstPrefix(fmt.Sprintf("%v", p.Get("first").String()))
	println(p.First().String())
	p.SetFirstPrefix(fmt.Sprintf("%v", p.Get("second").String()))
	println(p.First().String())
}

func testFooPair1_2(p *evaltest.FooPair) {
	println(p.First().Prefix)
	println(p.Second().Prefix)
	p.SetFirst(p.Second())
	println(p.First().Prefix)
	println(p.Second().Prefix)

	p2 := evaltest.NewFooPair(newFoo("a"), newFoo("b"))
	println(p2.First().Prefix)
	println(p2.Second().Prefix)
	println(p2.Get("first").Prefix)
	println(p2.Get("second").Prefix)
	p2.SetFirst(p2.Get("first"))
	println(p2.First().Prefix)
	println(p2.Second().Prefix)
	println(p2.Get("first").Prefix)
	println(p2.Get("second").Prefix)
	p2.SetFirst(p2.Get("second"))
	println(p2.First().Prefix)
	println(p2.Second().Prefix)
	println(p2.Get("first").Prefix)
	println(p2.Get("second").Prefix)
}

func testFooPair2_2() {
	p := evaltest.NewFooPair(newFoo("1"), newFoo("2"))
	println(p.First().Prefix)
	p.SetFirstPrefix(p.Get("first").Prefix)
	println(p.First().Prefix)
	p.SetFirstPrefix(p.Get("second").Prefix)
	println(p.First().Prefix)
}

func testFooPair3_2() {
	p := evaltest.NewFooPair(newFoo("1"), newFoo("2"))
	println(p.First().Prefix)
	p.SetFirstPrefix(fmt.Sprintf("%v", p.Get("first").Prefix))
	println(p.First().Prefix)
	p.SetFirstPrefix(fmt.Sprintf("%v", p.Get("second").Prefix))
	println(p.First().Prefix)
}

func main() {
	foo := evaltest.NewFoo("foo1")
	println(foo.Method1(10))
	println(evaltest.NewFoo("foo2").Method1(20))

	foo2 := newFoo("example")
	println(foo2.Method1(1032))
	println(newFoo("exampletwo").Method1(0))

	pair := evaltest.NewFooPair(newFoo("_1"), newFoo("_2"))

	testFooPair1(evaltest.NewFooPair(foo, foo2))
	testFooPair2()
	testFooPair3()
	testFooPair1_2(pair)
	testFooPair2_2()
	testFooPair3_2()
}
