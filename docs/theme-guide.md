# Goslide Theme CSS Specification

> 🤖 **LLM & AI Agent Directives**:  
> This document serves as the Single Source of Truth (SSOT) specification for LLMs and AI coding agents generating **Goslide Presentation Theme CSS stylesheets**.  
> When tasked with creating or customizing a theme, strictly adhere to the **4-Tier DOM Architecture** and **Hard Constraints** defined below to output deterministic, self-contained CSS.

---

## 1. 1920×1080 4-Tier Master DOM Architecture

Goslide operates on a strict, physical 1920×1080 (16:9) presentation canvas viewport, mirroring PowerPoint Slide Master isolation. The slide stage dynamically resizes via CSS `transform: scale(...)` in the browser.

```
┌────────────────────────────────────────────────────────┐  ▲
│ .slide-header / .slide-tracker (40px)                  │  │
├────────────────────────────────────────────────────────┤  │
│ .slide-title-box (120px) [default, two-cols only]      │  │
├────────────────────────────────────────────────────────┤  1080px
│                                                        │  │
│ .slide-body / .slide-content-box (840px ~ 960px)       │  │
│   (Markdown body, paragraphs, tables, code, media)     │  │
│                                                        │  │
├────────────────────────────────────────────────────────┤  │
│ .slide-footer (80px)                                   │  │
└────────────────────────────────────────────────────────┘  ▼
◄────────────────────── 1920px ──────────────────────────►
```

### 1.1 Complete Master HTML DOM Tree

```html
<section class="slide-card {{.Layout}} layout-{{.Layout}} {{.Classes}}">
  <!-- 0. Optional Background Dim Overlay -->
  <div class="slide-bg-dim"></div>

  <!-- 1. Running Header / Tracker (Height: 40px) -->
  <div class="slide-header slide-tracker">Category / Section / Deck Title</div>

  <!-- 2. Slide Title Placeholder Box (Height: 120px) -->
  <!-- Hidden via display: none in cover, section, and lead layouts -->
  <div class="slide-title-box">
    <h2 class="slide-title">Slide Title</h2>
    <p class="slide-subtitle">Slide Subtitle or Description</p>
  </div>

  <!-- 3. Slide Content / Body Flexbox Container (Height: 840px ~ 960px) -->
  <div class="slide-body slide-content-box">
    <!-- Rendered HTML Content from Markdown -->
  </div>

  <!-- 4. Slide Footer (Height: 80px) -->
  <div class="slide-footer">
    <span class="slide-author">Author / Org</span>
    <span class="slide-number">3 / 15</span>
  </div>
</section>
```

---

## 2. Hard Constraints for LLMs (Strict Rules)

Every generated theme CSS must comply with these architectural invariants:

### 🔴 Strict Don'ts
1. ❌ **NEVER apply `position: relative` to `.slide-card`, `.cover`, `.lead`, or `.has-bg-dim`.**  
   - Core presentation positioning relies on `.slide-card { position: absolute !important; top: 0; left: 0; }`.  
   - Overriding this to `position: relative` breaks the stack into normal document flow, causing the active slide to drop vertically off-screen to the bottom of the viewport.
2. ❌ **NEVER apply `justify-content: center` to `.slide-card` directly.**  
   - Applying flex centering to `.slide-card` displaces `.slide-header` and `.slide-footer` into the middle of the canvas.  
   - Vertical centering MUST always be targeted to `.slide-card.{layout} .slide-body`.
3. ❌ **NEVER apply fixed dimensions or scroll locks (`width: 1920px`, `height: 1080px`, `overflow: scroll`) to `html` or `body`.**  
   - Let the core stage runtime (`.goslide-stage`) manage responsive scaling.
4. ❌ **NEVER add blocking network font imports (`@import url(...)`).**  
   - Maintain offline portability and headless PDF/PPTX export reliability by using robust system font stacks (`-apple-system`, `BlinkMacSystemFont`, `Pretendard`, `Segoe UI`, `Roboto`, `sans-serif`).

### 🟢 Must-Haves
1. ✅ **Define theme colors and typography via `:root` CSS variables** for frictionless user customization.
2. ✅ **Explicitly override all 5 semantic layouts** (`cover`, `default`, `two-cols`, `section`, `lead`).
3. ✅ **Include complete component typography**: code blocks (`pre`, `code`), tables (`table`, `th`, `td`), callout alerts (`blockquote`), keyboard badges (`kbd`), and unordered/ordered lists.

