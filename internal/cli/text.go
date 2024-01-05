package cli

import (
	"fmt"
	"os"
	fp "path/filepath"
	"regexp"
	"strings"

	"github.com/RadhiFadlillah/vision-my-pdf/internal/cleaner"
	"github.com/RadhiFadlillah/vision-my-pdf/internal/vision"
	"golang.org/x/text/unicode/bidi"
)

var rxSpaces = regexp.MustCompile(` +`)

func savePagesAsText(tcl cleaner.Cleaner, pages []vision.Page, rootDir string) error {
	// Process each page
	for _, page := range pages {
		// Prepare output for this page
		imgName := cleanFileName(page.Image)
		textOutput := fp.Join(rootDir, "vision-text", imgName) + ".txt"

		// Build text for this page
		pageText := pageToText(page)
		pageText = tcl.Clean(pageText)

		// Save text to storage
		err := os.WriteFile(textOutput, []byte(pageText), os.ModePerm)
		if err != nil {
			return fmt.Errorf("save text failed for \"%s\": %w", imgName, err)
		}
	}

	return nil
}

func pageToText(page vision.Page) string {
	var sb strings.Builder
	for _, p := range page.Paragraphs {
		sb.WriteString(paragraphToText(p))
		sb.WriteString("\n\n")
	}
	return sb.String()
}

func paragraphToText(p vision.Paragraph) string {
	// Extract texts from lines
	var sb strings.Builder
	for _, l := range p.Lines {
		sb.WriteString(lineToText(l))
	}
	text := sb.String()

	// Clean up text
	text = rxSpaces.ReplaceAllString(text, " ")
	text = splitMixedDirParagraph(text)
	text = strings.TrimSpace(text)

	return text
}

func lineToText(l vision.Line) string {
	var sb strings.Builder
	for _, w := range l.Words {
		sb.WriteString(wordToText(w))
	}
	return sb.String()
}

func wordToText(w vision.Word) string {
	var sb strings.Builder
	sb.WriteString(w.Prefix)
	for _, s := range w.Symbols {
		sb.WriteString(s.Prefix)
		sb.WriteString(s.Text)
		sb.WriteString(s.Suffix)
	}
	sb.WriteString(w.Suffix)
	return sb.String()
}

func splitMixedDirParagraph(s string) string {
	var p bidi.Paragraph
	_, err := p.SetString(s)
	if err != nil {
		return s
	}

	orders, err := p.Order()
	if err != nil {
		return s
	}

	var sb strings.Builder
	for i := 0; i < orders.NumRuns(); i++ {
		r := orders.Run(i)
		sb.WriteString(r.String() + "\n")
	}

	return sb.String()
}
