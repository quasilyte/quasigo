package qruntime

import "strings"

type Func struct {
	StrConstants    []string
	ScalarConstants []uint64

	Codeptr *byte
	Code    []byte

	FrameSize  int
	FrameSlots byte
	NumParams  byte
	NumTemps   byte

	CanInline bool

	Name string
}

type FuncKey struct {
	Qualifier string
	Name      string
}

func (k FuncKey) String() string {
	if k.Qualifier != "" {
		return k.Qualifier + "." + k.Name
	}
	return k.Name
}

// typeAliases maps the predeclared type aliases to the
// spelling that is used by the compiler in the instantiation names.
var typeAliases = map[string]string{
	"byte":        "uint8",
	"rune":        "int32",
	"interface{}": "any",
}

// canonicalInstantiationName replaces the type aliases inside
// the generic func instantiation name type args, so both
// `If[[]byte]` and `If[[]uint8]` spellings can be used.
// Names without type args are returned as is.
func canonicalInstantiationName(name string) string {
	start := strings.IndexByte(name, '[')
	if start == -1 {
		return name
	}

	var buf strings.Builder
	buf.WriteString(name[:start])
	for i := start; i < len(name); i++ {
		switch name[i] {
		case '[', ']', ',', '*':
			typeName := name[start:i]
			if alias, ok := typeAliases[typeName]; ok {
				typeName = alias
			}
			buf.WriteString(typeName)
			buf.WriteByte(name[i])
			start = i + 1
		}
	}
	buf.WriteString(name[start:])
	return buf.String()
}

func Call(env *EvalEnv, fn *Func) Slot {
	eval(env, fn, env.slotbase)
	return env.result
}
