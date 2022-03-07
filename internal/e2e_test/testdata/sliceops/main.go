package main

//test:disasm_both
// main.makeIntSlice1 code=15 frame=48 (2 slots: 1 args, 0 locals, 1 temps)
//   LoadScalarConst arg0 = 8
//   Move arg1 = length
//   Move arg2 = length
//   CallNative tmp0 = builtin.makeSlice()
//   Return tmp0
func makeIntSlice1(length int) []int {
	return make([]int, length)
}

//test:disasm_both
// main.makeByteSlice1 code=15 frame=48 (2 slots: 1 args, 0 locals, 1 temps)
//   LoadScalarConst arg0 = 1
//   Move arg1 = length
//   Move arg2 = length
//   CallNative tmp0 = builtin.makeSlice()
//   Return tmp0
func makeByteSlice1(length int) []byte {
	return make([]byte, length)
}

//test:disasm_both
// main.makeBoolSlice1 code=15 frame=48 (2 slots: 1 args, 0 locals, 1 temps)
//   LoadScalarConst arg0 = 1
//   Move arg1 = length
//   Move arg2 = length
//   CallNative tmp0 = builtin.makeSlice()
//   Return tmp0
func makeBoolSlice1(length int) []bool {
	return make([]bool, length)
}

//test:disasm_both
// main.makeIntSlice2 code=15 frame=72 (3 slots: 2 args, 0 locals, 1 temps)
//   LoadScalarConst arg0 = 8
//   Move arg1 = length
//   Move arg2 = capacity
//   CallNative tmp0 = builtin.makeSlice()
//   Return tmp0
func makeIntSlice2(length, capacity int) []int {
	return make([]int, length, capacity)
}

//test:disasm_both
// main.makeByteSlice2 code=15 frame=72 (3 slots: 2 args, 0 locals, 1 temps)
//   LoadScalarConst arg0 = 1
//   Move arg1 = length
//   Move arg2 = capacity
//   CallNative tmp0 = builtin.makeSlice()
//   Return tmp0
func makeByteSlice2(length, capacity int) []byte {
	return make([]byte, length, capacity)
}

//test:disasm_both
// main.makeBoolSlice2 code=15 frame=72 (3 slots: 2 args, 0 locals, 1 temps)
//   LoadScalarConst arg0 = 1
//   Move arg1 = length
//   Move arg2 = capacity
//   CallNative tmp0 = builtin.makeSlice()
//   Return tmp0
func makeBoolSlice2(length, capacity int) []bool {
	return make([]bool, length, capacity)
}

//test:disasm_both
// main.intSliceLenCap code=19 frame=48 (2 slots: 1 args, 0 locals, 1 temps)
//   Len tmp0 = xs
//   Move arg0 = tmp0
//   CallVoidNative builtin.PrintInt()
//   Cap tmp0 = xs
//   Move arg0 = tmp0
//   CallVoidNative builtin.PrintInt()
//   ReturnVoid
func intSliceLenCap(xs []int) {
	println(len(xs))
	println(cap(xs))
}

//test:disasm_both
// main.byteSliceLenCap code=19 frame=48 (2 slots: 1 args, 0 locals, 1 temps)
//   Len tmp0 = xs
//   Move arg0 = tmp0
//   CallVoidNative builtin.PrintInt()
//   Cap tmp0 = xs
//   Move arg0 = tmp0
//   CallVoidNative builtin.PrintInt()
//   ReturnVoid
func byteSliceLenCap(xs []byte) {
	println(len(xs))
	println(cap(xs))
}

//test:disasm_both
// main.boolSliceLenCap code=19 frame=48 (2 slots: 1 args, 0 locals, 1 temps)
//   Len tmp0 = xs
//   Move arg0 = tmp0
//   CallVoidNative builtin.PrintInt()
//   Cap tmp0 = xs
//   Move arg0 = tmp0
//   CallVoidNative builtin.PrintInt()
//   ReturnVoid
func boolSliceLenCap(xs []bool) {
	println(len(xs))
	println(cap(xs))
}

//test:disasm_both
// main.intSliceIndexing code=8 frame=48 (2 slots: 2 args, 0 locals, 0 temps)
//   SliceIndexScalar64 arg0 = xs i
//   CallVoidNative builtin.PrintInt()
//   ReturnVoid
func intSliceIndexing(xs []int, i int) {
	println(xs[i])
}

//test:disasm_both
// main.byteSliceIndexing code=8 frame=48 (2 slots: 2 args, 0 locals, 0 temps)
//   SliceIndexScalar8 arg0 = xs i
//   CallVoidNative builtin.PrintByte()
//   ReturnVoid
func byteSliceIndexing(xs []byte, i int) {
	println(xs[i])
}

