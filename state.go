package v1

type Status string

const (
	Success Status = "success"
	Fail           = "fail"
	Normal         = "normal"
)

type State struct {
	id          int
	transitions []Transition
}

func (s *State) firstMatchingTransition(input rune) *State {
	for _, t := range s.transitions {
		if t.predicate(input) {
			return t.to
		}
	}

	return nil
}

func (s *State) isSuccessState() bool {
	if len(s.transitions) == 0 {
		return true
	}

	return false
}
