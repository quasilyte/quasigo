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

//test:irdump
// block0 [0]:
//   Zero temp1.v0
//   ScalarEq temp0.v0 = n temp1.v0
//   JumpZero L0 temp0.v0
// block1 [0]:
//   ReturnOne
// block2 (L0) [0]:
//   LoadScalarConst temp2.v0 = 1
//   IntSub64 temp1.v1 = n temp2.v0
//   Move arg0 = temp1.v1
//   Call temp0.v1 = main.isOdd
//   ReturnScalar temp0.v1
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

//test:disasm
// main.chain code=16 frame=96 (4 slots: 1 params, 3 locals)
//   Move arg0 = x
//   Call temp1 = main.chain2()
//   LoadScalarConst temp2 = 1
//   IntAdd64 temp0 = temp1 temp2
//   ReturnScalar temp0
//test:disasm_opt
// main.chain code=35 frame=264 (11 slots: 1 params, 10 locals)
//   Move temp2 = x
//   Move temp6 = temp2
//   LoadScalarConst temp9 = 1
//   IntAdd64 temp8 = temp6 temp9
//   Move temp5 = temp8
//   LoadScalarConst temp6 = 1
//   IntAdd64 temp4 = temp5 temp6
//   Move temp1 = temp4
//   LoadScalarConst temp2 = 1
//   IntAdd64 temp0 = temp1 temp2
//   ReturnScalar temp0
func chain(x int) int { return chain2(x) + 1 }
func chain2(x int) int { return chain3(x) + 1 }
func chain3(x int) int { return x + 1 }
