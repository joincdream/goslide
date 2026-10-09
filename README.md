<div align="center">

# Goslide 📽️

**A blazing-fast, single-binary Markdown presentation deck builder written in Pure Go.**

Convert Markdown documents into interactive HTML, vector PDF, and editable PPTX with zero external dependencies.  
Tailored for developers, educators, and YouTube screencast lecture recordings.

[English](README.md) | [한국어](README.ko.md) · [3-Min Quickstart](QUICKSTART.md)

</div>

---

## ✨ Features

- **🚀 Single Binary & Pure Go**: Zero CGO runtime dependencies. Cross-platform executable (`Linux`, `macOS`, `Windows`) with embedded themes via `embed.FS`.
- **🎨 Multi-Format Export**:
  - **Interactive HTML**: Self-contained, responsive web slides with offline support.
  - **Vector PDF**: Crisp, print-ready 16:9 borderless vector PDF powered by headless Chrome.
  - **High-Fidelity PPTX**: Pixel-perfect PowerPoint conversion preserving presenter speaker notes.
- **🖍️ Screencast & Live Annotation**: Built-in lightweight transparent canvas overlay designed for YouTube tutorial recording (`D` to toggle pen, `C` to clear).
- **🎯 Deterministic & Explicit Layouts**: Predictable layout rendering with explicit directives (`cover`, `two-cols`, `section`, `blank`, `default`). Zero guesswork or accidental layout shifts.
- **⚡ Real-Time Live Preview**: Ultra-light local dev server (`goslide serve`) using standard Server-Sent Events (SSE) hot reload, keeping your current slide position.
- **🖥️ Dual-Screen Presenter Console**: Synchronized dual-window presenter view (`P`), digital laser pointer (`L`), and spotlight focus (`S`).
- **💻 Rich Content Ready**: Native syntax highlighting via Chroma, KaTeX math formulas (`$...$`, `$$...$$`), and GitHub Flavored Markdown (GFM) tables.

---

## 🚀 Quick Start

> 💡 For a complete step-by-step walkthrough, see the **[3-Minute Quickstart Guide](QUICKSTART.md)**.

### 1. Installation

