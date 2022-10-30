package v7

import (
	"reflect"
	"testing"
)

func abcBuilder() *State {
	state1, state2, state3, state4 := &State{}, &State{}, &State{}, &State{}

	state1.addTransition(state2, Predicate{allowedChars: "a"}, "a")
	state2.addTransition(state3, Predicate{allowedChars: "b"}, "b")
	state3.addTransition(state4, Predicate{allowedChars: "c"}, "c")
	return state1
}

func aaaBuilder() *State {
	state1, state2, state3, state4 := &State{}, &State{}, &State{}, &State{}

	state1.addTransition(state2, Predicate{allowedChars: "a"}, "a")
	state2.addTransition(state3, Predicate{allowedChars: "a"}, "a")
	state3.addTransition(state4, Predicate{allowedChars: "a"}, "a")
	return state1
}

func aεbBuilder() *State {
	state1, state2, state3, state4 := &State{}, &State{}, &State{}, &State{}

	state1.addTransition(state2, Predicate{allowedChars: "a"}, "a")
	state2.addEpsilon(state3)
	state3.addTransition(state4, Predicate{allowedChars: "b"}, "b")
	return state1
}

func Test_DrawFSM(t *testing.T) {
	type test struct {
		name, expected string
		fsmBuilder     func() *State
	}

	tests := []test{
		{
			name:       "simple example",
			fsmBuilder: abcBuilder,
			expected: `graph LR
0((0)) --"a"--> 1((1))
1((1)) --"b"--> 2((2))
2((2)) --"c"--> 3((3))
style 3 stroke:green,stroke-width:4px;`,
		},
		{
			name:       "graph with epsilon",
			fsmBuilder: aεbBuilder,
			expected: `graph LR
0((0)) --"a"--> 1((1))
1((1)) -."ε".-> 2((2))
2((2)) --"b"--> 3((3))
style 3 stroke:green,stroke-width:4px;`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			drawing, _ := tt.fsmBuilder().Draw()

			if drawing != tt.expected {
				t.Fatalf("Expected drawing to be \n\"%s\", got\n\"%s\"", tt.expected, drawing)
			}
		})
	}
}

func Test_DrawSnapshot(t *testing.T) {
	type test struct {
		name, input, expected string
		fsmBuilder            func() *State
	}

	tests := []test{
		{
			name:       "initial snapshot",
			fsmBuilder: abcBuilder,
			input:      "",
			expected: `graph LR
0((0)) --"a"--> 1((1))
1((1)) --"b"--> 2((2))
2((2)) --"c"--> 3((3))
style 3 stroke:green,stroke-width:4px;
style 0 fill:#ff5555;`,
		},
		{
			name:       "after a single letter",
			fsmBuilder: abcBuilder,
			input:      "a",
			expected: `graph LR
0((0)) --"a"--> 1((1))
1((1)) --"b"--> 2((2))
2((2)) --"c"--> 3((3))
style 3 stroke:green,stroke-width:4px;
style 0 fill:#ff5555;
style 1 fill:#ff5555;`,
		},
		{
			name:       "all states highlighted",
			fsmBuilder: aaaBuilder,
			input:      "aaa",
			expected: `graph LR
0((0)) --"a"--> 1((1))
1((1)) --"a"--> 2((2))
2((2)) --"a"--> 3((3))
style 3 stroke:green,stroke-width:4px;
style 0 fill:#ff5555;
style 1 fill:#ff5555;
style 2 fill:#ff5555;
style 3 fill:#00ab41;`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := NewRunner(tt.fsmBuilder())
			for _, char := range tt.input {
				runner.Next(char)
				runner.Start()
			}
			snapshot := runner.drawSnapshot()

			if !reflect.DeepEqual(tt.expected, snapshot) {
				t.Fatalf("Expected drawing to be \n\"%v\"\ngot\n\"%v\"", tt.expected, snapshot)
			}
		})
	}
}
