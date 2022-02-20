package finite_state_machine

type runner struct {
	head      *StateLinked
	failState *StateLinked
	branches  *branchSet
}

func NewRunner(head *StateLinked) *runner {
	failState := &StateLinked{id: 0}

	return &runner{
		failState: failState,
		head:      head,
		branches:  newBranchSet(),
	}
}

func (r *runner) Next(input rune) StateType {
	// move along epsilon transitions first.
	// This is probably inefficient and could be moved into the main loop.
	r.processEpsilons()

	// move along regular transitions
	var nonFailedBranches = newBranchSet()
	for br := range r.branches.set {
		for _, destinationState := range br.matchingTransitions(input) {
			nonFailedBranches.add(destinationState)
		}
	}
	r.branches = nonFailedBranches

	// move along epsilon transitions after
	r.processEpsilons()

	return r.getTotalState()
}

func (r *runner) getTotalState() StateType {
	// if all branches have failed, return Fail
	if len(r.branches.set) == 0 {
		return Fail
	}
	// if any of the branches are success, return Success
	for b := range r.branches.set {
		if b.isSuccessState() {
			return Success
		}
	}
	// else, return normal
	return Normal
}

func (r *runner) processEpsilons() {
	// issue here is that multiple epsilons need to be processed in a chain
	var hasEpsilon = true
	for hasEpsilon {
		hasEpsilon = r.stepEpsilons()
	}
}

func (r *runner) stepEpsilons() (hasEpsilonAdvanced bool) {
	nextBranches := newBranchSet()
	for br := range r.branches.set {
		// if a branch contains an epsilon transition
		for _, t := range br.transitions {
			if t.epsilon && !r.branches.contains(t.to) {
				// add it to a branch
				nextBranches.add(t.to)
				hasEpsilonAdvanced = true
			}
		}
		nextBranches.add(br)
	}
	r.branches = nextBranches
	return hasEpsilonAdvanced
}

func (r *runner) Reset() {
	r.branches = newBranchSet()
	r.branches.add(r.head)
}

func (r *runner) onFailState(b *StateLinked) bool {
	return b == r.failState
}
