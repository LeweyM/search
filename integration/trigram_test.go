package integration

import (
	"search/src/trigram"
	"strings"
	"testing"
)

func TestTrigramIndexer(t *testing.T) {
	path := "../data/bible-in-pages"
	index := trigram.Index(path)

	testFileCandidatesAreTheSameAsGrep(t, index, path, "Shobek")
}

func testFileCandidatesAreTheSameAsGrep(t *testing.T, index *trigram.Indexer, path, query string) {
	files := index.Lookup(trigram.Query(query))
	fileMap := stringsToMap(files)
	grepResult := getDirectoryGrepResults(query, path)
	for k := range grepResult {
		fileParts := strings.Split(k, "-")
		file := path + "/" + strings.Join([]string{fileParts[0], fileParts[1], fileParts[2]}, "-")
		_, hasFile := fileMap[file]
		if !hasFile {
			t.Fatalf("Expected search to find [%s]", file)
		}
		delete(fileMap, file)
	}
	if len(fileMap) > 0 {
		t.Fatalf("Search found additional results to grep: %v", fileMap)
	}
}
