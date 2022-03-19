package integration

import (
	"search/src/trigram"
	"strings"
	"testing"
)

func TestTrigramIndexer(t *testing.T) {
	// search for word 'Shobek'
	path := "../data/bible-in-pages"

	index := trigram.Index(path)
	files := index.Lookup(trigram.Query("Shobek"))
	fileMap := stringsToMap(files)

	grepResult := getDirectoryGrepResults("Shobek", "/data/bible-in-pages")

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

	print(files)
}
