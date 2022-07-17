package ast

type stackLine struct {
	tail, head Node
}

type parser struct {
	stack []stackLine
}

func (p *parser) parse(input string) Ast {
	tokens := lex(input)
	p.startNewGroupOnStack()

	for i, token := range tokens {
		switch token.symbolType {
		case Character:
			p.concatToTail(tokens, i, CharacterLiteral{character: token.letter})
		case Pipe:
			node := p.pop()
			newGroup := Group{}

			group, isGroup := node.head.(*Group)
			_, isBranch := node.head.(*Branch)

			if isGroup {
				// wrap group in branch node
				node.head = &Branch{expressions: []Ast{
					group,
					&newGroup,
				}}
			} else if isBranch {
				// append new branch
				node.head.Append(&newGroup)
			} else {
				panic("must be either group or branch")
			}
			node.tail = &newGroup
			p.push(node)
		case LParen:
			p.startNewGroupOnStack()
		case RParen:
			inner := p.pop()
			p.concatToTail(tokens, i, inner.head)
		}
	}

	return p.pop().head
}

func (p *parser) startNewGroupOnStack() {
	newGroup := Group{}
	p.push(stackLine{
		tail: &newGroup,
		head: &newGroup,
	})
}

func (p *parser) concatToTail(tokens []symbol, i int, expression Ast) {
	if isModifier(tokens, i+1) {
		expression = ModifierExpression{
			expression: expression,
			modifier:   mapModifierTokenToAstModifier(tokens[i+1].symbolType),
		}
	}

	g := p.pop()
	g.tail.Append(expression)
	p.push(g)
}

func (p *parser) pop() stackLine {
	pop := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]

	return pop
}

func (p *parser) push(g stackLine) {
	p.stack = append(p.stack, g)
}

func mapModifierTokenToAstModifier(symbolType SymbolType) Modifier {
	m := map[SymbolType]Modifier{
		ZeroOrOne:  zeroOrOne,
		ZeroOrMore: zeroOrMany,
		OneOrMore:  OneOrMany,
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
