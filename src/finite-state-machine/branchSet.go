package finite_state_machine

type branchSet struct{ set map[*StateLinked]bool }

func newBranchSet() *branchSet {
	return &branchSet{set: make(map[*StateLinked]bool)}
}

func (b *branchSet) add(state *StateLinked) {
	b.set[state] = true
}

func (b *branchSet) contains(state *StateLinked) bool {
	return b.set[state]
}

func (b *branchSet) remove(state *StateLinked) {
	delete(b.set, state)
}

func (b *branchSet) getSet() map[*StateLinked]bool {
	return b.set
}
