//go:build main
// +build main

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"strings"
	"text/template"
)

var opcodePrototypes = []opcodeProto{
	{"LoadScalarConst", "op dst:wslot value:scalarindex"},
	{"LoadStrConst", "op dst:wslot value:strindex"},

	{"MoveScalar", "op dst:wslot src:rslot"},
	{"MoveStr", "op dst:wslot src:rslot"},
	{"MoveInterface", "op dst:wslot src:rslot"},
	{"MoveResult2", "op dst:wslot"},

	{"Not", "op dst:wslot x:rslot"},
	{"IsNil", "op dst:wslot x:rslot"},
	{"IsNotNil", "op dst:wslot x:rslot"},
	{"IsNilInterface", "op dst:wslot x:rslot"},
	{"IsNotNilInterface", "op dst:wslot x:rslot"},

	{"StrLen", "op dst:wslot str:rslot"},
	{"StrSlice", "op dst:wslot str:rslot from:rslot to:rslot"},
	{"StrSliceFrom", "op dst:wslot str:rslot from:rslot"},
	{"StrSliceTo", "op dst:wslot str:rslot to:rslot"},

	{"StrIndex", "op dst:wslot str:rslot index:rslot"},

	{"Concat", "op dst:wslot s1:rslot s2:rslot"},

	{"StrEq", "op dst:wslot s1:rslot s2:rslot"},
	{"StrNotEq", "op dst:wslot s1:rslot s2:rslot"},

	{"IntNeg", "op dst:wslot x:rslot"},

	{"ScalarEq", "op dst:wslot x:rslot y:rslot"},
	{"ScalarNotEq", "op dst:wslot x:rslot y:rslot"},

	{"IntGt", "op dst:wslot x:rslot y:rslot"},
	{"IntGtEq", "op dst:wslot x:rslot y:rslot"},
	{"IntLt", "op dst:wslot x:rslot y:rslot"},
	{"IntLtEq", "op dst:wslot x:rslot y:rslot"},

	{"IntAdd", "op dst:wslot x:rslot y:rslot"},
	{"IntSub", "op dst:wslot x:rslot y:rslot"},
	{"IntXor", "op dst:wslot x:rslot y:rslot"},
	{"IntMul", "op dst:wslot x:rslot y:rslot"},
	{"IntDiv", "op dst:wslot x:rslot y:rslot"},

	{"IntInc", "op x:rwslot"},
	{"IntDec", "op x:rwslot"},

	{"Jump", "op offset:offset"},
	{"JumpFalse", "op offset:offset cond:rslot"},
	{"JumpTrue", "op offset:offset cond:rslot"},

	{"Call", "op dst:wslot fn:funcid"},
	{"CallRecur", "op dst:wslot"},
	{"CallVoid", "op fn:funcid"},
	{"CallNative", "op dst:wslot fn:nativefuncid"},
	{"CallVoidNative", "op fn:nativefuncid"},

	{"PushVariadicBoolArg", "op x:rslot"},
	{"PushVariadicScalarArg", "op x:rslot"},
	{"PushVariadicStrArg", "op x:rslot"},
	{"PushVariadicInterfaceArg", "op x:rslot"},
	{"VariadicReset", "op"},

	{"ReturnVoid", "op"},
	{"ReturnFalse", "op"},
	{"ReturnTrue", "op"},
	{"ReturnStr", "op x:rslot"},
	{"ReturnScalar", "op x:rslot"},
	{"ReturnInterface", "op x:rslot"},
}

type opcodeProto struct {
	name string
	enc  string
}

type encodingInfo struct {
	width  int
	parts  int
	encdoc string
	args   string
	flags  string
}

type opcodeInfo struct {
	Opcode    byte
	Name      string
	Enc       string
	EncString string
	Width     int
	Flags     string
	Args      string
}

const stackUnchanged = ""

var fileTemplate = template.Must(template.New("opcodes.go").Parse(`// Code generated "gen_opcodes.go"; DO NOT EDIT.

package bytecode

const (
	OpInvalid Op = 0
{{ range .Opcodes }}
	// Encoding: {{.EncString}}
	Op{{ .Name }} Op = {{.Opcode}}
{{ end -}}
)

var opcodeInfoTable = [256]OpcodeInfo{
	OpInvalid: {Width: 1},

{{ range .Opcodes -}}
	Op{{.Name}}: {
		Width: {{.Width}},
		Flags: {{.Flags}},
		Args: []Argument{ {{.Args}} },
	},
{{ end }}
}
`))

