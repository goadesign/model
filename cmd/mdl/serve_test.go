package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"goa.design/model/mdl"
)

// minimalDesign returns a design with no views/elements sufficient for handlers
func minimalDesign() *mdl.Design { return &mdl.Design{} }

func TestServerHandlers(t *testing.T) {
	s := NewServer(minimalDesign())

	mux := http.NewServeMux()
	dir := t.TempDir()
	s.setupRoutesToMux(mux, "")

	// model.json
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/data/model.json", nil)
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("model.json status: %d", w.Code)
	}
	var tmp any
	if err := json.Unmarshal(w.Body.Bytes(), &tmp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	// layout.json (no files, should still be 200 JSON)
	// point server to temp dir for layout loading
	s.outDir = dir
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/data/layout.json", nil)
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("layout.json status: %d", w.Code)
	}
	if err := json.Unmarshal(w.Body.Bytes(), &tmp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	// save endpoint writes file
	body := bytes.NewBufferString("<svg><!--test--></svg>")
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodPost, "/data/save?id=view1", io.NopCloser(body))
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusAccepted {
		t.Fatalf("save status: %d", w.Code)
	}
	// file created
	saved := filepath.Join(dir, "view1.svg")
	content, err := os.ReadFile(saved)
	if err != nil {
		t.Fatalf("read saved SVG: %v", err)
	}
	if string(content) != "<svg><!--test--></svg>" {
		t.Fatalf("unexpected SVG content %q", content)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read output directory: %v", err)
	}
	if len(entries) != 1 || entries[0].Name() != "view1.svg" {
		t.Fatalf("expected only atomic target file, got %v", entries)
	}

	w = httptest.NewRecorder()
	r = httptest.NewRequest(
		http.MethodPost,
		"/data/save?id=..%2Foutside",
		bytes.NewBufferString("<svg/>"),
	)
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("unsafe save status: %d", w.Code)
	}
}
