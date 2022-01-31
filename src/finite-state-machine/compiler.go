package finite_state_machine

func Compile(input string) *StateLinked {
	symbols := lex(input)
	compiler := stackCompiler{}
	stateLinked := compiler.compile(symbols)
	return stateLinked
}

type stackCompiler struct {
	stack []*StateLinked
}

func (s *stackCompiler) pop() *StateLinked {
	i := len(s.stack) - 1
	state := s.stack[i]
	s.stack = s.stack[:i]
	return state
}

func (s *stackCompiler) push(linked *StateLinked) {
	s.stack = append(s.stack, linked)
}

func (s *stackCompiler) compile(symbols []symbol) *StateLinked {
	head := &StateLinked{
		transitions1: nil,
	}
	s.push(head) // starting state

	for len(symbols) > 0 {
		symbol := symbols[0]
		switch symbol.symbolType {
		case LParen:
			s.push(&StateLinked{})
		case RParen:
			s1 := s.pop()
			s2 := s.pop()
			s2Tail := tail(s2)
			s2Tail.merge(s1)
			s.push(s2)
		// branch
		case Pipe:
			s1 := s.pop()
			// transition with empty 'to' will be start of new branch
			s1.transitions1 = append([]transitionLinked{{to: nil}}, s1.transitions1...)
			s.push(s1)
		// concatenation
		default:
			s1 := s.pop()
			if symbol.symbolType == Character {
				predicate := func(r rune) bool { return r == symbol.letter }
				s.catenate(symbol, s1, predicate, "to -> "+string(symbol.letter), s1)
			} else if symbol.symbolType == AnyCharacter {
				predicate := func(r rune) bool { return true }
				s.catenate(symbol, tail(s1), predicate, "to -> .", s1)
			}
			s.push(s1)
		}
		symbols = symbols[1:]
	}
	return head
}

func (s *stackCompiler) catenate(symbol symbol, tail *StateLinked, predicate func(r rune) bool, description string, s1 *StateLinked) {
	if symbol.modifier == ZeroOrMore {
		tail.transitions1 = []transitionLinked{
			{to: nil},
			{to: tail, predicate: predicate, description: description},
		}
	} else {
		// will append to end of transition chain with index 0
		s.append(s1, &StateLinked{}, predicate, description)
	}
}

func tail(s *StateLinked) *StateLinked {
	head := s
	for len(head.transitions1) > 0 && head.transitions1[0].to != nil {
		head = head.transitions1[0].to
	}
	return head
}

func (s *stackCompiler) append(s1 *StateLinked, s2 *StateLinked, predicate Predicate, description string) {
	head := s1
	for len(head.transitions1) > 0 && head.transitions1[0].to != nil {
		head = head.transitions1[0].to
	}

	t := transitionLinked{
		to:          s2,
		predicate:   predicate,
		description: description,
	}

	if head.transitions1 != nil && head.transitions1[0].to == nil {
		head.transitions1[0] = t
	} else {
		head.transitions1 = append(head.transitions1, t)
	}
}

type modifier string

const (
	ZeroOrMore modifier = "ZeroOrMore"
)

type SymbolType int

const (
	Other SymbolType = iota
	AnyCharacter
	Pipe
	LParen
	RParen
	Character
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
		s.symbolType = AnyCharacter
	case '|':
		s.symbolType = Pipe
	default:
		s.symbolType = Character
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
