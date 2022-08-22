package v5

type myRegex struct {
	fsm *State
}

func NewMyRegex(re string) *myRegex {
	tokens := lex(re)
	parser := NewParser()
	ast := parser.Parse(tokens)
	state, _ := ast.compile()
	return &myRegex{fsm: state}
}

func (m *myRegex) MatchString(input string) bool {
	testRunner := NewRunner(m.fsm)
	return match(testRunner, []rune(input))
}

func match(runner *runner, input []rune) bool {
	runner.Reset()

	for _, character := range input {
		runner.Next(character)
		status := runner.GetStatus()

		if status == Fail {
			return match(runner, input[1:])
		}

		if status == Success {
			return true
		}
	}

	return runner.GetStatus() == Success
}
