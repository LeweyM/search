package finite_state_machine

import "fmt"

type builder struct {
	states []*StateLinked
}

var GlobalIdCounter = 0

func NewStateLinkedBuilder() *builder {
	var states []*StateLinked
	states = append(states, &StateLinked{id: 0}) // stand in for fail state
	return &builder{states: states}
}

func (b *builder) AddTransition(from, to int, letter rune) *builder {
	b.fillEmptyStatesTo(to)
	b.fillEmptyStatesTo(from)
	b.states[from].transitions1 = append(b.states[from].transitions1, transitionLinked{
		description: fmt.Sprintf("Matches: '%s'", string(letter)),
		to:          b.states[to],
		predicate:   func(input rune) bool { return input == letter },
	})
	return b
}

func (b *builder) AddWildTransition(from, to int) *builder {
	b.fillEmptyStatesTo(to)
	b.fillEmptyStatesTo(from)
	b.states[from].transitions1 = append(b.states[from].transitions1, transitionLinked{
		description: fmt.Sprintf("Matches anything"),
		to:          b.states[to],
		predicate:   func(input rune) bool { return true },
	})
	return b
}

func (b *builder) AddMachineTransition(from int, state *StateLinked) *builder {
	b.fillEmptyStatesTo(from)
	for _, t := range state.transitions1 {
		// when composing a transition, we merge the first transitions of the new state into the transition of the from state
		b.states[from].transitions1 = append(b.states[from].transitions1, transitionLinked{
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
			b.states = append(b.states, &StateLinked{id: GlobalIdCounter})
		}
	}
}

func (b *builder) Build() *StateLinked {
	return b.states[1]
}
