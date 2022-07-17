package ast

type parser struct {
}

func (p *parser) parse(input string) Ast {
	tokens := lex(input)

	g := Group{}

	for i, token := range tokens {
		if token.symbolType == Character {
			if isInBounds(tokens, i+1) && isModifier(tokens, i+1) {
				g.expressions = append(g.expressions, ModifierExpression{
					expression: CharacterLiteral{character: token.letter},
					modifier:   mapModifierTokenToAstModifier(tokens[i+1].symbolType),
				})
				i++
			} else {
				g.expressions = append(g.expressions, CharacterLiteral{character: token.letter})
			}
		}
	}

	return g
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
