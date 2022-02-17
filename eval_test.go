package quasigo_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/quasilyte/quasigo"
	"github.com/quasilyte/quasigo/internal/evaltest"
	"github.com/quasilyte/quasigo/qnative"
	"github.com/quasilyte/quasigo/stdlib/qfmt"
	"github.com/quasilyte/quasigo/stdlib/qstrconv"
	"github.com/quasilyte/quasigo/stdlib/qstrings"
)

func TestEval(t *testing.T) {
	type testCase struct {
		src    string
		result interface{}
	}

	exprTests := []testCase{
		// Const literals.
		{`1`, 1},
		{`"foo"`, "foo"},
		{`true`, true},
		{`false`, false},

		// Function args.
		{`b`, true},
		{`i`, 10},

		// Arith operators.
		{`5 + i`, 15},
		{`i + i`, 20},
		{`i - 5`, 5},
		{`5 - i`, -5},
		{`i * 3`, 30},
		{`i / 2`, 5},
		{`(i * 3) / 10`, 3},

		// String operators.
		{`s + s`, "foofoo"},

		// Bool operators.
		{`!b`, false},
		{`!!b`, true},
		{`i == 2`, false},
		{`i == 10`, true},
		{`i >= 10`, true},
		{`i >= 9`, true},
		{`i >= 11`, false},
		{`i > 10`, false},
		{`i > 9`, true},
		{`i > -1`, true},
		{`i < 10`, false},
		{`i < 11`, true},
		{`i <= 10`, true},
		{`i <= 11`, true},
		{`i != 2`, true},
		{`i != 10`, false},
		{`s != "foo"`, false},
		{`s != "bar"`, true},

		// || operator.
		{`i == 2 || i == 10`, true},
		{`i == 10 || i == 2`, true},
		{`i == 2 || i == 3 || i == 10`, true},
		{`i == 2 || i == 10 || i == 3`, true},
		{`i == 10 || i == 2 || i == 3`, true},
		{`!(i == 10 || i == 2 || i == 3)`, false},

		// && operator.
		{`i == 10 && s == "foo"`, true},
		{`i == 10 && s == "foo" && true`, true},
		{`i == 20 && s == "foo"`, false},
		{`i == 10 && s == "bar"`, false},
		{`i == 10 && s == "foo" && false`, false},

		// String slicing.
		{`s[:]`, "foo"},
		{`s[0:]`, "foo"},
		{`s[1:]`, "oo"},
		{`s[:1]`, "f"},
		{`s[:0]`, ""},
		{`s[1:2]`, "o"},
		{`s[1:3]`, "oo"},

		// Builtin len().
		{`len(s)`, 3},
		{`len(s) == 3`, true},
		{`len(s[1:])`, 2},

		// Slicing with len().
		{`s[:len(s)-1]`, "fo"},

		// Native func call.
		{`imul(2, 3)`, 6},
		{`idiv(9, 3)`, 3},
		{`idiv(imul(2, 3), 1 + 1)`, 3},

		// Nil checks.
		{`nilEface() == nil`, true},
		{`nil == nilEface()`, true},
		{`nilEface() != nil`, false},
		{`nil != nilEface()`, false},
	}

	/*
		exprTests := []testCase{
			// Accesing the fields.
			{`foo.Prefix`, "Hello"},



		}
	*/

	tests := []testCase{
		{`if b { return 1 }; return 0`, 1},
		{`if !b { return 1 }; return 0`, 0},
		{`if b { return 1 } else { return 0 }`, 1},
		{`if !b { return 1 } else { return 0 }`, 0},

		{`x := 2; if x == 2 { return "a" } else if x == 0 { return "b" }; return "c"`, "a"},
		{`x := 2; if x == 0 { return "a" } else if x == 2 { return "b" }; return "c"`, "b"},
		{`x := 2; if x == 0 { return "a" } else if x == 1 { return "b" }; return "c"`, "c"},
		{`x := 2; if x == 2 { return "a" } else if x == 0 { return "b" } else { return "c" }`, "a"},
		{`x := 2; if x == 0 { return "a" } else if x == 2 { return "b" } else { return "c" }`, "b"},
		{`x := 2; if x == 0 { return "a" } else if x == 1 { return "b" } else { return "c" }`, "c"},
		{`x := 0; if b { x = 5 } else { x = 50 }; return x`, 5},
		{`x := 0; if !b { x = 5 } else { x = 50 }; return x`, 50},
		{`x := 0; if b { x = 1 } else if x == 0 { x = 2 } else { x = 3 }; return x`, 1},
		{`x := 0; if !b { x = 1 } else if x == 0 { x = 2 } else { x = 3 }; return x`, 2},
		{`x := 0; if !b { x = 1 } else if x == 1 { x = 2 } else { x = 3 }; return x`, 3},

		{`x := true; y := !x; return y`, false},
		{`x := true; y := !x; return !y`, true},

		{`x := 1; y := 2; v1 := x + y; v2 := v1 + v1; return v1 + v2`, 9},

		{`x := 0; return x + 1`, 1},
		{`x := -10; return x + 1`, -9},
		{`x := 0; return x - 1`, -1},
		{`x := -10; return x - 1`, -11},
		{`x := 0; x++; return x`, 1},
		{`x := i; x++; return x`, 11},
		{`x := 0; x--; return x`, -1},
		{`x := i; x--; return x`, 9},

		{`j := 0; for { j = j + 1; break; }; return j`, 1},
		{`j := -5; for { if j > 0 { break }; j++; }; return j`, 1},
		{`j := -5; for { if j >= 0 { break }; j++; }; return j`, 0},
		{`j := 0; for j < 0 { j++; break; }; return j`, 0},
		{`j := -5; for j < 0 { j++ }; return j`, 0},
		{`j := -5; for j <= 0 { j++; }; return j`, 1},
		{`j := 0; for j < 100 { k := 0; for { if k > 40 { break }; k++; j++; } }; return j`, 123},
		{`j := 0; for j < 10000 { k := 0; for k < 10 { k++; j++; } }; return j`, 10000},

		// Multi-result native func call.
		{`quo, rem := idiv2(10, 3); return quo + rem`, 4},
		{`quo, rem := idiv2(10, 3); return quo == 3 && rem == 1`, true},
	}

	for _, test := range exprTests {
		test.src = `return ` + test.src
		tests = append(tests, test)
	}

	makePackageSource := func(body string, result interface{}) string {
		var returnType string
		switch result.(type) {
		case int:
			returnType = "int"
		case string:
			returnType = "string"
		case bool:
			returnType = "bool"
		default:
			t.Fatalf("unexpected result type: %T", result)
		}
		return `
		  package ` + testPackage + `
		  func target(i int, s string, b bool) ` + returnType + ` {
		    ` + body + `
		  }
		  func imul(x, y int) int
		  func idiv(x, y int) int
		  func idiv2(x, y int) (int, int)
		  func nilEface() interface{}
		  `
	}

	env := quasigo.NewEnv()

	env.AddNativeFunc(testPackage, "imul", func(ctx qnative.CallContext) {
		x := ctx.IntArg(0)
		y := ctx.IntArg(1)
		ctx.SetIntResult(x * y)
	})
	env.AddNativeFunc(testPackage, "idiv", func(ctx qnative.CallContext) {
		x := ctx.IntArg(0)
		y := ctx.IntArg(1)
		ctx.SetIntResult(x / y)
	})
	env.AddNativeFunc(testPackage, "idiv2", func(ctx qnative.CallContext) {
		x := ctx.IntArg(0)
		y := ctx.IntArg(1)
		quo := x / y
		rem := x % y
		ctx.SetIntResult(quo)
		ctx.SetIntResult2(rem)
	})
	env.AddNativeFunc(testPackage, "nilEface", func(ctx qnative.CallContext) {
		ctx.SetInterfaceResult(nil)
	})

	// env.AddNativeMethod(evaltestFoo, "Prefix", func(stack *quasigo.ValueStack) {
	// 	foo := stack.Pop().(*evaltest.Foo)
	// 	stack.Push(foo.Prefix)
	// })

	for i := range tests {
		test := tests[i]
		src := makePackageSource(test.src, test.result)
		parsed, err := parseGoFile(testPackage, src)
		if err != nil {
			t.Fatalf("parse %s: %v", test.src, err)
		}
		compiled, err := compileTestFunc(env, "target", parsed)
		if err != nil {
			t.Fatalf("compile %s: %v", test.src, err)
		}
		evalEnv := env.GetEvalEnv(1024)
		evalEnv.BindArgs(10, "foo", true, &evaltest.Foo{Prefix: "Hello"})
		result := quasigo.Call(evalEnv, compiled)
		var unboxedResult interface{}
		switch test.result.(type) {
		case int:
			unboxedResult = result.IntValue()
		case string:
			unboxedResult = result.StringValue()
		case bool:
			unboxedResult = result.BoolValue()
		default:
			t.Fatalf("can't unbox a value of type %T", test.result)
		}
		if unboxedResult != test.result {
			t.Fatalf("eval %s:\nhave: %#v\nwant: %#v", test.src, unboxedResult, test.result)
		}
	}
}

