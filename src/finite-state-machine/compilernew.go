package finite_state_machine

import (
	"fmt"
	"reflect"
	"search/src/ast"
)

type compilableAst interface {
	compile() (head *State, tail *State)
}

type CompilableAstGroupNode ast.Group
type CompilableAstCharacterLiteral ast.CharacterLiteral
type CompilableAstBranchNode ast.Branch
type CompilableAstModifierExpression ast.ModifierExpression

func CompileNEW(a ast.Ast) *State {
	d := getCompilableAst(a)

	head, _ := d.compile()
	return head
}

// CharacterLiteral
func (c *CompilableAstCharacterLiteral) compile() (head *State, tail *State) {
	startingState := State{}
	endState := State{}

	transition := Transition{
		to:          &endState,
		predicate:   func(input rune) bool { return input == c.Character },
		description: fmt.Sprintf("-- %s -->", string(c.Character)),
	}

	startingState.transitions = append(startingState.transitions, transition)
	return &startingState, &endState
}

// Group
func (b *CompilableAstGroupNode) compile() (head *State, tail *State) {
	startState := State{}
	currentTail := &startState

	for _, expression := range b.Expressions {
		c := getCompilableAst(expression)

		nextStateHead, nextStateTail := c.compile()
		currentTail.addEpsilonTransition(nextStateHead)
		currentTail = nextStateTail
	}

	return &startState, currentTail
}

// Branch
func (b *CompilableAstBranchNode) compile() (head *State, tail *State) {
	startState := State{}
	endState := State{}

	for _, expression := range b.Expressions {
		c := getCompilableAst(expression)
		nextStateHead, nextStateTail := c.compile()

		startState.addEpsilonTransition(nextStateHead)
		nextStateTail.addEpsilonTransition(&endState)
	}

	return &startState, &endState
}

// Modifier
func (m *CompilableAstModifierExpression) compile() (*State, *State) {
	start := &State{}
	end := &State{}

	head, tail := getCompilableAst(m.Expression).compile()

	if m.Modifier == ast.ZeroOrManyModifier {
		start.addEpsilonTransition(tail)
		tail.addEpsilonTransition(start)
	}
	start.addEpsilonTransition(head)
	tail.addEpsilonTransition(end)

	return start, end
}

/*
getCompilableAst is a factory method for transforming compilable version of ASTs.
*/
func getCompilableAst(a ast.Ast) compilableAst {
	group, ok := a.(*ast.Group)
	if ok {
		compilableAstGroupNode := CompilableAstGroupNode(*group)
		return &compilableAstGroupNode
	}

	characterLiteral, ok := a.(ast.CharacterLiteral)
	if ok {
		compilableAstCharacterLiteral := CompilableAstCharacterLiteral(characterLiteral)
		return &compilableAstCharacterLiteral
	}

	branch, ok := a.(*ast.Branch)
	if ok {
		compilableAstBranchNode := CompilableAstBranchNode(*branch)
		return &compilableAstBranchNode
	}

	modifier, ok := a.(ast.ModifierExpression)
	if ok {
		compilableAstModifierExpression := CompilableAstModifierExpression(modifier)
		return &compilableAstModifierExpression
	}

	panic(fmt.Sprintf("expression of type [%s] cannot be compiled. Implementation missing.", reflect.TypeOf(a)))
}

func (s *State) addEpsilonTransition(destination *State) {
	s.transitions = append(s.transitions, Transition{to: destination, predicate: func(r rune) bool { return true }, description: "epsilon", epsilon: true})
}
