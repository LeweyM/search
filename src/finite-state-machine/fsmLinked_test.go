package finite_state_machine

import (
	"fmt"
	"testing"
)

func TestLinkedDeeplyNestedCompiledMatcher(t *testing.T) {
	// equivalent to "abcd|fg"
	desc := "((abcd)|(fg))"
	//desc := "(((ab)(c)d)|(fg))" // TODO: Parenthesis for catenations
	for _, tt := range []struct {
		s               string
		expectedResults []result
	}{
		{"abcd", []result{{0, 3}}},
		{"fg", []result{{0, 1}}},
		{"abc", []result{}},
		{"f", []result{}},
	} {
		testCompiledMachine(t, desc, tt)
	}
}

// Overlapping branches can be reduced to single matching branches.
// One option is not to try to resolve this but to compile
// any overlaps into a single matcher.
//
// i.e. "(dog|dot)" => "do(g|t)"
func TestLinkedOverlappingBranchMatcher(t *testing.T) {
	// (1) -d-> (2) -o-> (3) -g-> (4!)
	//     -d-> (5) -o-> (6) -t-> (7!)
	desc := "(dog|dot)"
	m := NewStateLinkedBuilder(7).
		AddTransition(1, 2, 'd').
		AddTransition(2, 3, 'o').
		AddTransition(3, 4, 'g').SetSuccess(4).
		AddTransition(1, 5, 'd').
		AddTransition(5, 6, 'o').
		AddTransition(6, 7, 't').SetSuccess(7).
		Build()

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
		testMachine(t, desc, tt, m)
		testCompiledMachine(t, desc, tt)
	}
}

func TestOverlappingBranchComposableMatcher(t *testing.T) {
	// (1) -d-> (2) -o-> (3) -g-> (4!)
	//     -d-> (5) -o-> (6) -t-> (7!)

	// or, as composable machines
	// (9)(dog 1-2-3-4)!
	//    (dot 5-6-7-8)!
	desc := "dog|dot"
	dog := NewStateLinkedBuilder(4).
		AddTransition(1, 2, 'd').
		AddTransition(2, 3, 'o').
		AddTransition(3, 4, 'g').SetSuccess(4).
		Build()
	dot := NewStateLinkedBuilder(4).
		AddTransition(1, 2, 'd').
		AddTransition(2, 3, 'o').
		AddTransition(3, 4, 't').SetSuccess(4).
		Build()
	m := NewStateLinkedBuilder(2).
		AddMachineTransition(1, dog).
		AddMachineTransition(1, dot).
		Build()

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
		testMachine(t, desc, tt, m)
		testCompiledMachine(t, desc, tt)
	}
}

func TestLinkedBranchMatcher(t *testing.T) {
	// (1) -a-> (2) -b-> (3) -c-> (4!)
	//			    -d-> (5) -e---^
	desc := "a(bc|de)"
	m := NewStateLinkedBuilder(5).
		AddTransition(1, 2, 'a').
		AddTransition(2, 3, 'b').
		AddTransition(3, 4, 'c').SetSuccess(4).
		AddTransition(2, 5, 'd').
		AddTransition(5, 4, 'e').
		Build()

	for _, tt := range []struct {
		s               string
		expectedResults []result
	}{
		{"abc", []result{{0, 2}}},
		{"ade", []result{{0, 2}}},
		{"abd", []result{}},
		{"bc", []result{}},
	} {
		testMachine(t, desc, tt, m)
		testCompiledMachine(t, desc, tt)
	}
}

func TestManyRunnersLinked(t *testing.T) {
	// (1) -a-> (2r)<-* -a-> (3) -a-> (4) -a-> (5) -a-> (6) -a-> (7!)
	desc := "a*aaaaa"
	m := NewStateLinkedBuilder(7).
		AddTransition(1, 2, 'a').
		AddWildTransition(2, 2).
		AddTransition(2, 3, 'a').
		AddTransition(3, 4, 'a').
		AddTransition(4, 5, 'a').
		AddTransition(5, 6, 'a').
		AddTransition(6, 7, 'a').SetSuccess(7).
		Build()

	for _, tt := range []struct {
		s               string
		expectedResults []result
	}{
		{"aaaaaa", []result{{0, 5}}},
	} {
		testMachine(t, desc, tt, m)
		testCompiledMachine(t, desc, tt)
	}
}

