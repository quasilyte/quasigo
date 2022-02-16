// Code generated "gen_opcodes.go"; DO NOT EDIT.

package quasigo

const (
	opInvalid opcode = 0

	// Encoding: 0x01 dst:u8 value:u8 (width=3)
	opLoadScalarConst opcode = 1

	// Encoding: 0x02 dst:u8 value:u8 (width=3)
	opLoadStrConst opcode = 2

	// Encoding: 0x03 dst:u8 src:u8 (width=3)
	opMoveScalar opcode = 3

	// Encoding: 0x04 dst:u8 src:u8 (width=3)
	opMoveStr opcode = 4

	// Encoding: 0x05 dst:u8 src:u8 (width=3)
	opMoveInterface opcode = 5

	// Encoding: 0x06 dst:u8 (width=2)
	opMoveResult2 opcode = 6

	// Encoding: 0x07 dst:u8 x:u8 (width=3)
	opNot opcode = 7

	// Encoding: 0x08 dst:u8 x:u8 (width=3)
	opIsNil opcode = 8

	// Encoding: 0x09 dst:u8 x:u8 (width=3)
	opIsNotNil opcode = 9

	// Encoding: 0x0a dst:u8 x:u8 (width=3)
	opIsNilInterface opcode = 10

	// Encoding: 0x0b dst:u8 x:u8 (width=3)
	opIsNotNilInterface opcode = 11

	// Encoding: 0x0c dst:u8 str:u8 (width=3)
	opStrLen opcode = 12

	// Encoding: 0x0d dst:u8 str:u8 from:u8 to:u8 (width=5)
	opStrSlice opcode = 13

	// Encoding: 0x0e dst:u8 str:u8 from:u8 (width=4)
	opStrSliceFrom opcode = 14

	// Encoding: 0x0f dst:u8 str:u8 to:u8 (width=4)
	opStrSliceTo opcode = 15

	// Encoding: 0x10 dst:u8 s1:u8 s2:u8 (width=4)
	opConcat opcode = 16

	// Encoding: 0x11 dst:u8 s1:u8 s2:u8 (width=4)
	opStrEq opcode = 17

	// Encoding: 0x12 dst:u8 s1:u8 s2:u8 (width=4)
	opStrNotEq opcode = 18

	// Encoding: 0x13 dst:u8 x:u8 y:u8 (width=4)
	opIntEq opcode = 19

	// Encoding: 0x14 dst:u8 x:u8 y:u8 (width=4)
	opIntNotEq opcode = 20

	// Encoding: 0x15 dst:u8 x:u8 y:u8 (width=4)
	opIntGt opcode = 21

	// Encoding: 0x16 dst:u8 x:u8 y:u8 (width=4)
	opIntGtEq opcode = 22

	// Encoding: 0x17 dst:u8 x:u8 y:u8 (width=4)
	opIntLt opcode = 23

	// Encoding: 0x18 dst:u8 x:u8 y:u8 (width=4)
	opIntLtEq opcode = 24

	// Encoding: 0x19 dst:u8 x:u8 y:u8 (width=4)
	opIntAdd opcode = 25

	// Encoding: 0x1a dst:u8 x:u8 y:u8 (width=4)
	opIntSub opcode = 26

	// Encoding: 0x1b dst:u8 x:u8 y:u8 (width=4)
	opIntMul opcode = 27

	// Encoding: 0x1c dst:u8 x:u8 y:u8 (width=4)
	opIntDiv opcode = 28

	// Encoding: 0x1d x:u8 (width=2)
	opIntInc opcode = 29

	// Encoding: 0x1e x:u8 (width=2)
	opIntDec opcode = 30

	// Encoding: 0x1f offset:i16 (width=3)
	opJump opcode = 31

	// Encoding: 0x20 offset:i16 cond:u8 (width=4)
	opJumpFalse opcode = 32

	// Encoding: 0x21 offset:i16 cond:u8 (width=4)
	opJumpTrue opcode = 33

	// Encoding: 0x22 dst:u8 fn:u16 (width=4)
	opCall opcode = 34

	// Encoding: 0x23 dst:u8 (width=2)
	opCallRecur opcode = 35

	// Encoding: 0x24 dst:u8 fn:u16 (width=4)
	opCallNative opcode = 36

	// Encoding: 0x25 fn:u16 (width=3)
	opCallVoidNative opcode = 37

	// Encoding: 0x26 x:u8 (width=2)
	opPushVariadicBoolArg opcode = 38

	// Encoding: 0x27 x:u8 (width=2)
	opPushVariadicScalarArg opcode = 39

	// Encoding: 0x28 x:u8 (width=2)
	opPushVariadicStrArg opcode = 40

	// Encoding: 0x29 x:u8 (width=2)
	opPushVariadicInterfaceArg opcode = 41

	// Encoding: 0x2a (width=1)
	opVariadicReset opcode = 42

	// Encoding: 0x2b (width=1)
	opReturnVoid opcode = 43

	// Encoding: 0x2c (width=1)
	opReturnFalse opcode = 44

	// Encoding: 0x2d (width=1)
	opReturnTrue opcode = 45

	// Encoding: 0x2e x:u8 (width=2)
	opReturnStr opcode = 46

	// Encoding: 0x2f x:u8 (width=2)
	opReturnScalar opcode = 47

	// Encoding: 0x30 x:u8 (width=2)
	opReturnInterface opcode = 48
)