#### Download Binary (Recommended)
Download the pre-compiled archive for your OS (Linux, macOS, Windows) from [GitHub Releases](https://github.com/joincdream/goslide/releases) and extract the binary.

#### Install via Go (Go 1.22+)
```bash
go install github.com/joincdream/goslide/cmd/goslide@latest
```

#### Build from Source
```bash
git clone https://github.com/joincdream/goslide.git
cd goslide
make build
```

---

### 2. 3-Minute Interactive Demo & Live Reload

Goslide embeds a rich 14-slide showcase presentation and customizable themes directly inside the binary:

```bash
# 1. Unpack demo presentation (demo/) and project themes (themes/)
goslide demo

# 2. Launch live preview development server (opens browser automatically)
goslide serve demo/demo.md

# 3. Batch export to HTML, vector PDF, and PowerPoint PPTX simultaneously
goslide build demo/demo.md -f html,pdf,pptx

# 4. Initialize your own presentation
goslide init my-slide.md
```

---

### 3. AI / LLM Assisted Slide & Theme Creation

Leverage AI agents (ChatGPT, Claude, Cursor) to generate slide decks or design custom themes using official specifications:

#### 1) Slide Deck Generation
```bash
# Print AI prompt for slide deck generation (includes official DSL specification URL)
goslide prompt
```
Paste the output into ChatGPT or Claude with your topic (e.g. *"Create an 8-slide presentation on cloud-native microservices architecture"*). The AI will generate production-ready Goslide Markdown DSL instantly.  
👉 [Official Goslide Markdown DSL Specification (docs/dsl-guide.md)](docs/dsl-guide.md)

#### 2) Custom Theme CSS Generation
```bash
# Print AI prompt for custom CSS theme generation (includes official 4-Tier DOM architecture URL)
goslide prompt theme
```
Paste the output into your AI tool with your branding requirements (e.g. *"Modern corporate theme with navy blue accent, Inter typography, and light gray cards"*). The AI will craft a self-contained CSS stylesheet tailored for 1920×1080 fixed viewports.  
👉 [Official Goslide Theme CSS Architecture Specification (docs/theme-guide.md)](docs/theme-guide.md)

---

## ⌨️ Keyboard Shortcuts (Presentation & Screencast)

| Key | Description | Context |
| :---: | :--- | :--- |
| `→` / `Space` / `PageDown` / `J` | Next slide / fragment | Presentation |
| `←` / `PageUp` / `K` / `H` | Previous slide / fragment | Presentation |
| `Home` / `End` | Go to first / last slide | Presentation |
| `[Number]` + `Enter` | Jump directly to slide number | Presentation |
| `F` | Toggle fullscreen mode | Presentation |
| **`D`** | **Toggle transparent drawing pen** | **YouTube Recording / Live Screencast** |
| **`C`** | **Clear all drawings on active slide** | **YouTube Recording / Live Screencast** |
| `L` | Toggle digital laser pointer | Presenter Mode |
| `S` | Toggle spotlight mouse focus | Presenter Mode |
| `P` | Open pop-out dual-screen presenter console | Presenter Mode |
| `B` / `W` | Blackout / Whiteout screen curtain | Audience Engagement |
| `1` ~ `4` | Change pen / laser pointer color (Red, Blue, Green, Yellow) | Annotation / Pointer |
| `+` / `-` | Change pen / laser stroke width (Thin, Medium, Thick) | Annotation / Pointer |
| `?` | Open keyboard shortcuts help popup | Help |

---

## 📐 Directives & Layout Usage

### Frontmatter (Global Directives)

Declare presentation metadata in the YAML block at the top of the file:

```yaml
---
title: Presentation Title
theme: clean          # Themes: clean, dark, academic, cyber-dark, default
size: 16:9            # Aspect ratio: 16:9 (default) or 4:3
paginate: true        # Show slide numbers (e.g. 3 / 24)
header: Header Text
footer: Footer Text
autofit: true         # Prevent text overflow with dynamic font scaling
---
```

### Scoped Directives (Slide Level)

Control per-slide layout and styling using HTML comments:

- `<!-- layout: cover -->`: Centered layout for title and divider slides.
- `<!-- layout: two-cols -->` (or `<!-- split -->`): 2-column grid layout.
- `<!-- layout: section -->`: Section header layout.
- `<!-- layout: blank -->`: Zero-padding full canvas layout (ideal for diagrams/full-bleed images).
- `<!-- backgroundColor: #1e1e1e -->`: Override background color for active slide.
- `<!-- backgroundImage: url(bg.jpg) -->`: Apply background image.
- `<!-- note: Speaker notes content -->`: Speaker script recorded in presenter console and PPTX notes.

---

## 🎨 Theme Management & Custom Themes

Goslide is 100% CSS-driven, allowing you to craft custom branded themes with complete layout control.

### 1. Inspect & Export Built-in Themes
```bash
# List available built-in themes (clean, dark, academic, cyber-dark, default)
goslide theme list

# Export built-in theme CSS to local themes/ directory
goslide theme export clean themes/my-company.css
```

### 2. Automatic Custom Theme Discovery
Save your CSS file in the project's `themes/` folder (e.g., `themes/my-company.css`). Goslide discovers it automatically by name in Frontmatter:
```yaml
---
title: "Quarterly Business Review"
theme: my-company      # Automatically discovers themes/my-company.css
---
```

### 3. Master CSS Selector Architecture
| Selector | Purpose | Default Height / Properties |
| :--- | :--- | :--- |
| `.slide-card` | Slide viewport container | Background color, base font color, borders, shadow |
| `.slide-card .slide-header` | Top fixed tracker text | 40px |
| `.slide-card .slide-title-box` | Slide title & subtitle container | 120px |
| `.slide-card .slide-body` | Slide content flexbox container | 840px - 960px |
| `.slide-card .slide-footer` | Bottom footer & slide pagination | 80px |
| `.slide-card.layout-cover` | Cover title layout override | Centered H1 |
| `.slide-card.layout-two-cols` | 2-column grid layout | `.col-left`, `.col-right` |

👉 For full 4-tier DOM architecture and strict layout guidelines, see the **[Goslide Theme CSS Specification (docs/theme-guide.md)](docs/theme-guide.md)**.

---

## 🏗️ System Architecture

Goslide strictly adheres to a unidirectional, decoupled pipeline architecture:

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

Detailed specifications:
- [Core Architecture Specification](docs/okf/core-architecture.md)
- [Interfaces & Domain Model Contracts](docs/okf/contracts-interfaces.md)
- [Screencast Annotation Specification](docs/okf/screencast-annotation.md)
- [Hard Constraints & DoD Checklist](docs/okf/hard-constraints.md)

---

## 🛠️ Development & Contributing

```bash
# Run all quality checks (formatting, linter, complexity, race tests)
make check

# Run unit & race condition tests
make test-race

# Run linter
make lint

# Run performance benchmarks
make bench
```

---

## 👤 Author & Contact

**yundream (Joinc)**

- 🌐 **Website**: [https://www.joinc.co.kr](https://www.joinc.co.kr)
- 💼 **LinkedIn**: [linkedin.com/in/yundream](https://www.linkedin.com/in/yundream/)
- ✉️ **Email**: [yundream@gmail.com](mailto:yundream@gmail.com)
- 🐙 **GitHub**: [@joincdream](https://github.com/joincdream)

---

## 📄 License

This project is licensed under the [Apache 2.0 License](LICENSE).
