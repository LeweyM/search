package ast

type stackLine struct {
	tail, head Node
}

type parser struct {
	stack []stackLine
}

func (p *parser) parse(input string) Ast {
	tokens := lex(input)
	rootGroup := Group{}
	p.push(stackLine{
		tail: &rootGroup,
		head: &rootGroup,
	})

	for i, token := range tokens {
		if token.symbolType == Character {
			if isInBounds(tokens, i+1) && isModifier(tokens, i+1) {
				g := p.pop()
				g.tail.Append(ModifierExpression{
					expression: CharacterLiteral{character: token.letter},
					modifier:   mapModifierTokenToAstModifier(tokens[i+1].symbolType),
				})
				p.push(g)
				i++
			} else {
				g := p.pop()
				g.tail.Append(CharacterLiteral{character: token.letter})
				p.push(g)
			}
		}

		if token.symbolType == Pipe {
			node := p.pop()

			group, isGroup := node.head.(*Group)
			_, isBranch := node.head.(*Branch)
			if isGroup {
				// turn group into branch
				newGroup := &Group{}
				node.head = &Branch{expressions: []Ast{
					group,
					newGroup,
				}}
				node.tail = newGroup
			} else if isBranch {
				newGroup := Group{}
				node.head.Append(&newGroup)
				node.tail = &newGroup
			} else {
				panic("must be either group or branch")
			}

			p.push(node)
		}
	}

	return p.pop().head
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
	token := tokens[i]
	return token.symbolType == ZeroOrMore || token.symbolType == ZeroOrOne || token.symbolType == OneOrMore
}
