package v5

import "strings"

// single character

type SingleCharacterPredicate struct {
	character rune
}

func (a SingleCharacterPredicate) test(input rune) bool {
	return input == a.character
}

// disallow list

type DisallowListPredicate struct {
	disallowList string
}

func (d DisallowListPredicate) test(input rune) bool {
	return !strings.ContainsRune(d.disallowList, input)
}