var opcodeInfoTable = [256]opcodeInfo{
	opInvalid: {width: 1},

	opLoadScalarConst: {
		width: 3,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "value", kind: argkindScalarConst, offset: 2}},
	},
	opLoadStrConst: {
		width: 3,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "value", kind: argkindStrConst, offset: 2}},
	},
	opMoveScalar: {
		width: 3,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "src", kind: argkindSlot, offset: 2}},
	},
	opMoveStr: {
		width: 3,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "src", kind: argkindSlot, offset: 2}},
	},
	opMoveInterface: {
		width: 3,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "src", kind: argkindSlot, offset: 2}},
	},
	opMoveResult2: {
		width: 2,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1}},
	},
	opNot: {
		width: 3,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "x", kind: argkindSlot, offset: 2}},
	},
	opIsNil: {
		width: 3,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "x", kind: argkindSlot, offset: 2}},
	},
	opIsNotNil: {
		width: 3,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "x", kind: argkindSlot, offset: 2}},
	},
	opIsNilInterface: {
		width: 3,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "x", kind: argkindSlot, offset: 2}},
	},
	opIsNotNilInterface: {
		width: 3,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "x", kind: argkindSlot, offset: 2}},
	},
	opStrLen: {
		width: 3,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "str", kind: argkindSlot, offset: 2}},
	},
	opStrSlice: {
		width: 5,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "str", kind: argkindSlot, offset: 2},
			{name: "from", kind: argkindSlot, offset: 3},
			{name: "to", kind: argkindSlot, offset: 4}},
	},
	opStrSliceFrom: {
		width: 4,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "str", kind: argkindSlot, offset: 2},
			{name: "from", kind: argkindSlot, offset: 3}},
	},
	opStrSliceTo: {
		width: 4,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "str", kind: argkindSlot, offset: 2},
			{name: "to", kind: argkindSlot, offset: 3}},
	},
	opConcat: {
		width: 4,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "s1", kind: argkindSlot, offset: 2},
			{name: "s2", kind: argkindSlot, offset: 3}},
	},
	opStrEq: {
		width: 4,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "s1", kind: argkindSlot, offset: 2},
			{name: "s2", kind: argkindSlot, offset: 3}},
	},
	opStrNotEq: {
		width: 4,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "s1", kind: argkindSlot, offset: 2},
			{name: "s2", kind: argkindSlot, offset: 3}},
	},
	opIntEq: {
		width: 4,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "x", kind: argkindSlot, offset: 2},
			{name: "y", kind: argkindSlot, offset: 3}},
	},
	opIntNotEq: {
		width: 4,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "x", kind: argkindSlot, offset: 2},
			{name: "y", kind: argkindSlot, offset: 3}},
	},
	opIntGt: {
		width: 4,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "x", kind: argkindSlot, offset: 2},
			{name: "y", kind: argkindSlot, offset: 3}},
	},
	opIntGtEq: {
		width: 4,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "x", kind: argkindSlot, offset: 2},
			{name: "y", kind: argkindSlot, offset: 3}},
	},
	opIntLt: {
		width: 4,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "x", kind: argkindSlot, offset: 2},
			{name: "y", kind: argkindSlot, offset: 3}},
	},
	opIntLtEq: {
		width: 4,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "x", kind: argkindSlot, offset: 2},
			{name: "y", kind: argkindSlot, offset: 3}},
	},
	opIntAdd: {
		width: 4,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "x", kind: argkindSlot, offset: 2},
			{name: "y", kind: argkindSlot, offset: 3}},
	},
	opIntSub: {
		width: 4,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "x", kind: argkindSlot, offset: 2},
			{name: "y", kind: argkindSlot, offset: 3}},
	},
	opIntMul: {
		width: 4,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "x", kind: argkindSlot, offset: 2},
			{name: "y", kind: argkindSlot, offset: 3}},
	},
	opIntDiv: {
		width: 4,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "x", kind: argkindSlot, offset: 2},
			{name: "y", kind: argkindSlot, offset: 3}},
	},
	opIntInc: {
		width: 2,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "x", kind: argkindSlot, offset: 1}},
	},
	opIntDec: {
		width: 2,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "x", kind: argkindSlot, offset: 1}},
	},
	opJump: {
		width: 3,
		flags: 0,
		args: []opcodeArgument{
			{name: "offset", kind: argkindOffset, offset: 1}},
	},
	opJumpFalse: {
		width: 4,
		flags: 0,
		args: []opcodeArgument{
			{name: "offset", kind: argkindOffset, offset: 1},
			{name: "cond", kind: argkindSlot, offset: 3}},
	},
	opJumpTrue: {
		width: 4,
		flags: 0,
		args: []opcodeArgument{
			{name: "offset", kind: argkindOffset, offset: 1},
			{name: "cond", kind: argkindSlot, offset: 3}},
	},
	opCall: {
		width: 4,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "fn", kind: argkindFuncID, offset: 2}},
	},
	opCallRecur: {
		width: 2,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1}},
	},
	opCallNative: {
		width: 4,
		flags: opflagHasDst,
		args: []opcodeArgument{
			{name: "dst", kind: argkindSlot, offset: 1},
			{name: "fn", kind: argkindNativeFuncID, offset: 2}},
	},
	opCallVoidNative: {
		width: 3,
		flags: 0,
		args: []opcodeArgument{
			{name: "fn", kind: argkindNativeFuncID, offset: 1}},
	},
	opPushVariadicBoolArg: {
		width: 2,
		flags: 0,
		args: []opcodeArgument{
			{name: "x", kind: argkindSlot, offset: 1}},
	},
	opPushVariadicScalarArg: {
		width: 2,
		flags: 0,
		args: []opcodeArgument{
			{name: "x", kind: argkindSlot, offset: 1}},
	},
	opPushVariadicStrArg: {
		width: 2,
		flags: 0,
		args: []opcodeArgument{
			{name: "x", kind: argkindSlot, offset: 1}},
	},
	opPushVariadicInterfaceArg: {
		width: 2,
		flags: 0,
		args: []opcodeArgument{
			{name: "x", kind: argkindSlot, offset: 1}},
	},
	opVariadicReset: {
		width: 1,
		flags: 0,
		args:  []opcodeArgument{},
	},
	opReturnVoid: {
		width: 1,
		flags: 0,
		args:  []opcodeArgument{},
	},
	opReturnFalse: {
		width: 1,
		flags: 0,
		args:  []opcodeArgument{},
	},
	opReturnTrue: {
		width: 1,
		flags: 0,
		args:  []opcodeArgument{},
	},
	opReturnStr: {
		width: 2,
		flags: 0,
		args: []opcodeArgument{
			{name: "x", kind: argkindSlot, offset: 1}},
	},
	opReturnScalar: {
		width: 2,
		flags: 0,
		args: []opcodeArgument{
			{name: "x", kind: argkindSlot, offset: 1}},
	},
	opReturnInterface: {
		width: 2,
		flags: 0,
		args: []opcodeArgument{
			{name: "x", kind: argkindSlot, offset: 1}},
	},
}
