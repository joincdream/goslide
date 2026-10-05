package pptx

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/yundream/goslide/internal/browser"
	"github.com/yundream/goslide/internal/model"
	"github.com/yundream/goslide/internal/parser"
)

func TestBuildContentTypesXML(t *testing.T) {
	xml := string(buildContentTypesXML(2, []bool{true, false}))
	if !strings.Contains(xml, `PartName="/ppt/slides/slide1.xml"`) {
		t.Errorf("missing slide1 override in ContentTypes")
	}
	if !strings.Contains(xml, `PartName="/ppt/slides/slide2.xml"`) {
		t.Errorf("missing slide2 override in ContentTypes")
	}
	if !strings.Contains(xml, `PartName="/ppt/notesSlides/notesSlide1.xml"`) {
		t.Errorf("missing notesSlide1 override in ContentTypes")
	}
	if strings.Contains(xml, `PartName="/ppt/notesSlides/notesSlide2.xml"`) {
		t.Errorf("unexpected notesSlide2 override in ContentTypes")
	}
}

func TestBuildPresentationXML(t *testing.T) {
	xml := string(buildPresentationXML(2))
	if !strings.Contains(xml, `type="screen16x9"`) {
		t.Errorf("missing screen16x9 type in presentation.xml")
	}
	if !strings.Contains(xml, `cx="12192000" cy="6858000"`) {
		t.Errorf("missing 16:9 widescreen dimensions in presentation.xml")
	}
}

func TestBuildSlideXML(t *testing.T) {
	xml := string(buildSlideXML(nil))
	if !strings.Contains(xml, `<a:masterClrMapping/>`) {
		t.Errorf("missing masterClrMapping in slide.xml")
	}
	if !strings.Contains(xml, `cx="12192000" cy="6858000"`) {
		t.Errorf("missing 16:9 fullscreen image extent in slide.xml")
	}

	links := []SlideLink{
		{TargetURL: "https://youtube.com/watch?v=123", X: 100, Y: 200, W: 300, H: 400},
	}
	xmlWithLinks := string(buildSlideXML(links))
	if !strings.Contains(xmlWithLinks, `<a:hlinkClick r:id="rIdLink1"/>`) {
		t.Errorf("missing hyperlink click in slide.xml with links")
	}

	rels := string(buildSlideRelsXML(1, false, links))
	if !strings.Contains(rels, `Target="https://youtube.com/watch?v=123"`) {
		t.Errorf("missing hyperlink in slide rels")
	}
}

func TestBuildNotesSlideXML(t *testing.T) {
	notes := "Key takeaway: <Alpha> & \"Beta\""
	xml := string(buildNotesSlideXML(notes))
	if !strings.Contains(xml, `Key takeaway: &lt;Alpha&gt; &amp; &#34;Beta&#34;`) {
		t.Errorf("notes text was not properly XML escaped, got: %s", xml)
	}
	if !strings.Contains(xml, `<a:masterClrMapping/>`) {
		t.Errorf("missing masterClrMapping in notes.xml")
	}
}

func TestPPTXExporter_NilOrEmptyDeck(t *testing.T) {
	e := NewExporter()
	if err := e.Export(context.Background(), nil, "output.pptx"); !errors.Is(err, model.ErrExportFailed) {
		t.Errorf("expected ErrExportFailed for nil deck, got: %v", err)
	}

	emptyDeck := &model.Deck{}
	if err := e.Export(context.Background(), emptyDeck, "output.pptx"); !errors.Is(err, model.ErrExportFailed) {
		t.Errorf("expected ErrExportFailed for empty deck, got: %v", err)
	}
}

const samplePPTXMarkdown = `---
title: "PPTX Export Test"
theme: "clean"
---

<!-- _layout: cover -->
# First Slide: Cover

Welcome to PPTX Exporter

<!-- note: Speaker notes for slide 1.
Remember to introduce the project vision. -->

---

<!-- _layout: section -->
## Second Slide: No Notes

This slide does not have presenter notes.
`