func TestLinkedAnyCharacterMatcher(t *testing.T) {
	// (1) -a-> (2) -*-> (3) -b-> (4!)
	desc := "a.b"
	m := NewStateLinkedBuilder(4).
		AddTransition(1, 2, 'a').
		AddWildTransition(2, 3).
		AddTransition(3, 4, 'b').
		SetSuccess(4).
		Build()

	for _, tt := range []struct {
		s               string
		expectedResults []result
	}{
		{"", []result{}},
		{"ab", []result{}},
		{"azb", []result{{0, 2}}},
		{"acb", []result{{0, 2}}},
		{"azzzb", []result{}},
	} {
		testMachine(t, desc, tt, m)
		testCompiledMachine(t, desc, tt)
	}
}

func TestLinkedWildcardMatcher(t *testing.T) {
	// () -a-> (r)<-* -b-> (!)
	desc := "a.*b"
	m := NewStateLinkedBuilder(3).
		AddTransition(1, 2, 'a').
		AddWildTransition(2, 2).
		AddTransition(2, 3, 'b').
		SetSuccess(3).
		Build()

	for _, tt := range []struct {
		s               string
		expectedResults []result
	}{
		{"ab", []result{{0, 1}}},
		{"azb", []result{{0, 2}}},
		{"azzzb", []result{{0, 4}}},
		{"azzz", []result{}},
		{"ba", []result{}},
		{"aaaabbbb", []result{{0, 4}}},
		{"ababaccb", []result{{0, 1}, {2, 3}, {4, 7}}},
	} {
		testMachine(t, desc, tt, m)
		testCompiledMachine(t, desc, tt)
	}
}

func TestSimpleMatcher(t *testing.T) {
	//// () -a-> () -b-> () -c-> (!)
	desc := "abc"
	m := NewStateLinkedBuilder(4).
		AddTransition(1, 2, 'a').
		AddTransition(2, 3, 'b').
		AddTransition(3, 4, 'c').
		SetSuccess(4).
		Build()

	for _, tt := range []struct {
		s               string
		expectedResults []result
	}{
		{"abcdefg", []result{{0, 2}}},
		{"abcabc", []result{{0, 2}, {3, 5}}},
		{"ab", []result{}},
	} {
		testMachine(t, desc, tt, m)
		testCompiledMachine(t, desc, tt)
	}
}

func TestSimpleOverlappingMatcher(t *testing.T) {
	//// () -a-> () -a-> () -a-> (!)
	desc := "aaa"
	m := NewStateLinkedBuilder(4).
		AddTransition(1, 2, 'a').
		AddTransition(2, 3, 'a').
		AddTransition(3, 4, 'a').
		SetSuccess(4).
		Build()

	for _, tt := range []struct {
		s               string
		expectedResults []result
	}{
		{"aaa", []result{{0, 2}}},
		{"aab", []result{}},
		{"aaaaaa", []result{{0, 2}, {3, 5}}},
	} {
		testMachine(t, desc, tt, m)
		testCompiledMachine(t, desc, tt)
	}
}

func testCompiledMachine(t *testing.T, desc string, tt struct {
	s               string
	expectedResults []result
}) bool {
	return t.Run(fmt.Sprintf("FindAll for compiled('%s') in string '%s'", desc, tt.s), func(t *testing.T) {
		compiledMachine := Compile(desc)
		runner := NewRunner(compiledMachine)
		runner.Reset()
		testFindAll(t, tt.s, runner, tt.expectedResults)
	})
}

func testMachine(t *testing.T, desc string, tt struct {
	s               string
	expectedResults []result
}, m *StateLinked) bool {
	return t.Run(fmt.Sprintf("FindAll for '%s' in string '%s'", desc, tt.s), func(t *testing.T) {
		runner := NewRunner(m)
		runner.Reset()
		testFindAll(t, tt.s, runner, tt.expectedResults)
	})
}
