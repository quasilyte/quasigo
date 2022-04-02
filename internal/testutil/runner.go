package testutil

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/quasilyte/quasigo"
	"github.com/quasilyte/quasigo/internal/ir"
	"github.com/quasilyte/quasigo/internal/qruntime"
	"github.com/quasilyte/quasigo/qnative"
	"github.com/quasilyte/quasigo/stdlib/qfmt"
	"github.com/quasilyte/quasigo/stdlib/qstrconv"
	"github.com/quasilyte/quasigo/stdlib/qstrings"
)

type Runner struct {
	root *testing.T

	workdir string

	Targets []RunnerTarget

	NewEnv func() *quasigo.Env
}

type RunnerTarget struct {
	Name string
	Path string
}

func NewRunner(t *testing.T) *Runner {
	return &Runner{root: t}
}

func (r *Runner) Run() {
	{
		workdir, err := os.Getwd()
		if err != nil {
			r.root.Fatalf("getwd: %v", err)
		}
		r.workdir = workdir
	}

	if len(r.Targets) == 0 {
		r.Targets = r.testdataTargets()
		if len(r.Targets) == 0 {
			r.root.Fatalf("no targets provided and testdata is empty")
		}
	}

	for _, target := range r.Targets {
		r.root.Run(target.Name, func(t *testing.T) {
			r.runTarget(t, target)
		})
	}
}

func (r *Runner) testdataTargets() []RunnerTarget {
	files, err := os.ReadDir("testdata")
	if err != nil {
		r.root.Fatalf("find testdata targets: %v", err)
	}
	targets := make([]RunnerTarget, 0, len(files))
	for _, d := range files {
		if !d.IsDir() {
			continue
		}
		absPath := filepath.Join(r.workdir, "testdata", d.Name())
		targets = append(targets, RunnerTarget{
			Name: d.Name(),
			Path: absPath,
		})
	}
	return targets
}

func (r *Runner) runTarget(t *testing.T, target RunnerTarget) {
	fset := token.NewFileSet()
	packages, err := parser.ParseDir(fset, target.Path, nil, parser.ParseComments|parser.SkipObjectResolution)
	if err != nil {
		t.Fatalf("parse %s dir: %v", target.Name, err)
	}
	if len(packages) != 1 {
		t.Fatalf("%s: expected 1 package, found %d", target.Name, len(packages))
	}
	var pkg *ast.Package
	for k := range packages {
		pkg = packages[k]
		break
	}
	var pkgFiles []*ast.File
	for _, f := range pkg.Files {
		pkgFiles = append(pkgFiles, f)
	}
	typesConfig := &types.Config{Importer: importer.Default()}
	typesInfo := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Uses:  make(map[*ast.Ident]types.Object),
		Defs:  make(map[*ast.Ident]types.Object),
	}
	typesPkg, err := typesConfig.Check(pkg.Name, fset, pkgFiles, typesInfo)
	if err != nil {
		t.Fatalf("typecheck %s: %v", target.Name, err)
	}

	testPackage := &testPackage{
		typesPackage: typesPkg,
		typesInfo:    typesInfo,
		fset:         fset,
		files:        make([]*testFile, len(pkgFiles)),
	}
	for i := range pkgFiles {
		f := r.newTestFile(t, fset, pkgFiles[i])
		f.pkg = testPackage
		testPackage.files[i] = f
	}

	// Run Go only once.
	// This output will be used for both normal and optimized quasigo runs.
	var goResult string
	if typesPkg.Path() == "main" {
		relpath, err := filepath.Rel(r.workdir, target.Path)
		if err != nil {
			t.Fatalf("get relpath: %v", err)
		}
		goResult = r.runGo(t, "./"+relpath)
	}

	r.runTest(t, testPackage, goResult, false)
	r.runTest(t, testPackage, goResult, true)
}

