package v8

type runner struct {
	head         *State
	activeStates Set[*State]
}

func NewRunner(head *State) *runner {
	r := &runner{
		head:         head,
		activeStates: NewSet[*State](head),
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
	r.activeStates = NewSet[*State](r.head)
	r.advanceEpsilons()
}

func (r *runner) Start() {
	r.activeStates.add(r.head)
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
			r.activeStates.add(epsilon)
			r.activateConnectedEpsilons(epsilon)
		}
	}
}
