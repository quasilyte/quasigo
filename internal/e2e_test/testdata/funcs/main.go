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
// main.streq code=51 frame=144 (6 slots: 2 params, 4 locals)
//   Len temp1 = s1
//   Len temp2 = s2
//   ScalarNotEq temp0 = temp1 temp2
//   JumpZero L0 temp0
//   ReturnZero
// L0:
//   Zero temp0
//   Jump L1
// L3:
//   StrIndex temp2 = s1 temp0
//   StrIndex temp3 = s2 temp0
//   ScalarNotEq temp1 = temp2 temp3
//   JumpZero L2 temp1
//   ReturnZero
// L2:
//   IntInc temp0
// L1:
//   Len temp2 = s1
//   IntLt temp1 = temp0 temp2
//   JumpNotZero L3 temp1
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

//test:disasm
// main.fnv1 code=41 frame=120 (5 slots: 1 params, 4 locals)
//   LoadScalarConst temp0 = 2166136261
//   Zero temp1
//   Jump L0
// L1:
//   LoadScalarConst temp2 = 16777619
//   IntMul64 temp0 = temp0 temp2
//   StrIndex temp3 = s temp1
//   Move temp2 = temp3
//   IntXor temp0 = temp0 temp2
//   IntInc temp1
// L0:
//   Len temp3 = s
//   IntLt temp2 = temp1 temp3
//   JumpNotZero L1 temp2
//   ReturnScalar temp0
//
//test:disasm_opt
// main.fnv1 code=38 frame=120 (5 slots: 1 params, 4 locals)
//   LoadScalarConst temp0 = 2166136261
//   Zero temp1
//   Jump L0
// L1:
//   LoadScalarConst temp2 = 16777619
//   IntMul64 temp0 = temp0 temp2
//   StrIndex temp3 = s temp1
//   IntXor temp0 = temp0 temp3
//   IntInc temp1
// L0:
//   Len temp3 = s
//   IntLt temp2 = temp1 temp3
//   JumpNotZero L1 temp2
//   ReturnScalar temp0
func fnv1(s string) int {
	v := 0x811c9dc5
	for i := 0; i < len(s); i++ {
		v *= 0x01000193
		v ^= int(s[i])
	}
	return v
}

