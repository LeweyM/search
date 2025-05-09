package main

import (
	"context"
	"fmt"
	"github.com/LeweyM/search/src/search"
	"github.com/LeweyM/search/src/trigram"
	v10 "github.com/LeweyM/search/src/v10"
	v5 "github.com/LeweyM/search/src/v5"
	v6 "github.com/LeweyM/search/src/v6"
	v6BranchingIncomplete "github.com/LeweyM/search/src/v6-branching-incomplete"
	v6ParallelIncomplete "github.com/LeweyM/search/src/v6-parallel-incomplete"
	v7 "github.com/LeweyM/search/src/v7"
	v8 "github.com/LeweyM/search/src/v8"
	v8EpsilonConcatenation "github.com/LeweyM/search/src/v8-with-epsilon-concatenation"
	v9 "github.com/LeweyM/search/src/v9"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

const RedAnsi = "\u001b[31m"
const ResetAnsi = "\u001b[0m"

func main() {
	switch os.Args[1] {
	case "v5":
		v5.Main(os.Args[2:])
		return
	case "v6":
		v6.Main(os.Args[2:])
		return
	case "v6-branching-incomplete":
		v6BranchingIncomplete.Main(os.Args[2:])
		return
	case "v6-parallel-incomplete":
		v6ParallelIncomplete.Main(os.Args[2:])
		return
	case "v7":
		v7.Main(os.Args[2:])
		return
	case "v8":
		v8.Main(os.Args[2:])
		return
	case "v8-with-epsilon-concatenation":
		v8EpsilonConcatenation.Main(os.Args[2:])
		return
	case "v9":
		v9.Main(os.Args[2:])
		return
	case "v10":
		v10.Main(os.Args[2:])
		return
	default:
		panic(fmt.Sprintf("command '%s' not found", strings.Join(os.Args[1:], " ")))
	}

	//input := make(chan string)
	//inputOutput := os.Stdin
	//sc := screen.NewScreen(os.Stdout, input)
	//ctx := context.Background()
	//
	//sc.Run(ctx, inputOutput, func() { os.Exit(1) })
	//path := os.Args[1]
	//
	//file, err := os.Open(path)
	//if err != nil {
	//	panic(err)
	//}
	//fileInfo, err := file.Stat()
	//if err != nil {
	//	panic(err)
	//}
	//if fileInfo.IsDir() {
	//	listDir(ctx, input, sc, path)
	//} else {
	//	list(ctx, input, sc, path)
	//}
	return
}

type displayState struct {
	Ready        bool
	Lines        []string
	TotalResults int
	Done         bool
	Target       string
	Candidates   int
}

type Screen interface {
	SetLines(lines []string)
	Run(ctx context.Context, reader io.Reader, exit func())
	SetState(state interface{})
	SetTemplate(templateString string)
}

func listDir(ctx context.Context, input chan string, sc Screen, path string) {
	se := search.NewSearch(path)
	index := trigram.Index(path, false) // TODO: use a flag to write to json

	currentQuery := ""

	state := displayState{}
	templateString := `Input: {{ .Target }}
{{ if not .Ready }}Enter 3 letters or more to search.{{ else }}
Candidate Files: {{ .Candidates }}
{{ range $i, $line := .Lines }}
{{ $line }}{{ end }}
{{ if gt .TotalResults 10 }}{{ .TotalResults }} total results{{ end }}
{{ if not .Done }}... Searching{{ end }}{{ end }}`
	sc.SetTemplate(templateString)
	sc.SetState(state)

	results := make(chan search.ResultWithFile)
	var queryResults [][]search.ResultWithFile
	cancel, cancelFunc := context.WithCancel(ctx)
	previousLine := -1 // start at negative as there is no previous line at first
	previousFile := "" // start at empty as there is no previous file at first
	offset := 0
	for {
		select {
		// return if outer context is cancelled
		case <-ctx.Done():
			cancelFunc()
			return
		case t := <-input:
			switch t {
			case "LEFT", "RIGHT":
			case "DOWN":
				offset = min(len(queryResults)-10, offset+1)
				state = getStateNEW(queryResults, state, offset)
				sc.SetState(state)
			case "UP":
				offset = max(0, offset-1)
				state = getStateNEW(queryResults, state, offset)
				sc.SetState(state)
			default:
				state = displayState{Target: t}
				offset = 0
				previousLine = -1
				previousFile = ""
				queryResults = [][]search.ResultWithFile{}
				cancelFunc()
				cancel, cancelFunc = context.WithCancel(ctx)
				if len(t) >= 3 {
					state.Ready = true
					currentQuery = t
					candidates := index.Lookup(trigram.Query(t))
					state = updateCandidatesState(state, len(candidates))
					go se.SearchDirectoryRegexAsync(cancel, t, candidates, results)
				} else {
					state.Ready = false
				}
				sc.SetState(state)
			}
		case r := <-results:
			if r.Query != currentQuery {
				continue
			}
			// this assumes results from same line and file will come in sequentially
			if previousLine == r.LineNumber && r.File == previousFile {
				lastIndex := len(queryResults) - 1
				queryResults[lastIndex] = append(queryResults[lastIndex], r)
			} else {
				queryResults = append(queryResults, []search.ResultWithFile{r})
			}
			if r.Finished {
				state.Done = true
			} else {
				state = getStateNEW(queryResults, state, offset)
			}
			sc.SetState(state)
			previousLine = r.LineNumber
			previousFile = r.File
		}
	}
}

func updateCandidatesState(state displayState, candidates int) displayState {
	state.Candidates = candidates
	return state
}

func list(ctx context.Context, input chan string, sc Screen, path string) {
	se := search.NewSearch(path)
	se.LoadInMemory()
	se.LoadLinesInMemory()

	currentQuery := ""

	state := displayState{}
	templateString := `Input: {{ .Target }}
{{ if not .Ready }}Enter 3 letters or more to search.
{{ else }}{{ range $i, $line := .Lines }}
{{ $line }}{{ end }}
{{ if gt .TotalResults 10 }}{{ .TotalResults }} total results{{ end }}
{{ if not .Done }}... Searching{{ end }}{{ end }}`
	sc.SetTemplate(templateString)
	sc.SetState(state)

	results := make(chan search.Result)
	var queryResults [][]search.Result
	cancel, cancelFunc := context.WithCancel(ctx)
	previousLine := -1 // start at negative as there is no previous line at first
	offset := 0
	for {
		select {
		// return if outer context is cancelled
		case <-ctx.Done():
			cancelFunc()
			return
		case t := <-input:
			switch t {
			case "LEFT", "RIGHT":
			case "DOWN":
				offset = min(len(queryResults)-10, offset+1)
				state = getState(queryResults, state, offset)
				sc.SetState(state)
			case "UP":
				offset = max(0, offset-1)
				state = getState(queryResults, state, offset)
				sc.SetState(state)
			default:
				state = displayState{Target: t}
				offset = 0
				queryResults = [][]search.Result{}
				cancelFunc()
				cancel, cancelFunc = context.WithCancel(ctx)
				if len(t) >= 3 {
					state.Ready = true
					currentQuery = t
					go se.SearchRegex(cancel, t, results)
				} else {
					state.Ready = false
				}
				sc.SetState(state)
			}
		case r := <-results:
			if r.Query != currentQuery {
				continue
			}
			if previousLine == r.LineNumber {
				lastIndex := len(queryResults) - 1
				queryResults[lastIndex] = append(queryResults[lastIndex], r)
			} else {
				queryResults = append(queryResults, []search.Result{r})
			}
			if r.Finished {
				state.Done = true
			} else {
				state = getState(queryResults, state, offset)
			}
			sc.SetState(state)
			previousLine = r.LineNumber
		}
	}
}

func min(i, i2 int) int {
	if i < i2 {
		return i
	} else {
		return i2
	}
}

func max(i, i2 int) int {
	if i > i2 {
		return i
	} else {
		return i2
	}
}

func getStateNEW(queryResults [][]search.ResultWithFile, state displayState, offset int) displayState {
	if len(queryResults) > 10 {
		state.Lines = formatLinesNEW(queryResults[offset:offset+10], offset)
	} else {
		state.Lines = formatLinesNEW(queryResults, 0)
	}
	state.TotalResults = len(queryResults)
	return state
}

func formatLinesNEW(lineResults [][]search.ResultWithFile, offset int) []string {
	res := make([]string, 0, len(lineResults))
	for i, r := range lineResults {
		res = append(res, fmt.Sprintf(
			"%d: [file:%s] line-%s: \"%s\"",
			i+1+offset,
			r[0].File,
			strconv.Itoa(r[0].LineNumber),
			buildLine(r[0].LineContent, matchesFromResultsNEW(r)),
		))
	}
	return res
}

func matchesFromResultsNEW(results []search.ResultWithFile) (matches []search.Match) {
	for _, r := range results {
		matches = append(matches, r.Match)
	}
	return matches
}

func getState(queryResults [][]search.Result, state displayState, offset int) displayState {
	if len(queryResults) > 10 {
		state.Lines = formatLines(queryResults[offset:offset+10], offset)
	} else {
		state.Lines = formatLines(queryResults, 0)
	}
	state.TotalResults = len(queryResults)
	return state
}

func formatLines(lineResults [][]search.Result, offset int) []string {
	res := make([]string, 0, len(lineResults))
	for i, r := range lineResults {
		res = append(res, fmt.Sprintf("%d: line-%s: \"%s\"", i+1+offset, strconv.Itoa(r[0].LineNumber), buildLine(r[0].LineContent, matchesFromResults(r))))
	}
	return res
}

func matchesFromResults(results []search.Result) (matches []search.Match) {
	for _, r := range results {
		matches = append(matches, r.Match)
	}
	return matches
}

func buildLine(content string, matches []search.Match) string {
	var reducedMatches = reduceMatches(matches)
	var segments []string
	last := 0
	for _, m := range reducedMatches {
		segments = append(segments, content[last:m.Start])
		segments = append(segments, RedAnsi)
		segments = append(segments, content[m.Start:m.End+1])
		segments = append(segments, ResetAnsi)
		last = m.End + 1
	}
	segments = append(segments, content[last:])
	line := strings.Join(segments, "")
	line = strings.ReplaceAll(line, "\n", " \\n ")
	line = strings.TrimSpace(line)
	return line + ResetAnsi
}

func reduceMatches(matches []search.Match) []search.Match {
	var res []search.Match
	for _, match := range matches {
		i, isOverlapping := overlapAny(match, res)
		if isOverlapping {
			res[i] = merge(match, res[i])
		} else {
			res = append(res, match)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].End < res[j].End
	})
	return res
}

func merge(a, b search.Match) search.Match {
	return search.Match{
		Start: min(a.Start, b.Start),
		End:   max(a.End, b.End),
	}
}

func overlapAny(a search.Match, b []search.Match) (i int, ok bool) {
	for index, match := range b {
		if overlap(a, match) {
			return index, true
		}
	}
	return 0, false
}

func overlap(a, b search.Match) bool {
	// assumption: start is always less than end
	bStartsAfterA := b.Start > a.End
	aStartsAfterB := a.Start > b.End
	return !(bStartsAfterA || aStartsAfterB)
}
