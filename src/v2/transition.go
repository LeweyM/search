package v2

type Predicate func(input rune) bool

type Transition struct {
	to        *State
	from      *State
	predicate Predicate
}