package main

func printBytes(b []byte) {
	// TODO: use range stmt.
	for i := 0; i < len(b); i++ {
		println(b[i])
	}
}

//test:disasm
// main.indexByteAsInt code=11 frame=96 (4 slots: 1 params, 3 locals)
//   Zero temp2
//   SliceIndexScalar8 temp1 = b temp2
//   Move temp0 = temp1
//   ReturnScalar temp0
//
//test:disasm_opt
// main.indexByteAsInt code=8 frame=72 (3 slots: 1 params, 2 locals)
//   Zero temp1
//   SliceIndexScalar8 temp0 = b temp1
//   ReturnScalar temp0
func indexByteAsInt(b []byte) int {
	return int(b[0])
}

//test:disasm_both
// main.slicingBytesFromTo code=7 frame=96 (4 slots: 3 params, 1 locals)
//   BytesSlice temp0 = b from to
//   Return temp0
func slicingBytesFromTo(b []byte, from, to int) []byte {
	return b[from:to]
}

//test:disasm_both
// main.slicingBytesFrom code=6 frame=72 (3 slots: 2 params, 1 locals)
//   BytesSliceFrom temp0 = b from
//   Return temp0
func slicingBytesFrom(b []byte, from int) []byte {
	return b[from:]
}

//test:disasm_both
// main.slicingBytesTo code=6 frame=72 (3 slots: 2 params, 1 locals)
//   BytesSliceTo temp0 = b to
//   Return temp0
func slicingBytesTo(b []byte, to int) []byte {
	return b[:to]
}

//test:disasm_opt
// main.stringReverse code=62 frame=144 (6 slots: 1 params, 5 locals)
//   Move temp2 = s
//   Move temp3 = s
//   LoadScalarConst arg0 = 1
//   Move arg1 = temp2
//   Move arg2 = temp3
//   CallNative temp0 = builtin.makeSlice()
//   Zero temp1
//   LoadScalarConst temp4 = 1
//   IntSub64 temp2 = s temp4
//   Jump L0
// L1:
//   StrIndex temp3 = s temp2
//   SliceSetScalar8 temp0 temp1 temp3
//   IntInc temp1
//   IntDec temp2
// L0:
//   Zero temp4
//   IntGtEq temp3 = temp2 temp4
//   JumpNotZero L1 temp3
//   Move arg0 = temp0
//   CallNative temp2 = builtin.bytesToString()
//   ReturnStr temp2
//
//test:disasm
// main.stringReverse code=68 frame=144 (6 slots: 1 params, 5 locals)
//   LoadScalarConst temp1 = 1
//   Move temp2 = s
//   Move temp3 = s
//   Move arg0 = temp1
//   Move arg1 = temp2
//   Move arg2 = temp3
//   CallNative temp0 = builtin.makeSlice()
//   Zero temp1
//   Move temp3 = s
//   LoadScalarConst temp4 = 1
//   IntSub64 temp2 = temp3 temp4
//   Jump L0
// L1:
//   StrIndex temp3 = s temp2
//   SliceSetScalar8 temp0 temp1 temp3
//   IntInc temp1
//   IntDec temp2
// L0:
//   Zero temp4
//   IntGtEq temp3 = temp2 temp4
//   JumpNotZero L1 temp3
//   Move arg0 = temp0
//   CallNative temp2 = builtin.bytesToString()
//   ReturnStr temp2
func stringReverse(s string) string {
	out := make([]byte, len(s))
	j := 0
	for i := len(s) - 1; i >= 0; i-- {
		out[j] = s[i]
		j++
	}
	return string(out)
}

//test:disasm_opt
// main.removeChar code=65 frame=144 (6 slots: 2 params, 4 locals)
//   Move temp3 = s
//   LoadScalarConst arg0 = 1
//   Zero arg1
//   Move arg2 = temp3
//   CallNative temp0 = builtin.makeSlice()
//   Zero temp1
//   Jump L0
// L3:
//   StrIndex temp3 = s temp1
//   ScalarEq temp2 = temp3 ch
//   JumpZero L1 temp2
//   Jump L2
// L1:
//   Move arg0 = temp0
//   StrIndex arg1 = s temp1
//   CallNative temp0 = builtin.append8()
// L2:
//   IntInc temp1
// L0:
//   IntLt temp2 = temp1 s
//   JumpNotZero L3 temp2
//   Move arg0 = temp0
//   CallNative temp1 = builtin.bytesToString()
//   ReturnStr temp1
//
//test:disasm
// main.removeChar code=74 frame=144 (6 slots: 2 params, 4 locals)
//   LoadScalarConst temp1 = 1
//   Zero temp2
//   Move temp3 = s
//   Move arg0 = temp1
//   Move arg1 = temp2
//   Move arg2 = temp3
//   CallNative temp0 = builtin.makeSlice()
//   Zero temp1
//   Jump L0
// L3:
//   StrIndex temp3 = s temp1
//   ScalarEq temp2 = temp3 ch
//   JumpZero L1 temp2
//   Jump L2
// L1:
//   Move arg0 = temp0
//   StrIndex arg1 = s temp1
//   CallNative temp0 = builtin.append8()
// L2:
//   IntInc temp1
// L0:
//   Move temp3 = s
//   IntLt temp2 = temp1 temp3
//   JumpNotZero L3 temp2
//   Move arg0 = temp0
//   CallNative temp1 = builtin.bytesToString()
//   ReturnStr temp1
func removeChar(s string, ch byte) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == ch {
			continue
		}
		out = append(out, s[i])
	}
	return string(out)
}

