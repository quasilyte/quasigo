package main

func one(n int) int {
	println("one called")
	return n
}

func pair(n int) (int, int) {
	println("pair called")
	return n, n * 10
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
}
