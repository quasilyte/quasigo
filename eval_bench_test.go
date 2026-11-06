package quasigo_test

import (
	"fmt"
	"testing"

	"github.com/quasilyte/quasigo"
	"github.com/quasilyte/quasigo/internal/evaltest"
	"github.com/quasilyte/quasigo/internal/testutil"
	"github.com/quasilyte/quasigo/qnative"
	"github.com/quasilyte/quasigo/stdlib/qfmt"
)

const testPackage = "testpkg"

type benchTestCase struct {
	name   string
	src    string
	params string
	args   []any
}

// bindArgs pushes the args using their dynamic types.
// It only supports the types that are used in these tests.
func bindArgs(env *quasigo.EvalEnv, args ...any) {
	pusher := env.PrepareArgs()
	for _, arg := range args {
		switch arg := arg.(type) {
		case nil:
			pusher.PushInterface(nil)
		case bool:
			pusher.PushBool(arg)
		case byte:
			pusher.PushByte(arg)
		case int:
			pusher.PushInt(arg)
		case float64:
			pusher.PushFloat(arg)
		case string:
			pusher.PushString(arg)
		case []int:
			pusher.PushSlice(arg)
		case *evaltest.Point:
			pusher.PushPtr(arg)
		case []*evaltest.Point:
			pusher.PushSlice(arg)
		default:
			panic(fmt.Sprintf("bindArgs: unsupported %T arg", arg))
		}
	}
}

var benchmarksNoAlloc = []*benchTestCase{
	{
		name: `ReturnFalse`,
		src:  `return false`,
	},

	{
		name: `Return0`,
		src:  `return 0`,
	},

	{
		name: `ReturnInt`,
		src:  `return 384723`,
	},

	{
		name: `StrSlicing`,
		src:  `s1 := "hello"; s2 := s1[1:]; s3 := s2[:3]; return s3[1:2]`,
	},

	{
		name: `StrLen`,
		src:  `s1 := "hello"; return len(s1) + len(s1) + len(s1)`,
	},

	{
		name:   `ParamInt`,
		src:    `return x + y + z`,
		params: `x, y, z int`,
		args:   []any{10, 20, 30},
	},

	{
		name:   `IntSliceIndex`,
		src:    `xs[0] = 10; return xs[0]`,
		params: `xs []int`,
		args:   []any{[]int{0}},
	},

	{
		name: `LocalVars`,
		src:  `x := 1; y := x; return y`,
	},

	{
		name: `SimpleArith`,
		src:  `x := 0; y := x + 1; return x + y + x + 5`,
	},

	{
		name: `SimpleArith2`,
		src:  `x := 1; y := 2; v1 := x + y; v2 := v1 + v1; return v1 + v2`,
	},

	{
		name: "CondExpr",
		src:  `x := 10; y := 20; return (x == 0 || x > 0) && (y < 5 || y >= 10)`,
	},

	{
		name: `IfStmt`,
		src:  `x := 100; if x == 1 { x = 10 } else if x == 2 { x = 20 } else { x = 30 }; return x`,
	},

	{
		name: `CallNative`,
		src:  `return imul(1, 5) + imul(2, 2)`,
	},

	{
		name: `CallNative2`,
		src:  `return imul(imul(imul(1, 2), 3), 4)`,
	},

	{
		name: `Call0`,
		src:  `return fn0()`,
	},

	{
		name: `Call1`,
		src:  `return fn1(1)`,
	},

	{
		name: `Call2`,
		src:  `return fn2(1, 1)`,
	},

	{
		name: `CallNested`,
		src:  `return nestedcall(1)`,
	},

	{
		name: `CallLoop`,
		src:  `res := false; for i := 0; i < 100; i++ { res = nestedcall(i); res = fn2(1, 1); }; return res`,
	},

	{
		name: `Sqrt1000`,
		src:  `return sqrt(1000)`,
	},

	{
		name: `CallRet2`,
		src:  `a, b := ret2(1, 2); return a - b`,
	},

	{
		name: `CallRet2Str`,
		src:  `s, n := ret2str("a", 2); return s == "a" && n == 2`,
	},

	{
		name: `CallRet2Discard`,
		src:  `_, b := ret2(1, 2); return b`,
	},

	{
		name: `CallRet2Loop`,
		src:  `a := 1; b := 2; for i := 0; i < 100; i++ { a, b = ret2(a, b) }; return a - b`,
	},

	{
		name: `CallRecur5`,
		src:  `return recur(5)`,
	},

	{
		name: `CallRecur20`,
		src:  `return recur(20)`,
	},

	{
		name: `CallRecur40`,
		src:  `return recur(40)`,
	},

	{
		name: `CounterLoop`,
		src:  `j := 0; for j < 10000 { j++ }; return j`,
	},

	{
		name: `CounterLoopNested`,
		src:  `j := 0; for j < 10000 { k := 0; for k < 10 { k++; j++; } }; return j`,
	},
}

