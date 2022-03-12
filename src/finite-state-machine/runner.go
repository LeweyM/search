package finite_state_machine

type runner struct {
	head      *State
	failState *State
	branches  *branchSet
}

func NewRunner(head *State) *runner {
	failState := &State{id: 0}

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
		for _, t := range br.transitions {
			// if a branch contains an epsilon transition
			if t.epsilon {
				// and the destination has not yet been accounted for
				if !r.branches.contains(t.to) {
					// add its destination branches to the branch set
					nextBranches.add(t.to)
					hasEpsilonAdvanced = true
				}
			}
		}
		// then add the branch
		nextBranches.add(br)
	}
	r.branches = nextBranches
	return hasEpsilonAdvanced
}

func (r *runner) Reset() {
	r.branches = newBranchSet()
	r.branches.add(r.head)
	// process epsilons at the starting state
	r.processEpsilons()
}

func (r *runner) onFailState(b *State) bool {
	return b == r.failState
}
