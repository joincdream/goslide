package server

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/yundream/goslide/internal/model"
)

// SlidePatch represents an incremental HTML update for a single slide.
type SlidePatch struct {
	Index int    `json:"index"` // 1-based slide index (matches data-slide)
	HTML  string `json:"html"`  // Rendered <section class="slide-card ..."> HTML snippet
}

// DeckSnapshot stores the state of a successfully rendered deck for diffing.
type DeckSnapshot struct {
	Title       string
	SlideCount  int
	GlobalAttrs model.GlobalDirectives
	CustomCSS   string
	SlideHashes []string
}

// DiffResult describes whether a full reload or an incremental patch is needed.
type DiffResult struct {
	NeedsReload     bool
	ReloadReasonKey string // i18n translation key (e.g. "server.hmr.reload.theme")
	ChangedIndices  []int  // 0-based slide indices requiring re-render
}

// reloadRule defines a declarative rule for triggering a full reload.
type reloadRule struct {
	reasonKey string
	isChanged func(old *DeckSnapshot, new *model.Deck) bool
}

// globalReloadRules defines all conditions that mandate a full page reload.
var globalReloadRules = []reloadRule{
	{
		reasonKey: "server.hmr.reload.slide_count",
		isChanged: func(old *DeckSnapshot, new *model.Deck) bool {
			return old.SlideCount != len(new.Slides)
		},
	},
	{
		reasonKey: "server.hmr.reload.title",
		isChanged: func(old *DeckSnapshot, new *model.Deck) bool {
			return old.Title != new.Title
		},
	},
	{
		reasonKey: "server.hmr.reload.custom_css",
		isChanged: func(old *DeckSnapshot, new *model.Deck) bool {
			return old.CustomCSS != new.CustomCSS
		},
	},
	{
		reasonKey: "server.hmr.reload.theme",
		isChanged: func(old *DeckSnapshot, new *model.Deck) bool {
			return old.GlobalAttrs.Theme != new.GlobalAttrs.Theme
		},
	},
	{
		reasonKey: "server.hmr.reload.header",
		isChanged: func(old *DeckSnapshot, new *model.Deck) bool {
			return old.GlobalAttrs.Header != new.GlobalAttrs.Header
		},
	},
	{
		reasonKey: "server.hmr.reload.footer",
		isChanged: func(old *DeckSnapshot, new *model.Deck) bool {
			return old.GlobalAttrs.Footer != new.GlobalAttrs.Footer
		},
	},
	{
		reasonKey: "server.hmr.reload.layout",
		isChanged: func(old *DeckSnapshot, new *model.Deck) bool {
			return old.GlobalAttrs.Layout != new.GlobalAttrs.Layout
		},
	},
	{
		reasonKey: "server.hmr.reload.size",
		isChanged: func(old *DeckSnapshot, new *model.Deck) bool {
			return old.GlobalAttrs.Size != new.GlobalAttrs.Size
		},
	},
	{
		reasonKey: "server.hmr.reload.paginate",
		isChanged: func(old *DeckSnapshot, new *model.Deck) bool {
			return old.GlobalAttrs.Paginate != new.GlobalAttrs.Paginate
		},
	},
	{
		reasonKey: "server.hmr.reload.autofit",
		isChanged: func(old *DeckSnapshot, new *model.Deck) bool {
			return old.GlobalAttrs.Autofit != new.GlobalAttrs.Autofit
		},
	},
}

// ComputeSlideHash generates a deterministic SHA-256 hash representing the slide's content and style.
func ComputeSlideHash(s *model.Slide) string {
	if s == nil {
		return ""
	}
	h := sha256.New()
	_, _ = io.WriteString(h, s.RawContent)
	_, _ = io.WriteString(h, string(s.Layout))
	_, _ = io.WriteString(h, fmt.Sprintf("%+v", s.Directives))
	return hex.EncodeToString(h.Sum(nil))
}

// NewDeckSnapshot creates a snapshot of the deck for incremental diffing.
func NewDeckSnapshot(deck *model.Deck) *DeckSnapshot {
	if deck == nil {
		return &DeckSnapshot{}
	}
	hashes := make([]string, len(deck.Slides))
	for i, s := range deck.Slides {
		hashes[i] = ComputeSlideHash(s)
	}
	return &DeckSnapshot{
		Title:       deck.Title,
		SlideCount:  len(deck.Slides),
		GlobalAttrs: deck.GlobalAttrs,
		CustomCSS:   deck.CustomCSS,
		SlideHashes: hashes,
	}
}

// CompareSnapshots evaluates differences between previous snapshot and current deck.
func CompareSnapshots(oldSnap *DeckSnapshot, newDeck *model.Deck) DiffResult {
	if oldSnap == nil {
		return DiffResult{NeedsReload: true, ReloadReasonKey: "server.hmr.reload.initial"}
	}
	if newDeck == nil {
		return DiffResult{NeedsReload: true, ReloadReasonKey: "server.hmr.reload.initial"}
	}

	// 1. Evaluate declarative reload rules (global attribute changes)
	for _, rule := range globalReloadRules {
		if rule.isChanged(oldSnap, newDeck) {
			return DiffResult{
				NeedsReload:     true,
				ReloadReasonKey: rule.reasonKey,
			}
		}
	}

	// 2. Incremental slide hash comparison ($O(N)$)
	var changed []int
	for i, s := range newDeck.Slides {
		newHash := ComputeSlideHash(s)
		if i >= len(oldSnap.SlideHashes) || oldSnap.SlideHashes[i] != newHash {
			changed = append(changed, i)
		}
	}

	return DiffResult{
		NeedsReload:    false,
		ChangedIndices: changed,
	}
}
