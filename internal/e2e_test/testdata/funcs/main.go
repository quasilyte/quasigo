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

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

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
}
