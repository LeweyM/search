package search

import (
	"context"
	"os"
)

type search struct {
	filePath string
	content []byte
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

func (s *search) Search(ctx context.Context, target string, out chan int) {
	line := 0
	i := 0

	for _, ch := range s.content {
		select {
		case <-ctx.Done():
			return
		default:
			if ch == '\n' {
				line++
			}
			if i+len(target) < len(s.content) && string(s.content[i:i+len(target)]) == target {
				out <- line
			}
			i++
		}
	}
}

