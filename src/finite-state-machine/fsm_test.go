package finite_state_machine

import "testing"

func TestFsm(t *testing.T) {
	ss := startingState{}
	ss.nextState = newMatchesLetter(&ss, 'c', endState{})
	cMatcher := NewFsm(&ss)

	for _, tt := range []struct {
		s                  string
		finiteStateMachine *fsm
		expectedResults    []int
	}{
		{
			s:                  "abcdefg",
			finiteStateMachine: cMatcher,
			expectedResults:    []int{2},
		}, {
			s:                  "ccc",
			finiteStateMachine: cMatcher,
			expectedResults:    []int{0, 1, 2},
		}, {
			s:                  "abd",
			finiteStateMachine: cMatcher,
			expectedResults:    []int{},
		},
	}{
		test(t, tt.s, tt.finiteStateMachine, tt.expectedResults)
	}
}

func test(t *testing.T, s string, finiteStateMachine *fsm, expectedResults []int) {
	var results []int

	i := 0
	// not using iterator as i here as it counts bytes, not runes
	for _, char := range s {
		matches := finiteStateMachine.next(char)
		if matches {
			results = append(results, i)
		}
		i++
	}

	if len(results) != len(expectedResults) {
		t.Errorf("wrong number of results, expected %d, got %d", 1, len(results))
	}

	for j := range results {
		if results[j] != expectedResults[j] {
			t.Errorf("wrong result for result %d: expected %d, got %d", j, expectedResults[j], results[j])
		}
	}
}
