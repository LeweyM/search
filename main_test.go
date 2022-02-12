package main

import (
	"context"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"search/src/screen"
	"strings"
	"testing"
)

func BenchmarkList(b *testing.B) {
	state, _ := terminal.MakeRaw(0) // fd 0 is stdin
	for i := 0; i < b.N; i++ {
		testWithInput("aaaaaaaaaaaaaaaaaaaaaaaaaaa")
	}
	err := terminal.Restore(0, state)
	if err != nil {
		return
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
	list(ctx, input, sc)

}
