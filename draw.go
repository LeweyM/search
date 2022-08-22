package v5

import (
	"fmt"
	"strings"
)

type StateMap struct {
	nextId   int
	stateMap map[*State]int
}

func (sm *StateMap) getId(state *State) int {
	if sm.stateMap == nil {
		sm.stateMap = make(map[*State]int)
	}

	id, inMap := sm.stateMap[state]
	if !inMap {
		i := sm.nextId
		sm.stateMap[state] = sm.nextId
		sm.nextId++
		return i
	}

	return id
}

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
		res = append(res, fmt.Sprintf("%d((%d)) -- %s --> %d((%d))", currentId, currentId, transition.debugSymbol, destinationId, destinationId))

		res = append(res, draw(sm, transition.to)...)
	}

	return res
}
