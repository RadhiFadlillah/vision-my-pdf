package cleaner

import "golang.org/x/text/runes"

var Apostrophe = runes.Map(func(r rune) rune {
	switch r {
	case '\u0027', // Apostrophe
		'\u02BC', // Modifier Letter Apostrophe
		'\u055A', // Armenian Apostrophe
		'\u07F4', // Nko High Tone Apostrophe
		'\u07F5', // Nko Low Tone Apostrophe
		'\uFF07', // Fullwidth Apostrophe
		'\u2018', // Left Single Quotation Mark
		'\u2019', // Right Single Quotation Mark
		'\u201A', // Single Low-9 Quotation Mark
		'\u201B', // Single High-Reversed-9 Quotation Mark
		'\u275B', // Heavy Single Turned Comma Quotation Mark Ornament
		'\u275C', // Heavy Single Comma Quotation Mark Ornament
		'\u275F', // Heavy Low Single Comma Quotation Mark Ornament
		'\u02B9', // Modifier Letter Prime
		'\u2032', // Prime
		'\u2035', // Reversed Prime
		'\u02BB', // Modifier Letter Turned Comma
		'\uA78B', // Latin Capital Letter Saltillo
		'\uA78C', // Latin Small Letter Saltillo
		'\u0060', // Grave Accent
		'\u00B4', // Acute Accent
		'\u02CA', // Modifier Letter Acute Accent
		'\u02CB', // Modifier Letter Grave Accent
		'\u02CE', // Modifier Letter Low Grave Accent
		'\u02CF', // Modifier Letter Low Acute Accent
		'\u02F4': // Modifier Letter Middle Grave Accent
		return '\''
	default:
		return r
	}
})
