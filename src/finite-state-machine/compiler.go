package finite_state_machine

func Compile(input string) *StateLinked {
	n := len(input) + 1
	builder := NewStateLinkedBuilder(n)
	symbols := lex(input)
	prevStateNumber := 1
	for _, symbol := range symbols {
		if symbol.wild {
			if symbol.modifier == "ANY" {
				builder = builder.AddWildTransition(prevStateNumber, prevStateNumber)
			} else {
				builder = builder.AddWildTransition(prevStateNumber, prevStateNumber+1)
				prevStateNumber++
			}
		} else {
			builder = builder.AddTransition(prevStateNumber, prevStateNumber+1, symbol.letter)
			prevStateNumber++
		}
	}
	builder = builder.SetSuccess(prevStateNumber)
	stateLinked := builder.Build()
	return stateLinked
}

type modifier string

type symbol struct {
	wild     bool
	letter   rune
	modifier modifier
}

func lex(input string) []symbol {
	var symbols []symbol
	i := 0
	for i < len(input) {
		hasNext := i+1 < len(input)
		if hasNext {
			if isModifier(rune(input[i+1])) {
				symbols = append(symbols, lexRuneWithModifier(rune(input[i]), rune(input[i+1])))
				i += 2
			} else {
				symbols = append(symbols, lexRune(rune(input[i])))
				i += 1
			}
		} else {
			symbols = append(symbols, lexRune(rune(input[i])))
			i++
		}
	}
	return symbols
}

func lexRune(r rune) symbol {
	var s symbol
	if r == '.' {
		s.wild = true
	} else {
		s.letter = r
	}
	return s
}

func lexRuneWithModifier(r rune, mod rune) symbol {
	s := lexRune(r)
	if mod == '*' {
		s.modifier = "ANY"
	}
	return s
}

func isModifier(r rune) bool {
	return r == '*'
}
