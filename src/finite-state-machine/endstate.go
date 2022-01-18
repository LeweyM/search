package finite_state_machine

type endState struct {
}

func (m endState) test(r rune) finiteState {
	return m
}

func (m endState) isEndState() bool {
	return true
}
