# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [v0.1.1-test] - 2026-10-11

### Added
- **In-App Self-Update (`goslide update`)**:
  - Automatically check for the latest release via GitHub Releases API (`--check`).
  - Safe in-place atomic binary replacement across Linux, macOS, and Windows powered by `minio/selfupdate` (GOS-37, GOS-39).
  - SHA256 checksum verification before binary replacement.
  - Non-intrusive background update notification tip on CLI execution with 24-hour cache TTL.
- **Frontend i18n Support**:
  - Full Korean and English localization across Web presentation, annotation toolbar, and presenter console (GOS-34).
- **First-Time User Experience (FTUX)**:
  - Onboarding demo bundle (`goslide demo`) and template initialization (`goslide init`) packaged into release archives (GOS-31).
  - Multilingual quickstart guides (`QUICKSTART.md`, `QUICKSTART_en.md`).
- **Enhanced CLI & Server UX**:
  - Graceful Headless Chrome detection and troubleshooting guide on PDF/PPTX export failure (GOS-35, GOS-36).
  - Port binding conflict handling and auto-open browser fallback.

### Fixed
- Fixed HTML golden tests locale dependency by pinning test environment locale to ensure CI reproducibility.
- Resolved cyclomatic and cognitive complexity threshold violations across internal modules (GOS-32).

---

## [v0.1.0-test] - 2026-10-05

### Added
- **Core Markdown Presentation Engine**:
  - Pure Go static binary with CGO 0% and zero runtime OS dependencies.
  - Slide delimiters (`---`), frontmatter configuration, and directive parsing.
  - Syntax highlighting with Chroma and modular CSS theme system.
- **Multi-Format Export**:
  - Self-contained interactive HTML presentation slides with live reload (`goslide serve`).
  - High-resolution PDF export via headless browser integration.
  - PowerPoint (PPTX) export support.
- **Screencast & Annotation**:
  - Presenter console, interactive canvas drawing, and hotkey navigation.
- **Official Themes**:
  - Default, Academic, and Cyber-Dark themes with Mermaid diagram styling (GOS-29).
- **Cross-Platform Release Pipeline**:
  - Multi-arch automated packaging for Linux (`amd64`, `arm64`), macOS (`amd64`, `arm64`), and Windows (`amd64`, `arm64`) (GOS-33).
