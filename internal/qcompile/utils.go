package qcompile

import (
	"encoding/binary"
	"go/ast"
	"go/types"

	"github.com/quasilyte/quasigo/internal/bytecode"
)

func pickOp(cond bool, ifTrue, otherwise bytecode.Op) bytecode.Op {
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
	case types.Int, types.UntypedInt, types.Bool, types.UntypedBool, types.Uint8:
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

func typeIsByte(typ types.Type) bool {
	basic, ok := typ.Underlying().(*types.Basic)
	if !ok {
		return false
	}
	return basic.Kind() == types.Uint8
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

func identName(n ast.Expr) string {
	id, ok := n.(*ast.Ident)
	if ok {
		return id.Name
	}
	return ""
}
