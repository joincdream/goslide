# Goslide Official Documentation

Welcome to the Goslide Documentation Hub.  
This directory contains the official technical specifications for authoring Markdown presentations (DSL) and engineering custom presentation themes (CSS).

---

## 📚 Official Specifications (SSOT)

| Document | Purpose | Key Coverage |
|---|---|---|
| 📝 **[Slide Authoring & DSL Specification](dsl-guide.md)** | Content Authoring & Deck Generation | • 16:9 fixed viewport height invariants & content isolation<br>• YAML Frontmatter & 5 semantic layouts (`cover`, `two-cols`, etc.)<br>• Interactive reveals (`pause`), line highlighting, Mermaid & KaTeX<br>• Deterministic LLM output requirements |
| 🎨 **[Theme Development Specification](theme-guide.md)** | Styling & Visual Architecture | • 1920×1080 4-Tier Master DOM Architecture<br>• Strict positioning invariants (Never override `position: absolute`)<br>• Layout & component CSS selectors (code, tables, callouts, kbd)<br>• Production-ready Theme CSS Skeleton |

---

## 📂 Directory Role Taxonomy

Goslide enforces strict separation between documentation, executable examples, and theme assets:

* **`docs/` (Official Specifications)**:
  * [`docs/dsl-guide.md`](dsl-guide.md): Single Source of Truth for Slide Markdown DSL syntax.
  * [`docs/theme-guide.md`](theme-guide.md): Single Source of Truth for Presentation Theme CSS architecture.
* **`examples/` (Executable Presentation Decks)**:
  * [`examples/golden-sample/golden-deck.md`](../examples/golden-sample/golden-deck.md): Complete 5-slide production-grade presentation deck.
  * [`examples/demo/demo.md`](../examples/demo/demo.md): Interactive presentation feature & shortcut demo.
* **`themes/` (External CSS Themes)**:
  * [`themes/corporate.css`](../themes/corporate.css): Official enterprise business presentation theme.

---

## 💡 How to Prompt LLMs (Direct URL Referencing)

When prompting modern LLMs (ChatGPT, Claude, Cursor, Copilot), pass the document URL directly rather than copying and pasting text:

### 1. Generating a Presentation Slide Deck
```text
"Please refer to https://github.com/joincdream/goslide/blob/main/docs/dsl-guide.md 
 and generate a complete Goslide presentation deck for: [Your Topic / Outline / Document]."
```

### 2. Generating a Custom Theme CSS Stylesheet
```text
"Please refer to https://github.com/joincdream/goslide/blob/main/docs/theme-guide.md 
 and generate a complete Goslide theme CSS stylesheet for: [Your Brand Colors / Design Concept]."
```
