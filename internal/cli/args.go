package cli

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"
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

func getImageNames(dir string) ([]string, error) {
	// Fetch all `*_ocr.png`
	var names []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("dir \"%s\": %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		entryName := entry.Name()
		if isPNG(entryName) && strings.HasSuffix(entryName, "_ocr.png") {
			name := filepath.Join(dir, entryName)
			names = append(names, name)
		}
	}

	return names, nil
}

func isPNG(fName string) bool {
	ext := filepath.Ext(fName)
	mimeType := mime.TypeByExtension(ext)
	return mimeType == "image/png"
}
