package finite_state_machine

import (
	"search/src/ast"
	"testing"
)

type testNew struct {
	desc            string
	regex           string
	searchString    string
	expectedResults []localResult
}

func TestNewCompiler(t *testing.T) {
	tests := []testNew{
		// concatenation
		{desc: "simple string matching", regex: "aaa", searchString: "aaa", expectedResults: []localResult{{0, 2}}},
		{desc: "simple string not matching", regex: "aaa", searchString: "aab"},
		// wildcard characters '.'
		{desc: "a.b with empty", regex: "a.b", searchString: "", expectedResults: nil},
		{desc: "a.b with 'ab'", regex: "a.b", searchString: "ab", expectedResults: nil},
		{desc: "a.b with 'azb'", regex: "a.b", searchString: "azb", expectedResults: []localResult{{0, 2}}},
		{desc: "a.b with 'acb'", regex: "a.b", searchString: "acb", expectedResults: []localResult{{0, 2}}},
		{desc: "a.b with 'azzzb'", regex: "a.b", searchString: "azzzb", expectedResults: nil},
		// branching
		{desc: "branch matching first branch", regex: "cat|dog", searchString: "cat", expectedResults: []localResult{{0, 2}}},
		{desc: "branch matching second branch", regex: "cat|dog", searchString: "dog", expectedResults: []localResult{{0, 2}}},
		{desc: "abc|def|xyz with 'abc'", regex: "abc|def|xyz", searchString: "abc", expectedResults: []localResult{{0, 2}}},
		{desc: "abc|def|xyz with 'def'", regex: "abc|def|xyz", searchString: "def", expectedResults: []localResult{{0, 2}}},
		{desc: "abc|def|xyz with 'xyz'", regex: "abc|def|xyz", searchString: "xyz", expectedResults: []localResult{{0, 2}}},
		{desc: "abc|abx|aby|abz with 'abz'", regex: "abc|abx|aby|abz", searchString: "abz", expectedResults: []localResult{{0, 2}}},
		{desc: "abc|abx|aby|abz with 'abc'", regex: "abc|abx|aby|abz", searchString: "abc", expectedResults: []localResult{{0, 2}}},
		{desc: "abc|abx|aby|abz with 'abr'", regex: "abc|abx|aby|abz", searchString: "abr"},
		// *
		{desc: "a*b with 'ab'", regex: "a*b", searchString: "ab", expectedResults: []localResult{{0, 1}}},
		{desc: "a*b with 'aab'", regex: "a*b", searchString: "aab", expectedResults: []localResult{{0, 2}}},
		{desc: "a*b with 'aaab'", regex: "a*b", searchString: "aaab", expectedResults: []localResult{{0, 3}}},
		{desc: "a*b with 'b'", regex: "a*b", searchString: "b", expectedResults: []localResult{{0, 0}}},
		{desc: "a*b with 'bb'", regex: "a*b", searchString: "bb", expectedResults: []localResult{{0, 0}, {1, 1}}},
		{desc: "a*b with 'a'", regex: "a*b", searchString: "a"},
		{desc: "a*b with 'aa'", regex: "a*b", searchString: "aa"},
		{desc: "matching 0 'c's", regex: "abc*", searchString: "ab", expectedResults: []localResult{{0, 1}}},
		{desc: "matching 1 'c'", regex: "abc*", searchString: "abc", expectedResults: []localResult{{0, 2}}},
		{desc: "matching 2 'c's", regex: "abc*", searchString: "abcc", expectedResults: []localResult{{0, 3}}},
		// +
		{desc: "a+b with 'ab'", regex: "a+b", searchString: "ab", expectedResults: []localResult{{0, 1}}},
		{desc: "a+b with 'aab'", regex: "a+b", searchString: "aab", expectedResults: []localResult{{0, 2}}},
		{desc: "a+b with 'aaab'", regex: "a+b", searchString: "aaab", expectedResults: []localResult{{0, 3}}},
		{desc: "a+b with 'aazb'", regex: "a+b", searchString: "aazb"},
		{desc: "a+b with 'a'", regex: "a+b", searchString: "a"},
		{desc: "a+b with 'b'", regex: "a+b", searchString: "b"},
		{desc: "xy+ with 'x'", regex: "xy+", searchString: "x"},
		{desc: "xy+ with 'y'", regex: "xy+", searchString: "y"},
		{desc: "xy+ with 'xy'", regex: "xy+", searchString: "xy", expectedResults: []localResult{{0, 1}}},
		{desc: "xy+ with 'xyxy'", regex: "xy+", searchString: "xyxy", expectedResults: []localResult{{0, 1}, {2, 3}}},
		{desc: "xy+ with 'xyyy'", regex: "xy+", searchString: "xyyy", expectedResults: []localResult{{0, 3}}},
		// ?
		{desc: "a?b with 'ab'", regex: "a?b", searchString: "ab", expectedResults: []localResult{{0, 1}}},
		{desc: "a?b with 'b'", regex: "a?b", searchString: "b", expectedResults: []localResult{{0, 0}}},
		{desc: "a?b with 'a'", regex: "a?b", searchString: "a"},
		{desc: "cats? with 'cat'", regex: "cats?", searchString: "cat", expectedResults: []localResult{{0, 2}}},
		{desc: "cats? with 'cats'", regex: "cats?", searchString: "cats", expectedResults: []localResult{{0, 3}}},
		{desc: "held?p? with 'held'", regex: "held?p?", searchString: "held", expectedResults: []localResult{{0, 3}}},
		{desc: "held?p? with 'help'", regex: "held?p?", searchString: "help", expectedResults: []localResult{{0, 3}}},
		{desc: "held?p? with 'hel'", regex: "held?p?", searchString: "hel.", expectedResults: []localResult{{0, 2}}},
		{desc: "held?p? with 'helt'", regex: "held?p?", searchString: "helt", expectedResults: []localResult{{0, 2}}},
		{desc: "ab?c? with 'abz'", regex: "ab?c?", searchString: "abz", expectedResults: []localResult{{0, 1}}},
		// groups
		{desc: "(l|L)et with 'light'", regex: "(l|L)et", searchString: "light"},
		{desc: "(l|L)et with 'let'", regex: "(l|L)et", searchString: "let", expectedResults: []localResult{{0, 2}}},
		{desc: "(l|L)et with 'Let'", regex: "(l|L)et", searchString: "Let", expectedResults: []localResult{{0, 2}}},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			parser := ast.Parser{}
			tree := parser.Parse(tt.regex)
			compiledMachine := CompileNEW(tree)
			runner := NewRunner(compiledMachine)
			runner.Reset()
			testFindAll(t, tt.searchString, runner, tt.expectedResults)
		})
	}
}
