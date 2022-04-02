package main

//test:disasm
// main.test3withDefault code=56 frame=96 (4 slots: 1 params, 3 locals)
//   Move temp0 = x
//   LoadScalarConst temp2 = 10
//   ScalarEq temp1 = temp0 temp2
//   JumpZero L0 temp1
//   LoadStrConst temp1 = "ten"
//   ReturnStr temp1
// L0:
//   LoadScalarConst temp2 = 20
//   ScalarEq temp1 = temp0 temp2
//   JumpZero L1 temp1
//   LoadStrConst temp1 = "twenty"
//   ReturnStr temp1
// L1:
//   LoadScalarConst temp2 = 30
//   ScalarEq temp1 = temp0 temp2
//   JumpZero L2 temp1
//   LoadStrConst temp1 = "thirty"
//   ReturnStr temp1
// L2:
//   LoadStrConst temp1 = "?"
//   ReturnStr temp1
func test3withDefault(x byte) string {
	switch x {
	case 10:
		return "ten"
	case 20:
		return "twenty"
	case 30:
		return "thirty"
	default:
		return "?"
	}
}

func test5withDefault(x byte) string {
	switch x {
	case 75:
		return "e"
	case 60:
		return "d"
	case 45:
		return "c"
	case 30:
		return "b"
	case 15:
		return "a"
	default:
		return "?"
	}
}

func test5noDefault(x byte) string {
	res := "?"
	switch x {
	case 75:
		res = "e"
	case 60:
		res = "d"
	case 45:
		res = "c"
	case 30:
		res = "b"
	case 15:
		res = "a"
	}
	return res
}

//test:disasm
// main.test10withDefault code=202 frame=96 (4 slots: 1 params, 3 locals)
//   Move temp0 = x
//   Zero temp2
//   IntLt temp1 = temp0 temp2
//   JumpNotZero L0 temp1
//   LoadScalarConst temp2 = 90
//   IntGt temp1 = temp0 temp2
//   JumpNotZero L0 temp1
//   LoadScalarConst temp2 = 4
//   ScalarEq temp1 = temp0 temp2
//   JumpNotZero L1 temp1
//   IntGt temp1 = temp0 temp2
//   JumpNotZero L2 temp1
//   Zero temp2
//   ScalarEq temp1 = temp0 temp2
//   JumpNotZero L3 temp1
//   LoadScalarConst temp2 = 1
//   ScalarEq temp1 = temp0 temp2
//   JumpNotZero L4 temp1
//   LoadScalarConst temp2 = 2
//   ScalarEq temp1 = temp0 temp2
//   JumpNotZero L5 temp1
//   LoadScalarConst temp2 = 3
//   ScalarEq temp1 = temp0 temp2
//   JumpNotZero L6 temp1
//   Jump L0
// L2:
//   LoadScalarConst temp2 = 5
//   ScalarEq temp1 = temp0 temp2
//   JumpNotZero L7 temp1
//   LoadScalarConst temp2 = 6
//   ScalarEq temp1 = temp0 temp2
//   JumpNotZero L8 temp1
//   LoadScalarConst temp2 = 7
//   ScalarEq temp1 = temp0 temp2
//   JumpNotZero L9 temp1
//   LoadScalarConst temp2 = 8
//   ScalarEq temp1 = temp0 temp2
//   JumpNotZero L10 temp1
//   LoadScalarConst temp2 = 90
//   ScalarEq temp1 = temp0 temp2
//   JumpNotZero L11 temp1
//   Jump L0
// L3:
//   LoadStrConst temp1 = "0"
//   ReturnStr temp1
// L4:
//   LoadStrConst temp1 = "1"
//   ReturnStr temp1
// L5:
//   LoadStrConst temp1 = "2"
//   ReturnStr temp1
// L6:
//   LoadStrConst temp1 = "3"
//   ReturnStr temp1
// L1:
//   LoadStrConst temp1 = "4"
//   ReturnStr temp1
// L7:
//   LoadStrConst temp1 = "5"
//   ReturnStr temp1
// L8:
//   LoadStrConst temp1 = "6"
//   ReturnStr temp1
// L9:
//   LoadStrConst temp1 = "7"
//   ReturnStr temp1
// L10:
//   LoadStrConst temp1 = "8"
//   ReturnStr temp1
// L11:
//   LoadStrConst temp1 = "90"
//   ReturnStr temp1
// L0:
//   LoadStrConst temp1 = "?"
//   ReturnStr temp1
func test10withDefault(x byte) string {
	switch x {
	case 0:
		return "0"
	case 1:
		return "1"
	case 2:
		return "2"
	case 3:
		return "3"
	case 4:
		return "4"
	case 5:
		return "5"
	case 6:
		return "6"
	case 7:
		return "7"
	case 8:
		return "8"
	case 90:
		return "90"
	default:
		return "?"
	}
}

func test10noDefault(x byte) string {
	res := "?"
	switch x {
	case 0:
		res = "0"
	case 1:
		res = "1"
	case 2:
		res = "2"
	case 3:
		res = "3"
	case 4:
		res = "4"
	case 5:
		res = "5"
	case 6:
		res = "6"
	case 7:
		res = "7"
	case 8:
		res = "8"
	case 9:
		res = "9"
	}
	return res
}

