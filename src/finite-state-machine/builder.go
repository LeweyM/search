package finite_state_machine

import "fmt"

type builder struct {
	states []*State
}

var GlobalIdCounter = 0

func NewStateBuilder() *builder {
	var states []*State
	states = append(states, &State{id: 0}) // stand in for fail state
	return &builder{states: states}
}

func (b *builder) AddTransition(from, to int, letter rune) *builder {
	b.fillEmptyStatesTo(to)
	b.fillEmptyStatesTo(from)
	b.states[from].transitions = append(b.states[from].transitions, Transition{
		description: fmt.Sprintf("Matches: '%s'", string(letter)),
		to:          b.states[to],
		predicate:   func(input rune) bool { return input == letter },
	})
	return b
}

func (b *builder) AddWildTransition(from, to int) *builder {
	b.fillEmptyStatesTo(to)
	b.fillEmptyStatesTo(from)
	b.states[from].transitions = append(b.states[from].transitions, Transition{
		description: fmt.Sprintf("Matches anything"),
		to:          b.states[to],
		predicate:   func(input rune) bool { return true },
	})
	return b
}

func (b *builder) AddMachineTransition(from int, state *State) *builder {
	b.fillEmptyStatesTo(from)
	for _, t := range state.transitions {
		// when composing a transition, we merge the first transitions of the new state into the transition of the from state
		b.states[from].transitions = append(b.states[from].transitions, Transition{
			description: t.description,
			to:          t.to,
			predicate:   t.predicate,
		})
	}
	return b
}

func (b *builder) fillEmptyStatesTo(from int) {
	if from >= len(b.states) {
		for i := len(b.states); i <= from; i++ {
			GlobalIdCounter++
			b.states = append(b.states, &State{id: GlobalIdCounter})
		}
	}
}

func (b *builder) Build() *State {
	return b.states[1]
}
