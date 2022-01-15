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

	// loop forever
	for {}

}

func list(ctx context.Context, input chan string, sc screen.Screen) {
	se := search.NewSearch("./dict.txt")
	se.LoadInMemory()
	se.LoadLinesInMemory()

	results := make(chan search.Result, 10)
	waiting := make(chan bool, 10)
	reset := make(chan bool, 10)

	// run new search on each valid input
	cancel, cancelFunc := context.WithCancel(ctx)
	go func() {
		for t := range input {
			cancelFunc()
			reset <- true
			cancel, cancelFunc = context.WithCancel(ctx)

			if len(t) >= 5 {
				go se.Search(cancel, t, results)
			} else {
				waiting <- true
			}
		}
		defer cancelFunc()
	}()

	// update count on each result
	count := 0
	go func() {
		for r := range results {
			if count < 10 {
				sc.AddLine(fmt.Sprintf("%d: line-%s: \"%s\"", r.Count, strconv.Itoa(r.LineNumber), r.LineContent))
			} else {
				sc.SetLine(10, fmt.Sprintf("... %d more results not shown", count-10))
			}
			count++
		}
	}()

	// set to waiting
	go func() {
		for range waiting {
			count = 0
			sc.SetLines([]string{"Search for 5 more letters"})
		}
	}()

	go func() {
		for range reset {
			count = 0
			sc.SetLines([]string{})
		}
	}()
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
