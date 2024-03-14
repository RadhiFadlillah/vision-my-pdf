package vision

import (
	"image"
)

type Page struct {
	Image       string
	Paragraphs  []Paragraph
	BoundingBox image.Rectangle
}

func (p Page) Offset(pt image.Point) Page {
	p.BoundingBox = p.BoundingBox.Add(pt)
	for i, pa := range p.Paragraphs {
		p.Paragraphs[i] = pa.Offset(pt)
	}
	return p
}

type Paragraph struct {
	Lines       []Line `json:",omitempty"`
	BoundingBox image.Rectangle
}

func (pa Paragraph) Offset(pt image.Point) Paragraph {
	pa.BoundingBox = pa.BoundingBox.Add(pt)
	for i, l := range pa.Lines {
		pa.Lines[i] = l.Offset(pt)
	}
	return pa
}

type Line struct {
	Words       []Word `json:",omitempty"`
	BoundingBox image.Rectangle
}

func (l Line) Offset(pt image.Point) Line {
	l.BoundingBox = l.BoundingBox.Add(pt)
	for i, w := range l.Words {
		l.Words[i] = w.Offset(pt)
	}
	return l
}

type Word struct {
	Symbols     []Symbol `json:",omitempty"`
	Prefix      string   `json:",omitempty"`
	Suffix      string   `json:",omitempty"`
	BoundingBox image.Rectangle
}

func (w Word) Offset(pt image.Point) Word {
	w.BoundingBox = w.BoundingBox.Add(pt)
	for i, s := range w.Symbols {
		w.Symbols[i] = s.Offset(pt)
	}
	return w
}

type Symbol struct {
	Text        string `json:",omitempty"`
	Prefix      string `json:",omitempty"`
	Suffix      string `json:",omitempty"`
	BoundingBox image.Rectangle
}

func (s Symbol) Offset(pt image.Point) Symbol {
	s.BoundingBox = s.BoundingBox.Add(pt)
	return s
}
