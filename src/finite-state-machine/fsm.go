package finite_state_machine

type Status string

const (
	Success Status = "success"
	Fail           = "fail"
	Normal         = "normal"
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
}

func (s *State) addEpsilonTransition(destination *State) {
	t := Transition{
		to:          destination,
		description: "epsilon",
	}
	s.epsilons = append(s.epsilons, t)
}

/*
	-->(s)   (s2)-->(3)
becomes
	-->(s)-->(3)
	(s2)

note: this does not account for incoming transitions to (s2), these will
still point to s2 if they exist.

this should only be used when (s2) has no incoming transitions
*/
func (s *State) merge(s2 *State) {
	s2.checkHasNoIncomingTransitions(s2)

	for _, t := range s2.transitions {
		s.addTransition(t.to, t.predicate)
	}

	for _, t := range s2.epsilons {
		s.addEpsilonTransition(t.to)
	}
}

func (s *State) checkHasNoIncomingTransitions(target *State) {
	for _, t := range append(s.transitions, s.epsilons...) {
		if t.to == target {
			panic("state should have no incoming transitions")
		}

		(*State)(t.to).checkHasNoIncomingTransitions(target)
	}
}
