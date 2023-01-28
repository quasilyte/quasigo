package main

//test:irdump
// block0 [0]:
//   LoadStrConst temp0 = ""
//   Move temp2.v0 = s
//   LoadScalarConst temp3.v0 = 1
//   IntSub64 temp1 = temp2.v0 temp3.v0
//   Jump L2
// block1 (L3) [0]:
//   LoadScalarConst temp4.v0 = 1
//   IntAdd64 temp3.v1 = temp1 temp4.v0
//   StrSlice temp2.v1 = s temp1 temp3.v1
//   Concat temp0 = temp0 temp2.v1
// block2 (L1) [0]:
//   IntDec temp1
// block3 (L2) [0]:
//   Zero temp3.v2
//   IntGtEq temp2.v2 = temp1 temp3.v2
//   JumpNotZero L3 temp2.v2
// block4 (L0) [1]:
//   ReturnStr temp0
//   VarKill temp0
//
//test:disasm
// main.reverse code=46 frame=144 (6 slots: 1 params, 5 locals)
//   LoadStrConst temp0 = ""
//   Move temp2 = s
//   LoadScalarConst temp3 = 1
//   IntSub64 temp1 = temp2 temp3
//   Jump L0
// L1:
//   LoadScalarConst temp4 = 1
//   IntAdd64 temp3 = temp1 temp4
//   StrSlice temp2 = s temp1 temp3
//   Concat temp0 = temp0 temp2
//   IntDec temp1
// L0:
//   Zero temp3
//   IntGtEq temp2 = temp1 temp3
//   JumpNotZero L1 temp2
//   ReturnStr temp0
//
//test:disasm_opt
// main.reverse code=43 frame=144 (6 slots: 1 params, 5 locals)
//   LoadStrConst temp0 = ""
//   LoadScalarConst temp3 = 1
//   IntSub64 temp1 = s temp3
//   Jump L0
// L1:
//   LoadScalarConst temp4 = 1
//   IntAdd64 temp3 = temp1 temp4
//   StrSlice temp2 = s temp1 temp3
//   Concat temp0 = temp0 temp2
//   IntDec temp1
// L0:
//   Zero temp3
//   IntGtEq temp2 = temp1 temp3
//   JumpNotZero L1 temp2
//   ReturnStr temp0
func reverse(s string) string {
	out := ""
	for i := len(s) - 1; i >= 0; i-- {
		out += s[i : i+1]
	}
	return out
}

func isPalindrome(s string) bool {
	return s == reverse(s)
}

func isWordChar(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9') ||
		ch == '_'
}

func isIdent(s string) bool {
	if len(s) == 0 {
		return false
	}
	first := s[0]
	if !((first >= 'a' && first <= 'z') || (first >= 'A' && first <= 'Z') || first == '_') {
		return false
	}
	for i := 1; i < len(s); i++ {
		if !isWordChar(s[i]) {
			return false
		}
	}
	return true
}

func romanDigit(s string) int {
	if s == "I" {
		return 1
	}
	if s == "V" {
		return 5
	}
	if s == "X" {
		return 10
	}
	if s == "L" {
		return 50
	}
	if s == "C" {
		return 100
	}
	if s == "D" {
		return 500
	}
	if s == "M" {
		return 1000
	}
	return -1
}

func romanToInt(s string) int {
	if s == "" {
		return 0
	}
	num := 0
	lastint := 0
	total := 0
	for i := 0; i < len(s); i++ {
		char := s[len(s)-(i+1) : len(s)-i]
		num = romanDigit(char)
		if num < lastint {
			total = total - num
		} else {
			total = total + num
		}
		lastint = num
	}
	return total
}

func testindex() {
	println("abc"[0])
	s := "hello"
	println(s[0])
	println(s[len(s)-1])
	for i := 0; i < len(s); i++ {
		println(s[i])
	}
}

func testReverse() {
	println(reverse(""))
	println(reverse("a"))
	println(reverse("hello"))
	println(reverse("Longer text"))
}

func testIsPalindrome() {
	println(isPalindrome(""))
	println(isPalindrome("eye"))
	println(isPalindrome("redivider"))
	println(isPalindrome("meow"))
	println(isPalindrome("Longer text"))
}

func testRomanToInt() {
	println(romanToInt(""))
	println(romanToInt("I"))
	println(romanToInt("III"))
	println(romanToInt("IX"))
	println(romanToInt("XXI"))
	println(romanToInt("XXII"))
	println(romanToInt("XXVI"))
	println(romanToInt("XI"))
	println(romanToInt("LVIII"))
	println(romanToInt("MCMXCIV"))
	println(romanToInt("MCMXICIVI"))
}

func testIsWordChar() {
	println(isWordChar('0'))
	println(isWordChar('5'))
	println(isWordChar('1'))
	println(isWordChar('a'))
	println(isWordChar('b'))
	println(isWordChar('z'))
	println(isWordChar('A'))
	println(isWordChar('C'))
	println(isWordChar('Z'))
	println(isWordChar('_'))
	println(isWordChar('?'))
	println(isWordChar('%'))
	println(isWordChar(10))
	println(isWordChar('\r'))
}

func testIsIdent() {
	println(isIdent(""))
	println(isIdent("213"))
	println(isIdent("%#"))
	println(isIdent("aaa%"))
	println(isIdent("Hello, World"))
	println(isIdent("_aaa"))
	println(isIdent("ident"))
	println(isIdent("ident2"))
	println(isIdent("ident_2"))
	println(isIdent("IDENT2_"))
}

//test:disasm_both
// main.testStrSlice1 code=9 frame=72 (3 slots: 1 params, 2 locals)
//   LoadScalarConst temp1 = 1
//   StrSliceFrom temp0 = s temp1
//   ReturnStr temp0
func testStrSlice1(s string) string {
	return s[1:]
}

//test:disasm_both
// main.testStrSlice2 code=9 frame=72 (3 slots: 1 params, 2 locals)
//   LoadScalarConst temp1 = 1
//   StrSliceTo temp0 = s temp1
//   ReturnStr temp0
func testStrSlice2(s string) string {
	return s[:1]
}

//test:disasm_both
// main.testStrSlice3 code=13 frame=96 (4 slots: 1 params, 3 locals)
//   LoadScalarConst temp1 = 1
//   LoadScalarConst temp2 = 2
//   StrSlice temp0 = s temp1 temp2
//   ReturnStr temp0
func testStrSlice3(s string) string {
	return s[1:2]
}

func testBasicOps() {
	s1 := "hello"
	s2 := "world"
	println(testStrSlice1(s1))
	println(testStrSlice2(s1))
	println(testStrSlice3(s1))
	println(s1 < s2)
	println(s1 > s2)
	println(s2 < s1)
	println(s2 > s1)
	println(s1 < "")
	println(s1 > "")
	println(s2 < "")
	println(s2 > "")
	println(s1 < "OK")
	println(s1 > "OK")
	println(s2 < "OK")
	println(s2 > "OK")
}

func main() {
	testBasicOps()
	testindex()
	testReverse()
	testIsPalindrome()
	testRomanToInt()
	testIsWordChar()
	testIsIdent()
}
