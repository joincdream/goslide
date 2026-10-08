# Goslide Markdown DSL & Slide Authoring Specification

> 🤖 **LLM & AI Agent Directives**:  
> This document serves as the Single Source of Truth (SSOT) specification for LLMs and AI agents generating **Goslide Markdown Presentation Decks**.  
> When tasked with generating a slide deck from a topic, outline, or document, strictly adhere to the **Formatting Rules** and **Directive Catalog** defined below to output deterministic, error-free presentation markdown.

---

## 1. Core Authoring Principles & Format Constraints

* **16:9 Viewport Fixed-Height Principle**:
  * Slides render in a fixed, non-scrolling 1920×1080 canvas.
  * The total line count per slide (including blank lines) must **never exceed 10–12 lines**.
  * Multi-sentence paragraphs are strictly prohibited. Always use concise **bullet points with bold keywords (`- **Keyword**: Explanation`)**.
  * Restrict bullet points to a **maximum of 4 items** per slide.
* **Separation of Slide Content and Speaker Script**:
  * Slide viewports must remain high-signal, visual, and concise.
  * Move all extended explanations, background context, and verbatim talking points into the **`<!-- note: ... -->` block** at the bottom of the slide.
* **Slide Boundary Rules**:
  * Slides must be separated by an isolated **`---` on its own line**.
  * Never insert `---` inside fenced code blocks (```` ``` ````).
* **Directive Scope Rule (`_`)**:
  * Directives prefixed with an underscore (`_`) apply **locally to that single slide only** (e.g., `_layout`, `_backgroundImage`, `_paginate`, `_fragmentStyle`).
  * Directives without an underscore are inherited by all subsequent slides until overridden.

---

## 2. Global Configuration (YAML Frontmatter)

Place YAML Frontmatter at the very top of the markdown document:

```yaml
---
title: "Presentation Title"
author: "Author or Team Name"
theme: "corporate"  # corporate, clean, default, dark, or relative path to custom CSS
size: "16:9"        # 16:9 (default) or 4:3
paginate: true      # Enable slide page numbers (true/false)
header: "Global Header Text"
footer: "Global Footer Text"
---
```

---

## 3. Five Semantic Layouts (`_layout`)

Declare slide layout using the `<!-- _layout: {name} -->` directive:

| Layout | Directive | Purpose & Required Elements |
|---|---|---|
| **Title Cover** | `<!-- _layout: cover -->` | Main deck cover. Large H1 Title + H3 Subtitle + Presenter metadata |
| **Default Body** | `<!-- _layout: default -->` | Standard content slide (default if omitted). H2 Header + 3–4 concise bullets / tables / code |
| **Two Columns** | `<!-- _layout: two-cols -->` | Parallel comparison. **Must insert `<!-- split -->` on its own line between columns** |
| **Section Divider** | `<!-- _layout: section -->` | Mid-deck chapter transition. Large H1 Title alone or with a short quote |
| **Lead / Highlight** | `<!-- _layout: lead -->` | Vertically and horizontally centered takeaway message, Q&A, or Hands-on transition |
| **Blank / Canvas** | `<!-- _layout: blank -->` | Zero padding for full-bleed diagrams, fullscreen images, or large math equations |

---

## 4. Two-Column Split Syntax (`two-cols`)

When using `<!-- _layout: two-cols -->`, you **must insert `<!-- split -->` on its own line** between the left and right column contents:

```markdown
---
<!-- _layout: two-cols -->

## Architecture & Code Comparison

- **Unidirectional Pipeline**: Generates immutable IR from AST
- **Layout Decoupling**: Separates markdown content from presentation styling
- **Multi-Format Export**: Simultaneous HTML, PDF, and PPTX generation

<!-- split -->

```go
func RenderSlide(ctx context.Context, s *Slide) error {
    return pipeline.Execute(ctx, s)
}
```
```

---

## 5. Visual Styling Directives

Control slide appearance using inline HTML comment directives:

* **Background Image**: `<!-- _backgroundImage: "images/bg.png" -->`
* **Background Dim (Darkening Overlay)**: `<!-- _backgroundDim: 0.5 -->` (0.1 to 0.9, ensures text legibility over photo backgrounds)
* **Background Color**: `<!-- _backgroundColor: "#0f172a" -->`
* **Text Color**: `<!-- _color: "#ffffff" -->`
* **Theme Inversion Class**: `<!-- _class: invert -->`
* **Local Visibility Overrides**: `<!-- _header: "" -->`, `<!-- _footer: "" -->`, `<!-- _paginate: false -->`

---

## 6. Interactive Features & Markdown Extensions

### 6.1 Incremental Reveal (Step-by-Step Animation)
Insert `<!-- pause -->` between bullet points or paragraphs. Items following `<!-- pause -->` remain hidden until the presenter advances the slide:

```markdown
- First visible point
<!-- pause -->
- Second point (revealed on next keypress)
<!-- pause -->
- Third point (revealed on subsequent keypress)
```
* **Dim Fragment Animation**: Add `<!-- _fragmentStyle: dim -->` to pre-render unrevealed items at 25% opacity instead of fully hiding them.

### 6.2 Code Blocks & Specific Line Highlighting
Highlight key code lines while automatically dimming unselected lines using `{lines}` notation:

````markdown
```go {2,4-6}
func main() {
    ctx := context.Background() // Line 2: highlighted
    client := NewClient()       // Line 3: dimmed
    for token := range stream { // Lines 4-6: highlighted
        fmt.Println(token)
    }
}
```
````

### 6.3 Media Sizing & Centering
Control image dimensions and alignment directly within the alt-text attribute:

```markdown
![w:400 h:250](images/diagram.png)       <!-- Width 400px, height 250px -->
![w:60% center](images/architecture.png)  <!-- 60% viewport width, centered horizontally -->
![height:300px](images/screenshot.png)   <!-- Height 300px, auto aspect ratio -->
```

### 6.4 Responsive YouTube Embed
Paste standard watch or shortened URLs inside markdown image syntax for responsive 16:9 player embeds:

```markdown
![](https://www.youtube.com/watch?v=dQw4w9WgXcQ)
![](https://youtu.be/dQw4w9WgXcQ?t=1m30s) <!-- Supports start timestamps -->
```

### 6.5 GFM Alerts / Callout Boxes
Render visually distinct callout boxes with dedicated icons:

```markdown
> [!NOTE] Additional Context
> Helpful background info or reference material.

> [!TIP] Performance Optimization
> Best practice recommendations for optimal throughput.

> [!WARNING] Breaking Change
> Critical notices, deprecated parameters, or risk warnings.
```
*(Supported types: `[!NOTE]`, `[!TIP]`, `[!IMPORTANT]`, `[!WARNING]`, `[!CAUTION]`)*

### 6.6 Mermaid Diagrams & LaTeX Math (KaTeX)
* **Mermaid Flowcharts & Sequence Diagrams**:
  ````markdown
  ```mermaid
  graph LR
      A[Markdown] --> B[Goslide Engine]
      B --> C[HTML]
      B --> D[PDF]
      B --> E[PPTX]
  ```
  ````
* **LaTeX Math Equations**: Inline `$E=mc^2$` or block `$$\sum_{i=1}^n x_i$$`

---

## 7. Speaker Notes (Speaker Script)

Place extended explanations and talking points at the bottom of the slide. Visible only in Presenter Console (`P` key) and exported PPTX slide notes:

```markdown
<!-- note:
- Emphasize the latency benefits of item 2 in under 60 seconds.
- Remind attendees about the Q&A session at the end.
-->
```

---

## 8. Output Requirements for LLMs

When generating slide decks from user prompts:
1. **Single Markdown Code Block**: Output the entire presentation inside a single ```` ```markdown ```` code fence from YAML frontmatter to the final slide, without preamble or trailing conversational text.
2. **Mandatory Frontmatter**: Always include `title`, `author`, `theme: "corporate"`, and `paginate: true`.
3. **Strict Height Discipline**: Enforce max 10–12 lines per slide and max 4 bullet points.
4. **Mandatory Two-Column Delimiter**: Always place `<!-- split -->` on its own line between columns when using `<!-- _layout: two-cols -->`.
