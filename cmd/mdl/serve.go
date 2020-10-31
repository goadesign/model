package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"goa.design/model/mdl"
)

type (

	// Server implements a HTTP server with 4 endpoints:
	//
	//   * GET requests to "/" return the diagram editor single page app implemented in the "webapp" directory.
	//   * GET requests to "/data/model.json" return the JSON representation of the architecture model.
	//   * GET requests to "/data/layout.json" return the view element positions indexed by view id.
	//   * POST requests to "/data/save?id=<ID>" saves the SVG representation for the view with the given id.
	//     The request body must be a JSON representation of a SavedView data structure.
	//
	// Server is intended to provide the backend for the model single page app diagram editor.
	Server struct {
		design []byte
		lock   sync.Mutex
	}

	// Layout is position info saved for one view (diagram)
	Layout = map[string]interface{}

	// Layouts is a map from view key to the view Layout
	Layouts = map[string]Layout
)

// NewServer created a server that serves the given design.
func NewServer(d *mdl.Design) *Server {
	var s Server
	s.SetDesign(d)
	return &s
}

// Serve starts the HTTP server on localhost with the given port. outDir
// indicates where the view data structures are located. If devmode is true then
// the single page app is served directly from the source under the "webapp"
// directory. Otherwise it is served from the code embedded in the Go executable.
func (s *Server) Serve(outDir string, devmode bool, port int) error {

	if devmode {
		// in devmode (go run), serve the webapp from filesystem
		fs := http.FileSystem(http.Dir("./cmd/mdl/webapp/dist"))
		http.Handle("/", http.FileServer(fs))
	} else {
		// the TS/React webapp is embeded in the go executable using esc https://github.com/mjibson/esc
		// to update the webapp, run `make generate` in the root dir of the repo
		http.Handle("/", http.FileServer(FS(false)))
	}

	http.HandleFunc("/data/model.json", func(w http.ResponseWriter, r *http.Request) {
		s.lock.Lock()
		defer s.lock.Unlock()
		_, _ = w.Write(s.design)
	})

	http.HandleFunc("/data/layout.json", func(w http.ResponseWriter, r *http.Request) {
		s.lock.Lock()
		defer s.lock.Unlock()

		b, err := loadLayouts(outDir)
		if err != nil {
			fmt.Println(err)
		} else {
			_, _ = w.Write(b)
		}
	})

	http.HandleFunc("/data/save", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Param id is missing", http.StatusBadRequest)
			return
		}

		s.lock.Lock()
		defer s.lock.Unlock()

		svgFile := path.Join(outDir, id+".svg")
		f, err := os.Create(svgFile)
		if err != nil {
			msg := fmt.Sprintf("Saving failed, can't write to %s: %s!\n", svgFile, err)
			fmt.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		defer func() { _ = f.Close() }()
		_, _ = io.Copy(f, r.Body)

		w.WriteHeader(http.StatusAccepted)
	})

	// start the server
	fmt.Printf("Editor started. Open http://localhost:%d in your browser.\n", port)
	return http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil)
}

// SetDesign updates the design served by s.
//
// Note: it would have been more efficient to use the raw bytes read from the
// generated file instead of going through the unmarshal/marshal cycle however
// this approach is safer, makes it clearer and easier to compose. Also it is
// not expected that the model would need to be updated often.
func (s *Server) SetDesign(d *mdl.Design) {
	b, err := json.Marshal(d)
	if err != nil {
		panic("failed to serialize design: " + err.Error()) // bug
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	s.design = b
}

// loadLayouts lists out directory and reads layout info from SVG files
// for backwards compatibility, fallback to layout.json
func loadLayouts(dir string) ([]byte, error) {
	beginMark := []byte("<script type=\"application/json\"><![CDATA[")
	endMark := []byte("]]></script>")

	// first, read the fallback layout.json, then merge individual layouts from SVGs
	var layouts Layouts = make(map[string]Layout)
	lj := path.Join(dir, "layout.json")
	if fileExists(lj) {
		b, err := ioutil.ReadFile(lj)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(b, &layouts)
		if err != nil {
			return nil, err
		}
	}

	var svgFiles []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".svg") {
			svgFiles = append(svgFiles, f.Name())
		}
	}
	for _, file := range svgFiles {
		b, err := ioutil.ReadFile(path.Join(dir, file))
		if err != nil {
			return nil, err
		}

		// look for the first script block
		begin := bytes.Index(b, beginMark) + len(beginMark)
		end := bytes.Index(b, endMark)
		b = b[begin:end]

		var l Layout = make(map[string]interface{})
		err = json.Unmarshal(b, &l)
		if err != nil {
			return nil, err
		}
		id := file[:len(file)-4]
		layouts[id] = l["layout"].(Layout)
	}
	b, err := json.Marshal(layouts)
	if err != nil {
		return nil, err
	}
	return b, nil
}
