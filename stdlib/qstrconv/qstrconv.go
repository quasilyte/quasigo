package qstrconv

import (
	"strconv"

	"github.com/quasilyte/quasigo"
	"github.com/quasilyte/quasigo/qnative"
)

func ImportAll(env *quasigo.Env) {
	env.AddNativeFunc(`strconv`, `Atoi`, Atoi)
	env.AddNativeFunc(`strconv`, `Itoa`, Itoa)
}

func Atoi(ctx qnative.CallContext) {
	s := ctx.StringArg(0)
	v, err := strconv.Atoi(s)
	ctx.SetIntResult(v)
	ctx.SetInterfaceResult2(err)
}

func Itoa(ctx qnative.CallContext) {
	ctx.SetStringResult(strconv.Itoa(ctx.IntArg(0)))
}
