// Code generated "gen_opcodes.go"; DO NOT EDIT.

package bytecode

const (
	OpInvalid Op = 0

	// Encoding: 0x01 dst:u8 value:u8 (width=3)
	OpLoadScalarConst Op = 1

	// Encoding: 0x02 dst:u8 value:u8 (width=3)
	OpLoadStrConst Op = 2

	// Encoding: 0x03 dst:u8 src:u8 (width=3)
	OpMoveScalar Op = 3

	// Encoding: 0x04 dst:u8 src:u8 (width=3)
	OpMoveStr Op = 4

	// Encoding: 0x05 dst:u8 src:u8 (width=3)
	OpMoveInterface Op = 5

	// Encoding: 0x06 dst:u8 (width=2)
	OpMoveResult2 Op = 6

	// Encoding: 0x07 dst:u8 x:u8 (width=3)
	OpNot Op = 7

	// Encoding: 0x08 dst:u8 x:u8 (width=3)
	OpIsNil Op = 8

	// Encoding: 0x09 dst:u8 x:u8 (width=3)
	OpIsNotNil Op = 9

	// Encoding: 0x0a dst:u8 x:u8 (width=3)
	OpIsNilInterface Op = 10

	// Encoding: 0x0b dst:u8 x:u8 (width=3)
	OpIsNotNilInterface Op = 11

	// Encoding: 0x0c dst:u8 str:u8 (width=3)
	OpStrLen Op = 12

	// Encoding: 0x0d dst:u8 str:u8 from:u8 to:u8 (width=5)
	OpStrSlice Op = 13

	// Encoding: 0x0e dst:u8 str:u8 from:u8 (width=4)
	OpStrSliceFrom Op = 14

	// Encoding: 0x0f dst:u8 str:u8 to:u8 (width=4)
	OpStrSliceTo Op = 15

	// Encoding: 0x10 dst:u8 s1:u8 s2:u8 (width=4)
	OpConcat Op = 16

	// Encoding: 0x11 dst:u8 s1:u8 s2:u8 (width=4)
	OpStrEq Op = 17

	// Encoding: 0x12 dst:u8 s1:u8 s2:u8 (width=4)
	OpStrNotEq Op = 18

	// Encoding: 0x13 dst:u8 x:u8 y:u8 (width=4)
	OpIntEq Op = 19

	// Encoding: 0x14 dst:u8 x:u8 y:u8 (width=4)
	OpIntNotEq Op = 20

	// Encoding: 0x15 dst:u8 x:u8 y:u8 (width=4)
	OpIntGt Op = 21

	// Encoding: 0x16 dst:u8 x:u8 y:u8 (width=4)
	OpIntGtEq Op = 22

	// Encoding: 0x17 dst:u8 x:u8 y:u8 (width=4)
	OpIntLt Op = 23

	// Encoding: 0x18 dst:u8 x:u8 y:u8 (width=4)
	OpIntLtEq Op = 24

	// Encoding: 0x19 dst:u8 x:u8 y:u8 (width=4)
	OpIntAdd Op = 25

	// Encoding: 0x1a dst:u8 x:u8 y:u8 (width=4)
	OpIntSub Op = 26

	// Encoding: 0x1b dst:u8 x:u8 y:u8 (width=4)
	OpIntMul Op = 27

	// Encoding: 0x1c dst:u8 x:u8 y:u8 (width=4)
	OpIntDiv Op = 28

	// Encoding: 0x1d x:u8 (width=2)
	OpIntInc Op = 29

	// Encoding: 0x1e x:u8 (width=2)
	OpIntDec Op = 30

	// Encoding: 0x1f offset:i16 (width=3)
	OpJump Op = 31

	// Encoding: 0x20 offset:i16 cond:u8 (width=4)
	OpJumpFalse Op = 32

	// Encoding: 0x21 offset:i16 cond:u8 (width=4)
	OpJumpTrue Op = 33

	// Encoding: 0x22 dst:u8 fn:u16 (width=4)
	OpCall Op = 34

	// Encoding: 0x23 dst:u8 (width=2)
	OpCallRecur Op = 35

	// Encoding: 0x24 dst:u8 fn:u16 (width=4)
	OpCallNative Op = 36

	// Encoding: 0x25 fn:u16 (width=3)
	OpCallVoidNative Op = 37

	// Encoding: 0x26 x:u8 (width=2)
	OpPushVariadicBoolArg Op = 38

	// Encoding: 0x27 x:u8 (width=2)
	OpPushVariadicScalarArg Op = 39

	// Encoding: 0x28 x:u8 (width=2)
	OpPushVariadicStrArg Op = 40

	// Encoding: 0x29 x:u8 (width=2)
	OpPushVariadicInterfaceArg Op = 41

	// Encoding: 0x2a (width=1)
	OpVariadicReset Op = 42

	// Encoding: 0x2b (width=1)
	OpReturnVoid Op = 43

	// Encoding: 0x2c (width=1)
	OpReturnFalse Op = 44

	// Encoding: 0x2d (width=1)
	OpReturnTrue Op = 45

	// Encoding: 0x2e x:u8 (width=2)
	OpReturnStr Op = 46

	// Encoding: 0x2f x:u8 (width=2)
	OpReturnScalar Op = 47

	// Encoding: 0x30 x:u8 (width=2)
	OpReturnInterface Op = 48
)

