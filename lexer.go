package v4

type symbol int

const (
	AnyCharacter symbol = iota
	Pipe
	LParen
	RParen
	Character
	ZeroOrMore
	OneOrMore
	ZeroOrOne
)

type token struct {
	symbol symbol
	letter rune
}

func lex(input string) []token {
	var tokens []token
	for _, character := range input {
		tokens = append(tokens, lexRune(character))
	}
	return tokens
}

func lexRune(r rune) token {
	var s token
	switch r {
	case '(':
		s.symbol = LParen
	case ')':
		s.symbol = RParen
	case '.':
		s.symbol = AnyCharacter
	case '|':
		s.symbol = Pipe
	case '*':
		s.symbol = ZeroOrMore
	case '+':
		s.symbol = OneOrMore
	case '?':
		s.symbol = ZeroOrOne
	default:
		s.symbol = Character
		s.letter = r
	}
	return s
}
