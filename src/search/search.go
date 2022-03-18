package search

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
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
	Match       Match
}

type ResultWithFile struct {
	Result
	File string
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

func (s *search) SearchRegex(ctx context.Context, regex string, out chan Result) {
	state := finite_state_machine.Compile(regex)
	runner := finite_state_machine.NewRunner(state)
	resultChan := make(chan finite_state_machine.Result)

	go finite_state_machine.FindAllAsync(ctx, runner, string(s.content), resultChan)

	count := 0
	for result := range resultChan {
		out <- Result{
			LineNumber:  result.Line,
			LineContent: s.lines[result.Line-1],
			Match:       Match{Start: result.Start, End: result.End},
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

func (s *search) SearchDirectoryRegex(regex string) []ResultWithFile {
	fileChan := make(chan string)
	resultChan := make(chan ResultWithFile, 1000) //TODO: support more results
	defer close(fileChan)

	state := finite_state_machine.Compile(regex)

	// start workers
	for i := 0; i < runtime.NumCPU(); i++ {
		go s.worker(state, fileChan, resultChan)
	}

	// read directory
	dir, err := os.ReadDir(s.filePath)
	if err != nil {
		panic(err)
	}

	// load directories into workers
	for _, entry := range dir {
		if !entry.IsDir() {
			fileChan <- entry.Name()
		}
	}
	close(resultChan) //TODO: If files not finished loading this can lead to send on closed channel

	// compile results
	var res []ResultWithFile
	for result := range resultChan {
		result.Query = regex
		res = append(res, result)
	}

	return res
}

func (s *search) worker(fsm *finite_state_machine.State, in chan string, out chan ResultWithFile) {
	for fileName := range in {
		path := filepath.Join(s.filePath, fileName)

		file, err := os.ReadFile(path)
		if err != nil {
			panic("cannot read file")
		}

		var results []finite_state_machine.Result
		runner := finite_state_machine.NewRunner(fsm)
		results = finite_state_machine.FindAllWithLines(runner, string(file))

		lines := strings.Split(string(file), "\n")
		for _, result := range results {
			out <- ResultWithFile{
				Result: Result{
					LineNumber:  result.Line,
					LineContent: lines[result.Line-1],
					Match:       Match{Start: result.Start, End: result.End},
					Finished:    false,
				},
				File: fileName,
			}
		}
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
