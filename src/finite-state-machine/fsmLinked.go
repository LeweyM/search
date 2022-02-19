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

func NewEpsilon(to *StateLinked) transitionLinked {
	return transitionLinked{
		to:          to,
		predicate:   func(input rune) bool { return true },
		description: "epsilon",
		epsilon:     true,
	}
}

type branchSet struct{ set map[*StateLinked]bool }

func newBranchSet() *branchSet {
	return &branchSet{set: make(map[*StateLinked]bool)}
}

func (b *branchSet) add(state *StateLinked) {
	b.set[state] = true
}

func (b *branchSet) contains(state *StateLinked) bool {
	return b.set[state]
}

func (b *branchSet) remove(state *StateLinked) {
	delete(b.set, state)
}

func (b *branchSet) getSet() map[*StateLinked]bool {
	return b.set
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
	return matchingTransitions
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
		s.transitions1 = append(s.transitions1, t)
	}
}

type runner struct {
	head      *StateLinked
	failState *StateLinked
	branches  *branchSet
}

func NewRunner(head *StateLinked) *runner {
	failState := &StateLinked{id: 0}

	return &runner{
		failState: failState,
		head:      head,
		branches:  newBranchSet(),
	}
}

func (r *runner) Next(input rune) StateType {
	// move along epsilon transitions first.
	// This is probably inefficient and could be moved into the main loop.
	r.processEpsilons()

	// move along regular transitions
	var nonFailedBranches = newBranchSet()
	for br := range r.branches.set {
		for _, destinationState := range br.matchingTransitions(input) {
			nonFailedBranches.add(destinationState)
		}
	}

	r.branches = nonFailedBranches

	// move along epsilon transitions after
	r.processEpsilons()

	if len(r.branches.set) == 0 {
		return Fail
	}
	for b := range r.branches.set {
		if b.isSuccessState() {
			return Success
		}
	}
	return Normal
}

func (r *runner) processEpsilons() {
	nextBranches := newBranchSet()
	for br := range r.branches.set {
		// if a branch contains an epsilon transition
		for _, t := range br.transitions1 {
			if t.epsilon {
				// add it to a branch
				nextBranches.add(t.to)
			}
		}
		nextBranches.add(br)
	}

	r.branches = nextBranches
}

func (r *runner) Reset() {
	r.branches = newBranchSet()
	r.branches.add(r.head)
}

func (r *runner) onFailState(b *StateLinked) bool {
	return b == r.failState
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
