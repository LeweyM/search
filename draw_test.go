package v5

import "testing"

func TestState_Draw(t *testing.T) {
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
			parser := NewParser()

			tokens := lex(tt.input)
			ast := parser.Parse(tokens)
			fsm, _ := ast.compile()

			drawing, _ := fsm.Draw()

			if drawing != tt.expected {
				t.Fatalf("Expected drawing to be \n\"%s\", got\n\"%s\"", tt.expected, drawing)
			}
		})
	}
}
