package parser

import (
	"reflect"
	"testing"

	"github.com/yundream/goslide/internal/model"
)

func TestDirectiveManager_Scopes(t *testing.T) {
	dm := newDirectiveManager(model.GlobalDirectives{
		Theme:    "default",
		Layout:   model.LayoutDefault,
		Paginate: false,
	})

	// Slide 1: local _class and _backgroundColor, plus inheritable header
	s1Raw := `<!--
_class: lead
_backgroundColor: #112233
_backgroundDim: 0.5
header: "Chapter 1"
-->
# Slide 1 Content`

	res1 := dm.processSlide(s1Raw)
	if !reflect.DeepEqual(res1.Directives.Class, []string{"lead"}) {
		t.Errorf("slide 1 class mismatch: got %v, want [lead]", res1.Directives.Class)
	}
	if res1.Directives.BackgroundColor != "#112233" {
		t.Errorf("slide 1 bg color mismatch: got %q, want #112233", res1.Directives.BackgroundColor)
	}
	if res1.Directives.BackgroundDim != "0.5" {
		t.Errorf("slide 1 bg dim mismatch: got %q, want 0.5", res1.Directives.BackgroundDim)
	}
	if res1.Directives.Header != "Chapter 1" {
		t.Errorf("slide 1 header mismatch: got %q, want 'Chapter 1'", res1.Directives.Header)
	}
	if res1.CleanedContent != "# Slide 1 Content" {
		t.Errorf("slide 1 cleaned content mismatch: got %q", res1.CleanedContent)
	}

	// Slide 2: verify local directives did NOT leak, but header DID inherit
	s2Raw := "# Slide 2 Content"
	res2 := dm.processSlide(s2Raw)
	if len(res2.Directives.Class) != 0 {
		t.Errorf("slide 2 should not inherit local class, got %v", res2.Directives.Class)
	}
	if res2.Directives.BackgroundColor != "" {
		t.Errorf("slide 2 should not inherit local bg color, got %q", res2.Directives.BackgroundColor)
	}
	if res2.Directives.Header != "Chapter 1" {
		t.Errorf("slide 2 should inherit header 'Chapter 1', got %q", res2.Directives.Header)
	}

	// Slide 3: override layout locally to cover
	s3Raw := "<!-- _layout: cover -->\n# Title"
	res3 := dm.processSlide(s3Raw)
	if res3.Layout != model.LayoutCover {
		t.Errorf("slide 3 layout mismatch: got %v, want cover", res3.Layout)
	}

	// Slide 4: layout should revert to inherited (default)
	res4 := dm.processSlide("# Next")
	if res4.Layout != model.LayoutDefault {
		t.Errorf("slide 4 layout mismatch: got %v, want default", res4.Layout)
	}
}

func TestDirectiveManager_SpeakerNotes(t *testing.T) {
	dm := newDirectiveManager(model.GlobalDirectives{Layout: model.LayoutDefault})

	t.Run("single line note", func(t *testing.T) {
		input := "# Slide\n<!-- note: Remember to pause here. -->\nSome text."
		res := dm.processSlide(input)
		if res.Notes != "Remember to pause here." {
			t.Errorf("notes mismatch: got %q", res.Notes)
		}
		if res.CleanedContent != "# Slide\n\nSome text." {
			t.Errorf("cleaned content mismatch: got %q", res.CleanedContent)
		}
	})

	t.Run("multiline note", func(t *testing.T) {
		input := "# Slide\n<!--\nnote:\nLine 1 of speech.\nLine 2 of speech.\n-->"
		res := dm.processSlide(input)
		expected := "Line 1 of speech.\nLine 2 of speech."
		if res.Notes != expected {
			t.Errorf("multiline notes mismatch: got %q, want %q", res.Notes, expected)
		}
	})
}
