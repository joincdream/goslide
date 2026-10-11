package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yundream/goslide/internal/model"
)

// EventType defines the category of SSE message sent to clients.
type EventType string

const (
	EventReload  EventType = "reload"
	EventError   EventType = "error"
	EventPatch   EventType = "patch"
	EventWarning EventType = "warning"
)

// SSEMessage represents a structured payload broadcasted via Server-Sent Events.
type SSEMessage struct {
	Type        EventType          `json:"type"`
	Title       string             `json:"title,omitempty"`
	Message     string             `json:"message,omitempty"`
	File        string             `json:"file,omitempty"`
	Patches     []SlidePatch       `json:"patches,omitempty"`
	Diagnostics []model.Diagnostic `json:"diagnostics,omitempty"`
}

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	clientChan := make(chan string, 4)

	s.clientsMu.Lock()
	s.clients[clientChan] = struct{}{}
	s.clientsMu.Unlock()

	defer func() {
		s.clientsMu.Lock()
		delete(s.clients, clientChan)
		close(clientChan)
		s.clientsMu.Unlock()
	}()

	// Send initial handshake comment
	_, _ = fmt.Fprintf(w, ": connected\n\n")
	flusher.Flush()

	notify := r.Context().Done()
	for {
		select {
		case <-notify:
			return
		case msg, ok := <-clientChan:
			if !ok {
				return
			}
			_, err := fmt.Fprintf(w, "data: %s\n\n", msg)
			if err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

// BroadcastMessage sends a structured SSEMessage (marshaled to JSON) to all connected SSE clients.
func (s *Server) BroadcastMessage(msg SSEMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	s.clientsMu.RLock()
	defer s.clientsMu.RUnlock()

	for ch := range s.clients {
		select {
		case ch <- string(data):
		default:
			// Non-blocking drop if client buffer is full
		}
	}
}

// BroadcastReload notifies all clients to trigger a full page reload.
func (s *Server) BroadcastReload() {
	s.BroadcastMessage(SSEMessage{
		Type: EventReload,
	})
}

// BroadcastError sends a structured build/render error overlay message to all connected SSE clients.
func (s *Server) BroadcastError(title, msg, file string) {
	s.BroadcastMessage(SSEMessage{
		Type:    EventError,
		Title:   title,
		Message: msg,
		File:    file,
	})
}

// BroadcastPatches sends incremental slide patches to all connected SSE clients.
func (s *Server) BroadcastPatches(patches []SlidePatch) {
	s.BroadcastMessage(SSEMessage{
		Type:    EventPatch,
		Patches: patches,
	})
}

// BroadcastWarnings sends non-blocking DSL diagnostics as floating warning toasts to connected SSE clients.
func (s *Server) BroadcastWarnings(diagnostics []model.Diagnostic) {
	if len(diagnostics) == 0 {
		return
	}
	s.BroadcastMessage(SSEMessage{
		Type:        EventWarning,
		Diagnostics: diagnostics,
	})
}
