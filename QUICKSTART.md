<div align="center">

# 🚀 Goslide 3-Minute Quickstart Guide

[English](QUICKSTART.md) | [한국어](QUICKSTART.ko.md)

</div>

---

Welcome to **Goslide**! Goslide is a blazing-fast, single-binary Markdown presentation deck builder written in Pure Go.  
Follow these 5 simple steps to create your first presentation and export to PDF/PPTX in **under 3 minutes**.

---

## Step 1. Unpack Demo Environment (30 sec)

Run the following command in your terminal to unpack the showcase presentation (`demo/`) and project themes (`themes/`):

```bash
# Unpack demo presentation, image assets, and project themes
goslide demo
```

Generated directory structure:
- `demo/`: 14-slide rich showcase presentation (`demo.md`) and demo image assets
- `themes/`: Shared custom stylesheets (`clean.css`, `dark.css`, `academic.css`) for project-wide use

---

## Step 2. Live Preview & Hot Reload (60 sec)

```bash
# Start local live development server (opens default browser automatically)
goslide serve demo/demo.md
```

- **Live Hot Reload**: Open `demo/demo.md` in your text editor, edit content, and save (`Ctrl+S`). The browser refreshes seamlessly in under 300ms without losing your slide position!
- **Switch Themes**: Change `theme: "clean"` in the frontmatter of `demo/demo.md` to `"dark"` or `"academic"`, or customize the CSS files in `themes/`.
- **Presenter Shortcuts**:
  - `P`: Pop-out dual-screen presenter console (speaker notes, timer, next slide preview)
  - `D`: Toggle transparent screencast drawing pen (`C` to clear canvas)
  - `L`: Digital laser pointer mode
  - `S`: Spotlight mouse focus mode
  - `B` / `W`: Blackout / Whiteout screen curtain
  - `?`: Full keyboard shortcuts help modal

---

## Step 3. Multi-Format Batch Export (30 sec)

Export your presentation to interactive web slides, vector PDF, and PowerPoint all in one command:

```bash
# Batch build HTML, vector PDF, and PPTX simultaneously
goslide build demo/demo.md -f html,pdf,pptx
```

Output files:
- `demo/demo.html`: Offline, standalone interactive web slides
- `demo/demo.pdf`: Crisp, print-ready 16:9 borderless vector PDF
- `demo/demo.pptx`: Native PowerPoint presentation preserving speaker notes

---

## Step 4. Create Your Own Presentation (60 sec)

To start writing your own slide deck from scratch, run `init`:

```bash
# Initialize a new starter presentation
goslide init my-slide.md

# Start live preview
goslide serve my-slide.md
```

---

## Step 5. AI-Assisted Slide Generation (ChatGPT / Claude)

Want an AI agent to write your slides? Copy the dedicated Goslide system prompt directly from the CLI:

```bash
# Print AI system prompt to terminal
goslide prompt
```

Paste the output prompt into ChatGPT or Claude and describe your topic (e.g. *"Create a 7-slide technical talk on Go concurrency and channels"*). The AI will generate production-ready Goslide Markdown DSL instantly.

---

## 💡 Useful Commands

- `goslide theme list`: List all available built-in themes
- `goslide theme export clean my-style.css`: Export built-in theme CSS for custom branding
- `goslide --help`: Display full CLI documentation and options
