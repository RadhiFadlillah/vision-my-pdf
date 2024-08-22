package cli

import (
	"runtime"

	"github.com/urfave/cli/v2"
)

const (
	// Flag names for app worker and output
	_force       = "force"
	_worker      = "worker"
	_genDebug    = "gen-debug"
	_montageSize = "montage"

	// Flag names for OCR parameters
	_sortVertical = "sort-vertical"
	_mergeNewLine = "merge-newline"

	// Flag names for text cleaner
	_noDiacritic      = "no-diacritic"
	_noLatinDiacritic = "no-latin-diacritic"
	_normSubNumber    = "norm-sub-number"
	_normSuperNumber  = "norm-super-number"
	_normHyphen       = "norm-hyphen"
	_normDash         = "norm-dash"
	_normApostrophe   = "norm-apostrophe"
	_normQuotation    = "norm-quotation"

	// Flag names for shortcut for several other flags
	_normNumber = "norm-number"
	_normMark   = "norm-mark"
)

var appFlags = []cli.Flag{
	// Flags for app worker and output
	&cli.BoolFlag{
		Name:    _force,
		Aliases: []string{"f"},
		Usage:   "overwrite the existing OCR result",
	},
	&cli.Int64Flag{
		Name:    _worker,
		Aliases: []string{"j"},
		Usage:   "number of concurrent worker(s)",
		Value:   int64(runtime.GOMAXPROCS(0)),
	},
	&cli.BoolFlag{
		Name:    _genDebug,
		Aliases: []string{"gd"},
		Usage:   "generate debug image",
	},
	&cli.IntFlag{
		Name:    _montageSize,
		Aliases: []string{"m"},
		Usage:   "montage image size (must be between 1 and 5)",
		Value:   1,
	},

	// Flags for OCR parameters
	&cli.BoolFlag{
		Name:    _sortVertical,
		Aliases: []string{"sv"},
		Usage:   "sort paragraphs vertically, not work for multi column",
	},
	&cli.BoolFlag{
		Name:    _mergeNewLine,
		Aliases: []string{"mn"},
		Usage:   "merge newlines in a paragraph for text output",
	},

	// Flags for text cleaner
	&cli.BoolFlag{
		Name:    _noDiacritic,
		Aliases: []string{"nd"},
		Usage:   "remove all diacritics",
	},
	&cli.BoolFlag{
		Name:    _noLatinDiacritic,
		Aliases: []string{"nld"},
		Usage:   "remove only diacritics for Latin-script",
	},
	&cli.BoolFlag{
		Name:    _normSubNumber,
		Aliases: []string{"nsub"},
		Usage:   "replace subscript numbers with the normal one",
	},
	&cli.BoolFlag{
		Name:    _normSuperNumber,
		Aliases: []string{"nsuper"},
		Usage:   "replace superscript numbers with the normal one",
	},
	&cli.BoolFlag{
		Name:    _normHyphen,
		Aliases: []string{"nh"},
		Usage:   "normalize various hyphen symbols",
	},
	&cli.BoolFlag{
		Name:    _normDash,
		Aliases: []string{"nem"},
		Usage:   "normalize various em dash symbols",
	},
	&cli.BoolFlag{
		Name:    _normApostrophe,
		Aliases: []string{"na"},
		Usage:   "normalize various apostrophe marks",
	},
	&cli.BoolFlag{
		Name:    _normQuotation,
		Aliases: []string{"nq"},
		Usage:   "normalize various quotation marks",
	},

	// Flags as shortcut for several other flags
	&cli.BoolFlag{
		Name:    _normNumber,
		Aliases: []string{"nn"},
		Usage:   `alias for "--nsub --nsuper"`,
	},
	&cli.BoolFlag{
		Name:    _normMark,
		Aliases: []string{"nm"},
		Usage:   `alias for "--nh --nem --na --nq"`,
	},
}
