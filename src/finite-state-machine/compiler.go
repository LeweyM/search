package finite_state_machine

import "fmt"

const DEBUG = true

func Compile(input string) *State {
	symbols := lex(input)
	compiler := stackCompiler{}
	stateLinked := compiler.compile(symbols)
	return stateLinked
}

type stackCompiler struct {
	stack     []*State
	// tailStack keeps track of the tail - 1 of the current machine
	tailStack StateStack
}

func (s *stackCompiler) pop() *State {
	i := len(s.stack) - 1
	state := s.stack[i]
	s.stack = s.stack[:i]
	return state
}

func (s *stackCompiler) push(linked *State) {
	s.stack = append(s.stack, linked)
}

func (s *stackCompiler) compile(symbols []symbol) *State {
	s.tailStack = StateStack{}
	head := &State{
		transitions: nil,
	}
	// starting state
	s.push(head)
	s.tailStack.push(head)

	for len(symbols) > 0 {
		symbol := symbols[0]
		switch symbol.symbolType {
		case LParen:
			newState := &State{}
			s.push(newState)
			s.tailStack.push(newState)
		case RParen:
			s.tailStack.pop()
			t := s.tailStack.pop()
			if t.transitions != nil {
				s.tailStack.push(t.transitions[0].to)
			} else {
				s.tailStack.push(t)
			}
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
			s.push(outer)
		// branch
		case Pipe:
			s1 := s.pop()
			// transition with empty 'to' will be start of new branch
			s1.transitions = append([]Transition{{to: &State{empty: true}}}, s1.transitions...)
			s.push(s1)
		// concatenation
		case Character:
			s1 := s.pop()
			s1Tail := tail(s1)
			next := &State{}
			s.append(s1Tail, next, func(r rune) bool { return r == symbol.letter }, getDescription(symbol))
			s.push(s1)

			s.tailStack.pop()
			s.tailStack.push(s1Tail)
		case AnyCharacter: // '.'
			s1 := s.pop()
			s1Tail := tail(s1)
			next := &State{}
			s.append(s1Tail, next, func(r rune) bool { return true }, "to -> .")
			s.push(s1)

			s.tailStack.pop()
			s.tailStack.push(s1Tail)
		case ZeroOrMore: // '*'
			tail0 := s.tailStack.pop()
			s.tailStack.push(tail0)

			// (1) -x-> (2)
			// becomes:
			// (1) -x-> (2)
			//   --ep1->
			//   <--ep2-
			s1 := s.pop()
			s2 := tail(s1)
			epsilon1 := Transition{to: s2, predicate: func(r rune) bool { return true }, description: "epsilon", epsilon: true}
			epsilon2 := Transition{to: tail0, predicate: func(r rune) bool { return true }, description: "epsilon", epsilon: true}
			epsilon3 := Transition{to: &State{}, predicate: func(r rune) bool { return true }, description: "epsilon", epsilon: true}
			tail0.transitions = append(tail0.transitions, epsilon1)
			s2.transitions = append(s2.transitions, epsilon2)
			s2.transitions = append([]Transition{epsilon3}, s2.transitions...) // epsilon3 (to end state) should be first for tail searches
			s.push(s1)
		case OneOrMore: // '+'
			// (1) -a-> (2)				-- from a simple starting state
			// (1) -a-> (2) <-a- 		-- to a simple concatenation but with a recursive self definition
			s1 := s.pop()
			// grab the transition leading to the tail state. That is, grab the first transition from tail - 1.
			leadingTransition := tailN(s1, 1).transitions[0]
			// copy that transition to a secondary branch on the tail
			tail(s1).transitions = append([]Transition{{to: &State{empty: true}}}, leadingTransition)
			s.push(s1)
		case ZeroOrOne: // '?'
			// (1) -a-> (2)				-- from a simple starting state
			// becomes
			// (1) -a-> (2)
			//     -E->	(2)				-- add an epsilon transition to end state
			s1 := s.pop()
			tail0 := s.tailStack.pop()
			s1Tail := tail(s1)

			epsilon := Transition{to: s1Tail, predicate: func(r rune) bool { return true }, description: "epsilon", epsilon: true}
			tail0.transitions = append([]Transition{epsilon}, tail0.transitions...)

			s.tailStack.push(tail0)
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

func tail(s *State) *State {
	return tailN(s, 0)
}

func tails(s *State) []*State {
	if s.empty || len(s.transitions) == 0 {
		return []*State{s}
	}

	var l []*State
	for _, t := range s.transitions {
		l = append(l, tails(t.to)...)
	}

	return l
}

func tailN(s *State, lag int) *State {
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

func (s *stackCompiler) append(s1 *State, s2 *State, predicate Predicate, description string) {
	t := Transition{
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
