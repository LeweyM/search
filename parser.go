package v2

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(tokens []token) Node {
	group := Group{}

	for _, t := range tokens {
		switch t.symbol {
		case Character:
			group.Append(CharacterLiteral{Character: t.letter})
		}
	}

	return &group
}
