// Code generated "gen_opcodes.go"; DO NOT EDIT.

package bytecode

const (
	OpInvalid Op = 0

	// Encoding: 0x01 dst:u8 value:u8 (width=3)
	OpLoadScalarConst Op = 1

	// Encoding: 0x02 dst:u8 value:u8 (width=3)
	OpLoadStrConst Op = 2

	// Encoding: 0x03 dst:u8 (width=2)
	OpZero Op = 3

	// Encoding: 0x04 dst:u8 src:u8 (width=3)
	OpMove Op = 4

	// Encoding: 0x05 dst:u8 src:u8 (width=3)
	OpMove8 Op = 5

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
	OpLen Op = 12

	// Encoding: 0x0d dst:u8 str:u8 (width=3)
	OpCap Op = 13

	// Encoding: 0x0e dst:u8 str:u8 from:u8 to:u8 (width=5)
	OpStrSlice Op = 14

	// Encoding: 0x0f dst:u8 str:u8 from:u8 (width=4)
	OpStrSliceFrom Op = 15

	// Encoding: 0x10 dst:u8 str:u8 to:u8 (width=4)
	OpStrSliceTo Op = 16

	// Encoding: 0x11 dst:u8 str:u8 index:u8 (width=4)
	OpStrIndex Op = 17

	// Encoding: 0x12 dst:u8 slice:u8 index:u8 (width=4)
	OpSliceIndexScalar8 Op = 18

	// Encoding: 0x13 dst:u8 slice:u8 index:u8 (width=4)
	OpSliceIndexScalar64 Op = 19

	// Encoding: 0x14 slice:u8 index:u8 value:u8 (width=4)
	OpSliceSetScalar8 Op = 20

	// Encoding: 0x15 slice:u8 index:u8 value:u8 (width=4)
	OpSliceSetScalar64 Op = 21

	// Encoding: 0x16 dst:u8 s1:u8 s2:u8 (width=4)
	OpConcat Op = 22

	// Encoding: 0x17 dst:u8 s1:u8 s2:u8 (width=4)
	OpStrEq Op = 23

	// Encoding: 0x18 dst:u8 s1:u8 s2:u8 (width=4)
	OpStrNotEq Op = 24

	// Encoding: 0x19 dst:u8 s1:u8 s2:u8 (width=4)
	OpStrGt Op = 25

	// Encoding: 0x1a dst:u8 s1:u8 s2:u8 (width=4)
	OpStrLt Op = 26

	// Encoding: 0x1b dst:u8 x:u8 (width=3)
	OpIntNeg Op = 27

	// Encoding: 0x1c dst:u8 x:u8 y:u8 (width=4)
	OpScalarEq Op = 28

	// Encoding: 0x1d dst:u8 x:u8 y:u8 (width=4)
	OpScalarNotEq Op = 29

	// Encoding: 0x1e dst:u8 x:u8 y:u8 (width=4)
	OpIntGt Op = 30

	// Encoding: 0x1f dst:u8 x:u8 y:u8 (width=4)
	OpIntGtEq Op = 31

	// Encoding: 0x20 dst:u8 x:u8 y:u8 (width=4)
	OpIntLt Op = 32

	// Encoding: 0x21 dst:u8 x:u8 y:u8 (width=4)
	OpIntLtEq Op = 33

	// Encoding: 0x22 dst:u8 x:u8 y:u8 (width=4)
	OpIntAdd8 Op = 34

	// Encoding: 0x23 dst:u8 x:u8 y:u8 (width=4)
	OpIntAdd64 Op = 35

	// Encoding: 0x24 dst:u8 x:u8 y:u8 (width=4)
	OpIntSub8 Op = 36

	// Encoding: 0x25 dst:u8 x:u8 y:u8 (width=4)
	OpIntSub64 Op = 37

	// Encoding: 0x26 dst:u8 x:u8 y:u8 (width=4)
	OpIntMul8 Op = 38

	// Encoding: 0x27 dst:u8 x:u8 y:u8 (width=4)
	OpIntMul64 Op = 39

	// Encoding: 0x28 dst:u8 x:u8 y:u8 (width=4)
	OpIntXor Op = 40

	// Encoding: 0x29 dst:u8 x:u8 y:u8 (width=4)
	OpIntDiv Op = 41

	// Encoding: 0x2a dst:u8 x:u8 y:u8 (width=4)
	OpIntMod Op = 42

	// Encoding: 0x2b x:u8 (width=2)
	OpIntInc Op = 43

	// Encoding: 0x2c x:u8 (width=2)
	OpIntDec Op = 44

	// Encoding: 0x2d offset:i16 (width=3)
	OpJump Op = 45

	// Encoding: 0x2e offset:i16 cond:u8 (width=4)
	OpJumpZero Op = 46

	// Encoding: 0x2f offset:i16 cond:u8 (width=4)
	OpJumpNotZero Op = 47

	// Encoding: 0x30 value:u8 (width=2)
	OpJumpTable Op = 48

	// Encoding: 0x31 dst:u8 fn:u16 (width=4)
	OpCall Op = 49

	// Encoding: 0x32 dst:u8 (width=2)
	OpCallRecur Op = 50

	// Encoding: 0x33 fn:u16 (width=3)
	OpCallVoid Op = 51

	// Encoding: 0x34 dst:u8 fn:u16 (width=4)
	OpCallNative Op = 52

	// Encoding: 0x35 fn:u16 (width=3)
	OpCallVoidNative Op = 53

	// Encoding: 0x36 x:u8 (width=2)
	OpPushVariadicBoolArg Op = 54

	// Encoding: 0x37 x:u8 (width=2)
	OpPushVariadicScalarArg Op = 55

	// Encoding: 0x38 x:u8 (width=2)
	OpPushVariadicStrArg Op = 56

	// Encoding: 0x39 x:u8 (width=2)
	OpPushVariadicInterfaceArg Op = 57

	// Encoding: 0x3a (width=1)
	OpVariadicReset Op = 58

	// Encoding: 0x3b (width=1)
	OpReturnVoid Op = 59

	// Encoding: 0x3c (width=1)
	OpReturnZero Op = 60

	// Encoding: 0x3d (width=1)
	OpReturnOne Op = 61

	// Encoding: 0x3e x:u8 (width=2)
	OpReturnStr Op = 62

	// Encoding: 0x3f x:u8 (width=2)
	OpReturnScalar Op = 63

	// Encoding: 0x40 x:u8 (width=2)
	OpReturn Op = 64
)

