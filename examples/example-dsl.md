# Goslide DSL Specification & Golden Example Guide

[English](example-dsl.md) | [한국어](example-dsl.ko.md)

> This document serves as the Single Source of Truth (SSOT) specification for LLMs when generating **Goslide Markdown presentation decks**.
> It consists of three structured sections:
> 1. Semantic Rules & Constraints  
> 2. Directive Catalog  
> 3. Golden Example Deck  

---

## 1. Semantic Rules & Format Constraints

* **16:9 Viewport Fixed-Height Principle**: Slides are rendered within a fixed, non-scrolling screen (1920×1080).
  * The total line count per slide (including empty lines) must **not exceed 10–12 lines**.
  * Multi-sentence paragraphs are strictly prohibited; always use concise **bullet points with bold keywords (`- **Keyword**: Explanation`)**.
  * Restrict bullet points to a **maximum of 4 items** per slide.
* **Separation of Slide Content and Speaker Script**:
  * Keep the visual slide content strictly concise and high-signal.
  * Move all extended explanations, background context, and verbatim talking points into the **`<!-- note: ... -->` block** at the bottom of the slide.
* **Slide Boundary Rules**:
  * Slides must be separated by an isolated **`---` on its own line**.
  * Never insert `---` inside fenced code blocks (```` ``` ````).
* **Two-Column Split Requirement**:
  * When using `<!-- _layout: two-cols -->`, you **must insert `<!-- split -->` on its own line** between the left and right column contents.
  * Left column: typically 3–4 concise bullet points. Right column: supporting code snippet, diagram, or image.
* **Directive Scope Rule (`_`)**:
  * Directives prefixed with an underscore (`_`) apply **locally to that single slide only** (e.g., `_layout`, `_backgroundImage`, `_paginate`, `_fragmentStyle`).
  * Directives without an underscore are inherited by all subsequent slides.

---

## 2. Directive Catalog

### 2.1 Global Configuration (YAML Frontmatter)
```yaml
---
title: "Presentation Title"
author: "Author or Team Name"
theme: "clean"      # clean (minimalist white), default (tech standard), dark (dark mode)
size: "16:9"        # 16:9 (default) or 4:3
paginate: true      # Toggle slide page numbers (true/false)
header: "Header Text"  # Global top header
footer: "Footer Text"  # Global bottom footer
---
```

### 2.2 Five Semantic Layouts (`_layout`)
| Layout | Directive | Description & Required Elements |
| :--- | :--- | :--- |
| **Title Cover** | `<!-- _layout: cover -->` | Main deck cover. Vertically centered H1 Title + H3 Subtitle + Presenter info |
| **Section Divider** | `<!-- _layout: section -->` | Mid-deck chapter transition. H1 Title alone or with a short quote |
| **Two Columns** | `<!-- _layout: two-cols -->` | Parallel comparison. **Must include `<!-- split -->` between left and right columns** |
| **Blank / Canvas** | `<!-- _layout: blank -->` | Zero padding for full-bleed diagrams, images, or large math equations |
| **Default Body** | `<!-- _layout: default -->` | Standard content slide. H2 Header + 3–4 concise bullets / tables / code |

### 2.3 Visual Styling Directives
* **Background Image**: `<!-- _backgroundImage: url('images/bg.png') -->`
* **Background Dim (Darkening Overlay)**: `<!-- _backgroundDim: 0.5 -->` (0.1 to 0.9, ensures high text legibility over busy photos)
* **Background Color**: `<!-- _backgroundColor: #0f172a -->`
* **Text Color**: `<!-- _color: #ffffff -->`
* **Theme Inversion Class**: `<!-- _class: invert -->`
* **Local Visibility Overrides**: `<!-- _header: "" -->`, `<!-- _footer: "" -->`, `<!-- _paginate: false -->`

### 2.4 Speaker Notes
Placed at the bottom of the slide; visible only in presenter view (`P` key) and exported PPTX slide notes:
```markdown
<!-- note:
- Emphasize the latency benefits of item 2 in under 60 seconds.
- Remind attendees about the Q&A session at the end.
-->
```

### 2.5 Supported Markdown Extensions & Interactive Directives

#### 2.5.1 Incremental Reveal (Step-by-Step Animation)
Use `<!-- pause -->` between bullet points or paragraphs. Each item following `<!-- pause -->` remains hidden until the presenter advances the slide with Space or Arrow keys:
```markdown
- First visible point
<!-- pause -->
- Second point (revealed on next keypress)
<!-- pause -->
- Third point (revealed on subsequent keypress)
```
* **Dim Fragment Style**: Add `<!-- _fragmentStyle: dim -->` to pre-render unrevealed items at 25% opacity instead of fully hiding them, guiding audience anticipation.

#### 2.5.2 Code Blocks & Specific Line Highlighting
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

#### 2.5.3 Image Sizing & Centering
Control image dimensions and alignment directly within the alt-text attribute:
```markdown
![w:400 h:250](images/diagram.png)       <!-- Width 400px, height 250px -->
![w:60% center](images/architecture.png)  <!-- 60% viewport width, centered horizontally -->
![height:300px](images/screenshot.png)   <!-- Height 300px, auto aspect ratio -->
```

#### 2.5.4 Responsive YouTube Video Embed
Embed YouTube videos with responsive 16:9 aspect ratios by pasting standard watch or shortened URLs inside markdown image syntax:
```markdown
![](https://www.youtube.com/watch?v=dQw4w9WgXcQ)
![](https://youtu.be/dQw4w9WgXcQ?t=1m30s) <!-- Supports start timestamps -->
```

#### 2.5.5 GFM Alerts / Callout Boxes
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

#### 2.5.6 Mermaid Diagrams & LaTeX Math (KaTeX)
* **Mermaid Diagrams**: Native declarative flowcharts and sequence diagrams:
  ````markdown
  ```mermaid
  graph LR
      A[Markdown Source] --> B[Goslide Engine]
      B --> C[HTML]
      B --> D[PDF]
      B --> E[PPTX]
  ```
  ````
* **LaTeX Math**: Inline `$E=mc^2$` or display equation blocks `$$\sum_{i=1}^n x_i$$`.
* **Task Checklists**: `- [x] Completed milestone`, `- [ ] Pending item`.

---

## 3. Golden Example Deck

Below is a complete, production-ready 5-slide deck adhering to all rules and directives above. Replicate this exact structure, density, and formatting style when generating slides.

````markdown
---
title: "Modern Large Language Model Architecture"
author: "AI Core Architecture Lab"
theme: "clean"
size: "16:9"
paginate: true
header: "2026 AI Architecture Tech Seminar"
footer: "© 2026 Cloit Tech Architecture. All rights reserved."
---

<!-- _layout: cover -->
<!-- _paginate: false -->
<!-- _header: "" -->

# Modern Large Language Model Architecture
### From Transformer Core Mechanisms to Scalable Serving

**Speaker**: Cloit AI Architecture Lab  
**Date**: October 2026

<!-- note:
- Welcome attendees and introduce the session.
- Emphasize that this talk covers production-grade Go serving infrastructure beyond pure theory.
-->

---

<!-- _layout: section -->

# 01. Paradigm Shift: From RNN to Transformers
> *"General methods that leverage computation are ultimately the most effective."*

<!-- note:
- Section transition slide.
- Briefly remind the audience of Rich Sutton's Bitter Lesson and empirical scaling laws.
-->

---

<!-- _layout: two-cols -->

## Sequential vs Transformer Architecture

### Sequential Models (RNN / LSTM)
- **Sequential Bottleneck**: Hidden state dependency prevents GPU parallelization
- **Vanishing Context**: Gradient decay degrades long-range attention
- **Training Complexity**: $O(N)$ sequential steps bottleneck large corpora

<!-- split -->

### Transformer Decoder
- **Full Parallelization**: Simultaneous matrix operations across tensor cores
- **Self-Attention Mechanism**: $O(1)$ direct correlation path across all tokens
- **Predictable Scaling Law**: Capability grows reliably with compute budget

<!-- note:
- Guide audience through the left-to-right comparison.
- Highlight the parallelization advantage on the right column.
-->

---

<!-- _backgroundImage: url('assets/server-bg.png') -->
<!-- _backgroundDim: 0.7 -->
<!-- _color: white -->

## High-Throughput Go Streaming Proxy

- **Ultra-lightweight Goroutines**: Under 4KB baseline memory per active client
- **Native SSE Token Streaming**: Low-latency token delivery without heavy WebSocket overhead
- **Deterministic Single Binary**: Zero CGO dependencies for seamless cross-platform deployment

> [!TIP] Serving Optimization
> Leveraging Go's standard `net/http` and `http.Flusher` allows serving tens of thousands of concurrent streams with zero buffer bloat.

<!-- note:
- Showcase production architectural wins on this highlighted visual slide.
- Highlight the memory efficiency of Go over heavier alternatives.
-->

---

<!-- _fragmentStyle: dim -->

## Go Streaming Core Implementation

- **Low-Latency Streaming**: Immediate token flushing using standard `http.Flusher`
<!-- pause -->
- **Graceful Context Cancellation**: Client disconnects immediately release server resources

```go {3-5}
func StreamHandler(w http.ResponseWriter, r *http.Request) {
    flusher, _ := w.(http.Flusher)
    for token := range stream.Tokens() {
        fmt.Fprintf(w, "data: %s\n\n", token)
        flusher.Flush()
    }
}
```

<!-- note:
- Walk through lines 3 to 5 where the token flush occurs.
- Show how code line highlighting and incremental reveal direct audience attention.
-->
````
