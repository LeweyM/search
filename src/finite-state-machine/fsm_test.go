package finite_state_machine

import (
	context2 "context"
	"fmt"
	"testing"
)

type fsmTest struct {
	s               string
	expectedResults []localResult
}

type fsmTestWithLines struct {
	s               string
	expectedResults []Result
}

type compiledTest struct {
	regex           string
	input           string
	expectedResults []localResult
}

type compiledTestWithLines struct {
	regex           string
	input           string
	expectedResults []Result
}

func BenchmarkLinkedFSM(b *testing.B) {
	for i := 0; i < b.N; i++ {
		compiledMachine := Compile("abc")
		runner := NewRunner(compiledMachine)
		runner.Reset()
		FindAll(runner, "adbdbsbabc")
	}
}

func TestCharacterWithWildcardModifier(t *testing.T) {
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
		{regex: "abc*", input: "abc", expectedResults: []localResult{{0, 2}}},
		{regex: "abc*", input: "abcc", expectedResults: []localResult{{0, 3}}},
	} {
		testCompiledMachine(t, tt.regex, fsmTest{s: tt.input, expectedResults: tt.expectedResults})
	}
}

func TestMultiLines(t *testing.T) {
	for _, tt := range []compiledTestWithLines{
		{regex: "(dis)?like", input: "adultlike\nadultness", expectedResults: []Result{{
			Line:  1,
			Start: 5,
			End:   8,
		}}},
		{regex: "cat", input: "acca\nacxxxsabc\naccatura", expectedResults: []Result{{
			Line:  3,
			Start: 2,
			End:   4,
		}}},
		{regex: "cat", input: "ca\nxxx\ncat", expectedResults: []Result{{
			Line:  3,
			Start: 0,
			End:   2,
		}}},
	} {
		testCompiledMachineWithLines(t, tt.regex, fsmTestWithLines{s: tt.input, expectedResults: tt.expectedResults})
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
func TestZeroOrMoreOfAGroup(t *testing.T) {
	for _, tt := range []compiledTest{
		//{regex: "(ab)*", input: "", expectedResults: []localResult{{0, 0}}},
		{regex: "(ab)*", input: "ab", expectedResults: []localResult{{0, 1}}},
		{regex: "(ab)*", input: "abab", expectedResults: []localResult{{0, 3}}},
	} {
		testCompiledMachine(t, tt.regex, fsmTest{s: tt.input, expectedResults: tt.expectedResults})
	}

}

func TestCharacterWithOneOrMoreModifier(t *testing.T) {
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
		{regex: "xy+", input: "xyyy", expectedResults: []localResult{{0, 3}}},
	} {
		testCompiledMachine(t, tt.regex, fsmTest{s: tt.input, expectedResults: tt.expectedResults})
	}
}

func TestOneOrZeroOfGroup(t *testing.T) {
	for _, tt := range []compiledTest{
		// 1) -a-> (2) <-a- -b-> (3!)
		{regex: "a(bc)?d", input: "abcd", expectedResults: []localResult{{0, 3}}},
		{regex: "a(bc)?d", input: "ad", expectedResults: []localResult{{0, 1}}},
		{regex: "a(bc)?d", input: "abd", expectedResults: []localResult{}},
		{regex: "a(bc)?d", input: "abz", expectedResults: []localResult{}},

		{regex: "(dis)?like", input: "adultlike", expectedResults: []localResult{{5, 8}}},
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
		{regex: "cats?", input: "cats", expectedResults: []localResult{{0, 3}}},

		{regex: "held?p?", input: "held", expectedResults: []localResult{{0, 3}}},
		{regex: "held?p?", input: "help", expectedResults: []localResult{{0, 3}}},
		{regex: "held?p?", input: "hel.", expectedResults: []localResult{{0, 2}}},
		{regex: "held?p?", input: "helt", expectedResults: []localResult{{0, 2}}},

		{regex: "ab?c?", input: "abz", expectedResults: []localResult{{0, 1}}},
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
	m := NewStateBuilder().
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
	dog := NewStateBuilder().
		AddTransition(1, 2, 'd').
		AddTransition(2, 3, 'o').
		AddTransition(3, 4, 'g').
		Build()
	dot := NewStateBuilder().
		AddTransition(1, 2, 'd').
		AddTransition(2, 3, 'o').
		AddTransition(3, 4, 't').
		Build()
	m := NewStateBuilder().
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
	m := NewStateBuilder().
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
		{"aaaaaa", []localResult{{0, 5}}},
		{"aaaaaaa", []localResult{{0, 6}}},
	} {
		testCompiledMachine(t, regex, tt)
	}
}

func TestLinkedAnyCharacterMatcher(t *testing.T) {
	// (1) -a-> (2) -*-> (3) -b-> (4!)
	desc := "a.b"
	m := NewStateBuilder().
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
	desc := "a.*b"

	for _, tt := range []fsmTest{
		{"ab", []localResult{{0, 1}}},
		{"azb", []localResult{{0, 2}}},
		{"azzzb", []localResult{{0, 4}}},
		{"azzz", []localResult{}},
		{"ba", []localResult{}},
		{"aaaabbbb", []localResult{{0, 7}}},
		{"ababaccb", []localResult{{0, 7}}},
		{"a1b3456", []localResult{{0, 2}}},
	} {
		testCompiledMachine(t, desc, tt)
	}
}

func TestSimpleMatcher(t *testing.T) {
	//// () -a-> () -b-> () -c-> (!)
	desc := "abc"

	for _, tt := range []fsmTest{
		{"abcdefg", []localResult{{0, 2}}},
		{"abcabc", []localResult{{0, 2}, {3, 5}}},
		{"ab", []localResult{}},
	} {
		testCompiledMachine(t, desc, tt)
	}
}

func TestSimpleOverlappingMatcher(t *testing.T) {
	//// () -a-> () -a-> () -a-> (!)
	desc := "aaa"
	m := NewStateBuilder().
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

func testCompiledMachineWithLines(t *testing.T, regex string, tt fsmTestWithLines) bool {
	return t.Run(fmt.Sprintf("FindAll for compiled('%s') in string '%s'", regex, tt.s), func(t *testing.T) {
		compiledMachine := Compile(regex)
		runner := NewRunner(compiledMachine)
		runner.Reset()
		testFindAllWithLines(t, tt.s, runner, tt.expectedResults)
	})
}

func testMachine(t *testing.T, desc string, tt fsmTest, m *State) bool {
	return t.Run(fmt.Sprintf("FindAll for '%s' in string '%s'", desc, tt.s), func(t *testing.T) {
		runner := NewRunner(m)
		runner.Reset()
		testFindAll(t, tt.s, runner, tt.expectedResults)
	})
}

func testFindAll(t *testing.T, s string, finiteStateMachine Machine, expectedResults []localResult) {
	results := FindAll(finiteStateMachine, s)

	if len(results) != len(expectedResults) {
		t.Fatalf("wrong number of results for string '%s', expected %+v, got %+v", s, len(expectedResults), len(results))
	}

	for j := range results {
		if results[j] != expectedResults[j] {
			t.Fatalf("wrong Result for string '%s': expected %d, got %d", s, expectedResults[j], results[j])
		}
	}
}

func testFindAllWithLines(t *testing.T, s string, finiteStateMachine Machine, expectedResults []Result) {
	results := FindAllWithLines(context2.TODO(), finiteStateMachine, s)

	if len(results) != len(expectedResults) {
		t.Fatalf("wrong number of results for string '%s', expected %+v, got %+v", s, len(expectedResults), len(results))
	}

	for j := range results {
		if results[j] != expectedResults[j] {
			t.Fatalf("wrong Result for string '%s': expected %d, got %d", s, expectedResults[j], results[j])
		}
	}
}