func TestNoAllocs(t *testing.T) {
	for _, test := range benchmarksNoAlloc {
		env, compiled := compileBenchFunc(t, test.params, test.src)
		evalEnv := env.GetEvalEnv(2 * 4096)
		bindArgs(evalEnv, test.args...)

		allocs := testing.AllocsPerRun(5, func() {
			quasigo.Call(evalEnv, compiled)
		})
		if int(allocs) != 0 {
			t.Errorf("%s does allocate (%d allocs)", test.name, int(allocs))
		}
	}
}

func TestPushArgs(t *testing.T) {
	points := evaltest.NewPoints(3)
	points[1] = nil
	tests := []struct {
		params string
		src    string
		args   []any
		want   int
	}{
		{`b bool`, `if b { return 1 }; return 0`, []any{true}, 1},
		{`b byte`, `return int(b)`, []any{byte('A')}, 65},
		{`x float64`, `if x > 1.5 { return 1 }; return 0`, []any{2.5}, 1},
		{`s string, i int, xs []int`, `return len(s) + xs[i]`, []any{"abc", 1, []int{10, 20}}, 23},
		{`err error`, `if err == nil { return 1 }; return 0`, []any{nil}, 1},
		{`p *evaltest.Point`, `return p.X`, []any{evaltest.NewPoint(7, 70)}, 7},
		{`xs []*evaltest.Point, i int`, `return xs[i].X`, []any{points, 2}, 2},
		{`xs []*evaltest.Point`, `if xs[1] == nil { return len(xs) }; return 0`, []any{points}, 3},
	}
	for _, test := range tests {
		env, compiled := compileBenchFunc(t, test.params, test.src)
		evalEnv := env.GetEvalEnv(4096)
		bindArgs(evalEnv, test.args...)
		have := quasigo.Call(evalEnv, compiled).IntValue()
		if have != test.want {
			t.Errorf("%s: have %d, want %d", test.src, have, test.want)
		}
	}
}

func TestPushArgsNoAllocs(t *testing.T) {
	env, compiled := compileBenchFunc(t, `s string, i int, xs []int, p *evaltest.Point`, `return len(s) + xs[i] + p.X`)
	evalEnv := env.GetEvalEnv(4096)
	xs := []int{10, 20}
	p := evaltest.NewPoint(7, 70)
	var have int
	allocs := testing.AllocsPerRun(10, func() {
		args := evalEnv.PrepareArgs()
		args.PushString("abc")
		args.PushInt(1)
		args.PushSlice(xs)
		args.PushPtr(p)
		have = quasigo.Call(evalEnv, compiled).IntValue()
	})
	if int(allocs) != 0 {
		t.Errorf("args pushing does allocate (%d allocs)", int(allocs))
	}
	if have != 30 {
		t.Errorf("have %d, want 30", have)
	}
}

