package integration

import (
	context2 "context"
	"fmt"
	"os"
	"search/src/search"
	"search/src/trigram"
	"testing"
)

type test struct {
	path, regex string
}

// for a specific case like this, the trigram index can filter out almost all the files where we
// shouldn't search.

// found 1 candidates out of 608 files

// 19,494,433 ns/op - with index
// vs
// 5,953,637,518 ns/op - without index
// == 305x speedup!
func BenchmarkTrigramIndexedDirectorySearchEasy(b *testing.B) {
	compareIndexedAndNotIndexedPerformance(b, test{path: "../data/bible-in-pages", regex: "Shobek"})
}

// for a common word such as this, the trigram index doesn't filter many files so the results are
// less dramatic.

// found 157 candidates out of 608 files

// 2,595,302,207 ns/op - with index
// vs
// 6,172,375,438 ns/op - without index
// == 2-3x speedup
func BenchmarkTrigramIndexedDirectorySearchMedium(b *testing.B) {
	compareIndexedAndNotIndexedPerformance(b, test{path: "../data/bible-in-pages", regex: "god"})
}

// as this case is uses a regex search too small for trigram filtering, we would expect the results to be
// more or less the same

// found 608 candidates out of 608 files

// 7,137,193,041 ns/op - with index
// vs
// 6,320,655,498 ns/op - without index
// == no speedup (actually a little slower)
func BenchmarkTrigramIndexedDirectorySearchHard(b *testing.B) {
	compareIndexedAndNotIndexedPerformance(b, test{path: "../data/bible-in-pages", regex: "ab"})
}

func compareIndexedAndNotIndexedPerformance(b *testing.B, t test) {
	index := trigram.Index(t.path)
	lookup := index.Lookup(trigram.Query(t.regex))

	dir, err := os.ReadDir(t.path)
	if err != nil {
		panic(err)
	}

	b.Logf("Searching in %d files.", len(lookup))
	b.Run(fmt.Sprintf("test with regex: With Index: '%s'", t.regex), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			searchWithIndex(index, t)
		}
	})

	b.Logf("Searching in %d files.", len(dir))
	b.Run(fmt.Sprintf("test with regex: Without Index: '%s'", t.regex), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			search.NewSearch(t.path).SearchDirectoryRegex(context2.TODO(), t.regex)
		}
	})
}
