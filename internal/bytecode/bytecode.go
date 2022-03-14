package bytecode

//go:generate stringer -type=Op -trimprefix=Op
type Op byte

func (op Op) HasDst() bool {
	return opcodeInfoTable[op].Flags&FlagHasDst != 0
}

func (op Op) Width() int {
	return int(opcodeInfoTable[op].Width)
}

func (op Op) Args() []Argument {
	return opcodeInfoTable[op].Args
}

func (op Op) IsJump() bool {
	switch op {
	case OpJump, OpJumpNotZero, OpJumpZero, OpJumpTable:
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
	Offset uint8
	Flags  ArgumentFlags
}

func (a Argument) IsWriteSlot() bool { return a.Flags&FlagIsWrite != 0 }
func (a Argument) IsReadSlot() bool  { return a.Flags&FlagIsRead != 0 }

type ArgumentFlags uint8

const (
	FlagIsWrite ArgumentFlags = 1 << iota
	FlagIsRead
)

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
