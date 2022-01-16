package screen

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type state struct {
	lines []string
	input string
}

func (s state) Equals(other state) bool {
	return s.input == other.input && s.linesAreEqual(other)
}

func (a state) linesAreEqual(b state) bool {
	if len(a.lines) != len(b.lines) {
		return false
	}
	for i, v := range a.lines {
		if v != b.lines[i] {
			return false
		}
	}
	return true
}

type screen struct {
	writer         io.Writer
	lines          []string
	input          string
	InputChan      chan string
	linesChan      chan []string
	output         chan string
	m              sync.Mutex
}

type Screen interface {
	SetLines(lines []string)
	Run(ctx context.Context)
}

func NewScreen(writer io.Writer, out chan string) *screen {
	return &screen{
		writer:         writer,
		lines:          []string{},
		input:          "",
		InputChan:      make(chan string),
		linesChan:      make(chan []string),
		output:         out,
		m:              sync.Mutex{},
	}
}

func (s *screen) SetLines(lines []string) {
	s.linesChan <- lines
}

func (s *screen) Run(ctx context.Context) {
	go s.readInput(ctx, bufio.NewReader(os.Stdin))
	go s.updateStateAndScreen(ctx)
	ticker := time.NewTicker(50 * time.Millisecond)
	var oldState state
	go func() {
		for range ticker.C {
			nextState := state{
				lines: s.lines,
				input: s.input,
			}
			if !nextState.Equals(oldState) {
				s.refresh(ctx)
			}
			oldState = nextState
		}
	}()
}

func (s *screen) updateStateAndScreen(ctx context.Context) {
	for {
		select {
		case lines := <-s.linesChan:
			s.m.Lock()
			s.lines = lines
			s.m.Unlock()
		case input := <-s.InputChan:
			s.m.Lock()
			s.input = input
			s.m.Unlock()
			s.output <- s.input
		}
	}
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
					s.InputChan <- s.input[0 : len(s.input)-1]
				}
				continue
			}
			// enter
			if r == 13 {
				s.InputChan <- ""
				continue
			}
			if string(r) != "" {
				s.InputChan <- s.input + string(r)
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