func main() {
	opcodes := make([]opcodeInfo, len(opcodePrototypes))
	for i, proto := range opcodePrototypes {
		opcode := byte(i + 1)
		encInfo := decodeEnc(proto.enc)
		var encString string
		if encInfo.encdoc == "" {
			encString = fmt.Sprintf("0x%02x (width=%d)", opcode, encInfo.width)
		} else {
			encString = fmt.Sprintf("0x%02x %s (width=%d)",
				opcode, encInfo.encdoc, encInfo.width)
		}

		opcodes[i] = opcodeInfo{
			Opcode:    opcode,
			Name:      proto.name,
			Enc:       proto.enc,
			EncString: encString,
			Width:     encInfo.width,
			Flags:     encInfo.flags,
			Args:      encInfo.args,
		}
	}

	var buf bytes.Buffer
	err := fileTemplate.Execute(&buf, map[string]interface{}{
		"Opcodes": opcodes,
	})
	if err != nil {
		log.Panicf("execute template: %v", err)
	}
	writeFile("opcodes.gen.go", buf.Bytes())
}

func decodeEnc(enc string) encodingInfo {
	fields := strings.Fields(enc)
	width := 1 // opcode is uint8

	opfield := fields[0]
	if opfield != "op" {
		panic(fmt.Sprintf("parse %s: expected 'op', found '%s'", opfield))
	}

	argfields := fields[1:]

	var encdocParts []string
	var argList []string
	var argFlagList []string
	hasDst := false
	argOffset := 1
	for i, f := range argfields {
		argFlagList = argFlagList[:0]
		parts := strings.Split(f, ":")
		var typ string
		if len(parts) == 2 {
			typ = parts[1]
		} else {
			panic(fmt.Sprintf("parse %s: can't decode %s field: expected 2 parts", enc, f))
		}
		argName := parts[0]
		argType := ""
		encType := ""
		argWidth := 0

		switch typ {
		case "wslot", "rwslot":
			if i != 0 {
				panic(fmt.Sprintf("parse %s: dst arg at i=%d", enc, i))
			}
			if typ == "wslot" {
				argFlagList = append(argFlagList, "FlagIsWrite")
			} else {
				argFlagList = append(argFlagList, "FlagIsWrite", "FlagIsRead")
			}
			hasDst = true
			argType = "ArgSlot"
			encType = "u8"
			argWidth = 1
		case "rslot":
			argFlagList = append(argFlagList, "FlagIsRead")
			argType = "ArgSlot"
			encType = "u8"
			argWidth = 1
		case "strindex":
			argType = "ArgStrConst"
			encType = "u8"
			argWidth = 1
		case "scalarindex":
			argType = "ArgScalarConst"
			encType = "u8"
			argWidth = 1
		case "offset":
			argType = "ArgOffset"
			encType = "i16"
			argWidth = 2
		case "funcid":
			argType = "ArgFuncID"
			encType = "u16"
			argWidth = 2
		case "nativefuncid":
			argType = "ArgNativeFuncID"
			encType = "u16"
			argWidth = 2
		default:
			panic(fmt.Sprintf("unknown op argument type: %s", typ))
		}
		argFlags := "0"
		if len(argFlagList) != 0 {
			argFlags = strings.Join(argFlagList, " | ")
		}
		encdocParts = append(encdocParts, argName+":"+encType)
		argList = append(argList, fmt.Sprintf("{Name: %q, Kind: %s, Offset: %d, Flags: %s}", argName, argType, argOffset, argFlags))
		width += argWidth
		argOffset += argWidth
	}

	var flagList []string
	if hasDst {
		flagList = append(flagList, "FlagHasDst")
	}

	flagsString := "0"
	if len(flagList) != 0 {
		flagsString = strings.Join(flagList, " | ")
	}

	argsString := ""
	if len(argList) != 0 {
		argsString = "\n" + strings.Join(argList, ",\n")
	}

	return encodingInfo{
		width:  width,
		encdoc: strings.Join(encdocParts, " "),
		parts:  len(fields),
		flags:  flagsString,
		args:   argsString,
	}
}

func writeFile(filename string, data []byte) {
	pretty, err := format.Source(data)
	if err != nil {
		log.Panicf("gofmt: %v", err)
	}
	if err := ioutil.WriteFile(filename, pretty, 0666); err != nil {
		log.Panicf("write %s: %v", filename, err)
	}
}
