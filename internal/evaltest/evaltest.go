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

func NilEface() any { return nil }

func NilFoo() *Foo { return nil }

func NilFooAsEface() any { return (*Foo)(nil) }

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

type Point struct {
	X int
	Y int
}

func NewPoint(x, y int) *Point {
	return &Point{X: x, Y: y}
}

func NewPoints(n int) []*Point {
	points := make([]*Point, n)
	for i := range points {
		points[i] = &Point{X: i, Y: i * 10}
	}
	return points
}

// SumX returns the sum of the X fields, nil elements are skipped.
func SumX(points []*Point) int {
	sum := 0
	for _, p := range points {
		if p != nil {
			sum += p.X
		}
	}
	return sum
}
