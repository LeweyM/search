package screen

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"sync"
)

type screen struct {
	writer         io.Writer
	lines          []string
	input          string
	InputChan      chan string
	linesChan      chan []string
	output         chan string
	m              sync.Mutex
	appendLineChan chan string
	modifyLineChan chan modification
}

type modification struct {
	n    int
	line string
}

type Screen interface {
	AddLine(line string)
	SetLines(lines []string)
	Run(ctx context.Context)
	SetLine(i int, line string)
}

func NewScreen(writer io.Writer, out chan string) *screen {
	return &screen{
		writer:         writer,
		lines:          []string{},
		input:          "",
		InputChan:      make(chan string),
		modifyLineChan: make(chan modification),
		appendLineChan: make(chan string, 100),
		linesChan:      make(chan []string),
		output:         out,
		m:              sync.Mutex{},
	}
}

func (s *screen) AddLine(line string) {
	s.appendLineChan <- line
}

func (s *screen) SetLines(lines []string) {
	s.linesChan <- lines
}

func (s *screen) SetLine(lineNum int, line string) {
	s.modifyLineChan <- modification{n: lineNum, line: line}
}

func (s *screen) Run(ctx context.Context) {
	go s.readInput(ctx, bufio.NewReader(os.Stdin))
	go s.updateStateAndScreen(ctx)
}

func (s *screen) updateStateAndScreen(ctx context.Context) {
	for {
		select {
		case lineMod := <-s.modifyLineChan:
			s.m.Lock()
			if lineMod.n >= len(s.lines) {
				i := len(s.lines)
				for i < lineMod.n {
					s.lines = append(s.lines, "")
					i++
				}
				s.lines = append(s.lines, lineMod.line)
			} else {
				s.lines[lineMod.n] = lineMod.line
			}
			s.m.Unlock()
			if s.refresh(ctx) {
				return
			}
		case line := <-s.appendLineChan:
			s.m.Lock()
			s.lines = append(s.lines, line)
			s.m.Unlock()
			if s.refresh(ctx) {
				return
			}
		case lines := <-s.linesChan:
			s.m.Lock()
			s.lines = lines
			s.m.Unlock()
			if s.refresh(ctx) {
				return
			}
		case input := <-s.InputChan:
			s.m.Lock()
			s.input = input
			s.m.Unlock()
			s.output <- s.input
			if s.refresh(ctx) {
				return
			}
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
