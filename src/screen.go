package screen

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
)

type screen struct {
	writer    io.Writer
	lines     []string
	input     string
	inputChan chan string
	linesChan chan []string
}

func NewScreen(writer io.Writer) *screen {
	return &screen{
		writer:    writer,
		lines:     []string{},
		input:     "",
		inputChan: make(chan string),
		linesChan: make(chan []string),
	}
}

func (s *screen) SetLines(lines []string) {
	s.linesChan <- lines
}

func (s *screen) Run(ctx context.Context) {
	go s.readInput(ctx, bufio.NewReader(os.Stdin))
	go s.run(ctx)
}

func (s *screen) run(ctx context.Context) {
	go func() {
		for lines := range s.linesChan {
			s.lines = lines
			if s.refresh(ctx) {
				return
			}
		}
	}()

	go func() {
		for input := range s.inputChan {
			s.input = input
			if s.refresh(ctx) {
				return
			}
		}
	}()
}

func (s *screen) refresh(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		s.setCursorPosition(0, 0)
		s.clearScreen()
		s.printScreen()
	}
	return false
}

func (s *screen) printScreen() {
	fmt.Fprint(s.writer, "screen: ")
	fmt.Fprint(s.writer, "\r\n")
	fmt.Fprint(s.writer, s.input)
	fmt.Fprint(s.writer, "\r\n")

	for _, line := range s.lines {
		fmt.Fprintf(s.writer, "\r\n%s", line)
	}
}

func (s *screen) readInput(ctx context.Context, in *bufio.Reader) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			r, _, err := in.ReadRune()
			if err != nil {
				panic(err)
			}
			if r == 'q' {
				panic(fmt.Errorf("user exit program"))
			}
			// backspace
			if r == 127 {
				if len(s.input) > 0 {
					s.inputChan <- s.input[0 : len(s.input)-1]
				}
				continue
			}
			// enter
			if r == 13 {
				s.inputChan <- ""
				continue
			}
			if string(r) != "" {
				s.inputChan <- s.input + string(r)
				continue
			}
		}
	}
}

func (s *screen) setCursorPosition(y, x int) {
	fmt.Fprint(s.writer, fmt.Sprintf("\033[%d;%dH", y, x))
}

func (s *screen) clearScreen() {
	fmt.Fprintf(s.writer, "\033[H\033[2J")
}