---

## 3. Semantic Mapping & Architecture Matrix (YAML Specification)

The following machine-readable YAML specification establishes the definitive 1:1 binding between Markdown DSL syntax, HTML DOM targets, affected slide regions, and semantic design directives:

```yaml
semantic_mapping_matrix:
  - dsl_directive: "<!-- _layout: cover -->"
    target_selectors:
      card: ".slide-card.cover, .slide-card.layout-cover"
      body: ".slide-card.cover .slide-body, .slide-card.layout-cover .slide-body"
      hidden_elements: [".slide-title-box", ".slide-header", ".slide-footer"]
    affected_region: "Entire slide canvas (headers & footers suppressed)"
    semantic_purpose: "Deck Identity & Authority (First Impression)"
    design_guidelines:
      - "Establish presentation tone and visual authority immediately"
      - "Maximize contrast between the slide canvas and the Hero Title (H1)"
      - "Set H1 font-size to 3.5rem–4.0rem with margin-top: auto"
      - "Anchor presenter metadata cleanly to the floor with margin-top: auto"
    strict_invariants:
      - "Never declare position: relative on the card"
      - "Set .slide-body height to 1000px with box-sizing: border-box"

  - dsl_directive: "<!-- _layout: default -->"
    target_selectors:
      card: ".slide-card.default, .slide-card.layout-default"
      title_box: ".slide-card.default .slide-title-box"
      body: ".slide-card.default .slide-body"
    affected_region: "Standard 4-tier layout (120px title box + 840px body)"
    semantic_purpose: "Cognitive Clarity & Readability"
    design_guidelines:
      - "Optimized for scannable structured content and clear information hierarchy"
      - "Provide a distinct accent line on .slide-title-box::after to demarcate the topic"
      - "Maintain consistent 1.4rem–1.6rem body typography"
    strict_invariants:
      - "Title box must retain 120px height; body must retain 840px height"

  - dsl_directive: "<!-- _layout: two-cols -->"
    secondary_token: "<!-- split -->"
    target_selectors:
      container: ".slide-card.two-cols .two-cols, .slide-body .two-cols"
      columns: [".col-left", ".col-right"]
    affected_region: "Body internal flex split (Left and Right columns)"
    semantic_purpose: "Comparative Symmetry & Evidence Coupling"
    design_guidelines:
      - "Guide left-to-right cognitive flow: Concept/Theory on left, Code/Proof on right"
      - "Maintain a visual gutter of gap: 2rem or 48px between columns"
    strict_invariants:
      - "Both .col-left and .col-right MUST declare flex: 1; min-width: 0;"
      - "Never invent imaginary classes like .split or .two-columns"

  - dsl_directive: "<!-- _layout: section -->"
    target_selectors:
      card: ".slide-card.section, .slide-card.layout-section"
      body: ".slide-card.section .slide-body"
    affected_region: "Mid-deck chapter transition (Centered H1, title box hidden)"
    semantic_purpose: "Cognitive Reset & Narrative Punctuation"
    design_guidelines:
      - "Signal clear narrative transition between presentation chapters"
      - "Use distinct background contrast or a thick primary accent border (8–12px)"
      - "Render H1 at 3.0rem–3.5rem with high visual weight"

  - dsl_directive: "<!-- _layout: lead -->"
    target_selectors:
      card: ".slide-card.lead, .slide-card.layout-lead"
      body: ".slide-card.lead .slide-body, .slide-card.layout-lead .slide-body"
      hidden_elements: [".slide-title-box"]
    affected_region: "Vertically and horizontally centered body container (960px)"
    semantic_purpose: "Singular Punchline & Hands-on Transition"
    design_guidelines:
      - "Deliver an unmissable single takeaway, quote, or hands-on workshop transition"
      - "Center both H1 (3.5rem) and supporting H3/quote (1.6rem)"
    strict_invariants:
      - "Center alignment MUST be declared on .slide-body, NEVER on .slide-card"
      - "Declare display: flex !important; justify-content: center !important; align-items: center !important; on .slide-body"

  - dsl_directive: "<!-- _layout: blank -->"
    target_selectors:
      card: ".slide-card.blank, .slide-card.layout-blank"
      body: ".slide-card.blank .slide-body"
      hidden_elements: [".slide-title-box", ".slide-header", ".slide-footer"]
    affected_region: "Full-bleed canvas (Padding 0, headers suppressed)"
    semantic_purpose: "Full-Bleed Visual Immersion"
    design_guidelines:
      - "Unobstructed full-screen canvas for complex system architecture diagrams and math"
    strict_invariants:
      - "Padding on .slide-card and .slide-body must be 0 !important"

  - dsl_directive: "<!-- _backgroundImage --> & <!-- _backgroundDim -->"
    target_selectors:
      overlay: ".slide-bg-dim"
      active_card: ".slide-card.has-bg-dim"
    affected_region: "Z-index 1 overlay between background image and body content"
    semantic_purpose: "Legibility Protection & Contrast Assurance"
    design_guidelines:
      - "Guarantees minimum 4.5:1 WCAG contrast for white text over busy photography"
    strict_invariants:
      - "NEVER declare position: relative on .slide-card.has-bg-dim"
      - ".slide-bg-dim uses position: absolute; inset: 0; pointer-events: none;"

  - dsl_directive: "<!-- pause -->"
    target_selectors:
      item: ".slide-card .fragment"
      revealed_item: ".slide-card .fragment.visible"
      dimmed_item: ".slide-card.has-fragment-dim .fragment:not(.visible)"
    affected_region: "Individual list items, paragraphs, code blocks, or table rows"
    semantic_purpose: "Paced Attention & Anti-Spoiler Storytelling"
    design_guidelines:
      - "Keep the audience synchronized with the speaker's vocal progression"
      - "Use opacity: 0; transition: opacity 0.25s; for clean reveals"

  - dsl_directive: "Fenced code blocks with line highlights (e.g. ```go {2,4-6})"
    target_selectors:
      block: ".slide-card pre"
      inline: ".slide-card :not(pre) > code"
      highlighted_line: ".chroma .hl"
    affected_region: "Code display elements"
    semantic_purpose: "Engineering Focus & Walkthrough"
    design_guidelines:
      - "Subdue non-highlighted lines to guide eye to active execution logic"
      - "Use clean monospaced font with line-height: 1.45–1.55"

  - dsl_directive: "GFM Callouts (> [!NOTE], > [!TIP], > [!WARNING])"
    target_selectors:
      callout: ".slide-card .markdown-alert, .slide-card blockquote"
      types:
        note: ".markdown-alert-note"
        tip: ".markdown-alert-tip"
        warning: ".markdown-alert-warning"
        important: ".markdown-alert-important"
    affected_region: "Callout block elements"
    semantic_purpose: "Epistemic Status & Best Practice Highlighting"
    design_guidelines:
      - "Set distinct left border (4px) with semantic hues: blue (note), emerald (tip), amber (warning)"
