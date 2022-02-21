package qopt_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/quasilyte/quasigo"
	"github.com/quasilyte/quasigo/internal/testutil"
	"github.com/quasilyte/quasigo/qnative"
)

const testPackage = `testpkg`

func TestOptimize(t *testing.T) {
	tests := map[string][]string{
		// TODO: use only 1 temp here.
		`x := 10; return x + 1`: {
			`testpkg.f code=12 frame=168 (7 slots: 4 args, 1 locals, 2 temps)`,
			`  LoadScalarConst x = 10`,
			`  LoadScalarConst tmp1 = 1`,
			`  IntAdd tmp0 = x tmp1`,
			`  ReturnScalar tmp0`,
		},
		`x := 10; return x - 1`: {
			`testpkg.f code=12 frame=168 (7 slots: 4 args, 1 locals, 2 temps)`,
			`  LoadScalarConst x = 10`,
			`  LoadScalarConst tmp1 = 1`,
			`  IntSub tmp0 = x tmp1`,
			`  ReturnScalar tmp0`,
		},

		// TODO: remove redundant local->local stores.
		// Maybe also delete redundant local slots.
		`x := 1; y := x; return y`: {
			`testpkg.f code=8 frame=144 (6 slots: 4 args, 2 locals, 0 temps)`,
			`  LoadScalarConst x = 1`,
			`  MoveScalar y = x`,
			`  ReturnScalar y`,
		},

		// TODO: remove redundant local->tmp->local stores.
		`x1 := 1; x2 := 2; x3 := 3; x4 := 4; return imul(imul(imul(x1, x2), x3), x4)`: {
			`testpkg.f code=50 frame=312 (13 slots: 4 args, 4 locals, 5 temps)`,
			`  LoadScalarConst x1 = 1`,
			`  LoadScalarConst x2 = 2`,
			`  LoadScalarConst x3 = 3`,
			`  LoadScalarConst x4 = 4`,
			`  MoveScalar arg0 = x1`,
			`  MoveScalar arg1 = x2`,
			`  CallNative tmp2 = testpkg.imul()`,
			`  MoveScalar tmp3 = x3`,
			`  MoveScalar arg0 = tmp2`,
			`  MoveScalar arg1 = tmp3`,
			`  CallNative tmp1 = testpkg.imul()`,
			`  MoveScalar tmp4 = x4`,
			`  MoveScalar arg0 = tmp1`,
			`  MoveScalar arg1 = tmp4`,
			`  CallNative tmp0 = testpkg.imul()`,
			`  ReturnScalar tmp0`,
		},

		`return concat(concat(concat("1", "2"), "3"), "4")`: {
			`testpkg.f code=32 frame=168 (7 slots: 4 args, 0 locals, 3 temps)`,
			`  LoadStrConst arg0 = "1"`,
			`  LoadStrConst arg1 = "2"`,
			`  CallNative tmp2 = testpkg.concat()`,
			`  MoveStr arg0 = tmp2`,
			`  LoadStrConst arg1 = "3"`,
			`  CallNative tmp1 = testpkg.concat()`,
			`  MoveStr arg0 = tmp1`,
			`  LoadStrConst arg1 = "4"`,
			`  CallNative tmp0 = testpkg.concat()`,
			`  ReturnStr tmp0`,
		},

		`return imul(imul(imul(1, 2), 3), 4)`: {
			`testpkg.f code=32 frame=168 (7 slots: 4 args, 0 locals, 3 temps)`,
			`  LoadScalarConst arg0 = 1`,
			`  LoadScalarConst arg1 = 2`,
			`  CallNative tmp2 = testpkg.imul()`,
			`  MoveScalar arg0 = tmp2`,
			`  LoadScalarConst arg1 = 3`,
			`  CallNative tmp1 = testpkg.imul()`,
			`  MoveScalar arg0 = tmp1`,
			`  LoadScalarConst arg1 = 4`,
			`  CallNative tmp0 = testpkg.imul()`,
			`  ReturnScalar tmp0`,
		},

		`x1 := 1; x2 := 2; x3 := 3; x4 := 4; return imul(x1, imul(x2, imul(x3, x4)))`: {
			`testpkg.f code=50 frame=312 (13 slots: 4 args, 4 locals, 5 temps)`,
			`  LoadScalarConst x1 = 1`,
			`  LoadScalarConst x2 = 2`,
			`  LoadScalarConst x3 = 3`,
			`  LoadScalarConst x4 = 4`,
			`  MoveScalar tmp1 = x1`,
			`  MoveScalar tmp3 = x2`,
			`  MoveScalar arg0 = x3`,
			`  MoveScalar arg1 = x4`,
			`  CallNative tmp4 = testpkg.imul()`,
			`  MoveScalar arg0 = tmp3`,
			`  MoveScalar arg1 = tmp4`,
			`  CallNative tmp2 = testpkg.imul()`,
			`  MoveScalar arg0 = tmp1`,
			`  MoveScalar arg1 = tmp2`,
			`  CallNative tmp0 = testpkg.imul()`,
			`  ReturnScalar tmp0`,
		},

		`x := 10; y := 20; return (x == 0 || x > 0) && (y < 5 || y >= 10)`: {
			`testpkg.f code=48 frame=264 (11 slots: 4 args, 2 locals, 5 temps)`,
			`  LoadScalarConst x = 10`,
			`  LoadScalarConst y = 20`,
			`  LoadScalarConst tmp1 = 0`,
			`  ScalarEq tmp0 = x tmp1`,
			`  JumpTrue L0 tmp0`,
			`  LoadScalarConst tmp2 = 0`,
			`  IntGt tmp0 = x tmp2`,
			`L0:`,
			`  JumpFalse L1 tmp0`,
			`  LoadScalarConst tmp3 = 5`,
			`  IntLt tmp0 = y tmp3`,
			`  JumpTrue L1 tmp0`,
			`  LoadScalarConst tmp4 = 10`,
			`  IntGtEq tmp0 = y tmp4`,
			`L1:`,
			`  ReturnScalar tmp0`,
		},
	}

	makePackageSource := func(body string) string {
		return `
		  package ` + testPackage + `
		  func f(i int, s string, b bool, err error) interface{} {
			` + body + `
		  }
		  func concat(x, y string) string
		  func imul(x, y int) int
		  func sprintf(format string, args ...interface{}) string
		  `
	}

	for testSrc, disasmLines := range tests {
		env := quasigo.NewEnv()
		env.AddNativeFunc(testPackage, "concat", func(ctx qnative.CallContext) {
			panic("should not be called")
		})
		env.AddNativeFunc(testPackage, "imul", func(ctx qnative.CallContext) {
			panic("should not be called")
		})
		env.AddNativeFunc(testPackage, "sprintf", func(ctx qnative.CallContext) {
			panic("should not be called")
		})
		src := makePackageSource(testSrc)
		parsed, err := testutil.ParseGoFile(testPackage, src)
		if err != nil {
			t.Fatalf("parse %s: %v", testSrc, err)
		}
		compiled, err := testutil.CompileOptTestFile(env, "f", testPackage, parsed)
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
