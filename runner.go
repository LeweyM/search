package v7

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

	nextActiveStates := Set[*State]{}
	for activeState := range r.activeStates {
		for _, nextState := range activeState.matchingTransitions(input) {
			nextActiveStates.add(nextState)
		}
	}
	r.activeStates = nextActiveStates
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
}

func (r *runner) Start() {
	r.activeStates.add(r.head)
}
