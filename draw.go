package v5

import (
	"fmt"
	"strings"
)

func (s *State) Draw() (graph string, nodeSet OrderedSet[*State]) {
	// initialize sets
	transitionSet := OrderedSet[Transition]{}
	nodeSet = OrderedSet[*State]{}

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

	// draw outline around success nodes
	for _, state := range nodeSet.list() {
		if state.isSuccessState() {
			output = append(output, fmt.Sprintf("style %d stroke:green,stroke-width:4px;", nodeSet.getIndex(state)))
		}
	}

	return strings.Join(output, "\n"), nodeSet
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
}

// drawSnapshot will draw a mermaid graph from the FSM, as well as color the current node.
func (r runner) drawSnapshot() string {
	graph, nodeSet := r.head.Draw()
	switch r.GetStatus() {
	case Normal:
		graph += fmt.Sprintf("\nstyle %d fill:#ff5555;", nodeSet.getIndex(r.current))
	case Success:
		graph += fmt.Sprintf("\nstyle %d fill:#00ab41;", nodeSet.getIndex(r.current))
	}

	return graph
}
