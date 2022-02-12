package main

import (
	"context"
	"fmt"
	"io"
	"os"
	screen "search/src/screen"
	"search/src/search"
	"strconv"
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
	ready  bool
	lines  []string
	done   bool
	target string
}

type Screen interface {
	SetLines(lines []string)
	Run(ctx context.Context, reader io.Reader, exit func())
}

func list(ctx context.Context, input chan string, sc Screen) {
	se := search.NewSearch("./dict.txt")
	se.LoadInMemory()
	se.LoadLinesInMemory()

	currentQuery := ""

	state := displayState{}
	show := func(dState displayState) {
		if !dState.ready {
			sc.SetLines([]string{"Enter 3 letters or more to search."})
			return
		}

		for i, line := range dState.lines {
			if len(line) > 50 {
				dState.lines[i] = dState.lines[i][:50] + "..."
			}
		}

		if len(dState.lines) > 10 {
			display := append(dState.lines[:10], fmt.Sprintf("... %d more results not shown for query \"%s\"", len(dState.lines)-10, dState.target))
			if dState.done {
				sc.SetLines(append(display, fmt.Sprintf("Search finished. %d matching entries found.", len(dState.lines))))
			} else {
				sc.SetLines(append(display, "... Searching"))
			}
		} else {
			if dState.done {
				sc.SetLines(append(dState.lines, fmt.Sprintf("Search finished. %d matching entries found.", len(dState.lines))))
			} else {
				sc.SetLines(append(dState.lines, "... Searching"))
			}
		}
	}
	show(state)

	results := make(chan search.Result)
	cancel, cancelFunc := context.WithCancel(ctx)

	for {
		select {
		// return if outer context is cancelled
		case <-ctx.Done():
			cancelFunc()
			return
		case t := <-input:
			state = displayState{}
			cancelFunc()
			cancel, cancelFunc = context.WithCancel(ctx)
			if len(t) >= 3 {
				state.ready = true
				state.target = t
				currentQuery = t
				go se.SearchRegex(cancel, t, results)
			} else {
				state.ready = false
			}
			show(state)
		case r := <-results:
			if r.Query == currentQuery {
				if r.Finished {
					state.done = true
				} else {
					state.lines = append(state.lines, fmt.Sprintf("%d: line-%s: \"%s\"", len(state.lines), strconv.Itoa(r.LineNumber), r.LineContent))
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
