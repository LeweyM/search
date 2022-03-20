package integration

import (
	"fmt"
	"testing"
)

// Slow tests which compare search results to grep search results.
func TestIntegration(t *testing.T) {
	type test struct {
		path, regex string
	}

	tests := []test{
		{path: "../data/bible/bible.txt", regex: "(G|g)od"},
		{path: "../data/bible/bible.txt", regex: "heaven.*hell"},
		{path: "../data/bible/bible.txt", regex: "gos?"},
		{path: "../data/bible/bible.txt", regex: "go+d"},
		{path: "../data/bible/bible.txt", regex: "(beast|burden)"},
	}

	for _, t2 := range tests {
		t.Run(fmt.Sprintf("test file search against grep with regex: '%s'", t2.regex), func(t *testing.T) {
			grepResultsMap := getGrepResults(t2.path, t2.regex)
			searchResultsMap := getSearchResults(t2.path, t2.regex)
			for k, v := range grepResultsMap {
				content, hasLine := searchResultsMap[k]
				if !hasLine {
					t.Fatalf("Expected search to find line %d with content %s", k, v)
				}
				if content != v {
					t.Fatalf("Expected line %d to have content '%s', but instead had '%s'", k, v, content)
				}
			}
		})
	}
}

func TestDirectorySearch(t *testing.T) {
	type test struct {
		path, regex string
	}

	tests := []test{
		{path: "../data/bible-in-pages", regex: "(G|g)od"},
		{path: "../data/bible-in-pages", regex: "heaven.*hell"},
		{path: "../data/bible-in-pages", regex: "gos?"},
		{path: "../data/bible-in-pages", regex: "go+d"},
		{path: "../data/bible-in-pages", regex: "(beast|burden)"},
	}

	for _, t2 := range tests {
		t.Run(fmt.Sprintf("test directory search against grep with regex: '%s'", t2.regex), func(t *testing.T) {
			resultsMap := getDirectorySearchResults(t2.path, t2.regex)
			grepResultsMap := getDirectoryGrepResults(t2.regex, t2.path)
			for k, v := range grepResultsMap {
				content, hasLine := resultsMap[k]
				if !hasLine {
					t.Fatalf("Expected search to find %s with content %s", k, v)
				}
				if content != v {
					t.Fatalf("Expected %s to have content '%s', but instead had '%s'", k, v, content)
				}
				delete(grepResultsMap, k)
			}
			if len(grepResultsMap) > 0 {
				t.Fatalf("Search found additional results to grep: %v", grepResultsMap)
			}
		})
	}
}
