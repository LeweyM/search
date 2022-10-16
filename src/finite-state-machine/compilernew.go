package finite_state_machine

import (
	"fmt"
	"github.com/leweyM/search/src/ast"
	"reflect"
)

type compilableAst interface {
	compile() (head *State, tail *State)
}

type CompilableAstGroupNode ast.Group
type CompilableAstCharacterLiteral ast.CharacterLiteral
type CompilableAstWildcardCharacterLiteral ast.WildcardCharacterLiteral
type CompilableAstBranchNode ast.Branch
type CompilableAstModifierExpression ast.ModifierExpression

// used to link together legacy code
func Compile(input string) *State {
	p := ast.Parser{}
	tree := p.Parse(input)
	return CompileNEW(tree)
}

func CompileNEW(a ast.Ast) *State {
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

// WildcardCharacterLiteral
func (c CompilableAstWildcardCharacterLiteral) compile() (head *State, tail *State) {
	startingState := State{}
	endState := State{}

	startingState.addTransition(&endState, func(input rune) bool { return true })
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

	switch m.Modifier {
	case ast.ZeroOrManyModifier:
		start.addEpsilonTransition(tail)
		tail.addEpsilonTransition(start)
	case ast.OneOrManyModifier:
		tail.addEpsilonTransition(start)
	case ast.ZeroOrOneModifier:
		start.addEpsilonTransition(tail)
	}
	start.addEpsilonTransition(head)
	tail.addEpsilonTransition(end)

	return start, end
}

/*
getCompilableAst is a factory method for transforming compilable version of ASTs.
*/
func getCompilableAst(a ast.Ast) compilableAst {
	switch node := a.(type) {
	case *ast.Group:
		compilableAstGroupNode := CompilableAstGroupNode(*node)
		return &compilableAstGroupNode
	case *ast.Branch:
		compilableAstBranchNode := CompilableAstBranchNode(*node)
		return &compilableAstBranchNode
	case ast.CharacterLiteral:
		compilableAstCharacterLiteral := CompilableAstCharacterLiteral(node)
		return &compilableAstCharacterLiteral
	case ast.WildcardCharacterLiteral:
		compilableAstWildcardCharacterLiteral := CompilableAstWildcardCharacterLiteral(node)
		return &compilableAstWildcardCharacterLiteral
	case ast.ModifierExpression:
		compilableAstModifierExpression := CompilableAstModifierExpression(node)
		return &compilableAstModifierExpression
	default:
		panic(fmt.Sprintf("expression of type [%s] cannot be compiled. Implementation missing.", reflect.TypeOf(a)))
	}
}
