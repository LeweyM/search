package v9

type Parser struct {
	fsmStack []CompositeNode
	tokens   []token
}

func NewParser(tokens []token) *Parser {
	return &Parser{
		tokens:   tokens,
		fsmStack: []CompositeNode{},
	}
}

func (p *Parser) Parse() Node {
	p.pushNewGroup()

	for i, t := range p.tokens {
		switch t.symbol {
		case Character:
			node := p.pop()
			node.Append(p.wrapWithModifier(i, CharacterLiteral{Character: t.letter}))
			p.push(node)
		case AnyCharacter:
			node := p.pop()
			node.Append(p.wrapWithModifier(i, WildcardLiteral{}))
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
			outer.Append(p.wrapWithModifier(i, inner))
			p.push(outer)
		}
	}

	return p.pop()
}

func (p *Parser) peekAhead(i int) (bool, token) {
	nextIndex := i + 1

	if nextIndex >= len(p.tokens) {
		return false, token{}
	}

	return true, p.tokens[nextIndex]
}

func (p *Parser) wrapWithModifier(i int, child Node) Node {
	ok, nextToken := p.peekAhead(i)
	if ok {
		switch nextToken.symbol {
		case ZeroOrMore:
			return ZeroOrMoreModifier{Child: child}
		}
	}

	return child
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
