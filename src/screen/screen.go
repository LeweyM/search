package screen

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"strings"
	"text/template"
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
	writer    io.Writer
	InputChan chan rune
	linesChan chan []string
	output    chan string
	template  *template.Template
}

func (s *screen) SetState(state interface{}) {
	var bf bytes.Buffer
	err := s.template.Execute(&bf, state)
	if err != nil {
		panic(err)
	}
	s2 := bf.String()
	split := strings.Split(s2, "\n")
	s.linesChan <- split
}

func (s *screen) SetTemplate(templateString string) {
	parsedTemplate, err := template.New("screenTemplate").Parse(templateString)
	if err != nil {
		panic(err)
	}
	s.template = parsedTemplate
}

func NewScreen(writer io.Writer, out chan string) *screen {
	return &screen{
		writer:    writer,
		InputChan: make(chan rune, 100),
		linesChan: make(chan []string, 100),
		output:    out,
	}
}

func (s *screen) SetLines(lines []string) {
	s.linesChan <- lines
}

func (s *screen) Run(ctx context.Context, inputStream io.Reader, exit func()) {
	fmt.Print("\033[?25l") // hide cursor
	terminal.MakeRaw(0) // fd 0 is stdin

	go s.readInput(ctx, time.NewTicker(10*time.Millisecond), bufio.NewReader(inputStream), exit)
	go s.update(ctx, time.NewTicker(50*time.Millisecond))
}

func (s *screen) readInput(ctx context.Context, ticker *time.Ticker, in io.RuneReader, exit func()) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r, _, err := in.ReadRune()
			if err != nil {
				panic(err)
			} // exit program on key "Q"
			if r == 'Q' {
				exit()
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
			// backspace
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
			//// non-alphanumeric numbers
			//if r < 65 || r > 122 {
			//	continue
			//}
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
