package v5

import (
	"fmt"
	"strings"
)

func (s *State) Draw() string {
	res := []string{
		"graph LR",
	}

	stateMap := StateMap{}
	for _, transition := range s.transitions {
		res = append(res, drawVertex(TransitionSet{}, stateMap, transition)...)
	}

	return strings.Join(res, "\n")
}

func drawVertex(visited TransitionSet, sm StateMap, t Transition) []string {
	res := []string{}

	if visited.has(t) {
		return res
	}

	toVisit := []Transition{}
	for _, transition := range t.to.transitions {
		toVisit = append(toVisit, transition)
	}
	for _, incoming := range t.to.incoming {
		for _, transition := range incoming.transitions {
			toVisit = append(toVisit, transition)
		}
	}

	// add the transition to the list
	fromId := sm.getId(t.from)
	toId := sm.getId(t.to)
	res = append(res, fmt.Sprintf("%d((%d)) --\"%s\"--> %d((%d))", fromId, fromId, t.debugSymbol, toId, toId))

	visited.set(t)

	// recursively add all the transitions from the child nodes and flatten the list
	for _, transition := range toVisit {
		res = append(res, drawVertex(visited, sm, transition)...)
	}

	return res
}

// StateMap is used to cache ids for States
type StateMap struct {
	nextId   int
	stateMap map[*State]int
}

func (sm *StateMap) has(state *State) bool {
	if sm.stateMap == nil {
		return false
	}

	_, has := sm.stateMap[state]
	return has
}

// getId will return an incrementing numerical id for each state
func (sm *StateMap) getId(state *State) int {
	// initialize map
	if sm.stateMap == nil {
		sm.stateMap = make(map[*State]int)
	}

	id, inMap := sm.stateMap[state]
	if !inMap {
		// if state is not in the cache, store the state and increment the id
		i := sm.nextId
		sm.stateMap[state] = sm.nextId
		sm.nextId++
		return i
	}

	// else, just return the cached id
	return id
}

type comparableTransition struct {
	symbol string
	to     destination
	from   destination
}

type TransitionSet map[comparableTransition]bool

func (ts *TransitionSet) set(transition Transition) {
	(*ts)[comparableTransition{
		symbol: transition.debugSymbol,
		to:     transition.to,
		from:   transition.from,
	}] = true
}

func (ts *TransitionSet) has(transition Transition) bool {
	return (*ts)[comparableTransition{
		symbol: transition.debugSymbol,
		from:   transition.from,
		to:     transition.to,
	}]
}
