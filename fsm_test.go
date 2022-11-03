package v10

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

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

		// group
		{"word followed by group", "1(|)", "0"},
		{"empty group concatenation", "(()0)0", "0"},
		{"group followed by word", "(|)1", "0"},

		// zero or one
		{"simple zero or one with 0 '?' match", "ab?c", "ac"},
		{"simple zero or one with one '?' matches", "ab?c", "abc"},
		{"simple zero or one too many '?' matches", "ab?c", "abbc"},

		// one or more
		{"simple one or more with 0 '+' matches", "ab+c", "ac"},
		{"simple one or more with one '+' matches", "ab+c", "abc"},
		{"simple one or more with many '+' matches", "ab+c", "abbbbc"},

		// zero or more
		{"simple zero or more with 0 '*' matches", "ab*c", "ac"},
		{"simple zero or more with one '*' matches", "ab*c", "abc"},
		{"simple zero or more with many '*' matches", "ab*c", "abbbbc"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("without reducers - %s", tt.name), func(t *testing.T) {
			compareWithGoStdLib(t, NewMyRegex(tt.regex), tt.regex, tt.input)
		})
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("with epsilon reducer - %s", tt.name), func(t *testing.T) {
			compareWithGoStdLib(t, NewMyRegex(tt.regex, &epsilonReducer{}), tt.regex, tt.input)
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
		if strings.ContainsAny(regex, "[]{}$^\\") {
			t.Skip()
		}

		if strings.Contains(regex, "(?") {
			// '?' on its own is used for special group constructs, which we're not implementing.
			t.Skip()
		}

		_, err := regexp.Compile(regex)
		if err != nil {
			t.Skip()
		}
		compareWithGoStdLib(t, NewMyRegex(regex, &epsilonReducer{}), regex, input)
	})
}

func compareWithGoStdLib(t *testing.T, myRegex *myRegex, regex, input string) {
	t.Helper()

	result := myRegex.MatchString(input)
	goRegexMatch := regexp.MustCompile(regex).MatchString(input)

	if result != goRegexMatch {
		t.Fatalf(
			"Mismatch - \nRegex: '%s' (as bytes: %x), \nInput: '%s' (as bytes: %x) \n-> \nGo Regex Pkg: '%t', \nOur regex result: '%v'",
			regex,
			[]byte(regex),
			input,
			[]byte(input),
			goRegexMatch,
			result,
		)
	}
}
