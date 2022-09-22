package v9

import (
	"fmt"
	"sort"
	"strings"
)

type reducer struct {
	root               *State
	characters         []rune
	transitionMap      map[rune][]Transition
	compoundStates     compoundStateSet
	defaultTransitions []Transition
}

func newReducer(root *State) *reducer {
	return &reducer{root: root}
}

/*
algorithm:

gather the overlapping transitions for a set of states, such that we have all transitions for each
input character. There should also be a "." character which will represent all other characters not in the map.

i.e.
(0) -.-> (1) -b-> (3)

	-a-> (2) -c->

	{
		"a": [TransitionA(0-.->1), TransitionB(0-a->2)]
		"b": [TransitionA(0-.->1), TransitionC(1-b->3)]
		"c": [TransitionA(0-.->1), TransitionD(2-c->3)]
	}

Then, we add root state to a queue.

Iterating through the queue, with s == state

 0. get ep-closure of s. c = ep-closure
    0.a if s is a composite state, c = union of composed states closures
 1. look at each character in transition map, and collect all destinations where the origin is in c.
 2. combine the destinations into a composite state and, if not already visited, add to queue. I.e. 0(a) => State(1-2), 0(b) => nothing, 0(.) => State(1)
 3. collect a list of transitions from the steps above for s.
 4. remove original transitions from s, and replace with transitions collected in step 3.
 5. take a new state from queue

this should lead to the following given the first example

(0) -a-----> (1-2) -b-> (3)

	-[.-a]->       -c-> (3)

other format:
(0) -a-----> (1-2)
(0) -[.-a]-> (1-2) // all characters other than 'a'
(1-2) -b---> (3)
(1-2) -c---> (3)
*/
func (r *reducer) reduce() {
	visited := NewSet[*State]()
	visited.add(r.root)
	queue := []*State{}
	queue = append(queue, r.root)

	for len(queue) > 0 {
		s := queue[0]
		queue = queue[1:]

		newStates := r.reduceStateEpsilons(s)
		for _, newState := range newStates {
			if !visited.has(newState) {
				visited.add(newState)
				queue = append(queue, newState)
			}
		}
	}
}

/*
	 given a set of states, create a map of inputs to a list of transitions which all match those inputs.
	 i.e.
		(0) -.-> (1) -b-> (3)
			-a-> (2)

		{
			"a": [TransitionA(0-.->1), TransitionB(0-a->2)]
			"b": [TransitionA(0-.->1), TransitionC(1-b->3)]
			".": [TransitionA(0-.->1)]
		}
*/

// reduces state and returns next states to process
func (r *reducer) reduceState(s *State) (nextStates []*State) {
	closure := r.getClosureOfState(s)
	var nextTransitions []Transition

	// group the transitions per character and create a new transition, possibly using a compound state
	for char, transitions := range closure.transitionsFromCharacters {
		destinations := mapFn[Transition, *State](transitions, func(t Transition) *State { return t.to })
		setOfDestination := NewSet[*State](destinations...)

		if setOfDestination.size() > 0 {
			newState := r.compoundStates.getCompoundState(setOfDestination.list()...)

			nextTransitions = append(nextTransitions, Transition{
				debugSymbol: string(char),
				to:          newState,
				from:        s,
				predicate:   Predicate{allowedChars: string(char)},
			})
			nextStates = append(nextStates, newState)
		}
	}

	if closure.isFinal {
		s.isSuccess = true
	}

	// add default transitions, if there are for this clojure, but removing overlapping characters
	for _, t := range closure.wildcardTransitions {
		defaultT := t
		for _, tt := range nextTransitions {
			t.predicate.disallowedChars += tt.predicate.allowedChars
		}
		nextTransitions = append(nextTransitions, defaultT)
	}

	s.transitions = nil
	s.epsilons = nil
	s.transitions = nextTransitions

	return nextStates
}

func (r *reducer) reduceStateEpsilons(s *State) (nextStates []*State) {
	closure := r.getClosureOfState(s)
	var nextTransitions []Transition

	for _, transitions := range closure.transitionsFromCharacters {
		nextTransitions = append(
			nextTransitions,
			mapFn[Transition](transitions, func(t Transition) Transition { return t.From(s) })...,
		)
		for _, t := range transitions {
			nextStates = append(nextStates, t.to)
		}
	}

	if closure.isFinal {
		s.isSuccess = true
	}

	nextTransitions = append(
		nextTransitions,
		mapFn[Transition](closure.wildcardTransitions, func(t Transition) Transition { return t.From(s) })...,
	)

	for _, t := range closure.wildcardTransitions {
		nextStates = append(nextStates, t.to)
	}

	s.transitions = nil
	s.epsilons = nil
	s.transitions = nextTransitions

	return nextStates
}

func mapFn[T comparable, R any](items []T, mapFunc func(item T) R) []R {
	res := make([]R, len(items))

	for i := range items {
		res[i] = mapFunc(items[i])
	}

	return res
}

type compoundStateSet struct {
	compoundMap map[string]*State
	innerMap    map[*State][]*State
}

func (c *compoundStateSet) get(states ...*State) *State {
	return c.compoundMap[c.getCompoundKey(states)]
}

func (c *compoundStateSet) add(states ...*State) *State {
	if len(states) == 0 {
		return nil
	}

	if c.compoundMap == nil {
		c.compoundMap = make(map[string]*State)
	}
	if c.innerMap == nil {
		c.innerMap = make(map[*State][]*State)
	}

	var newState *State
	if len(states) == 1 {
		newState = states[0]
	} else {
		newState = &State{}
	}

	c.compoundMap[c.getCompoundKey(states)] = newState
	c.innerMap[newState] = states
	return newState
}

func (c *compoundStateSet) getCompoundState(states ...*State) *State {
	state := c.get(states...)
	if state != nil {
		return state
	}

	return c.add(states...)
}

func (c *compoundStateSet) getInnerStates(state *State) []*State {
	states, ok := c.innerMap[state]
	if !ok {
		return []*State{state}
	}
	return states
}

func (c *compoundStateSet) getCompoundKey(states []*State) string {
	hashes := []string{}
	for _, state := range states {
		hashes = append(hashes, fmt.Sprintf("%p", state))
	}
	sort.Strings(hashes)
	return strings.Join(hashes, "-")
}

type closure struct {
	isFinal                   bool
	states                    []*State
	transitionsFromCharacters map[rune][]Transition
	wildcardTransitions       []Transition
}

func (r *reducer) getClosureOfState(state *State) closure {
	var traverse func(state *State, visited Set[*State])

	traverse = func(state *State, visited Set[*State]) {
		if visited.has(state) {
			return
		}
		visited.add(state)

		for _, e := range state.epsilons {
			traverse(e, visited)
		}
	}

	states := NewSet[*State]()
	for _, s := range r.compoundStates.getInnerStates(state) {
		traverse(s, states)
	}

	// isFinal
	isFinal := false
	for s := range states {
		if s.isSuccessState() {
			isFinal = true
			break
		}
	}

	c := closure{
		states:                    states.list(),
		transitionsFromCharacters: make(map[rune][]Transition),
		wildcardTransitions:       make([]Transition, 0),
		isFinal:                   isFinal,
	}

	for s := range states {
		for _, t := range s.transitions {
			if t.predicate.disallowedChars == "" {
				// normal
				for _, r := range t.predicate.allowedChars {
					c.transitionsFromCharacters[r] = append(c.transitionsFromCharacters[r], t)
				}
			} else {
				// wildcard
				c.wildcardTransitions = append(c.wildcardTransitions, t)
			}
		}
	}

	return c
}