func (r *Runner) runTest(t *testing.T, pkg *testPackage, goResult string, optimize bool) {
	env := r.newTestEnv()

	compileContext := &quasigo.CompileContext{
		Env:      env.handle,
		Package:  pkg.typesPackage,
		Types:    pkg.typesInfo,
		Sizes:    types.SizesFor("gc", runtime.GOARCH),
		Fset:     pkg.fset,
		Optimize: optimize,
		Static:   true,
	}
	checkDisasm := false
	checkIR := false
	checkConstants := false
	for _, f := range pkg.files {
		if len(f.GetConstantsChecks(optimize)) != 0 {
			checkConstants = true
		}
		if len(f.GetIRDumpChecks(optimize)) != 0 {
			checkIR = true
		}
		if len(f.GetDisasmChecks(optimize)) != 0 {
			checkDisasm = true
		}
	}
	var irdumpResults map[string]string
	if checkIR {
		irdumpResults = make(map[string]string)
		compileContext.TestingContext = &compilerTestingContext{
			OnFuncIR: func(fn *ir.Func) {
				dump := ir.Dump(fn)
				irdumpResults[fn.Name] = dump
			},
		}
	}
	for _, f := range pkg.files {
		r.compileQuasigo(t, compileContext, f.syntax)
	}

	suffix := ""
	if optimize {
		suffix = "_opt"
	}

	if checkIR {
		t.Run("irdump", func(t *testing.T) {
			for _, f := range pkg.files {
				r.checkIR(t, f, irdumpResults, f.GetIRDumpChecks(optimize))
			}
		})
	}

	if checkDisasm {
		t.Run("disasm"+suffix, func(t *testing.T) {
			for _, f := range pkg.files {
				r.checkDisasm(t, f, env, f.GetDisasmChecks(optimize))
			}
		})
	}

	if checkConstants {
		t.Run("constants"+suffix, func(t *testing.T) {
			for _, f := range pkg.files {
				r.checkConstants(t, f, env, f.GetConstantsChecks(optimize))
			}
		})
	}

	if pkg.typesPackage.Path() == "main" {
		t.Run("exec"+suffix, func(t *testing.T) {
			quasigoResult := r.runQuasigo(t, env)
			if diff := cmp.Diff(quasigoResult, goResult); diff != "" {
				t.Fatalf("output mismatch (-have +want):\n%s", diff)
			}
		})
	}
}

func (r *Runner) checkIR(t *testing.T, f *testFile, irdump map[string]string, checks []disasmCheck) {
	for _, c := range checks {
		funcKey := qruntime.FuncKey{
			Qualifier: f.pkg.typesPackage.Path(),
			Name:      c.funcName,
		}
		dump, ok := irdump[funcKey.String()]
		if !ok {
			t.Fatalf("can't find IR dump for %s", funcKey)
		}
		have := splitLines(dump)
		want := splitLines(c.expected)
		if diff := cmp.Diff(have, want); diff != "" {
			t.Errorf("%s:%d: irdump mismatch (-have +want):\n%s", f.name, c.line, diff)
			fmt.Println("For copy/paste:")
			for _, l := range strings.Split(dump, "\n") {
				if l == "" {
					continue
				}
				fmt.Printf("// %s\n", l)
			}
			t.FailNow()
		}
	}
}

func (r *Runner) checkConstants(t *testing.T, f *testFile, env *testEnv, checks []disasmCheck) {
	for _, c := range checks {
		fn := env.handle.GetFunc(f.pkg.typesPackage.Path(), c.funcName)

		var buf strings.Builder
		buf.WriteString("  scalar constants: [")
		for i, v := range fn.ScalarConstants() {
			if i != 0 {
				buf.WriteByte(' ')
			}
			buf.WriteString(strconv.FormatInt(int64(v), 10))
		}
		buf.WriteString("]\n")
		buf.WriteString("  string constants: [")
		for i, v := range fn.StringConstants() {
			if i != 0 {
				buf.WriteByte(' ')
			}
			buf.WriteByte('"')
			buf.WriteString(v)
			buf.WriteByte('"')
		}
		buf.WriteString("]\n")

		have := splitLines(buf.String())
		want := splitLines(c.expected)
		if diff := cmp.Diff(have, want); diff != "" {
			t.Errorf("%s:%d: constants mismatch (-have +want):\n%s", f.name, c.line, diff)
			fmt.Println("For copy/paste:")
			for _, l := range strings.Split(buf.String(), "\n") {
				if l == "" {
					continue
				}
				fmt.Printf("// %s\n", l)
			}
			t.FailNow()
		}
	}
}

