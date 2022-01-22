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

type wildTransition int

type state struct {
	transitions    []transition
	wildTransition *wildTransition
	stateType      StateType
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
}

func (m *machine) AddWildTransition(from, to int) *machine {
	w := wildTransition(to)
	m.states[from].wildTransition = &w
	return m
}
