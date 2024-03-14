package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"time"

	"github.com/RadhiFadlillah/vision-my-pdf/internal/montage"
	"github.com/RadhiFadlillah/vision-my-pdf/internal/vision"
	"github.com/sirupsen/logrus"
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

		// Get image paths and other relevant files
		imagePaths, oldFiles, err := getRelevantFiles(rootDir)
		if err != nil {
			return err
		}

		// If there are no image, stop
		if len(imagePaths) == 0 {
			return fmt.Errorf("no image detected")
		}

		// Prepare output dirs
		now := time.Now().Format("20060102150405")
		cacheDir := filepath.Join(rootDir, "vision-cache")
		debugDir := filepath.Join(rootDir, "vision-debug")
		backupDir := filepath.Join(rootDir, fmt.Sprintf("vision-backup-%s", now))

		outputDirs := []string{cacheDir}
		if len(oldFiles) > 0 {
			outputDirs = append(outputDirs, backupDir)
		}
		if c.Bool(_genDebug) {
			outputDirs = append(outputDirs, debugDir)
		}

		err = prepareOutputDirs(outputDirs...)
		if err != nil {
			return err
		}

		// Adjust montage size
		montageSize := c.Int(_montageSize)
		if montageSize < 1 {
			montageSize = 1
		} else if montageSize > 5 {
			montageSize = 5
		}

		// Filter images to be montaged
		rewriteOutput := c.Bool(_force)
		var montageQueue []string
		for _, imgPath := range imagePaths {
			// Create absolute path to image
			absPath, err := filepath.Abs(imgPath)
			if err != nil {
				absPath = imgPath
			}

			// Check if OCR cache for this image exists
			imgName := cleanFileName(imgPath)
			ocrOutput := filepath.Join(cacheDir, imgName+".json")

			if fileExist(ocrOutput) {
				page, err := decodePageFile(ocrOutput)
				if err == nil && page != nil && !rewriteOutput {
					logrus.Warnf("skipped \"%s\": already converted", imgName)
					continue
				}
			}

			// Save this image in the queue to be montaged
			montageQueue = append(montageQueue, absPath)
		}

		// Generate montages
		var montages []montage.Montage
		queueSize := len(montageQueue)

		for i := 0; i < queueSize; i += montageSize {
			limit := i + montageSize
			if limit > queueSize {
				limit = queueSize
			}

			montage, err := montage.Create(montageQueue[i:limit]...)
			if err != nil {
				return err
			}

			montages = append(montages, montage)
			logrus.Printf("generate montage for %s", montage.Name())
		}

		// Run OCR concurrently
		pages, err := runOCR(montages, cacheDir, nWorker)
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

		// Save the old text and HOCR to backup dir
		if len(oldFiles) > 0 {
			for _, of := range oldFiles {
				err = copyFile(of, backupDir)
				if err != nil {
					return err
				}
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
			err = saveDebugImages(pages, debugDir)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
func decodePageFile(path string) (*vision.Page, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var page vision.Page
	err = json.NewDecoder(f).Decode(&page)
	if err != nil {
		return nil, err
	}

	return &page, nil
}
