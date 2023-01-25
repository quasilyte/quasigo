package quasigo_test

import (
	"testing"

	"github.com/quasilyte/quasigo"
	"github.com/quasilyte/quasigo/internal/testutil"
	"github.com/quasilyte/quasigo/qnative"
	"github.com/quasilyte/quasigo/stdlib/qfmt"
)

const testPackage = "testpkg"

type benchTestCase struct {
	name   string
	src    string
	params string
	args   []interface{}
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
		args:   []interface{}{10, 20, 30},
	},

	{
		name:   `IntSliceIndex`,
		src:    `xs[0] = 10; return xs[0]`,
		params: `xs []int`,
		args:   []interface{}{[]int{0}},
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
		evalEnv := env.GetEvalEnv(4096)
		evalEnv.BindArgs(test.args...)

		allocs := testing.AllocsPerRun(5, func() {
			quasigo.Call(evalEnv, compiled)
		})
		if int(allocs) != 0 {
			t.Errorf("%s does allocate (%d allocs)", test.name, int(allocs))
		}
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
			evalEnv := env.GetEvalEnv(4096)
			evalEnv.BindArgs(test.args...)
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
var _ = fmt.Sprintf
func fn0() bool { return false }
func fn1(x int) bool { return false }
func fn2(x, y int) bool { return false }

func recur(x int) bool {
	if x == 0 {
		return true
	}
	return recur(x - 1)
}

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

func f(` + paramsSig + `) interface{} {
	` + body + `
}
func imul(x, y int) int
`
	}

	env := quasigo.NewEnv()
	env.AddNativeFunc(testPackage, "imul", func(ctx qnative.CallContext) {
		x := ctx.IntArg(0)
		y := ctx.IntArg(1)
		ctx.SetIntResult(x * y)
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
