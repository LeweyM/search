package v7

type Status string

const (
	Success Status = "success"
	Fail           = "fail"
	Normal         = "normal"
)

type State struct {
	transitions []Transition
	epsilons    []*State
}

func (s *State) matchingTransitions(input rune) []*State {
	var res []*State
	for _, t := range s.transitions {
		if t.predicate.test(input) {
			res = append(res, t.to)
		}
	}
	return res
}

func (s *State) isSuccessState() bool {
	if len(s.transitions) == 0 && len(s.epsilons) == 0 {
		return true
	}

	return false
}

// helper function to add a transition to State.
func (s *State) addTransition(destination *State, predicate Predicate, debugSymbol string) {
	t := Transition{
		debugSymbol: debugSymbol,
		to:          destination,
		from:        s,
		predicate:   predicate,
	}
	s.transitions = append(s.transitions, t)
}

func (s *State) addEpsilon(destination *State) {
	s.epsilons = append(s.epsilons, destination)
}

// adds the transitions of other State (s2) to this State (s).
func (s *State) merge(s2 *State) {
	for _, t := range s2.transitions {
		// 1. copy s2 transitions to s
		s.addTransition(t.to, t.predicate, t.debugSymbol)
	}

	// 2. remove s2
	s2.delete()
}

func (s *State) delete() {
	s.transitions = nil
	s.epsilons = nil
}
