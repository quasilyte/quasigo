package qcompile_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/quasilyte/quasigo"
	"github.com/quasilyte/quasigo/internal/testutil"
	"github.com/quasilyte/quasigo/qnative"
)

const testPackage = `testpkg`

func TestCompileError(t *testing.T) {
	tests := []struct {
		src string
		err string
	}{
		{
			src: `return manyargs(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)`,
			err: `native funcs can't have more than 6 args, got 10`,
		},

		{
			src: `return If(b, 1, 2)`,
		},
		{
			src: `return If[int](b, 1, 2)`,
		},
		{
			src: `return If(b, make([]byte, 0), nil)`,
		},
		{
			src: `return If[int32](b, 1, 2)`,
		},
		{
			src: `return If[rune](b, 1, 2)`,
		},
		{
			src: `return If[any](b, 1, 2)`,
		},
		{
			src: `return If[myInt](b, 1, 2)`,
		},
		{
			src: `return If(b, myInt(1), 2)`,
		},
		{
			src: `return If(b, myInts(nil), nil)`,
			err: `no native binding for testpkg.If[[]int] instantiation`,
		},
		{
			src: `return If[byte](b, 1, 2)`,
			err: `no native binding for testpkg.If[uint8] instantiation`,
		},
		{
			src: `return If(b, "a", "b")`,
			err: `no native binding for testpkg.If[string] instantiation`,
		},
		{
			src: `return If(b, map[int]int{}, nil)`,
			err: `can't call If: map[int]int type arg is not supported`,
		},
	}

	makePackageSource := func(body string) string {
		return `
		  package ` + testPackage + `
		  func add1(x int) int { return x + 1 }
		  func concat(s1, s2 string) string { return s1 + s2 }
		  func f(i int, s string, b bool, err error) any {
			` + body + `
		  }
		  func imul(x, y int) int
		  func idiv2(x, y int) (int, int)
		  func manyargs(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9 int) int
		  func sprintf(format string, args ...any) string
		  func If[T any](cond bool, onTrue, onFalse T) T { var zero T; return zero }
		  type myInt = int
		  type myInts = []myInt
		  `
	}

	for i := range tests {
		test := tests[i]
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			env := quasigo.NewEnv()
			env.AddNativeFunc(testPackage, "manyargs", func(ctx qnative.CallContext) {
				panic("should not be called")
			})
			env.AddNativeFunc(testPackage, "imul", func(ctx qnative.CallContext) {
				panic("should not be called")
			})
			env.AddNativeFunc(testPackage, "idiv2", func(ctx qnative.CallContext) {
				panic("should not be called")
			})
			env.AddNativeFunc(testPackage, "sprintf", func(ctx qnative.CallContext) {
				panic("should not be called")
			})
			env.AddNativeFunc(testPackage, "If[int]", func(ctx qnative.CallContext) {
				panic("should not be called")
			})
			env.AddNativeFunc(testPackage, "If[[]uint8]", func(ctx qnative.CallContext) {
				panic("should not be called")
			})
			env.AddNativeFunc(testPackage, "If[rune]", func(ctx qnative.CallContext) {
				panic("should not be called")
			})
			env.AddNativeFunc(testPackage, "If[interface{}]", func(ctx qnative.CallContext) {
				panic("should not be called")
			})
			env.AddNativeFunc("builtin", "PrintInt", func(ctx qnative.CallContext) {
				panic("should not be called")
			})
			env.AddNativeFunc("builtin", "PrintString", func(ctx qnative.CallContext) {
				panic("should not be called")
			})
			env.AddNativeMethod(`error`, `Error`, func(ctx qnative.CallContext) {
				panic("should not be called")
			})
			src := makePackageSource(test.src)
			parsed, err := testutil.ParseGoFile(testPackage, src)
			if err != nil {
				t.Fatalf("parse %s: %v", test.src, err)
			}
			_, err = testutil.CompileTestFile(env, "f", testPackage, parsed)
			want := "<empty error>"
			if test.err != "" {
				want = test.err
			}
			have := "<empty error>"
			if err != nil {
				have = err.Error()
			}
			if !strings.Contains(have, want) {
				t.Fatalf("errors mismatched:\nhave: %q\nwant: %q", have, want)
			}
		})

	}
}

func TestCompileErrorMultiResult(t *testing.T) {
	tests := []struct {
		src string
		err string
	}{
		{
			src: `func f() (int, int, int) { return 1, 2, 3 }`,
			err: `at most 2 value results are supported, f has 3`,
		},

		{
			src: `func f() (int, map[int]int) { return 1, nil }`,
			err: `function result type: map[int]int is not supported`,
		},

		{
			src: `func f() (a, b int) { return }`,
			err: `'naked' return statements are not allowed`,
		},

		{
			src: `
				func two() (int, int) { return 1, 2 }
				func take(x, y int) int { return x + y }
				func f() int { return take(two()) }`,
			err: `can't pass tuple as a func argument`,
		},

		{
			src: `func f() int { a, b := 1, 2; return a + b }`,
			err: `only single right operand is allowed in assignments`,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			env := quasigo.NewEnv()
			src := "package " + testPackage + "\n" + test.src
			parsed, err := testutil.ParseGoFile(testPackage, src)
			if err != nil {
				t.Fatalf("parse %s: %v", test.src, err)
			}
			_, err = testutil.CompileTestFile(env, "f", testPackage, parsed)
			if err == nil {
				t.Fatalf("expected an error containing %q, got none", test.err)
			}
			if !strings.Contains(err.Error(), test.err) {
				t.Fatalf("errors mismatched:\nhave: %q\nwant: %q", err.Error(), test.err)
			}
		})
	}
}

func TestCompileErrorGenericFunc(t *testing.T) {
	src := "package " + testPackage + "\n" + `func id[T any](x T) T { return x }`
	parsed, err := testutil.ParseGoFile(testPackage, src)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	_, err = testutil.CompileTestFunc(quasigo.NewEnv(), "id", parsed)
	want := `id: generic funcs can only be bound as natives`
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Fatalf("errors mismatched:\nhave: %v\nwant: %q", err, want)
	}
}
