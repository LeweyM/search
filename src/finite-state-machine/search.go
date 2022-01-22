package finite_state_machine

type result struct {
	start, end int
}

func FindAll(finiteStateMachine *machine, searchString string) []result {
	var results []result

	start := 0
	end := 0
	// not using iterator as i here as it counts bytes, not runes
	for _, char := range searchString {
		currentState := finiteStateMachine.Next(char)
		if currentState == Success {
			results = append(results, result{start: start, end: end})
			finiteStateMachine.Reset()
			start = end + 1
		}
		if currentState == Fail {
			finiteStateMachine.Reset()
			start = end + 1
		}
		end++
	}

	return results
}
