package v8

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
		{"empty regex", "", "abc"},
		{"non matching string", "abc", "xxx"},
		{"matching string", "abc", "abc"},
		{"partial matching string", "abc", "ab"},
		{"nested expressions", "a(b(d))c", "abdc"},
		{"substring match with reset needed", "aA", "aaA"},
		{"substring match without reset needed", "B", "ABA"},
		{"multibyte characters", "Ȥ", "Ȥ"},
		{
			"complex multibyte characters",
			string([]byte{0xef, 0xbf, 0xbd, 0x30}),
			string([]byte{0xcc, 0x87, 0x30}),
		},
		// wildcard
		{"wildcard regex matching", "ab.", "abc"},
		{"wildcard regex not matching", "ab.", "ab"},
		{"wildcards matching newlines", "..0", "0\n0"},

		// branch
		{"branch matching 1st branch", "ab|cd", "ab"},
		{"branch matching 2nd branch", "ab|cd", "cd"},
		{"branch not matching", "ab|cd", "ac"},
		{"branch with shared characters", "dog|dot", "dog"},
		{"branch with shared characters", "dog|dot", "dot"},
		{"branch with right side empty", "a|", ""},
		{"branch with left side empty", "|a", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			myRegex := NewMyRegex(tt.regex)
			result := myRegex.MatchString(tt.input)

			goRegexMatch := regexp.MustCompile(tt.regex).MatchString(tt.input)
			t.Logf("Compiled state machine:\n%v\n", myRegex.DebugFSM())

			if result != goRegexMatch {
				t.Fatalf(
					"Mismatch - \nRegex: '%s' (as bytes: %x), \nInput: '%s' (as bytes: %x) \n-> \nGo Regex Pkg: '%t', \nOur regex result: '%v'",
					tt.regex,
					[]byte(tt.regex),
					tt.input,
					[]byte(tt.input),
					goRegexMatch,
					result,
				)
			}
		})
	}
}

func FuzzFSM(f *testing.F) {
	f.Add("ab|cd|ef", "abc")
	f.Add("abc", "abc")
	f.Add("abc", "")
	f.Add("abc", "xxx")
	f.Add("ca(t)(s)", "dog")

	f.Fuzz(func(t *testing.T, regex, input string) {
		if strings.ContainsAny(regex, "[]{}$^*+?\\") {
			t.Skip()
		}

		compiledGoRegex, err := regexp.Compile(regex)
		if err != nil {
			t.Skip()
		}

		result := NewMyRegex(regex).MatchString(input)
		goRegexMatch := compiledGoRegex.MatchString(input)

		if result != goRegexMatch {
			t.Fatalf(
				"Mismatch - \nRegex: '%s' (as bytes: %x), \nInput: '%s' (as bytes: %x) \n-> \nGo Regex Pkg: '%t', \nOur regex result: '%v'",
				regex,
				[]byte(regex),
				input,
				[]byte(input),
				goRegexMatch,
				result)
		}
	})
}
