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

func (i *Indexer) Lookup(q *query) []string {
	if q.any {
		return i.fileMap
	}

	trigrams := q.Trigrams()
	files := make([]string, 0)

	var trigramResults [][]int
	for _, trigram := range trigrams {
		fileIndices, ok := i.trigramToFiles[trigram]
		if !ok {
			continue
		}
		trigramResults = append(trigramResults, fileIndices)
	}

	for _, fileIndex := range intersect(trigramResults) {
		files = append(files, i.fileMap[fileIndex])
	}
	return files
}

// intersectPair assumes that a b are both sorted and that there are no duplicates
// [0, 2, 5, 7]
// [3, 4, 5, 6, 7, 8]
// => [5, 7]

// algorithm:

// two pointers, march the lowest,
// if they point to the same value, add and march both
// if you reach the end of either list, return
func intersectPair(A []int, B []int) (res []int) {
	if len(A) == 0 || len(B) == 0 {
		return res
	}
	a, b := 0, 0
	for {
		if A[a] == B[b] {
			res = append(res, A[a])
		}
		if A[a] > B[b] {
			b++
		} else {
			a++
		}
		if a >= len(A) || b >= len(B) {
			return res
		}
	}
}

func intersect(list [][]int) (res []int) {
	for _, item := range list {
		if res == nil {
			res = item
		} else {
			res = intersectPair(res, item)
		}
	}
	return res
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

			if len(t.trigramToFiles[trigram]) == 0 || t.trigramToFiles[trigram][len(t.trigramToFiles[trigram])-1] != fileIndex {
				t.trigramToFiles[trigram] = append(t.trigramToFiles[trigram], fileIndex)
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
