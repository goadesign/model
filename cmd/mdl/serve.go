package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
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

	// SavedView is the data structure created and updated by the single page
	// app for each design view.
	SavedView struct {
		Layout interface{} `json:"layout,omitempty"`
		SVG    string      `json:"svg,omitempty"`
	}
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
	layoutFile := path.Join(outDir, "layout.json")

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
		if fileExists(layoutFile) {
			http.ServeFile(w, r, layoutFile)
		} else {
			fmt.Fprint(w, "{}")
		}
	})

	http.HandleFunc("/data/save", func(w http.ResponseWriter, r *http.Request) {
		var savedData SavedView
		err := json.NewDecoder(r.Body).Decode(&savedData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Param id is missing", http.StatusBadRequest)
			return
		}

		s.lock.Lock()
		defer s.lock.Unlock()

		data := make(map[string]interface{})
		if fileExists(layoutFile) {
			file, err := ioutil.ReadFile(layoutFile)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// if the file contains garbage, ignore it's content
			_ = json.Unmarshal(file, &data)
		}
		data[id] = savedData.Layout

		out, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = ioutil.WriteFile(layoutFile, out, 0644)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		svgFile := path.Join(outDir, id+".svg")
		_ = ioutil.WriteFile(svgFile, []byte(savedData.SVG), 0644)

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
