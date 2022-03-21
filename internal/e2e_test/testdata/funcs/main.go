package main

func ten() int { return 10 }

func helloWorld() {
	println("Hello, world!")
}

func add1(x int) int {
	return x + 1
}

func Len(s string) int {
	return len(s)
}

func concat(s1, s2 string) string {
	return s1 + s2
}

func concat3(s1, s2, s3 string) string {
	return concat(concat(s1, s2), s3)
}

//test:disasm_both
// main.streq code=51 frame=144 (6 slots: 2 args, 1 locals, 3 temps)
//   Len temp1 = s1
//   Len temp2 = s2
//   ScalarNotEq temp0 = temp1 temp2
//   JumpZero L0 temp0
//   ReturnZero
// L0:
//   Zero i
//   Jump L1
// L3:
//   StrIndex temp1 = s1 i
//   StrIndex temp2 = s2 i
//   ScalarNotEq temp0 = temp1 temp2
//   JumpZero L2 temp0
//   ReturnZero
// L2:
//   IntInc i
// L1:
//   Len temp1 = s1
//   IntLt temp0 = i temp1
//   JumpNotZero L3 temp0
//   ReturnOne
func streq(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

//test:disasm_both
// main.fnv1 code=41 frame=120 (5 slots: 1 args, 2 locals, 2 temps)
//   LoadScalarConst v = 2166136261
//   Zero i
//   Jump L0
// L1:
//   LoadScalarConst temp0 = 16777619
//   IntMul64 v = v temp0
//   StrIndex temp1 = s i
//   Move temp0 = temp1
//   IntXor v = v temp0
//   IntInc i
// L0:
//   Len temp1 = s
//   IntLt temp0 = i temp1
//   JumpNotZero L1 temp0
//   ReturnScalar v
func fnv1(s string) int {
	v := 0x811c9dc5
	for i := 0; i < len(s); i++ {
		v *= 0x01000193
		v ^= int(s[i])
	}
	return v
}

//test:disasm_both
// main.isNumericString code=50 frame=168 (7 slots: 1 args, 1 locals, 5 temps)
//   Zero i
//   Jump L0
// L3:
//   StrIndex temp1 = s i
//   LoadScalarConst temp2 = 48
//   IntLt temp0 = temp1 temp2
//   JumpNotZero L1 temp0
//   StrIndex temp3 = s i
//   LoadScalarConst temp4 = 57
//   IntGt temp0 = temp3 temp4
// L1:
//   JumpZero L2 temp0
//   ReturnZero
// L2:
//   IntInc i
// L0:
//   Len temp1 = s
//   IntLt temp0 = i temp1
//   JumpNotZero L3 temp0
//   ReturnOne
func isNumericString(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

//test:disasm
// main.atoi code=95 frame=240 (10 slots: 1 args, 3 locals, 6 temps)
//   Len temp1 = s
//   Zero temp2
//   ScalarEq temp0 = temp1 temp2
//   JumpZero L0 temp0
//   ReturnZero
// L0:
//   Zero result
//   Zero sign
//   Zero i
//   Zero temp2
//   StrIndex temp1 = s temp2
//   LoadScalarConst temp3 = 45
//   ScalarEq temp0 = temp1 temp3
//   JumpZero L1 temp0
//   LoadScalarConst sign = 1
//   LoadScalarConst i = 1
// L1:
//   Jump L2
// L3:
//   LoadScalarConst temp1 = 10
//   IntMul64 temp0 = result temp1
//   StrIndex temp4 = s i
//   LoadScalarConst temp5 = 48
//   IntSub8 temp3 = temp4 temp5
//   Move temp2 = temp3
//   IntAdd64 result = temp0 temp2
//   IntInc i
// L2:
//   Len temp1 = s
//   IntLt temp0 = i temp1
//   JumpNotZero L3 temp0
//   JumpZero L4 sign
//   IntNeg temp0 = result
//   ReturnScalar temp0
// L4:
//   ReturnScalar result
//
//test:disasm_opt
// main.atoi code=89 frame=240 (10 slots: 1 args, 3 locals, 6 temps)
//   Len temp1 = s
//   JumpNotZero L0 temp1
//   ReturnZero
// L0:
//   Zero result
//   Zero sign
//   Zero i
//   Zero temp2
//   StrIndex temp1 = s temp2
//   LoadScalarConst temp3 = 45
//   ScalarEq temp0 = temp1 temp3
//   JumpZero L1 temp0
//   LoadScalarConst sign = 1
//   LoadScalarConst i = 1
// L1:
//   Jump L2
// L3:
//   LoadScalarConst temp1 = 10
//   IntMul64 temp0 = result temp1
//   StrIndex temp4 = s i
//   LoadScalarConst temp5 = 48
//   IntSub8 temp3 = temp4 temp5
//   Move temp2 = temp3
//   IntAdd64 result = temp0 temp2
//   IntInc i
// L2:
//   Len temp1 = s
//   IntLt temp0 = i temp1
//   JumpNotZero L3 temp0
//   JumpZero L4 sign
//   IntNeg temp0 = result
//   ReturnScalar temp0
// L4:
//   ReturnScalar result
func atoi(s string) int {
	if len(s) == 0 {
		return 0
	}
	result := 0
	sign := false
	i := 0
	if s[0] == '-' {
		sign = true
		i = 1
	}
	for i < len(s) {
		result = result*10 + int(s[i]-'0')
		i++
	}
	if sign {
		return -result
	}
	return result
}

//test:disasm_both
// main.countByte code=36 frame=144 (6 slots: 2 args, 2 locals, 2 temps)
//   Zero n
//   Zero i
//   Jump L0
// L2:
//   StrIndex temp1 = s i
//   ScalarEq temp0 = temp1 b
//   JumpZero L1 temp0
//   IntInc n
// L1:
//   IntInc i
// L0:
//   Len temp1 = s
//   IntLt temp0 = i temp1
//   JumpNotZero L2 temp0
//   ReturnScalar n
func countByte(s string, b byte) int {
	n := 0
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			n++
		}
	}
	return n
}

