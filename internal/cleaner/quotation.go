package cleaner

import "golang.org/x/text/runes"

var Quotation = runes.Map(func(r rune) rune {
	switch r {
	case '\u0022', // Quotation Mark
		'\u02EE', // Modifier Letter Double Apostrophe
		'\u201C', // Left Double Quotation Mark
		'\u201D', // Right Double Quotation Mark
		'\u201E', // Double Low-9 Quotation Mark
		'\u201F', // Double High-Reversed-9 Quotation Mark
		'\u275D', // Heavy Double Turned Comma Quotation Mark Ornament
		'\u275E', // Heavy Double Comma Quotation Mark Ornament
		'\u2760', // Heavy Low Double Comma Quotation Mark Ornament
		'\u2E42', // Double Low-Reversed-9 Quotation Mark
		'\u301D', // Reversed Double Prime Quotation Mark
		'\u301E', // Double Prime Quotation Mark
		'\u301F', // Low Double Prime Quotation Mark
		'\uFF02', // Fullwidth Quotation Mark
		'\u02BA', // Modifier Letter Double Prime
		'\u2033', // Double Prime
		'\u2036', // Reversed Double Prime
		'\u02DD', // Double Acute Accent
		'\u02F5', // Modifier Letter Middle Double Grave Accent
		'\u02F6': // Modifier Letter Middle Double Acute Accent
		return '"'
	default:
		return r
	}
})
