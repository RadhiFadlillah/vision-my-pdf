package cli

import (
	"github.com/RadhiFadlillah/vision-my-pdf/internal/cleaner"
	"github.com/urfave/cli/v2"
	"golang.org/x/text/transform"
)

func prepareTextCleaner(c *cli.Context) cleaner.Cleaner {
	var cleaners []transform.Transformer
	addCleaner := func(c ...transform.Transformer) {
		cleaners = append(cleaners, c...)
	}

	// Parse flags for diacritics
	if c.Bool(_noDiacritic) {
		addCleaner(cleaner.Diacritic)
	} else if c.Bool(_noLatinDiacritic) {
		addCleaner(cleaner.LatinDiacritic)
	}

	// Parse flags for number normalization
	if c.Bool(_normNumber) {
		addCleaner(cleaner.SubscriptNumber, cleaner.SuperscriptNumber)
	} else {
		if c.Bool(_normSubNumber) {
			addCleaner(cleaner.SubscriptNumber)
		}
		if c.Bool(_normSuperNumber) {
			addCleaner(cleaner.SubscriptNumber)
		}
	}

	// Parse flags for mark normalization
	if c.Bool(_normMark) {
		addCleaner(cleaner.Hyphen, cleaner.Dash, cleaner.Apostrophe, cleaner.Quotation)
	} else {
		if c.Bool(_normHyphen) {
			addCleaner(cleaner.Hyphen)
		}
		if c.Bool(_normDash) {
			addCleaner(cleaner.Dash)
		}
		if c.Bool(_normApostrophe) {
			addCleaner(cleaner.Apostrophe)
		}
		if c.Bool(_normQuotation) {
			addCleaner(cleaner.Quotation)
		}
	}

	// Return the final cleaner
	return cleaner.New(cleaners...)
}
