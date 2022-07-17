package ast

type Ast interface {
}

type Node interface {
	Expressions() []Ast
	Append(Ast)
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

func (g *Group) Expressions() []Ast {
	return nil
}

func (g *Group) Append(ast Ast) {
	g.expressions = append(g.expressions, ast)
}

type Branch struct {
	expressions []Ast
}

func (g *Branch) Append(ast Ast) {
	g.expressions = append(g.expressions, ast)
}

func (g *Branch) Expressions() []Ast {
	return nil
}

type CharacterLiteral struct {
	character rune
}

type ModifierExpression struct {
	modifier   Modifier
	expression Ast
}
