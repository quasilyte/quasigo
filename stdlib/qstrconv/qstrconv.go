package qstrconv

import (
	"strconv"

	"github.com/quasilyte/quasigo"
)

func ImportAll(env *quasigo.Env) {
	env.AddNativeFunc(`strconv`, `Atoi`, Atoi)
	env.AddNativeFunc(`strconv`, `Itoa`, Itoa)
}

func Atoi(ctx quasigo.NativeCallContext) {
	s := ctx.StringArg(0)
	v, err := strconv.Atoi(s)
	ctx.SetIntResult(v)
	ctx.SetInterfaceResult2(err)
}

func Itoa(ctx quasigo.NativeCallContext) {
	ctx.SetStringResult(strconv.Itoa(ctx.IntArg(0)))
}
