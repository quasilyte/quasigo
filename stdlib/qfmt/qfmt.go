package qfmt

import (
	"fmt"

	"github.com/quasilyte/quasigo"
	"github.com/quasilyte/quasigo/qnative"
)

func ImportAll(env *quasigo.Env) {
	env.AddNativeFunc(`fmt`, `Sprintf`, Sprintf)
}

func Sprintf(ctx qnative.CallContext) {
	format := ctx.StringArg(0)
	args := ctx.VariadicArg()
	ctx.SetStringResult(fmt.Sprintf(format, args...))
}
