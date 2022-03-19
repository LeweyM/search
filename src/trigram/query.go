package trigram

type query struct {
	any      bool
	trigrams []string
}

func Query(exp string) *query {
	q := query{}
	q.compile(exp)

	return &q
}

func (q *query) Trigrams() []string {
	return q.trigrams
}

func (q *query) compile(exp string) {
	if len(exp) < 3 {
		q.any = true
		return
	}

	if len(exp) == 3 {
		q.trigrams = []string{exp}
		return
	}

	var trigrams []string
	for i := 0; i < len(exp)-2; i++ {
		trigrams = append(trigrams, exp[i:i+3])
	}

	q.trigrams = trigrams
	return
}
