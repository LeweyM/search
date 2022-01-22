package finite_state_machine

import "fmt"

type Builder struct {
	states []*matchesLetter
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (f *Builder) State(state matchesLetter) *Builder {
	if len(f.states) > 0 {
		f.lastState().nextState = &state
	}
	f.states = append(f.states, &state)
	state = state.Base(f.states[0])
	return f
}

func (f *Builder) End() *Builder {
	s := f.lastState()
	s.endState = true
	return f
}

func (f *Builder) Recursive() *Builder {
	s := f.lastState()
	s.nextState = s
	return f
}

func (f *Builder) To(i int) *Builder {
	f.lastState().nextState = f.states[i]
	return f
}

func (f *Builder) Build() (finiteState, error) {
	for i, state := range f.states {
		if state.base == nil {
			return nil, fmt.Errorf("state(%d) does not have a base state", i)
		}
		if state.nextState == nil {
			return nil, fmt.Errorf("state(%d) does not have a next state", i)
		}
	}
	return f.states[0], nil
}

func (f *Builder) lastState() *matchesLetter {
	return f.states[len(f.states)-1]
}
