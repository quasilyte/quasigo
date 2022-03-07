package main

//test:disasm_both
// main.reverse code=47 frame=144 (6 slots: 1 args, 2 locals, 3 temps)
//   LoadStrConst out = ""
//   Len tmp0 = s
//   LoadScalarConst tmp1 = 1
//   IntSub i = tmp0 tmp1
//   Jump L0
// L1:
//   LoadScalarConst tmp2 = 1
//   IntAdd tmp1 = i tmp2
//   StrSlice tmp0 = s i tmp1
//   Concat out = out tmp0
//   IntDec i
// L0:
//   LoadScalarConst tmp1 = 0
//   IntGtEq tmp0 = i tmp1
//   JumpNotZero L1 tmp0
//   ReturnStr out
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

func main() {
	testindex()
	testReverse()
	testIsPalindrome()
	testRomanToInt()
	testIsWordChar()
	testIsIdent()
}
