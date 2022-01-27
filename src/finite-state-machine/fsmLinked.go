package finite_state_machine

type letterTransitionLinked struct {
	to    *stateLinked
	input rune
}

type wildTransitionLinked struct {
	to *stateLinked
}

type stateLinked struct {
	id              int
	transitions     []letterTransitionLinked
	wildTransitions []wildTransitionLinked
	stateType       StateType
}

type destination *stateLinked

func (s *stateLinked) matchingTransitions(input rune) []destination {
	var matchingTransitions []destination
	for _, t := range s.transitions {
		if t.input == input {
			matchingTransitions = append(matchingTransitions, t.to)
		}
	}
	for _, t := range s.wildTransitions {
		matchingTransitions = append(matchingTransitions, t.to)
	}
	return append(matchingTransitions)
}

type branch struct {
	currState *stateLinked
}

type runner struct {
	head      *stateLinked
	failState *stateLinked
	branches  []branch
}

func NewRunner(head *stateLinked) *runner {
	failState := &stateLinked{id: 0, stateType: Fail}

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
	states []*stateLinked
}

func NewStateLinkedBuilder(n int) *builder {
	var states []*stateLinked
	states = append(states, &stateLinked{id: 0, stateType: Fail}) // stand in for fail state

	for i := 1; i <= n; i++ {
		states = append(states, &stateLinked{id: i, stateType: Normal})
	}
	return &builder{states: states}
}

func (b *builder) AddTransition(from, to int, letter rune) *builder {
	if from >= len(b.states) || to >= len(b.states) {
		panic("Cannot set a transition for a state outside of range")
	}
	b.states[from].transitions = append(b.states[from].transitions, letterTransitionLinked{
		to:    b.states[to],
		input: letter,
	})
	return b
}

func (b *builder) SetSuccess(n int) *builder {
	b.states[n].stateType = Success
	return b
}

func (b *builder) Build() *stateLinked {
	return b.states[1]
}

func (b *builder) AddWildTransition(from, to int) *builder {
	if from >= len(b.states) || to >= len(b.states) {
		panic("Cannot set a transition for a state outside of range")
	}
	b.states[from].wildTransitions = append(b.states[from].wildTransitions, wildTransitionLinked{
		to: b.states[to],
	})
	return b
}
