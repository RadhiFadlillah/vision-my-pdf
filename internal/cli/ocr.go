package cli

import (
	"context"
	"fmt"
	fp "path/filepath"
	"sort"
	"sync"

	"github.com/RadhiFadlillah/vision-my-pdf/internal/montage"
	"github.com/RadhiFadlillah/vision-my-pdf/internal/vision"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

func runOCR(montages []montage.Montage, outputDir string, nWorker int64) ([]vision.Page, error) {
	// Prepare concurrent helper
	var wg sync.WaitGroup
	var mut sync.Mutex
	ctx := context.Background()
	sem := semaphore.NewWeighted(nWorker)

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
	for _, montage := range montages {
		// Prepare output for this image
		montageName := cleanFileName(montage.Name())
		ocrOutput := fp.Join(outputDir, montageName+".json")

		// Acquire semaphore
		wg.Add(1)
		if err := sem.Acquire(ctx, 1); err != nil {
			err = fmt.Errorf("ocr semaphore error: %w", err)
			return nil, err
		}

		// Run OCR
		montage := montage
		go func() {
			// Make sure to release semaphore
			defer func() {
				wg.Done()
				sem.Release(1)
			}()

			// Parse image
			pages, err := vision.ParseMontage(ctx, montage)
			if err != nil {
				msg := fmt.Errorf("ocr failed for \"%s\": %w", montageName, err)
				logrus.Warn(msg)
				saveError(err)
				return
			}

			if len(pages) == 0 {
				logrus.Warnf("ocr found no text in \"%s\"", montageName)
				return
			}

			// Save parse result to file
			for _, page := range pages {
				if err = saveOcrRaw(ocrOutput, page); err != nil {
					msg := fmt.Errorf("save ocr result failed for \"%s\": %w", page.Image, err)
					logrus.Warn(msg)
					saveError(err)
					return
				}

				savePage(page)
				logrus.Printf("converted \"%s\"", page.Image)
			}
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
