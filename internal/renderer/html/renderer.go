package html

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"strconv"
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
		if r.themeMgr != nil {
			r.themeMgr.SetBaseDir(dir)
		}
	}
}

// WithThemeManager injects a custom theme Manager instance.
func WithThemeManager(mgr *theme.Manager) Option {
	return func(r *HTMLRenderer) {
		r.themeMgr = mgr
		if r.themeMgr != nil && r.baseDir != "" {
			r.themeMgr.SetBaseDir(r.baseDir)
		}
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
	mgr := theme.NewManager(nil)
	mgr.SetBaseDir(".")

	r := &HTMLRenderer{
		themeMgr: mgr,
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
		Theme:       chosenTheme,
		ComposedCSS: template.CSS(composedCSS), // nolint:gosec
		CoreJS:      template.JS(coreJS),       // nolint:gosec
		Slides:      slideViews,
	}

	if err := tmpl.Execute(w, data); err != nil {
		return fmt.Errorf("failed to execute master template: %w", err)
	}

	return nil
}

// RenderSlide renders a single slide into its self-contained <section> HTML element.
func (r *HTMLRenderer) RenderSlide(ctx context.Context, deck *model.Deck, slideIndex int, w io.Writer) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("%w: %w", model.ErrCanceled, err)
	}
	if deck == nil {
		return errors.New("cannot render nil deck")
	}
	if slideIndex < 0 || slideIndex >= len(deck.Slides) {
		return fmt.Errorf("invalid slide index: %d (total: %d)", slideIndex, len(deck.Slides))
	}

	tmpl, err := parseSlideCardTemplate()
	if err != nil {
		return fmt.Errorf("failed to parse slide card template: %w", err)
	}

	s := deck.Slides[slideIndex]
	view, err := r.buildSlideView(slideIndex, s, deck)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(w, view); err != nil {
		return fmt.Errorf("failed to execute slide card template for slide %d: %w", slideIndex, err)
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
	bgColor := s.Directives.BackgroundColor
	bgImage := cleanBackgroundImagePath(s.Directives.BackgroundImage)
	isBgGradient := false

	if strings.Contains(bgImage, "-gradient(") {
		isBgGradient = true
	} else if strings.Contains(bgColor, "-gradient(") {
		bgImage = bgColor
		bgColor = ""
		isBgGradient = true
	} else if bgImage != "" {
		bgImage = filepath.ToSlash(bgImage)
	}

	if r.standalone {
		var err error
		htmlContent, leftHTML, rightHTML, err = r.bundleSlide(s)
		if err != nil {
			return slideTemplateData{}, err
		}
		if bgImage != "" && !isBgGradient {
			bgImage = r.bundler.BundleSingleImage(bgImage)
		}
	}

	bgImageCSS := formatBackgroundImageCSS(bgImage, isBgGradient)

	header := s.Directives.Header
	if header == "" {
		header = deck.GlobalAttrs.Header
	}

	footer := s.Directives.Footer
	if footer == "" {
		footer = deck.GlobalAttrs.Footer
	}

	return slideTemplateData{
		Index:        s.Index,
		IsFirst:      i == 0,
		Layout:       string(s.Layout),
		Classes:      strings.Join(s.Directives.Class, " "),
		BgColor:      template.CSS(bgColor), // nolint:gosec
		BgImage:      bgImageCSS,            // nolint:gosec
		IsBgGradient: isBgGradient,
		BgDim:        template.CSS(resolveBgDim(s.Directives.BackgroundDim, s.Directives.Class)), // nolint:gosec
		Color:        template.CSS(s.Directives.Color),                                           // nolint:gosec
		Header:       header,
		TitleHTML:    template.HTML(s.TitleHTML), // nolint:gosec
		Footer:       footer,
		Paginate:     s.Directives.Paginate || deck.GlobalAttrs.Paginate,
		Autofit:      s.Directives.Autofit || deck.GlobalAttrs.Autofit,
		HTMLContent:  template.HTML(htmlContent), // nolint:gosec
		LeftHTML:     template.HTML(leftHTML),    // nolint:gosec
		RightHTML:    template.HTML(rightHTML),   // nolint:gosec
		Notes:        s.Notes,
	}, nil
}

func resolveBgDim(rawDim string, classes []string) string {
	rawDim = strings.TrimSpace(rawDim)
	if rawDim == "" {
		for _, c := range classes {
			if strings.EqualFold(c, "dim") {
				return "rgba(0, 0, 0, 0.5)"
			}
		}
		return ""
	}

	if strings.HasPrefix(rawDim, "rgba(") || strings.HasPrefix(rawDim, "rgb(") || strings.HasPrefix(rawDim, "#") {
		return rawDim
	}

	clean := strings.TrimSuffix(rawDim, "%")
	if val, err := strconv.ParseFloat(clean, 64); err == nil {
		if strings.HasSuffix(rawDim, "%") {
			val = val / 100.0
		}
		if val < 0 {
			val = 0
		} else if val > 1 {
			val = 1
		}
		return fmt.Sprintf("rgba(0, 0, 0, %g)", val)
	}

	return rawDim
}

func cleanBackgroundImagePath(src string) string {
	trimmed := strings.TrimSpace(src)
	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "url(") && strings.HasSuffix(lower, ")") {
		inner := strings.TrimSpace(trimmed[4 : len(trimmed)-1])
		if len(inner) >= 2 {
			first := inner[0]
			last := inner[len(inner)-1]
			if (first == '\'' && last == '\'') || (first == '"' && last == '"') {
				inner = strings.TrimSpace(inner[1 : len(inner)-1])
			}
		}
		return inner
	}
	return trimmed
}

func formatBackgroundImageCSS(bgImage string, isBgGradient bool) template.CSS {
	if bgImage == "" {
		return ""
	}
	if isBgGradient {
		return template.CSS(bgImage) // nolint:gosec
	}
	safePath := strings.ReplaceAll(bgImage, "'", "%27")
	return template.CSS(fmt.Sprintf("url('%s')", safePath)) // nolint:gosec
}
