package v6

type Parser struct {
	tokens []token
}

func NewParser(tokens []token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() Node {
	var root CompositeNode
	root = &Group{}

	for _, t := range p.tokens {
		switch t.symbol {
		case Character:
			root.Append(CharacterLiteral{Character: t.letter})
		case AnyCharacter:
			root.Append(WildcardLiteral{})
		case Pipe:
			switch b := root.(type) {
			case *Branch:
				b.Split()
			default:
				root = &Branch{ChildNodes: []Node{root, &Group{}}}
			}
		}
	}

	return root
}
