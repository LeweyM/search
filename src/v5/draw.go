package v5

import (
	"fmt"
	"strings"
)

func (s *State) Draw() string {
	// initialize sets
	transitionSet := OrderedSet[Transition]{}
	nodeSet := OrderedSet[*State]{}

	// collect transitions
	visitNodes(s, &transitionSet, &nodeSet)

	output := []string{
		"graph LR",
	}

	// draw transitions
	for _, t := range transitionSet.list() {
		fromId := nodeSet.getIndex(t.from)
		toId := nodeSet.getIndex(t.to)
		output = append(output, fmt.Sprintf("%d((%d)) --\"%s\"--> %d((%d))", fromId, fromId, t.debugSymbol, toId, toId))
	}
	return strings.Join(output, "\n")
}

func visitNodes(
	node *State,
	transitions *OrderedSet[Transition],
	visited *OrderedSet[*State],
) {
	// 1. If the current node has already been visited, stop.
	if visited.has(node) {
		return
	}

	// 2. Add the transitions from this node to a set of transitions.
	for _, transition := range node.transitions {
		transitions.add(transition)
	}

	// 3. Mark the current node as visited.
	visited.add(node)

	// 4. Recur on the destination node of every outgoing transition.
	for _, transition := range node.transitions {
		destinationNode := transition.to
		visitNodes(destinationNode, transitions, visited)
	}
	// 5. Recur on the source node of every incoming transition.
	for _, sourceNode := range node.incoming {
		visitNodes(sourceNode, transitions, visited)
	}
}

// OrderedSet maintains an ordered set of unique items of type <T>
type OrderedSet[T comparable] struct {
	set       map[T]int
	nextIndex int
}

func (o *OrderedSet[T]) add(t T) {
	if o.set == nil {
		o.set = make(map[T]int)
	}

	if !o.has(t) {
		o.set[t] = o.nextIndex
		o.nextIndex++
	}
}

func (o *OrderedSet[T]) has(t T) bool {
	_, hasItem := o.set[t]
	return hasItem
}

func (o *OrderedSet[T]) list() []T {
	size := len(o.set)
	list := make([]T, size)

	for t, i := range o.set {
		list[i] = t
	}

	return list
}

func (o *OrderedSet[T]) getIndex(t T) int {
	return o.set[t]
}
