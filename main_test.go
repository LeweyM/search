package main

import (
	"context"
	"os"
	"search/src/screen"
	"strings"
	"testing"
)

func BenchmarkList(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testWithInput("aaaaaaaaaaaaaaaaaaaaaaaaaaa") // ~280058584 ns/op
	}
}

func testWithInput(inputString string) {
	input := make(chan string)
	sc := screen.NewScreen(os.Stdout, input)
	ctx, cancelFunc := context.WithCancel(context.Background())
	exit := func() {
		cancelFunc()
	}
	sc.Run(ctx, strings.NewReader(inputString+"Q"), exit)
	list(ctx, input, sc, "./data/bible/bible.txt")
}
