package finite_state_machine

import "fmt"

const DEBUG = false

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
			// s1 (1) -a-> (2) -b-> (3)     -- main branch being processed inside the parens
			// s2 (4) -x-> (5)              -- outer branch which has been pushed to the stack
			//
			// becomes
			//
			// (4) -x-> (5/1) -a-> (2) -b-> (3)  -- the end of 2nd branch is merged with beginning of 1st branch

			s1 := s.pop()
			s2 := s.pop()
			s2Tail := tail(s2)
			s2Tail.merge(s1)
			// peek ahead
			if len(symbols) > 1 && symbols[1].symbolType == ZeroOrMore {
				// s1 (1) -a-> (2) -b-> (3)     -- main branch being processed inside the parens
				// s2 (4) -x-> (5)              -- outer branch which has been pushed to the stack
				//
				// becomes
				//
				// (4) -x-> (5/1) -a-> (2) -b-> (3)  -- the end of 2nd branch is merged with beginning of 1st branch
				//                ------------>      -- unconditional direct path to end state if 0 (ab)s
				//								     -- ? should there be a recursive path back also? It will not be rung as it is greedy however...
				s2Tail.transitions1 = append(s2Tail.transitions1, transitionLinked{to: tail(s1), predicate: func(input rune) bool { return true }, description: "to -> ."})
				symbols = symbols[1:]
			}
			s.push(s2)
		// branch
		case Pipe:
			s1 := s.pop()
			// transition with empty 'to' will be start of new branch
			s1.transitions1 = append([]transitionLinked{{to: &StateLinked{empty: true}}}, s1.transitions1...)
			s.push(s1)
		// concatenation
		case Character:
			s1 := s.pop()
			s2Tail := tail(s1)
			next := &StateLinked{}
			s.append(s2Tail, next, func(r rune) bool { return r == symbol.letter }, getDescription(symbol))
			s.push(s1)
		case AnyCharacter:
			s1 := s.pop()
			s.append(tail(s1), &StateLinked{}, func(r rune) bool { return true }, "to -> .")
			s.push(s1)
		case ZeroOrMore:
			// (1) -x-> (2) - we start with a simple starting state and a condition to transfer to state(2) on an 'x'
			//
			// becomes:
			//
			// (1) -> (2)   - the starting state now unconditionally goes to the next state, as there may be 0 'x's
			//    <-x       - the old condition is used for the recursive loop, as there may be many 'x's
			s1 := s.pop()
			// make a tail of s1 with a loop to itself
			oneBeforeTail := tailN(s1, 1)
			oneBeforeTail.transitions1[0].to = oneBeforeTail
			// make sure the main branch is nil
			oneBeforeTail.transitions1 = append([]transitionLinked{{to: &StateLinked{empty: true}}}, oneBeforeTail.transitions1...)
			s.push(s1)
		case OneOrMore:
			// (1) -a-> (2)				-- from a simple starting state
			// (1) -a-> (2) <-a- 		-- to a simple concatenation but with a recursive self definition
			s1 := s.pop()
			// grab the transition leading to the tail state. That is, grab the first transition from tail - 1.
			leadingTransition := tailN(s1, 1).transitions1[0]
			// copy that transition to a secondary branch on the tail
			tail(s1).transitions1 = append([]transitionLinked{{to: &StateLinked{empty: true}}}, leadingTransition)
			s.push(s1)
		case ZeroOrOne:
			// (1) -a-> (2)				-- from a simple starting state
			// (1) -E-> (2)
			//     -a->					-- add an epsilon transition to state 2
			s1 := s.pop()
			tail1 := tailN(s1, 1)
			epsilon := transitionLinked{to: tail(s1), predicate: func(r rune) bool { return true }, description: "epsilon", epsilon: true}
			tail1.transitions1 = append(tail1.transitions1, epsilon)
			s.push(s1)
		}
		symbols = symbols[1:]
	}
	return head
}

func getDescription(symbol symbol) string {
	var desc string
	if DEBUG {
		desc = fmt.Sprintf("to -> %d", symbol.letter)
	} else {
		desc = "to -> letter"
	}
	return desc
}

func tail(s *StateLinked) *StateLinked {
	return tailN(s, 0)
}

func tailN(s *StateLinked, lag int) *StateLinked {
	head := s
	behind := s
	for len(head.transitions1) > 0 && !head.transitions1[0].to.empty {
		head = head.transitions1[0].to
		if lag > 0 {
			lag--
		} else {
			behind = behind.transitions1[0].to
		}
	}
	return behind
}

func (s *stackCompiler) append(s1 *StateLinked, s2 *StateLinked, predicate Predicate, description string) {
	t := transitionLinked{
		to:          s2,
		predicate:   predicate,
		description: description,
	}

	if len(s1.transitions1) > 0 && s1.transitions1[0].to.empty {
		s1.transitions1[0] = t
	} else {
		s1.transitions1 = append(s1.transitions1, t)
	}
}

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
