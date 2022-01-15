package screen

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

type screen struct {
	writer io.Writer
	lines []string
	input string
}

func NewScreen(writer io.Writer, lines []string) *screen {
	return &screen{writer: writer, lines: lines}
}

func (s *screen) SetLines(lines []string) {
	s.lines = lines
}

func (s *screen) Run(ctx context.Context) {
	in := bufio.NewReader(os.Stdin)
	go func() {
		for {
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
					s.input = s.input[0:len(s.input)-1]
				}
				continue
			}
			// enter
			if r == 13 {
				s.input = ""
				continue
			}
			if string(r) != "" {
				s.input += string(r)
				continue
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			s.setCursorPosition(0, 0)
			s.clearScreen()
			fmt.Fprint(s.writer, "screen: ")
			fmt.Fprint(s.writer, "\r\n")
			fmt.Fprint(s.writer, s.input)
			fmt.Fprint(s.writer, "\r\n")

			for _, line := range s.lines {
				fmt.Fprintf(s.writer, "\r\n%s", line)
			}

			time.Sleep(time.Millisecond * 50)
		}
	}
}

func (s *screen) setCursorPosition(y, x int) {
	fmt.Fprint(s.writer, fmt.Sprintf("\033[%d;%dH", y, x))
}

func (s *screen) clearScreen() {
	fmt.Fprintf(s.writer, "\033[H\033[2J")
}

func (s *screen) readInput(in *bufio.Reader) (error, rune) {
	r, _, err := in.ReadRune()
	if err != nil {
		return err, 'a'
	}
	if r == 'q' {
		return fmt.Errorf("user exit program"), 'a'
	}

	return nil, r
}
