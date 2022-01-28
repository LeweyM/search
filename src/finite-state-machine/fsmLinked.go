package finite_state_machine

type Predicate func(input rune) bool

type transitionLinked struct {
	// to: a pointer to the next state
	to destination
	// predicate: a function to determine if the runner should move to the next state
	predicate Predicate
}

type StateLinked struct {
	id           int
	transitions1 []transitionLinked
	stateType    StateType
}

type destination *StateLinked

func (s *StateLinked) matchingTransitions(input rune) []destination {
	var matchingTransitions []destination
	for _, t := range s.transitions1 {
		if t.predicate(input) {
			matchingTransitions = append(matchingTransitions, t.to)
		}
	}
	return append(matchingTransitions)
}

type branch struct {
	currState *StateLinked
}

type runner struct {
	head      *StateLinked
	failState *StateLinked
	branches  []branch
}

func NewRunner(head *StateLinked) *runner {
	failState := &StateLinked{id: 0, stateType: Fail}

	return &runner{
		failState: failState,
		head:      head,
		branches:  []branch{{currState: head}},
	}
}

func (r *runner) Next(input rune) StateType {
	activeBranches := r.activeBranches()

	for bIndex, b := range activeBranches {
		matchingTransitions := b.currState.matchingTransitions(input)
		// if no transitions are possible, the branch has failed
		if len(matchingTransitions) == 0 {
			r.branches[bIndex].currState = r.failState
			continue
		}
		// if there is only one transition, we move
		if len(matchingTransitions) == 1 {
			r.branches[bIndex].currState = matchingTransitions[0]
			continue
		}
		// if there are multiple transitions, we branch
		r.branches[bIndex].currState = matchingTransitions[0]
		for i := 1; i < len(matchingTransitions); i++ {
			// TODO: Trim duplicate branches here. Maybe store as a set of 'to' pointers
			r.branches = append(r.branches, branch{
				currState: matchingTransitions[i],
			})
		}
	}

	// recount active branches here as the number could have changed
	activeBranches = r.activeBranches()
	if len(activeBranches) == 0 {
		return Fail
	}
	for _, b := range activeBranches {
		if b.currState.stateType == Success {
			return Success
		}
	}
	return Normal
}

func (r *runner) Reset() {
	r.branches = []branch{{currState: r.head}}
}

func (r *runner) activeBranches() []branch {
	var activeBranches []branch
	for _, b := range r.branches {
		if b.currState.stateType == Normal || b.currState.stateType == Success {
			activeBranches = append(activeBranches, b)
		}
	}
	return activeBranches
}

type builder struct {
	states []*StateLinked
}

var GlobalIdCounter = 0

func NewStateLinkedBuilder(n int) *builder {
	var states []*StateLinked
	states = append(states, &StateLinked{id: 0, stateType: Fail}) // stand in for fail state

	for i := 1; i <= n; i++ {
		GlobalIdCounter++
		states = append(states, &StateLinked{id: GlobalIdCounter, stateType: Normal})
	}
	return &builder{states: states}
}

func (b *builder) AddTransition(from, to int, letter rune) *builder {
	if from >= len(b.states) || to >= len(b.states) {
		panic("Cannot set a transition for a state outside of range")
	}
	b.states[from].transitions1 = append(b.states[from].transitions1, transitionLinked{
		to:        b.states[to],
		predicate: func(input rune) bool { return input == letter },
	})
	return b
}

func (b *builder) AddWildTransition(from, to int) *builder {
	if from >= len(b.states) || to >= len(b.states) {
		panic("Cannot set a transition for a state outside of range")
	}
	b.states[from].transitions1 = append(b.states[from].transitions1, transitionLinked{
		to:        b.states[to],
		predicate: func(input rune) bool { return true },
	})
	return b
}

func (b *builder) AddMachineTransition(from int, state *StateLinked) *builder {
	if from >= len(b.states) {
		panic("Cannot set a transition for a state outside of range")
	}
	for _, t := range state.transitions1 {
		// when composing a transition, we merge the first transitions of the new state into the transition of the from state
		b.states[from].transitions1 = append(b.states[from].transitions1, transitionLinked{
			to:        t.to,
			predicate: t.predicate,
		})
	}
	return b
}

func (b *builder) SetSuccess(n int) *builder {
	b.states[n].stateType = Success
	return b
}

func (b *builder) Build() *StateLinked {
	return b.states[1]
}
