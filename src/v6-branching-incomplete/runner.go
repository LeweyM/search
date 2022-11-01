package v6BranchingIncomplete

type runner struct {
	head    *State
	current *State
}

func NewRunner(head *State) *runner {
	r := &runner{
		head:    head,
		current: head,
	}

	return r
}

func (r *runner) Next(input rune) {
	if r.current == nil {
		return
	}

	// move to next matching transition
	r.current = r.current.firstMatchingTransition(input)
}

func (r *runner) GetStatus() Status {
	// if the current state is nil, return Fail
	if r.current == nil {
		return Fail
	}

	// if the current state has no transitions from it, return Success
	if r.current.isSuccessState() {
		return Success
	}

	// else, return normal
	return Normal
}

func (r *runner) Reset() {
	r.current = r.head
}