//test:disasm_opt
// main.tolower code=80 frame=144 (6 slots: 1 params, 5 locals)
//   Move temp2 = s
//   Move temp3 = s
//   LoadScalarConst arg0 = 1
//   Move arg1 = temp2
//   Move arg2 = temp3
//   CallNative temp0 = builtin.makeSlice()
//   Zero temp1
//   Jump L0
// L3:
//   StrIndex temp2 = s temp1
//   LoadScalarConst temp4 = 65
//   IntGtEq temp3 = temp2 temp4
//   JumpZero L1 temp3
//   LoadScalarConst temp4 = 90
//   IntLtEq temp3 = temp2 temp4
// L1:
//   JumpZero L2 temp3
//   LoadScalarConst temp3 = 32
//   IntAdd8 temp2 = temp2 temp3
// L2:
//   SliceSetScalar8 temp0 temp1 temp2
//   IntInc temp1
// L0:
//   IntLt temp2 = temp1 s
//   JumpNotZero L3 temp2
//   Move arg0 = temp0
//   CallNative temp1 = builtin.bytesToString()
//   ReturnStr temp1
//
//test:disasm
// main.tolower code=86 frame=144 (6 slots: 1 params, 5 locals)
//   LoadScalarConst temp1 = 1
//   Move temp2 = s
//   Move temp3 = s
//   Move arg0 = temp1
//   Move arg1 = temp2
//   Move arg2 = temp3
//   CallNative temp0 = builtin.makeSlice()
//   Zero temp1
//   Jump L0
// L3:
//   StrIndex temp2 = s temp1
//   LoadScalarConst temp4 = 65
//   IntGtEq temp3 = temp2 temp4
//   JumpZero L1 temp3
//   LoadScalarConst temp4 = 90
//   IntLtEq temp3 = temp2 temp4
// L1:
//   JumpZero L2 temp3
//   LoadScalarConst temp3 = 32
//   IntAdd8 temp2 = temp2 temp3
// L2:
//   SliceSetScalar8 temp0 temp1 temp2
//   IntInc temp1
// L0:
//   Move temp3 = s
//   IntLt temp2 = temp1 temp3
//   JumpNotZero L3 temp2
//   Move arg0 = temp0
//   CallNative temp1 = builtin.bytesToString()
//   ReturnStr temp1
func tolower(s string) string {
	out := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch >= 'A' && ch <= 'Z' {
			ch += 32
		}
		out[i] = ch
	}
	return string(out)
}

func testStringReverse() {
	println(stringReverse(""))
	println(stringReverse("123"))
	println(stringReverse("hello, world!"))
}

func testRemoveChar() {
	println(removeChar("", 'x'))
	println(removeChar("000", '0'))
	println(removeChar("123", '1'))
	println(removeChar("hello, world!", 'l'))
}

func testToLower() {
	println(tolower(""))
	println(tolower("abc"))
	println(tolower("SDhusdGYASGdsdc cx"))
	println(tolower("Hello, world!"))
	println(tolower("HELLO, WORLD!"))
}

func testSlicing() {
	b := make([]byte, 5)
	for i := 0; i < len(b); i++ {
		b[i] = byte(i+10) * 2
	}
	printBytes(b)
	printBytes(b[:])
	printBytes(b[0:])
	printBytes(b[:len(b)])
	printBytes(b[:len(b)-1])
	printBytes(b[1:])
	printBytes(b[:1])
	printBytes(slicingBytesFromTo(b, 1, 3))
	printBytes(slicingBytesFrom(b, 1))
	printBytes(slicingBytesTo(b, 3))
}

func testIndexing() {
	b := make([]byte, 3)
	b[0] = 10
	b[1] = 20
	b[2] = 30
	println(indexByteAsInt(b))
}

func main() {
	testIndexing()
	testStringReverse()
	testRemoveChar()
	testToLower()
	testSlicing()
}
