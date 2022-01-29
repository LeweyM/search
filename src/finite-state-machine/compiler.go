package finite_state_machine

func Compile(input string) *StateLinked {
	symbols := lex(input)
	fragment, _ := compileFragment(symbols)
	return fragment
}

func compileFragment(symbols []symbol) (*StateLinked, int) {
	var branches []*StateLinked
	builder := NewStateLinkedBuilder(len(symbols) + 1)
	symbolCounter := 0
	prevSymbolIndex := 1
	Loop: for len(symbols) > 0 {
		symbol := symbols[0]
		switch symbol.symbolType {
		case LParen:
			innerFragment, matchingParenIndex := compileFragment(symbols[1:])
			builder.AddMachineTransition(prevSymbolIndex, innerFragment)
			prevSymbolIndex = matchingParenIndex+1
			symbols = symbols[prevSymbolIndex:]
		case RParen:
			break Loop
		case Pipe:
			builder = builder.SetSuccess(prevSymbolIndex)
			branches = append(branches, builder.Build())
			builder = NewStateLinkedBuilder(len(symbols) + 1)
			prevSymbolIndex = 1
		case Wild:
			if symbol.modifier == ZeroOrMore {
				builder = builder.AddWildTransition(prevSymbolIndex, prevSymbolIndex)
			} else {
				builder = builder.AddWildTransition(prevSymbolIndex, prevSymbolIndex+1)
				prevSymbolIndex++
			}
		default:
			builder = builder.AddTransition(prevSymbolIndex, prevSymbolIndex+1, symbol.letter)
			prevSymbolIndex++
		}
		symbolCounter++
		symbols = symbols[1:]
	}
	builder = builder.SetSuccess(prevSymbolIndex)

	for _, b := range branches {
		builder.AddMachineTransition(1, b)
	}

	stateLinked := builder.Build()
	return stateLinked, symbolCounter
}

type modifier string
const (
	ZeroOrMore modifier = "ZeroOrMore"
)

type SymbolType int
const (
	Other SymbolType = iota
	Wild
	Pipe
	LParen
	RParen
)

type symbol struct {
	symbolType SymbolType
	letter     rune
	modifier   modifier
}

func lex(input string) []symbol {
	var symbols []symbol
	i := 0
	for i < len(input) {
		curr := rune(input[i])
		hasNext := i+1 < len(input)
		if hasNext {
			next := rune(input[i+1])
			if isModifier(next) {
				symbols = append(symbols, lexRuneWithModifier(curr, next))
				i += 2
			} else {
				symbols = append(symbols, lexRune(curr))
				i += 1
			}
		} else {
			symbols = append(symbols, lexRune(curr))
			i++
		}
	}
	return symbols
}

func lexRune(r rune) symbol {
	var s symbol
	switch r {
	case '(':
		s.symbolType = LParen
	case ')':
		s.symbolType = RParen
	case '.':
		s.symbolType = Wild
	case '|':
		s.symbolType = Pipe
	default:
		s.letter = r
	}
	return s
}

func lexRuneWithModifier(r rune, mod rune) symbol {
	s := lexRune(r)
	if mod == '*' {
		s.modifier = ZeroOrMore
	}
	return s
}

func isModifier(r rune) bool {
	return r == '*'
}