//test:disasm_both
// main.isNumericString code=50 frame=168 (7 slots: 1 params, 6 locals)
//   Zero temp0
//   Jump L0
// L3:
//   StrIndex temp2 = s temp0
//   LoadScalarConst temp3 = 48
//   IntLt temp1 = temp2 temp3
//   JumpNotZero L1 temp1
//   StrIndex temp4 = s temp0
//   LoadScalarConst temp5 = 57
//   IntGt temp1 = temp4 temp5
// L1:
//   JumpZero L2 temp1
//   ReturnZero
// L2:
//   IntInc temp0
// L0:
//   Len temp2 = s
//   IntLt temp1 = temp0 temp2
//   JumpNotZero L3 temp1
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
// main.atoi code=95 frame=240 (10 slots: 1 params, 9 locals)
//   Len temp1 = s
//   Zero temp2
//   ScalarEq temp0 = temp1 temp2
//   JumpZero L0 temp0
//   ReturnZero
// L0:
//   Zero temp0
//   Zero temp1
//   Zero temp2
//   Zero temp5
//   StrIndex temp4 = s temp5
//   LoadScalarConst temp6 = 45
//   ScalarEq temp3 = temp4 temp6
//   JumpZero L1 temp3
//   LoadScalarConst temp1 = 1
//   LoadScalarConst temp2 = 1
// L1:
//   Jump L2
// L3:
//   LoadScalarConst temp4 = 10
//   IntMul64 temp3 = temp0 temp4
//   StrIndex temp7 = s temp2
//   LoadScalarConst temp8 = 48
//   IntSub8 temp6 = temp7 temp8
//   Move temp5 = temp6
//   IntAdd64 temp0 = temp3 temp5
//   IntInc temp2
// L2:
//   Len temp4 = s
//   IntLt temp3 = temp2 temp4
//   JumpNotZero L3 temp3
//   JumpZero L4 temp1
//   IntNeg temp3 = temp0
//   ReturnScalar temp3
// L4:
//   ReturnScalar temp0
//
//test:disasm_opt
// main.atoi code=86 frame=240 (10 slots: 1 params, 9 locals)
//   Len temp1 = s
//   JumpNotZero L0 temp1
//   ReturnZero
// L0:
//   Zero temp0
//   Zero temp1
//   Zero temp2
//   Zero temp5
//   StrIndex temp4 = s temp5
//   LoadScalarConst temp6 = 45
//   ScalarEq temp3 = temp4 temp6
//   JumpZero L1 temp3
//   LoadScalarConst temp1 = 1
//   LoadScalarConst temp2 = 1
// L1:
//   Jump L2
// L3:
//   LoadScalarConst temp4 = 10
//   IntMul64 temp3 = temp0 temp4
//   StrIndex temp7 = s temp2
//   LoadScalarConst temp8 = 48
//   IntSub8 temp6 = temp7 temp8
//   IntAdd64 temp0 = temp3 temp6
//   IntInc temp2
// L2:
//   Len temp4 = s
//   IntLt temp3 = temp2 temp4
//   JumpNotZero L3 temp3
//   JumpZero L4 temp1
//   IntNeg temp3 = temp0
//   ReturnScalar temp3
// L4:
//   ReturnScalar temp0
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
// main.countByte code=36 frame=144 (6 slots: 2 params, 4 locals)
//   Zero temp0
//   Zero temp1
//   Jump L0
// L2:
//   StrIndex temp3 = s temp1
//   ScalarEq temp2 = temp3 b
//   JumpZero L1 temp2
//   IntInc temp0
// L1:
//   IntInc temp1
// L0:
//   Len temp3 = s
//   IntLt temp2 = temp1 temp3
//   JumpNotZero L2 temp2
//   ReturnScalar temp0
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
// main.hasPrefix code=27 frame=168 (7 slots: 2 params, 5 locals)
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
// main.factorial code=29 frame=120 (5 slots: 1 params, 4 locals)
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

//test:disasm
// main.charsum code=42 frame=192 (8 slots: 1 params, 7 locals)
//   Zero temp0
//   Zero temp1
//   Jump L0
// L1:
//   Zero temp5
//   StrIndex temp4 = s temp5
//   LoadScalarConst temp6 = 48
//   IntSub8 temp3 = temp4 temp6
//   Move temp2 = temp3
//   IntAdd64 temp0 = temp0 temp2
//   IntInc temp1
// L0:
//   Len temp3 = s
//   IntLt temp2 = temp1 temp3
//   JumpNotZero L1 temp2
//   ReturnScalar temp0
//
//test:disasm_opt
// main.charsum code=39 frame=192 (8 slots: 1 params, 7 locals)
//   Zero temp0
//   Zero temp1
//   Jump L0
// L1:
//   Zero temp5
//   StrIndex temp4 = s temp5
//   LoadScalarConst temp6 = 48
//   IntSub8 temp3 = temp4 temp6
//   IntAdd64 temp0 = temp0 temp3
//   IntInc temp1
// L0:
//   Len temp3 = s
//   IntLt temp2 = temp1 temp3
//   JumpNotZero L1 temp2
//   ReturnScalar temp0
func charsum(s string) int {
	sum := 0
	for i := 0; i < len(s); i++ {
		sum += int(s[0] - '0')
	}
	return sum
}

//test:disasm
// main.substringIndex code=53 frame=168 (7 slots: 2 params, 5 locals)
//   Move temp0 = s
//   Zero temp1
//   Jump L0
// L2:
//   Len temp4 = sub
//   StrSliceTo temp3 = temp0 temp4
//   StrEq temp2 = temp3 sub
//   JumpZero L1 temp2
//   ReturnScalar temp1
// L1:
//   IntInc temp1
//   LoadScalarConst temp2 = 1
//   StrSliceFrom temp0 = temp0 temp2
// L0:
//   Len temp3 = temp0
//   Len temp4 = sub
//   IntGtEq temp2 = temp3 temp4
//   JumpNotZero L2 temp2
//   LoadScalarConst temp2 = -1
//   ReturnScalar temp2
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
