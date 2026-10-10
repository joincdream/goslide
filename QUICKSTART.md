<div align="center">

# Goslide 3-Minute Quickstart Guide

[English](QUICKSTART.md) | [한국어](QUICKSTART.ko.md)

</div>

---

Welcome to **Goslide**! Goslide is a blazing-fast, single-binary Markdown presentation deck builder written in Pure Go.  
Follow these steps to create your first presentation and export to PDF/PPTX in **under 3 minutes**.

---

## Step 0. Download & Install Binary (30 sec)

Goslide is a single static binary with zero CGO or external runtime dependencies. Download the appropriate pre-compiled artifact for your system from [GitHub Releases](https://github.com/joincdream/goslide/releases):

| Operating System | CPU Architecture / Target Device | Recommended Artifact |
| :--- | :--- | :--- |
| **macOS** | **Apple Silicon (M1 / M2 / M3 / M4)** | `goslide_darwin_arm64.tar.gz` |
| **macOS** | Intel Processor | `goslide_darwin_amd64.tar.gz` |
| **Linux** | 64-bit Intel / AMD (x86_64) | `goslide_linux_amd64.tar.gz` |
| **Linux** | 64-bit ARM (aarch64 / Raspberry Pi 4, 5) | `goslide_linux_arm64.tar.gz` |
| **Windows** | 64-bit Intel / AMD (x64) | `goslide_windows_amd64.zip` |
| **Windows** | 64-bit ARM (Snapdragon X Elite, etc.) | `goslide_windows_arm64.zip` |

### Linux & macOS One-Liner Setup
```bash
# 1. Extract the downloaded archive (e.g., Linux amd64)
tar -xzf goslide_linux_amd64.tar.gz

# 2. Make executable and install into user PATH
chmod +x goslide
mkdir -p ~/.local/bin && cp goslide ~/.local/bin/

# 3. Verify installation
goslide --help
```

> Note (macOS users): If blocked by macOS Gatekeeper ("unidentified developer"), run:  
> `xattr -d com.apple.quarantine goslide`

### Windows Setup
- Extract the downloaded `.zip` file, move `goslide.exe` to a preferred directory (e.g., `C:\bin`), and add that path to your system's `Path` environment variable.

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

## Useful Commands & Tips

- `goslide theme list`: List all available built-in themes
- `goslide theme export clean my-style.css`: Export built-in theme CSS for custom branding
- `goslide --help`: Display full CLI documentation and options

### Verify Download Integrity (SHA-256)
```bash
# Linux
sha256sum -c checksums.txt --ignore-missing

# macOS
shasum -a 256 -c checksums.txt --ignore-missing
```
