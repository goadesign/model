package main

//go:generate esc -o webapp.go -pkg main -prefix webapp/dist webapp/dist/

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func main() {
	projectDir := flag.String("project-dir", "./", "Directory that contains model.json file."+
		" The model.layout.json file will also be saved here.")
	port := flag.Int("port", 8080, "Local HTTP port.")
	devmode := os.Getenv("DEVMODE") == "1"

	flag.Parse()

	dir, err := filepath.Abs(*projectDir)
	if err != nil {
		fail(err.Error())
	}

	modelFile := path.Join(dir, "model.json")
	layoutFile := path.Join(dir, "model.layout.json")
	if !fileExists(modelFile) {
		fail("Specified input file %s does not exist.", modelFile)
	}

	fmt.Printf("Reading from:      %s\nWriting layout to: %s\n", modelFile, layoutFile)

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
		http.ServeFile(w, r, modelFile)
	})
	http.HandleFunc("/data/model.layout.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, layoutFile)
	})

	// save handler: merges layout data for view identified by query param "id" into model.layout.json
	// if model.layout.json does not exist, it is created
	http.HandleFunc("/data/save", func(w http.ResponseWriter, r *http.Request) {
		var savedData []interface{}
		err := json.NewDecoder(r.Body).Decode(&savedData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, ok := r.URL.Query()["id"]
		if !ok {
			http.Error(w, "Param id is missing", http.StatusBadRequest)
			return
		}
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
		view := make(map[string]interface{})
		view["elements"] = savedData
		data[id[0]] = view

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
		w.WriteHeader(http.StatusAccepted)
	})

	// start the server
	fmt.Printf("Editor started. Open http://localhost:%d in your browser.\n", *port)
	err = http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", *port), nil)
	fail(err.Error())
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