```

---



## 4. Component Selectors Specification

```css
/* Code Blocks (Chroma Syntax Highlighter Integration) */
.slide-card pre {
  padding: 1rem 1.25rem;
  border-radius: 8px;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.95rem;
  line-height: 1.5;
  overflow-x: auto;
}

/* Inline Code */
.slide-card :not(pre) > code {
  padding: 0.15rem 0.4rem;
  border-radius: 4px;
  font-size: 0.88em;
}

/* Tables (Clean Striped Layout) */
.slide-card table {
  width: 100%;
  border-collapse: collapse;
  margin: 1rem 0;
  font-size: 0.95em;
}
.slide-card th {
  font-weight: 600;
  text-align: left;
  padding: 0.75rem 1rem;
}
.slide-card td {
  padding: 0.65rem 1rem;
}

/* GFM Alerts & Blockquotes */
.slide-card blockquote {
  border-left: 4px solid var(--accent-color);
  padding: 0.75rem 1.25rem;
  margin: 1rem 0;
  border-radius: 0 8px 8px 0;
}

/* Keyboard Shortcut Badges */
.slide-card kbd {
  display: inline-block;
  padding: 2px 8px;
  font-family: ui-monospace, monospace;
  font-size: 0.82em;
  border-radius: 4px;
  box-shadow: 0 2px 0 rgba(0, 0, 0, 0.15);
}

/* Lists */
.slide-card ul, .slide-card ol {
  margin: 1rem 0;
  padding-left: 1.5rem;
  line-height: 1.6;
}
.slide-card li {
  margin-bottom: 0.5rem;
}

/* Incremental Reveal (Fragment Animations) */
.slide-card .fragment {
  opacity: 0;
  transition: opacity 0.25s ease-in-out;
}
.slide-card .fragment.visible {
  opacity: 1;
}
.slide-card.has-fragment-dim .fragment:not(.visible) {
  opacity: 0.25; /* Pre-render unrevealed items at 25% opacity */
}

