package v9

type Parser struct {
	fsmStack []CompositeNode
}

func NewParser() *Parser {
	return &Parser{fsmStack: []CompositeNode{}}
}

func (p *Parser) Parse(tokens []token) Node {
	p.pushNewGroup()

	for _, t := range tokens {
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
		case LParen:
			p.pushNewGroup()
		case RParen:
			inner := p.pop()
			outer := p.pop()
			outer.Append(inner)
			p.push(outer)
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
