package ast

type Ast interface {
}

type Group struct {
	expressions []Ast
}

type CharacterLiteral struct {
	character rune
}
