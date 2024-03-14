package vision

import (
	"bytes"
	"context"
	"image/png"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
	visionpb "cloud.google.com/go/vision/v2/apiv1/visionpb"
	"github.com/RadhiFadlillah/vision-my-pdf/internal/montage"
)

func ParseMontage(ctx context.Context, montage montage.Montage) ([]Page, error) {
	// Open vision client API
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, err
	}

	// Make sure image is not empty
	bounds := montage.Image.Bounds().Size()
	if valid := bounds.X > 1 && bounds.Y > 1; !valid {
		return nil, nil
	}

	// Encode to the new reader
	var buf bytes.Buffer
	err = png.Encode(&buf, montage.Image)
	if err != nil {
		return nil, err
	}

	// Decode visionImg for Google vision
	r := bytes.NewReader(buf.Bytes())
	visionImg, err := vision.NewImageFromReader(r)
	if err != nil {
		return nil, err
	}

	// Look for texts within image
	annotations, err := client.DetectDocumentText(ctx, visionImg, nil)
	if err != nil {
		return nil, err
	}

	if annotations == nil {
		return nil, nil
	}

	// Extract each paragraphs from OCR result
	var paragraphs []Paragraph
	for _, page := range annotations.Pages {
		for _, block := range page.Blocks {
			for _, paragraph := range block.Paragraphs {
				p := parseParagraph(paragraph)
				paragraphs = append(paragraphs, p)
			}
		}
	}

	// Split paragraphs to each page
	var pages []Page
	var paragraphCursor int

	for i, imgPath := range montage.Paths {
		// Create the page
		pageYLimit := montage.YLimits[i]
		page := Page{
			Image:       imgPath,
			BoundingBox: montage.Bounds[i],
		}

		// Save paragraphs for current page
		for paragraphCursor < len(paragraphs) {
			paragraph := paragraphs[paragraphCursor]
			if paragraph.BoundingBox.Max.Y >= pageYLimit {
				break
			}

			page.Paragraphs = append(page.Paragraphs, paragraph)
			paragraphCursor++
		}

		pages = append(pages, page)
	}

	return pages, nil
}

func parseParagraph(paragraph *visionpb.Paragraph) Paragraph {
	// Prepare result
	result := Paragraph{
		BoundingBox: bpToRect(paragraph.BoundingBox),
	}

	// Process each word inside it
	var words []Word
	for _, word := range paragraph.Words {
		parsedWord := parseWord(word)

		// If parsed word is empty, skip
		if len(parsedWord.Symbols) == 0 {
			continue
		}

		// If there are no separator between parsed word and the last one, merge it
		if nWord := len(words); nWord > 0 {
			lastWord := words[nWord-1]
			if lastWord.Suffix == "" && parsedWord.Prefix == "" {
				lastWord.Symbols = append(lastWord.Symbols, parsedWord.Symbols...)
				lastWord.BoundingBox = lastWord.BoundingBox.Union(parsedWord.BoundingBox)
				lastWord.Suffix = parsedWord.Suffix
				words[nWord-1] = lastWord
				continue
			}
		}

		// Save the parsed word
		words = append(words, parsedWord)
	}

	// Separate word to lines
	var lineWords [][]Word
	var cursor int

	for i := 1; i < len(words); i++ {
		currentWord := words[i]
		previousWord := words[i-1]

		// Move suffix from previous to current word
		currentWord.Prefix = previousWord.Suffix + currentWord.Prefix

		// If current word has new line, save the new line
		if strings.Contains(currentWord.Prefix, "\n") {
			newLineWords := append([]Word{}, words[cursor:i]...)
			lineWords = append(lineWords, newLineWords)
			cursor = i
		}
	}

	// Save the leftover word
	if cursor < len(words) {
		newLineWords := append([]Word{}, words[cursor:]...)
		lineWords = append(lineWords, newLineWords)
	}

	// Save words as lines
	for _, lw := range lineWords {
		// Skip if line is empty
		if len(lw) == 0 {
			continue
		}

		// Create bounding box
		box := lw[0].BoundingBox
		for i := 1; i < len(lw); i++ {
			box = box.Union(lw[i].BoundingBox)
		}

		result.Lines = append(result.Lines, Line{
			Words:       lw,
			BoundingBox: box,
		})
	}

	return result
}

func parseWord(word *visionpb.Word) Word {
	// Prepare result
	result := Word{
		BoundingBox: bpToRect(word.BoundingBox),
	}

	// Process each symbol inside it
	for _, symbol := range word.Symbols {
		// Check if symbol has break
		prefix, suffix := createBreakCharacter(symbol)

		// Save the symbol
		s := Symbol{
			Text:        symbol.Text,
			Prefix:      prefix,
			Suffix:      suffix,
			BoundingBox: bpToRect(symbol.BoundingBox),
		}

		// Add symbol to word
		result.Symbols = append(result.Symbols, s)
	}

	// Take the prefix and suffix from symbols
	if nSymbols := len(result.Symbols); nSymbols > 0 {
		result.Prefix = result.Symbols[0].Prefix
		result.Symbols[0].Prefix = ""

		result.Suffix = result.Symbols[nSymbols-1].Suffix
		result.Symbols[nSymbols-1].Suffix = ""
	}

	return result
}

func createBreakCharacter(symbol *visionpb.Symbol) (prefix, suffix string) {
	// Make sure symbol has property
	prop := symbol.Property
	if prop == nil {
		return
	}

	// Make sure break detected
	db := prop.DetectedBreak
	if db == nil {
		return
	}

	// Create break character
	var bc string

	switch db.Type {
	case visionpb.TextAnnotation_DetectedBreak_HYPHEN:
		bc = "-\n"
	case visionpb.TextAnnotation_DetectedBreak_LINE_BREAK,
		visionpb.TextAnnotation_DetectedBreak_EOL_SURE_SPACE:
		bc = "\n"
	default:
		bc = " "
	}

	// Put the break where it belong
	if db.IsPrefix {
		prefix = bc
	} else {
		suffix = bc
	}

	return
}
