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
* **Directive Scope Rule (`_`)**:
  * Directives prefixed with an underscore (`_`) apply **locally to that single slide only** (e.g., `_layout`, `_backgroundImage`, `_paginate`).
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
* **Background Dim (Darkening Overlay)**: `<!-- _backgroundDim: 0.5 -->` (0.1 to 0.9, ensures high text legibility)
* **Background Color**: `<!-- _backgroundColor: #0f172a -->`
* **Text Color**: `<!-- _color: #ffffff -->`
* **Theme Inversion Class**: `<!-- _class: invert -->`
* **Local Visibility Overrides**: `<!-- _header: "" -->`, `<!-- _footer: "" -->`, `<!-- _paginate: false -->`

### 2.4 Speaker Notes
Placed at the bottom of the slide; visible only in presenter view (`P` key) and exported PPTX slide notes:
```markdown
<!-- note:
- Emphasize the latency benefits of item 2 in under 60 seconds.
-->
```

### 2.5 Supported Markdown Extensions
* **Fenced Code Blocks & Line Highlighting**: ````go {2-4,7}```` (Native Chroma syntax highlighting with line focus; dims unselected lines)
* **Incremental Reveal (Fragments)**: `<!-- pause -->` (Steps through bullets or paragraphs one-by-one with space/arrow keys)
* **Fragment Dim Style**: `<!-- _fragmentStyle: dim -->` (Pre-renders unrevealed bullets at 25% opacity before focus)
* **LaTeX Math (KaTeX)**: Inline `$E=mc^2$` or block `$$\text{Attention}(Q, K, V) = \text{softmax}\left(\frac{QK^T}{\sqrt{d_k}}\right)V$$`
* **GFM Callouts**: `> [!NOTE]`, `> [!TIP]`, `> [!WARNING]`, `> [!IMPORTANT]`
* **Task Checklists**: `- [x] Completed task`, `- [ ] Planned task`

---

## 3. Golden Example Deck

Below is a complete, production-ready 4-slide deck adhering to all rules and directives above. Replicate this exact structure, density, and formatting style when generating slides.

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
