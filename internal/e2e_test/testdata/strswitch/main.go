package main

import (
	"fmt"
)

func test3withDefault(x string) int {
	switch x {
	case "10":
		return 1
	case "20":
		return 2
	case "30":
		return 3
	default:
		return -1
	}
}

func test5withDefault(x string) int {
	switch x {
	case "75":
		return 1
	case "60":
		return 2
	case "45":
		return 3
	case "30":
		return 4
	case "15":
		return 5
	default:
		return -1
	}
}

func test5noDefault(x string) int {
	res := -1
	switch x {
	case "75":
		res = 1
	case "60":
		res = 2
	case "45":
		res = 3
	case "30":
		res = 4
	case "15":
		res = 5
	}
	return res
}

// main.test10withDefault code=204 frame=96 (4 slots: 1 args, 1 locals, 2 temps)
//   Move auto0 = x
//   LoadStrConst temp1 = "0"
//   StrLt temp0 = auto0 temp1
//   JumpNotZero L0 temp0
//   LoadStrConst temp1 = "9"
//   StrGt temp0 = auto0 temp1
//   JumpNotZero L0 temp0
//   LoadStrConst temp1 = "4"
//   StrEq temp0 = auto0 temp1
//   JumpNotZero L1 temp0
//   StrGt temp0 = auto0 temp1
//   JumpNotZero L2 temp0
//   LoadStrConst temp1 = "0"
//   StrEq temp0 = auto0 temp1
//   JumpNotZero L3 temp0
//   LoadStrConst temp1 = "1"
//   StrEq temp0 = auto0 temp1
//   JumpNotZero L4 temp0
//   LoadStrConst temp1 = "2"
//   StrEq temp0 = auto0 temp1
//   JumpNotZero L5 temp0
//   LoadStrConst temp1 = "3"
//   StrEq temp0 = auto0 temp1
//   JumpNotZero L6 temp0
//   Jump L0
// L2:
//   LoadStrConst temp1 = "5"
//   StrEq temp0 = auto0 temp1
//   JumpNotZero L7 temp0
//   LoadStrConst temp1 = "6"
//   StrEq temp0 = auto0 temp1
//   JumpNotZero L8 temp0
//   LoadStrConst temp1 = "7"
//   StrEq temp0 = auto0 temp1
//   JumpNotZero L9 temp0
//   LoadStrConst temp1 = "8"
//   StrEq temp0 = auto0 temp1
//   JumpNotZero L10 temp0
//   LoadStrConst temp1 = "9"
//   StrEq temp0 = auto0 temp1
//   JumpNotZero L11 temp0
//   Jump L0
// L3:
//   LoadScalarConst temp0 = 0
//   ReturnScalar temp0
// L4:
//   LoadScalarConst temp0 = 1
//   ReturnScalar temp0
// L5:
//   LoadScalarConst temp0 = 2
//   ReturnScalar temp0
// L6:
//   LoadScalarConst temp0 = 3
//   ReturnScalar temp0
// L1:
//   LoadScalarConst temp0 = 4
//   ReturnScalar temp0
// L7:
//   LoadScalarConst temp0 = 5
//   ReturnScalar temp0
// L8:
//   LoadScalarConst temp0 = 6
//   ReturnScalar temp0
// L9:
//   LoadScalarConst temp0 = 7
//   ReturnScalar temp0
// L10:
//   LoadScalarConst temp0 = 8
//   ReturnScalar temp0
// L11:
//   LoadScalarConst temp0 = 9
//   ReturnScalar temp0
// L0:
//   LoadScalarConst temp0 = -1
//   ReturnScalar temp0
func test10withDefault(x string) int {
	switch x {
	case "0":
		return 0
	case "1":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	case "4":
		return 4
	case "5":
		return 5
	case "6":
		return 6
	case "7":
		return 7
	case "8":
		return 8
	case "9":
		return 9
	default:
		return -1
	}
}

func test10noDefault(x string) int {
	res := -1
	switch x {
	case "0":
		res = 0
	case "1":
		res = 1
	case "2":
		res = 2
	case "3":
		res = 3
	case "4":
		res = 4
	case "5":
		res = 5
	case "6":
		res = 6
	case "7":
		res = 7
	case "8":
		res = 8
	case "9":
		res = 9
	}
	return res
}

