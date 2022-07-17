package ast

import (
	"fmt"
	"strings"
)

const indentationFormat = "----"

type Ast interface {
	string(indent int) string
}

type Node interface {
	Append(Ast)
	string(indent int) string
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

func (g *Group) Append(ast Ast) {
	g.expressions = append(g.expressions, ast)
}

func (g *Group) String() string {
	return g.string(0)
}

func (g *Group) string(indent int) string {
	children := []string{}
	for _, child := range g.expressions {
		children = append(children, child.string(indent+1))
	}

	indentation := strings.Repeat(indentationFormat, indent)

	return fmt.Sprintf("\n%sGroup {%+v\n%s}", indentation, children, indentation)
}

type Branch struct {
	expressions []Ast
}

func (b *Branch) Append(ast Ast) {
	b.expressions = append(b.expressions, ast)
}

func (b *Branch) String() string {
	return fmt.Sprintf("Root Branch %s", b.string(0))
}

func (b *Branch) string(indent int) string {
	children := []string{}
	for _, child := range b.expressions {
		children = append(children, child.string(indent+1))
	}

	indentation := strings.Repeat(indentationFormat, indent)

	return fmt.Sprintf("\n%sBranch {%+v\n%s}", indentation, children, indentation)
}

type CharacterLiteral struct {
	character rune
}

func (c CharacterLiteral) string(indent int) string {
	indentation := strings.Repeat(indentationFormat, indent)

	return fmt.Sprintf("\n%sLiteral {%s}\n%s", indentation, string(c.character), indentation)
}

type ModifierExpression struct {
	modifier   Modifier
	expression Ast
}

//func (m ModifierExpression) String() string {
//	return fmt.Sprintf("[modifier: %v, character: %+v]", m.modifier, m.expression)
//}

func (m ModifierExpression) string(indent int) string {
	indentation := strings.Repeat(indentationFormat, indent)

	return fmt.Sprintf("\n%sModifier {%+v} expression {%+v}", indentation, m.modifier, m.expression.string(indent+1))
}
