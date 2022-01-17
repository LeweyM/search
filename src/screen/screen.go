package screen

import (
	"bufio"
	"context"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"time"
)

type state struct {
	lines []string
	input string
}

func (a state) Equals(other state) bool {
	return a.input == other.input && a.linesAreEqual(other)
}

func (a state) linesAreEqual(other state) bool {
	if len(a.lines) != len(other.lines) {
		return false
	}
	for i, v := range a.lines {
		if v != other.lines[i] {
			return false
		}
	}
	return true
}

type screen struct {
	writer         io.Writer
	InputChan      chan rune
	linesChan      chan []string
	output         chan string
}

type Screen interface {
	SetLines(lines []string)
	Run(ctx context.Context)
}

func NewScreen(writer io.Writer, out chan string) *screen {
	return &screen{
		writer:         writer,
		InputChan:      make(chan rune),
		linesChan:      make(chan []string),
		output:         out,
	}
}

func (s *screen) SetLines(lines []string) {
	s.linesChan <- lines
}

func (s *screen) Run(ctx context.Context) {
	fmt.Print("\033[?25l") // hide cursor
	terminal.MakeRaw(0)    // fd 0 is stdin

	go s.readInput(ctx, bufio.NewReader(os.Stdin))
	go s.update(ctx, time.NewTicker(100*time.Millisecond))
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
			// exit program on key "Q"
			if r == 'Q' {
				os.Exit(0)
			}
			s.InputChan <- r
		}
	}
}

func (s *screen) update(ctx context.Context, ticker *time.Ticker) {
	var linesL []string
	var inputL string
	hasChanged := true

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case r := <-s.InputChan:
			if r == 127 {
				if len(inputL) > 0 {
					inputL = inputL[0 : len(inputL)-1]
					s.output <- inputL
				}
				continue
			}
			// enter
			if r == 13 {
				inputL = ""
				s.output <- inputL
				continue
			}
			inputL = inputL + string(r)
			hasChanged = false
			s.output <- inputL
		case <-ticker.C:
			if hasChanged {
				s.refresh(ctx, state{
					lines: linesL,
					input: inputL,
				})
			}
			hasChanged = false
		case lines := <-s.linesChan:
			linesL = lines
			hasChanged = true
		}
	}
}

func (s *screen) refresh(ctx context.Context, st state) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		s.setCursorPosition(0, 0)
		s.clearScreen()
		s.printScreen(st)
	}
	return false
}

func (s *screen) printScreen(st state) {
	fmt.Fprint(s.writer, "screen: ")
	fmt.Fprint(s.writer, "\r\n")
	fmt.Fprint(s.writer, st.input)
	fmt.Fprint(s.writer, "\r\n")

	for _, line := range st.lines {
		fmt.Fprintf(s.writer, "\r\n%s", line)
	}
}

func (s *screen) setCursorPosition(y, x int) {
	fmt.Fprint(s.writer, fmt.Sprintf("\033[%d;%dH", y, x))
}

func (s *screen) clearScreen() {
	fmt.Fprintf(s.writer, "\033[H\033[2J")
}
