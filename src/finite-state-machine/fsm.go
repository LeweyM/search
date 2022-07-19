package finite_state_machine

type StateType int

const (
	Success StateType = iota
	Fail
	Normal
)

type Predicate func(input rune) bool

type destination *State

type Transition struct {
	// to: a pointer to the next state
	to destination
	// predicate: a function to determine if the runner should move to the next state
	predicate   Predicate
	description string
}

type State struct {
	id          int
	transitions []Transition
	epsilons    []Transition
	incoming    []Transition
}

func (s *State) matchingTransitions(input rune) []destination {
	var matchingTransitions []destination
	for _, t := range s.transitions {
		if t.predicate(input) {
			matchingTransitions = append(matchingTransitions, t.to)
		}
	}

	return matchingTransitions
}

func (s *State) isSuccessState() bool {
	if len(s.transitions) == 0 && len(s.epsilons) == 0 {
		return true
	}

	return false
}

func (s *State) addTransition(destination *State, predicate Predicate) {
	t := Transition{
		to:        destination,
		predicate: predicate,
	}
	s.transitions = append(s.transitions, t)
	destination.incoming = append(destination.incoming, t)
}

func (s *State) addEpsilonTransition(destination *State) {
	t := Transition{
		to:          destination,
		description: "epsilon",
	}
	s.epsilons = append(s.epsilons, t)
	destination.incoming = append(destination.incoming, t)
}
