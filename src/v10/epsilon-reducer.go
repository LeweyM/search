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
	// 0. if this state has already been reduced, return to avoid infinite recursive loops.
	if visited.has(s) {
		return
	}
	visited.add(s)

	// 1. Collect the states of the epsilon closure.
	closure := s.getEpsilonClosure()

	// 2. Collect the transitions of all states within the closure.
	closureTransitions := collectTransitions(closure)

	// 3. Remove any epsilon transitions from the state.
	s.epsilons = nil

	// 4. Replace the transitions of the state with the closure transitions.
	s.transitions = nil
	for _, t := range closureTransitions {
		s.addTransition(t.to, t.predicate, t.debugSymbol)
	}

	// 5.  If any state in the closure is a success state, make the state a success state.
	for state := range closure {
		if state.isSuccessState() {
			s.SetSuccess()
			break
		}
	}

	// 6. Recur on connected states.
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