func (r *Runner) checkDisasm(t *testing.T, f *testFile, env *testEnv, checks []disasmCheck) {
	for _, c := range checks {
		fn := env.handle.GetFunc(f.pkg.typesPackage.Path(), c.funcName)
		disasmOutput := quasigo.Disasm(env.handle, fn)
		have := splitLines(disasmOutput)
		want := splitLines(c.expected)
		if diff := cmp.Diff(have, want); diff != "" {
			t.Errorf("%s:%d: disasm mismatch (-have +want):\n%s", f.name, c.line, diff)
			fmt.Println("For copy/paste:")
			for _, l := range strings.Split(disasmOutput, "\n") {
				if l == "" {
					continue
				}
				fmt.Printf("// %s\n", l)
			}
			t.FailNow()
		}
	}
}

func (r *Runner) compileQuasigo(t *testing.T, ctx *quasigo.CompileContext, f *ast.File) {
	filename := filepath.Base(ctx.Fset.Position(f.Pos()).Filename)
	for _, decl := range f.Decls {
		decl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if decl.Body == nil {
			continue
		}
		fn, err := quasigo.Compile(ctx, decl)
		if err != nil {
			t.Fatalf("%s: compile %s: %v", filename, decl.Name, err)
		}
		ctx.Env.AddFunc(ctx.Package.Path(), decl.Name.String(), fn)
	}
}

func (r *Runner) newTestFile(t *testing.T, fset *token.FileSet, syntax *ast.File) *testFile {
	f := &testFile{
		name:   filepath.Base(fset.Position(syntax.Pos()).Filename),
		syntax: syntax,
	}

	for _, decl := range syntax.Decls {
		decl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if decl.Doc != nil {
			r.loadFuncComments(t, fset, f, decl.Name.Name, decl.Doc)
		}
	}

	return f
}

func (r *Runner) newTestEnv() *testEnv {
	var env *quasigo.Env
	if r.NewEnv != nil {
		env = r.NewEnv()
	} else {
		env = quasigo.NewEnv()
	}
	stdout := &bytes.Buffer{}
	{
		env.AddNativeFunc(`builtin`, `PrintString`, func(ctx qnative.CallContext) {
			fmt.Fprintln(stdout, ctx.StringArg(0))
		})
		env.AddNativeFunc(`builtin`, `PrintByte`, func(ctx qnative.CallContext) {
			fmt.Fprintln(stdout, ctx.ByteArg(0))
		})
		env.AddNativeFunc(`builtin`, `PrintInt`, func(ctx qnative.CallContext) {
			fmt.Fprintln(stdout, ctx.IntArg(0))
		})
		env.AddNativeFunc(`builtin`, `PrintBool`, func(ctx qnative.CallContext) {
			fmt.Fprintln(stdout, ctx.BoolArg(0))
		})
	}

	env.AddNativeMethod(`error`, `Error`, func(ctx qnative.CallContext) {
		err := ctx.InterfaceArg(0).(error)
		ctx.SetStringResult(err.Error())
	})

	qstrings.ImportAll(env)
	qstrconv.ImportAll(env)
	qfmt.ImportAll(env)

	return &testEnv{
		handle: env,
		stdout: stdout,
	}
}

