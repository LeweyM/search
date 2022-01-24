package finite_state_machine

const InitialState = 1

type StateType int

const (
	Success StateType = iota
	Fail
	Normal
)

type transition struct {
	to    int
	input rune
}

type machineTransition struct {
	next *machine
	to   int
}

type wildTransition int

type state struct {
	transitions        []transition
	wildTransition     *wildTransition
	stateType          StateType
	machineTransitions []machineTransition
}

type machine struct {
	states       []state
	currentState int
}

func NewMachine(n int) *machine {
	var states []state
	states = append(states, state{
		transitions: []transition{},
		stateType:   Fail,
	})
	for i := 0; i < n; i++ {
		var transitions []transition
		states = append(states, state{
			transitions: transitions,
			stateType:   Normal,
		})
	}
	return &machine{
		states:       states,
		currentState: InitialState,
	}
}

func (m *machine) AddTransition(from, to int, input rune) *machine {
	m.states[from].transitions = append(m.states[from].transitions, transition{
		to:    to,
		input: input,
	})
	return m
}

func (m *machine) SetSuccess(state int) *machine {
	m.states[state].stateType = Success
	return m
}

func (m *machine) Next(input rune) StateType {
	var hasTransition bool
	for _, t := range m.states[m.currentState].transitions {
		if t.input == input {
			hasTransition = true
			m.currentState = t.to
			break
		}
	}

	// if state has machine transitions, delegate to those sub machines
	currState := m.currentState
	failingSubstatesCount := 0
	for i := range m.states[currState].machineTransitions {
		stateType := m.states[currState].machineTransitions[i].next.Next(input)
		if stateType == Fail {
			failingSubstatesCount++
			if failingSubstatesCount == len(m.states[currState].machineTransitions) {
				hasTransition = false
				break
			}
		} else if stateType == Success {
			m.currentState = m.states[currState].machineTransitions[i].to
			// jumps at the first successful sub machine
			hasTransition = true
			break
		}
		hasTransition = true
	}

	// Decision: normal transitions take precedence over wild transitions
	if !hasTransition && m.states[m.currentState].wildTransition != nil {
		m.currentState = int(*m.states[m.currentState].wildTransition)
		hasTransition = true
	}

	// having no transition means the fsm enters a failed state
	if !hasTransition {
		m.currentState = 0
	}

	return m.states[m.currentState].stateType
}
func (m *machine) Reset() {
	m.currentState = InitialState

	// Also reset all sub machines
	for i := range m.states {
		for j := range m.states[i].machineTransitions {
			m.states[i].machineTransitions[j].next.Reset()
		}
	}
}

func (m *machine) AddWildTransition(from, to int) *machine {
	w := wildTransition(to)
	m.states[from].wildTransition = &w
	return m
}

func (m *machine) AddMachineTransition(from, to int, next *machine) *machine {
	m.states[from].machineTransitions = append(m.states[from].machineTransitions, machineTransition{
		next: next,
		to:   to,
	})
	return m
}