func test21withDefault(x byte) string {
	switch x {
	case 1:
		return "<0>"
	case 2:
		return "<1>"
	case 0:
		return "<2>"
	case 5:
		return "<3>"
	case 10:
		return "<4>"
	case 15:
		return "<5>"
	case 20:
		return "<6>"
	case 25:
		return "<7>"
	case 30:
		return "<8>"
	case 35:
		return "<9>"
	case 40:
		return "<10>"
	case 45:
		return "<11>"
	case 50:
		return "<12>"
	case 55:
		return "<13>"
	case 60:
		return "<14>"
	case 65:
		return "<15>"
	case 70:
		return "<16>"
	case 75:
		return "<17>"
	case 80:
		return "<18>"
	case 85:
		return "<19>"
	case 90:
		return "<20>"
	default:
		return "?"
	}
}

func test21noDefault(x byte) string {
	res := "?"
	switch x {
	case 1:
		res = "<0>"
	case 2:
		res = "<1>"
	case 0:
		res = "<2>"
	case 5:
		res = "<3>"
	case 10:
		res = "<4>"
	case 15:
		res = "<5>"
	case 20:
		res = "<6>"
	case 25:
		res = "<7>"
	case 30:
		res = "<8>"
	case 35:
		res = "<9>"
	case 40:
		res = "<10>"
	case 45:
		res = "<11>"
	case 50:
		res = "<12>"
	case 55:
		res = "<13>"
	case 60:
		res = "<14>"
	case 65:
		res = "<15>"
	case 70:
		res = "<16>"
	case 75:
		res = "<17>"
	case 80:
		res = "<18>"
	case 85:
		res = "<19>"
	case 90:
		res = "<20>"
	}
	return res
}

func test40withDefault(x byte) string {
	switch x {
	case 0:
		return "<0>"
	case 2:
		return "<1>"
	case 4:
		return "<2>"
	case 6:
		return "<3>"
	case 8:
		return "<4>"
	case 10:
		return "<5>"
	case 12:
		return "<6>"
	case 14:
		return "<7>"
	case 16:
		return "<8>"
	case 18:
		return "<9>"
	case 20:
		return "<10>"
	case 22:
		return "<11>"
	case 24:
		return "<12>"
	case 26:
		return "<13>"
	case 28:
		return "<14>"
	case 30:
		return "<15>"
	case 32:
		return "<16>"
	case 34:
		return "<17>"
	case 36:
		return "<18>"
	case 38:
		return "<19>"
	case 40:
		return "<20>"
	case 42:
		return "<21>"
	case 44:
		return "<22>"
	case 46:
		return "<23>"
	case 48:
		return "<24>"
	case 50:
		return "<25>"
	case 52:
		return "<26>"
	case 54:
		return "<27>"
	case 56:
		return "<28>"
	case 58:
		return "<29>"
	case 60:
		return "<30>"
	case 62:
		return "<31>"
	case 64:
		return "<32>"
	case 66:
		return "<33>"
	case 68:
		return "<34>"
	case 70:
		return "<35>"
	case 72:
		return "<36>"
	case 74:
		return "<37>"
	case 76:
		return "<38>"
	case 78:
		return "<39>"
	default:
		return "?"
	}
}

func test40noDefault(x byte) string {
	res := "?"
	switch x {
	case 0:
		res = "<0>"
	case 2:
		res = "<1>"
	case 4:
		res = "<2>"
	case 6:
		res = "<3>"
	case 8:
		res = "<4>"
	case 10:
		res = "<5>"
	case 12:
		res = "<6>"
	case 14:
		res = "<7>"
	case 16:
		res = "<8>"
	case 18:
		res = "<9>"
	case 20:
		res = "<10>"
	case 22:
		res = "<11>"
	case 24:
		res = "<12>"
	case 26:
		res = "<13>"
	case 28:
		res = "<14>"
	case 30:
		res = "<15>"
	case 32:
		res = "<16>"
	case 34:
		res = "<17>"
	case 36:
		res = "<18>"
	case 38:
		res = "<19>"
	case 40:
		res = "<20>"
	case 42:
		res = "<21>"
	case 44:
		res = "<22>"
	case 46:
		res = "<23>"
	case 48:
		res = "<24>"
	case 50:
		res = "<25>"
	case 52:
		res = "<26>"
	case 54:
		res = "<27>"
	case 56:
		res = "<28>"
	case 58:
		res = "<29>"
	case 60:
		res = "<30>"
	case 62:
		res = "<31>"
	case 64:
		res = "<32>"
	case 66:
		res = "<33>"
	case 68:
		res = "<34>"
	case 70:
		res = "<35>"
	case 72:
		res = "<36>"
	case 74:
		res = "<37>"
	case 76:
		res = "<38>"
	case 78:
		res = "<39>"
	}
	return res
}

func main() {
	b := byte(0)
	for i := -5; i < 50; i++ {
		b = byte(i)
		println(test3withDefault(b))
		println(test5noDefault(b))
		println(test5withDefault(b))
		println(test10noDefault(b))
		println(test10withDefault(b))
	}
	for i2 := -25; i2 < 100; i2++ {
		b = byte(i2)
		println(test21withDefault(b))
		println(test21noDefault(b))
		println(test40withDefault(b))
		println(test40noDefault(b))
	}
}
