package ast

type SymbolType int

const (
	AnyCharacter SymbolType = iota
	Pipe
	LParen
	RParen
	Character
	ZeroOrMore
	OneOrMore
	ZeroOrOne
)

type symbol struct {
	symbolType SymbolType
	letter     rune
}

func lex(input string) []symbol {
	var symbols []symbol
	i := 0
	for i < len(input) {
		symbols = append(symbols, lexRune(rune(input[i])))
		i++
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
		s.symbolType = AnyCharacter
	case '|':
		s.symbolType = Pipe
	case '*':
		s.symbolType = ZeroOrMore
	case '+':
		s.symbolType = OneOrMore
	case '?':
		s.symbolType = ZeroOrOne
	default:
		s.symbolType = Character
		s.letter = r
	}
	return s
}
