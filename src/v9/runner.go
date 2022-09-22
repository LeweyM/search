package v9

type runner struct {
	head         *State
	activeStates Set[*State]
	heads        []*State
}

func NewRunner(head *State) *runner {
	r := &runner{
		head:         head,
		activeStates: NewSet[*State](connectedStates(head)...),
		heads:        append(connectedStates(head), head),
	}

	return r
}

func (r *runner) Next(input rune) {
	if r.activeStates.size() == 0 {
		return
	}

	r.advanceEpsilons()

	nextActiveStates := Set[*State]{}
	for activeState := range r.activeStates {
		for _, nextState := range activeState.matchingTransitions(input) {
			nextActiveStates.add(nextState)
		}
	}
	r.activeStates = nextActiveStates

	r.advanceEpsilons()

	//// remove isolated nodes
	//for _, activeState := range r.activeStates.list() {
	//	if len(activeState.transitions) == 0 && len(activeState.epsilons) == 0 && !activeState.isSuccessState() {
	//		activeState.delete()
	//		r.activeStates.remove(activeState)
	//	}
	//}
}

func (r *runner) GetStatus() Status {
	// if there are no actives states, return Fail
	if r.activeStates.size() == 0 {
		return Fail
	}

	// if any of the active states is a success state, return Success
	for state := range r.activeStates {
		if state.isSuccessState() {
			return Success
		}
	}

	// else, return normal
	return Normal
}

func (r *runner) Reset() {
	r.activeStates = NewSet(r.heads...)
	r.advanceEpsilons()
}

func (r *runner) Start() {
	for _, state := range r.heads {
		r.activateState(state)
	}
	r.advanceEpsilons()
}

func (r *runner) advanceEpsilons() {
	for state := range r.activeStates {
		r.activateConnectedEpsilons(state)
	}
}

func (r *runner) activateConnectedEpsilons(state *State) {
	for _, epsilon := range state.epsilons {
		if !r.activeStates.has(epsilon) {
			r.activateState(epsilon)
			r.activateConnectedEpsilons(epsilon)
		}
	}
}

func (r *runner) activateState(state *State) {
	//state.fullyReduceEpsilons3()
	//// remove isolated nodes
	//for _, activeState := range r.activeStates.list() {
	//	if len(activeState.transitions) == 0 && len(activeState.epsilons) == 0 && !activeState.isSuccessState() {
	//		activeState.delete()
	//		r.activeStates.remove(activeState)
	//	}
	//}
	r.activeStates.add(state)
}