var opcodeInfoTable = [256]OpcodeInfo{
	OpInvalid: {Width: 1},

	OpLoadScalarConst: {
		Width: 3,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "value", Kind: ArgScalarConst, Offset: 2, Flags: 0}},
	},
	OpLoadStrConst: {
		Width: 3,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "value", Kind: ArgStrConst, Offset: 2, Flags: 0}},
	},
	OpZero: {
		Width: 2,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite}},
	},
	OpMove: {
		Width: 3,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "src", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead}},
	},
	OpMove8: {
		Width: 3,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "src", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead}},
	},
	OpMoveResult2: {
		Width: 2,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite}},
	},
	OpNot: {
		Width: 3,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "x", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead}},
	},
	OpIsNil: {
		Width: 3,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "x", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead}},
	},
	OpIsNotNil: {
		Width: 3,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "x", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead}},
	},
	OpIsNilInterface: {
		Width: 3,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "x", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead}},
	},
	OpIsNotNilInterface: {
		Width: 3,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "x", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead}},
	},
	OpLen: {
		Width: 3,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "str", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead}},
	},
	OpCap: {
		Width: 3,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "str", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead}},
	},
	OpStrSlice: {
		Width: 5,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "str", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "from", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead},
			{Name: "to", Kind: ArgSlot, Offset: 4, Flags: FlagIsRead}},
	},
	OpStrSliceFrom: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "str", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "from", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpStrSliceTo: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "str", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "to", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpStrIndex: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "str", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "index", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpSliceIndexScalar8: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "slice", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "index", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpSliceIndexScalar64: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "slice", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "index", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpSliceSetScalar8: {
		Width: 4,
		Flags: 0,
		Args: []Argument{
			{Name: "slice", Kind: ArgSlot, Offset: 1, Flags: FlagIsRead},
			{Name: "index", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "value", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpSliceSetScalar64: {
		Width: 4,
		Flags: 0,
		Args: []Argument{
			{Name: "slice", Kind: ArgSlot, Offset: 1, Flags: FlagIsRead},
			{Name: "index", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "value", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpConcat: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "s1", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "s2", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpStrEq: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "s1", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "s2", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpStrNotEq: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "s1", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "s2", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpStrGt: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "s1", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "s2", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpStrLt: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "s1", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "s2", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpIntNeg: {
		Width: 3,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "x", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead}},
	},
	OpScalarEq: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "x", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "y", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpScalarNotEq: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "x", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "y", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpIntGt: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "x", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "y", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpIntGtEq: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "x", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "y", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpIntLt: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "x", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "y", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpIntLtEq: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "x", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "y", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpIntAdd8: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "x", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "y", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpIntAdd64: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "x", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "y", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpIntSub8: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "x", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "y", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpIntSub64: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "x", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "y", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpIntMul8: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "x", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "y", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpIntMul64: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "x", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "y", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpIntXor: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "x", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "y", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpIntDiv: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "x", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "y", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpIntMod: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "x", Kind: ArgSlot, Offset: 2, Flags: FlagIsRead},
			{Name: "y", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpIntInc: {
		Width: 2,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "x", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite | FlagIsRead}},
	},
	OpIntDec: {
		Width: 2,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "x", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite | FlagIsRead}},
	},
	OpJump: {
		Width: 3,
		Flags: 0,
		Args: []Argument{
			{Name: "offset", Kind: ArgOffset, Offset: 1, Flags: 0}},
	},
	OpJumpZero: {
		Width: 4,
		Flags: 0,
		Args: []Argument{
			{Name: "offset", Kind: ArgOffset, Offset: 1, Flags: 0},
			{Name: "cond", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpJumpNotZero: {
		Width: 4,
		Flags: 0,
		Args: []Argument{
			{Name: "offset", Kind: ArgOffset, Offset: 1, Flags: 0},
			{Name: "cond", Kind: ArgSlot, Offset: 3, Flags: FlagIsRead}},
	},
	OpJumpTable: {
		Width: 2,
		Flags: 0,
		Args: []Argument{
			{Name: "value", Kind: ArgSlot, Offset: 1, Flags: FlagIsRead}},
	},
	OpCall: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "fn", Kind: ArgFuncID, Offset: 2, Flags: 0}},
	},
	OpCallRecur: {
		Width: 2,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite}},
	},
	OpCallVoid: {
		Width: 3,
		Flags: 0,
		Args: []Argument{
			{Name: "fn", Kind: ArgFuncID, Offset: 1, Flags: 0}},
	},
	OpCallNative: {
		Width: 4,
		Flags: FlagHasDst,
		Args: []Argument{
			{Name: "dst", Kind: ArgSlot, Offset: 1, Flags: FlagIsWrite},
			{Name: "fn", Kind: ArgNativeFuncID, Offset: 2, Flags: 0}},
	},
	OpCallVoidNative: {
		Width: 3,
		Flags: 0,
		Args: []Argument{
			{Name: "fn", Kind: ArgNativeFuncID, Offset: 1, Flags: 0}},
	},
	OpPushVariadicBoolArg: {
		Width: 2,
		Flags: 0,
		Args: []Argument{
			{Name: "x", Kind: ArgSlot, Offset: 1, Flags: FlagIsRead}},
	},
	OpPushVariadicScalarArg: {
		Width: 2,
		Flags: 0,
		Args: []Argument{
			{Name: "x", Kind: ArgSlot, Offset: 1, Flags: FlagIsRead}},
	},
	OpPushVariadicStrArg: {
		Width: 2,
		Flags: 0,
		Args: []Argument{
			{Name: "x", Kind: ArgSlot, Offset: 1, Flags: FlagIsRead}},
	},
	OpPushVariadicInterfaceArg: {
		Width: 2,
		Flags: 0,
		Args: []Argument{
			{Name: "x", Kind: ArgSlot, Offset: 1, Flags: FlagIsRead}},
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
	OpReturnZero: {
		Width: 1,
		Flags: 0,
		Args:  []Argument{},
	},
	OpReturnOne: {
		Width: 1,
		Flags: 0,
		Args:  []Argument{},
	},
	OpReturnStr: {
		Width: 2,
		Flags: 0,
		Args: []Argument{
			{Name: "x", Kind: ArgSlot, Offset: 1, Flags: FlagIsRead}},
	},
	OpReturnScalar: {
		Width: 2,
		Flags: 0,
		Args: []Argument{
			{Name: "x", Kind: ArgSlot, Offset: 1, Flags: FlagIsRead}},
	},
	OpReturn: {
		Width: 2,
		Flags: 0,
		Args: []Argument{
			{Name: "x", Kind: ArgSlot, Offset: 1, Flags: FlagIsRead}},
	},
}
