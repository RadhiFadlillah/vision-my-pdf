package cli

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path/filepath"
	fp "path/filepath"
	"strings"

	"github.com/RadhiFadlillah/vision-my-pdf/internal/vision"
)

func prepareTemporaryDirs(rootDir string) error {
	dirs := []string{
		"vision-cache",
		"vision-debug",
		"vision-text",
		"vision-hocr",
	}

	for _, dir := range dirs {
		dir = filepath.Join(rootDir, dir)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			err = fmt.Errorf("temporary dir \"%s\": %w", dir, err)
			return err
		}
	}

	return nil
}

func cleanFileName(fName string) string {
	fName = fp.Base(fName)
	fName = strings.TrimSuffix(fName, fp.Ext(fName))
	return fName
}

func fileExist(f string) bool {
	fs, err := os.Stat(f)
	return err == nil && !fs.IsDir()
}

func saveOcrRaw(path string, page *vision.Page) error {
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	return json.NewEncoder(dst).Encode(&page)
}

func getMidPoint(rect image.Rectangle) image.Point {
	x := rect.Min.X + rect.Dx()/2
	y := rect.Min.Y + rect.Dy()/2
	return image.Pt(x, y)
}
