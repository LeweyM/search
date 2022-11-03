package v5

import (
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compareWithGoStdLib(t, NewMyRegex(tt.regex), tt.regex, tt.input)
		})
	}
}

func FuzzFSM(f *testing.F) {
	f.Add("abc", "abc")
	f.Add("abc", "")
	f.Add("abc", "xxx")
	f.Add("ca(t)(s)", "dog")

	f.Fuzz(func(t *testing.T, regex, input string) {
		if strings.ContainsAny(regex, "[]{}$^|*+?\\") {
			t.Skip()
		}

		_, err := regexp.Compile(regex)
		if err != nil {
			t.Skip()
		}

		compareWithGoStdLib(t, NewMyRegex(regex), regex, input)
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
