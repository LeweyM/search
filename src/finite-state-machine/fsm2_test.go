package finite_state_machine

import (
	"fmt"
	"testing"
)

type fs2Machine struct {
	description string
	fs          *machine
}

// Overlapping branches can be reduced to single matching branches.
// One option is not to try to resolve this but to compile
// any overlaps into a single matcher.
//
// i.e. "(dog|dot)" => "do(g|t)"
func TestOverlappingBranchMatcher(t *testing.T) {
	// (1) -d-> (2) -o-> (3) -g-> (4!)
	//     -d-> (5) -o-> (6) -t-> (7!)

	// or, as composable machines
	// (1) -> (dog) -> (2)
	//     -> (dot) ----^
	desc := "(dog|dot)"
	dog := NewMachine(4).
		AddTransition(1, 2, 'd').
		AddTransition(2, 3, 'o').
		AddTransition(3, 4, 'g').SetSuccess(4)
	dot := NewMachine(4).
		AddTransition(1, 2, 'd').
		AddTransition(2, 3, 'o').
		AddTransition(3, 4, 't').SetSuccess(4)
	m := NewMachine(2).
		AddMachineTransition(1, 2, dog).
		AddMachineTransition(1, 2, dot).SetSuccess(2)

	for _, tt := range []struct {
		s               string
		expectedResults []result
	}{
		{"dog", []result{{0, 2}}},
		{"dot", []result{{0, 2}}},
		{"dox", []result{}},
		{"doxdog", []result{{3, 5}}},
		{"doxdot", []result{{3, 5}}},
		{"dodot", []result{{2, 4}}},
		{"dodog", []result{{2, 4}}},
	} {
		t.Run(fmt.Sprintf("FindAll for '%s' in string '%s'", desc, tt.s), func(t *testing.T) {
			m.Reset()
			testFindAll(t, tt.s, m, tt.expectedResults)
		})
	}
}

func TestBranchMatcher(t *testing.T) {
	// (1) -a-> (2) -b-> (3) -c-> (4!)
	//			    -d-> (5) -e---^
	desc := "a(bc|de)"
	m := NewMachine(5).
		AddTransition(1, 2, 'a').
		AddTransition(2, 3, 'b').
		AddTransition(3, 4, 'c').SetSuccess(4).
		AddTransition(2, 5, 'd').
		AddTransition(5, 4, 'e')

	for _, tt := range []struct {
		s               string
		expectedResults []result
	}{
		{"abc", []result{{0, 2}}},
		{"ade", []result{{0, 2}}},
		{"abd", []result{}},
	} {
		t.Run(fmt.Sprintf("FindAll for '%s' in string '%s'", desc, tt.s), func(t *testing.T) {
			m.Reset()
			testFindAll(t, tt.s, m, tt.expectedResults)
		})
	}
}

func TestWildcardMatcher(t *testing.T) {
	// () -a-> (r)<-* -b-> (!)
	desc := "a.*b"
	m := NewMachine(3).
		AddTransition(1, 2, 'a').
		AddWildTransition(2, 2).
		AddTransition(2, 3, 'b').
		SetSuccess(3)

	for _, tt := range []struct {
		s               string
		expectedResults []result
	}{
		{"azzzb", []result{{0, 4}}},
		{"azzz", []result{}},
		{"ba", []result{}},
		{"aaaabbbb", []result{{0, 4}}},
		{"ababaccb", []result{{0, 1}, {2, 3}, {4, 7}}},
	} {
		t.Run(fmt.Sprintf("FindAll for '%s' in string '%s'", desc, tt.s), func(t *testing.T) {
			m.Reset()
			testFindAll(t, tt.s, m, tt.expectedResults)
		})
	}
}

func TestFsm2(t *testing.T) {
	for _, tt := range []struct {
		s                  string
		finiteStateMachine fs2Machine
		expectedResults    []result
	}{
		// notice it does not recognize overlapping matches. E.g. "AAAaaa" & "aAAAaa" & "aaAAAa" & "aaaAAA". Only "AAAaaa" and "aaaAAA" are recognized.
		{
			s:                  "aaaaaa",
			finiteStateMachine: aaaMatcher2(),
			expectedResults:    []result{{0, 2}, {3, 5}},
		},
		{
			s:                  "abaaa",
			finiteStateMachine: aaaMatcher2(),
			expectedResults:    []result{{2, 4}},
		},
		{
			s:                  "ab",
			finiteStateMachine: abcMatcher2(),
			expectedResults:    []result{},
		}, {
			s:                  "abcdefg",
			finiteStateMachine: abcMatcher2(),
			expectedResults:    []result{{0, 2}},
		}, {
			s:                  "abcabc",
			finiteStateMachine: abcMatcher2(),
			expectedResults:    []result{{0, 2}, {3, 5}},
		}, {
			s:                  "abcdefg",
			finiteStateMachine: cMatcher2(),
			expectedResults:    []result{{2, 2}},
		}, {
			s:                  "ccc",
			finiteStateMachine: cMatcher2(),
			expectedResults:    []result{{0, 0}, {1, 1}, {2, 2}},
		}, {
			s:                  "abd",
			finiteStateMachine: cMatcher2(),
			expectedResults:    []result{},
		},
	} {
		t.Run(fmt.Sprintf("FindAll for %s in string '%s'", tt.finiteStateMachine.description, tt.s), func(t *testing.T) {
			testFindAll(t, tt.s, tt.finiteStateMachine.fs, tt.expectedResults)
		})
	}
}

func testFindAll(t *testing.T, s string, finiteStateMachine *machine, expectedResults []result) {
	results := FindAll(finiteStateMachine, s)

	if len(results) != len(expectedResults) {
		t.Fatalf("wrong number of results for string '%s', expected %+v, got %+v", s, len(expectedResults), len(results))
	}

	for j := range results {
		if results[j] != expectedResults[j] {
			t.Fatalf("wrong result for string '%s': expected %d, got %d", s, expectedResults[j], results[j])
		}
	}
}

// () -a-> () -a-> () -a-> (!) <-a
//  <-----------------------
func aaaMatcher2() fs2Machine {
	m := NewMachine(4).
		AddTransition(1, 2, 'a').
		AddTransition(2, 3, 'a').
		AddTransition(3, 4, 'a').
		SetSuccess(4)
	return fs2Machine{description: "aaa", fs: m}
}

// () -a-> () -b-> () -c-> (!) -a
//          ^-------------------
func abcMatcher2() fs2Machine {
	m := NewMachine(4).
		AddTransition(1, 2, 'a').
		AddTransition(2, 3, 'b').
		AddTransition(3, 4, 'c').
		SetSuccess(4)
	return fs2Machine{description: "abc", fs: m}
}

// () -c-> (!) <--c
func cMatcher2() fs2Machine {
	m := NewMachine(2).
		AddTransition(1, 2, 'c').
		SetSuccess(2)
	return fs2Machine{description: "c", fs: m}
}
