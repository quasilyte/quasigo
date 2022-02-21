package main

func ten() int { return 10 }

func helloWorld() {
	println("Hello, world!")
}

func add1(x int) int {
	return x + 1
}

func strlen(s string) int {
	return len(s)
}

func concat(s1, s2 string) string {
	return s1 + s2
}

func concat3(s1, s2, s3 string) string {
	return concat(concat(s1, s2), s3)
}

//test:disasm_both
// main.streq code=52 frame=144 (6 slots: 2 args, 1 locals, 3 temps)
//   StrLen tmp1 = s1
//   StrLen tmp2 = s2
//   ScalarNotEq tmp0 = tmp1 tmp2
//   JumpFalse L0 tmp0
//   ReturnFalse
// L0:
//   LoadScalarConst i = 0
//   Jump L1
// L3:
//   StrIndex tmp1 = s1 i
//   StrIndex tmp2 = s2 i
//   ScalarNotEq tmp0 = tmp1 tmp2
//   JumpFalse L2 tmp0
//   ReturnFalse
// L2:
//   IntInc i
// L1:
//   StrLen tmp1 = s1
//   IntLt tmp0 = i tmp1
//   JumpTrue L3 tmp0
//   ReturnTrue
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
// main.fnv1 code=42 frame=120 (5 slots: 1 args, 2 locals, 2 temps)
//   LoadScalarConst v = 2166136261
//   LoadScalarConst i = 0
//   Jump L0
// L1:
//   LoadScalarConst tmp0 = 16777619
//   IntMul v = v tmp0
//   StrIndex tmp1 = s i
//   MoveScalar tmp0 = tmp1
//   IntXor v = v tmp0
//   IntInc i
// L0:
//   StrLen tmp1 = s
//   IntLt tmp0 = i tmp1
//   JumpTrue L1 tmp0
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
// main.atoi code=104 frame=240 (10 slots: 1 args, 3 locals, 6 temps)
//   StrLen tmp1 = s
//   LoadScalarConst tmp2 = 0
//   ScalarEq tmp0 = tmp1 tmp2
//   JumpFalse L0 tmp0
//   LoadScalarConst tmp0 = 0
//   ReturnScalar tmp0
// L0:
//   LoadScalarConst result = 0
//   LoadScalarConst sign = 0
//   LoadScalarConst i = 0
//   LoadScalarConst tmp2 = 0
//   StrIndex tmp1 = s tmp2
//   LoadScalarConst tmp3 = 45
//   ScalarEq tmp0 = tmp1 tmp3
//   JumpFalse L1 tmp0
//   LoadScalarConst sign = 1
//   LoadScalarConst i = 1
// L1:
//   Jump L2
// L3:
//   LoadScalarConst tmp1 = 10
//   IntMul tmp0 = result tmp1
//   StrIndex tmp4 = s i
//   LoadScalarConst tmp5 = 48
//   IntSub tmp3 = tmp4 tmp5
//   MoveScalar tmp2 = tmp3
//   IntAdd result = tmp0 tmp2
//   IntInc i
// L2:
//   StrLen tmp1 = s
//   IntLt tmp0 = i tmp1
//   JumpTrue L3 tmp0
//   JumpFalse L4 sign
//   IntNeg tmp0 = result
//   ReturnScalar tmp0
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
// main.countByte code=38 frame=144 (6 slots: 2 args, 2 locals, 2 temps)
//   LoadScalarConst n = 0
//   LoadScalarConst i = 0
//   Jump L0
// L2:
//   StrIndex tmp1 = s i
//   ScalarEq tmp0 = tmp1 b
//   JumpFalse L1 tmp0
//   IntInc n
// L1:
//   IntInc i
// L0:
//   StrLen tmp1 = s
//   IntLt tmp0 = i tmp1
//   JumpTrue L2 tmp0
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
//   StrLen tmp1 = s
//   StrLen tmp2 = prefix
//   IntGtEq tmp0 = tmp1 tmp2
//   JumpFalse L0 tmp0
//   StrLen tmp4 = prefix
//   StrSliceTo tmp3 = s tmp4
//   StrEq tmp0 = tmp3 prefix
// L0:
//   ReturnScalar tmp0
func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

//test:disasm_opt
// main.factorial code=34 frame=120 (5 slots: 1 args, 0 locals, 4 temps)
//   LoadScalarConst tmp1 = 0
//   IntLtEq tmp0 = x tmp1
//   JumpFalse L0 tmp0
//   LoadScalarConst tmp0 = 1
//   ReturnScalar tmp0
// L0:
//   LoadScalarConst tmp3 = 1
//   IntSub tmp2 = x tmp3
//   MoveScalar arg0 = tmp2
//   CallRecur tmp1
//   IntMul tmp0 = x tmp1
//   ReturnScalar tmp0
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
// main.charsum code=45 frame=192 (8 slots: 1 args, 2 locals, 5 temps)
//   LoadScalarConst sum = 0
//   LoadScalarConst i = 0
//   Jump L0
// L1:
//   LoadScalarConst tmp3 = 0
//   StrIndex tmp2 = s tmp3
//   LoadScalarConst tmp4 = 48
//   IntSub tmp1 = tmp2 tmp4
//   MoveScalar tmp0 = tmp1
//   IntAdd sum = sum tmp0
//   IntInc i
// L0:
//   StrLen tmp1 = s
//   IntLt tmp0 = i tmp1
//   JumpTrue L1 tmp0
//   ReturnScalar sum
func charsum(s string) int {
	sum := 0
	for i := 0; i < len(s); i++ {
		sum += int(s[0] - '0')
	}
	return sum
}

//test:disasm
// main.substringIndex code=54 frame=168 (7 slots: 2 args, 2 locals, 3 temps)
//   MoveStr head = s
//   LoadScalarConst i = 0
//   Jump L0
// L2:
//   StrLen tmp2 = sub
//   StrSliceTo tmp1 = head tmp2
//   StrEq tmp0 = tmp1 sub
//   JumpFalse L1 tmp0
//   ReturnScalar i
// L1:
//   IntInc i
//   LoadScalarConst tmp0 = 1
//   StrSliceFrom head = head tmp0
// L0:
//   StrLen tmp1 = head
//   StrLen tmp2 = sub
//   IntGtEq tmp0 = tmp1 tmp2
//   JumpTrue L2 tmp0
//   LoadScalarConst tmp0 = -1
//   ReturnScalar tmp0
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
	println(fnv1("examp9wqu8 rwy7ayd7s8yd S&CY&W"))
	println(fnv1("some text that will definitely cause the overflow"))
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

func main() {
	helloWorld()
	println(ten())
	println(add1(ten()))
	println(strlen("hello"))
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
}