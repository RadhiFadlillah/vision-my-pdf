package vision

import (
	"bytes"
	"context"
	"image/png"
	"path/filepath"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
	visionpb "cloud.google.com/go/vision/v2/apiv1/visionpb"
	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/imgio"
)

func ParseImage(ctx context.Context, file string, keepHyphen bool) (*Page, error) {
	// Open vision client API
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, err
	}

	// Open file
	img, err := imgio.Open(file)
	if err != nil {
		return nil, err
	}

	// Make sure image is not empty
	bounds := img.Bounds().Size()
	if valid := bounds.X > 1 && bounds.Y > 1; !valid {
		return nil, nil
	}

	// Invert image, since Google vision seems to yield better performance
	// with white text on black background.
	inverted := effect.Invert(img)

	// Encode to the new reader
	var buf bytes.Buffer
	err = png.Encode(&buf, inverted)
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

	// Create absolute path to file
	absPath, err := filepath.Abs(file)
	if err != nil {
		absPath = file
	}

	// Parse page
	result := Page{Image: absPath}
	for _, page := range annotations.Pages {
		for _, block := range page.Blocks {
			for _, paragraph := range block.Paragraphs {
				p := parseParagraph(paragraph, keepHyphen)
				result.Paragraphs = append(result.Paragraphs, p)
			}
		}
	}

	return &result, nil
}

func parseParagraph(paragraph *visionpb.Paragraph, keepHyphen bool) Paragraph {
	// Prepare result
	result := Paragraph{
		BoundingBox: bpToRect(paragraph.BoundingBox),
	}

	// Process each word inside it
	var words []Word
	for _, word := range paragraph.Words {
		parsedWord := parseWord(word, keepHyphen)
		if len(parsedWord.Symbols) > 0 {
			words = append(words, parsedWord)
		}
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

func parseWord(word *visionpb.Word, keepHyphen bool) Word {
	// Prepare result
	result := Word{
		BoundingBox: bpToRect(word.BoundingBox),
	}

	// Process each symbol inside it
	for _, symbol := range word.Symbols {
		// Check if symbol has break
		prefix, suffix := createBreakCharacter(symbol, keepHyphen)

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

func createBreakCharacter(symbol *visionpb.Symbol, keepHyphen bool) (prefix, suffix string) {
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
		bc = "\n"
		if keepHyphen {
			bc = "-\n"
		}
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
