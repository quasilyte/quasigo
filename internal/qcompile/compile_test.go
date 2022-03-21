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
			`  LoadScalarConst temp0 = 1138`,
			`  ReturnScalar temp0`,
		},

		`return "ok"`: {
			`testpkg.f code=5 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  LoadStrConst temp0 = "ok"`,
			`  ReturnStr temp0`,
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
			`  LoadScalarConst temp1 = 1`,
			`  IntAdd64 temp0 = x temp1`,
			`  ReturnScalar temp0`,
		},
		`x := 10; return x - 1`: {
			`testpkg.f code=12 frame=168 (7 slots: 4 args, 1 locals, 2 temps)`,
			`  LoadScalarConst x = 10`,
			`  LoadScalarConst temp1 = 1`,
			`  IntSub64 temp0 = x temp1`,
			`  ReturnScalar temp0`,
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
			`  IntAdd64 temp0 = v1 v2`,
			`  ReturnScalar temp0`,
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
			`  LoadScalarConst temp1 = 2`,
			`  ScalarNotEq temp0 = i temp1`,
			`  JumpZero L0 temp0`,
			`  LoadStrConst temp0 = "a"`,
			`  ReturnStr temp0`,
			`L0:`,
			`  JumpZero L1 b`,
			`  LoadStrConst temp0 = "b"`,
			`  ReturnStr temp0`,
			`L1:`,
			`  LoadStrConst temp0 = "c"`,
			`  ReturnStr temp0`,
		},

		`j := -5; for { if j > 0 { break }; j++; }; return j`: {
			`testpkg.f code=23 frame=168 (7 slots: 4 args, 1 locals, 2 temps)`,
			`  LoadScalarConst j = -5`,
			`L2:`,
			`  Zero temp1`,
			`  IntGt temp0 = j temp1`,
			`  JumpZero L0 temp0`,
			`  Jump L1`,
			`L0:`,
			`  IntInc j`,
			`  Jump L2`,
			`L1:`,
			`  ReturnScalar j`,
		},

		`return len(s)`: {
			`testpkg.f code=5 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  Len temp0 = s`,
			`  ReturnScalar temp0`,
		},

		`return len(s) >= 0`: {
			`testpkg.f code=11 frame=168 (7 slots: 4 args, 0 locals, 3 temps)`,
			`  Len temp1 = s`,
			`  Zero temp2`,
			`  IntGtEq temp0 = temp1 temp2`,
			`  ReturnScalar temp0`,
		},

		`return s[:]`: {
			`testpkg.f code=5 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  Move temp0 = s`,
			`  ReturnStr temp0`,
		},

		`return s[:][:][:]`: {
			`testpkg.f code=5 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  Move temp0 = s`,
			`  ReturnStr temp0`,
		},

		`return s[0:]`: {
			`testpkg.f code=5 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  Move temp0 = s`,
			`  ReturnStr temp0`,
		},

		`return s[0:len(s)]`: {
			`testpkg.f code=5 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  Move temp0 = s`,
			`  ReturnStr temp0`,
		},

		`return s[:len(s)]`: {
			`testpkg.f code=5 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  Move temp0 = s`,
			`  ReturnStr temp0`,
		},

		// TODO: optimize.
		`return !(i == 0)`: {
			`testpkg.f code=11 frame=168 (7 slots: 4 args, 0 locals, 3 temps)`,
			`  Zero temp2`,
			`  ScalarEq temp1 = i temp2`,
			`  Not temp0 = temp1`,
			`  ReturnScalar temp0`,
		},

		`return s[1:]`: {
			`testpkg.f code=9 frame=144 (6 slots: 4 args, 0 locals, 2 temps)`,
			`  LoadScalarConst temp1 = 1`,
			`  StrSliceFrom temp0 = s temp1`,
			`  ReturnStr temp0`,
		},

		`return s[:1]`: {
			`testpkg.f code=9 frame=144 (6 slots: 4 args, 0 locals, 2 temps)`,
			`  LoadScalarConst temp1 = 1`,
			`  StrSliceTo temp0 = s temp1`,
			`  ReturnStr temp0`,
		},

		`return s[1:2]`: {
			`testpkg.f code=13 frame=168 (7 slots: 4 args, 0 locals, 3 temps)`,
			`  LoadScalarConst temp1 = 1`,
			`  LoadScalarConst temp2 = 2`,
			`  StrSlice temp0 = s temp1 temp2`,
			`  ReturnStr temp0`,
		},

		// TODO: optimize.
		`return i + 0`: {
			`testpkg.f code=8 frame=144 (6 slots: 4 args, 0 locals, 2 temps)`,
			`  Zero temp1`,
			`  IntAdd64 temp0 = i temp1`,
			`  ReturnScalar temp0`,
		},

		// TODO: emit inc for +1.
		`x := 0; x += 1; return x`: {
			`testpkg.f code=11 frame=144 (6 slots: 4 args, 1 locals, 1 temps)`,
			`  Zero x`,
			`  LoadScalarConst temp0 = 1`,
			`  IntAdd64 x = x temp0`,
			`  ReturnScalar x`,
		},

		// TODO: optimize.
		`if !b { return 10 }; return 20`: {
			`testpkg.f code=17 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  Not temp0 = b`,
			`  JumpZero L0 temp0`,
			`  LoadScalarConst temp0 = 10`,
			`  ReturnScalar temp0`,
			`L0:`,
			`  LoadScalarConst temp0 = 20`,
			`  ReturnScalar temp0`,
		},

		// TODO: optimize.
		`x := i; cond := x == 0; for !cond { cond = x == 0; x-- }; return 10`: {
			`testpkg.f code=32 frame=168 (7 slots: 4 args, 2 locals, 1 temps)`,
			`  Move x = i`,
			`  Zero temp0`,
			`  ScalarEq cond = x temp0`,
			`  Jump L0`,
			`L1:`,
			`  Zero temp0`,
			`  ScalarEq cond = x temp0`,
			`  IntDec x`,
			`L0:`,
			`  Not temp0 = cond`,
			`  JumpNotZero L1 temp0`,
			`  LoadScalarConst temp0 = 10`,
			`  ReturnScalar temp0`,
		},

		`return len("x")`: {
			`testpkg.f code=1 frame=96 (4 slots: 4 args, 0 locals, 0 temps)`,
			`  ReturnOne`,
		},

		`return i == 10 || i == 2`: {
			`testpkg.f code=20 frame=168 (7 slots: 4 args, 0 locals, 3 temps)`,
			`  LoadScalarConst temp1 = 10`,
			`  ScalarEq temp0 = i temp1`,
			`  JumpNotZero L0 temp0`,
			`  LoadScalarConst temp2 = 2`,
			`  ScalarEq temp0 = i temp2`,
			`L0:`,
			`  ReturnScalar temp0`,
		},

		`return i == 10 && s == "foo"`: {
			`testpkg.f code=20 frame=168 (7 slots: 4 args, 0 locals, 3 temps)`,
			`  LoadScalarConst temp1 = 10`,
			`  ScalarEq temp0 = i temp1`,
			`  JumpZero L0 temp0`,
			`  LoadStrConst temp2 = "foo"`,
			`  StrEq temp0 = s temp2`,
			`L0:`,
			`  ReturnScalar temp0`,
		},

		`x := 10; y := 20; return (x == 0 || x > 0) && (y < 5 || y >= 10)`: {
			`testpkg.f code=46 frame=264 (11 slots: 4 args, 2 locals, 5 temps)`,
			`  LoadScalarConst x = 10`,
			`  LoadScalarConst y = 20`,
			`  Zero temp1`,
			`  ScalarEq temp0 = x temp1`,
			`  JumpNotZero L0 temp0`,
			`  Zero temp2`,
			`  IntGt temp0 = x temp2`,
			`L0:`,
			`  JumpZero L1 temp0`,
			`  LoadScalarConst temp3 = 5`,
			`  IntLt temp0 = y temp3`,
			`  JumpNotZero L1 temp0`,
			`  LoadScalarConst temp4 = 10`,
			`  IntGtEq temp0 = y temp4`,
			`L1:`,
			`  ReturnScalar temp0`,
		},

		`return imul(i, 5) == 10`: {
			`testpkg.f code=19 frame=168 (7 slots: 4 args, 0 locals, 3 temps)`,
			`  Move arg0 = i`,
			`  LoadScalarConst arg1 = 5`,
			`  CallNative temp1 = testpkg.imul()`,
			`  LoadScalarConst temp2 = 10`,
			`  ScalarEq temp0 = temp1 temp2`,
			`  ReturnScalar temp0`,
		},

		`x, y := idiv2(30, 4); return x + y`: {
			`testpkg.f code=18 frame=168 (7 slots: 4 args, 2 locals, 1 temps)`,
			`  LoadScalarConst arg0 = 30`,
			`  LoadScalarConst arg1 = 4`,
			`  CallNative x = testpkg.idiv2()`,
			`  MoveResult2 y`,
			`  IntAdd64 temp0 = x y`,
			`  ReturnScalar temp0`,
		},

		`return add1(10)`: {
			`testpkg.f code=9 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  LoadScalarConst arg0 = 10`,
			`  Call temp0 = testpkg.add1()`,
			`  ReturnScalar temp0`,
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
			`  IntAdd64 temp1 = x x`,
			`  IntAdd64 temp0 = temp1 x`,
			`  ReturnScalar temp0`,
		},

		`x := 1; return x + x`: {
			`testpkg.f code=9 frame=144 (6 slots: 4 args, 1 locals, 1 temps)`,
			`  LoadScalarConst x = 1`,
			`  IntAdd64 temp0 = x x`,
			`  ReturnScalar temp0`,
		},

		`return err.Error()`: {
			`testpkg.f code=9 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  Move arg0 = err`,
			`  CallNative temp0 = error.Error()`,
			`  ReturnStr temp0`,
		},

		`formatString := "%s:%d"; filename := "file.go"; line := 14; return sprintf(formatString, filename, line)`: {
			`testpkg.f code=29 frame=216 (9 slots: 4 args, 3 locals, 2 temps)`,
			`  LoadStrConst formatString = "%s:%d"`,
			`  LoadStrConst filename = "file.go"`,
			`  LoadScalarConst line = 14`,
			`  Move arg0 = formatString`,
			`  VariadicReset`,
			`  Move temp1 = filename`,
			`  PushVariadicStrArg temp1`,
			`  Move temp1 = line`,
			`  PushVariadicScalarArg temp1`,
			`  CallNative temp0 = testpkg.sprintf()`,
			`  ReturnStr temp0`,
		},

		`return sprintf("%s:%d", "file", 10)`: {
			`testpkg.f code=20 frame=144 (6 slots: 4 args, 0 locals, 2 temps)`,
			`  LoadStrConst arg0 = "%s:%d"`,
			`  VariadicReset`,
			`  LoadStrConst temp1 = "file"`,
			`  PushVariadicStrArg temp1`,
			`  LoadScalarConst temp1 = 10`,
			`  PushVariadicScalarArg temp1`,
			`  CallNative temp0 = testpkg.sprintf()`,
			`  ReturnStr temp0`,
		},

		`return err == nil`: {
			`testpkg.f code=5 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  IsNilInterface temp0 = err`,
			`  ReturnScalar temp0`,
		},

		`return nil == err`: {
			`testpkg.f code=5 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  IsNilInterface temp0 = err`,
			`  ReturnScalar temp0`,
		},

		`return err != nil`: {
			`testpkg.f code=5 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  IsNotNilInterface temp0 = err`,
			`  ReturnScalar temp0`,
		},

		`return nil != err`: {
			`testpkg.f code=5 frame=120 (5 slots: 4 args, 0 locals, 1 temps)`,
			`  IsNotNilInterface temp0 = err`,
			`  ReturnScalar temp0`,
		},

		`return imul(imul(1, 2), imul(3, 4))`: {
			`testpkg.f code=32 frame=168 (7 slots: 4 args, 0 locals, 3 temps)`,
			`  LoadScalarConst arg0 = 1`,
			`  LoadScalarConst arg1 = 2`,
			`  CallNative temp1 = testpkg.imul()`,
			`  LoadScalarConst arg0 = 3`,
			`  LoadScalarConst arg1 = 4`,
			`  CallNative temp2 = testpkg.imul()`,
			`  Move arg0 = temp1`,
			`  Move arg1 = temp2`,
			`  CallNative temp0 = testpkg.imul()`,
			`  ReturnScalar temp0`,
		},

		`return imul(imul(imul(1, 2), 3), 4)`: {
			`testpkg.f code=38 frame=216 (9 slots: 4 args, 0 locals, 5 temps)`,
			`  LoadScalarConst arg0 = 1`,
			`  LoadScalarConst arg1 = 2`,
			`  CallNative temp2 = testpkg.imul()`,
			`  LoadScalarConst temp3 = 3`,
			`  Move arg0 = temp2`,
			`  Move arg1 = temp3`,
			`  CallNative temp1 = testpkg.imul()`,
			`  LoadScalarConst temp4 = 4`,
			`  Move arg0 = temp1`,
			`  Move arg1 = temp4`,
			`  CallNative temp0 = testpkg.imul()`,
			`  ReturnScalar temp0`,
		},

		`x1 := 1; x2 := 2; x3 := 3; x4 := 4; return imul(imul(imul(x1, x2), x3), x4)`: {
			`testpkg.f code=50 frame=312 (13 slots: 4 args, 4 locals, 5 temps)`,
			`  LoadScalarConst x1 = 1`,
			`  LoadScalarConst x2 = 2`,
			`  LoadScalarConst x3 = 3`,
			`  LoadScalarConst x4 = 4`,
			`  Move arg0 = x1`,
			`  Move arg1 = x2`,
			`  CallNative temp2 = testpkg.imul()`,
			`  Move temp3 = x3`,
			`  Move arg0 = temp2`,
			`  Move arg1 = temp3`,
			`  CallNative temp1 = testpkg.imul()`,
			`  Move temp4 = x4`,
			`  Move arg0 = temp1`,
			`  Move arg1 = temp4`,
			`  CallNative temp0 = testpkg.imul()`,
			`  ReturnScalar temp0`,
		},

		`x1 := 1; x2 := 2; x3 := 3; x4 := 4; return imul(x1, imul(x2, imul(x3, x4)))`: {
			`testpkg.f code=50 frame=312 (13 slots: 4 args, 4 locals, 5 temps)`,
			`  LoadScalarConst x1 = 1`,
			`  LoadScalarConst x2 = 2`,
			`  LoadScalarConst x3 = 3`,
			`  LoadScalarConst x4 = 4`,
			`  Move temp1 = x1`,
			`  Move temp3 = x2`,
			`  Move arg0 = x3`,
			`  Move arg1 = x4`,
			`  CallNative temp4 = testpkg.imul()`,
			`  Move arg0 = temp3`,
			`  Move arg1 = temp4`,
			`  CallNative temp2 = testpkg.imul()`,
			`  Move arg0 = temp1`,
			`  Move arg1 = temp2`,
			`  CallNative temp0 = testpkg.imul()`,
			`  ReturnScalar temp0`,
		},

		`return concat("x", sprintf("%d", 10))`: {
			`testpkg.f code=28 frame=192 (8 slots: 4 args, 0 locals, 4 temps)`,
			`  LoadStrConst temp1 = "x"`,
			`  LoadStrConst arg0 = "%d"`,
			`  VariadicReset`,
			`  LoadScalarConst temp3 = 10`,
			`  PushVariadicScalarArg temp3`,
			`  CallNative temp2 = testpkg.sprintf()`,
			`  Move arg0 = temp1`,
			`  Move arg1 = temp2`,
			`  Call temp0 = testpkg.concat()`,
			`  ReturnStr temp0`,
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
