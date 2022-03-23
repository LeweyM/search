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
	terminal.MakeRaw(0)    // fd 0 is stdin

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

const (
	UP    = "UP"
	DOWN  = "DOWN"
	LEFT  = "LEFT"
	RIGHT = "RIGHT"
)

func (s *screen) update(ctx context.Context, ticker *time.Ticker) {
	var inputL string
	hasChanged := true
	currentState := state{}
	nextState := state{}

	flag1 := false
	flag2 := false

	s.clearScreen()

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
					nextState.input = inputL
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
			// arrows
			if r == 27 {
				flag1 = true
			}
			if flag1 && r == 91 {
				flag2 = true
				continue
			}
			if r >= 65 && r <= 68 && flag1 && flag2 {
				switch r {
				// TODO: This introduces a bug as "UP" input is same as arrow "UP"
				case 65:
					s.output <- UP
				case 66:
					s.output <- DOWN
				case 67:
					s.output <- RIGHT
				case 68:
					s.output <- LEFT
				}
				flag1 = false
				flag2 = false
				continue
			}
			// alphanumeric numbers
			if r >= 32 && r <= 127 {
				inputL = inputL + string(r)
				nextState.input = inputL
				hasChanged = false
				s.output <- inputL
			}
		case <-ticker.C:
			if hasChanged {
				s.refresh(ctx, currentState, nextState)
				currentState = nextState
			}
			hasChanged = false
		case lines := <-s.linesChan:
			nextState.lines = lines
			hasChanged = true
		}
	}
}

func (s *screen) swapOutStates(currentState, nextState state) {
	s.setCursorPosition(1, 0)
	rowChanged := false
	for i, nextStateLine := range nextState.lines {
		if i >= len(currentState.lines) {
			fmt.Fprintf(s.writer, "%s", nextStateLine)
			s.clearRestOfLine()
			fmt.Fprintf(s.writer, "\n\r")
			continue
		}

		width, _, err := terminal.GetSize(0)
		if err != nil {
			panic(err)
		}
		// TODO: This is pretty good for now, but some lines might have a longer length than len(line) because of invisible chars, such as escape ansi codes for color.
		if (len(nextStateLine) / width) > 0 {
			rowChanged = true
		}

		currentStateLine := currentState.lines[i]
		// if state has not changed, do nothing
		if !rowChanged && currentStateLine == nextStateLine {
			fmt.Fprintf(s.writer, "\r\n")
			continue
		}

		fmt.Fprintf(s.writer, "%s", nextStateLine)
		s.clearRestOfLine()
		fmt.Fprintf(s.writer, "\n\r")
	}
	s.clearRestOfScreen()
}

func (s *screen) clearRestOfScreen() {
	fmt.Fprintf(s.writer, "\u001b[0J")
}

func (s *screen) clearRestOfLine() {
	fmt.Fprint(s.writer, "\u001b[0K")
}

func (s *screen) refresh(ctx context.Context, currentState, nextState state) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		s.swapOutStates(currentState, nextState)
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
