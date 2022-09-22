package v9

import "fmt"

type Status string

const (
	Success Status = "success"
	Fail           = "fail"
	Normal         = "normal"
)

type State struct {
	transitions []Transition
	epsilons    []*State
	incoming    []*State
	isSuccess   bool
}

func (s *State) matchingTransitions(input rune) []*State {
	var res []*State
	for _, t := range s.transitions {
		// if there is a match
		if t.predicate.test(input) {
			// activate the destination
			res = append(res, t.to)
			// reduce epsilon transitions at destination
			//res = append(res, t.to.fullyReduceEpsilons()...)
		}
	}
	return res
}

// represents a single transition which, when extended by epsilons, leads to a destination
type journey struct {
	origin, destination *State
	Transition          Transition
}

func (s *State) fullyReduceEpsilons() []*State {
	var res []*State
	// for each of the epsilons destinations
	for _, epsilonDestination := range s.epsilons {

		//get all journeys going to the destination
		journeysToEpsilonDestination := journeysTo(epsilonDestination, NewSet[*State]())

		for _, j := range journeysToEpsilonDestination {
			// add transition from start of each journey to destination
			j.origin.addTransition(j.destination, j.Transition.predicate, j.Transition.debugSymbol)
			// and activate each journey
			res = append(res, j.destination)
		}

		// and finally remove the epsilon
		s.removeEpsilon(epsilonDestination)
	}

	return res
}

func (s *State) fullyReduceEpsilons3() {
	//get all journeys going to the destination
	journeysToEpsilonDestination := journeysTo(s, NewSet[*State]())

	for _, j := range journeysToEpsilonDestination {
		// add transition from start of each journey to destination
		j.origin.addTransition(s, j.Transition.predicate, j.Transition.debugSymbol)
	}

	// and finally remove all epsilons to state
	for _, epsilonOrigin := range s.incomingEpsilons() {
		epsilonOrigin.removeEpsilon(s)
		// remove isolated nodes
		if len(epsilonOrigin.transitions) == 0 && len(epsilonOrigin.epsilons) == 0 && !epsilonOrigin.isSuccessState() {
			epsilonOrigin.delete()
			// note: this might leave some active states even after they have been removed.
			// example solution:
			// todo: return list of deleted states so they can be removed from active list.
		}
	}
}

func (s *State) fullyReduceEpsilons2() {
	for _, epsilonOrigin := range s.incomingEpsilons() {

		//get all journeys going to the destination
		journeysToEpsilonDestination := journeysTo(epsilonOrigin, NewSet[*State]())

		for _, j := range journeysToEpsilonDestination {
			// add transition from start of each journey to destination
			j.origin.addTransition(s, j.Transition.predicate, j.Transition.debugSymbol)
		}

		// and finally remove the epsilon
		epsilonOrigin.removeEpsilon(s)
		// remove isolated nodes
		if len(epsilonOrigin.transitions) == 0 && len(epsilonOrigin.epsilons) == 0 && !epsilonOrigin.isSuccessState() {
			epsilonOrigin.delete()
			// note: this might leave some active states even after they have been removed.
			// example solution:
			// todo: return list of deleted states so they can be removed from active list.
		}
	}
}

func journeysTo(destination *State, visited Set[*State]) []journey {
	if visited.has(destination) {
		return nil
	}
	visited.add(destination)

	var res []journey
	for _, s := range destination.incoming {
		for _, t := range s.transitions {
			if t.to == destination {
				res = append(res, journey{
					origin:      t.from,
					destination: destination,
					Transition: Transition{
						debugSymbol: t.debugSymbol,
						to:          destination,
						from:        s,
						predicate:   t.predicate,
					},
				})
			}
		}

		for _, e := range s.epsilons {
			if e == destination {
				journeysToEpsilon := journeysTo(s, visited)
				for _, j := range journeysToEpsilon {
					res = append(res, journey{
						origin:      j.origin,
						destination: destination,
						Transition: Transition{
							debugSymbol: j.Transition.debugSymbol,
							to:          destination,
							from:        j.origin,
							predicate:   j.Transition.predicate,
						},
					})
				}
			}
		}
	}

	return res
}

func filter[T comparable](res []T, to T) []T {
	var newS []T
	for _, s := range res {
		if s != to {
			newS = append(newS, s)
		}
	}
	return newS
}

func recurGetTerminalState(set Set[*State], s *State) {
	if set.has(s) {
		return
	}
	set.add(s)

	for _, epsilon := range s.epsilons {
		recurGetTerminalState(set, epsilon)
	}
}

func connectedStates(epsilon *State) []*State {
	s := Set[*State]{}

	recurGetTerminalState(s, epsilon)

	return s.list()
}

func (s *State) isSuccessState() bool {
	return s.isSuccess
}

// helper function to add a transition to State.
func (s *State) addTransition(destination *State, predicate Predicate, debugSymbol string) {
	t := Transition{
		debugSymbol: debugSymbol,
		to:          destination,
		from:        s,
		predicate:   predicate,
	}

	// todo: optimize. Set?
	for _, tr := range s.transitions {
		if tr == t {
			return
		}
	}
	s.transitions = append(s.transitions, t)
	s.addIncoming(destination)
}

func (s *State) addIncoming(destination *State) {
	destination.incoming = append(destination.incoming, s)
}

func (s *State) addEpsilon(destination *State) {
	s.epsilons = append(s.epsilons, destination)
	s.addIncoming(destination)
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

	for _, incoming := range s.incoming {
		incoming.removeEpsilon(s)
		incoming.removeTransition(s)
	}

	// 2. remove the outgoing transitions
	s.transitions = nil
	s.epsilons = nil
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

func (s *State) removeEpsilon(target *State) {
	newEpsilons := []*State{}

	for _, epsilon := range s.epsilons {
		if epsilon != target {
			newEpsilons = append(newEpsilons, epsilon)
		} else {
			target.removeIncoming(epsilon)
		}
	}

	s.epsilons = newEpsilons
}

func (s *State) removeTransition(s2 *State) {
	newTransitions := []Transition{}
	for _, transition := range s.transitions {
		if transition.to != s2 {
			newTransitions = append(newTransitions, transition)
		}
	}
	s.transitions = newTransitions
}

func (s *State) setAsSuccess() {
	s.isSuccess = true
}

func (s *State) incomingEpsilons() []*State {
	var res []*State
	for _, state := range s.incoming {
		for _, epsilon := range state.epsilons {
			if epsilon == s {
				res = append(res, state)
			}
		}
	}
	return res
}
