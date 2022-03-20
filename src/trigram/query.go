package trigram

import "strings"

type query struct {
	trigrams  []queryable
	rootQuery queryable
	index     *Indexer
}

func Query(exp string) *query {
	q := query{}
	q.rootQuery = q.compile(exp)

	return &q
}

func (q *query) Lookup(index *Indexer) []int {
	return q.rootQuery.lookup(index)
}

// compile has a recursive definition
//
// for concatenations, ors are used...
//
// compile(abcde) == compile(abcd) OR compile(cde)
// compile(cde) == trigram(cde)
// compile(abcd) == compile(abc) OR compile(bcd)
// compile(abc) == trigram(abc)
// compile(bcd) == trigram(bcd)
// so
// compile(abcde) == trigram(abc) OR trigram(bcd) OR trigram(cde)
//
// for pipes, ands are used...
//
// compile(abc|defg) == compile(abc) AND compile(defg)
// compile(abc|defg) == compile(abc) AND compile(def) OR compile(efg)
// compile(abc|defg) == trigram(abc) AND trigram(def) OR trigram(efg)
func (q *query) compile(exp string) queryable {
	if len(exp) < 3 {
		return &any{}
	}

	if len(exp) == 3 {
		return &trigram{val: exp}
	}

	containsPipe := strings.Contains(exp, "|")
	// no pipe
	if containsPipe {
		pipeSplit := strings.SplitN(exp, "|", 2)
		return &and{
			a: q.compile(pipeSplit[0]),
			b: q.compile(pipeSplit[1]),
		}
	} else {
		return &or{
			a: q.compile(exp[0:3]),
			b: q.compile(exp[1:]),
		}
	}
}

type queryable interface {
	lookup(indexer *Indexer) []int
}

type or struct {
	a, b queryable
}

func (o *or) lookup(indexer *Indexer) []int {
	return intersectPair(o.a.lookup(indexer), o.b.lookup(indexer))
}

type trigram struct {
	val string
}

func (t *trigram) lookup(indexer *Indexer) []int {
	return indexer.trigramToFiles[t.val]
}

type and struct {
	a, b queryable
}

func (a *and) lookup(indexer *Indexer) []int {
	return unionPair(a.a.lookup(indexer), a.b.lookup(indexer))
}

type any struct{}

func (a *any) lookup(indexer *Indexer) []int {
	res := make([]int, len(indexer.fileMap))
	for i := range res {
		res[i] = i
	}
	return res
}
