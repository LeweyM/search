package v5

import (
	"reflect"
	"testing"
)

func Test_DrawFSM(t *testing.T) {
	type test struct {
		input, expected string
	}

	tests := []test{
		{
			input: "abc",
			expected: `graph LR
0((0)) --"a"--> 1((1))
1((1)) --"b"--> 2((2))
2((2)) --"c"--> 3((3))`,
		},
		{
			input: "a b",
			expected: `graph LR
0((0)) --"a"--> 1((1))
1((1)) --" "--> 2((2))
2((2)) --"b"--> 3((3))`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			drawing := NewMyRegex(tt.input).DebugFSM()

			if drawing != tt.expected {
				t.Fatalf("Expected drawing to be \n\"%s\", got\n\"%s\"", tt.expected, drawing)
			}
		})
	}
}

func TestState_DebugMatch(t *testing.T) {
	type test struct {
		name, regex, input string
		expected           []debugStep
	}

	tests := []test{
		{
			name:  "normal with match",
			regex: "abc",
			input: "abc",
			expected: []debugStep{
				{runnerDrawing: `graph LR
0((0)) --"a"--> 1((1))
1((1)) --"b"--> 2((2))
2((2)) --"c"--> 3((3))
style 0 fill:#ff5555;`, currentCharacterIndex: 0},
				{runnerDrawing: `graph LR
0((0)) --"a"--> 1((1))
1((1)) --"b"--> 2((2))
2((2)) --"c"--> 3((3))
style 1 fill:#ff5555;`, currentCharacterIndex: 1},
				{runnerDrawing: `graph LR
0((0)) --"a"--> 1((1))
1((1)) --"b"--> 2((2))
2((2)) --"c"--> 3((3))
style 2 fill:#ff5555;`, currentCharacterIndex: 2},
				{runnerDrawing: `graph LR
0((0)) --"a"--> 1((1))
1((1)) --"b"--> 2((2))
2((2)) --"c"--> 3((3))
style 3 fill:#00ab41;`, currentCharacterIndex: 3},
			},
		},
		{
			name:  "backtracking with match",
			regex: "aab",
			input: "aaab",
			expected: []debugStep{
				{runnerDrawing: `graph LR
0((0)) --"a"--> 1((1))
1((1)) --"a"--> 2((2))
2((2)) --"b"--> 3((3))
style 0 fill:#ff5555;`, currentCharacterIndex: 0},
				{runnerDrawing: `graph LR
0((0)) --"a"--> 1((1))
1((1)) --"a"--> 2((2))
2((2)) --"b"--> 3((3))
style 1 fill:#ff5555;`, currentCharacterIndex: 1},
				{runnerDrawing: `graph LR
0((0)) --"a"--> 1((1))
1((1)) --"a"--> 2((2))
2((2)) --"b"--> 3((3))
style 2 fill:#ff5555;`, currentCharacterIndex: 2},
				{runnerDrawing: `graph LR
0((0)) --"a"--> 1((1))
1((1)) --"a"--> 2((2))
2((2)) --"b"--> 3((3))`, currentCharacterIndex: 3},
				{runnerDrawing: `graph LR
0((0)) --"a"--> 1((1))
1((1)) --"a"--> 2((2))
2((2)) --"b"--> 3((3))
style 0 fill:#ff5555;`, currentCharacterIndex: 1},
				{runnerDrawing: `graph LR
0((0)) --"a"--> 1((1))
1((1)) --"a"--> 2((2))
2((2)) --"b"--> 3((3))
style 1 fill:#ff5555;`, currentCharacterIndex: 2},
				{runnerDrawing: `graph LR
0((0)) --"a"--> 1((1))
1((1)) --"a"--> 2((2))
2((2)) --"b"--> 3((3))
style 2 fill:#ff5555;`, currentCharacterIndex: 3},
				{runnerDrawing: `graph LR
0((0)) --"a"--> 1((1))
1((1)) --"a"--> 2((2))
2((2)) --"b"--> 3((3))
style 3 fill:#00ab41;`, currentCharacterIndex: 4},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			regex := NewMyRegex(tt.regex)
			steps := regex.DebugMatch(tt.input)

			if !reflect.DeepEqual(tt.expected, steps) {
				t.Fatalf("Expected drawing to be \n\"%v\"\ngot\n\"%v\"", tt.expected, steps)
			}
		})
	}
}
