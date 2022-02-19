package testutil

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"

	"github.com/quasilyte/quasigo"
)

type ParsedTestFile struct {
	Ast   *ast.File
	Pkg   *types.Package
	Types *types.Info
	Fset  *token.FileSet
}

func ParseGoFile(pkgPath, src string) (*ParsedTestFile, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		return nil, err
	}
	typechecker := &types.Config{
		Importer: importer.ForCompiler(fset, "source", nil),
	}
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Uses:  make(map[*ast.Ident]types.Object),
		Defs:  make(map[*ast.Ident]types.Object),
	}
	pkg, err := typechecker.Check(pkgPath, fset, []*ast.File{file}, info)
	result := &ParsedTestFile{
		Ast:   file,
		Pkg:   pkg,
		Types: info,
		Fset:  fset,
	}
	return result, err
}

func CompileOptTestFile(env *quasigo.Env, targetFunc, pkgPath string, parsed *ParsedTestFile) (quasigo.Func, error) {
	return compileTestFile(env, targetFunc, pkgPath, parsed, true)
}

func CompileTestFile(env *quasigo.Env, targetFunc, pkgPath string, parsed *ParsedTestFile) (quasigo.Func, error) {
	return compileTestFile(env, targetFunc, pkgPath, parsed, false)
}

func compileTestFile(env *quasigo.Env, targetFunc, pkgPath string, parsed *ParsedTestFile, opt bool) (quasigo.Func, error) {
	var resultFunc quasigo.Func
	for _, decl := range parsed.Ast.Decls {
		decl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if decl.Body == nil {
			continue
		}
		ctx := &quasigo.CompileContext{
			Env:      env,
			Package:  parsed.Pkg,
			Types:    parsed.Types,
			Fset:     parsed.Fset,
			Optimize: opt,
		}
		fn, err := quasigo.Compile(ctx, decl)
		if err != nil {
			return resultFunc, fmt.Errorf("compile %s func: %v", decl.Name, err)
		}
		if decl.Name.String() == targetFunc {
			resultFunc = fn
		} else {
			env.AddFunc(pkgPath, decl.Name.String(), fn)
		}
	}
	return resultFunc, nil
}

func CompileTestFunc(env *quasigo.Env, fn string, parsed *ParsedTestFile) (quasigo.Func, error) {
	var target *ast.FuncDecl
	for _, decl := range parsed.Ast.Decls {
		decl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if decl.Name.String() == fn {
			target = decl
			break
		}
	}
	if target == nil {
		return quasigo.Func{}, fmt.Errorf("test function %s not found", fn)
	}

	ctx := &quasigo.CompileContext{
		Env:     env,
		Package: parsed.Pkg,
		Types:   parsed.Types,
		Fset:    parsed.Fset,
	}
	return quasigo.Compile(ctx, target)
}
