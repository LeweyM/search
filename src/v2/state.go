package v2

import "fmt"

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
	predicate Predicate
}

type State struct {
	id          int
	transitions []Transition
	incoming    []*State
}

func (s *State) firstMatchingTransition(input rune) destination {
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
	destination.incoming = append(destination.incoming, s)
}

// adds the transitions of other State (s2) to this State (s).
//
// warning: do not use if State s2 has any incoming transitions.
func (s *State) merge(s2 *State) {
	if len(s2.incoming) != 0 {
		panic(fmt.Sprintf("State (%+v) cannot be merged if it has any incoming transitions. It has incoming transitions from the following states; %+v", *s2, s.incoming))
	}

	for _, t := range s2.transitions {
		s.addTransition(t.to, t.predicate)
	}
}