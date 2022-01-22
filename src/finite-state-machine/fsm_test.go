package finite_state_machine

import (
	"fmt"
	"testing"
)

type fsMachine struct {
	description string
	fs          finiteState
}

func TestFsm(t *testing.T) {
	for _, tt := range []struct {
		s                  string
		finiteStateMachine fsMachine
		expectedResults    []int
	}{
		{
			s:                  "aaaa",
			finiteStateMachine: aaaMatcher(t),
			expectedResults:    []int{2, 3},
		},
		{
			s:                  "ab",
			finiteStateMachine: abcMatcher(t),
			expectedResults:    []int{},
		}, {
			s:                  "abcdefg",
			finiteStateMachine: abcMatcher(t),
			expectedResults:    []int{2}, //gets end of result, TODO: get beginning and end
		}, {
			s:                  "abcabc",
			finiteStateMachine: abcMatcher(t),
			expectedResults:    []int{2, 5},
		}, {
			s:                  "abcdefg",
			finiteStateMachine: cMatcher(t),
			expectedResults:    []int{2},
		}, {
			s:                  "ccc",
			finiteStateMachine: cMatcher(t),
			expectedResults:    []int{0, 1, 2},
		}, {
			s:                  "abd",
			finiteStateMachine: cMatcher(t),
			expectedResults:    []int{},
		},
	} {
		t.Run(fmt.Sprintf("Search for %s in string '%s'", tt.finiteStateMachine.description, tt.s), func(t *testing.T) {
			test(t, tt.s, tt.finiteStateMachine.fs, tt.expectedResults)
		})
	}
}

func test(t *testing.T, s string, finiteStateMachine finiteState, expectedResults []int) {
	var results []int

	i := 0
	// not using iterator as i here as it counts bytes, not runes
	for _, char := range s {
		finiteStateMachine = finiteStateMachine.test(char)
		if finiteStateMachine.isEndState() {
			results = append(results, i)
		}
		i++
	}

	if len(results) != len(expectedResults) {
		t.Fatalf("wrong number of results, expected %d, got %d", len(expectedResults), len(results))
	}

	for j := range results {
		if results[j] != expectedResults[j] {
			t.Fatalf("wrong result for result %d: expected %d, got %d", j, expectedResults[j], results[j])
		}
	}
}

// () -a-> () -a-> () -a-> (!) <-a
//  <-----------------------
func aaaMatcher(t *testing.T) fsMachine {
	state, err := NewBuilder().
		State(matchesLetter{letter: 'a'}).
		State(matchesLetter{letter: 'a'}).
		State(matchesLetter{letter: 'a'}).
		State(matchesLetter{letter: 'a'}).End().Recursive().
		Build()
	if err != nil {
		t.Fatalf("Cannot build finite state machine: %v", err)
	}
	return fsMachine{description: "aaa", fs: state}
}

// () -a-> () -b-> () -c-> (!) -a
//          ^-------------------
func abcMatcher(t *testing.T) fsMachine {
	state, err := NewBuilder().
		State(matchesLetter{letter: 'a'}).
		State(matchesLetter{letter: 'b'}).
		State(matchesLetter{letter: 'c'}).
		State(matchesLetter{letter: 'a'}).End().To(1).
		Build()
	if err != nil {
		t.Fatalf("Cannot build finite state machine: %v", err)
	}
	return fsMachine{description: "abc", fs: state}
}

// () -c-> (!) <--c
func cMatcher(t *testing.T) fsMachine {
	state, err := NewBuilder().
		State(matchesLetter{letter: 'c'}).
		State(matchesLetter{letter: 'c'}).End().Recursive().
		Build()
	if err != nil {
		t.Fatalf("Cannot build finite state machine: %v", err)
	}
	return fsMachine{description: "c", fs: state}
}
