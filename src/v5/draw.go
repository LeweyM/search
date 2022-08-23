package v5

import (
	"fmt"
	"strings"
)

func (s *State) Draw() string {
	sm := StateIDMap{}

	// collect transitions
	transitionSet := TransitionSet{}
	visitNodes(s, &transitionSet, make(map[*State]bool))

	output := []string{
		"graph LR",
	}

	// draw transitions
	for _, t := range transitionSet.list() {
		fromId := sm.getId(t.from)
		toId := sm.getId(t.to)
		output = append(output, fmt.Sprintf("%d((%d)) --\"%s\"--> %d((%d))", fromId, fromId, t.debugSymbol, toId, toId))
	}
	return strings.Join(output, "\n")
}

func visitNodes(node *State, transitions *TransitionSet, visited map[*State]bool) {
	// if already visited the node, return
	if visited[node] {
		return
	}

	// add transitions from node
	for _, transition := range node.transitions {
		transitions.set(transition)
	}

	visited[node] = true

	// visit all neighbours
	for _, transition := range node.transitions {
		visitNodes(transition.to, transitions, visited)
	}
	for _, incomingNode := range node.incoming {
		visitNodes(incomingNode, transitions, visited)
	}
}

// StateIDMap is used to cache ids for States
type StateIDMap struct {
	nextId   int
	stateMap map[*State]int
}

func (sm *StateIDMap) has(state *State) bool {
	if sm.stateMap == nil {
		return false
	}

	_, has := sm.stateMap[state]
	return has
}

// getId will return an incrementing numerical id for each state
func (sm *StateIDMap) getId(state *State) int {
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

// TransitionSet maintains an ordered set of unique Transitions
type TransitionSet struct {
	transitionSet  map[Transition]bool
	transitionList []Transition
}

func (ts *TransitionSet) set(t Transition) {
	hasTransition := ts.transitionSet[t]

	if !hasTransition {
		if ts.transitionSet == nil {
			ts.transitionSet = make(map[Transition]bool)
		}
		ts.transitionSet[t] = true
		ts.transitionList = append(ts.transitionList, t)
	}
}

func (ts *TransitionSet) has(t Transition) bool {
	return ts.transitionSet[t]
}

func (ts *TransitionSet) list() []Transition {
	return ts.transitionList[:]
}
