package finite_state_machine

type finiteState interface {
	test(r rune) finiteState
	isEndState() bool
}