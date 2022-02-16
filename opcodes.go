package quasigo

//go:generate stringer -type=opcode -trimprefix=op
type opcode byte

func (op opcode) HasDst() bool {
	return opcodeInfoTable[op].flags&opflagHasDst != 0
}

func (op opcode) Args() []opcodeArgument {
	return opcodeInfoTable[op].args
}

func (op opcode) IsJump() bool {
	switch op {
	case opJump, opJumpTrue, opJumpFalse:
		return true
	default:
		return false
	}
}

type opcodeFlags uint8

const (
	opflagHasDst opcodeFlags = 1 << iota
)

type opcodeArgKind uint8

const (
	argkindSlot opcodeArgKind = iota
	argkindStrConst
	argkindScalarConst
	argkindOffset
	argkindFuncID
	argkindNativeFuncID
)

type opcodeArgument struct {
	name   string
	kind   opcodeArgKind
	offset int
}

type opcodeInfo struct {
	width uint8
	flags opcodeFlags
	args  []opcodeArgument
}
