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
	ZeroOrOneModifier Modifier = iota
	ZeroOrManyModifier
	OneOrManyModifier
)

type Group struct {
	Expressions []Ast
}

func (g *Group) Append(ast Ast) {
	g.Expressions = append(g.Expressions, ast)
}

func (g *Group) String() string {
	return g.string(0)
}

func (g *Group) string(indent int) string {
	children := []string{}
	for _, child := range g.Expressions {
		children = append(children, child.string(indent+1))
	}

	indentation := strings.Repeat(indentationFormat, indent)

	return fmt.Sprintf("\n%sGroup {%+v\n%s}", indentation, children, indentation)
}

type Branch struct {
	Expressions []Ast
}

func (b *Branch) Append(ast Ast) {
	b.Expressions = append(b.Expressions, ast)
}

func (b *Branch) String() string {
	return fmt.Sprintf("Root Branch %s", b.string(0))
}

func (b *Branch) string(indent int) string {
	children := []string{}
	for _, child := range b.Expressions {
		children = append(children, child.string(indent+1))
	}

	indentation := strings.Repeat(indentationFormat, indent)

	return fmt.Sprintf("\n%sBranch {%+v\n%s}", indentation, children, indentation)
}

type CharacterLiteral struct {
	Character rune
}

func (c CharacterLiteral) string(indent int) string {
	indentation := strings.Repeat(indentationFormat, indent)

	return fmt.Sprintf("\n%sLiteral {%s}\n%s", indentation, string(c.Character), indentation)
}

type WildcardCharacterLiteral struct{}

func (c WildcardCharacterLiteral) string(indent int) string {
	indentation := strings.Repeat(indentationFormat, indent)

	return fmt.Sprintf("\n%sWildcard Literal\n%s", indentation, indentation)
}

type ModifierExpression struct {
	Modifier   Modifier
	Expression Ast
}

func (m ModifierExpression) string(indent int) string {
	indentation := strings.Repeat(indentationFormat, indent)

	return fmt.Sprintf("\n%sModifier {%+v} Expression {%+v}", indentation, m.Modifier, m.Expression.string(indent+1))
}
