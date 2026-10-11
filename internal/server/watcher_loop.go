package server

import (
	"bytes"
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/yundream/goslide/internal/i18n"
	"github.com/yundream/goslide/internal/model"
)

func (s *Server) watchLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.watcher.Events():
			s.handleFileChange(ctx)
		case err, ok := <-s.watcher.Errors():
			if !ok {
				return
			}
			log.Printf("[goslide] Watcher error: %v", err)
		}
	}
}

func (s *Server) handleFileChange(ctx context.Context) {
	start := time.Now()

	deck, err := s.rebuildAndWarmCache(ctx)
	if err != nil {
		return
	}

	diff := s.reconcileSnapshot(deck)
	s.dispatchDiff(ctx, deck, diff, start)
}

func (s *Server) rebuildAndWarmCache(ctx context.Context) (*model.Deck, error) {
	raw, readErr := os.ReadFile(s.absFile)
	if readErr != nil {
		log.Printf("[goslide] Warning: failed to read file on change: %v", readErr)
		return nil, readErr
	}

	deck, parseErr := s.parser.Parse(ctx, bytes.NewReader(raw))
	if parseErr != nil {
		log.Printf("[goslide] Warning: parse failed on file change: %v", parseErr)
		s.BroadcastError("Markdown Syntax Error", parseErr.Error(), filepath.Base(s.absFile))
		return nil, parseErr
	}

	s.logDiagnostics(deck.Diagnostics)
	s.BroadcastWarnings(deck.Diagnostics)

	// Pre-render full HTML to warm cache and update lastHTML
	var fullBuf bytes.Buffer
	if renderErr := s.renderer.Render(ctx, deck, &fullBuf); renderErr != nil {
		log.Printf("[goslide] Warning: render failed on file change: %v", renderErr)
		s.BroadcastError("Render Error", renderErr.Error(), filepath.Base(s.absFile))
		return nil, renderErr
	}
	renderedBytes := injectSSEScript(fullBuf.Bytes())
	s.lastHTMLMu.Lock()
	s.lastHTML = renderedBytes
	s.lastHTMLMu.Unlock()

	return deck, nil
}

func (s *Server) reconcileSnapshot(deck *model.Deck) DiffResult {
	s.lastSnapshotMu.RLock()
	oldSnap := s.lastSnapshot
	s.lastSnapshotMu.RUnlock()

	diff := CompareSnapshots(oldSnap, deck)

	s.lastSnapshotMu.Lock()
	s.lastSnapshot = NewDeckSnapshot(deck)
	s.lastSnapshotMu.Unlock()

	return diff
}

func (s *Server) dispatchDiff(ctx context.Context, deck *model.Deck, diff DiffResult, start time.Time) {
	bundle, _ := i18n.GetDefaultBundle()

	if diff.NeedsReload {
		reason := diff.ReloadReasonKey
		if bundle != nil {
			reason = bundle.T(diff.ReloadReasonKey)
		}
		log.Printf("[goslide] Rebuild completed in %v. Full reload triggered (%s)...", time.Since(start), reason)
		s.BroadcastReload()
		return
	}

	if len(diff.ChangedIndices) > 0 {
		s.dispatchSlidePatches(ctx, deck, diff.ChangedIndices, start, bundle)
		return
	}

	log.Printf("[goslide] File saved with no slide content changes (%v)", time.Since(start))
}

func (s *Server) dispatchSlidePatches(ctx context.Context, deck *model.Deck, indices []int, start time.Time, bundle *i18n.Bundle) {
	var patches []SlidePatch
	for _, idx := range indices {
		var patchBuf bytes.Buffer
		if err := s.renderer.RenderSlide(ctx, deck, idx, &patchBuf); err == nil {
			patches = append(patches, SlidePatch{
				Index: idx + 1,
				HTML:  patchBuf.String(),
			})
		} else {
			log.Printf("[goslide] Warning: failed to render slide patch %d: %v", idx+1, err)
		}
	}
	if len(patches) > 0 {
		if bundle != nil {
			log.Print(bundle.T("server.hmr.patch.log", patches[0].Index, time.Since(start)))
		} else {
			log.Printf("[goslide] Slide %d incrementally patched in %v", patches[0].Index, time.Since(start))
		}
		s.BroadcastPatches(patches)
	}
}

func (s *Server) logDiagnostics(diagnostics []model.Diagnostic) {
	for _, diag := range diagnostics {
		if len(diag.Candidates) > 0 {
			log.Printf("[goslide] ⚠️  Slide %d (Line %d): %s (Raw: %q, Available: %v)",
				diag.SlideIndex, diag.Line, diag.Message, diag.RawSnippet, diag.Candidates)
		} else {
			log.Printf("[goslide] ⚠️  Slide %d (Line %d): %s (Raw: %q)",
				diag.SlideIndex, diag.Line, diag.Message, diag.RawSnippet)
		}
	}
}
