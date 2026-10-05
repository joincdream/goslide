package parser

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	gmparser "github.com/yuin/goldmark/parser"
	gmhtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yundream/goslide/internal/model"
)

var _ model.Parser = (*Parser)(nil)

// Parser parses markdown source into an immutable presentation Deck.
type Parser struct {
	gm goldmark.Markdown
}

// NewParser creates a new Parser initialized with GFM, Chroma highlight, and HTML rendering extensions.
func NewParser() *Parser {
	gm := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			newChromaHighlightExtension("dracula"),
		),
		goldmark.WithParserOptions(
			gmparser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			gmhtml.WithUnsafe(),
		),
	)
	return &Parser{gm: gm}
}

// Parse parses markdown content from an io.Reader into an immutable Deck.
func (p *Parser) Parse(ctx context.Context, r io.Reader) (*model.Deck, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", model.ErrCanceled, err)
	}

	if r == nil {
		return nil, fmt.Errorf("%w: nil reader provided", model.ErrCanceled)
	}

	rawBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read markdown source: %w", err)
	}

	content := normalizeNewlines(rawBytes)

	if ctxErr := ctx.Err(); ctxErr != nil {
		return nil, fmt.Errorf("%w: %v", model.ErrCanceled, ctxErr)
	}

	fm, err := extractFrontmatter(content)
	if err != nil {
		return nil, err
	}

	slides, err := p.buildSlides(ctx, fm)
	if err != nil {
		return nil, err
	}

	deck := &model.Deck{
		Title:       fm.Title,
		Author:      fm.Author,
		CreatedAt:   time.Now(),
		GlobalAttrs: fm.GlobalAttrs,
		CustomCSS:   fm.CustomCSS,
		Slides:      slides,
	}

	return deck, nil
}

func normalizeNewlines(rawBytes []byte) string {
	content := string(bytes.ReplaceAll(rawBytes, []byte("\r\n"), []byte("\n")))
	return string(bytes.ReplaceAll([]byte(content), []byte("\r"), []byte("\n")))
}

func (p *Parser) buildSlides(ctx context.Context, fm *FrontmatterResult) ([]*model.Slide, error) {
	slideChunks := splitSlides(fm.Body)
	dm := newDirectiveManager(fm.GlobalAttrs)
	slides := make([]*model.Slide, 0, len(slideChunks))

	for i, chunk := range slideChunks {
		if loopErr := ctx.Err(); loopErr != nil {
			return nil, fmt.Errorf("%w: %v", model.ErrCanceled, loopErr)
		}

		parsed := dm.processSlide(chunk)
		titleHTML, htmlContent, leftHTML, rightHTML, err := renderSlideContent(p.gm, parsed.Layout, parsed.CleanedContent)
		if err != nil {
			return nil, fmt.Errorf("failed to render slide %d: %w", i+1, err)
		}

		slides = append(slides, &model.Slide{
			Index:       i + 1,
			Layout:      parsed.Layout,
			Directives:  parsed.Directives,
			RawContent:  chunk,
			TitleHTML:   titleHTML,
			HTMLContent: htmlContent,
			Notes:       parsed.Notes,
			LeftHTML:    leftHTML,
			RightHTML:   rightHTML,
		})
	}

	return slides, nil
}
