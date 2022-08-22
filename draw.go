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
	res = append(res, draw(stateMap, s)...)

	return strings.Join(res, "\n")
}

func draw(sm StateMap, s *State) []string {
	res := []string{}

	currentId := sm.getId(s)
	for _, transition := range s.transitions {
		destinationId := sm.getId(transition.to)

		// add the transition to the list
		res = append(res, fmt.Sprintf("%d((%d)) -- %s --> %d((%d))", currentId, currentId, transition.debugSymbol, destinationId, destinationId))

		// recursively add all the transitions from the child nodes and flatten the list
		res = append(res, draw(sm, transition.to)...)
	}

	return res
}

// StateMap is used to cache ids for States
type StateMap struct {
	nextId   int
	stateMap map[*State]int
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
