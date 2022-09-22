package v9

import (
	"fmt"
	"sort"
	"strings"
)

func (s *State) Draw(nodeSet *OrderedSet[*State]) (graph string) {
	// initialize sets
	transitionSet := OrderedSet[Transition]{}

	// collect transitions
	visitNodes(s, &transitionSet, nodeSet, &OrderedSet[*State]{})

	output := []string{
		"graph LR",
	}

	// draw transitions
	for _, t := range transitionSet.list() {
		fromId := nodeSet.getIndex(t.from)
		toId := nodeSet.getIndex(t.to)
		if t.debugSymbol == "ε" {
			output = append(output, fmt.Sprintf("%d((%d)) -.\"%s\".-> %d((%d))", fromId, fromId, t.debugSymbol, toId, toId))
		} else {
			output = append(output, fmt.Sprintf("%d((%d)) --\"%s\"--> %d((%d))", fromId, fromId, t.debugSymbol, toId, toId))
		}
	}

	for _, state := range nodeSet.list() {
		if state.isSuccessState() {
			output = append(output, fmt.Sprintf("style %d stroke:green,stroke-width:4px;", nodeSet.getIndex(state)))
		}
	}

	return strings.Join(output, "\n")
}

func visitNodes(
	node *State,
	transitions *OrderedSet[Transition],
	nodes *OrderedSet[*State],
	visited *OrderedSet[*State],
) {
	// 1. If the current node has already been visited, stop.
	if visited.has(node) {
		return
	}

	// 2.i. Add transitions for the nodes epsilon transitions
	for _, epsilon := range node.epsilons {
		transitions.add(Transition{
			debugSymbol: "ε",
			from:        node,
			to:          epsilon,
			predicate:   Predicate{},
		})
	}
	// 2.ii Add the transitions from this node to a set of transitions.
	for _, transition := range node.transitions {
		transitions.add(transition)
	}

	// 3. Mark the current node as visited.
	visited.add(node)
	nodes.add(node)

	// 4.i. Recur on every epsilon.
	for _, epsilon := range node.epsilons {
		visitNodes(epsilon, transitions, nodes, visited)
	}
	// 4. Recur on the destination node of every outgoing transition.
	for _, transition := range node.transitions {
		destinationNode := transition.to
		visitNodes(destinationNode, transitions, nodes, visited)
	}
	//// 5. Recur on the source node of every incoming transition.
	//for _, sourceNode := range node.incoming {
	//	visitNodes(sourceNode, transitions, nodes, visited)
	//}
}

// drawSnapshot will draw a mermaid graph from the FSM, as well as color the active nodes.
func (r runner) drawSnapshot(nodeSet *OrderedSet[*State]) string {
	graph := r.head.Draw(nodeSet)
	activeStates := getSortedActiveStates(r.activeStates.list(), nodeSet)

	for _, state := range activeStates {
		nodeLabel := nodeSet.getIndex(state)
		if state.isSuccessState() {
			graph += fmt.Sprintf("\nstyle %d fill:#00ab41;", nodeLabel)
		} else {
			graph += fmt.Sprintf("\nstyle %d fill:#ff5555;", nodeLabel)
		}
	}

	return graph
}

func getSortedActiveStates(activeStates []*State, nodeSet *OrderedSet[*State]) []*State {
	byAscendingNodeLabel := func(i, j int) bool {
		return nodeSet.getIndex(activeStates[i]) < nodeSet.getIndex(activeStates[j])
	}
	sort.Slice(activeStates, byAscendingNodeLabel)
	return activeStates
}
