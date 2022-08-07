package v2

type stackLine struct {
	tail, head Node
}

type Parser struct {
	stack []stackLine
}

func NewParser() *Parser {
	return &Parser{stack: []stackLine{}}
}

func (p *Parser) Parse(tokens []token) Ast {
	p.pushNewGroup()

	for _, t := range tokens {
		switch t.symbolType {
		case Character:
			node := p.pop()
			node.tail.Append(CharacterLiteral{Character: t.letter})
			p.push(node)
		}
	}

	return p.pop().head
}

func (p *Parser) pushNewGroup() {
	newGroup := Group{}
	p.push(stackLine{
		tail: &newGroup,
		head: &newGroup,
	})
}
func (p *Parser) pop() stackLine {
	pop := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]

	return pop
}

func (p *Parser) push(g stackLine) {
	p.stack = append(p.stack, g)
}
