package main

// Funcs are intentionally defined in the "wrong" order:
// callers go before callees.

func main() {
	println(sumTo(10))
	println(isEven(10))
	println(isEven(7))
	println(isOdd(10))
	println(isOdd(7))
	println(pair())
	println(greet("world"))
	println(chain(1))
}

func sumTo(n int) int {
	if n == 0 {
		return 0
	}
	return n + sumTo(n-1)
}

// Mutual recursion.
func isEven(n int) bool {
	if n == 0 {
		return true
	}
	return isOdd(n - 1)
}

func isOdd(n int) bool {
	if n == 0 {
		return false
	}
	return isEven(n - 1)
}

// Forward call with a multi-value result and a void forward call.
func pair() int {
	a, b := divmod(17, 5)
	log("pair")
	return a*100 + b
}

func divmod(x, y int) (int, int) {
	return x / y, x % y
}

func log(s string) {
	println(s)
}

func greet(name string) string {
	return prefix() + name
}

func prefix() string { return "hello, " }

// A chain where everything is a forward call.
//test:irdump
// block0 [0]:
//   Move arg0 = x
//   Call temp1.v0 = main.chain2
//   LoadScalarConst temp2.v0 = 1
//   IntAdd64 temp0.v0 = temp1.v0 temp2.v0
//   ReturnScalar temp0.v0
func chain(x int) int { return chain2(x) + 1 }
func chain2(x int) int { return chain3(x) + 1 }
func chain3(x int) int { return x + 1 }