var opcodeInfoTable = [256]OpcodeInfo{
	OpInvalid: {Width: 1},

	OpLoadScalarConst: {
		Width: 3,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "value", Kind: ArgScalarConst, Offset: 2}},
	},
	OpLoadStrConst: {
		Width: 3,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "value", Kind: ArgStrConst, Offset: 2}},
	},
	OpMoveScalar: {
		Width: 3,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "src", Kind: ArgSlot, Offset: 2}},
	},
	OpMoveStr: {
		Width: 3,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "src", Kind: ArgSlot, Offset: 2}},
	},
	OpMoveInterface: {
		Width: 3,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "src", Kind: ArgSlot, Offset: 2}},
	},
	OpMoveResult2: {
		Width: 2,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1}},
	},
	OpNot: {
		Width: 3,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "x", Kind: ArgSlot, Offset: 2}},
	},
	OpIsNil: {
		Width: 3,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "x", Kind: ArgSlot, Offset: 2}},
	},
	OpIsNotNil: {
		Width: 3,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "x", Kind: ArgSlot, Offset: 2}},
	},
	OpIsNilInterface: {
		Width: 3,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "x", Kind: ArgSlot, Offset: 2}},
	},
	OpIsNotNilInterface: {
		Width: 3,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "x", Kind: ArgSlot, Offset: 2}},
	},
	OpStrLen: {
		Width: 3,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "str", Kind: ArgSlot, Offset: 2}},
	},
	OpStrSlice: {
		Width: 5,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "str", Kind: ArgSlot, Offset: 2},
			{Name: "from", Kind: ArgSlot, Offset: 3},
			{Name: "to", Kind: ArgSlot, Offset: 4}},
	},
	OpStrSliceFrom: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "str", Kind: ArgSlot, Offset: 2},
			{Name: "from", Kind: ArgSlot, Offset: 3}},
	},
	OpStrSliceTo: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "str", Kind: ArgSlot, Offset: 2},
			{Name: "to", Kind: ArgSlot, Offset: 3}},
	},
	OpConcat: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "s1", Kind: ArgSlot, Offset: 2},
			{Name: "s2", Kind: ArgSlot, Offset: 3}},
	},
	OpStrEq: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "s1", Kind: ArgSlot, Offset: 2},
			{Name: "s2", Kind: ArgSlot, Offset: 3}},
	},
	OpStrNotEq: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "s1", Kind: ArgSlot, Offset: 2},
			{Name: "s2", Kind: ArgSlot, Offset: 3}},
	},
	OpIntEq: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "x", Kind: ArgSlot, Offset: 2},
			{Name: "y", Kind: ArgSlot, Offset: 3}},
	},
	OpIntNotEq: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "x", Kind: ArgSlot, Offset: 2},
			{Name: "y", Kind: ArgSlot, Offset: 3}},
	},
	OpIntGt: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "x", Kind: ArgSlot, Offset: 2},
			{Name: "y", Kind: ArgSlot, Offset: 3}},
	},
	OpIntGtEq: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "x", Kind: ArgSlot, Offset: 2},
			{Name: "y", Kind: ArgSlot, Offset: 3}},
	},
	OpIntLt: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "x", Kind: ArgSlot, Offset: 2},
			{Name: "y", Kind: ArgSlot, Offset: 3}},
	},
	OpIntLtEq: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "x", Kind: ArgSlot, Offset: 2},
			{Name: "y", Kind: ArgSlot, Offset: 3}},
	},
	OpIntAdd: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "x", Kind: ArgSlot, Offset: 2},
			{Name: "y", Kind: ArgSlot, Offset: 3}},
	},
	OpIntSub: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "x", Kind: ArgSlot, Offset: 2},
			{Name: "y", Kind: ArgSlot, Offset: 3}},
	},
	OpIntMul: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "x", Kind: ArgSlot, Offset: 2},
			{Name: "y", Kind: ArgSlot, Offset: 3}},
	},
	OpIntDiv: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "x", Kind: ArgSlot, Offset: 2},
			{Name: "y", Kind: ArgSlot, Offset: 3}},
	},
	OpIntInc: {
		Width: 2,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "x", Kind: ArgSlot, Offset: 1}},
	},
	OpIntDec: {
		Width: 2,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "x", Kind: ArgSlot, Offset: 1}},
	},
	OpJump: {
		Width: 3,
		Flags: 0,
		Args: []Argument{
			{Name: "offset", Kind: ArgOffset, Offset: 1}},
	},
	OpJumpFalse: {
		Width: 4,
		Flags: 0,
		Args: []Argument{
			{Name: "offset", Kind: ArgOffset, Offset: 1},
			{Name: "cond", Kind: ArgSlot, Offset: 3}},
	},
	OpJumpTrue: {
		Width: 4,
		Flags: 0,
		Args: []Argument{
			{Name: "offset", Kind: ArgOffset, Offset: 1},
			{Name: "cond", Kind: ArgSlot, Offset: 3}},
	},
	OpCall: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "fn", Kind: ArgFuncID, Offset: 2}},
	},
	OpCallRecur: {
		Width: 2,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1}},
	},
	OpCallNative: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1},
			{Name: "fn", Kind: ArgNativeFuncID, Offset: 2}},
	},
	OpCallVoidNative: {
		Width: 3,
		Flags: 0,
		Args: []Argument{
			{Name: "fn", Kind: ArgNativeFuncID, Offset: 1}},
	},
	OpPushVariadicBoolArg: {
		Width: 2,
		Flags: 0,
		Args: []Argument{
			{Name: "x", Kind: ArgSlot, Offset: 1}},
	},
	OpPushVariadicScalarArg: {
		Width: 2,
		Flags: 0,
		Args: []Argument{
			{Name: "x", Kind: ArgSlot, Offset: 1}},
	},
	OpPushVariadicStrArg: {
		Width: 2,
		Flags: 0,
		Args: []Argument{
			{Name: "x", Kind: ArgSlot, Offset: 1}},
	},
	OpPushVariadicInterfaceArg: {
		Width: 2,
		Flags: 0,
		Args: []Argument{
			{Name: "x", Kind: ArgSlot, Offset: 1}},
	},
	OpVariadicReset: {
		Width: 1,
		Flags: 0,
		Args:  []Argument{},
	},
	OpReturnVoid: {
		Width: 1,
		Flags: 0,
		Args:  []Argument{},
	},
	OpReturnFalse: {
		Width: 1,
		Flags: 0,
		Args:  []Argument{},
	},
	OpReturnTrue: {
		Width: 1,
		Flags: 0,
		Args:  []Argument{},
	},
	OpReturnStr: {
		Width: 2,
		Flags: 0,
		Args: []Argument{
			{Name: "x", Kind: ArgSlot, Offset: 1}},
	},
	OpReturnScalar: {
		Width: 2,
		Flags: 0,
		Args: []Argument{
			{Name: "x", Kind: ArgSlot, Offset: 1}},
	},
	OpReturnInterface: {
		Width: 2,
		Flags: 0,
		Args: []Argument{
			{Name: "x", Kind: ArgSlot, Offset: 1}},
	},
}
