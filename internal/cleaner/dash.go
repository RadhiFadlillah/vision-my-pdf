package cleaner

import "golang.org/x/text/runes"

var Dash = runes.Map(func(r rune) rune {
	switch r {
	case '\u2014', // Em Dash
		'\u2E3A', // Two-Em Dash
		'\u2E3B', // Three-Em Dash
		'\uFE58': // Small Em Dash
		return '\u2014'
	default:
		return r
	}
})
