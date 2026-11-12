package qruntime

import (
	"testing"
)

func TestCanonicalInstantiationName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"f", "f"},
		{"byte", "byte"},
		{"If[int]", "If[int]"},
		{"If[byte]", "If[uint8]"},
		{"If[uint8]", "If[uint8]"},
		{"If[rune]", "If[int32]"},
		{"If[[]byte]", "If[[]uint8]"},
		{"If[*byte]", "If[*uint8]"},
		{"If[[][]rune]", "If[[][]int32]"},
		{"If[interface{}]", "If[any]"},
		{"If[[]interface{}]", "If[[]any]"},
		{"Pair[byte,rune]", "Pair[uint8,int32]"},
		{"Pair[string,[]byte]", "Pair[string,[]uint8]"},
		{"If[*example.com/byte.Foo]", "If[*example.com/byte.Foo]"},
		{"If[example.com/x.byte]", "If[example.com/x.byte]"},
	}

	for _, test := range tests {
		have := canonicalInstantiationName(test.name)
		if have != test.want {
			t.Errorf("canonicalInstantiationName(%q):\nhave: %q\nwant: %q", test.name, have, test.want)
		}
	}
}
