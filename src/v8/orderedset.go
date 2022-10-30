package v8

import (
	"sort"
)

// OrderedSet maintains an ordered set of unique items of type <T>
type OrderedSet[T comparable] struct {
	set       map[T]int
	nextIndex int
}

func (o *OrderedSet[T]) add(ts ...T) {
	if o.set == nil {
		o.set = make(map[T]int)
	}

	for _, t := range ts {
		if !o.has(t) {
			o.set[t] = o.nextIndex
			o.nextIndex++
		}
	}
}

func (o *OrderedSet[T]) has(t T) bool {
	_, hasItem := o.set[t]
	return hasItem
}

func (o *OrderedSet[T]) list() []T {
	size := len(o.set)
	list := make([]T, size)

	i := 0
	for t := range o.set {
		list[i] = t
		i++
	}

	sort.Slice(list, func(i, j int) bool {
		return o.getIndex(list[i]) < o.getIndex(list[j])
	})

	return list
}

func (o *OrderedSet[T]) getIndex(t T) int {
	return o.set[t]
}
