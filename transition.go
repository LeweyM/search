package v5

import (
	"strings"
)

type Transition struct {
	debugSymbol string
	// to: a pointer to the next state
	to        *State
	from      *State
	predicate Predicate
}

type Predicate struct {
	allowedChars    string
	disallowedChars string
}

func (p Predicate) test(input rune) bool {
	if p.allowedChars != "" && p.disallowedChars != "" {
		panic("must be mutually exclusive")
	}

	if len(p.allowedChars) > 0 {
		return strings.ContainsRune(p.allowedChars, input)
	}
	if len(p.disallowedChars) > 0 {
		return !strings.ContainsRune(p.disallowedChars, input)
	}
	return false
}
