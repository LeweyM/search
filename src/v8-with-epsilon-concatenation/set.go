package v8EpsilonConcatenation

type Set[t comparable] map[t]struct{}

func NewSet[t comparable](items ...t) Set[t] {
	set := Set[t](make(map[t]struct{}))
	for _, item := range items {
		set.add(item)
	}
	return set
}

func (s *Set[t]) add(item t) {
	(*s)[item] = struct{}{}
}

func (s *Set[t]) remove(item t) {
	delete(*s, item)
}

func (s *Set[t]) has(item t) bool {
	_, ok := (*s)[item]
	return ok
}

func (s *Set[t]) list() []t {
	set := *s
	list := make([]t, len(set))
	i := 0
	for item := range set {
		list[i] = item
		i++
	}
	return list
}

func (s *Set[t]) size() int {
	return len(*s)
}
