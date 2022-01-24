package finite_state_machine

type result struct {
	start, end int
}

func FindAll(finiteStateMachine *machine, searchString string) []result {
	var results []result

	start := 0
	end := 0
	runes := []rune(searchString)
	hasRerunFail := false
	for end < len(runes) {
		char := runes[end]
		currentState := finiteStateMachine.Next(char)
		switch currentState {
		case Success:
			results = append(results, result{start: start, end: end})
			finiteStateMachine.Reset()
			end++
			start = end
			break
		case Fail:
			finiteStateMachine.Reset()
			// in the case that a search fails, we want to rerun that char once in case the char that
			// fails one match is the beginning of another match
			if !hasRerunFail {
				hasRerunFail = true
			} else {
				end++
				hasRerunFail = false
			}
			start = end
			break
		default:
			end++
		}
	}
	return results
}
