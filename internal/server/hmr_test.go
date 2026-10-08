package server

import (
	"reflect"
	"testing"

	"github.com/yundream/goslide/internal/model"
)

func createBaseDeck() *model.Deck {
	return &model.Deck{
		Title: "Tech Conference 2026",
		GlobalAttrs: model.GlobalDirectives{
			Theme:    "default",
			Layout:   model.LayoutDefault,
			Size:     model.Ratio16x9,
			Paginate: true,
			Header:   "Goslide Tech Talk",
			Footer:   "Confidential",
			Autofit:  false,
		},
		CustomCSS: "body { font-family: sans-serif; }",
		Slides: []*model.Slide{
			{
				Index:      1,
				Layout:     model.LayoutCover,
				RawContent: "# Welcome\nGoslide Intro",
				Directives: model.SlideDirectives{Layout: model.LayoutCover},
			},
			{
				Index:      2,
				Layout:     model.LayoutDefault,
				RawContent: "## Architecture\nPipeline design",
				Directives: model.SlideDirectives{Layout: model.LayoutDefault},
			},
			{
				Index:      3,
				Layout:     model.LayoutTwoCols,
				RawContent: "## Benchmark\nSpeed comparison",
				Directives: model.SlideDirectives{Layout: model.LayoutTwoCols},
			},
		},
	}
}

func TestCompareSnapshots_Initial(t *testing.T) {
	deck := createBaseDeck()
	res := CompareSnapshots(nil, deck)

	if !res.NeedsReload {
		t.Fatalf("expected NeedsReload=true on initial load, got false")
	}
	if res.ReloadReasonKey != "server.hmr.reload.initial" {
		t.Errorf("expected reason key 'server.hmr.reload.initial', got '%s'", res.ReloadReasonKey)
	}
}

func TestCompareSnapshots_GlobalRulesTable(t *testing.T) {
	base := createBaseDeck()
	oldSnap := NewDeckSnapshot(base)

	tests := []struct {
		name           string
		modify         func(d *model.Deck)
		expectedReason string
	}{
		{
			name: "slide count increased",
			modify: func(d *model.Deck) {
				d.Slides = append(d.Slides, &model.Slide{Index: 4, RawContent: "## New Slide"})
			},
			expectedReason: "server.hmr.reload.slide_count",
		},
		{
			name: "slide count decreased",
			modify: func(d *model.Deck) {
				d.Slides = d.Slides[:2]
			},
			expectedReason: "server.hmr.reload.slide_count",
		},
		{
			name: "title changed",
			modify: func(d *model.Deck) {
				d.Title = "Updated Conference Title"
			},
			expectedReason: "server.hmr.reload.title",
		},
		{
			name: "custom css changed",
			modify: func(d *model.Deck) {
				d.CustomCSS = "body { background: #000; }"
			},
			expectedReason: "server.hmr.reload.custom_css",
		},
		{
			name: "theme changed",
			modify: func(d *model.Deck) {
				d.GlobalAttrs.Theme = "dark"
			},
			expectedReason: "server.hmr.reload.theme",
		},
		{
			name: "header changed",
			modify: func(d *model.Deck) {
				d.GlobalAttrs.Header = "New Global Header"
			},
			expectedReason: "server.hmr.reload.header",
		},
		{
			name: "footer changed",
			modify: func(d *model.Deck) {
				d.GlobalAttrs.Footer = "New Global Footer"
			},
			expectedReason: "server.hmr.reload.footer",
		},
		{
			name: "layout changed",
			modify: func(d *model.Deck) {
				d.GlobalAttrs.Layout = model.LayoutSection
			},
			expectedReason: "server.hmr.reload.layout",
		},
		{
			name: "size changed",
			modify: func(d *model.Deck) {
				d.GlobalAttrs.Size = model.Ratio4x3
			},
			expectedReason: "server.hmr.reload.size",
		},
		{
			name: "paginate flag changed",
			modify: func(d *model.Deck) {
				d.GlobalAttrs.Paginate = false
			},
			expectedReason: "server.hmr.reload.paginate",
		},
		{
			name: "autofit flag changed",
			modify: func(d *model.Deck) {
				d.GlobalAttrs.Autofit = true
			},
			expectedReason: "server.hmr.reload.autofit",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			newDeck := createBaseDeck()
			tc.modify(newDeck)

			res := CompareSnapshots(oldSnap, newDeck)
			if !res.NeedsReload {
				t.Fatalf("expected NeedsReload=true for %s, got false", tc.name)
			}
			if res.ReloadReasonKey != tc.expectedReason {
				t.Errorf("expected reason key '%s', got '%s'", tc.expectedReason, res.ReloadReasonKey)
			}
		})
	}
}

