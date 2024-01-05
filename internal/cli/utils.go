package cli

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"mime"
	"os"
	fp "path/filepath"
	"strings"

	"github.com/RadhiFadlillah/vision-my-pdf/internal/vision"
)

func getRootDir(args []string) (string, error) {
	// We only accept one arguments
	if len(args) != 1 {
		return "", fmt.Errorf("need exactly one arg")
	}

	// Make sure the argument exist and is directory
	arg := args[0]
	fs, err := os.Stat(arg)
	if err != nil {
		return "", fmt.Errorf("arg \"%s\": %w", arg, err)
	}

	if !fs.IsDir() {
		return "", fmt.Errorf("arg is not dir")
	}

	return arg, nil
}

func getRelevantFiles(dir string) (images, oldFiles []string, err error) {
	// Fetch entries in this dir
	entries, err := os.ReadDir(dir)
	if err != nil {
		err = fmt.Errorf("dir \"%s\": %w", dir, err)
		return
	}

	// Extract relevant files
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		entryName := fp.Join(dir, entry.Name())
		switch {
		case isPNG(entryName) && strings.HasSuffix(entryName, "_ocr.png"):
			images = append(images, entryName)
		case strings.HasSuffix(entryName, "_ocr_hocr.hocr"),
			strings.HasSuffix(entryName, "_ocr_hocr.txt"):
			oldFiles = append(oldFiles, entryName)
		}
	}

	return
}

func isPNG(fName string) bool {
	ext := fp.Ext(fName)
	mimeType := mime.TypeByExtension(ext)
	return mimeType == "image/png"
}

func prepareOutputDirs(dirs ...string) error {
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			err = fmt.Errorf("output dir \"%s\": %w", fp.Base(dir), err)
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

func copyFile(srcPath string, dstDir string) error {
	var err error

	// Open source file
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	// Create destination file
	dstPath := fp.Join(dstDir, fp.Base(srcPath))
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy the file
	_, err = io.Copy(dst, src)
	return err
}
