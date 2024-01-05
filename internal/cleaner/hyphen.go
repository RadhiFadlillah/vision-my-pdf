package cleaner

import "golang.org/x/text/runes"

var Hyphen = runes.Map(func(r rune) rune {
	switch r {
	case '\u002D', // Hyphen-Minus
		'\u2010', // Hyphen
		'\u2011', // Non-Breaking Hyphen
		'\uFE63', // Small Hyphen-Minus
		'\uFF0D', // Fullwidth Hyphen-Minus
		'\u2012', // Figure Dash
		'\u2013': // En Dash
		return '-'
	default:
		return r
	}
})
