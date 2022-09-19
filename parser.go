package v7

type Parser struct {
	fsmStack []CompositeNode
	tokens   []token
}

func NewParser(tokens []token) *Parser {
	return &Parser{fsmStack: []CompositeNode{}, tokens: tokens}
}

func (p *Parser) Parse() Node {
	p.pushNewGroup()

	for _, t := range p.tokens {
		switch t.symbol {
		case Character:
			node := p.pop()
			node.Append(CharacterLiteral{Character: t.letter})
			p.push(node)
		case AnyCharacter:
			node := p.pop()
			node.Append(WildcardLiteral{})
			p.push(node)
		case Pipe:
			node := p.pop()
			switch b := node.(type) {
			case *Branch:
				b.Split()
			default:
				node = &Branch{ChildNodes: []Node{node, &Group{}}}
			}
			p.push(node)
		}
	}

	return p.pop()
}

func (p *Parser) pushNewGroup() {
	p.push(&Group{})
}

func (p *Parser) pop() CompositeNode {
	pop := p.fsmStack[len(p.fsmStack)-1]
	p.fsmStack = p.fsmStack[:len(p.fsmStack)-1]

	return pop
}

func (p *Parser) push(g CompositeNode) {
	p.fsmStack = append(p.fsmStack, g)
}
