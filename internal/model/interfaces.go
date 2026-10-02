package model

import (
	"context"
	"io"
)

// Parser parses markdown source and builds an immutable presentation Deck.
type Parser interface {
	Parse(ctx context.Context, r io.Reader) (*Deck, error)
}

// Renderer renders Deck into a specific stream format (primarily HTML).
type Renderer interface {
	Render(ctx context.Context, deck *Deck, w io.Writer) error
}

// Exporter writes Deck into external file formats (PDF, PPTX) on disk.
type Exporter interface {
	Export(ctx context.Context, deck *Deck, outputPath string) error
}
