package finite_state_machine

type matchesLetter struct {
	base      finiteState
	letter    rune
	nextState finiteState
	endState  bool
}

func (m matchesLetter) test(r rune) finiteState {
	if r == m.letter {
		return m.nextState
	} else {
		return m.base
	}
}

func (m matchesLetter) isEndState() bool {
	return m.endState
}

func (m matchesLetter) End() matchesLetter {
	return matchesLetter{
		base:      m.base,
		letter:    m.letter,
		nextState: m.nextState,
		endState:  true,
	}
}

func (m matchesLetter) Base(state finiteState) matchesLetter {
	return matchesLetter{
		base:      state,
		letter:    m.letter,
		nextState: m.nextState,
		endState:  m.endState,
	}
}