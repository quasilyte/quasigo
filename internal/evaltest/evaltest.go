package evaltest

import (
	"fmt"
)

// This package is used for quasigo testing.

type Foo struct {
	Prefix string
}

func (foo *Foo) String() string { return foo.Prefix }

func (foo *Foo) Method1(x int) string { return foo.Prefix + fmt.Sprint(x) }

func NewFoo(prefix string) *Foo { return &Foo{Prefix: prefix} }

func NilEface() interface{} { return nil }

func NilFoo() *Foo { return nil }

func NilFooAsEface() interface{} { return (*Foo)(nil) }

type FooPair struct {
	first  *Foo
	second *Foo
}

func NewFooPair(x, y *Foo) *FooPair {
	return &FooPair{first: x, second: y}
}

func (p *FooPair) SetFirst(x *Foo) { p.first = x }

func (p *FooPair) SetFirstPrefix(s string) { p.first.Prefix = s }

func (p *FooPair) Get(key string) *Foo {
	if key == "first" {
		return p.first
	}
	if key == "second" {
		return p.second
	}
	return nil
}

func (p *FooPair) First() *Foo  { return p.first }
func (p *FooPair) Second() *Foo { return p.second }
