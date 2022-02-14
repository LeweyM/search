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
	Ready bool
	Lines  []string
	Done   bool
	Target string
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
... Searching
{{ end }}
	`
	sc.SetTemplate(templateString)

	show := func(dState displayState) {
		sc.SetState(dState)
	}
	show(state)

	results := make(chan search.Result)
	cancel, cancelFunc := context.WithCancel(ctx)

	resultCounter := 0
	for {
		select {
		// return if outer context is cancelled
		case <-ctx.Done():
			cancelFunc()
			return
		case t := <-input:
			state = displayState{}
			resultCounter = 0
			cancelFunc()
			cancel, cancelFunc = context.WithCancel(ctx)
			state.Target = t
			if len(t) >= 3 {
				state.Ready = true
				currentQuery = t
				go se.SearchRegex(cancel, t, results)
			} else {
				state.Ready = false
			}
			show(state)
		case r := <-results:
			if r.Query == currentQuery {
				resultCounter++
				if r.Finished {
					state.Done = true
				} else {
					if len(state.Lines) < 10 {
						line := strings.ReplaceAll(r.LineContent, "\n", " \\n ")
						state.Lines = append(state.Lines, fmt.Sprintf("%d: line-%s: \"%s\"", len(state.Lines), strconv.Itoa(r.LineNumber), line))
					} else {
						if len(state.Lines) < 11 {
							state.Lines = append(state.Lines, "1 more element not shown")
						}
						state.Lines[10] =  fmt.Sprintf("%d more elements not shown", resultCounter)

					}
				}
				show(state)
			}
		}
	}
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
