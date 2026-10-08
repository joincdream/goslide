package parser

import (
	"regexp"
	"strings"

	"github.com/yundream/goslide/internal/model"
)

var htmlCommentRegex = regexp.MustCompile(`<!--(?s:(.*?))-->`)

type directiveManager struct {
	inheritedDirectives model.SlideDirectives
	inheritedLayout     model.LayoutType
}

func newDirectiveManager(global model.GlobalDirectives) *directiveManager {
	return &directiveManager{
		inheritedDirectives: model.SlideDirectives{
			Header:   global.Header,
			Footer:   global.Footer,
			Paginate: global.Paginate,
			Autofit:  global.Autofit,
		},
		inheritedLayout: global.Layout,
	}
}

type slideParseResult struct {
	CleanedContent string
	Directives     model.SlideDirectives
	Layout         model.LayoutType
	Notes          string
	Diagnostics    []model.Diagnostic
}

// processSlide parses and strips comments, updating slide attributes, state, and collecting diagnostics.
func (dm *directiveManager) processSlide(slideIndex int, rawContent string) slideParseResult {
	diagnostics := inspectSlideComments(slideIndex, rawContent)

	currentDirectives := dm.inheritedDirectives
	currentLayout := dm.inheritedLayout
	var noteChunks []string

	cleanedContent := htmlCommentRegex.ReplaceAllStringFunc(rawContent, func(fullMatch string) string {
		submatches := htmlCommentRegex.FindStringSubmatch(fullMatch)
		if len(submatches) < 2 {
			return fullMatch
		}
		commentBody := strings.TrimSpace(submatches[1])

		// Preserve 2-column split marker and incremental build pause marker
		if strings.EqualFold(commentBody, "split") || strings.EqualFold(commentBody, "pause") {
			return fullMatch
		}

		// Handle speaker notes
		if isNoteComment(commentBody) {
			noteText := extractNoteContent(commentBody)
			if noteText != "" {
				noteChunks = append(noteChunks, noteText)
			}
			return "" // Strip from markdown
		}

		// Handle directives
		parsed, isDirective := parseDirectiveBlock(commentBody)
		if isDirective {
			dm.applyDirectives(parsed, &currentDirectives, &currentLayout)
			return "" // Strip from markdown
		}

		return fullMatch
	})

	return slideParseResult{
		CleanedContent: strings.TrimSpace(cleanedContent),
		Directives:     currentDirectives,
		Layout:         currentLayout,
		Notes:          strings.Join(noteChunks, "\n\n"),
		Diagnostics:    diagnostics,
	}
}

func isNoteComment(body string) bool {
	lower := strings.ToLower(body)
	return strings.HasPrefix(lower, "note:") ||
		strings.HasPrefix(lower, "note\n") ||
		strings.HasPrefix(lower, "note\r\n") ||
		lower == "note"
}

func extractNoteContent(body string) string {
	lower := strings.ToLower(body)
	if strings.HasPrefix(lower, "note:") {
		return strings.TrimSpace(body[5:])
	}
	if strings.HasPrefix(lower, "note") {
		return strings.TrimSpace(body[4:])
	}
	return strings.TrimSpace(body)
}

type parsedDirective struct {
	isLocal bool
	key     string
	value   string
}

func parseDirectiveBlock(body string) ([]parsedDirective, bool) {
	lines := strings.Split(body, "\n")
	var result []parsedDirective

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		colonIdx := strings.Index(trimmed, ":")
		if colonIdx == -1 {
			return nil, false
		}

		key := strings.TrimSpace(trimmed[:colonIdx])
		value := unquoteValue(trimmed[colonIdx+1:])

		isLocal := strings.HasPrefix(key, "_")
		cleanKey := strings.TrimPrefix(key, "_")

		if !isKnownDirectiveKey(cleanKey) {
			return nil, false
		}

		result = append(result, parsedDirective{
			isLocal: isLocal,
			key:     cleanKey,
			value:   value,
		})
	}

	return result, len(result) > 0
}

func unquoteValue(s string) string {
	trimmed := strings.TrimSpace(s)
	if len(trimmed) >= 2 {
		first := trimmed[0]
		last := trimmed[len(trimmed)-1]
		if (first == '"' && last == '"') || (first == '\'' && last == '\'') {
			return trimmed[1 : len(trimmed)-1]
		}
	}
	return trimmed
}

type directiveHandler func(
	value string,
	target *model.SlideDirectives,
	layoutTarget *model.LayoutType,
)

var directiveRegistry = map[string]directiveHandler{
	"class": func(value string, target *model.SlideDirectives, layoutTarget *model.LayoutType) {
		classes := strings.Fields(value)
		target.Class = classes
		for _, c := range classes {
			if strings.EqualFold(c, "lead") || strings.EqualFold(c, "quote") {
				*layoutTarget = model.LayoutLead
				break
			}
		}
	},
	"background": func(value string, target *model.SlideDirectives, _ *model.LayoutType) {
		target.BackgroundColor = value
	},
	"backgroundColor": func(value string, target *model.SlideDirectives, _ *model.LayoutType) {
		target.BackgroundColor = value
	},
	"backgroundImage": func(value string, target *model.SlideDirectives, _ *model.LayoutType) {
		target.BackgroundImage = value
	},
	"backgroundDim": func(value string, target *model.SlideDirectives, _ *model.LayoutType) {
		target.BackgroundDim = value
	},
	"color": func(value string, target *model.SlideDirectives, _ *model.LayoutType) {
		target.Color = value
	},
	"header": func(value string, target *model.SlideDirectives, _ *model.LayoutType) {
		target.Header = value
	},
	"footer": func(value string, target *model.SlideDirectives, _ *model.LayoutType) {
		target.Footer = value
	},
	"paginate": func(value string, target *model.SlideDirectives, _ *model.LayoutType) {
		target.Paginate = (value == "true")
	},
	"layout": func(value string, _ *model.SlideDirectives, layoutTarget *model.LayoutType) {
		*layoutTarget = model.NormalizeLayout(value)
	},
	"autofit": func(value string, target *model.SlideDirectives, _ *model.LayoutType) {
		target.Autofit = (value == "true" || value == "1" || value == "on")
	},
	"fragmentStyle": func(value string, target *model.SlideDirectives, _ *model.LayoutType) {
		if strings.EqualFold(value, "dim") {
			target.Class = append(target.Class, "dim-fragments")
		}
	},
}

func isKnownDirectiveKey(key string) bool {
	_, ok := directiveRegistry[key]
	return ok
}

func (dm *directiveManager) applyDirectives(
	directives []parsedDirective,
	current *model.SlideDirectives,
	currentLayout *model.LayoutType,
) {
	for _, d := range directives {
		applySingleDirective(d.key, d.value, current, currentLayout)
		if !d.isLocal {
			applySingleDirective(d.key, d.value, &dm.inheritedDirectives, &dm.inheritedLayout)
		}
	}
}

func applySingleDirective(
	key, value string,
	target *model.SlideDirectives,
	layoutTarget *model.LayoutType,
) {
	if handler, ok := directiveRegistry[key]; ok {
		handler(value, target, layoutTarget)
	}
}
