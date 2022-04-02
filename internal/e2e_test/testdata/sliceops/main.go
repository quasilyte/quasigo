package main

//test:disasm_both
// main.makeIntSlice1 code=15 frame=48 (2 slots: 1 params, 1 locals)
//   LoadScalarConst arg0 = 8
//   Move arg1 = length
//   Move arg2 = length
//   CallNative temp0 = builtin.makeSlice()
//   Return temp0
func makeIntSlice1(length int) []int {
	return make([]int, length)
}

//test:disasm_both
// main.makeFloatSlice1 code=15 frame=48 (2 slots: 1 params, 1 locals)
//   LoadScalarConst arg0 = 8
//   Move arg1 = length
//   Move arg2 = length
//   CallNative temp0 = builtin.makeSlice()
//   Return temp0
func makeFloatSlice1(length int) []float64 {
	return make([]float64, length)
}

//test:disasm_both
// main.makeByteSlice1 code=15 frame=48 (2 slots: 1 params, 1 locals)
//   LoadScalarConst arg0 = 1
//   Move arg1 = length
//   Move arg2 = length
//   CallNative temp0 = builtin.makeSlice()
//   Return temp0
func makeByteSlice1(length int) []byte {
	return make([]byte, length)
}

//test:disasm_both
// main.makeBoolSlice1 code=15 frame=48 (2 slots: 1 params, 1 locals)
//   LoadScalarConst arg0 = 1
//   Move arg1 = length
//   Move arg2 = length
//   CallNative temp0 = builtin.makeSlice()
//   Return temp0
func makeBoolSlice1(length int) []bool {
	return make([]bool, length)
}

//test:disasm_both
// main.makeIntSlice2 code=15 frame=72 (3 slots: 2 params, 1 locals)
//   LoadScalarConst arg0 = 8
//   Move arg1 = length
//   Move arg2 = capacity
//   CallNative temp0 = builtin.makeSlice()
//   Return temp0
func makeIntSlice2(length, capacity int) []int {
	return make([]int, length, capacity)
}

//test:disasm_both
// main.makeFloatSlice2 code=15 frame=72 (3 slots: 2 params, 1 locals)
//   LoadScalarConst arg0 = 8
//   Move arg1 = length
//   Move arg2 = capacity
//   CallNative temp0 = builtin.makeSlice()
//   Return temp0
func makeFloatSlice2(length, capacity int) []float64 {
	return make([]float64, length, capacity)
}

//test:disasm_both
// main.makeByteSlice2 code=15 frame=72 (3 slots: 2 params, 1 locals)
//   LoadScalarConst arg0 = 1
//   Move arg1 = length
//   Move arg2 = capacity
//   CallNative temp0 = builtin.makeSlice()
//   Return temp0
func makeByteSlice2(length, capacity int) []byte {
	return make([]byte, length, capacity)
}

//test:disasm_both
// main.makeBoolSlice2 code=15 frame=72 (3 slots: 2 params, 1 locals)
//   LoadScalarConst arg0 = 1
//   Move arg1 = length
//   Move arg2 = capacity
//   CallNative temp0 = builtin.makeSlice()
//   Return temp0
func makeBoolSlice2(length, capacity int) []bool {
	return make([]bool, length, capacity)
}

//test:disasm_both
// main.intSliceLenCap code=19 frame=48 (2 slots: 1 params, 1 locals)
//   Len temp0 = xs
//   Move arg0 = temp0
//   CallVoidNative builtin.PrintInt()
//   Cap temp0 = xs
//   Move arg0 = temp0
//   CallVoidNative builtin.PrintInt()
//   ReturnVoid
func intSliceLenCap(xs []int) {
	println(len(xs))
	println(cap(xs))
}

//test:disasm_both
// main.floatSliceLenCap code=19 frame=48 (2 slots: 1 params, 1 locals)
//   Len temp0 = xs
//   Move arg0 = temp0
//   CallVoidNative builtin.PrintInt()
//   Cap temp0 = xs
//   Move arg0 = temp0
//   CallVoidNative builtin.PrintInt()
//   ReturnVoid
func floatSliceLenCap(xs []float64) {
	println(len(xs))
	println(cap(xs))
}

//test:disasm_both
// main.byteSliceLenCap code=19 frame=48 (2 slots: 1 params, 1 locals)
//   Len temp0 = xs
//   Move arg0 = temp0
//   CallVoidNative builtin.PrintInt()
//   Cap temp0 = xs
//   Move arg0 = temp0
//   CallVoidNative builtin.PrintInt()
//   ReturnVoid
func byteSliceLenCap(xs []byte) {
	println(len(xs))
	println(cap(xs))
}

//test:disasm_both
// main.boolSliceLenCap code=19 frame=48 (2 slots: 1 params, 1 locals)
//   Len temp0 = xs
//   Move arg0 = temp0
//   CallVoidNative builtin.PrintInt()
//   Cap temp0 = xs
//   Move arg0 = temp0
//   CallVoidNative builtin.PrintInt()
//   ReturnVoid
func boolSliceLenCap(xs []bool) {
	println(len(xs))
	println(cap(xs))
}

