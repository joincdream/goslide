<div align="center">

# Goslide 📽️

**A blazing-fast, single-binary Markdown presentation deck builder written in Pure Go.**

Convert Markdown to interactive HTML, vector PDF, and editable PPTX with zero external dependencies.  
Tailored for developers, educators, and YouTube screencast lecture recordings.

[English](README.md) | [한국어](README.ko.md)

</div>

---

## ✨ Features

- **🚀 Single Binary & Pure Go**: Zero CGO runtime dependencies. Cross-platform executable (`linux`, `darwin`, `windows`) with embedded themes via `embed.FS`.
- **🎨 Multi-Format Export**:
  - **Interactive HTML**: Self-contained, responsive web slides with keyboard navigation.
  - **Vector PDF**: Crisp, print-ready 16:9 vector PDF generation powered by headless Chrome.
  - **High-Fidelity PPTX**: 100% pixel-perfect PowerPoint conversion preserving presenter speaker notes.
- **🖍️ Screencast & Live Annotation**: Built-in lightweight transparent canvas overlay designed for YouTube tutorial recording (`D` to toggle pen, `C` to clear).
- **🎯 Deterministic & Explicit Layouts**: Predictable layout rendering with explicit directives (`cover`, `two-cols`, `default`). Zero guesswork or accidental layout shifts.
- **⚡ Real-Time Live Preview**: Ultra-light local dev server (`goslide serve`) using standard Server-Sent Events (SSE) hot reload, keeping your current slide position.
- **🖥️ Dual-Screen Presenter Console**: Synchronized dual-window presenter view (`P`), digital laser pointer (`L`), and spotlight focus (`S`).
- **💻 Rich Content Ready**: Native syntax highlighting via Chroma, KaTeX math formulas (`$...$`, `$$...$$`), and GitHub Flavored Markdown (GFM) tables.

---

## 🚀 Quick Start

### Installation

```bash
# Using go install (Go 1.22+)
go install github.com/yundream/goslide/cmd/goslide@latest

# Or build from source
git clone https://github.com/yundream/goslide.git
cd goslide
make build
```

### Create Your First Slide

Create a markdown file `presentation.md`:

````markdown
---
title: "Modern Cloud Architecture"
theme: default
paginate: true
header: "Goslide Presentation"
footer: "© 2026 Cloud IT"
---

# Modern Cloud Architecture
### Scalable & Deterministic Systems

<!-- layout: cover -->

---

## 📌 Microservices Pipeline

- Go-based microservices with gRPC
- Real-time event streaming via NATS
- Kubernetes automated orchestration

<!-- layout: two-cols -->

<!-- split -->

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Goslide!")
}
```

<!-- note: Explain the decoupled architecture and performance benefits -->
````

### 🤖 Generating Slides with AI / LLMs

You can instantly generate production-ready Goslide presentations using ChatGPT, Claude, or any LLM:

> "Please create a slide deck about [Your Topic] in Markdown, referencing the specification at:  
> https://raw.githubusercontent.com/yundream/goslide/main/examples/example-dsl.md"

👉 For full semantic rules, directive catalogs, and golden samples, see the [Goslide DSL Specification & Example Guide](examples/example-dsl.md).

### Build Slides

```bash
# Generate standalone interactive HTML
goslide build presentation.md -o index.html

# Generate 16:9 Vector PDF
goslide build presentation.md -f pdf -o presentation.pdf

# Generate Editable PPTX
goslide build presentation.md -f pptx -o presentation.pptx
```

### Start Live Preview Server

```bash
# Start local development server with SSE hot-reload
goslide serve presentation.md --port 8080
```

Open `http://localhost:8080` in your browser. Whenever you save `presentation.md`, the browser automatically updates without losing your active slide index.

---

## ⌨️ Presentation & Screencast Shortcuts

| Key | Action | Context |
| :---: | :--- | :--- |
| `→` / `Space` / `PageDown` | Next slide | Presentation |
| `←` / `PageUp` | Previous slide | Presentation |
| `Home` / `End` | First / Last slide | Presentation |
| `[Number]` + `Enter` | Jump directly to slide number | Presentation |
| `F` | Toggle fullscreen mode | Presentation |
| **`D`** | **Toggle drawing pen on/off** | **Screencast / Annotation** |
| **`C`** | **Clear all canvas drawings** | **Screencast / Annotation** |
| `P` | Open dual-window presenter console | Presenter Tools |
| `L` | Toggle laser pointer | Presenter Tools |
| `S` | Toggle spotlight focus mode | Presenter Tools |

---

## 📐 Directives & Layouts

### Frontmatter (Global Directives)

Set presentation-level configuration at the top YAML block:

```yaml
---
title: My Presentation
theme: default        # default, clean, dark
paginate: true        # Show slide numbers (e.g. 3 / 24)
header: Header Text   # Optional top header
footer: Footer Text   # Optional bottom footer
size: 16:9            # 16:9 (default) or 4:3
---
```

### Slide Directives

Control individual slide behavior using HTML comments:

- `<!-- layout: cover -->`: Centered hero slide for titles and transitions.
- `<!-- layout: two-cols -->` (or `<!-- split -->`): Automatically splits contents into a 2-column layout.
- `<!-- backgroundColor: #1e1e1e -->`: Override background color for the current slide.
- `<!-- note: Speaker notes here -->`: Presenter speaker notes (synced to presenter console and PPTX).

---

## 🏗️ Architecture

Goslide follows a unidirectional, low-coupling pipeline:

```
[Markdown Source]
       │
       ▼ (1. Parse)
[Immutable Slide IR] (Deck, Slide, Directives)
       │
       ▼ (2. Transform & Resolve)
[Themed Presentation Model]
       │
       ├───────────────────┼───────────────────┐
       ▼ (3a. Render)      ▼ (3b. Export)      ▼ (3c. Export)
  HTML Renderer       PDF Exporter        PPTX Exporter
  (html/template)       (chromedp)        (OPC zip builder)
```

For detailed specifications and architectural documentation:
- [Core Architecture](docs/okf/core-architecture.md)
- [Interfaces & Contracts](docs/okf/contracts-interfaces.md)
- [Screencast Annotation Specification](docs/okf/screencast-annotation.md)
- [Hard Constraints & Guidelines](docs/okf/hard-constraints.md)

---

## 🛠️ Development

```bash
# Run unit & race condition tests
make test-race

# Run linter
make lint

# Run performance benchmark
make bench
```

---

## 📄 License

This project is licensed under the [Apache 2.0 License](LICENSE).
