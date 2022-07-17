package ast

type parser struct {
}

func (p *parser) parse(input string) Ast {
	g := Group{}

	for _, char := range input {
		g.expressions = append(g.expressions, CharacterLiteral{character: char})
	}

	return g
}
