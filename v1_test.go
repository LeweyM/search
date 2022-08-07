package v1

import (
	"testing"
)

func TestHandmadeFSM(t *testing.T) {
	// handMade
	startState := State{}
	stateA := State{}
	stateB := State{}
	stateC := State{}

	startState.transitions = append(startState.transitions, Transition{
		to:        &stateA,
		predicate: func(input rune) bool { return input == 'a' },
	})

	stateA.transitions = append(stateA.transitions, Transition{
		to:        &stateB,
		predicate: func(input rune) bool { return input == 'b' },
	})

	stateB.transitions = append(stateB.transitions, Transition{
		to:        &stateC,
		predicate: func(input rune) bool { return input == 'c' },
	})

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
			testRunner := NewRunner(&startState)

			for _, character := range tt.input {
				testRunner.Next(character)
			}

			result := testRunner.GetStatus()
			if tt.expectedStatus != result {
				t.Fatalf("Expected FSM to have final state of '%v', got '%v'", tt.expectedStatus, result)
			}
		})
	}
}
