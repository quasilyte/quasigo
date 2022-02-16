package bytecode

//go:generate stringer -type=Op -trimprefix=Op
type Op byte

func (op Op) HasDst() bool {
	return opcodeInfoTable[op].Flags&FlagHasDst != 0
}

func (op Op) Args() []Argument {
	return opcodeInfoTable[op].Args
}

func (op Op) IsJump() bool {
	switch op {
	case OpJump, OpJumpTrue, OpJumpFalse:
		return true
	default:
		return false
	}
}

type OpcodeFlags uint8

const (
	FlagHasDst OpcodeFlags = 1 << iota
)

type ArgumentKind uint8

const (
	ArgSlot ArgumentKind = iota
	ArgStrConst
	ArgScalarConst
	ArgOffset
	ArgFuncID
	ArgNativeFuncID
)

type Argument struct {
	Name   string
	Kind   ArgumentKind
	Offset int
}

type OpcodeInfo struct {
	Width uint8
	Flags OpcodeFlags
	Args  []Argument
}

func Walk(code []byte, fn func(pc int, op Op)) {
	pc := 0
	for pc < len(code) {
		op := Op(code[pc])
		fn(pc, op)
		pc += int(opcodeInfoTable[op].Width)
	}
}
