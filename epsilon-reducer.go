package v10

// epsilonReducer will turn an epsilon-NFA into an NFA. It does this by collecting the transitions
// of all the states in a given state's epsilon closure (the set of states connected by epsilons)
// and applying those to state.
type epsilonReducer struct{}

func (e *epsilonReducer) reduce(s *State) {
	states := NewSet[*State]()
	e.reduceEpsilons(s, &states)
}

func (e *epsilonReducer) reduceEpsilons(s *State, visited *Set[*State]) {
	// if this state has already been reduced, return.
	if visited.has(s) {
		return
	}
	visited.add(s)

	// collect all the transitions of the current states clojure.
	closure := s.getEpsilonClosure()
	closureTransitions := collectTransitions(closure)

	// remove the current states transitions and epsilons.
	s.transitions, s.epsilons = nil, nil

	// replace the current states transitions with all the transitions of the closure.
	for _, t := range closureTransitions {
		s.addTransition(t.to, t.predicate, t.debugSymbol)
	}

	// if any of the states in the closure was a success state, make the current state a success state.
	for state := range closure {
		if state.isSuccessState() {
			s.SetSuccess()
		}
	}

	// recur on the states connected by the new transitions.
	for _, transition := range s.transitions {
		e.reduceEpsilons(transition.to, visited)
	}
}

func collectTransitions(states Set[*State]) []Transition {
	var transitions []Transition
	for state := range states {
		transitions = append(transitions, state.transitions...)
	}
	return transitions
}
