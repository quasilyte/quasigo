package qcompile_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/quasilyte/quasigo"
	"github.com/quasilyte/quasigo/internal/testutil"
	"github.com/quasilyte/quasigo/qnative"
)

func TestCompile(t *testing.T) {
	tests := map[string][]string{
		// We perform const-folding for simple expressions,
		// so there should be no actual evaluations here.
		`return 40 + 549 * 2`: {
			`testpkg.f code=5 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  LoadScalarConst tmp0 = 1138`,
			`  ReturnScalar tmp0`,
		},

		`return "ok"`: {
			`testpkg.f code=5 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  LoadStrConst tmp0 = "ok"`,
			`  ReturnStr tmp0`,
		},

		// No redundant copy (move) is generated in this example.
		// Return over var loads directly from the slot.
		`x := 10; return x`: {
			`testpkg.f code=5 frame=120 (5 slots: 4 args, 1 locals, 0 temps)`,
			`  LoadScalarConst x = 10`,
			`  ReturnScalar x`,
		},

		// A lot of redundant operations here.
		// We probably need a separate optimization pass.
		`x := 10; return x + 1`: {
			`testpkg.f code=12 frame=168 (7 slots: 4 args, 1 locals, 2 temps)`,
			`  LoadScalarConst x = 10`,
			`  LoadScalarConst tmp1 = 1`,
			`  IntAdd64 tmp0 = x tmp1`,
			`  ReturnScalar tmp0`,
		},
		`x := 10; return x - 1`: {
			`testpkg.f code=12 frame=168 (7 slots: 4 args, 1 locals, 2 temps)`,
			`  LoadScalarConst x = 10`,
			`  LoadScalarConst tmp1 = 1`,
			`  IntSub64 tmp0 = x tmp1`,
			`  ReturnScalar tmp0`,
		},

		`x := true; y := !x; return y`: {
			`testpkg.f code=8 frame=144 (6 slots: 4 args, 2 locals, 0 temps)`,
			`  LoadScalarConst x = 1`,
			`  Not y = x`,
			`  ReturnScalar y`,
		},

		`x := 1; y := x; return y`: {
			`testpkg.f code=8 frame=144 (6 slots: 4 args, 2 locals, 0 temps)`,
			`  LoadScalarConst x = 1`,
			`  Move y = x`,
			`  ReturnScalar y`,
		},

		`x := 0; x++; return x`: {
			`testpkg.f code=6 frame=120 (5 slots: 4 args, 1 locals, 0 temps)`,
			`  Zero x`,
			`  IntInc x`,
			`  ReturnScalar x`,
		},
		`x := 0; x--; return x`: {
			`testpkg.f code=6 frame=120 (5 slots: 4 args, 1 locals, 0 temps)`,
			`  Zero x`,
			`  IntDec x`,
			`  ReturnScalar x`,
		},

		`x := 1; y := 2; v1 := x + y; v2 := v1 + v1; return v1 + v2`: {
			`testpkg.f code=20 frame=216 (9 slots: 4 args, 4 locals, 1 temps)`,
			`  LoadScalarConst x = 1`,
			`  LoadScalarConst y = 2`,
			`  IntAdd64 v1 = x y`,
			`  IntAdd64 v2 = v1 v1`,
			`  IntAdd64 tmp0 = v1 v2`,
			`  ReturnScalar tmp0`,
		},

		`if b { return 1 }; return 0`: {
			`testpkg.f code=6 frame=96 (4 slots: 4 args, 0 locals, 0 temps)`,
			`  JumpZero L0 b`,
			`  ReturnOne`,
			`L0:`,
			`  ReturnZero`,
		},

		`if b { return 1 } else { return 0 }`: {
			`testpkg.f code=6 frame=96 (4 slots: 4 args, 0 locals, 0 temps)`,
			`  JumpZero L0 b`,
			`  ReturnOne`,
			`L0:`,
			`  ReturnZero`,
		},

		`x := 0; if b { x = 5 } else { x = 50 }; return x`: {
			`testpkg.f code=17 frame=120 (5 slots: 4 args, 1 locals, 0 temps)`,
			`  Zero x`,
			`  JumpZero L0 b`,
			`  LoadScalarConst x = 5`,
			`  Jump L1`,
			`L0:`,
			`  LoadScalarConst x = 50`,
			`L1:`,
			`  ReturnScalar x`,
		},

		`if i != 2 { return "a" } else if b { return "b" }; return "c"`: {
			`testpkg.f code=30 frame=144 (6 slots: 4 args, 0 locals, 2 temps)`,
			`  LoadScalarConst tmp1 = 2`,
			`  ScalarNotEq tmp0 = i tmp1`,
			`  JumpZero L0 tmp0`,
			`  LoadStrConst tmp0 = "a"`,
			`  ReturnStr tmp0`,
			`L0:`,
			`  JumpZero L1 b`,
			`  LoadStrConst tmp0 = "b"`,
			`  ReturnStr tmp0`,
			`L1:`,
			`  LoadStrConst tmp0 = "c"`,
			`  ReturnStr tmp0`,
		},

		`j := -5; for { if j > 0 { break }; j++; }; return j`: {
			`testpkg.f code=23 frame=168 (7 slots: 4 args, 1 locals, 2 temps)`,
			`  LoadScalarConst j = -5`,
			`L2:`,
			`  Zero tmp1`,
			`  IntGt tmp0 = j tmp1`,
			`  JumpZero L0 tmp0`,
			`  Jump L1`,
			`L0:`,
			`  IntInc j`,
			`  Jump L2`,
			`L1:`,
			`  ReturnScalar j`,
		},

		`return len(s)`: {
			`testpkg.f code=5 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  Len tmp0 = s`,
			`  ReturnScalar tmp0`,
		},

		`return len(s) >= 0`: {
			`testpkg.f code=11 frame=168 (7 slots: 4 args, 0 locals, 3 temps)`,
			`  Len tmp1 = s`,
			`  Zero tmp2`,
			`  IntGtEq tmp0 = tmp1 tmp2`,
			`  ReturnScalar tmp0`,
		},

		`return s[:]`: {
			`testpkg.f code=5 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  Move tmp0 = s`,
			`  ReturnStr tmp0`,
		},

		`return s[:][:][:]`: {
			`testpkg.f code=5 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  Move tmp0 = s`,
			`  ReturnStr tmp0`,
		},

		`return s[0:]`: {
			`testpkg.f code=5 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  Move tmp0 = s`,
			`  ReturnStr tmp0`,
		},

		`return s[0:len(s)]`: {
			`testpkg.f code=5 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  Move tmp0 = s`,
			`  ReturnStr tmp0`,
		},

		`return s[:len(s)]`: {
			`testpkg.f code=5 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  Move tmp0 = s`,
			`  ReturnStr tmp0`,
		},

		// TODO: optimize.
		`return !(i == 0)`: {
			`testpkg.f code=11 frame=168 (7 slots: 4 args, 0 locals, 3 temps)`,
			`  Zero tmp2`,
			`  ScalarEq tmp1 = i tmp2`,
			`  Not tmp0 = tmp1`,
			`  ReturnScalar tmp0`,
		},

		`return s[1:]`: {
			`testpkg.f code=9 frame=144 (6 slots: 4 args, 0 locals, 2 temps)`,
			`  LoadScalarConst tmp1 = 1`,
			`  StrSliceFrom tmp0 = s tmp1`,
			`  ReturnStr tmp0`,
		},

		`return s[:1]`: {
			`testpkg.f code=9 frame=144 (6 slots: 4 args, 0 locals, 2 temps)`,
			`  LoadScalarConst tmp1 = 1`,
			`  StrSliceTo tmp0 = s tmp1`,
			`  ReturnStr tmp0`,
		},

		`return s[1:2]`: {
			`testpkg.f code=13 frame=168 (7 slots: 4 args, 0 locals, 3 temps)`,
			`  LoadScalarConst tmp1 = 1`,
			`  LoadScalarConst tmp2 = 2`,
			`  StrSlice tmp0 = s tmp1 tmp2`,
			`  ReturnStr tmp0`,
		},

		// TODO: optimize.
		`return i + 0`: {
			`testpkg.f code=8 frame=144 (6 slots: 4 args, 0 locals, 2 temps)`,
			`  Zero tmp1`,
			`  IntAdd64 tmp0 = i tmp1`,
			`  ReturnScalar tmp0`,
		},

		// TODO: emit inc for +1.
		`x := 0; x += 1; return x`: {
			`testpkg.f code=11 frame=144 (6 slots: 4 args, 1 locals, 1 temps)`,
			`  Zero x`,
			`  LoadScalarConst tmp0 = 1`,
			`  IntAdd64 x = x tmp0`,
			`  ReturnScalar x`,
		},

		// TODO: optimize.
		`if !b { return 10 }; return 20`: {
			`testpkg.f code=17 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  Not tmp0 = b`,
			`  JumpZero L0 tmp0`,
			`  LoadScalarConst tmp0 = 10`,
			`  ReturnScalar tmp0`,
			`L0:`,
			`  LoadScalarConst tmp0 = 20`,
			`  ReturnScalar tmp0`,
		},

		// TODO: optimize.
		`x := i; cond := x == 0; for !cond { cond = x == 0; x-- }; return 10`: {
			`testpkg.f code=32 frame=168 (7 slots: 4 args, 2 locals, 1 temps)`,
			`  Move x = i`,
			`  Zero tmp0`,
			`  ScalarEq cond = x tmp0`,
			`  Jump L0`,
			`L1:`,
			`  Zero tmp0`,
			`  ScalarEq cond = x tmp0`,
			`  IntDec x`,
			`L0:`,
			`  Not tmp0 = cond`,
			`  JumpNotZero L1 tmp0`,
			`  LoadScalarConst tmp0 = 10`,
			`  ReturnScalar tmp0`,
		},

		`return len("x")`: {
			`testpkg.f code=1 frame=96 (4 slots: 4 args, 0 locals, 0 temps)`,
			`  ReturnOne`,
		},

		`return i == 10 || i == 2`: {
			`testpkg.f code=20 frame=168 (7 slots: 4 args, 0 locals, 3 temps)`,
			`  LoadScalarConst tmp1 = 10`,
			`  ScalarEq tmp0 = i tmp1`,
			`  JumpNotZero L0 tmp0`,
			`  LoadScalarConst tmp2 = 2`,
			`  ScalarEq tmp0 = i tmp2`,
			`L0:`,
			`  ReturnScalar tmp0`,
		},

		`return i == 10 && s == "foo"`: {
			`testpkg.f code=20 frame=168 (7 slots: 4 args, 0 locals, 3 temps)`,
			`  LoadScalarConst tmp1 = 10`,
			`  ScalarEq tmp0 = i tmp1`,
			`  JumpZero L0 tmp0`,
			`  LoadStrConst tmp2 = "foo"`,
			`  StrEq tmp0 = s tmp2`,
			`L0:`,
			`  ReturnScalar tmp0`,
		},

		`x := 10; y := 20; return (x == 0 || x > 0) && (y < 5 || y >= 10)`: {
			`testpkg.f code=46 frame=264 (11 slots: 4 args, 2 locals, 5 temps)`,
			`  LoadScalarConst x = 10`,
			`  LoadScalarConst y = 20`,
			`  Zero tmp1`,
			`  ScalarEq tmp0 = x tmp1`,
			`  JumpNotZero L0 tmp0`,
			`  Zero tmp2`,
			`  IntGt tmp0 = x tmp2`,
			`L0:`,
			`  JumpZero L1 tmp0`,
			`  LoadScalarConst tmp3 = 5`,
			`  IntLt tmp0 = y tmp3`,
			`  JumpNotZero L1 tmp0`,
			`  LoadScalarConst tmp4 = 10`,
			`  IntGtEq tmp0 = y tmp4`,
			`L1:`,
			`  ReturnScalar tmp0`,
		},

		`return imul(i, 5) == 10`: {
			`testpkg.f code=19 frame=168 (7 slots: 4 args, 0 locals, 3 temps)`,
			`  Move arg0 = i`,
			`  LoadScalarConst arg1 = 5`,
			`  CallNative tmp1 = testpkg.imul()`,
			`  LoadScalarConst tmp2 = 10`,
			`  ScalarEq tmp0 = tmp1 tmp2`,
			`  ReturnScalar tmp0`,
		},

		`x, y := idiv2(30, 4); return x + y`: {
			`testpkg.f code=18 frame=168 (7 slots: 4 args, 2 locals, 1 temps)`,
			`  LoadScalarConst arg0 = 30`,
			`  LoadScalarConst arg1 = 4`,
			`  CallNative x = testpkg.idiv2()`,
			`  MoveResult2 y`,
			`  IntAdd64 tmp0 = x y`,
			`  ReturnScalar tmp0`,
		},

		`return add1(10)`: {
			`testpkg.f code=9 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  LoadScalarConst arg0 = 10`,
			`  Call tmp0 = testpkg.add1()`,
			`  ReturnScalar tmp0`,
		},

		`println("s"); return 0`: {
			`testpkg.f code=7 frame=96 (4 slots: 4 args, 0 locals, 0 temps)`,
			`  LoadStrConst arg0 = "s"`,
			`  CallVoidNative builtin.PrintString()`,
			`  ReturnZero`,
		},

		`println(540); return 0`: {
			`testpkg.f code=7 frame=96 (4 slots: 4 args, 0 locals, 0 temps)`,
			`  LoadScalarConst arg0 = 540`,
			`  CallVoidNative builtin.PrintInt()`,
			`  ReturnZero`,
		},

		`x := 1; return x + x + x`: {
			`testpkg.f code=13 frame=168 (7 slots: 4 args, 1 locals, 2 temps)`,
			`  LoadScalarConst x = 1`,
			`  IntAdd64 tmp1 = x x`,
			`  IntAdd64 tmp0 = tmp1 x`,
			`  ReturnScalar tmp0`,
		},

		`x := 1; return x + x`: {
			`testpkg.f code=9 frame=144 (6 slots: 4 args, 1 locals, 1 temps)`,
			`  LoadScalarConst x = 1`,
			`  IntAdd64 tmp0 = x x`,
			`  ReturnScalar tmp0`,
		},

		`return err.Error()`: {
			`testpkg.f code=9 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  Move arg0 = err`,
			`  CallNative tmp0 = error.Error()`,
			`  ReturnStr tmp0`,
		},

		`formatString := "%s:%d"; filename := "file.go"; line := 14; return sprintf(formatString, filename, line)`: {
			`testpkg.f code=29 frame=216 (9 slots: 4 args, 3 locals, 2 temps)`,
			`  LoadStrConst formatString = "%s:%d"`,
			`  LoadStrConst filename = "file.go"`,
			`  LoadScalarConst line = 14`,
			`  Move arg0 = formatString`,
			`  VariadicReset`,
			`  Move tmp1 = filename`,
			`  PushVariadicStrArg tmp1`,
			`  Move tmp1 = line`,
			`  PushVariadicScalarArg tmp1`,
			`  CallNative tmp0 = testpkg.sprintf()`,
			`  ReturnStr tmp0`,
		},

		`return sprintf("%s:%d", "file", 10)`: {
			`testpkg.f code=20 frame=144 (6 slots: 4 args, 0 locals, 2 temps)`,
			`  LoadStrConst arg0 = "%s:%d"`,
			`  VariadicReset`,
			`  LoadStrConst tmp1 = "file"`,
			`  PushVariadicStrArg tmp1`,
			`  LoadScalarConst tmp1 = 10`,
			`  PushVariadicScalarArg tmp1`,
			`  CallNative tmp0 = testpkg.sprintf()`,
			`  ReturnStr tmp0`,
		},

		`return err == nil`: {
			`testpkg.f code=5 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  IsNilInterface tmp0 = err`,
			`  ReturnScalar tmp0`,
		},

		`return nil == err`: {
			`testpkg.f code=5 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  IsNilInterface tmp0 = err`,
			`  ReturnScalar tmp0`,
		},

		`return err != nil`: {
			`testpkg.f code=5 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  IsNotNilInterface tmp0 = err`,
			`  ReturnScalar tmp0`,
		},

		`return nil != err`: {
			`testpkg.f code=5 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  IsNotNilInterface tmp0 = err`,
			`  ReturnScalar tmp0`,
		},

		`return imul(imul(1, 2), imul(3, 4))`: {
			`testpkg.f code=32 frame=168 (7 slots: 4 args, 0 locals, 3 temps)`,
			`  LoadScalarConst arg0 = 1`,
			`  LoadScalarConst arg1 = 2`,
			`  CallNative tmp1 = testpkg.imul()`,
			`  LoadScalarConst arg0 = 3`,
			`  LoadScalarConst arg1 = 4`,
			`  CallNative tmp2 = testpkg.imul()`,
			`  Move arg0 = tmp1`,
			`  Move arg1 = tmp2`,
			`  CallNative tmp0 = testpkg.imul()`,
			`  ReturnScalar tmp0`,
		},

		`return imul(imul(imul(1, 2), 3), 4)`: {
			`testpkg.f code=38 frame=216 (9 slots: 4 args, 0 locals, 5 temps)`,
			`  LoadScalarConst arg0 = 1`,
			`  LoadScalarConst arg1 = 2`,
			`  CallNative tmp2 = testpkg.imul()`,
			`  LoadScalarConst tmp3 = 3`,
			`  Move arg0 = tmp2`,
			`  Move arg1 = tmp3`,
			`  CallNative tmp1 = testpkg.imul()`,
			`  LoadScalarConst tmp4 = 4`,
			`  Move arg0 = tmp1`,
			`  Move arg1 = tmp4`,
			`  CallNative tmp0 = testpkg.imul()`,
			`  ReturnScalar tmp0`,
		},

		`x1 := 1; x2 := 2; x3 := 3; x4 := 4; return imul(imul(imul(x1, x2), x3), x4)`: {
			`testpkg.f code=50 frame=312 (13 slots: 4 args, 4 locals, 5 temps)`,
			`  LoadScalarConst x1 = 1`,
			`  LoadScalarConst x2 = 2`,
			`  LoadScalarConst x3 = 3`,
			`  LoadScalarConst x4 = 4`,
			`  Move arg0 = x1`,
			`  Move arg1 = x2`,
			`  CallNative tmp2 = testpkg.imul()`,
			`  Move tmp3 = x3`,
			`  Move arg0 = tmp2`,
			`  Move arg1 = tmp3`,
			`  CallNative tmp1 = testpkg.imul()`,
			`  Move tmp4 = x4`,
			`  Move arg0 = tmp1`,
			`  Move arg1 = tmp4`,
			`  CallNative tmp0 = testpkg.imul()`,
			`  ReturnScalar tmp0`,
		},

		`x1 := 1; x2 := 2; x3 := 3; x4 := 4; return imul(x1, imul(x2, imul(x3, x4)))`: {
			`testpkg.f code=50 frame=312 (13 slots: 4 args, 4 locals, 5 temps)`,
			`  LoadScalarConst x1 = 1`,
			`  LoadScalarConst x2 = 2`,
			`  LoadScalarConst x3 = 3`,
			`  LoadScalarConst x4 = 4`,
			`  Move tmp1 = x1`,
			`  Move tmp3 = x2`,
			`  Move arg0 = x3`,
			`  Move arg1 = x4`,
			`  CallNative tmp4 = testpkg.imul()`,
			`  Move arg0 = tmp3`,
			`  Move arg1 = tmp4`,
			`  CallNative tmp2 = testpkg.imul()`,
			`  Move arg0 = tmp1`,
			`  Move arg1 = tmp2`,
			`  CallNative tmp0 = testpkg.imul()`,
			`  ReturnScalar tmp0`,
		},

		`return concat("x", sprintf("%d", 10))`: {
			`testpkg.f code=28 frame=192 (8 slots: 4 args, 0 locals, 4 temps)`,
			`  LoadStrConst tmp1 = "x"`,
			`  LoadStrConst arg0 = "%d"`,
			`  VariadicReset`,
			`  LoadScalarConst tmp3 = 10`,
			`  PushVariadicScalarArg tmp3`,
			`  CallNative tmp2 = testpkg.sprintf()`,
			`  Move arg0 = tmp1`,
			`  Move arg1 = tmp2`,
			`  Call tmp0 = testpkg.concat()`,
			`  ReturnStr tmp0`,
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
		  func sprintf(format string, args ...interface{}) string
		  `
	}

	for testSrc, disasmLines := range tests {
		env := quasigo.NewEnv()
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
		src := makePackageSource(testSrc)
		parsed, err := testutil.ParseGoFile(testPackage, src)
		if err != nil {
			t.Fatalf("parse %s: %v", testSrc, err)
		}
		compiled, err := testutil.CompileTestFile(env, "f", testPackage, parsed)
		if err != nil {
			t.Fatal(err)
		}
		if compiled.IsNil() {
			t.Fatal("can't find f function")
		}
		want := disasmLines
		have := strings.Split(quasigo.Disasm(env, compiled), "\n")

		have = have[:len(have)-1] // Drop an empty line
		if diff := cmp.Diff(have, want); diff != "" {
			t.Errorf("compile %s (-have +want):\n%s", testSrc, diff)
			fmt.Println("For copy/paste:")
			for _, l := range have {
				fmt.Printf("  `%s`,\n", l)
			}
			continue
		}
	}
}
