package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	fp "path/filepath"
	"sort"
	"sync"

	"github.com/RadhiFadlillah/vision-my-pdf/internal/vision"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

type OcrOptions struct {
	NWorker       int64
	RewriteOutput bool
	KeepHyphen    bool
}

func runOCR(rootDir string, opts OcrOptions) ([]vision.Page, error) {
	// Get image paths
	imagePaths, err := getImageNames(rootDir)
	if err != nil {
		return nil, err
	}

	// If there are no image, stop
	if len(imagePaths) == 0 {
		return nil, fmt.Errorf("no image detected")
	}

	// Prepare concurrent helper
	var wg sync.WaitGroup
	var mut sync.Mutex
	ctx := context.Background()
	sem := semaphore.NewWeighted(opts.NWorker)

	// Prepare output and helper functions
	var errors []error
	var pages []vision.Page

	saveError := func(err error) {
		mut.Lock()
		defer mut.Unlock()
		errors = append(errors, err)
	}

	savePage := func(p vision.Page) {
		mut.Lock()
		defer mut.Unlock()
		pages = append(pages, p)
	}

	// Run OCR concurrently
	for _, imgPath := range imagePaths {
		// Prepare output for this image
		imgName := cleanFileName(imgPath) + ".json"
		ocrOutput := fp.Join(rootDir, "vision-cache", imgName)

		// Check if output exists
		if fileExist(ocrOutput) {
			page, err := decodePageFile(ocrOutput)
			if err == nil && page != nil && !opts.RewriteOutput {
				savePage(*page)
				logrus.Warnf("skipped \"%s\": already converted", imgPath)
				continue
			}
		}

		// Acquire semaphore
		wg.Add(1)
		if err := sem.Acquire(ctx, 1); err != nil {
			err = fmt.Errorf("ocr semaphore error: %w", err)
			return nil, err
		}

		// Run OCR
		imgPath := imgPath
		go func() {
			// Make sure to release semaphore
			defer func() {
				wg.Done()
				sem.Release(1)
			}()

			// Parse image
			page, err := vision.ParseImage(ctx, imgPath, opts.KeepHyphen)
			if err != nil {
				msg := fmt.Errorf("ocr failed for \"%s\": %w", imgPath, err)
				logrus.Warn(msg)
				saveError(err)
				return
			}

			if page == nil {
				logrus.Warnf("ocr found no text in \"%s\"", imgPath)
				return
			}

			// Save parse result to file
			if err = saveOcrRaw(ocrOutput, page); err != nil {
				msg := fmt.Errorf("save ocr result failed for \"%s\": %w", imgPath, err)
				logrus.Warn(msg)
				saveError(err)
				return
			}

			savePage(*page)
			logrus.Printf("converted \"%s\"", imgPath)
		}()
	}

	// Wait until all goroutine finished
	wg.Wait()

	// Print all error
	if nError := len(errors); nError > 0 {
		for _, err := range errors {
			logrus.Errorln(err)
		}
		return nil, fmt.Errorf("ocr fail with %d error(s)", nError)
	}

	// Sort pages by its file name
	sort.SliceStable(pages, func(a, b int) bool {
		return pages[a].Image < pages[b].Image
	})

	logrus.Print("ocr finished")
	return pages, nil
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