/* Density Modifiers */
.slide-card.compact .slide-body {
  font-size: 0.88rem;
}
.slide-card.dense .slide-body {
  font-size: 0.78rem;
}
.slide-card.compact .slide-body pre,
.slide-card.dense .slide-body pre {
  padding: 0.6rem 0.8rem;
  font-size: 0.82em;
}
```

---

## 5. Standard Theme CSS Skeleton (Reference Baseline)

LLMs must use this skeleton as the baseline architecture when authoring new themes:

```css
/* ==========================================================================
   Goslide Theme: [Theme Name]
   ========================================================================== */

:root {
  /* Brand Accents & Palette */
  --primary-color: #0f172a;
  --secondary-color: #1e293b;
  --accent-color: #2563eb;
  --accent-hover: #1d4ed8;

  /* Surfaces & Typography */
  --bg-card: #ffffff;
  --bg-stage: #f1f5f9;
  --bg-subtle: #f8fafc;
  --text-main: #1e293b;
  --text-muted: #64748b;
  --border-color: #e2e8f0;

  /* Typography Stack */
  --font-family: -apple-system, BlinkMacSystemFont, "Pretendard", "Segoe UI", Roboto, sans-serif;
  --font-code: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}

/* Base Viewport & Card Canvas */
body {
  font-family: var(--font-family);
  background-color: var(--bg-stage);
  color: var(--text-main);
}

.slide-card {
  background-color: var(--bg-card);
  color: var(--text-main);
  border: 1px solid var(--border-color);
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.08);
  /* INVARIANT: Never add position: relative here */
}

/* 1. Cover Layout */
.slide-card.cover,
.slide-card.layout-cover {
  /* Custom background gradients or dark modes */
  color: #f8fafc;
  overflow: hidden;
}
.slide-card.cover .slide-title-box,
.slide-card.cover .slide-header,
.slide-card.cover .slide-footer,
.slide-card.layout-cover .slide-title-box,
.slide-card.layout-cover .slide-header,
.slide-card.layout-cover .slide-footer {
  display: none !important;
}
.slide-card.cover .slide-body,
.slide-card.layout-cover .slide-body {
  display: flex !important;
  flex-direction: column !important;
  justify-content: flex-start !important;
  flex: 0 0 1000px !important;
  height: 1000px !important;
  box-sizing: border-box !important;
  padding: 100px 140px 60px 140px !important;
  text-align: center !important;
}
.slide-card.cover h1,
.slide-card.layout-cover h1 {
  margin-top: auto !important;
  margin-bottom: 0.5rem !important;
  font-size: 3.75rem !important;
  font-weight: 800 !important;
  line-height: 1.2 !important;
}
.slide-card.cover p,
.slide-card.layout-cover p {
  margin-top: auto !important;
  margin-bottom: 0 !important;
  text-align: left !important;
}

/* 2. Lead Layout */
.slide-card.lead .slide-title-box,
.slide-card.layout-lead .slide-title-box {
  display: none !important;
}
.slide-card.lead .slide-body,
.slide-card.layout-lead .slide-body {
  display: flex !important;
  flex-direction: column !important;
  justify-content: center !important;
  align-items: center !important;
  text-align: center !important;
  flex: 0 0 960px !important;
  height: 960px !important;
  box-sizing: border-box !important;
  padding: 60px 140px !important;
}
.slide-card.lead h1,
.slide-card.layout-lead h1 {
  font-size: 3.5rem !important;
  font-weight: 800 !important;
  margin-bottom: 1.25rem !important;
  color: inherit !important;
}

/* 3. Two Columns Layout */
.slide-card .two-cols {
  display: flex;
  gap: 2rem;
  width: 100%;
}
.slide-card .two-cols .col-left,
.slide-card .two-cols .col-right {
  flex: 1;
  min-width: 0;
}
```

---

## 6. Output Requirements for LLMs

When generating theme CSS from user requests:
1. **Single Self-Contained CSS Block**: Output the entire CSS inside a single ```` ```css ```` code fence without omitted placeholders.
2. **Variable-Driven Customization**: Expose all branding colors and typography under `:root`.
3. **Strict Zero-Displacement**: Never introduce `position: relative` or margins that displace `.slide-card` from its fixed (0, 0) canvas coordinate.