func test21withDefault(x string) int {
	switch x {
	case "-10":
		return 0
	case "-5":
		return 1
	case "0":
		return 2
	case "5":
		return 3
	case "10":
		return 4
	case "15":
		return 5
	case "20":
		return 6
	case "25":
		return 7
	case "30":
		return 8
	case "35":
		return 9
	case "40":
		return 10
	case "45":
		return 11
	case "50":
		return 12
	case "55":
		return 13
	case "60":
		return 14
	case "65":
		return 15
	case "70":
		return 16
	case "75":
		return 17
	case "80":
		return 18
	case "85":
		return 19
	case "90":
		return 20
	default:
		return -1
	}
}

func test21noDefault(x string) int {
	res := -1
	switch x {
	case "-10":
		res = 0
	case "-5":
		res = 1
	case "0":
		res = 2
	case "5":
		res = 3
	case "10":
		res = 4
	case "15":
		res = 5
	case "20":
		res = 6
	case "25":
		res = 7
	case "30":
		res = 8
	case "35":
		res = 9
	case "40":
		res = 10
	case "45":
		res = 11
	case "50":
		res = 12
	case "55":
		res = 13
	case "60":
		res = 14
	case "65":
		res = 15
	case "70":
		res = 16
	case "75":
		res = 17
	case "80":
		res = 18
	case "85":
		res = 19
	case "90":
		res = 20
	}
	return res
}

func test40withDefault(x string) int {
	switch x {
	case "0":
		return 0
	case "2":
		return 1
	case "4":
		return 2
	case "6":
		return 3
	case "8":
		return 4
	case "10":
		return 5
	case "12":
		return 6
	case "14":
		return 7
	case "16":
		return 8
	case "18":
		return 9
	case "20":
		return 10
	case "22":
		return 11
	case "24":
		return 12
	case "26":
		return 13
	case "28":
		return 14
	case "30":
		return 15
	case "32":
		return 16
	case "34":
		return 17
	case "36":
		return 18
	case "38":
		return 19
	case "40":
		return 20
	case "42":
		return 21
	case "44":
		return 22
	case "46":
		return 23
	case "48":
		return 24
	case "50":
		return 25
	case "52":
		return 26
	case "54":
		return 27
	case "56":
		return 28
	case "58":
		return 29
	case "60":
		return 30
	case "62":
		return 31
	case "64":
		return 32
	case "66":
		return 33
	case "68":
		return 34
	case "70":
		return 35
	case "72":
		return 36
	case "74":
		return 37
	case "76":
		return 38
	case "78":
		return 39
	default:
		return -1
	}
}

func test40noDefault(x string) int {
	res := -1
	switch x {
	case "0":
		res = 0
	case "2":
		res = 1
	case "4":
		res = 2
	case "6":
		res = 3
	case "8":
		res = 4
	case "10":
		res = 5
	case "12":
		res = 6
	case "14":
		res = 7
	case "16":
		res = 8
	case "18":
		res = 9
	case "20":
		res = 10
	case "22":
		res = 11
	case "24":
		res = 12
	case "26":
		res = 13
	case "28":
		res = 14
	case "30":
		res = 15
	case "32":
		res = 16
	case "34":
		res = 17
	case "36":
		res = 18
	case "38":
		res = 19
	case "40":
		res = 20
	case "42":
		res = 21
	case "44":
		res = 22
	case "46":
		res = 23
	case "48":
		res = 24
	case "50":
		res = 25
	case "52":
		res = 26
	case "54":
		res = 27
	case "56":
		res = 28
	case "58":
		res = 29
	case "60":
		res = 30
	case "62":
		res = 31
	case "64":
		res = 32
	case "66":
		res = 33
	case "68":
		res = 34
	case "70":
		res = 35
	case "72":
		res = 36
	case "74":
		res = 37
	case "76":
		res = 38
	case "78":
		res = 39
	}
	return res
}

func main() {
	for i := -5; i < 50; i++ {
		s := fmt.Sprintf("%d", i)
		println(test3withDefault(s))
		println(test5noDefault(s))
		println(test5withDefault(s))
		println(test10noDefault(s))
		println(test10withDefault(s))
	}
	for i2 := -25; i2 < 100; i2++ {
		s2 := fmt.Sprintf("%d", i2)
		println(test21withDefault(s2))
		println(test21noDefault(s2))
		println(test40withDefault(s2))
		println(test40noDefault(s2))
	}
}
