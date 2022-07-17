package ast

type Ast interface {
}

type Modifier int

const (
	zeroOrOne Modifier = iota
	zeroOrMany
	OneOrMany
)

type Group struct {
	expressions []Ast
}

type Branch struct {
	expressions []Ast
}

type CharacterLiteral struct {
	character rune
}

type ModifierExpression struct {
	modifier   Modifier
	expression Ast
}
