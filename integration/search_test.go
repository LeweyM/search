package integration

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"search/src/search"
	"search/src/trigram"
	"strconv"
	"strings"
	"testing"
)

type grepResult struct {
	file    string
	line    int
	content string
	match   string
}

// Slow tests which compare search results to grep search results.
func TestIntegration(t *testing.T) {
	type test struct {
		searchPath, grepPath, regex string
	}

	tests := []test{
		{searchPath: "../data/bible/bible.txt", grepPath: "/data/bible/", regex: "(G|g)od"},
		{searchPath: "../data/bible/bible.txt", grepPath: "/data/bible/", regex: "heaven.*hell"},
		{searchPath: "../data/bible/bible.txt", grepPath: "/data/bible/", regex: "gos?"},
		{searchPath: "../data/bible/bible.txt", grepPath: "/data/bible/", regex: "go+d"},
		{searchPath: "../data/bible/bible.txt", grepPath: "/data/bible/", regex: "(beast|burden)"},
	}

	for _, t2 := range tests {
		t.Run(fmt.Sprintf("test against grep with regex: '%s'", t2.regex), func(t *testing.T) {
			grepResultsMap := getGrepResults(t2.grepPath, t2.regex)
			searchResultsMap := getSearchResults(t2.searchPath, t2.regex)
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
		{path: "/data/bible-in-pages", regex: "(G|g)od"},
		{path: "/data/bible-in-pages", regex: "heaven.*hell"},
		{path: "/data/bible-in-pages", regex: "gos?"},
		{path: "/data/bible-in-pages", regex: "go+d"},
		{path: "/data/bible-in-pages", regex: "(beast|burden)"},
	}

	for _, t2 := range tests {
		t.Run(fmt.Sprintf("test against grep with regex: '%s'", t2.regex), func(t *testing.T) {
			trigram.Index(".." + t2.path)
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

func getDirectoryGrepResults(regex string, path string) map[string]string {
	sanitizedGrep := sanitize(regex)
	grepResults := grep(sanitizedGrep, path)
	grepResultsMap := make(map[string]string)
	for _, res := range grepResults {
		grepResultsMap[res.file+"-"+strconv.Itoa(res.line)] = res.content
	}
	return grepResultsMap
}

func getDirectorySearchResults(path string, regex string) map[string]string {
	newSearch := search.NewSearch(".." + path)
	results := newSearch.SearchDirectoryRegex(regex)
	resultsMap := make(map[string]string)
	for _, result := range results {
		if result.Finished {
			break
		}
		key := result.File + "-" + strconv.Itoa(result.LineNumber)
		_, has := resultsMap[key]
		// only take the first result for a line
		if !has {
			resultsMap[key] = result.LineContent[result.Match.Start : result.Match.End+1]
		}
	}
	return resultsMap
}

func getGrepResults(path string, regex string) map[int]string {
	sanitizedGrep := sanitize(regex)
	grepResults := grep(sanitizedGrep, path)
	grepResultsMap := make(map[int]string)
	for _, res := range grepResults {
		grepResultsMap[res.line] = res.content
	}
	return grepResultsMap
}

func getSearchResults(filePath string, regex string) map[int]string {
	newSearch := search.NewSearch(filePath)
	newSearch.LoadInMemory()
	newSearch.LoadLinesInMemory()

	out := make(chan search.Result)
	go newSearch.SearchRegex(context.Background(), regex, out)

	searchResultsMap := make(map[int]string)
	for result := range out {
		if result.Finished {
			break
		}
		_, has := searchResultsMap[result.LineNumber]
		// only take the first result for a line
		if !has {
			searchResultsMap[result.LineNumber] = result.LineContent[result.Match.Start : result.Match.End+1]
		}
	}
	return searchResultsMap
}

func sanitize(regex string) string {
	var res string
	escapableCharacters := map[rune]bool{
		'(': true,
		')': true,
		'|': true,
		'+': true,
		'?': true,
	}
	for _, char := range regex {
		if escapableCharacters[char] {
			res += "\\" + string(char)
		} else {
			res += string(char)
		}
	}
	return res
}

func grep(regex, path string) []grepResult {
	out := bytes.Buffer{}
	cmd := exec.Command("grep", "-nro", regex, ".."+path)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	s := out.String()
	lines := strings.Split(s, "\n")
	results := make([]grepResult, 0, len(lines))
	for _, line := range lines {
		after := strings.SplitAfter(line, fmt.Sprintf(path+"/"))
		if len(after) < 2 {
			continue
		}
		res := after[1]
		data := strings.SplitN(res, ":", 3)
		page := data[0]
		line, err := strconv.Atoi(data[1])
		if err != nil {
			panic(err)
		}
		content := data[2]
		results = append(results, grepResult{
			file:    page,
			line:    line,
			content: content,
		})
	}
	return results
}
