package quasigo

import (
	"encoding/binary"
	"go/ast"
	"go/types"
)

func pickOp(cond bool, ifTrue, otherwise opcode) opcode {
	if cond {
		return ifTrue
	}
	return otherwise
}

func put16(code []byte, pos, value int) {
	binary.LittleEndian.PutUint16(code[pos:], uint16(value))
}

func typeIsScalar(typ types.Type) bool {
	basic, ok := typ.Underlying().(*types.Basic)
	if !ok {
		return false
	}
	switch basic.Kind() {
	case types.Int, types.UntypedInt, types.Bool, types.UntypedBool:
		return true
	default:
		return false
	}
}

func typeIsBool(typ types.Type) bool {
	basic, ok := typ.Underlying().(*types.Basic)
	if !ok {
		return false
	}
	return basic.Info()&types.IsBoolean != 0
}

func typeIsInt(typ types.Type) bool {
	basic, ok := typ.Underlying().(*types.Basic)
	if !ok {
		return false
	}
	switch basic.Kind() {
	case types.Int, types.UntypedInt:
		return true
	default:
		return false
	}
}

func typeIsPointer(typ types.Type) bool {
	_, ok := typ.Underlying().(*types.Pointer)
	return ok
}

func typeIsInterface(typ types.Type) bool {
	_, ok := typ.Underlying().(*types.Interface)
	return ok
}

func typeIsString(typ types.Type) bool {
	basic, ok := typ.Underlying().(*types.Basic)
	if !ok {
		return false
	}
	return basic.Info()&types.IsString != 0
}

func walkBytecode(code []byte, fn func(pc int, op opcode)) {
	pc := 0
	for pc < len(code) {
		op := opcode(code[pc])
		fn(pc, op)
		pc += int(opcodeInfoTable[op].width)
	}
}

func identName(n ast.Expr) string {
	id, ok := n.(*ast.Ident)
	if ok {
		return id.Name
	}
	return ""
}
