package integration

import (
	"fmt"
	"search/src/trigram"
	"strings"
	"testing"
)

func TestTrigramIndexer(t *testing.T) {
	type test struct {
		path, query string
	}

	tests := []test{
		{path: "../data/bible-in-pages", query: "Shobek"},
		{path: "../data/bible-in-pages", query: "god"},
	}

	for _, t2 := range tests {
		t.Run(fmt.Sprintf("[Trigram Indexer] Test that file candiates are the same as grep for query (%s)", t2.query), func(t *testing.T) {
			testFileCandidatesAreTheSameAsGrep(t, t2.path, t2.query)
		})
	}
}

func testFileCandidatesAreTheSameAsGrep(t *testing.T, path, query string) {
	index := trigram.Index(path)
	files := index.Lookup(trigram.Query(query))
	trigramResults := stringsToMap(files)
	grepResult := getDirectoryGrepResults(query, path)
	grepPagesMap := toPagesMap(grepResult, path)
	for grepFile := range grepPagesMap {
		_, hasFile := trigramResults[grepFile]
		if !hasFile {
			t.Fatalf("Expected search to find [%s]", grepFile)
		}
		delete(trigramResults, grepFile)
	}
	if len(trigramResults) > 0 {
		t.Fatalf("Search found additional results to grep: %v", trigramResults)
	}
}

func toPagesMap(result map[string]string, path string) map[string]struct{} {
	res := make(map[string]struct{}, len(result))
	for fileAndLine, _ := range result {
		res[buildFileName(fileAndLine, path)] = struct{}{}
	}
	return res
}

func buildFileName(k string, path string) string {
	fileParts := strings.Split(k, "-")
	file := path + "/" + strings.Join([]string{fileParts[0], fileParts[1], fileParts[2]}, "-")
	return file
}