func TestEvalFile(t *testing.T) {
	files, err := ioutil.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}

	runGo := func(main string) (string, error) {
		out, err := exec.Command("go", "run", main).CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("%v: %s", err, out)
		}
		return string(out), nil
	}

	runQuasigo := func(main string) (string, error) {
		src, err := os.ReadFile(main)
		if err != nil {
			return "", err
		}
		env := quasigo.NewEnv()
		parsed, err := parseGoFile("main", string(src))
		if err != nil {
			return "", fmt.Errorf("parse: %v", err)
		}

		var stdout bytes.Buffer
		env.AddNativeFunc(`builtin`, `PrintString`, func(ctx qnative.CallContext) {
			fmt.Fprintln(&stdout, ctx.StringArg(0))
		})
		env.AddNativeFunc(`builtin`, `PrintInt`, func(ctx qnative.CallContext) {
			fmt.Fprintln(&stdout, ctx.IntArg(0))
		})
		env.AddNativeFunc(`builtin`, `PrintBool`, func(ctx qnative.CallContext) {
			fmt.Fprintln(&stdout, ctx.BoolArg(0))
		})

		env.AddNativeMethod(`error`, `Error`, func(ctx qnative.CallContext) {
			err := ctx.InterfaceArg(0).(error)
			ctx.SetStringResult(err.Error())
		})

		qstrings.ImportAll(env)
		qstrconv.ImportAll(env)
		qfmt.ImportAll(env)
		registerEvaltestPackage(env)

		mainFunc, err := compileTestFile(env, "main", "main", parsed)
		if err != nil {
			return "", err
		}
		if mainFunc.IsNil() {
			return "", errors.New("can't find main() function")
		}

		evalEnv := env.GetEvalEnv(4096)
		quasigo.Call(evalEnv, mainFunc)
		return stdout.String(), nil
	}

	runTest := func(t *testing.T, mainFile string) {
		goResult, err := runGo(mainFile)
		if err != nil {
			t.Fatalf("run go: %v", err)
		}
		quasigoResult, err := runQuasigo(mainFile)
		if err != nil {
			t.Fatalf("run quasigo: %v", err)
		}
		if diff := cmp.Diff(quasigoResult, goResult); diff != "" {
			t.Errorf("output mismatch:\nhave (+): `%s`\nwant (-): `%s`\ndiff: %s", quasigoResult, goResult, diff)
		}
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		mainFile := filepath.Join("testdata", f.Name(), "main.go")
		t.Run(f.Name(), func(t *testing.T) {
			runTest(t, mainFile)
		})
	}
}

