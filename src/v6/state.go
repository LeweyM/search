package v6

import "fmt"

type Status string

const (
	Success Status = "success"
	Fail           = "fail"
	Normal         = "normal"
)

type State struct {
	transitions []Transition
	incoming    []*State
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
	if len(s.transitions) == 0 {
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
		// 1. copy s2 transitions to s
		s.addTransition(t.to, t.predicate, t.debugSymbol)
	}

	// 2. remove s2
	s2.delete()
}

func (s *State) delete() {
	// 1. remove s from incoming of connected nodes.
	for _, t := range s.transitions {
		t.to.removeIncoming(s)
	}

	// 2. remove the outgoing transitions
	s.transitions = nil
}

func (s *State) removeIncoming(target *State) {
	var newIncoming []*State
	for _, state := range s.incoming {
		if target != state {
			newIncoming = append(newIncoming, state)
		}
	}
	s.incoming = newIncoming
}