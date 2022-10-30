package v2

type Status string

const (
	Success Status = "success"
	Fail           = "fail"
	Normal         = "normal"
)

type Predicate func(input rune) bool

type Transition struct {
	// to: a pointer to the next state
	to *State
	// predicate: a function to determine if the runner should move to the next state
	predicate Predicate
}

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
		predicate: predicate,
	}
	s.transitions = append(s.transitions, t)
}

// adds the transitions of other State (s2) to this State (s).
func (s *State) merge(s2 *State) {
	for _, t := range s2.transitions {
		// 1. copy s2 transitions to s
		s.addTransition(t.to, t.predicate)
	}

	// 2. remove s2
	s2.delete()
}

func (s *State) delete() {
	s.transitions = nil
}
