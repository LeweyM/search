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
)

type grepResult struct {
	file    string
	line    int
	content string
	match   string
}

func searchWithIndex(index *trigram.Indexer, t struct {
	path, regex string
}) {

	fileCandidates := index.Lookup(trigram.Query(t.regex))
	newSearch := search.NewSearch(t.path)
	var results []search.ResultWithFile
	for _, fileCandidate := range fileCandidates {
		filename := strings.Split(fileCandidate, "/")
		results = append(results, newSearch.SearchFile(context.TODO(), t.regex, filename[3])...)
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
	newSearch := search.NewSearch(path)
	results := newSearch.SearchDirectoryRegex(context.TODO(), regex)
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
	cmd := exec.Command("grep", "-nro", regex, path)
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

func stringsToMap(files []string) map[string]struct{} {
	res := make(map[string]struct{}, len(files))
	for _, file := range files {
		res[file] = struct{}{}
	}
	return res
}
