package v3

import (
	"regexp"
	"strings"
	"testing"
)

func TestCompiledFSM(t *testing.T) {
	parser := NewParser()

	tokens := lex("abc")
	ast := parser.Parse(tokens)
	startState, _ := ast.compile()

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
			testRunner := NewRunner(startState)

			for _, character := range tt.input {
				testRunner.Next(character)
			}

			result := testRunner.GetStatus()
			if tt.expectedStatus != result {
				t.Fatalf("Expected FSM to have final state of '%v', got '%v'", tt.expectedStatus, result)
			}

			goRegex := regexp.MustCompile("abc")
			match := goRegex.Match([]byte(tt.input))

			if tt.expectedStatus == Success {
				if !match {
					t.Fatalf("Expected FSM to have same result as GORegex pkg. Expected '%v', got '%v'", tt.expectedStatus, match)
				}
			} else {
				if match {
					t.Fatalf("Expected FSM to have same result as GORegex pkg. Expected '%v', got '%v'", tt.expectedStatus, match)
				}
			}
		})
	}
}

func TestFSMAgainstGoRegexPkg(t *testing.T) {
	type test struct {
		name  string
		regex string
		input string
	}

	tests := []test{
		{"empty string", "abc", ""},
		{"non matching string", "abc", "xxx"},
		{"matching string", "abc", "abc"},
		{"partial matching string", "abc", "ab"},

		{"nested expressions", "a(b(d))c", "abdc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchRegex(tt.regex, tt.input)

			goRegexMatch := regexp.MustCompile(tt.regex).MatchString(tt.input)

			if (result == Success && !goRegexMatch) || (result != Success && goRegexMatch) {
				t.Fatalf("Mismatch - Regex: '%s', Input: '%s' -> Go Regex Pkg: '%t', Our regex result: '%v'", tt.regex, tt.input, goRegexMatch, result)
			}
		})
	}
}

func FuzzFSM(f *testing.F) {
	f.Add("abc", "abc")
	f.Add("abc", "")
	f.Add("abc", "xxx")
	f.Add("ca(t)(s)", "dog")

	f.Fuzz(func(t *testing.T, regex, input string) {
		if strings.ContainsAny(input, "Ȥ") {
			t.Skip()
		}

		if strings.ContainsAny(regex, "$^|*+?.\\") {
			t.Skip()
		}

		compiledGoRegex, err := regexp.Compile(regex)
		if err != nil {
			t.Skip()
		}

		result := matchRegex(regex, input)
		goRegexMatch := compiledGoRegex.MatchString(input)

		if (result == Success && !goRegexMatch) || (result == Fail && goRegexMatch) {
			t.Fatalf("Mismatch - Regex: '%s', Input: '%s' -> Go Regex Pkg: '%t', Our regex result: '%v'", regex, input, goRegexMatch, result)
		}
	})
}

func matchRegex(regex, input string) Status {
	parser := NewParser()
	tokens := lex(regex)
	ast := parser.Parse(tokens)
	startState, _ := ast.compile()
	testRunner := NewRunner(startState)

	for _, character := range input {
		testRunner.Next(character)
		status := testRunner.GetStatus()
		if status == Fail {
			testRunner.Reset()
			continue
		}

		if status != Normal {
			return status
		}
	}

	return testRunner.GetStatus()
}
