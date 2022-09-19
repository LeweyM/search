package v2

type Parser struct {
	tokens []token
}

func NewParser(tokens []token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() Node {
	group := Group{}

	for _, t := range p.tokens {
		switch t.symbol {
		case Character:
			group.Append(CharacterLiteral{Character: t.letter})
		}
	}

	return &group
}
