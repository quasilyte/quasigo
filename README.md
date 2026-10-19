# quasigo

`quasigo` is a Go subset interpreter written in Go.

This is mostly a distraction project for me.

It was used in some real-world projects though:

* [go-ruleguard](https://github.com/quasilyte/go-ruleguard) - as an engine to evaluate user-defined linting rules
* [RunePact](https://store.steampowered.com/app/4614330/RunePact) - as a scripting language for the game

## Project Goals

In short: simple and efficient.

* Support only the minimal Go subset needed for game scripting and alike
* Be as fast as feasible without large code size
* Do not allocate when avoidable, keep the runtime footptint small
* Enable full data sharing with a hosting app (in/out serialization)
* Easy to integrate into Go apps
* Fast script->Go native calls

It can be used as an expression evaluation engine too, but maybe it's an overkill for that.

quasigo optimizes for low overhead per calling a small scripting func from a host app. Meaning hundreds of smaller funcs can be called inside an app in a short time. It is not optimized to be a good engine to run large scripts or even apps.
