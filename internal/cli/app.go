package cli

import (
	"runtime"
	"sort"

	"github.com/urfave/cli/v2"
)

func NewApp() *cli.App {
	return &cli.App{
		Name:      "vision-my-pdf",
		Usage:     "generate HOCR using Google Vision API, to be used with OCRmyPDF",
		UsageText: "vision-my-pdf [flags] ocrmypdf-dir",
		Flags:     appFlags,
		Action:    appActionHandler(),
	}
}

func appActionHandler() cli.ActionFunc {
	return func(c *cli.Context) error {
		// Check number of workers
		nWorker := c.Int64(_worker)
		if nWorker <= 0 {
			nWorker = int64(runtime.GOMAXPROCS(0))
		}

		// Get root dir
		rootDir, err := getRootDir(c.Args().Slice())
		if err != nil {
			return err
		}

		// Prepare temporary dirs
		err = prepareTemporaryDirs(rootDir)
		if err != nil {
			return err
		}

		// Run OCR concurrently
		pages, err := runOCR(rootDir, OcrOptions{
			NWorker:       nWorker,
			RewriteOutput: c.Bool(_force)})
		if err != nil {
			return err
		}

		// If needed, sort paragraph vertically
		if c.Bool(_sortVertical) {
			for i := range pages {
				sort.SliceStable(pages[i].Paragraphs, func(a, b int) bool {
					rectA := pages[i].Paragraphs[a].BoundingBox
					rectB := pages[i].Paragraphs[b].BoundingBox
					midA := getMidPoint(rectA)
					midB := getMidPoint(rectB)
					return midB.Y > midA.Y
				})
			}
		}

		// Create text from OCR page
		tcl := prepareTextCleaner(c)
		err = savePagesAsText(tcl, pages, rootDir)
		if err != nil {
			return err
		}

		// Create HOCR
		err = savePagesAsHOCR(tcl, pages, rootDir)
		if err != nil {
			return err
		}

		// Generate debug images
		if c.Bool(_genDebug) {
			err = saveDebugImages(pages, rootDir)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
