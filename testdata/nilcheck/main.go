package main

import (
	"github.com/quasilyte/quasigo/internal/evaltest"
)

func main() {
	println(evaltest.NilEface() == nil)
	println(evaltest.NilFoo() == nil)
	println(evaltest.NilFooAsEface() == nil)
}
