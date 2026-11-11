package main

import "strings"

func one(n int) int {
	println("one called")
	return n
}

func pair(n int) (int, int) {
	println("pair called")
	return n, n * 10
}

// small is likely to be inlined.
func small(n int) int {
	println("small called")
	return n + 1
}

func smallBranches(n int) int {
	if n > 0 {
		return 1
	}
	return 0
}

// The inlined body should not move the result anywhere.

//test:disasm_opt
// main.inlinedDiscard code=23 frame=120 (5 slots: 1 params, 4 locals)
//   Move arg0 = x
//   CallVoid main.smallBranches()
//   Move temp0 = x
//   LoadScalarConst temp3 = 1
//   IntAdd64 temp2 = temp0 temp3
//   Move arg0 = x
//   CallVoidNative builtin.PrintInt()
//   ReturnVoid
func inlinedDiscard(x int) {
	smallBranches(x)
	plusOne(x)
	println(x)
}

func plusOne(x int) int { return x + 1 }

// A non-void recursive call in a stmt context.
func recCount(n int) int {
	if n == 0 {
		return 0
	}
	println(n)
	recCount(n - 1)
	return n
}

// A void recursive call should not clobber the frame of the caller.
func recVoid(n int) {
	if n == 0 {
		return
	}
	println(n)
	recVoid(n - 1)
	println(n)
}

func main() {
	{
		// The result is never read after the assignment,
		// but the call has a side effect and should not be removed.
		x := 1
		println(x)
		x = one(2)
	}

	{
		// Only the second result is read, but both are produced by one call.
		a, b := pair(1)
		println(a)
		println(b)
		a, b = pair(2)
		println(b)
	}

	{
		recVoid(3)
	}

	{
		// Non-void calls in a stmt context, the results are discarded.
		x := 5
		one(1)
		pair(1)
		small(1)
		smallBranches(x)
		inlinedDiscard(x)
		recCount(2)
		println(x)
		strings.HasPrefix("foo", "f")
		println(strings.HasPrefix("foo", "f"))
	}
}
