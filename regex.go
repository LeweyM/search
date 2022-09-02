package v5

type myRegex struct {
	fsm *State
}

type debugStep struct {
	runnerDrawing         string
	currentCharacterIndex int
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
	return match(testRunner, []rune(input), dummyChan[debugStep](), 0)
}

func (m *myRegex) DebugFSM() string {
	graph, _ := m.fsm.Draw()
	return graph
}

func (m *myRegex) DebugMatch(input string) []debugStep {
	testRunner := NewRunner(m.fsm)
	debugStepChan := make(chan debugStep)
	go func() {
		match(testRunner, []rune(input), debugStepChan, 0)
		close(debugStepChan)
	}()
	var debugSteps []debugStep
	for step := range debugStepChan {
		debugSteps = append(debugSteps, step)
	}

	return debugSteps
}

func match(runner *runner, input []rune, debugChan chan debugStep, offset int) bool {
	runner.Reset()
	debugChan <- debugStep{runnerDrawing: runner.drawCurrentState(), currentCharacterIndex: offset}

	for i, character := range input {
		runner.Next(character)
		debugChan <- debugStep{runnerDrawing: runner.drawCurrentState(), currentCharacterIndex: offset + i + 1}
		status := runner.GetStatus()

		if status == Fail {
			return match(runner, input[1:], debugChan, offset+1)
		}

		if status == Success {
			return true
		}
	}

	return runner.GetStatus() == Success
}

func dummyChan[T any]() chan T {
	dummyChan := make(chan T)
	go func() {
		for {
			<-dummyChan
		}
	}()
	return dummyChan
}