//test:disasm_both
// main.intSliceIndexing code=8 frame=48 (2 slots: 2 params, 0 locals)
//   SliceIndexScalar64 arg0 = xs i
//   CallVoidNative builtin.PrintInt()
//   ReturnVoid
func intSliceIndexing(xs []int, i int) {
	println(xs[i])
}

//test:disasm_both
// main.floatSliceIndexing code=8 frame=48 (2 slots: 2 params, 0 locals)
//   SliceIndexScalar64 arg0 = xs i
//   CallVoidNative builtin.PrintFloat()
//   ReturnVoid
func floatSliceIndexing(xs []float64, i int) {
	println(xs[i])
}

//test:disasm_both
// main.byteSliceIndexing code=8 frame=48 (2 slots: 2 params, 0 locals)
//   SliceIndexScalar8 arg0 = xs i
//   CallVoidNative builtin.PrintByte()
//   ReturnVoid
func byteSliceIndexing(xs []byte, i int) {
	println(xs[i])
}

//test:disasm_both
// main.boolSliceIndexing code=8 frame=48 (2 slots: 2 params, 0 locals)
//   SliceIndexScalar8 arg0 = xs i
//   CallVoidNative builtin.PrintBool()
//   ReturnVoid
func boolSliceIndexing(xs []bool, i int) {
	println(xs[i])
}

//test:disasm_both
// main.intSliceAssign code=5 frame=72 (3 slots: 3 params, 0 locals)
//   SliceSetScalar64 xs i value
//   ReturnVoid
func intSliceAssign(xs []int, i, value int) {
	xs[i] = value
}

//test:disasm_both
// main.floatSliceAssign code=5 frame=72 (3 slots: 3 params, 0 locals)
//   SliceSetScalar64 xs i value
//   ReturnVoid
func floatSliceAssign(xs []float64, i int, value float64) {
	xs[i] = value
}

//test:disasm_both
// main.byteSliceAssign code=5 frame=72 (3 slots: 3 params, 0 locals)
//   SliceSetScalar8 xs i value
//   ReturnVoid
func byteSliceAssign(xs []byte, i int, value byte) {
	xs[i] = value
}

//test:disasm_both
// main.boolSliceAssign code=5 frame=72 (3 slots: 3 params, 0 locals)
//   SliceSetScalar8 xs i value
//   ReturnVoid
func boolSliceAssign(xs []bool, i int, value bool) {
	xs[i] = value
}

//test:disasm_both
// main.intSliceAppend code=12 frame=72 (3 slots: 2 params, 1 locals)
//   Move arg0 = xs
//   Move arg1 = value
//   CallNative temp0 = builtin.append64()
//   Return temp0
func intSliceAppend(xs []int, value int) []int {
	out := append(xs, value)
	return out
}

//test:disasm_both
// main.floatSliceAppend code=12 frame=72 (3 slots: 2 params, 1 locals)
//   Move arg0 = xs
//   Move arg1 = value
//   CallNative temp0 = builtin.append64()
//   Return temp0
func floatSliceAppend(xs []float64, value float64) []float64 {
	out := append(xs, value)
	return out
}

//test:disasm_both
// main.byteSliceAppend code=12 frame=48 (2 slots: 2 params, 0 locals)
//   Move arg0 = xs
//   Move arg1 = value
//   CallNative xs = builtin.append8()
//   Return xs
func byteSliceAppend(xs []byte, value byte) []byte {
	xs = append(xs, value)
	return xs
}

//test:disasm_both
// main.boolSliceAppend code=12 frame=72 (3 slots: 2 params, 1 locals)
//   Move arg0 = xs
//   Move arg1 = value
//   CallNative temp0 = builtin.append8()
//   Return temp0
func boolSliceAppend(xs []bool, value bool) []bool {
	return append(xs, value)
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

func divFloatSlices(xs, ys []float64) []float64 {
	result := make([]float64, len(xs))
	for i := 0; i < len(xs); i++ {
		result[i] = xs[i] / ys[i]
	}
	return result
}

func testFloatSlice() {
	floatSliceLenCap(makeFloatSlice1(10))
	floatSliceLenCap(makeFloatSlice2(3, 11))

	s := make([]float64, 1, 3)
	floatSliceIndexing(s, 0)
	floatSliceAssign(s, 0, 152948)
	floatSliceIndexing(s, 0)
	s[0] = -1
	println(s[0])
	s[0] = 1
	println(s[0])
	s[0] = 143.6
	println(s[0])
	for i := 5; i <= 10; i++ {
		println(len(s))
		println(cap(s))
		s = floatSliceAppend(s, float64(i))
		println(s[len(s)-2])
		println(s[len(s)-1])
	}
	println(len(s))
	println(cap(s))

	{
		xs := make([]float64, 3)
		xs[0] = 124.5
		xs[1] = 493.2
		xs[2] = 294.0
		ys := make([]float64, 3)
		ys[0] = 24.5
		ys[1] = 1
		ys[2] = 2.5
		result := divFloatSlices(xs, ys)
		for i := 0; i < len(result); i++ {
			println(result[i])
		}
	}
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
	testFloatSlice()
	testByteSlice()
	testBoolSlice()
}