func TestPPTXExporter_Export(t *testing.T) {
	browserPath, err := browser.FindChrome("")
	if err != nil {
		t.Skipf("skipping PPTX export test: chrome not available: %v", err)
	}

	p := parser.NewParser()
	deck, err := p.Parse(context.Background(), strings.NewReader(samplePPTXMarkdown))
	if err != nil {
		t.Fatalf("failed to parse sample markdown: %v", err)
	}

	outputPath := filepath.Join(t.TempDir(), "presentation.pptx")
	exporter := NewExporter(
		WithTheme("clean"),
		WithBrowserPath(browserPath),
		WithTimeout(30*time.Second),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	if err := exporter.Export(ctx, deck, outputPath); err != nil {
		t.Fatalf("Export() failed: %v", err)
	}

	verifyZipHeader(t, outputPath)
	entries := openZipEntries(t, outputPath)
	verifyExpectedZipEntries(t, entries)
	verifyNotesSlideContent(t, entries)
	verifySlideImageHeader(t, entries)
}

func verifyZipHeader(t *testing.T, path string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if !bytes.HasPrefix(data, []byte("PK\x03\x04")) {
		t.Errorf("file does not start with ZIP magic header")
	}
}

func openZipEntries(t *testing.T, path string) map[string]*zip.File {
	t.Helper()
	zr, err := zip.OpenReader(path)
	if err != nil {
		t.Fatalf("failed to open zip file: %v", err)
	}
	t.Cleanup(func() { _ = zr.Close() })

	entries := make(map[string]*zip.File, len(zr.File))
	for _, f := range zr.File {
		entries[f.Name] = f
	}
	return entries
}

func verifyExpectedZipEntries(t *testing.T, entries map[string]*zip.File) {
	t.Helper()
	expectedFiles := []string{
		"[Content_Types].xml",
		"_rels/.rels",
		"ppt/presentation.xml",
		"ppt/_rels/presentation.xml.rels",
		"ppt/slides/slide1.xml",
		"ppt/slides/_rels/slide1.xml.rels",
		"ppt/media/slide1.png",
		"ppt/notesSlides/notesSlide1.xml",
		"ppt/notesSlides/_rels/notesSlide1.xml.rels",
		"ppt/slides/slide2.xml",
		"ppt/slides/_rels/slide2.xml.rels",
		"ppt/media/slide2.png",
	}

	for _, name := range expectedFiles {
		if _, ok := entries[name]; !ok {
			t.Errorf("missing expected zip entry: %s", name)
		}
	}

	if _, ok := entries["ppt/notesSlides/notesSlide2.xml"]; ok {
		t.Errorf("slide 2 should not have notesSlide2.xml")
	}
}

func verifyNotesSlideContent(t *testing.T, entries map[string]*zip.File) {
	t.Helper()
	notesFile, ok := entries["ppt/notesSlides/notesSlide1.xml"]
	if !ok {
		t.Fatalf("notesSlide1.xml missing")
	}
	rc, err := notesFile.Open()
	if err != nil {
		t.Fatalf("failed to read notesSlide1.xml: %v", err)
	}
	defer func() { _ = rc.Close() }()

	notesBytes, _ := io.ReadAll(rc)
	if !strings.Contains(string(notesBytes), "Speaker notes for slide 1.") {
		t.Errorf("expected speaker notes text in notesSlide1.xml, got: %s", string(notesBytes))
	}
}

func verifySlideImageHeader(t *testing.T, entries map[string]*zip.File) {
	t.Helper()
	pngFile, ok := entries["ppt/media/slide1.png"]
	if !ok {
		t.Fatalf("slide1.png missing")
	}
	prc, err := pngFile.Open()
	if err != nil {
		t.Fatalf("failed to read slide1.png: %v", err)
	}
	defer func() { _ = prc.Close() }()

	pngHeader := make([]byte, 8)
	_, _ = io.ReadFull(prc, pngHeader)
	if !bytes.Equal(pngHeader, []byte("\x89PNG\r\n\x1a\n")) {
		t.Errorf("slide1.png does not start with PNG magic header: %x", pngHeader)
	}
}
