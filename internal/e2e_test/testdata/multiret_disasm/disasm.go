package multiretdisasm

// Functions that return 2 values use a single Return2 op that copies
// both values as whole slots, so it works for any result types.

//test:disasm_both
// multiretdisasm.swap code=3 frame=48 (2 slots: 2 params, 0 locals)
//   Return2 y x
func swap(x, y int) (int, int) {
	return y, x
}

//test:disasm_both
// multiretdisasm.consts code=9 frame=48 (2 slots: 0 params, 2 locals)
//   LoadScalarConst temp0 = 10
//   LoadScalarConst temp1 = 20
//   Return2 temp0 temp1
func consts() (int, int) {
	return 10, 20
}

//test:disasm_both
// multiretdisasm.strBool code=21 frame=72 (3 slots: 1 params, 2 locals)
//   JumpZero L0 b
//   LoadStrConst temp0 = "yes"
//   LoadScalarConst temp1 = 1
//   Return2 temp0 temp1
// L0:
//   LoadStrConst temp0 = "no"
//   Zero temp1
//   Return2 temp0 temp1
func strBool(b bool) (string, bool) {
	if b {
		return "yes", true
	}
	return "no", false
}

//test:disasm_both
// multiretdisasm.forward code=15 frame=96 (4 slots: 2 params, 2 locals)
//   Move arg0 = x
//   Move arg1 = y
//   Call temp0 = multiretdisasm.swap()
//   MoveResult2 temp1
//   Return2 temp0 temp1
func forward(x, y int) (int, int) {
	return swap(x, y)
}

// The first result is stored by the call itself,
// the second one is moved right after the call.

//test:disasm_both
// multiretdisasm.define code=18 frame=72 (3 slots: 0 params, 3 locals)
//   LoadScalarConst arg0 = 1
//   LoadScalarConst arg1 = 2
//   Call temp0 = multiretdisasm.swap()
//   MoveResult2 temp1
//   IntSub64 temp2 = temp0 temp1
//   ReturnScalar temp2
func define() int {
	a, b := swap(1, 2)
	return a - b
}

//test:disasm_both
// multiretdisasm.assign code=18 frame=72 (3 slots: 2 params, 1 locals)
//   Move arg0 = a
//   Move arg1 = b
//   Call a = multiretdisasm.swap()
//   MoveResult2 b
//   IntSub64 temp0 = a b
//   ReturnScalar temp0
func assign(a, b int) int {
	a, b = swap(a, b)
	return a - b
}

// A blank identifier discards the result.

//test:disasm_both
// multiretdisasm.blankFirst code=13 frame=24 (1 slots: 0 params, 1 locals)
//   LoadScalarConst arg0 = 1
//   LoadScalarConst arg1 = 2
//   CallVoid multiretdisasm.swap()
//   MoveResult2 temp0
//   ReturnScalar temp0
func blankFirst() int {
	_, b := swap(1, 2)
	return b
}

//test:disasm_both
// multiretdisasm.blankSecond code=12 frame=24 (1 slots: 0 params, 1 locals)
//   LoadScalarConst arg0 = 1
//   LoadScalarConst arg1 = 2
//   Call temp0 = multiretdisasm.swap()
//   ReturnScalar temp0
func blankSecond() int {
	a, _ := swap(1, 2)
	return a
}

//test:disasm_both
// multiretdisasm.blankBoth code=10 frame=0 (0 slots: 0 params, 0 locals)
//   LoadScalarConst arg0 = 1
//   LoadScalarConst arg1 = 2
//   CallVoid multiretdisasm.swap()
//   ReturnVoid
func blankBoth() {
	_, _ = swap(1, 2)
}

//test:disasm
// multiretdisasm.recur code=33 frame=96 (4 slots: 1 params, 3 locals)
//   Zero temp1
//   ScalarEq temp0 = n temp1
//   JumpZero L0 temp0
//   Zero temp0
//   Zero temp1
//   Return2 temp0 temp1
// L0:
//   LoadScalarConst temp2 = 1
//   IntSub64 temp1 = n temp2
//   Move arg0 = temp1
//   CallRecurVoid
//   MoveResult2 temp0
//   Return2 n temp0
//
//test:disasm_opt
// multiretdisasm.recur code=27 frame=96 (4 slots: 1 params, 3 locals)
//   JumpNotZero L0 n
//   Zero temp0
//   Zero temp1
//   Return2 temp0 temp1
// L0:
//   LoadScalarConst temp2 = 1
//   IntSub64 temp1 = n temp2
//   Move arg0 = temp1
//   CallRecurVoid
//   MoveResult2 temp0
//   Return2 n temp0
func recur(n int) (int, int) {
	if n == 0 {
		return 0, 0
	}
	_, b := recur(n - 1)
	return n, b
}
