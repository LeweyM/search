package finite_state_machine

func Compile(input string) *StateLinked {
	symbols := lex(input)
	n := len(symbols) + 1
	var branches []*StateLinked
	builder := NewStateLinkedBuilder(n)
	prevStateNumber := 1
	for _, symbol := range symbols {
		switch symbol.symbolType {
		case LParen, RParen:
			break // skip for now
		case Pipe:
			builder = builder.SetSuccess(prevStateNumber)
			branches = append(branches, builder.Build())
			prevStateNumber = 1
		case Wild:
			if symbol.modifier == ZeroOrMore {
				builder = builder.AddWildTransition(prevStateNumber, prevStateNumber)
			} else {
				builder = builder.AddWildTransition(prevStateNumber, prevStateNumber+1)
				prevStateNumber++
			}
		default:
			builder = builder.AddTransition(prevStateNumber, prevStateNumber+1, symbol.letter)
			prevStateNumber++
		}
	}
	builder = builder.SetSuccess(prevStateNumber)

	for _, b := range branches {
		builder.AddMachineTransition(1, b)
	}

	stateLinked := builder.Build()
	return stateLinked
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
