package qstrings

import (
	"strings"

	"github.com/quasilyte/quasigo"
	"github.com/quasilyte/quasigo/qnative"
)

func ImportAll(env *quasigo.Env) {
	env.AddNativeFunc(`strings`, `Replace`, Replace)
	env.AddNativeFunc(`strings`, `ReplaceAll`, ReplaceAll)
	env.AddNativeFunc(`strings`, `TrimPrefix`, TrimPrefix)
	env.AddNativeFunc(`strings`, `TrimSuffix`, TrimSuffix)
	env.AddNativeFunc(`strings`, `HasPrefix`, HasPrefix)
	env.AddNativeFunc(`strings`, `HasSuffix`, HasSuffix)
	env.AddNativeFunc(`strings`, `Contains`, Contains)
}

func Replace(ctx qnative.CallContext) {
	s := ctx.StringArg(0)
	oldPart := ctx.StringArg(1)
	newPart := ctx.StringArg(2)
	n := ctx.IntArg(3)
	ctx.SetStringResult(strings.Replace(s, oldPart, newPart, n))
}

func ReplaceAll(ctx qnative.CallContext) {
	s := ctx.StringArg(0)
	oldPart := ctx.StringArg(1)
	newPart := ctx.StringArg(2)
	ctx.SetStringResult(strings.ReplaceAll(s, oldPart, newPart))
}

func TrimPrefix(ctx qnative.CallContext) {
	s := ctx.StringArg(0)
	prefix := ctx.StringArg(1)
	ctx.SetStringResult(strings.TrimPrefix(s, prefix))
}

func TrimSuffix(ctx qnative.CallContext) {
	s := ctx.StringArg(0)
	suffix := ctx.StringArg(1)
	ctx.SetStringResult(strings.TrimSuffix(s, suffix))
}

func HasPrefix(ctx qnative.CallContext) {
	s := ctx.StringArg(0)
	prefix := ctx.StringArg(1)
	ctx.SetBoolResult(strings.HasPrefix(s, prefix))
}

func HasSuffix(ctx qnative.CallContext) {
	s := ctx.StringArg(0)
	suffix := ctx.StringArg(1)
	ctx.SetBoolResult(strings.HasSuffix(s, suffix))
}

func Contains(ctx qnative.CallContext) {
	s := ctx.StringArg(0)
	substr := ctx.StringArg(1)
	ctx.SetBoolResult(strings.Contains(s, substr))
}