func TestCompareSnapshots_IncrementalSlideDiff(t *testing.T) {
	base := createBaseDeck()
	oldSnap := NewDeckSnapshot(base)

	t.Run("no changes", func(t *testing.T) {
		newDeck := createBaseDeck()
		res := CompareSnapshots(oldSnap, newDeck)

		if res.NeedsReload {
			t.Errorf("expected NeedsReload=false, got true (reason: %s)", res.ReloadReasonKey)
		}
		if len(res.ChangedIndices) != 0 {
			t.Errorf("expected 0 changed slides, got %v", res.ChangedIndices)
		}
	})

	t.Run("single slide content changed", func(t *testing.T) {
		newDeck := createBaseDeck()
		newDeck.Slides[1].RawContent = "## Architecture\nPipeline design with incremental caching"

		res := CompareSnapshots(oldSnap, newDeck)
		if res.NeedsReload {
			t.Fatalf("expected NeedsReload=false, got true (reason: %s)", res.ReloadReasonKey)
		}
		expectedIndices := []int{1}
		if !reflect.DeepEqual(res.ChangedIndices, expectedIndices) {
			t.Errorf("expected changed indices %v, got %v", expectedIndices, res.ChangedIndices)
		}
	})

	t.Run("slide local directives changed", func(t *testing.T) {
		newDeck := createBaseDeck()
		newDeck.Slides[0].Directives.BackgroundColor = "#112233"

		res := CompareSnapshots(oldSnap, newDeck)
		if res.NeedsReload {
			t.Fatalf("expected NeedsReload=false, got true (reason: %s)", res.ReloadReasonKey)
		}
		expectedIndices := []int{0}
		if !reflect.DeepEqual(res.ChangedIndices, expectedIndices) {
			t.Errorf("expected changed indices %v, got %v", expectedIndices, res.ChangedIndices)
		}
	})

	t.Run("multiple slides changed simultaneously", func(t *testing.T) {
		newDeck := createBaseDeck()
		newDeck.Slides[0].RawContent = "# Welcome to Goslide v2"
		newDeck.Slides[2].RawContent = "## Benchmark\nSub-millisecond parsing"

		res := CompareSnapshots(oldSnap, newDeck)
		if res.NeedsReload {
			t.Fatalf("expected NeedsReload=false, got true (reason: %s)", res.ReloadReasonKey)
		}
		expectedIndices := []int{0, 2}
		if !reflect.DeepEqual(res.ChangedIndices, expectedIndices) {
			t.Errorf("expected changed indices %v, got %v", expectedIndices, res.ChangedIndices)
		}
	})
}

func TestComputeSlideHash(t *testing.T) {
	s1 := &model.Slide{
		Index:      1,
		Layout:     model.LayoutDefault,
		RawContent: "Hello World",
	}
	s2 := &model.Slide{
		Index:      1,
		Layout:     model.LayoutDefault,
		RawContent: "Hello World",
	}
	s3 := &model.Slide{
		Index:      1,
		Layout:     model.LayoutDefault,
		RawContent: "Hello World 2",
	}

	h1 := ComputeSlideHash(s1)
	h2 := ComputeSlideHash(s2)
	h3 := ComputeSlideHash(s3)

	if h1 == "" || h1 != h2 {
		t.Errorf("expected identical hashes for same content, got h1=%s, h2=%s", h1, h2)
	}
	if h1 == h3 {
		t.Errorf("expected different hashes for different content")
	}
	if ComputeSlideHash(nil) != "" {
		t.Errorf("expected empty hash for nil slide")
	}
}
