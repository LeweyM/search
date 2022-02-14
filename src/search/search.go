package search

import (
	"context"
	"os"
	finite_state_machine "search/src/finite-state-machine"
	"strings"
)

const CHAR_PADDING int = 50

type search struct {
	filePath string
	content  []byte
	lines    []string
}

type Match struct {
	Start, End int
}

type Result struct {
	LineNumber  int
	LineContent string
	Count       int
	Query       string
	Finished    bool
	Result      Match
}

func NewSearch(filePath string) *search {
	return &search{filePath: filePath}
}

func (s *search) LoadInMemory() {
	file, err := os.ReadFile(s.filePath)
	if err != nil {
		panic("cannot read file")
	}
	s.content = file
}

func (s *search) LoadLinesInMemory() {
	file, err := os.ReadFile(s.filePath)
	if err != nil {
		panic("cannot read file")
	}
	var sb strings.Builder
	for _, ch := range file {
		if ch == '\n' {
			s2 := sb.String()
			s.lines = append(s.lines, s2)
			sb.Reset()
		} else {
			sb.WriteByte(ch)
		}
	}
	s.lines = append(s.lines, sb.String())
	return
}

func (s *search) Count(ctx context.Context, target string, out chan int) {
	count := 0
	i := 0

	for range s.content {
		select {
		case <-ctx.Done():
			return
		default:
			if i+len(target) < len(s.content) && string(s.content[i:i+len(target)]) == target {
				count++
				out <- count
			}
			i++
		}
	}
	out <- count
}

func (s *search) Search(ctx context.Context, target string, out chan Result) {
	line := 0
	i := 0
	count := 0

	for _, ch := range s.content {
		select {
		case <-ctx.Done():
			return
		default:
			if ch == '\n' {
				line++
			}
			if i+len(target) <= len(s.content) && string(s.content[i:i+len(target)]) == target {
				//time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
				out <- Result{
					Finished:    false,
					Query:       target,
					Count:       count,
					LineNumber:  line + 1,
					LineContent: s.lines[line],
				}
				count++
			}
			i++
		}
	}
	out <- Result{
		Finished: true,
		Query:    target,
	}
}

func (s *search) SearchRegex(ctx context.Context, regex string, out chan Result) {
	state := finite_state_machine.Compile(regex)
	runner := finite_state_machine.NewRunner(state)
	resultChan := make(chan finite_state_machine.Result, 10)

	go finite_state_machine.FindAllAsync(ctx, runner, string(s.content), resultChan)

	count := 0
	for result := range resultChan {
		start := s.sampleStart(result)
		end := s.sampleEnd(result)
		out <- Result{
			LineNumber:  result.Line,
			LineContent: string(s.content)[start:end],
			Result:      Match{Start: result.Start - start, End: (result.Start - start) + (result.End - result.Start)},
			Count:       count,
			Query:       regex,
			Finished:    false,
		}
		count++
	}

	out <- Result{
		Finished: true,
	}
}

func (s *search) sampleEnd(result finite_state_machine.Result) int {
	end := result.End + 1 + CHAR_PADDING
	if end >= len(string(s.content)) {
		return result.End + 1
	} else {
		return end
	}
}

func (s *search) sampleStart(result finite_state_machine.Result) int {
	start := result.Start - CHAR_PADDING
	if start < 0 {
		return 0
	} else {
		return start
	}
}
