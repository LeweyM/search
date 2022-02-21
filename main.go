package main

import (
	"context"
	"fmt"
	"io"
	"os"
	screen "search/src/screen"
	"search/src/search"
	"strconv"
	"strings"
)

const RED_ANSI = "\u001b[31m"
const RESET_ANSI = "\u001b[0m"

func main() {
	input := make(chan string)
	io := os.Stdin
	sc := screen.NewScreen(os.Stdout, input)
	ctx := context.Background()

	sc.Run(ctx, io, func() { os.Exit(1) })

	list(ctx, input, sc)
	//count(ctx, input, sc)
}

type displayState struct {
	Ready        bool
	Lines        []string
	TotalResults int
	Done         bool
	Target       string
}

type Screen interface {
	SetLines(lines []string)
	Run(ctx context.Context, reader io.Reader, exit func())
	SetState(state interface{})
	SetTemplate(templateString string)
}

func list(ctx context.Context, input chan string, sc Screen) {
	se := search.NewSearch("./bible.txt")
	se.LoadInMemory()
	se.LoadLinesInMemory()

	currentQuery := ""

	state := displayState{}
	templateString := `Input: {{ .Target }}
{{ if not .Ready }}Enter 3 letters or more to search.
{{ else }}{{ range $i, $line := .Lines }}
{{ $line }}{{ end }}
{{ if gt .TotalResults 10 }}{{ .TotalResults }} total results{{ end }}
... Searching{{ end }}`
	sc.SetTemplate(templateString)
	sc.SetState(state)

	results := make(chan search.Result)
	var queryResults []search.Result
	cancel, cancelFunc := context.WithCancel(ctx)
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
				queryResults = []search.Result{}
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
			queryResults = append(queryResults, r)
			if r.Finished {
				state.Done = true
			} else {
				state = getState(queryResults, state, offset)
			}
			sc.SetState(state)
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

func getState(queryResults []search.Result, state displayState, offset int) displayState {
	if len(queryResults) > 10 {
		state.Lines = formatLines(queryResults[offset:offset+10], offset)
	} else {
		state.Lines = formatLines(queryResults, 0)
	}
	state.TotalResults = len(queryResults)
	return state
}

func formatLines(results []search.Result, offset int) []string {
	res := make([]string, 0, len(results))
	for i, r := range results {
		res = append(res, fmt.Sprintf("%d: line-%s: \"%s\"", i+1+offset, strconv.Itoa(r.LineNumber), buildLine(r)))
	}
	return res
}

func buildLine(r search.Result) string {
	beforeMatch := r.LineContent[:r.Result.Start]
	match := r.LineContent[r.Result.Start : r.Result.End+1]
	afterMatch := r.LineContent[r.Result.End+1:]
	line := beforeMatch + RED_ANSI + match + RESET_ANSI + afterMatch
	line = strings.ReplaceAll(line, "\n", " \\n ")
	if len(line) > 300 {
		line = line[:300] + RESET_ANSI
	}
	return line
}

func count(ctx context.Context, input chan string, sc Screen) {
	se := search.NewSearch("./dict.txt")
	se.LoadInMemory()

	results := make(chan int, 0)
	waiting := make(chan bool, 0)

	// run new search on each valid input
	go func() {
		cancel, cancelFunc := context.WithCancel(ctx)
		for t := range input {
			cancelFunc()
			cancel, cancelFunc = context.WithCancel(ctx)
			if len(t) >= 3 {
				go se.Count(cancel, t, results)
			} else {
				waiting <- true
			}
		}
		defer cancelFunc()
	}()

	// update count on each result
	go func() {
		for r := range results {
			sc.SetLines([]string{"count: " + strconv.Itoa(r)})
		}
	}()

	// set to waiting
	go func() {
		for range waiting {
			sc.SetLines([]string{"Search for 3 more letters"})
		}
	}()
}
