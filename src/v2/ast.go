package v2

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

type CharacterLiteral struct {
	Character rune
}

func (c CharacterLiteral) string(indent int) string {
	indentation := strings.Repeat(indentationFormat, indent)

	return fmt.Sprintf("\n%sLiteral {%s}\n%s", indentation, string(c.Character), indentation)
}
