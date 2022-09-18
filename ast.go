package v9

import (
	"fmt"
	"strings"
)

type Node interface {
	compile() (head *State, tail *State)
	string(indentation int) string
}

type CompositeNode interface {
	Node
	Append(node Node)
}

type Group struct {
	ChildNodes []Node
}

func (g *Group) Append(node Node) {
	g.ChildNodes = append(g.ChildNodes, node)
}

type Branch struct {
	ChildNodes []Node
}

func (b *Branch) Append(node Node) {
	for i := len(b.ChildNodes) - 1; i > 0; i-- {
		switch n := b.ChildNodes[i].(type) {
		case CompositeNode:
			n.Append(node)
			return
		}
	}

	panic("should have at least one composite node child")
}

func (b *Branch) Split() {
	b.ChildNodes = append(b.ChildNodes, &Group{})
}

type CharacterLiteral struct {
	Character rune
}

type WildcardLiteral struct{}

/* Compiler methods */

func (b *Branch) compile() (head *State, tail *State) {
	startState := &State{}
	endState := &State{}
	for _, expression := range b.ChildNodes {
		nextStateHead, tail := expression.compile()
		startState.addEpsilon(nextStateHead)
		tail.addEpsilon(endState)
	}
	return startState, endState
}

func (g *Group) compile() (head *State, tail *State) {
	startState := State{}
	currentTail := &startState

	for _, expression := range g.ChildNodes {
		nextStateHead, nextStateTail := expression.compile()
		_, isChar := expression.(CharacterLiteral)
		if isChar {
			currentTail.merge(nextStateHead)
		} else {
			currentTail.addEpsilon(nextStateHead)
		}
		currentTail = nextStateTail
	}

	return &startState, currentTail
}

func (l CharacterLiteral) compile() (head *State, tail *State) {
	startingState := State{}
	endState := State{}

	startingState.addTransition(&endState, Predicate{allowedChars: string(l.Character)}, string(l.Character))
	return &startingState, &endState
}

func (w WildcardLiteral) compile() (head *State, tail *State) {
	startingState := State{}
	endState := State{}

	startingState.addTransition(&endState, Predicate{disallowedChars: "\n"}, ".")
	return &startingState, &endState
}

/* Debug methods */

func (g *Group) String() string {
	return "\n" + g.string(0)
}

func (b *Branch) String() string {
	return "\n" + b.string(0)
}

func (g *Group) string(indentation int) string {
	return compositeToString("Group", g.ChildNodes, indentation)
}

func (b *Branch) string(indentation int) string {
	return compositeToString("Branch", b.ChildNodes, indentation)
}

func (l CharacterLiteral) string(indentation int) string {
	padding := strings.Repeat("--", indentation)
	return fmt.Sprintf("%sCharacterLiteral('%s')", padding, string(l.Character))
}

func (w WildcardLiteral) string(indentation int) string {
	padding := strings.Repeat("--", indentation)
	return fmt.Sprintf("%sWildcardCharacterLiteral", padding)
}

func compositeToString(title string, children []Node, indentation int) string {
	padding := strings.Repeat("--", indentation)
	res := padding + title
	for _, node := range children {
		res += fmt.Sprintf("\n%s%s", padding, node.string(indentation+1))
	}
	return res
}