//test:disasm_opt
// main.hasPrefix code=27 frame=168 (7 slots: 2 args, 0 locals, 5 temps)
//   Len temp1 = s
//   Len temp2 = prefix
//   IntGtEq temp0 = temp1 temp2
//   JumpZero L0 temp0
//   Len temp4 = prefix
//   StrSliceTo temp3 = s temp4
//   StrEq temp0 = temp3 prefix
// L0:
//   ReturnScalar temp0
func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

//test:disasm_opt
// main.factorial code=29 frame=120 (5 slots: 1 args, 0 locals, 4 temps)
//   Zero temp1
//   IntLtEq temp0 = x temp1
//   JumpZero L0 temp0
//   ReturnOne
// L0:
//   LoadScalarConst temp3 = 1
//   IntSub64 temp2 = x temp3
//   Move arg0 = temp2
//   CallRecur temp1
//   IntMul64 temp0 = x temp1
//   ReturnScalar temp0
func factorial(x int) int {
	if x <= 0 {
		return 1
	}
	return x * factorial(x-1)
}

func testFactorial() {
	i := 0
	for i < 10 {
		println(factorial(i))
		i++
	}
	println(factorial(4) + factorial(8))
}

//test:disasm_both
// main.charsum code=42 frame=192 (8 slots: 1 args, 2 locals, 5 temps)
//   Zero sum
//   Zero i
//   Jump L0
// L1:
//   Zero temp3
//   StrIndex temp2 = s temp3
//   LoadScalarConst temp4 = 48
//   IntSub8 temp1 = temp2 temp4
//   Move temp0 = temp1
//   IntAdd64 sum = sum temp0
//   IntInc i
// L0:
//   Len temp1 = s
//   IntLt temp0 = i temp1
//   JumpNotZero L1 temp0
//   ReturnScalar sum
func charsum(s string) int {
	sum := 0
	for i := 0; i < len(s); i++ {
		sum += int(s[0] - '0')
	}
	return sum
}

//test:disasm
// main.substringIndex code=53 frame=168 (7 slots: 2 args, 2 locals, 3 temps)
//   Move head = s
//   Zero i
//   Jump L0
// L2:
//   Len temp2 = sub
//   StrSliceTo temp1 = head temp2
//   StrEq temp0 = temp1 sub
//   JumpZero L1 temp0
//   ReturnScalar i
// L1:
//   IntInc i
//   LoadScalarConst temp0 = 1
//   StrSliceFrom head = head temp0
// L0:
//   Len temp1 = head
//   Len temp2 = sub
//   IntGtEq temp0 = temp1 temp2
//   JumpNotZero L2 temp0
//   LoadScalarConst temp0 = -1
//   ReturnScalar temp0
func substringIndex(s, sub string) int {
	head := s
	i := 0
	for len(head) >= len(sub) {
		if head[:len(sub)] == sub {
			return i
		}
		i++
		head = head[1:]
	}
	return -1
}

