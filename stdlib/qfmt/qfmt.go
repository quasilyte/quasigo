package qfmt

import (
	"fmt"

	"github.com/quasilyte/quasigo"
)

func ImportAll(env *quasigo.Env) {
	env.AddNativeFunc(`fmt`, `Sprintf`, Sprintf)
}

func Sprintf(ctx quasigo.NativeCallContext) {
	format := ctx.StringArg(0)
	args := ctx.VariadicArg()
	ctx.SetStringResult(fmt.Sprintf(format, args...))
}
