package parser

import (
	"bytes"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type chromaHighlightExtension struct {
	theme string
}

func newChromaHighlightExtension(theme string) *chromaHighlightExtension {
	if theme == "" {
		theme = "dracula"
	}
	return &chromaHighlightExtension{theme: theme}
}

func (e *chromaHighlightExtension) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(newChromaRenderer(e.theme), 100),
	))
}

type chromaRenderer struct {
	theme string
}

func newChromaRenderer(theme string) *chromaRenderer {
	return &chromaRenderer{theme: theme}
}

func (r *chromaRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindFencedCodeBlock, r.renderFencedCodeBlock)
	reg.Register(ast.KindCodeBlock, r.renderCodeBlock)
}

func (r *chromaRenderer) renderFencedCodeBlock(
	w util.BufWriter, source []byte, node ast.Node, entering bool,
) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n, ok := node.(*ast.FencedCodeBlock)
	if !ok {
		return ast.WalkContinue, nil
	}
	language := string(n.Language(source))

	var buf bytes.Buffer
	lines := n.Lines()
	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		buf.Write(line.Value(source))
	}

	return r.formatCode(w, language, buf.String())
}

func (r *chromaRenderer) renderCodeBlock(
	w util.BufWriter, source []byte, node ast.Node, entering bool,
) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n, ok := node.(*ast.CodeBlock)
	if !ok {
		return ast.WalkContinue, nil
	}
	var buf bytes.Buffer
	lines := n.Lines()
	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		buf.Write(line.Value(source))
	}

	return r.formatCode(w, "", buf.String())
}

func (r *chromaRenderer) formatCode(
	w util.BufWriter, language, code string,
) (ast.WalkStatus, error) {
	lexer := lexers.Get(language)
	if lexer == nil {
		lexer = lexers.Analyse(code)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	style := styles.Get(r.theme)
	if style == nil {
		style = styles.Fallback
	}

	formatter := chromahtml.New(
		chromahtml.WithClasses(false),
		chromahtml.TabWidth(4),
	)

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return ast.WalkContinue, err
	}

	if err := formatter.Format(w, style, iterator); err != nil {
		return ast.WalkContinue, err
	}

	return ast.WalkSkipChildren, nil
}
