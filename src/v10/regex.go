package v10

type Reducer interface {
	reduce(s *State)
}

type myRegex struct {
	fsm *State
}

type debugStep struct {
	runnerDrawing         string
	currentCharacterIndex int
}

func NewMyRegex(re string, reducers ...Reducer) *myRegex {
	tokens := lex(re)
	parser := NewParser(tokens)
	ast := parser.Parse()
	state, endState := ast.compile()
	endState.SetSuccess()
	for _, reducer := range reducers {
		reducer.reduce(state)
	}

	return &myRegex{fsm: state}
}

func (m *myRegex) MatchString(input string) bool {
	runner := NewRunner(m.fsm)
	return match(runner, []rune(input), nil, 0)
}

func (m *myRegex) DebugFSM() string {
	graph, _ := m.fsm.Draw()
	return graph
}

func (m *myRegex) DebugMatch(input string) []debugStep {
	runner := NewRunner(m.fsm)
	debugStepChan := make(chan debugStep)
	go func() {
		match(runner, []rune(input), debugStepChan, 0)
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
	if debugChan != nil {
		debugChan <- debugStep{runnerDrawing: runner.drawSnapshot(), currentCharacterIndex: offset}
	}

	for i, character := range input {
		runner.Next(character)
		runner.Start()
		if debugChan != nil {
			debugChan <- debugStep{runnerDrawing: runner.drawSnapshot(), currentCharacterIndex: offset + i + 1}
		}
		status := runner.GetStatus()

		if status == Success {
			return true
		}
	}

	return runner.GetStatus() == Success
}
