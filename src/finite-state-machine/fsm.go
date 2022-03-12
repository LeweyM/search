package finite_state_machine

type StateType int

const (
	Success StateType = iota
	Fail
	Normal
)

type Predicate func(input rune) bool

type Transition struct {
	// to: a pointer to the next state
	to destination
	// predicate: a function to determine if the runner should move to the next state
	predicate   Predicate
	description string
	epsilon     bool
}

func NewEpsilon(to *State) Transition {
	return Transition{
		to:          to,
		predicate:   func(input rune) bool { return true },
		description: "epsilon",
		epsilon:     true,
	}
}

type State struct {
	empty        bool
	id          int
	transitions []Transition
}

type destination *State

func (s *State) matchingTransitions(input rune) []destination {
	var matchingTransitions []destination
	for _, t := range s.transitions {
		if t.predicate != nil && t.predicate(input) && !t.epsilon { // TODO: clean up epsilon check here
			matchingTransitions = append(matchingTransitions, t.to)
		}
	}
	return matchingTransitions
}

func (s *State) isSuccessState() bool {
	if len(s.transitions) == 0 {
		return true
	} else {
		// not efficient
		for _, linked := range s.transitions {
			if linked.to.empty {
				return true
			}
		}
		return false
	}
}

func (s *State) merge(s2 *State) {
	if s2.transitions[0].to.empty {
		s2.transitions = s2.transitions[1:]
	}
	for _, t := range s2.transitions {
		// when composing a transition, we merge the first transitions of the new state into the transition of the from state
		s.transitions = append(s.transitions, t)
	}
}
