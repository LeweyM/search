package v5

type Node interface {
	compile() (head *State, tail *State)
}

type CompositeNode interface {
	Node
	Append(node Node)
}

type Group struct {
	ChildNodes []Node
}

type CharacterLiteral struct {
	Character rune
}

func (g *Group) Append(node Node) {
	g.ChildNodes = append(g.ChildNodes, node)
}

type WildcardLiteral struct{}

/* Compiler methods */

func (g *Group) compile() (head *State, tail *State) {
	startState := State{}
	currentTail := &startState

	for _, expression := range g.ChildNodes {
		nextStateHead, nextStateTail := expression.compile()
		currentTail.merge(nextStateHead)
		currentTail = nextStateTail
	}

	return &startState, currentTail
}

func (l CharacterLiteral) compile() (head *State, tail *State) {
	startingState := State{}
	endState := State{}

	startingState.addTransition(&endState, func(input rune) bool { return input == l.Character }, string(l.Character))
	return &startingState, &endState
}

func (w WildcardLiteral) compile() (head *State, tail *State) {
	startingState := State{}
	endState := State{}

	startingState.addTransition(&endState, func(input rune) bool { return input != '\n' }, ".")
	return &startingState, &endState
}
