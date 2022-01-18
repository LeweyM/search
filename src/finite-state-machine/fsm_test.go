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
			finiteStateMachine: AAAMatcher(),
			expectedResults:    []int{2, 3},
		},
		{
			s:                  "ab",
			finiteStateMachine: ABCMatcher(),
			expectedResults:    []int{},
		}, {
			s:                  "abcdefg",
			finiteStateMachine: ABCMatcher(),
			expectedResults:    []int{2}, //gets end of result, TODO: get beginning and end
		}, {
			s:                  "abcdefg",
			finiteStateMachine: CMatcher(),
			expectedResults:    []int{2},
		}, {
			s:                  "ccc",
			finiteStateMachine: CMatcher(),
			expectedResults:    []int{0, 1, 2},
		}, {
			s:                  "abd",
			finiteStateMachine: CMatcher(),
			expectedResults:    []int{},
		},
	} {
		t.Run(fmt.Sprintf("Search for %s in string '%s'", tt.finiteStateMachine.description, tt.s), func(t *testing.T) {
			test(t, tt.s, tt.finiteStateMachine.fs, tt.expectedResults)
		})
	}
}

func AAAMatcher() fsMachine {
	return fsMachine{description: "aaa", fs: buildAAAMatcher()}
}

// () -a-> () -a-> () -a-> (!) <-a
//  <-----------------------
func buildAAAMatcher() finiteState {
	aRecursive := matchesLetter{letter: 'a'}.End()
	aRecursive.nextState = &aRecursive
	aMatcher := matchesLetter{letter: 'a', nextState: &aRecursive}
	aaMatcher := matchesLetter{letter: 'a', nextState: &aMatcher}
	aaaMatcher := matchesLetter{letter: 'a', nextState: &aaMatcher}
	aRecursive = aRecursive.Base(&aaMatcher)
	aaaMatcher = aaaMatcher.Base(&aaMatcher)
	aaMatcher = aaMatcher.Base(&aaMatcher)
	aMatcher = aMatcher.Base(&aaMatcher)
	return aaaMatcher
}

func ABCMatcher() fsMachine {
	return fsMachine{description: "abc", fs: buildABCMatcher()}
}

// () -a-> () -b-> () -c-> (!)
func buildABCMatcher() finiteState {
	abcExitMatcher := matchesLetter{letter: 'a'}.End()
	cMatcher := matchesLetter{letter: 'c', nextState: &abcExitMatcher}
	bcMatcher := matchesLetter{letter: 'b', nextState: &cMatcher}
	abcMatcher := matchesLetter{letter: 'a', nextState: &bcMatcher}
	abcExitMatcher.nextState = &bcMatcher
	cMatcher = cMatcher.Base(&abcMatcher)
	bcMatcher = bcMatcher.Base(&abcMatcher)
	abcMatcher = abcMatcher.Base(&abcMatcher)
	abcExitMatcher = abcExitMatcher.Base(&abcMatcher)
	return abcMatcher
}

func CMatcher() fsMachine {
	return fsMachine{description: "c", fs: buildCMatcher()}
}

// () -c-> (!) <--c
func buildCMatcher() finiteState {
	cMatcher := matchesLetter{letter: 'c'}.End()
	cMatcher.nextState = &cMatcher
	rootState := matchesLetter{letter: 'c', nextState: &cMatcher}
	rootState = rootState.Base(&rootState)
	cMatcher = cMatcher.Base(&rootState)
	return rootState
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
