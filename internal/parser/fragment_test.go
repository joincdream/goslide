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
	if !strings.Contains(content, `<li class="fragment" data-fragment-index="1">Bullet 2</li>`) {
		t.Errorf("expected li with fragment 1 in HTMLContent, got:\n%s", content)
	}
	if !strings.Contains(content, `<li class="fragment" data-fragment-index="2">Bullet 3</li>`) {
		t.Errorf("expected li with fragment 2 in HTMLContent, got:\n%s", content)
	}
	if strings.Contains(content, "<!-- pause -->") {
		t.Errorf("pause comment should not remain in HTMLContent:\n%s", content)
	}
	if ulCount := strings.Count(content, "<ul"); ulCount != 1 {
		t.Errorf("expected exactly 1 <ul> tag (merged list), but got %d in:\n%s", ulCount, content)
	}
}

func TestNormalizeListFragments_OrderedList(t *testing.T) {
	md := `---
title: Ordered Pause Test
---

1. First Step
<!-- pause -->
2. Second Step
<!-- pause -->
3. Third Step
`

	p := NewParser()
	deck, err := p.Parse(context.Background(), strings.NewReader(md))
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	content := deck.Slides[0].HTMLContent
	if olCount := strings.Count(content, "<ol"); olCount != 1 {
		t.Errorf("expected exactly 1 <ol> tag (merged list), but got %d in:\n%s", olCount, content)
	}
	if !strings.Contains(content, `<li class="fragment" data-fragment-index="1">Second Step</li>`) {
		t.Errorf("expected li with fragment 1 in HTMLContent, got:\n%s", content)
	}
	if !strings.Contains(content, `<li class="fragment" data-fragment-index="2">Third Step</li>`) {
		t.Errorf("expected li with fragment 2 in HTMLContent, got:\n%s", content)
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

func TestNormalizeListFragments_RealWorldSlide16(t *testing.T) {
	md := `---
title: Real World Demo Test
---

<!-- _layout: section -->
<!-- _fragmentStyle: dim -->

# Engineering Conclusions & Roadmap
### "Moving Beyond Simple Prompts toward Autonomous Agents & Domain Optimization"

* **Deterministic Verification** — Schema validation gates to eliminate hallucinations
<!-- pause -->
* **Enterprise Private Serving** — Open-weight models for TCO & data sovereignty
<!-- pause -->
* **Real-Time Observability** — Live TTFT latency & token throughput metrics
`

	p := NewParser()
	deck, err := p.Parse(context.Background(), strings.NewReader(md))
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	content := deck.Slides[0].HTMLContent
	if ulCount := strings.Count(content, "<ul"); ulCount != 1 {
		t.Errorf("expected exactly 1 <ul> tag for slide 16, got %d in:\n%s", ulCount, content)
	}

	expectedLi1 := `<li class="fragment" data-fragment-index="1"><strong>Enterprise Private Serving</strong>`
	if !strings.Contains(content, expectedLi1) {
		t.Errorf("expected fragment 1 on Enterprise Private Serving, got:\n%s", content)
	}

	expectedLi2 := `<li class="fragment" data-fragment-index="2"><strong>Real-Time Observability</strong>`
	if !strings.Contains(content, expectedLi2) {
		t.Errorf("expected fragment 2 on Real-Time Observability, got:\n%s", content)
	}
}
