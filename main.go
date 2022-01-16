package main

import (
	"context"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	screen "search/src/screen"
	"search/src/search"
	"strconv"
)

func main() {
	fmt.Print("\033[?25l") // hide cursor
	terminal.MakeRaw(0)    // fd 0 is stdin

	input := make(chan string)
	sc := screen.NewScreen(os.Stdout, input)
	ctx := context.Background()

	sc.Run(ctx)

	list(ctx, input, sc)
	//count(ctx, input, sc)
}

type displayState struct {
	ready    bool
	lines    []string
}

func list(ctx context.Context, input chan string, sc screen.Screen) {
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
		if len(dState.lines) > 10 {
			sc.SetLines(append(dState.lines[:10], fmt.Sprintf("... %d more results not shown", len(dState.lines) - 10)))
		} else {
			sc.SetLines(dState.lines)
		}
	}
	show(state)

	results := make(chan search.Result)
	cancel, cancelFunc := context.WithCancel(ctx)

	for {
		select {
		case t := <-input:
			state.lines = []string{}
			cancelFunc()
			cancel, cancelFunc = context.WithCancel(ctx)
			if len(t) >= 3 {
				state.ready = true
				currentQuery = t
				go se.Search(cancel, t, results)
			} else {
				state.ready = false
			}
			show(state)
		case r := <-results:
			if r.Query == currentQuery {
				state.lines = append(state.lines, fmt.Sprintf("%d: line-%s: \"%s\"", len(state.lines), strconv.Itoa(r.LineNumber), r.LineContent))
				show(state)
			}
		}
	}
}

func count(ctx context.Context, input chan string, sc screen.Screen) {
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
