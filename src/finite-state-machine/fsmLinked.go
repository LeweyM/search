package finite_state_machine

import "fmt"

type Predicate func(input rune) bool

type transitionLinked struct {
	// to: a pointer to the next state
	to destination
	// predicate: a function to determine if the runner should move to the next state
	predicate   Predicate
	description string
	epsilon     bool
}

type StateLinked struct {
	empty        bool
	id           int
	transitions1 []transitionLinked
}

type destination *StateLinked

func (s *StateLinked) matchingTransitions(input rune) []destination {
	var matchingTransitions []destination
	for _, t := range s.transitions1 {
		if t.predicate != nil && t.predicate(input) {
			matchingTransitions = append(matchingTransitions, t.to)
		}
	}
	return append(matchingTransitions)
}

func (s *StateLinked) isSuccessState() bool {
	if len(s.transitions1) == 0 {
		return true
	} else {
		// not efficient
		for _, linked := range s.transitions1 {
			if linked.to.empty {
				return true
			}
		}
		return false
	}
}

func (s *StateLinked) merge(s2 *StateLinked) {
	if s2.transitions1[0].to.empty {
		s2.transitions1 = s2.transitions1[1:]
	}
	for _, t := range s2.transitions1 {
		// when composing a transition, we merge the first transitions of the new state into the transition of the from state
		s.transitions1 = append(s.transitions1, transitionLinked{
			description: t.description,
			to:          t.to,
			predicate:   t.predicate,
		})
	}
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
	failState := &StateLinked{id: 0}

	return &runner{
		failState: failState,
		head:      head,
		branches:  []branch{{currState: head}},
	}
}

func (r *runner) Next(input rune) StateType {
	// move along epsilon transitions first.
	// This is probably inefficient and could be moved into the main loop.
	r.processEpsilons(input, r.branches)

	// move along regular transitions
	var nonFailedBranches []branch

	for bIndex, b := range r.branches {
		matchingTransitions := b.currState.matchingTransitions(input)
		// if no transitions are possible, the branch has failed
		if len(matchingTransitions) == 0 {
			// remove failed branch
			continue
		}
		// if there is only one transition, we move
		if len(matchingTransitions) == 1 {
			br := r.branches[bIndex]
			br.currState = matchingTransitions[0]
			nonFailedBranches = append(nonFailedBranches, br)
			continue
		}
		// if there are multiple transitions, we branch
		br := r.branches[bIndex]
		br.currState = matchingTransitions[0]
		nonFailedBranches = append(nonFailedBranches, br)

		for i := 1; i < len(matchingTransitions); i++ {
			// TODO: Trim duplicate branches here. Maybe store as a set of 'to' pointers
			nonFailedBranches = append(nonFailedBranches, branch{
				currState: matchingTransitions[i],
			})
		}
	}
	r.branches = nonFailedBranches

	// move along epsilon transitions after
	r.processEpsilons(input, r.branches)

	if len(r.branches) == 0 {
		return Fail
	}
	for _, b := range r.branches {
		if b.currState.isSuccessState() {
			return Success
		}
	}
	return Normal
}

func (r *runner) processEpsilons(input rune, branches []branch) {
	for bIndex := range branches {
		matchingTransitions := r.branches[bIndex].currState.matchingTransitions(input)

		for _, t := range r.branches[bIndex].currState.transitions1 {
			if t.epsilon {
				if len(matchingTransitions) == 1 {
					r.branches[bIndex].currState = t.to
				} else if len(matchingTransitions) > 1 {
					r.branches = append(r.branches, branch{
						currState: t.to,
					})
				}
			}
		}
	}
}

func (r *runner) Reset() {
	r.branches = []branch{{currState: r.head}}
}

func (r *runner) onFailState(b branch) bool {
	return b.currState == r.failState
}

type builder struct {
	states []*StateLinked
}

var GlobalIdCounter = 0

func NewStateLinkedBuilder() *builder {
	var states []*StateLinked
	states = append(states, &StateLinked{id: 0}) // stand in for fail state
	return &builder{states: states}
}

func (b *builder) AddTransition(from, to int, letter rune) *builder {
	b.fillEmptyStatesTo(to)
	b.fillEmptyStatesTo(from)
	b.states[from].transitions1 = append(b.states[from].transitions1, transitionLinked{
		description: fmt.Sprintf("Matches: '%s'", string(letter)),
		to:          b.states[to],
		predicate:   func(input rune) bool { return input == letter },
	})
	return b
}

func (b *builder) AddWildTransition(from, to int) *builder {
	b.fillEmptyStatesTo(to)
	b.fillEmptyStatesTo(from)
	b.states[from].transitions1 = append(b.states[from].transitions1, transitionLinked{
		description: fmt.Sprintf("Matches anything"),
		to:          b.states[to],
		predicate:   func(input rune) bool { return true },
	})
	return b
}

func (b *builder) AddMachineTransition(from int, state *StateLinked) *builder {
	b.fillEmptyStatesTo(from)
	for _, t := range state.transitions1 {
		// when composing a transition, we merge the first transitions of the new state into the transition of the from state
		b.states[from].transitions1 = append(b.states[from].transitions1, transitionLinked{
			description: t.description,
			to:          t.to,
			predicate:   t.predicate,
		})
	}
	return b
}

func (b *builder) fillEmptyStatesTo(from int) {
	if from >= len(b.states) {
		for i := len(b.states); i <= from; i++ {
			GlobalIdCounter++
			b.states = append(b.states, &StateLinked{id: GlobalIdCounter})
		}
	}
}

func (b *builder) Build() *StateLinked {
	return b.states[1]
}
