package parser

import (
	"context"
	"strings"
	"testing"
)

func TestTransformFragments_ListItems(t *testing.T) {
	input := `<ul>
<li>Item 1</li>
<!-- pause -->
<li>Item 2</li>
<!-- pause -->
<li>Item 3</li>
</ul>`

	output := transformFragments(input)

	if strings.Contains(output, "<!-- pause -->") {
		t.Errorf("expected <!-- pause --> to be removed, got:\n%s", output)
	}

	expectedItem2 := `<li class="fragment" data-fragment-index="1">Item 2</li>`
	if !strings.Contains(output, expectedItem2) {
		t.Errorf("expected output to contain %q, got:\n%s", expectedItem2, output)
	}

	expectedItem3 := `<li class="fragment" data-fragment-index="2">Item 3</li>`
	if !strings.Contains(output, expectedItem3) {
		t.Errorf("expected output to contain %q, got:\n%s", expectedItem3, output)
	}
}

func TestTransformFragments_Paragraphs(t *testing.T) {
	input := `<p>First paragraph</p>
<!-- pause -->
<p>Second paragraph</p>`

	output := transformFragments(input)

	expectedP := `<p class="fragment" data-fragment-index="1">Second paragraph</p>`
	if !strings.Contains(output, expectedP) {
		t.Errorf("expected output to contain %q, got:\n%s", expectedP, output)
	}
}

func TestParser_EndToEndFragment(t *testing.T) {
	md := `---
title: Fragment Test
---

## Slide with Pause

- Bullet 1
<!-- pause -->
- Bullet 2
<!-- pause -->
- Bullet 3
`

	p := NewParser()
	deck, err := p.Parse(context.Background(), strings.NewReader(md))
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if len(deck.Slides) != 1 {
		t.Fatalf("expected 1 slide, got %d", len(deck.Slides))
	}

	content := deck.Slides[0].HTMLContent
	if !strings.Contains(content, `class="fragment" data-fragment-index="1"`) {
		t.Errorf("expected fragment 1 in HTMLContent, got:\n%s", content)
	}
	if !strings.Contains(content, `class="fragment" data-fragment-index="2"`) {
		t.Errorf("expected fragment 2 in HTMLContent, got:\n%s", content)
	}
	if strings.Contains(content, "<!-- pause -->") {
		t.Errorf("pause comment should not remain in HTMLContent:\n%s", content)
	}
}

func TestParser_DimFragmentsDirective(t *testing.T) {
	md := `---
title: Dim Test
---

<!-- _fragmentStyle: dim -->
## Slide with Dim Fragments

- Item 1
<!-- pause -->
- Item 2
`

	p := NewParser()
	deck, err := p.Parse(context.Background(), strings.NewReader(md))
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	foundDim := false
	for _, c := range deck.Slides[0].Directives.Class {
		if c == "dim-fragments" {
			foundDim = true
			break
		}
	}
	if !foundDim {
		t.Errorf("expected 'dim-fragments' in slide classes, got: %v", deck.Slides[0].Directives.Class)
	}
}
