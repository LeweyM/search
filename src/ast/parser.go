package ast

type stackLine struct {
	tail, head Node
}

type Parser struct {
	stack []stackLine
}

func (p *Parser) Parse(input string) Ast {
	tokens := lex(input)
	p.pushNewGroup()

	for i, token := range tokens {
		switch token.symbolType {
		case AnyCharacter:
			node := p.pop()
			node.tail.Append(p.wrapWithModifier(tokens, i, WildcardCharacterLiteral{}))
			p.push(node)
		case Character:
			node := p.pop()
			node.tail.Append(p.wrapWithModifier(tokens, i, CharacterLiteral{Character: token.letter}))
			p.push(node)
		case Pipe:
			node := p.pop()
			newGroup := Group{}
			node.head = &Branch{Expressions: []Ast{
				node.head,
				&newGroup,
			}}
			node.tail = &newGroup
			p.push(node)
		case LParen:
			p.pushNewGroup()
		case RParen:
			inner := p.pop()
			outer := p.pop()
			outer.tail.Append(p.wrapWithModifier(tokens, i, inner.head))
			p.push(outer)
		}
	}

	return p.pop().head
}

func (p *Parser) wrapWithBranch(node Ast, newGroup *Group) *Branch {
	return &Branch{Expressions: []Ast{
		node,
		newGroup,
	}}
}

func (p *Parser) pushNewGroup() {
	newGroup := Group{}
	p.push(stackLine{
		tail: &newGroup,
		head: &newGroup,
	})
}

func (p *Parser) wrapWithModifier(tokens []symbol, i int, expression Ast) Ast {
	if isModifier(tokens, i+1) {
		expression = ModifierExpression{
			Expression: expression,
			Modifier:   mapModifierTokenToAstModifier(tokens[i+1].symbolType),
		}
	}
	return expression
}

func (p *Parser) pop() stackLine {
	pop := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]

	return pop
}

func (p *Parser) push(g stackLine) {
	p.stack = append(p.stack, g)
}

func mapModifierTokenToAstModifier(symbolType SymbolType) Modifier {
	m := map[SymbolType]Modifier{
		ZeroOrOne:  ZeroOrOneModifier,
		ZeroOrMore: ZeroOrManyModifier,
		OneOrMore:  OneOrManyModifier,
	}

	return m[symbolType]
}

func isInBounds(tokens []symbol, i int) bool {
	return i < len(tokens)
}

func isModifier(tokens []symbol, i int) bool {
	if !isInBounds(tokens, i) {
		return false
	}

	token := tokens[i]
	return token.symbolType == ZeroOrMore || token.symbolType == ZeroOrOne || token.symbolType == OneOrMore
}
