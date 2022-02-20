package finite_state_machine

import "fmt"

const DEBUG = true

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
		transitions: nil,
	}
	s.push(head) // starting state

	for len(symbols) > 0 {
		symbol := symbols[0]
		switch symbol.symbolType {
		case LParen:
			s.push(&StateLinked{})
		case RParen:
			// x(ab)
			// inner (1) -a-> (2) -b-> (3)     -- main branch being processed inside the parens
			// outer (4) -x-> (5)              -- outer branch which has been pushed to the stack
			//
			// becomes
			//
			// (4) -x-> (5/1) -a-> (2) -b-> (3)  -- the end of 2nd branch is merged with beginning of 1st branch

			// 1. join the ends of the inner branches with epsilons
			inner := s.pop()
			innerTails := tails(inner)
			if len(innerTails) > 1 {
				for _, innerTail := range innerTails[1:] {
					firstInnerTail := innerTails[0]
					innerTail.transitions = append(innerTail.transitions, NewEpsilon(firstInnerTail))
				}
			}

			// 2. merge the end of the outer branch to the beginning of the inner branch
			outer := s.pop()
			outerTail := tail(outer)
			outerTail.merge(inner)
			// peek ahead
			if len(symbols) > 1 && symbols[1].symbolType == ZeroOrMore {
				// inner (1) -a-> (2) -b-> (3)     -- main branch being processed inside the parens
				// outer (4) -x-> (5)              -- outer branch which has been pushed to the stack
				//
				// becomes
				//
				// (4) -x-> (5/1) -a-> (2) -b-> (3)  -- the end of 2nd branch is merged with beginning of 1st branch
				//                ------------>      -- unconditional direct path to end state if 0 (ab)s
				//								     -- ? should there be a recursive path back also? It will not be rung as it is greedy however...
				outerTail.transitions = append(outerTail.transitions, TransitionLinked{to: tail(inner), predicate: func(input rune) bool { return true }, description: "to -> ."})
				symbols = symbols[1:]
			}
			s.push(outer)
		// branch
		case Pipe:
			s1 := s.pop()
			// transition with empty 'to' will be start of new branch
			s1.transitions = append([]TransitionLinked{{to: &StateLinked{empty: true}}}, s1.transitions...)
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
			oneBeforeTail.transitions[0].to = oneBeforeTail
			// make sure the main branch is nil
			oneBeforeTail.transitions = append([]TransitionLinked{{to: &StateLinked{empty: true}}}, oneBeforeTail.transitions...)
			s.push(s1)
		case OneOrMore:
			// (1) -a-> (2)				-- from a simple starting state
			// (1) -a-> (2) <-a- 		-- to a simple concatenation but with a recursive self definition
			s1 := s.pop()
			// grab the transition leading to the tail state. That is, grab the first transition from tail - 1.
			leadingTransition := tailN(s1, 1).transitions[0]
			// copy that transition to a secondary branch on the tail
			tail(s1).transitions = append([]TransitionLinked{{to: &StateLinked{empty: true}}}, leadingTransition)
			s.push(s1)
		case ZeroOrOne:
			// (1) -a-> (2)				-- from a simple starting state
			// (1) -E-> (2)
			//     -a->					-- add an epsilon transition to state 2
			s1 := s.pop()
			tail1 := tailN(s1, 1)
			epsilon := TransitionLinked{to: tail(s1), predicate: func(r rune) bool { return true }, description: "epsilon", epsilon: true}
			tail1.transitions = append(tail1.transitions, epsilon)
			s.push(s1)
		}
		symbols = symbols[1:]
	}
	return head
}

func getDescription(symbol symbol) string {
	var desc string
	if DEBUG {
		desc = fmt.Sprintf("to -> %s", string(symbol.letter))
	} else {
		desc = "to -> letter"
	}
	return desc
}

func tail(s *StateLinked) *StateLinked {
	return tailN(s, 0)
}

func tails(s *StateLinked) []*StateLinked {
	if s.empty || len(s.transitions) == 0 {
		return []*StateLinked{s}
	}

	var l []*StateLinked
	for _, t := range s.transitions {
		l = append(l, tails(t.to)...)
	}

	return l
}

func tailN(s *StateLinked, lag int) *StateLinked {
	head := s
	behind := s
	for len(head.transitions) > 0 && !head.transitions[0].to.empty {
		head = head.transitions[0].to
		if lag > 0 {
			lag--
		} else {
			behind = behind.transitions[0].to
		}
	}
	return behind
}

func (s *stackCompiler) append(s1 *StateLinked, s2 *StateLinked, predicate Predicate, description string) {
	t := TransitionLinked{
		to:          s2,
		predicate:   predicate,
		description: description,
	}

	if len(s1.transitions) > 0 && s1.transitions[0].to.empty {
		s1.transitions[0] = t
	} else {
		s1.transitions = append(s1.transitions, t)
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
