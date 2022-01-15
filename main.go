package main

import (
	"context"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	screen "search/src"
	"strconv"
	"time"
)

func main() {
	fmt.Print("\033[?25l") // hide cursor
	terminal.MakeRaw(0) // fd 0 is stdin

	screen := screen.NewScreen(os.Stdout)
	ctx := context.Background()

	screen.Run(ctx)

	count := 0
	for {
		count++
		screen.SetLines([]string{"bobo", "something: " + strconv.Itoa(count), "foo"})
		time.Sleep(100 * time.Millisecond)
	}
}