func BenchmarkEval(b *testing.B) {
	var tests = []*benchTestCase{
		{
			name: `CallNativeVariadic0`,
			src:  `return fmt.Sprintf("no formatting")`,
		},
		{
			name: `CallNativeVariadic1`,
			src:  `return fmt.Sprintf("Hello, %s!", "world")`,
		},
		{
			name: `CallNativeVariadic2`,
			src:  `return fmt.Sprintf("%s:%d", "foo.go", 105)`,
		},
		{
			name: `CallNativeVariadic3`,
			src:  `message := "bad"; return fmt.Sprintf("%s:%d: %s", "foo.go", 105, message)`,
		},
		{
			name: "CallNativeVariadicInts",
			src:  `return variadicInts(1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000)`,
		},
		{
			name: "CallNativeVariadicFloats",
			src:  `return variadicFloats(1000.1, 2000.2, 3000.3, 4000.4, 5000.5, 6000.6, 7000.7, 8000.8)`,
		},
	}

	tests = append(tests, benchmarksNoAlloc...)

	runBench := func(b *testing.B, env *quasigo.EvalEnv, fn quasigo.Func) {
		for i := 0; i < b.N; i++ {
			_ = quasigo.Call(env, fn)
		}
	}

	for _, test := range tests {
		test := test
		b.Run(test.name, func(b *testing.B) {
			env, compiled := compileBenchFunc(b, test.params, test.src)
			evalEnv := env.GetEvalEnv(2 * 4096)
			bindArgs(evalEnv, test.args...)
			b.ResetTimer()
			runBench(b, evalEnv, compiled)
		})
	}
}

func compileBenchFunc(t testing.TB, paramsSig, bodySrc string) (*quasigo.Env, quasigo.Func) {
	makePackageSource := func(body string) string {
		return `
package ` + testPackage + `
import "fmt"
import "github.com/quasilyte/quasigo/internal/evaltest"
var _ = fmt.Sprintf
var _ = evaltest.NewPoint
func fn0() bool { return false }
func fn1(x int) bool { return false }
func fn2(x, y int) bool { return false }
func variadicInts(xs ...int) int { return 0 }
func variadicFloats(xs ...float64) float64 { return 0 }

func recur(x int) bool {
	if x == 0 {
		return true
	}
	return recur(x - 1)
}

func ret2(x, y int) (int, int) { return y, x }
func ret2str(s string, n int) (string, int) { return s, n }

func sqrt(x int) int {
	if x == 0 || x == 1 {
		return x
	}
	start := 1
	end := x / 2
	result := 0
	for start <= end {
		mid := (start + end) / 2
		sqr := mid * mid
		if sqr == x {
			return mid
		}
		if sqr <= x {
			start = mid + 1
			result = mid
		} else {
			end = mid - 1
		}
	}
	return result
}

func _nestedcall4(x int) bool { return false }
func _nestedcall3(x int) bool { return _nestedcall4(x) }
func _nestedcall2(x int) bool { return _nestedcall3(x) }
func _nestedcall1(x int) bool { return _nestedcall2(x) }
func nestedcall(x int) bool { return _nestedcall1(x) }

func f(` + paramsSig + `) any {
	` + body + `
}
func imul(x, y int) int
`
	}

	env := quasigo.NewEnv()
	env.AddNativeFunc(testPackage, "variadicInts", func(ctx qnative.CallContext) {
		ctx.SetIntResult(11)
	})
	env.AddNativeFunc(testPackage, "variadicFloats", func(ctx qnative.CallContext) {
		ctx.SetFloatResult(22)
	})
	env.AddNativeFunc(testPackage, "imul", func(ctx qnative.CallContext) {
		x := ctx.IntArg(0)
		y := ctx.IntArg(1)
		ctx.SetIntResult(x * y)
	})
	env.AddNativeMethod(`*github.com/quasilyte/quasigo/internal/evaltest.Point`, `X`, func(ctx qnative.CallContext) {
		ctx.SetIntResult(ctx.PtrArg[evaltest.Point](0).X)
	})
	qfmt.ImportAll(env)
	src := makePackageSource(bodySrc)
	parsed, err := testutil.ParseGoFile(testPackage, src)
	if err != nil {
		t.Fatalf("parse %s: %v", bodySrc, err)
	}
	compiled, err := testutil.CompileOptTestFile(env, "f", testPackage, parsed)
	if err != nil {
		t.Fatalf("compile %s: %v", bodySrc, err)
	}
	return env, compiled
}
