package cleaner

import "golang.org/x/text/runes"

var SuperscriptNumber = runes.Map(func(r rune) rune {
	if replacement, exist := superscriptNumbers[r]; exist {
		return replacement
	}
	return r
})

var superscriptNumbers = map[rune]rune{
	'\u2070': '0', // Superscript Zero,
	'\u00B9': '1', // Superscript One,
	'\u00B2': '2', // Superscript Two,
	'\u00B3': '3', // Superscript Three,
	'\u2074': '4', // Superscript Four,
	'\u2075': '5', // Superscript Five,
	'\u2076': '6', // Superscript Six,
	'\u2077': '7', // Superscript Seven,
	'\u2078': '8', // Superscript Eight,
	'\u2079': '9', // Superscript Nine,
}