func testSubstringIndex() {
	println(substringIndex("", ""))
	println(substringIndex("hello", "h"))
	println(substringIndex("h", "hello"))
	println(substringIndex("hello, world", "world"))
	println(substringIndex("hello, world", "hello"))
	println(substringIndex("hello, world", ","))
	println(substringIndex("abc", "a"))
	println(substringIndex("abc", "b"))
	println(substringIndex("abc", "c"))
	println(substringIndex("a", "abc"))
	println(substringIndex("b", "abc"))
	println(substringIndex("c", "abc"))
}

func testAtoi() {
	println(atoi(""))
	println(atoi("1"))
	println(atoi("255"))
	println(atoi("-1"))
	println(atoi("-127"))
	println(atoi("3438"))
	println(atoi("139"))
	println(atoi("-19224"))
	println(atoi("9380000"))
	println(atoi("-93100110"))
}

func testStreq() {
	println(streq("", ""))
	println(streq("1", ""))
	println(streq("", "1"))
	println(streq("abc", "abc"))
	println(streq("hello", "holla"))
	println(streq("123", "124"))
}

func testFnv1() {
	println(fnv1(""))
	println(fnv1("0"))
	println(fnv1("x"))
	println(fnv1("hello"))
	println(fnv1("2834"))
	println(fnv1("dsiua9uqw"))
	println(fnv1("Hello, world!"))
	println(fnv1("aaaaaaaaaaaaaaaa"))
	println(fnv1("aaaaaaaaaaaaaaaaaaaa"))
	println(fnv1("examp9wqu8 rwy7ayd7s8yd S&CY&W"))
	println(fnv1("some text that will definitely cause the overflow"))
	println(fnv1("Lorem ipsum dolor sit amet, consectetur adipiscing elit"))
	println(fnv1("1, 2, Fizz, 4, Buzz, Fizz, 7, 8, Fizz, Buzz, 11, Fizz, 13, 14, Fizz Buzz, 16, 17, Fizz, 19, Buzz, Fizz, 22, 23, Fizz, Buzz, 26, Fizz, 28, 29, Fizz Buzz"))
}

func testCountByte() {
	println(countByte("foo", '0'))
	println(countByte("foo", 'f'))
	println(countByte("foo", 'o'))
	println(countByte("foo", 0))
	println(countByte("Hello, world", 'o'))
	println(countByte("Hello, world", 'z'))
	println(countByte("Hello, world", ' '))
	println(countByte("Hello, world", 'l'))
}

func testCharsum() {
	println(charsum(""))
	println(charsum("0"))
	println(charsum("foo"))
	println(charsum("hello, world"))
	println(charsum("some longer string for the test purposes"))
	println(charsum("329i4i24923"))
	println(charsum("1-010329 8*$#&Q YW&FWQ&Dsahdsyds "))
}

func testNumericString() {
	println(isNumericString(""))
	println(isNumericString("0"))
	println(isNumericString("1392"))
	println(isNumericString("13922183"))
	println(isNumericString("a"))
	println(isNumericString("xasid9"))
	println(isNumericString("28382x"))
}

func main() {
	helloWorld()
	println(ten())
	println(add1(ten()))
	println(Len("hello"))
	println(concat("foo", "bar"))
	println(concat3("", "", ""))
	println(concat3("x", "", ""))
	println(concat3("", "x", ""))
	println(concat3("", "", "x"))
	println(concat3("a", "b", "c"))
	println(concat3("hello", "world", ""))
	println(hasPrefix("", ""))
	println(hasPrefix("", "hello"))
	println(hasPrefix("hello", ""))
	println(hasPrefix("hello", "hello"))
	println(hasPrefix("hello, world", "hello"))
	testFactorial()
	testSubstringIndex()
	testAtoi()
	testStreq()
	testFnv1()
	testCountByte()
	testCharsum()
	testNumericString()
}
