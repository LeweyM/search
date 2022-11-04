package v2

import (
	"testing"
)

func TestCompiledFSM(t *testing.T) {
	tokens := lex("abc")
	parser := NewParser(tokens)
	ast := parser.Parse()
	startState, _ := ast.compile()

	type test struct {
		name           string
		input          string
		expectedStatus Status
	}

	tests := []test{
		{"empty string", "", Normal},
		{"non matching string", "xxx", Fail},
		{"matching string", "abc", Success},
		{"partial matching string", "ab", Normal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := NewRunner(startState)

			for _, character := range tt.input {
				runner.Next(character)
			}

			result := runner.GetStatus()
			if tt.expectedStatus != result {
				t.Fatalf("Expected FSM to have final state of '%v', got '%v'", tt.expectedStatus, result)
			}
		})
	}
}
