package finite_state_machine

import (
	"fmt"
	"testing"
)

type fsmTest struct {
	s               string
	expectedResults []localResult
}

type compiledTest struct {
	regex           string
	input           string
	expectedResults []localResult
}

func BenchmarkLinkedFSM(b *testing.B) {
	for i := 0; i < b.N; i++ {
		compiledMachine := Compile("abc")
		runner := NewRunner(compiledMachine)
		runner.Reset()
		FindAll(runner, "adbdbsbabc")
	}
}

func TestCharacterwithwildcardmodifier(t *testing.T) {
	for _, tt := range []compiledTest{
		// (1) <-a- -b-> (2!)
		{regex: "a*b", input: "ab", expectedResults: []localResult{{0, 1}}},
		{regex: "a*b", input: "aab", expectedResults: []localResult{{0, 2}}},
		{regex: "a*b", input: "aaab", expectedResults: []localResult{{0, 3}}},
		{regex: "a*b", input: "b", expectedResults: []localResult{{0, 0}}},
		{regex: "a*b", input: "bb", expectedResults: []localResult{{0, 0}, {1, 1}}},
		{regex: "a*b", input: "a"},
		{regex: "a*b", input: "aa"},
		{regex: "abc*", input: "ab", expectedResults: []localResult{{0, 1}}},
		{regex: "abc*", input: "abc", expectedResults: []localResult{{0, 1}}},  // don't match full string as they are greedy
		{regex: "abc*", input: "abcc", expectedResults: []localResult{{0, 1}}}, // don't match full string as they are greedy
	} {
		testCompiledMachine(t, tt.regex, fsmTest{s: tt.input, expectedResults: tt.expectedResults})
	}
}
func TestDeeplyNestedCatenation(t *testing.T) {
	for _, tt := range []compiledTest{
		{regex: "(((ab)(c)d)|(fg))", input: "abcd", expectedResults: []localResult{{0, 3}}},
		{regex: "(((ab)(c)d)|(fg))", input: "fg", expectedResults: []localResult{{0, 1}}},
		{regex: "(((ab)(c)d)|(fg))", input: "abc"},
		{regex: "(((ab)(c)d)|(fg))", input: "f"},
	} {
		testCompiledMachine(t, tt.regex, fsmTest{s: tt.input, expectedResults: tt.expectedResults})
	}
}
func TestMultiplePipeBranches(t *testing.T) {
	for _, tt := range []compiledTest{
		{regex: "abc|def|xyz", input: "abc", expectedResults: []localResult{{0, 2}}},
		{regex: "abc|def|xyz", input: "def", expectedResults: []localResult{{0, 2}}},
		{regex: "abc|def|xyz", input: "xyz", expectedResults: []localResult{{0, 2}}},
		{regex: "abc|abx|aby|abz", input: "abz", expectedResults: []localResult{{0, 2}}},
		{regex: "abc|abx|aby|abz", input: "abc", expectedResults: []localResult{{0, 2}}},
		{regex: "abc|abx|aby|abz", input: "abr"},

		{regex: "(l|L)et", input: "light"},
		{regex: "(l|L)et", input: "let", expectedResults: []localResult{{0, 2}}},
		{regex: "(l|L)et", input: "Let", expectedResults: []localResult{{0, 2}}},
	} {
		testCompiledMachine(t, tt.regex, fsmTest{s: tt.input, expectedResults: tt.expectedResults})
	}
}
func TestZeroormoreofagroup(t *testing.T) {
	for _, tt := range []compiledTest{
		{regex: "(ab)*", input: "", expectedResults: []localResult{{0, 0}}},
		{regex: "(ab)*", input: "ab", expectedResults: []localResult{{0, 0}, {1, 1}, {2, 2}}},                   // too greedy for interesting results
		{regex: "(ab)*", input: "abab", expectedResults: []localResult{{0, 0}, {1, 1}, {2, 2}, {3, 3}, {4, 4}}}, // too greedy for interesting results
	} {
		testCompiledMachine(t, tt.regex, fsmTest{s: tt.input, expectedResults: tt.expectedResults})
	}

}
func TestCharacterwithoneormoremodifier(t *testing.T) {
	for _, tt := range []compiledTest{
		// 1) -a-> (2) <-a- -b-> (3!)
		{regex: "a+b", input: "ab", expectedResults: []localResult{{0, 1}}},
		{regex: "a+b", input: "aab", expectedResults: []localResult{{0, 2}}},
		{regex: "a+b", input: "aaab", expectedResults: []localResult{{0, 3}}},
		{regex: "a+b", input: "aazb"},
		{regex: "a+b", input: "a"},
		{regex: "a+b", input: "b"},

		{regex: "xy+", input: "x"},
		{regex: "xy+", input: "y"},
		{regex: "xy+", input: "xy", expectedResults: []localResult{{0, 1}}},
		{regex: "xy+", input: "xyxy", expectedResults: []localResult{{0, 1}, {2, 3}}},
		{regex: "xy+", input: "xyyy", expectedResults: []localResult{{0, 1}}}, // too greedy, will grab the first match which is (0,1) instead of (0,3)
	} {
		testCompiledMachine(t, tt.regex, fsmTest{s: tt.input, expectedResults: tt.expectedResults})
	}

}
func TestCharacterWithZeroOrOneModifier(t *testing.T) {
	for _, tt := range []compiledTest{
		//
		{regex: "a?b", input: "ab", expectedResults: []localResult{{0, 1}}},
		{regex: "a?b", input: "b", expectedResults: []localResult{{0, 0}}},
		{regex: "a?b", input: "a"},

		{regex: "cats?", input: "cat", expectedResults: []localResult{{0, 2}}},
		{regex: "cats?", input: "cats", expectedResults: []localResult{{0, 2}}}, // too greedy

		{regex: "held?p?", input: "held", expectedResults: []localResult{{0, 2}}}, // too greedy
		{regex: "held?p?", input: "help", expectedResults: []localResult{{0, 2}}}, // too greedy
		{regex: "held?p?", input: "hel.", expectedResults: []localResult{{0, 2}}},
		{regex: "held?p?", input: "helt", expectedResults: []localResult{{0, 2}}},
	} {
		testCompiledMachine(t, tt.regex, fsmTest{s: tt.input, expectedResults: tt.expectedResults})
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
	m := NewStateLinkedBuilder().
		AddTransition(1, 2, 'd').
		AddTransition(2, 3, 'o').
		AddTransition(3, 4, 'g').
		AddTransition(1, 5, 'd').
		AddTransition(5, 6, 'o').
		AddTransition(6, 7, 't').
		Build()

	for _, tt := range []fsmTest{
		{"dog", []localResult{{0, 2}}},
		{"dot", []localResult{{0, 2}}},
		{"dox", []localResult{}},
		{"doxdog", []localResult{{3, 5}}},
		{"doxdot", []localResult{{3, 5}}},
		{"dodot", []localResult{{2, 4}}},
		{"dodog", []localResult{{2, 4}}},
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
	dog := NewStateLinkedBuilder().
		AddTransition(1, 2, 'd').
		AddTransition(2, 3, 'o').
		AddTransition(3, 4, 'g').
		Build()
	dot := NewStateLinkedBuilder().
		AddTransition(1, 2, 'd').
		AddTransition(2, 3, 'o').
		AddTransition(3, 4, 't').
		Build()
	m := NewStateLinkedBuilder().
		AddMachineTransition(1, dog).
		AddMachineTransition(1, dot).
		Build()

	for _, tt := range []fsmTest{
		{"dog", []localResult{{0, 2}}},
		{"dot", []localResult{{0, 2}}},
		{"dox", []localResult{}},
		{"doxdog", []localResult{{3, 5}}},
		{"doxdot", []localResult{{3, 5}}},
		{"dodot", []localResult{{2, 4}}},
		{"dodog", []localResult{{2, 4}}},
	} {
		testMachine(t, desc, tt, m)
		testCompiledMachine(t, desc, tt)
	}
}

func TestLinkedBranchMatcher(t *testing.T) {
	// (1) -a-> (2) -b-> (3) -c-> (4!)
	//			    -d-> (5) -e---^
	desc := "a(bc|de)"
	m := NewStateLinkedBuilder().
		AddTransition(1, 2, 'a').
		AddTransition(2, 3, 'b').
		AddTransition(3, 4, 'c').
		AddTransition(2, 5, 'd').
		AddTransition(5, 4, 'e').
		Build()

	for _, tt := range []fsmTest{
		{"abc", []localResult{{0, 2}}},
		{"ade", []localResult{{0, 2}}},
		{"abd", []localResult{}},
		{"bc", []localResult{}},
	} {
		testMachine(t, desc, tt, m)
		testCompiledMachine(t, desc, tt)
	}
}

func TestManyRunnersLinked(t *testing.T) {
	// (1) -a-> (2r)<-* -a-> (3) -a-> (4) -a-> (5) -a-> (6) -a-> (7!)
	regex := "a*aaaaa"

	for _, tt := range []fsmTest{
		{"aaaaa", []localResult{{0, 4}}},
		{"aaaaaa", []localResult{{0, 4}}},
		{"aaaaaaa", []localResult{{0, 4}}},
	} {
		testCompiledMachine(t, regex, tt)
	}
}

func TestLinkedAnyCharacterMatcher(t *testing.T) {
	// (1) -a-> (2) -*-> (3) -b-> (4!)
	desc := "a.b"
	m := NewStateLinkedBuilder().
		AddTransition(1, 2, 'a').
		AddWildTransition(2, 3).
		AddTransition(3, 4, 'b').
		Build()

	for _, tt := range []fsmTest{
		{"", []localResult{}},
		{"ab", []localResult{}},
		{"azb", []localResult{{0, 2}}},
		{"acb", []localResult{{0, 2}}},
		{"azzzb", []localResult{}},
	} {
		testMachine(t, desc, tt, m)
		testCompiledMachine(t, desc, tt)
	}
}

func TestLinkedWildcardMatcher(t *testing.T) {
	// () -a-> (r)<-* -b-> (!)
	desc := "a.*b"
	m := NewStateLinkedBuilder().
		AddTransition(1, 2, 'a').
		AddWildTransition(2, 2).
		AddTransition(2, 3, 'b').
		Build()

	for _, tt := range []fsmTest{
		{"ab", []localResult{{0, 1}}},
		{"azb", []localResult{{0, 2}}},
		{"azzzb", []localResult{{0, 4}}},
		{"azzz", []localResult{}},
		{"ba", []localResult{}},
		{"aaaabbbb", []localResult{{0, 4}}},
		{"ababaccb", []localResult{{0, 1}, {2, 3}, {4, 7}}},
	} {
		testMachine(t, desc, tt, m)
		testCompiledMachine(t, desc, tt)
	}
}

func TestSimpleMatcher(t *testing.T) {
	//// () -a-> () -b-> () -c-> (!)
	desc := "abc"
	m := NewStateLinkedBuilder().
		AddTransition(1, 2, 'a').
		AddTransition(2, 3, 'b').
		AddTransition(3, 4, 'c').
		Build()

	for _, tt := range []fsmTest{
		{"abcdefg", []localResult{{0, 2}}},
		{"abcabc", []localResult{{0, 2}, {3, 5}}},
		{"ab", []localResult{}},
	} {
		testMachine(t, desc, tt, m)
		testCompiledMachine(t, desc, tt)
	}
}

func TestSimpleOverlappingMatcher(t *testing.T) {
	//// () -a-> () -a-> () -a-> (!)
	desc := "aaa"
	m := NewStateLinkedBuilder().
		AddTransition(1, 2, 'a').
		AddTransition(2, 3, 'a').
		AddTransition(3, 4, 'a').
		Build()

	for _, tt := range []fsmTest{
		{"aaa", []localResult{{0, 2}}},
		{"aab", []localResult{}},
		{"aaaaaa", []localResult{{0, 2}, {3, 5}}},
	} {
		testMachine(t, desc, tt, m)
		testCompiledMachine(t, desc, tt)
	}
}

func testCompiledMachine(t *testing.T, regex string, tt fsmTest) bool {
	return t.Run(fmt.Sprintf("FindAll for compiled('%s') in string '%s'", regex, tt.s), func(t *testing.T) {
		compiledMachine := Compile(regex)
		runner := NewRunner(compiledMachine)
		runner.Reset()
		testFindAll(t, tt.s, runner, tt.expectedResults)
	})
}

func testMachine(t *testing.T, desc string, tt fsmTest, m *StateLinked) bool {
	return t.Run(fmt.Sprintf("FindAll for '%s' in string '%s'", desc, tt.s), func(t *testing.T) {
		runner := NewRunner(m)
		runner.Reset()
		testFindAll(t, tt.s, runner, tt.expectedResults)
	})
}
