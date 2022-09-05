package v5

import (
	"reflect"
	"testing"
)

func Test_DrawFSM(t *testing.T) {
	type test struct {
		name, regex, expected string
	}

	tests := []test{
		{
			name:  "simple example",
			regex: "abc",
			expected: `graph LR
0((0)) --"a"--> 1((1))
1((1)) --"b"--> 2((2))
2((2)) --"c"--> 3((3))`,
		},
		{
			name:  "example with whitespace",
			regex: "a b",
			expected: `graph LR
0((0)) --"a"--> 1((1))
1((1)) --" "--> 2((2))
2((2)) --"b"--> 3((3))`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			drawing := NewMyRegex(tt.regex).DebugFSM()

			if drawing != tt.expected {
				t.Fatalf("Expected drawing to be \n\"%s\", got\n\"%s\"", tt.expected, drawing)
			}
		})
	}
}

func Test_DrawSnapshot(t *testing.T) {
	type test struct {
		name, regex, input, expected string
	}

	tests := []test{
		{
			name:  "initial snapshot",
			regex: "abc",
			input: "",
			expected: `graph LR
0((0)) --"a"--> 1((1))
1((1)) --"b"--> 2((2))
2((2)) --"c"--> 3((3))
style 0 fill:#ff5555;`,
		},
		{
			name:  "after a single letter",
			regex: "abc",
			input: "a",
			expected: `graph LR
0((0)) --"a"--> 1((1))
1((1)) --"b"--> 2((2))
2((2)) --"c"--> 3((3))
style 1 fill:#ff5555;`,
		},
		{
			name:  "last state highlighted",
			regex: "aaa",
			input: "aaa",
			expected: `graph LR
0((0)) --"a"--> 1((1))
1((1)) --"a"--> 2((2))
2((2)) --"a"--> 3((3))
style 3 fill:#00ab41;`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := lex(tt.regex)
			parser := NewParser()
			ast := parser.Parse(tokens)
			state, _ := ast.compile()
			runner := NewRunner(state)
			for _, char := range tt.input {
				runner.Next(char)
			}
			snapshot := runner.drawSnapshot()

			if !reflect.DeepEqual(tt.expected, snapshot) {
				t.Fatalf("Expected drawing to be \n\"%v\"\ngot\n\"%v\"", tt.expected, snapshot)
			}
		})
	}
}
