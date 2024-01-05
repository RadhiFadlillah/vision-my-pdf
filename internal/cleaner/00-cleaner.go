package cleaner

import (
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Cleaner struct {
	transform.Transformer
}

func New(cleaners ...transform.Transformer) Cleaner {
	var transformers []transform.Transformer
	transformers = append(transformers, norm.NFKD)
	transformers = append(transformers, cleaners...)
	transformers = append(transformers, norm.NFKC)
	return Cleaner{transform.Chain(transformers...)}
}

func (c Cleaner) Clean(s string) string {
	cleaned, _, err := transform.String(c.Transformer, s)
	if err != nil {
		return s
	} else {
		return cleaned
	}
}
