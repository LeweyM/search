package trigram

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
)

type Indexer struct {
	fileMap []string
	m       map[string][]int
}

func newIndexer() *Indexer {
	return &Indexer{
		fileMap: make([]string, 0),
		m:       make(map[string][]int),
	}
}

func Index(dirPath string) *Indexer {
	// read directory
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}

	indexer := newIndexer()

	fileIndex := 0
	for _, entry := range dir {
		if !entry.IsDir() {
			path := filepath.Join(dirPath, entry.Name())

			open, err := os.Open(path)
			if err != nil {
				panic("cannot read file")
			}

			indexer.fileMap = append(indexer.fileMap, path)
			indexer.index(bufio.NewReader(open), fileIndex)
			fileIndex++
		}
	}

	return indexer
}

func (t *Indexer) Lookup(s string) []string {
	if len(s) != 3 {
		panic("can only lookup a trigram")
	}

	fileIndices, ok := t.m[s]
	if !ok {
		return []string{}
	}

	files := make([]string, len(fileIndices))
	for i, fileIndex := range fileIndices {
		files[i] = t.fileMap[fileIndex]
	}
	return files
}

func (t *Indexer) index(reader io.RuneReader, fileIndex int) {
	trigram, err := readThree(reader)
	if err == io.EOF {
		return
	} else if err != nil {
		panic(err)
	}

	for {
		if c, _, err := reader.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		} else {
			if c == '\n' {
				t.index(reader, fileIndex)
				return
			}
			trigram = trigram[1:3] + string(c)

			if len(t.m[trigram]) == 0 || t.m[trigram][len(t.m[trigram])-1] != fileIndex {
				t.m[trigram] = append(t.m[trigram], fileIndex)
			}
		}
	}
}

func readThree(reader io.RuneReader) (string, error) {
	res := make([]rune, 3)
	for i := 0; i < 3; i++ {
		c, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return "", err
			} else {
				panic(err)
			}
		}
		res[i] = c
	}
	return string(res), nil
}

func hasIndexed() {

}

func lookup() {

}
