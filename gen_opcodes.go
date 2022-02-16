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

	{"Concat", "op dst:wslot s1:rslot s2:rslot"},

	{"StrEq", "op dst:wslot s1:rslot s2:rslot"},
	{"StrNotEq", "op dst:wslot s1:rslot s2:rslot"},

	{"IntEq", "op dst:wslot x:rslot y:rslot"},
	{"IntNotEq", "op dst:wslot x:rslot y:rslot"},
	{"IntGt", "op dst:wslot x:rslot y:rslot"},
	{"IntGtEq", "op dst:wslot x:rslot y:rslot"},
	{"IntLt", "op dst:wslot x:rslot y:rslot"},
	{"IntLtEq", "op dst:wslot x:rslot y:rslot"},

	{"IntAdd", "op dst:wslot x:rslot y:rslot"},
	{"IntSub", "op dst:wslot x:rslot y:rslot"},
	{"IntMul", "op dst:wslot x:rslot y:rslot"},
	{"IntDiv", "op dst:wslot x:rslot y:rslot"},

	{"IntInc", "op x:rwslot"},
	{"IntDec", "op x:rwslot"},

	{"Jump", "op offset:offset"},
	{"JumpFalse", "op offset:offset cond:rslot"},
	{"JumpTrue", "op offset:offset cond:rslot"},

	{"Call", "op dst:wslot fn:funcid"},
	{"CallRecur", "op dst:wslot"},
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

package quasigo

const (
	opInvalid opcode = 0
{{ range .Opcodes }}
	// Encoding: {{.EncString}}
	op{{ .Name }} opcode = {{.Opcode}}
{{ end -}}
)

var opcodeInfoTable = [256]opcodeInfo{
	opInvalid: {width: 1},

{{ range .Opcodes -}}
	op{{.Name}}: {
		width: {{.Width}},
		flags: {{.Flags}},
		args: []opcodeArgument{ {{.Args}} },
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
	hasDst := false
	argOffset := 1
	for i, f := range argfields {
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
			hasDst = true
			argType = "argkindSlot"
			encType = "u8"
			argWidth = 1
		case "rslot":
			argType = "argkindSlot"
			encType = "u8"
			argWidth = 1
		case "strindex":
			argType = "argkindStrConst"
			encType = "u8"
			argWidth = 1
		case "scalarindex":
			argType = "argkindScalarConst"
			encType = "u8"
			argWidth = 1
		case "offset":
			argType = "argkindOffset"
			encType = "i16"
			argWidth = 2
		case "funcid":
			argType = "argkindFuncID"
			encType = "u16"
			argWidth = 2
		case "nativefuncid":
			argType = "argkindNativeFuncID"
			encType = "u16"
			argWidth = 2
		default:
			panic(fmt.Sprintf("unknown op argument type: %s", typ))
		}
		encdocParts = append(encdocParts, argName+":"+encType)
		argList = append(argList, fmt.Sprintf("{name: %q, kind: %s, offset: %d}", argName, argType, argOffset))
		width += argWidth
		argOffset += argWidth
	}

	var flagList []string
	if hasDst {
		flagList = append(flagList, "opflagHasDst")
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
