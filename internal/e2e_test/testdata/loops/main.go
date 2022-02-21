package main

const one int = 1

//test:disasm_both
// main.cstyle1 code=26 frame=72 (3 slots: 0 args, 1 locals, 2 temps)
//   LoadScalarConst i = 0
//   Jump L0
// L1:
//   MoveScalar arg0 = i
//   CallVoidNative builtin.PrintInt()
//   IntInc i
// L0:
//   LoadScalarConst tmp1 = 3
//   IntLt tmp0 = i tmp1
//   JumpTrue L1 tmp0
//   ReturnVoid
func cstyle1() {
	for i := 0; i < 3; i++ {
		println(i)
	}
}

// TODO: opt: remove dead store to i in the beginning of this func.
//
//test:disasm_both
// main.cstyle2 code=35 frame=72 (3 slots: 0 args, 1 locals, 2 temps)
//   LoadScalarConst i = 10
//   LoadScalarConst i = 0
//   Jump L0
// L1:
//   MoveScalar arg0 = i
//   CallVoidNative builtin.PrintInt()
//   IntInc i
// L0:
//   LoadScalarConst tmp1 = 3
//   IntLt tmp0 = i tmp1
//   JumpTrue L1 tmp0
//   MoveScalar arg0 = i
//   CallVoidNative builtin.PrintInt()
//   ReturnVoid
func cstyle2() {
	i := 10
	for i = 0; i < 3; i++ {
		println(i)
	}
	println(i)
}

//test:disasm_both
// main.cstyle3 code=46 frame=72 (3 slots: 0 args, 1 locals, 2 temps)
//   LoadScalarConst i = 0
//   Jump L0
// L3:
//   MoveScalar arg0 = i
//   CallVoidNative builtin.PrintInt()
//   LoadScalarConst tmp1 = 5
//   IntGt tmp0 = i tmp1
//   JumpFalse L1 tmp0
//   Jump L2
// L1:
//   LoadStrConst arg0 = "after continue"
//   CallVoidNative builtin.PrintString()
// L2:
//   IntInc i
// L0:
//   LoadScalarConst tmp1 = 10
//   IntLt tmp0 = i tmp1
//   JumpTrue L3 tmp0
//   ReturnVoid
func cstyle3() {
	for i := 0; i < 10; i++ {
		println(i)
		if i > 5 {
			continue
		}
		println("after continue")
	}
}

//test:disasm_both
// main.cstyle4 code=46 frame=72 (3 slots: 0 args, 1 locals, 2 temps)
//   LoadScalarConst i = 10
//   Jump L0
// L3:
//   MoveScalar arg0 = i
//   CallVoidNative builtin.PrintInt()
//   LoadScalarConst tmp1 = 5
//   IntLtEq tmp0 = i tmp1
//   JumpFalse L1 tmp0
//   Jump L2
// L1:
//   LoadStrConst arg0 = "after break"
//   CallVoidNative builtin.PrintString()
//   IntDec i
// L0:
//   LoadScalarConst tmp1 = 0
//   IntGtEq tmp0 = i tmp1
//   JumpTrue L3 tmp0
// L2:
//   ReturnVoid
func cstyle4() {
	i := 10
	for ; i >= 0; i-- {
		println(i)
		if i <= 5 {
			break
		}
		println("after break")
	}
}

//test:disasm_both
// main.cstyle5 code=35 frame=72 (3 slots: 0 args, 1 locals, 2 temps)
//   LoadScalarConst i = 10
// L2:
//   MoveScalar arg0 = i
//   CallVoidNative builtin.PrintInt()
//   LoadScalarConst tmp1 = 5
//   IntLtEq tmp0 = i tmp1
//   JumpFalse L0 tmp0
//   Jump L1
// L0:
//   LoadStrConst arg0 = "after break"
//   CallVoidNative builtin.PrintString()
//   IntDec i
//   Jump L2
// L1:
//   ReturnVoid
func cstyle5() {
	i := 10
	for ; ; i-- {
		println(i)
		if i <= 5 {
			break
		}
		println("after break")
	}
}

//test:disasm_both
// main.cstyle6 code=35 frame=72 (3 slots: 0 args, 1 locals, 2 temps)
//   LoadScalarConst i = 10
// L2:
//   MoveScalar arg0 = i
//   CallVoidNative builtin.PrintInt()
//   LoadScalarConst tmp1 = 5
//   IntLtEq tmp0 = i tmp1
//   JumpFalse L0 tmp0
//   Jump L1
// L0:
//   LoadStrConst arg0 = "after break"
//   CallVoidNative builtin.PrintString()
//   IntDec i
//   Jump L2
// L1:
//   ReturnVoid
func cstyle6() {
	for i := 10; ; i-- {
		println(i)
		if i <= 5 {
			break
		}
		println("after break")
	}
}

//test:disasm_both
// main.cstyle7 code=35 frame=72 (3 slots: 0 args, 1 locals, 2 temps)
//   LoadScalarConst i = 0
// L2:
//   IntInc i
//   MoveScalar arg0 = i
//   CallVoidNative builtin.PrintInt()
//   LoadScalarConst tmp1 = 5
//   IntEq tmp0 = i tmp1
//   JumpFalse L0 tmp0
//   Jump L1
// L0:
//   LoadStrConst arg0 = "after break"
//   CallVoidNative builtin.PrintString()
//   Jump L2
// L1:
//   ReturnVoid
func cstyle7() {
	for i := 0; ; {
		i++
		println(i)
		if i == 5 {
			break
		}
		println("after break")
	}
}

func testWhile() {
	// While-style loops.
	{
		i := 0
		for i < 5 {
			println(i)
			i++
		}
	}
	{
		i2 := 2
		for i2 < 2 {
			println(i2)
			i2 += one
		}
	}
}

func main() {
	testWhile()
	cstyle1()
	cstyle2()
	cstyle3()
	cstyle4()
	cstyle5()
	cstyle6()
	cstyle7()
}
