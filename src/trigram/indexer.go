package trigram

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
)

type Indexer struct {
	fileMap        []string
	trigramToFiles map[string][]int
}

func newIndexer() *Indexer {
	return &Indexer{
		fileMap:        make([]string, 0),
		trigramToFiles: make(map[string][]int),
	}
}

func Index(dirPath string) *Indexer {
	indexer := newIndexer()
	indexer.readDirectory(dirPath)

	fileIndex := 0
	for _, path := range indexer.fileMap {
		open, err := os.Open(path)
		if err != nil {
			panic("cannot read file")
		}

		indexer.index(bufio.NewReader(open), fileIndex)
		fileIndex++
	}

	return indexer
}

func (i *Indexer) Lookup(q *query) []string {
	if q.any {
		return i.fileMap
	}

	filesIndices := q.Lookup(i)

	return i.getFilePathsFromIndices(filesIndices)
}

func (i *Indexer) getFilePathsFromIndices(filesIndices []int) []string {
	var filesPaths []string
	for _, fileIndex := range filesIndices {
		filePath := i.fileMap[fileIndex]
		filesPaths = append(filesPaths, filePath)
	}
	return filesPaths
}

func (i *Indexer) readDirectory(dirPath string) {
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}

	for _, entry := range dir {
		if !entry.IsDir() {
			path := filepath.Join(dirPath, entry.Name())
			i.fileMap = append(i.fileMap, path)
		}
	}
}

func (i *Indexer) index(reader io.RuneReader, fileIndex int) {
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
				i.index(reader, fileIndex)
				return
			}
			trigram = trigram[1:3] + string(c)

			if len(i.trigramToFiles[trigram]) == 0 || i.trigramToFiles[trigram][len(i.trigramToFiles[trigram])-1] != fileIndex {
				i.trigramToFiles[trigram] = append(i.trigramToFiles[trigram], fileIndex)
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
