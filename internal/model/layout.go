package model

import "strings"

// ContentFormat defines how the slide body content is structured.
type ContentFormat string

const (
	FormatStandard ContentFormat = "standard"
	FormatTwoCols  ContentFormat = "two-cols"
	FormatCentered ContentFormat = "centered"
	FormatFull     ContentFormat = "full"
)

// MasterLayoutSpec defines the immutable structural specification of a slide master category.
type MasterLayoutSpec struct {
	Type          LayoutType
	SeparateTitle bool
	ContentFormat ContentFormat
	HasTracker    bool
	HasFooter     bool
	ContentHeight int // 840, 960, 1000, 1080
}

// LayoutRegistry is the single source of truth (SSOT) defining all PPT slide masters.
var LayoutRegistry = map[LayoutType]MasterLayoutSpec{
	LayoutDefault: {
		Type:          LayoutDefault,
		SeparateTitle: true,
		ContentFormat: FormatStandard,
		HasTracker:    true,
		HasFooter:     true,
		ContentHeight: 840,
	},
	LayoutTwoCols: {
		Type:          LayoutTwoCols,
		SeparateTitle: true,
		ContentFormat: FormatTwoCols,
		HasTracker:    true,
		HasFooter:     true,
		ContentHeight: 840,
	},
	LayoutCover: {
		Type:          LayoutCover,
		SeparateTitle: false,
		ContentFormat: FormatCentered,
		HasTracker:    false,
		HasFooter:     true,
		ContentHeight: 1000,
	},
	LayoutSection: {
		Type:          LayoutSection,
		SeparateTitle: false,
		ContentFormat: FormatCentered,
		HasTracker:    false,
		HasFooter:     true,
		ContentHeight: 1000,
	},
	LayoutLead: {
		Type:          LayoutLead,
		SeparateTitle: false,
		ContentFormat: FormatCentered,
		HasTracker:    true,
		HasFooter:     true,
		ContentHeight: 960,
	},
	LayoutBlank: {
		Type:          LayoutBlank,
		SeparateTitle: false,
		ContentFormat: FormatFull,
		HasTracker:    false,
		HasFooter:     false,
		ContentHeight: 1080,
	},
}

// LayoutAliases maps alternative names, aliases, and shorthand syntax to canonical LayoutType.
var LayoutAliases = map[string]LayoutType{
	"default":  LayoutDefault,
	"standard": LayoutDefault,
	"cover":    LayoutCover,
	"title":    LayoutCover,
	"section":  LayoutSection,
	"chapter":  LayoutSection,
	"two-cols": LayoutTwoCols,
	"twocols":  LayoutTwoCols,
	"lead":     LayoutLead,
	"quote":    LayoutLead,
	"callout":  LayoutLead,
	"blank":    LayoutBlank,
	"empty":    LayoutBlank,
}

// GetLayoutSpec returns the master specification for the given layout type, falling back to LayoutDefault.
func GetLayoutSpec(t LayoutType) MasterLayoutSpec {
	if spec, ok := LayoutRegistry[t]; ok {
		return spec
	}
	return LayoutRegistry[LayoutDefault]
}

// NormalizeLayout normalizes arbitrary layout string input to canonical LayoutType.
func NormalizeLayout(val string) LayoutType {
	cleaned := strings.ToLower(strings.TrimSpace(val))
	if canonical, ok := LayoutAliases[cleaned]; ok {
		return canonical
	}
	return LayoutDefault
}
