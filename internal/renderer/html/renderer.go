package html

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/yundream/goslide/internal/model"
	"github.com/yundream/goslide/internal/theme"
)

// Option configures an HTMLRenderer.
type Option func(*HTMLRenderer)

// WithTheme sets the theme name to use (overrides frontmatter theme).
func WithTheme(name string) Option {
	return func(r *HTMLRenderer) {
		r.themeName = name
	}
}

// WithCustomCSS sets the path to an external CSS stylesheet.
func WithCustomCSS(path string) Option {
	return func(r *HTMLRenderer) {
		r.themePath = path
	}
}

// WithStandalone sets whether to inline local image assets as Base64 Data URIs.
func WithStandalone(standalone bool) Option {
	return func(r *HTMLRenderer) {
		r.standalone = standalone
	}
}

// WithBaseDir sets the base directory for resolving relative image assets.
func WithBaseDir(dir string) Option {
	return func(r *HTMLRenderer) {
		r.baseDir = dir
		r.bundler = NewAssetBundler(dir)
	}
}

// WithThemeManager injects a custom theme Manager instance.
func WithThemeManager(mgr *theme.Manager) Option {
	return func(r *HTMLRenderer) {
		r.themeMgr = mgr
	}
}

// HTMLRenderer renders a model.Deck into a standalone interactive HTML document.
type HTMLRenderer struct {
	themeMgr   *theme.Manager
	themeName  string
	themePath  string
	standalone bool
	baseDir    string
	bundler    *AssetBundler
}

// NewRenderer creates a new HTMLRenderer with the given options.
func NewRenderer(opts ...Option) *HTMLRenderer {
	r := &HTMLRenderer{
		themeMgr: theme.NewManager(nil),
		baseDir:  ".",
		bundler:  NewAssetBundler("."),
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

// Render synthesizes the presentation Deck into an HTML stream.
func (r *HTMLRenderer) Render(ctx context.Context, deck *model.Deck, w io.Writer) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("%w: %w", model.ErrCanceled, err)
	}
	if deck == nil {
		return errors.New("cannot render nil deck")
	}

	chosenTheme := r.resolveTheme(deck)
	composedCSS, err := r.themeMgr.ComposeFullCSS(chosenTheme, r.themePath, deck.CustomCSS)
	if err != nil {
		return fmt.Errorf("failed to compose presentation styles: %w", err)
	}

	coreJS, err := r.themeMgr.GetCoreJS()
	if err != nil {
		return fmt.Errorf("failed to load core presentation js: %w", err)
	}

	tmpl, err := parseMasterTemplate()
	if err != nil {
		return fmt.Errorf("failed to parse master template: %w", err)
	}

	slideViews := make([]slideTemplateData, 0, len(deck.Slides))
	for i, s := range deck.Slides {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("%w: %w", model.ErrCanceled, err)
		}
		view, err := r.buildSlideView(i, s, deck)
		if err != nil {
			return err
		}
		slideViews = append(slideViews, view)
	}

	title := deck.Title
	if title == "" {
		title = "Goslide Presentation"
	}

	data := documentTemplateData{
		Title:       title,
		ComposedCSS: template.CSS(composedCSS), // nolint:gosec
		CoreJS:      template.JS(coreJS),       // nolint:gosec
		Slides:      slideViews,
	}

	if err := tmpl.Execute(w, data); err != nil {
		return fmt.Errorf("failed to execute master template: %w", err)
	}

	return nil
}

func (r *HTMLRenderer) resolveTheme(deck *model.Deck) string {
	if r.themeName != "" {
		return r.themeName
	}
	if deck.GlobalAttrs.Theme != "" {
		return deck.GlobalAttrs.Theme
	}
	return theme.DefaultTheme
}

func (r *HTMLRenderer) bundleSlide(s *model.Slide) (string, string, string, error) {
	content, err := r.bundler.BundleImages(s.HTMLContent)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to bundle images for slide %d: %w", s.Index, err)
	}

	left, err := r.bundler.BundleImages(s.LeftHTML)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to bundle left column images for slide %d: %w", s.Index, err)
	}

	right, err := r.bundler.BundleImages(s.RightHTML)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to bundle right column images for slide %d: %w", s.Index, err)
	}

	return content, left, right, nil
}

func (r *HTMLRenderer) buildSlideView(i int, s *model.Slide, deck *model.Deck) (slideTemplateData, error) {
	htmlContent := s.HTMLContent
	leftHTML := s.LeftHTML
	rightHTML := s.RightHTML

	if r.standalone {
		var err error
		htmlContent, leftHTML, rightHTML, err = r.bundleSlide(s)
		if err != nil {
			return slideTemplateData{}, err
		}
	}

	header := s.Directives.Header
	if header == "" {
		header = deck.GlobalAttrs.Header
	}

	footer := s.Directives.Footer
	if footer == "" {
		footer = deck.GlobalAttrs.Footer
	}

	return slideTemplateData{
		Index:       s.Index,
		IsFirst:     i == 0,
		Layout:      string(s.Layout),
		Classes:     strings.Join(s.Directives.Class, " "),
		BgColor:     s.Directives.BackgroundColor,
		BgImage:     s.Directives.BackgroundImage,
		Color:       s.Directives.Color,
		Header:      header,
		Footer:      footer,
		Paginate:    s.Directives.Paginate || deck.GlobalAttrs.Paginate,
		HTMLContent: template.HTML(htmlContent), // nolint:gosec
		LeftHTML:    template.HTML(leftHTML),    // nolint:gosec
		RightHTML:   template.HTML(rightHTML),   // nolint:gosec
		Notes:       s.Notes,
	}, nil
}
