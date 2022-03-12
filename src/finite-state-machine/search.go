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
	lineStart := 0
	runes := append([]rune(searchString)) // we add a 'NULL' 0 rune at the End so that even empty string inputs are processed.
	hasRerunFail := false
	hasStartedMatch := false
	lastSuccessIndex := 0
	for end < len(runes) {
		select {
		case <-ctx.Done():
			return
		default:
			char := runes[end]
			if char == '\n' {
				// if result found, return until end of line
				if hasStartedMatch {
					out <- Result{Start: start - lineStart, End: lastSuccessIndex - lineStart, Line: lineCounter}
					hasStartedMatch = false
				}
				// Like grep, don't search for matches across lines.
				finiteStateMachine.Reset()
				lineCounter++
				end++
				lineStart = end
				start = end
				continue
			}
			currentState := finiteStateMachine.Next(char)
			switch currentState {
			case Success:
				if !hasStartedMatch {
					hasStartedMatch = true
				}
				lastSuccessIndex = end
				end++
				break
			case Fail:
				if hasStartedMatch {
					out <- Result{Start: start - lineStart, End: lastSuccessIndex - lineStart, Line: lineCounter}
					hasStartedMatch = false
				}
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

	if hasStartedMatch {
		out <- Result{Start: start - lineStart, End: lastSuccessIndex - lineStart, Line: lineCounter}
	}

	if len([]rune(searchString)) == 0 {
		currentState := finiteStateMachine.Next(0)
		if currentState == Success {
			out <- Result{Start: 0, End: 0, Line: 0}
		}
	}
}

type localResult struct {
	start, end int
}

type localResultWithLines struct {
	line, start, end int
}

func FindAllWithLines(finiteStateMachine Machine, searchString string) []localResultWithLines {
	var results []localResultWithLines
	resultChan := make(chan Result, 100)
	FindAllAsync(context.TODO(), finiteStateMachine, searchString, resultChan)
	for res := range resultChan {
		results = append(results, localResultWithLines{
			start: res.Start,
			end:   res.End,
			line:  res.Line,
		})
	}
	return results
}

func FindAll(finiteStateMachine Machine, searchString string) []localResult {
	var results []localResult
	resultChan := make(chan Result, 100)
	FindAllAsync(context.TODO(), finiteStateMachine, searchString, resultChan)
	for res := range resultChan {
		results = append(results, localResult{
			start: res.Start,
			end:   res.End,
		})
	}
	return results
}
