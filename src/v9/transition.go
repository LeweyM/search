package v9

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

func (t Transition) invert() Transition {
	return Transition{
		debugSymbol: t.debugSymbol,
		to:          t.from,
		from:        t.to,
		predicate:   t.predicate,
	}
}

func (t Transition) From(s *State) Transition {
	t.from = s
	return t
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
