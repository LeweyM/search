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
		{path: "../data/bible-in-pages", query: "sadness"},
		{path: "../data/bible-in-pages", query: "beast|burden|bad"},
		{path: "../data/bible-in-pages", query: "z"},
	}

	for _, t2 := range tests {
		t.Run(fmt.Sprintf("[Trigram Indexer] Test that file candiates are the same as grep for query (%s)", t2.query), func(t *testing.T) {
			testFileCandidatesAreTheSameAsGrep(t, t2.path, t2.query)
		})
	}
}

func testFileCandidatesAreTheSameAsGrep(t *testing.T, path, query string) {
	trigramResultSet := stringsToMap(trigram.Index(path).Lookup(trigram.Query(query)))
	grepPagesSet := toPagesMap(getDirectoryGrepResults(query, path), path)
	for grepFile := range grepPagesSet {
		_, hasFile := trigramResultSet[grepFile]
		if !hasFile {
			// The trigram index cannot produce false negatives. If it does not provide a candidate
			// which would contain a positive match, the index is incorrect.
			t.Fatalf("Trigram lookup failed to provide candidate [%s]", grepFile)
		}
		delete(trigramResultSet, grepFile)
	}
	if len(trigramResultSet) > 0 {
		// The trigram index can produce false positives, so it's not incorrect for it to produce candidates
		// which do not match the search.
		t.Logf("Trigram lookup found [%d] additional candiates which did not contain matches: %v", len(trigramResultSet), trigramResultSet)
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
