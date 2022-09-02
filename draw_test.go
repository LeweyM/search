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
		regex    string
		input    string
		expected []string
	}

	tests := []test{
		{
			regex: "abc",
			input: "abc",
			expected: []string{`graph LR
0((0)) --"a"--> 1((1))
1((1)) --"b"--> 2((2))
2((2)) --"c"--> 3((3))
style 0 fill:#ff5555;`,
				`graph LR
0((0)) --"a"--> 1((1))
1((1)) --"b"--> 2((2))
2((2)) --"c"--> 3((3))
style 1 fill:#ff5555;`,
				`graph LR
0((0)) --"a"--> 1((1))
1((1)) --"b"--> 2((2))
2((2)) --"c"--> 3((3))
style 2 fill:#ff5555;`,
				`graph LR
0((0)) --"a"--> 1((1))
1((1)) --"b"--> 2((2))
2((2)) --"c"--> 3((3))
style 3 fill:#00ab41;`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			regex := NewMyRegex(tt.regex)
			stepDrawings := regex.DebugMatch(tt.input)
			result := mapToGraphDrawings(stepDrawings)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Fatalf("Expected drawing to be \n\"%s\"\ngot\n\"%s\"", tt.expected, result)
			}
		})
	}
}

func mapToGraphDrawings(drawings []StepDrawing) []string {
	res := make([]string, len(drawings))
	for i := range drawings {
		res[i] = drawings[i].graphDrawing
	}
	return res
}
