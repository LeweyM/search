package finite_state_machine

type StateStack []*State

func (s *StateStack) pop() *State {
	i := len(*s) - 1
	state := (*s)[i]
	*s = (*s)[:i]
	return state
}

func (s *StateStack) peek() *State {
	i := len(*s) - 1
	return (*s)[i]
}

func (s *StateStack) push(linked *State) {
	*s = append(*s, linked)
}
