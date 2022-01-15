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

	se := search.NewSearch("./2mb-sample.txt")
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

	// loop forever
	for {}

}