//test:disasm_both
// main.boolSliceIndexing code=8 frame=48 (2 slots: 2 args, 0 locals, 0 temps)
//   SliceIndexScalar8 arg0 = xs i
//   CallVoidNative builtin.PrintBool()
//   ReturnVoid
func boolSliceIndexing(xs []bool, i int) {
	println(xs[i])
}

//test:disasm_both
// main.intSliceAssign code=5 frame=72 (3 slots: 3 args, 0 locals, 0 temps)
//   SliceSetScalar64 xs i value
//   ReturnVoid
func intSliceAssign(xs []int, i, value int) {
	xs[i] = value
}

//test:disasm_both
// main.byteSliceAssign code=5 frame=72 (3 slots: 3 args, 0 locals, 0 temps)
//   SliceSetScalar8 xs i value
//   ReturnVoid
func byteSliceAssign(xs []byte, i int, value byte) {
	xs[i] = value
}

//test:disasm_both
// main.boolSliceAssign code=5 frame=72 (3 slots: 3 args, 0 locals, 0 temps)
//   SliceSetScalar8 xs i value
//   ReturnVoid
func boolSliceAssign(xs []bool, i int, value bool) {
	xs[i] = value
}

//test:disasm_both
// main.intSliceAppend code=12 frame=72 (3 slots: 2 args, 1 locals, 0 temps)
//   Move arg0 = xs
//   Move arg1 = value
//   CallNative out = builtin.append64()
//   Return out
func intSliceAppend(xs []int, value int) []int {
	out := append(xs, value)
	return out
}

//test:disasm_both
// main.byteSliceAppend code=12 frame=72 (3 slots: 2 args, 1 locals, 0 temps)
//   Move arg0 = xs
//   Move arg1 = value
//   CallNative out = builtin.append8()
//   Return out
func byteSliceAppend(xs []byte, value byte) []byte {
	out := append(xs, value)
	return out
}

//test:disasm_both
// main.boolSliceAppend code=12 frame=72 (3 slots: 2 args, 1 locals, 0 temps)
//   Move arg0 = xs
//   Move arg1 = value
//   CallNative out = builtin.append8()
//   Return out
func boolSliceAppend(xs []bool, value bool) []bool {
	out := append(xs, value)
	return out
}

func testIntSlice() {
	intSliceLenCap(makeIntSlice1(10))
	intSliceLenCap(makeIntSlice2(3, 11))

	s := make([]int, 1, 3)
	intSliceIndexing(s, 0)
	intSliceAssign(s, 0, 152948)
	intSliceIndexing(s, 0)
	s[0] = -1
	println(s[0])
	for i := 5; i <= 10; i++ {
		println(len(s))
		println(cap(s))
		s = intSliceAppend(s, i)
		println(s[len(s)-2])
		println(s[len(s)-1])
	}
	println(len(s))
	println(cap(s))
}

func testByteSliceOverflow() {
	elems := make([]byte, 0, 1024)
	for i := 0; i < cap(elems); i++ {
		elems = append(elems, byte(i))
		println(elems[len(elems)-1])
	}
}

func testByteSlice() {
	byteSliceLenCap(makeByteSlice1(10))
	byteSliceLenCap(makeByteSlice2(3, 11))

	s := make([]byte, 1, 3)
	byteSliceIndexing(s, 0)
	byteSliceAssign(s, 0, 32)
	byteSliceIndexing(s, 0)
	s[0] = 100
	println(s[0])
	for i := 5; i <= 10; i++ {
		println(len(s))
		println(cap(s))
		s = byteSliceAppend(s, byte(i))
		println(s[len(s)-2])
		println(s[len(s)-1])
	}
	println(len(s))
	println(cap(s))

	testByteSliceOverflow()
}

func testBoolSlice() {
	boolSliceLenCap(makeBoolSlice1(10))
	boolSliceLenCap(makeBoolSlice2(3, 11))

	s := make([]bool, 1, 3)
	boolSliceIndexing(s, 0)
	boolSliceAssign(s, 0, true)
	boolSliceIndexing(s, 0)
	s[0] = false
	println(s[0])
	for i := 5; i <= 10; i++ {
		println(len(s))
		println(cap(s))
		s = boolSliceAppend(s, true)
		println(s[len(s)-2])
		println(s[len(s)-1])
	}
	println(len(s))
	println(cap(s))
}

func main() {
	testIntSlice()
	testByteSlice()
	testBoolSlice()
}
