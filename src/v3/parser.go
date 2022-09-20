package v3

type Parser struct {
	tokens []token
}

func NewParser(tokens []token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() Node {
	root := Group{}

	for _, t := range p.tokens {
		switch t.symbol {
		case Character:
			root.Append(CharacterLiteral{Character: t.letter})
		}
	}

	return &root
}
