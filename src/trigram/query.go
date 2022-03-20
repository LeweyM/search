package trigram

type query struct {
	any       bool
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

func (q *query) compile(exp string) queryable {
	if len(exp) < 3 {
		q.any = true
		return nil // TODO
	}

	if len(exp) == 3 {
		return &trigram{
			val: exp,
		}
	}

	return &or{
		a: &trigram{
			val: exp[0:3],
		},
		b: q.compile(exp[1:]),
	}
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
