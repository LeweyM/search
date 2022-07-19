package finite_state_machine

type branchSet struct {
	set map[*State]bool
}

func newBranchSet() *branchSet {
	return &branchSet{set: make(map[*State]bool)}
}

func (b *branchSet) add(state *State) {
	b.set[state] = true
}

func (b *branchSet) contains(state *State) bool {
	return b.set[state]
}

func (b *branchSet) remove(state *State) {
	delete(b.set, state)
}

func (b *branchSet) getSet() map[*State]bool {
	return b.set
}
