package v3

type Status string

const (
	Success Status = "success"
	Fail           = "fail"
	Normal         = "normal"
)

type State struct {
	id          int
	transitions []Transition
}

func (s *State) firstMatchingTransition(input rune) *State {
	for _, t := range s.transitions {
		if t.predicate(input) {
			return t.to
		}
	}

	return nil
}

func (s *State) isSuccessState() bool {
	if len(s.transitions) == 0 {
		return true
	}

	return false
}

// helper function to add a transition to State.
func (s *State) addTransition(destination *State, predicate Predicate) {
	t := Transition{
		to:        destination,
		from:      s,
		predicate: predicate,
	}
	s.transitions = append(s.transitions, t)
}

// adds the transitions of other State (s2) to this State (s).
func (s *State) merge(s2 *State) {
	for _, t := range s2.transitions {
		s.addTransition(t.to, t.predicate)
	}
}
