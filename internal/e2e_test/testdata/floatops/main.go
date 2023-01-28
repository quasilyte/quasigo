package main

//test:disasm_both
// main.intToFloat code=5 frame=48 (2 slots: 1 params, 1 locals)
//   ConvIntToFloat temp0 = x
//   ReturnScalar temp0
func intToFloat(x int) float64 {
	return float64(x)
}

//test:disasm_both
// main.floatadd code=6 frame=72 (3 slots: 2 params, 1 locals)
//   FloatAdd64 temp0 = x y
//   ReturnScalar temp0
func floatadd(x, y float64) float64 {
	return x + y
}

//test:disasm_both
// main.floatsub code=6 frame=72 (3 slots: 2 params, 1 locals)
//   FloatSub64 temp0 = x y
//   ReturnScalar temp0
func floatsub(x, y float64) float64 {
	return x - y
}

//test:disasm_both
// main.floateq code=6 frame=72 (3 slots: 2 params, 1 locals)
//   ScalarEq temp0 = x y
//   ReturnScalar temp0
func floateq(x, y float64) bool {
	return x == y
}

//test:disasm_both
// main.floatmul code=6 frame=72 (3 slots: 2 params, 1 locals)
//   FloatMul64 temp0 = x y
//   ReturnScalar temp0
func floatmul(x, y float64) float64 {
	return x * y
}

//test:disasm_both
// main.floatdiv code=6 frame=72 (3 slots: 2 params, 1 locals)
//   FloatDiv64 temp0 = x y
//   ReturnScalar temp0
func floatdiv(x, y float64) float64 {
	return x / y
}

//test:disasm_both
// main.abs code=17 frame=72 (3 slots: 1 params, 2 locals)
//   Zero temp1
//   FloatLt temp0 = x temp1
//   JumpZero L0 temp0
//   FloatNeg temp0 = x
//   ReturnScalar temp0
// L0:
//   ReturnScalar x
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

//test:disasm_both
// main.sqrt code=49 frame=144 (6 slots: 1 params, 5 locals)
//   LoadScalarConst temp0 = 4607182418800017408
//   Jump L0
// L1:
//   FloatDiv64 temp2 = x temp0
//   FloatAdd64 temp1 = temp0 temp2
//   LoadScalarConst temp2 = 4611686018427387904
//   FloatDiv64 temp0 = temp1 temp2
// L0:
//   FloatDiv64 temp4 = x temp0
//   FloatSub64 temp3 = temp4 temp0
//   Move arg0 = temp3
//   Call temp2 = main.abs()
//   LoadScalarConst temp3 = 4532020583610935537
//   FloatGt temp1 = temp2 temp3
//   JumpNotZero L1 temp1
//   ReturnScalar temp0
func sqrt(x float64) float64 {
	y := 1.0
	for abs(x/y-y) > 0.00001 {
		y = (y + x/y) / 2
	}
	return y
}

func testArithOps(x, y float64) {
	println(floatadd(x, y))
	println(floatmul(x, y))
	println(floateq(x, y))
	println(floatsub(x, y))
}

func testBasicOps() {
	testArithOps(0, 0)
	testArithOps(1.4, -1.4)
	testArithOps(0.001, 0.001)
	testArithOps(1.3, 0)
	testArithOps(14.4, 1.5)
	testArithOps(124.6, 1)
	testArithOps(1.5, 249)

	println(floatdiv(1.5, 1.0))
	println(floatdiv(1.5, 0.4))
	println(floatdiv(0.3, 4.6))
	println(floatdiv(324, 5))

	println(abs(0))
	println(abs(14))
	println(abs(4.2))
	println(abs(-428.3))
	println(abs(-1.4))

	println(sqrt(0))
	println(sqrt(1))
	println(sqrt(6))
	println(sqrt(16))
	println(sqrt(36))
	println(sqrt(36.5))
	println(sqrt(81))
	println(sqrt(81.7))
	println(sqrt(81.2385))
	println(sqrt(900.101))
	println(sqrt(1024.5))

	println(intToFloat(0))
	println(intToFloat(35))
	println(intToFloat(46))
	println(intToFloat(-46))
	println(intToFloat(438))
	println(intToFloat(8328438))
	println(intToFloat(9993281472773771))
}

func main() {
	testBasicOps()
}
