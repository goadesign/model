package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"sync"
)

type Server struct {
	modelLock sync.Mutex
	model     []byte
}

func (s Server) setModel(m []byte) {
	s.modelLock.Lock()
	s.model = m
	s.modelLock.Unlock()
}

func (s Server) Serve(outDir string, devmode bool, port int) error {
	layoutFile := path.Join(outDir, "layout.json")

	if devmode {
		// in devmode (go run), serve the webapp from filesystem
		fs := http.FileSystem(http.Dir("./cmd/stz-edit/webapp/dist"))
		http.Handle("/", http.FileServer(fs))
	} else {
		// the TS/React webapp is embeded in the go executable using esc https://github.com/mjibson/esc
		// to update the webapp, run `make generate` in the root dir of the repo
		http.Handle("/", http.FileServer(FS(false)))
	}

	http.HandleFunc("/data/model.json", func(w http.ResponseWriter, r *http.Request) {
		s.modelLock.Lock()
		_, _ = w.Write(s.model)
		s.modelLock.Unlock()
	})
	http.HandleFunc("/data/layout.json", func(w http.ResponseWriter, r *http.Request) {
		if fileExists(layoutFile) {
			http.ServeFile(w, r, layoutFile)
		} else {
			fmt.Fprint(w, "{}")
		}
	})

	http.HandleFunc("/data/save", func(w http.ResponseWriter, r *http.Request) {
		savedData := mdl{}
		err := json.NewDecoder(r.Body).Decode(&savedData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		tmp, ok := r.URL.Query()["id"]
		if !ok {
			http.Error(w, "Param id is missing", http.StatusBadRequest)
			return
		}
		id := tmp[0]

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

type mdl struct {
	Layout interface{} `json:"layout,omitempty"`
	SVG    string      `json:"svg,omitempty"`
}