func registerEvaltestPackage(env *quasigo.Env) {
	const evaltestPkgPath = `github.com/quasilyte/quasigo/internal/evaltest`
	const evaltestFoo = `*` + evaltestPkgPath + `.Foo`
	const evaltestFooPair = `*` + evaltestPkgPath + `.FooPair`

	env.AddNativeFunc(evaltestPkgPath, "NilFooAsEface", func(ctx qnative.CallContext) {
		ctx.SetInterfaceResult((*evaltest.Foo)(nil))
	})
	env.AddNativeFunc(evaltestPkgPath, "NewFoo", func(ctx qnative.CallContext) {
		prefix := ctx.StringArg(0)
		ctx.SetInterfaceResult(&evaltest.Foo{Prefix: prefix})
	})
	env.AddNativeFunc(evaltestPkgPath, "NilFoo", func(ctx qnative.CallContext) {
		ctx.SetInterfaceResult((*evaltest.Foo)(nil))
	})
	env.AddNativeFunc(evaltestPkgPath, "NilEface", func(ctx qnative.CallContext) {
		ctx.SetInterfaceResult(nil)
	})
	env.AddNativeFunc(evaltestPkgPath, "NewFooPair", func(ctx qnative.CallContext) {
		x := ctx.InterfaceArg(0).(*evaltest.Foo)
		y := ctx.InterfaceArg(1).(*evaltest.Foo)
		ctx.SetInterfaceResult(evaltest.NewFooPair(x, y))
	})

	env.AddNativeMethod(evaltestFooPair, "SetFirst", func(ctx qnative.CallContext) {
		p := ctx.InterfaceArg(0).(*evaltest.FooPair)
		x := ctx.InterfaceArg(1).(*evaltest.Foo)
		p.SetFirst(x)
	})
	env.AddNativeMethod(evaltestFooPair, "SetFirstPrefix", func(ctx qnative.CallContext) {
		p := ctx.InterfaceArg(0).(*evaltest.FooPair)
		prefix := ctx.StringArg(1)
		p.SetFirstPrefix(prefix)
	})
	env.AddNativeMethod(evaltestFooPair, "First", func(ctx qnative.CallContext) {
		p := ctx.InterfaceArg(0).(*evaltest.FooPair)
		ctx.SetInterfaceResult(p.First())
	})
	env.AddNativeMethod(evaltestFooPair, "Second", func(ctx qnative.CallContext) {
		p := ctx.InterfaceArg(0).(*evaltest.FooPair)
		ctx.SetInterfaceResult(p.Second())
	})
	env.AddNativeMethod(evaltestFooPair, "Get", func(ctx qnative.CallContext) {
		p := ctx.InterfaceArg(0).(*evaltest.FooPair)
		key := ctx.StringArg(1)
		ctx.SetInterfaceResult(p.Get(key))
	})

	env.AddNativeMethod(evaltestFoo, "Method1", func(ctx qnative.CallContext) {
		foo := ctx.InterfaceArg(0).(*evaltest.Foo)
		x := ctx.IntArg(1)
		ctx.SetStringResult(foo.Prefix + fmt.Sprint(x))
	})
	env.AddNativeMethod(evaltestFoo, "Prefix", func(ctx qnative.CallContext) {
		foo := ctx.InterfaceArg(0).(*evaltest.Foo)
		ctx.SetStringResult(foo.Prefix)
	})
	env.AddNativeMethod(evaltestFoo, "String", func(ctx qnative.CallContext) {
		foo := ctx.InterfaceArg(0).(*evaltest.Foo)
		ctx.SetStringResult(foo.Prefix)
	})
}
