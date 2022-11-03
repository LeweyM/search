package v10

type Status string

const (
	Success Status = "success"
	Fail           = "fail"
	Normal         = "normal"
)

type State struct {
	transitions []Transition
	epsilons    []*State
	success     bool
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
	return s.success
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

func (s *State) SetSuccess() {
	s.success = true
}

func (s *State) getEpsilonClosure() Set[*State] {
	set := NewSet[*State](s)

	s.traverseEpsilons(set)

	return set
}

func (s *State) traverseEpsilons(states Set[*State]) {
	for _, state := range s.epsilons {
		if !states.has(state) {
			states.add(state)
			state.traverseEpsilons(states)
		}
	}
}
