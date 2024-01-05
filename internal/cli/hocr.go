package cli

import (
	"fmt"
	"image"
	"os"
	fp "path/filepath"
	"regexp"
	"strings"

	"github.com/RadhiFadlillah/vision-my-pdf/internal/cleaner"
	"github.com/RadhiFadlillah/vision-my-pdf/internal/vision"
	"github.com/go-shiori/dom"
)

var rxSymbolOnly = regexp.MustCompile(`^[^\p{L}\p{N}\s]+$`)

func savePagesAsHOCR(tcl cleaner.Cleaner, pages []vision.Page, rootDir string) error {
	// Process each page
	for _, page := range pages {
		// Prepare output for this page
		imgName := cleanFileName(page.Image)
		textOutput := fp.Join(rootDir, "vision-hocr", imgName) + "_hocr.hocr"

		// Build HOCR for this page
		pageHOCR := pageToHOCR(tcl, page)

		// Save text to storage
		err := os.WriteFile(textOutput, []byte(pageHOCR), os.ModePerm)
		if err != nil {
			return fmt.Errorf("save HOCR failed for \"%s\": %w", imgName, err)
		}
	}

	return nil
}

func pageToHOCR(tcl cleaner.Cleaner, page vision.Page) string {
	// Prepare counter
	var paragraphCounter int
	var lineCounter int
	var wordCounter int

	// Create HTML document
	doc := dom.CreateElement("html")
	dom.SetAttribute(doc, "xmlns", "http://www.w3.org/1999/xhtml")
	dom.SetAttribute(doc, "xml:lang", "en")
	dom.SetAttribute(doc, "lang", "en")

	// Prepare head and put it in document
	head := dom.CreateElement("head")
	dom.AppendChild(doc, head)

	// Prepare metas and put it in head
	meta1 := dom.CreateElement("meta")
	dom.SetAttribute(meta1, "http-equiv", "Content-Type")
	dom.SetAttribute(meta1, "content", "text/html;charset=utf-8")
	dom.AppendChild(head, meta1)

	meta2 := dom.CreateElement("meta")
	dom.SetAttribute(meta2, "name", "ocr-system")
	dom.SetAttribute(meta2, "content", "Google Vision")
	dom.AppendChild(head, meta2)

	meta3 := dom.CreateElement("meta")
	dom.SetAttribute(meta3, "name", "ocr-capabilities")
	dom.SetAttribute(meta3, "content", "ocr_page ocr_carea ocr_par ocr_line ocrx_word")
	dom.AppendChild(head, meta3)

	// Prepare body and put it in document
	body := dom.CreateElement("body")
	dom.AppendChild(doc, body)

	// Process the page
	// Create element for page, then put it to body
	divPage := dom.CreateElement("div")
	dom.SetAttribute(divPage, "class", "ocr_page")
	dom.SetAttribute(divPage, "id", "page_1")
	dom.SetAttribute(divPage, "title", rectToString(page.BoundingBox))
	dom.AppendChild(body, divPage)

	// Process each paragraph
	for _, p := range page.Paragraphs {
		paragraphCounter++

		// Create element for c-area, then put it to page
		divCarea := dom.CreateElement("div")
		dom.SetAttribute(divCarea, "class", "ocr_carea")
		dom.SetAttribute(divCarea, "id", fmt.Sprintf("block_1_%d", paragraphCounter))
		dom.SetAttribute(divCarea, "title", rectToString(p.BoundingBox))
		dom.AppendChild(divPage, divCarea)

		// Create element for paragraph, then put it to c-area
		pPar := dom.CreateElement("p")
		dom.SetAttribute(pPar, "class", "ocr_par")
		dom.SetAttribute(pPar, "id", fmt.Sprintf("par_1_%d", paragraphCounter))
		dom.SetAttribute(pPar, "title", rectToString(p.BoundingBox))
		dom.AppendChild(divCarea, pPar)

		// Process each line in paragraph
		for _, l := range p.Lines {
			lineCounter++

			// Create element for line, then put it to paragraph
			spanLine := dom.CreateElement("span")
			dom.SetAttribute(spanLine, "class", "ocr_line")
			dom.SetAttribute(spanLine, "id", fmt.Sprintf("line_1_%d", lineCounter))
			dom.SetAttribute(spanLine, "title", rectToString(l.BoundingBox))
			dom.AppendChild(pPar, spanLine)

			// Process each word in line
			for _, w := range l.Words {
				wordCounter++

				// Create element for word, then put it to line
				spanWord := dom.CreateElement("span")
				dom.SetAttribute(spanWord, "class", "ocrx_word")
				dom.SetAttribute(spanWord, "id", fmt.Sprintf("word_1_%d", wordCounter))
				dom.SetAttribute(spanWord, "title", rectToString(w.BoundingBox))
				dom.AppendChild(spanLine, spanWord)

				// Get the previous sibling of span word
				prevSpanWord := dom.PreviousElementSibling(spanWord)

				// Set word text
				wordText := wordToText(w)
				wordText = tcl.Clean(wordText)
				wordText = strings.TrimSpace(wordText)

				// If its text only contains symbol, put it in the previous span
				if rxSymbolOnly.MatchString(wordText) && prevSpanWord != nil {
					dom.SetTextContent(prevSpanWord, dom.TextContent(prevSpanWord)+wordText)
					spanLine.RemoveChild(spanWord)
				} else {
					dom.SetTextContent(spanWord, wordText)
				}
			}
		}
	}

	// Return the final string
	return dom.OuterHTML(doc)
}

func rectToString(rect image.Rectangle) string {
	return fmt.Sprintf("bbox %d %d %d %d",
		rect.Min.X, rect.Min.Y,
		rect.Max.X, rect.Max.Y)
}
