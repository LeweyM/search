package finite_state_machine

type matchesLetter struct {
	base      finiteState
	letter    rune
	nextState finiteState
}

func newMatchesLetter(base finiteState, letter rune, nextState finiteState) *matchesLetter {
	return &matchesLetter{base: base, letter: letter, nextState: nextState}
}

func (m matchesLetter) test(r rune) finiteState {
	if r == m.letter {
		return m.nextState
	} else {
		return m.base
	}
}

func (m matchesLetter) isEndState() bool {
	return false
}
