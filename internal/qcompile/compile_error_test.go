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
			src: `return sprintf("%q", sprintf("str"))`,
			err: `can't call testpkg.sprintf: nested variadic calls are not supported`,
		},

		{
			src: `return manyargs(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)`,
			err: `native funcs can't have more than 6 args, got 10`,
		},
	}

	makePackageSource := func(body string) string {
		return `
		  package ` + testPackage + `
		  func add1(x int) int { return x + 1 }
		  func concat(s1, s2 string) string { return s1 + s2 }
		  func f(i int, s string, b bool, err error) interface{} {
			` + body + `
		  }
		  func imul(x, y int) int
		  func idiv2(x, y int) (int, int)
		  func manyargs(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9 int) int
		  func sprintf(format string, args ...interface{}) string
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
