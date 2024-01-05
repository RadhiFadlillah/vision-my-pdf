package cleaner

import "golang.org/x/text/runes"

var SubscriptNumber = runes.Map(func(r rune) rune {
	if replacement, exist := subscriptNumbers[r]; exist {
		return replacement
	}
	return r
})

var subscriptNumbers = map[rune]rune{
	'\u2080': '0', // Subscript Zero,
	'\u2081': '1', // Subscript One,
	'\u2082': '2', // Subscript Two,
	'\u2083': '3', // Subscript Three,
	'\u2084': '4', // Subscript Four,
	'\u2085': '5', // Subscript Five,
	'\u2086': '6', // Subscript Six,
	'\u2087': '7', // Subscript Seven,
	'\u2088': '8', // Subscript Eight,
	'\u2089': '9', // Subscript Nine,
}