func (r *Runner) loadFuncComments(t *testing.T, fset *token.FileSet, f *testFile, funcName string, cg *ast.CommentGroup) {
	var disasmBoth strings.Builder
	var disasmBothLine int
	var disasm strings.Builder
	var disasmLine int
	var disasmOpt strings.Builder
	var disasmOptLine int
	var irdump strings.Builder
	var irdumpLine int
	var constants strings.Builder
	var constantsLine int
	var constantsOpt strings.Builder
	var constantsOptLine int

	var currentSection *strings.Builder
	for _, c := range cg.List {
		if !strings.HasPrefix(c.Text, "//") {
			continue
		}
		s := strings.TrimSpace(strings.TrimPrefix(c.Text, "//"))
		switch {
		case strings.HasPrefix(s, "test:constants_opt"):
			currentSection = &constantsOpt
			constantsOptLine = fset.Position(c.Pos()).Line
			continue
		case strings.HasPrefix(s, "test:constants"):
			currentSection = &constants
			constantsLine = fset.Position(c.Pos()).Line
			continue
		case strings.HasPrefix(s, "test:irdump"):
			currentSection = &irdump
			irdumpLine = fset.Position(c.Pos()).Line
			continue
		case strings.HasPrefix(s, "test:disasm_opt"):
			currentSection = &disasmOpt
			disasmOptLine = fset.Position(c.Pos()).Line
			continue
		case strings.HasPrefix(s, "test:disasm_both"):
			currentSection = &disasmBoth
			disasmBothLine = fset.Position(c.Pos()).Line
			continue
		case strings.HasPrefix(s, "test:disasm"):
			currentSection = &disasm
			disasmLine = fset.Position(c.Pos()).Line
			continue
		case s == "":
			currentSection = nil
			continue
		}
		if currentSection == nil {
			continue
		}
		currentSection.WriteString(s)
		currentSection.WriteByte('\n')
	}

	if disasmBoth.Len() != 0 {
		if disasmOpt.Len() != 0 || disasm.Len() != 0 {
			t.Fatalf("used disasm_both with other disasm test directive")
		}
		disasm.WriteString(disasmBoth.String())
		disasmLine = disasmBothLine
		disasmOpt.WriteString(disasmBoth.String())
		disasmOptLine = disasmBothLine
	}
	if disasm.Len() != 0 {
		f.disasm = append(f.disasm, disasmCheck{
			line:     disasmLine,
			funcName: funcName,
			expected: disasm.String(),
		})
	}
	if disasmOpt.Len() != 0 {
		f.disasmOpt = append(f.disasmOpt, disasmCheck{
			line:     disasmOptLine,
			funcName: funcName,
			expected: disasmOpt.String(),
		})
	}
	if irdump.Len() != 0 {
		f.irdump = append(f.irdump, disasmCheck{
			line:     irdumpLine,
			funcName: funcName,
			expected: irdump.String(),
		})
	}
	if constants.Len() != 0 {
		f.constants = append(f.constants, disasmCheck{
			line:     constantsLine,
			funcName: funcName,
			expected: constants.String(),
		})
	}
	if constantsOpt.Len() != 0 {
		f.constantsOpt = append(f.constantsOpt, disasmCheck{
			line:     constantsOptLine,
			funcName: funcName,
			expected: constantsOpt.String(),
		})
	}
}

func (r *Runner) runGo(t *testing.T, main string) string {
	out, err := exec.Command("go", "run", main).CombinedOutput()
	if err != nil {
		t.Fatalf("run go: %v: %s", err, out)
	}
	return string(out)
}

func (r *Runner) runQuasigo(t *testing.T, env *testEnv) string {
	env.stdout.Reset()
	mainFunc := env.handle.GetFunc("main", "main")
	if mainFunc.IsNil() {
		t.Fatalf("can't find main function")
	}
	evalEnv := env.handle.GetEvalEnv(4096)
	quasigo.Call(evalEnv, mainFunc)
	return env.stdout.String()
}

func splitLines(s string) []string {
	lines := strings.Split(s, "\n")
	if len(lines) == 0 {
		return lines
	}
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

type testEnv struct {
	handle *quasigo.Env
	stdout *bytes.Buffer
}

type testPackage struct {
	typesInfo    *types.Info
	typesPackage *types.Package
	fset         *token.FileSet
	files        []*testFile
}

type testFile struct {
	name         string
	syntax       *ast.File
	disasm       []disasmCheck
	disasmOpt    []disasmCheck
	irdump       []disasmCheck
	constants    []disasmCheck
	constantsOpt []disasmCheck
	pkg          *testPackage
}

func (f *testFile) GetDisasmChecks(optimize bool) []disasmCheck {
	if optimize {
		return f.disasmOpt
	}
	return f.disasm
}

func (f *testFile) GetIRDumpChecks(optimize bool) []disasmCheck {
	if optimize {
		return f.irdump
	}
	return nil
}

func (f *testFile) GetConstantsChecks(optimize bool) []disasmCheck {
	if optimize {
		return f.constantsOpt
	}
	return f.constants
}

type disasmCheck struct {
	line     int
	funcName string
	expected string
}

type compilerTestingContext struct {
	OnFuncIR func(*ir.Func)
}

func (ctx *compilerTestingContext) FuncIR(fn *ir.Func) {
	ctx.OnFuncIR(fn)
}
