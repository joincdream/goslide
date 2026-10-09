package model

import "time"

// SizeRatio defines the aspect ratio of the presentation.
type SizeRatio string

const (
	Ratio16x9 SizeRatio = "16:9"
	Ratio4x3  SizeRatio = "4:3"
)

// LayoutType defines the structural layout of a slide.
type LayoutType string

const (
	LayoutDefault LayoutType = "default"
	LayoutCover   LayoutType = "cover"
	LayoutSection LayoutType = "section"
	LayoutTwoCols LayoutType = "two-cols"
	LayoutLead    LayoutType = "lead"
	LayoutBlank   LayoutType = "blank"
)

// Deck represents the entire immutable presentation deck.
type Deck struct {
	Title       string           `json:"title" yaml:"title"`
	Author      string           `json:"author" yaml:"author"`
	CreatedAt   time.Time        `json:"created_at" yaml:"created_at"`
	GlobalAttrs GlobalDirectives `json:"global_attributes" yaml:",inline"`
	CustomCSS   string           `json:"custom_css,omitempty" yaml:"custom_css"`
	Slides      []*Slide         `json:"slides" yaml:"-"`
	Diagnostics []Diagnostic     `json:"diagnostics,omitempty" yaml:"-"`
}

// GlobalDirectives holds deck-level global attributes defined in Frontmatter.
type GlobalDirectives struct {
	Theme    string     `json:"theme" yaml:"theme"`
	Lang     string     `json:"lang,omitempty" yaml:"lang"`
	Layout   LayoutType `json:"layout" yaml:"layout"`
	Size     SizeRatio  `json:"size" yaml:"size"`
	Paginate bool       `json:"paginate" yaml:"paginate"`
	Header   string     `json:"header,omitempty" yaml:"header"`
	Footer   string     `json:"footer,omitempty" yaml:"footer"`
	Autofit  bool       `json:"autofit" yaml:"autofit"`
}

// SlideDirectives holds scoped directives applied to a single slide.
type SlideDirectives struct {
	Layout          LayoutType `json:"layout,omitempty"`
	Class           []string   `json:"class,omitempty"`
	BackgroundColor string     `json:"background_color,omitempty"`
	BackgroundImage string     `json:"background_image,omitempty"`
	BackgroundDim   string     `json:"background_dim,omitempty"`
	Color           string     `json:"color,omitempty"`
	Header          string     `json:"header,omitempty"`
	Footer          string     `json:"footer,omitempty"`
	Paginate        bool       `json:"paginate"`
	Autofit         bool       `json:"autofit"`
}

// Slide represents a single presentation slide.
type Slide struct {
	Index       int             `json:"index"`
	Layout      LayoutType      `json:"layout"`
	Directives  SlideDirectives `json:"directives"`
	RawContent  string          `json:"-"`
	TitleHTML   string          `json:"title_html,omitempty"`
	HTMLContent string          `json:"html_content"`
	Notes       string          `json:"notes,omitempty"`
	LeftHTML    string          `json:"left_html,omitempty"`
	RightHTML   string          `json:"right_html,omitempty"`
	Diagnostics []Diagnostic    `json:"diagnostics,omitempty"`
}
