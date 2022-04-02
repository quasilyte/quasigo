package main

const one int = 1

//test:disasm_both
// main.cstyle1 code=25 frame=72 (3 slots: 0 params, 3 locals)
//   Zero temp0
//   Jump L0
// L1:
//   Move arg0 = temp0
//   CallVoidNative builtin.PrintInt()
//   IntInc temp0
// L0:
//   LoadScalarConst temp2 = 3
//   IntLt temp1 = temp0 temp2
//   JumpNotZero L1 temp1
//   ReturnVoid
func cstyle1() {
	for i := 0; i < 3; i++ {
		println(i)
	}
}

// TODO: opt: remove dead store to i in the beginning of this func.
//
//test:disasm_both
// main.cstyle2 code=34 frame=72 (3 slots: 0 params, 3 locals)
//   LoadScalarConst temp0 = 10
//   Zero temp0
//   Jump L0
// L1:
//   Move arg0 = temp0
//   CallVoidNative builtin.PrintInt()
//   IntInc temp0
// L0:
//   LoadScalarConst temp2 = 3
//   IntLt temp1 = temp0 temp2
//   JumpNotZero L1 temp1
//   Move arg0 = temp0
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
// main.cstyle3 code=45 frame=72 (3 slots: 0 params, 3 locals)
//   Zero temp0
//   Jump L0
// L3:
//   Move arg0 = temp0
//   CallVoidNative builtin.PrintInt()
//   LoadScalarConst temp2 = 5
//   IntGt temp1 = temp0 temp2
//   JumpZero L1 temp1
//   Jump L2
// L1:
//   LoadStrConst arg0 = "after continue"
//   CallVoidNative builtin.PrintString()
// L2:
//   IntInc temp0
// L0:
//   LoadScalarConst temp2 = 10
//   IntLt temp1 = temp0 temp2
//   JumpNotZero L3 temp1
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
// main.cstyle4 code=45 frame=72 (3 slots: 0 params, 3 locals)
//   LoadScalarConst temp0 = 10
//   Jump L0
// L3:
//   Move arg0 = temp0
//   CallVoidNative builtin.PrintInt()
//   LoadScalarConst temp2 = 5
//   IntLtEq temp1 = temp0 temp2
//   JumpZero L1 temp1
//   Jump L2
// L1:
//   LoadStrConst arg0 = "after break"
//   CallVoidNative builtin.PrintString()
//   IntDec temp0
// L0:
//   Zero temp2
//   IntGtEq temp1 = temp0 temp2
//   JumpNotZero L3 temp1
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
// main.cstyle5 code=35 frame=72 (3 slots: 0 params, 3 locals)
//   LoadScalarConst temp0 = 10
// L2:
//   Move arg0 = temp0
//   CallVoidNative builtin.PrintInt()
//   LoadScalarConst temp2 = 5
//   IntLtEq temp1 = temp0 temp2
//   JumpZero L0 temp1
//   Jump L1
// L0:
//   LoadStrConst arg0 = "after break"
//   CallVoidNative builtin.PrintString()
//   IntDec temp0
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
// main.cstyle6 code=37 frame=96 (4 slots: 1 params, 3 locals)
//   IntDec n
//   LoadScalarConst temp0 = 10
// L2:
//   Move arg0 = temp0
//   CallVoidNative builtin.PrintInt()
//   LoadScalarConst temp2 = 5
//   IntLtEq temp1 = temp0 temp2
//   JumpZero L0 temp1
//   Jump L1
// L0:
//   LoadStrConst arg0 = "after break"
//   CallVoidNative builtin.PrintString()
//   IntDec temp0
//   Jump L2
// L1:
//   ReturnVoid
func cstyle6(n int) {
	n--
	for i := 10; ; i-- {
		println(i)
		if i <= 5 {
			break
		}
		println("after break")
	}
}

//test:disasm_both
// main.cstyle7 code=34 frame=72 (3 slots: 0 params, 3 locals)
//   Zero temp0
// L2:
//   IntInc temp0
//   Move arg0 = temp0
//   CallVoidNative builtin.PrintInt()
//   LoadScalarConst temp2 = 5
//   ScalarEq temp1 = temp0 temp2
//   JumpZero L0 temp1
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

//test:irdump
// block0 (L1) [0]:
//   LoadScalarConst temp0 = -5
// block1 [0]:
//   Zero temp2.v0
//   IntGt temp1.v0 = temp0 temp2.v0
//   JumpZero L2 temp1.v0
// block2 [0]:
//   Jump L0
// block3 (L2) [0]:
//   IntInc temp0
//   Move arg0 = temp0
//   CallVoidNative builtin.PrintInt
//   Jump L1
// block4 (L0) [0]:
//   ReturnVoid
//
//test:disasm_both
// main.while1 code=28 frame=72 (3 slots: 0 params, 3 locals)
//   LoadScalarConst temp0 = -5
// L2:
//   Zero temp2
//   IntGt temp1 = temp0 temp2
//   JumpZero L0 temp1
//   Jump L1
// L0:
//   IntInc temp0
//   Move arg0 = temp0
//   CallVoidNative builtin.PrintInt()
//   Jump L2
// L1:
//   ReturnVoid
func while1() {
	j := -5
	for {
		if j > 0 {
			break
		}
		j++
		println(j)
	}
}

func testWhile() {
	// While-style loops.
	while1()
	{
		i := 0
		for i < 5 {
			println(i)
			i++
		}
	}
	{
		i2 := 0
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
	cstyle6(15)
	cstyle7()
}
