package finite_state_machine

type fsm struct {
	head finiteState
	current finiteState
}

func NewFsm(head finiteState) *fsm {
	return &fsm{head: head, current: head}
}

type finiteState interface {
	test(r rune) finiteState
	isEndState() bool
}

func (f *fsm) next(r rune) bool {
	f.current = f.current.test(r)
	isInEndState := f.current.isEndState()
	if isInEndState {
		f.current = f.head
	}
	return isInEndState
}

type startingState struct {
	nextState finiteState
}

func (s startingState) test(r rune) finiteState {
	return s.nextState.test(r)
}

func (s startingState) isEndState() bool {
	return false
}
