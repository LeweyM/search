package v2

import (
	"fmt"
	"reflect"
)

type compilableAst interface {
	compile() (head *State, tail *State)
}

type CompilableAstGroupNode Group
type CompilableAstCharacterLiteral CharacterLiteral

func Compile(a Ast) *State {
	d := getCompilableAst(a)

	head, _ := d.compile()
	return head
}

// CharacterLiteral
func (c *CompilableAstCharacterLiteral) compile() (head *State, tail *State) {
	startingState := State{}
	endState := State{}

	startingState.addTransition(&endState, func(input rune) bool { return input == c.Character })
	return &startingState, &endState
}

// Group
func (b *CompilableAstGroupNode) compile() (head *State, tail *State) {
	startState := State{}
	currentTail := &startState

	for _, expression := range b.Expressions {
		c := getCompilableAst(expression)

		nextStateHead, nextStateTail := c.compile()
		currentTail.merge(nextStateHead)
		currentTail = nextStateTail
	}

	return &startState, currentTail
}

/*
getCompilableAst is a factory method for transforming compilable version of ASTs.
*/
func getCompilableAst(a Ast) compilableAst {
	switch node := a.(type) {
	case *Group:
		compilableAstGroupNode := CompilableAstGroupNode(*node)
		return &compilableAstGroupNode
	case CharacterLiteral:
		compilableAstCharacterLiteral := CompilableAstCharacterLiteral(node)
		return &compilableAstCharacterLiteral
	default:
		panic(fmt.Sprintf("expression of type [%s] cannot be compiled. Implementation missing.", reflect.TypeOf(a)))
	}
}
