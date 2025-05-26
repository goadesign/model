package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"goa.design/model/mdl"
	model "goa.design/model/pkg"
)

//go:embed webapp/dist/*
var distFS embed.FS

type (
	// Server implements a HTTP server with 4 endpoints for the model diagram editor
	Server struct {
		design []byte
		lock   sync.RWMutex
		outDir string
	}

	// Layout represents position info saved for one view (diagram)
	Layout = map[string]any
)

// NewServer creates a server that serves the given design
func NewServer(d *mdl.Design) *Server {
	s := &Server{}
	s.SetDesign(d)
	return s
}

// Serve starts the HTTP server on localhost with the given port
func (s *Server) Serve(outDir, devDistPath string, port int) error {
	s.outDir = outDir

	s.setupRoutes(devDistPath)

	server := &http.Server{
		Addr:              fmt.Sprintf("127.0.0.1:%d", port),
		ReadHeaderTimeout: 3 * time.Second,
	}

	fmt.Printf("mdl %s, editor started. Open http://localhost:%d in your browser.\n", model.Version(), port)
	return server.ListenAndServe()
}

// setupRoutes configures all HTTP routes
func (s *Server) setupRoutes(devDistPath string) {
	// Serve static files
	if devDistPath != "" {
		// Development mode: serve from filesystem
		fs := http.FileSystem(http.Dir(devDistPath))
		http.Handle("/", http.FileServer(fs))
	} else {
		// Production mode: serve from embedded files
		sub, _ := fs.Sub(distFS, "webapp/dist")
		http.Handle("/", http.FileServer(http.FS(sub)))
	}

	// API endpoints
	http.HandleFunc("/data/model.json", s.handleModelData)
	http.HandleFunc("/data/layout.json", s.handleLayoutData)
	http.HandleFunc("/data/save", s.handleSave)
}

// handleModelData serves the JSON representation of the architecture model
func (s *Server) handleModelData(w http.ResponseWriter, r *http.Request) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	w.Write(s.design)
}

// handleLayoutData serves the view element positions indexed by view id
func (s *Server) handleLayoutData(w http.ResponseWriter, r *http.Request) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	layouts, err := s.loadLayouts()
	if err != nil {
		s.handleError(w, fmt.Errorf("failed to load layouts: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(layouts)
}

// handleSave saves the SVG representation for a view
func (s *Server) handleSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		s.handleError(w, fmt.Errorf("missing id parameter"))
		return
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	if err := s.saveSVG(id, r.Body); err != nil {
		s.handleError(w, fmt.Errorf("failed to save SVG: %w", err))
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// saveSVG saves the SVG content to a file
func (s *Server) saveSVG(id string, body io.Reader) error {
	svgFile := path.Join(s.outDir, id+".svg")
	f, err := os.Create(svgFile)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, body)
	return err
}

// SetDesign updates the design served by the server
func (s *Server) SetDesign(d *mdl.Design) {
	b, err := json.Marshal(d)
	if err != nil {
		panic("failed to serialize design: " + err.Error()) // This should never happen
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	s.design = b
}

// handleError writes the error to stderr and returns an HTTP error response
func (s *Server) handleError(w http.ResponseWriter, err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// loadLayouts reads layout information from SVG files and fallback layout.json
func (s *Server) loadLayouts() ([]byte, error) {
	layouts := make(map[string]Layout)

	// Load fallback layout.json for backwards compatibility
	if err := s.loadLayoutJSON(layouts); err != nil {
		return nil, err
	}

	// Load individual layouts from SVG files
	if err := s.loadLayoutsFromSVGs(layouts); err != nil {
		return nil, err
	}

	return json.Marshal(layouts)
}

// loadLayoutJSON loads the fallback layout.json file
func (s *Server) loadLayoutJSON(layouts map[string]Layout) error {
	layoutFile := path.Join(s.outDir, "layout.json")
	if !fileExists(layoutFile) {
		return nil // No fallback file, that's okay
	}

	b, err := os.ReadFile(layoutFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &layouts)
}

// loadLayoutsFromSVGs extracts layout information from SVG files
func (s *Server) loadLayoutsFromSVGs(layouts map[string]Layout) error {
	files, err := os.ReadDir(s.outDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".svg") {
			continue
		}

		if err := s.loadLayoutFromSVG(f.Name(), layouts); err != nil {
			return err
		}
	}

	return nil
}

// loadLayoutFromSVG extracts layout information from a single SVG file
func (s *Server) loadLayoutFromSVG(filename string, layouts map[string]Layout) error {
	const (
		beginMark = "<script type=\"application/json\"><![CDATA["
		endMark   = "]]></script>"
	)

	b, err := os.ReadFile(path.Join(s.outDir, filename))
	if err != nil {
		return err
	}

	// Find the JSON script block
	beginBytes := []byte(beginMark)
	endBytes := []byte(endMark)

	begin := bytes.Index(b, beginBytes)
	if begin == -1 {
		return nil // No layout data in this SVG
	}
	begin += len(beginBytes)

	end := bytes.Index(b, endBytes)
	if end == -1 {
		return fmt.Errorf("malformed SVG: missing end marker in %s", filename)
	}

	layoutData := b[begin:end]

	var data map[string]any
	if err := json.Unmarshal(layoutData, &data); err != nil {
		return fmt.Errorf("invalid JSON in SVG %s: %w", filename, err)
	}

	if layout, ok := data["layout"].(map[string]any); ok {
		id := strings.TrimSuffix(filename, ".svg")
		layouts[id] = layout
	}

	return nil
}

// fileExists checks if a file exists and is not a directory
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		fail("Can't read FileInfo: %s", err)
	}
	return !info.IsDir()
}
