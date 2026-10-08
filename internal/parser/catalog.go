package parser

import "strings"

// AvailableDirectives returns the list of supported slide-level directive keys (with '_' prefix).
func AvailableDirectives() []string {
	return []string{
		"_layout",
		"_class",
		"_background",
		"_backgroundColor",
		"_backgroundImage",
		"_backgroundDim",
		"_color",
		"_header",
		"_footer",
		"_paginate",
		"_autofit",
		"_fragmentStyle",
	}
}

// knownDirectiveCanonical maps lowercase directive keys to their canonical form.
var knownDirectiveCanonical = map[string]string{
	"layout":          "_layout",
	"class":           "_class",
	"background":      "_background",
	"backgroundcolor": "_backgroundColor",
	"backgroundimage": "_backgroundImage",
	"backgrounddim":   "_backgroundDim",
	"color":           "_color",
	"header":          "_header",
	"footer":          "_footer",
	"paginate":        "_paginate",
	"autofit":         "_autofit",
	"fragmentstyle":   "_fragmentStyle",
}

// MatchCanonicalDirective returns canonical directive with '_' prefix if recognized.
func MatchCanonicalDirective(key string) (string, bool) {
	lower := strings.ToLower(strings.TrimSpace(key))
	lower = strings.TrimPrefix(lower, "_")
	canonical, ok := knownDirectiveCanonical[lower]
	return canonical, ok
}
