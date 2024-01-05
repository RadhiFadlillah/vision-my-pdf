package vision

import (
	"image"
)

type Page struct {
	Image       string
	Paragraphs  []Paragraph
	BoundingBox image.Rectangle
}

type Paragraph struct {
	Lines       []Line `json:",omitempty"`
	BoundingBox image.Rectangle
}

type Line struct {
	Words       []Word `json:",omitempty"`
	BoundingBox image.Rectangle
}

type Word struct {
	Symbols     []Symbol `json:",omitempty"`
	Prefix      string   `json:",omitempty"`
	Suffix      string   `json:",omitempty"`
	BoundingBox image.Rectangle
}

type Symbol struct {
	Text        string `json:",omitempty"`
	Prefix      string `json:",omitempty"`
	Suffix      string `json:",omitempty"`
	BoundingBox image.Rectangle
}
