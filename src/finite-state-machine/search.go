package finite_state_machine

import "context"

type Result struct {
	Line, Start, End int
}

type Machine interface {
	Next(input rune) StateType
	Reset()
}

func FindAllAsync(ctx context.Context, finiteStateMachine Machine, searchString string, out chan Result) {
	defer close(out)
	lineCounter := 1
	start := 0
	end := 0
	runes := append([]rune(searchString), 0) // we add a 'NULL' 0 rune at the End so that even empty string inputs are processed.
	hasRerunFail := false
	for end < len(runes) {
		select {
		case <-ctx.Done():
			return
		default:
			char := runes[end]
			if !hasRerunFail && char == '\n' {
				lineCounter++
			}
			currentState := finiteStateMachine.Next(char)
			switch currentState {
			case Success:
				out <- Result{Start: start, End: end, Line: lineCounter}
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
	}
}

type localResult struct {
	start, end int
}

func FindAll(finiteStateMachine Machine, searchString string) []localResult {
	var results []localResult
	resultChan := make(chan Result, 10)
	FindAllAsync(context.TODO(), finiteStateMachine, searchString, resultChan)
	for res := range resultChan {
		results = append(results, localResult{
			start: res.Start,
			end:   res.End,
		})
	}
	return results
}
