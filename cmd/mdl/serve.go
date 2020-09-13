package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/jaschaephraim/lrserver"
	"goa.design/model/mdl"
	"golang.org/x/tools/go/packages"
)

func serve(pkg, config, out string, port int, debug bool) error {
	// Generate diagrams, return right away on error
	views, err := loadViews(pkg, out, debug)
	if err != nil {
		return err
	}

	// Watch model design and regenerate on change
	pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedFiles}, pkg)
	if err != nil {
		return err
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()
	if err = watcher.Add(filepath.Dir(pkgs[0].GoFiles[0])); err != nil {
		return err
	}

	// Create live reload server and hookup to watcher
	lr := lrserver.New(lrserver.DefaultName, lrserver.DefaultPort)
	lr.SetStatusLog(nil)
	lr.SetErrorLog(nil)
	go func() {
		if err := lr.ListenAndServe(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}()
	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				if !strings.HasPrefix(filepath.Base(ev.Name), tmpDirPrefix) {
					views, err = loadViews(pkg, out, debug)
					lr.Reload(ev.Name)
				}
			case err := <-watcher.Errors:
				fmt.Fprintln(os.Stderr, err.Error())
			}
		}
	}()

	// Serve generated diagrams
	listindex := template.Must(template.New("listindex").Parse(listindexT))
	errindex := template.Must(template.New("errindex").Parse(errindexT))
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if err != nil {
			errindex.Execute(w, err.Error())
			return
		}
		view, ok := views[req.URL.Path[1:]]
		if !ok {
			if err = listindex.Execute(w, listData(views)); err != nil {
				errindex.Execute(w, err.Error())
			}
			return
		}
		data := &ViewData{
			Key:                 view.Key,
			Title:               view.Title,
			Description:         view.Description,
			Version:             view.Version,
			MermaidSource:       template.JS(view.Mermaid),
			MermaidLegendSource: template.JS(view.Legend),
			MermaidConfig:       template.JS(config),
			CSS:                 template.CSS(DefaultCSS),
		}
		if err = indexTmpl.Execute(w, data); err != nil {
			errindex.Execute(w, err.Error())
			return
		}
	})

	fmt.Printf("[Model] listening on :%d\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func listData(views map[string]*mdl.RenderedView) []*ViewData {
	keys := make([]string, len(views))
	i := 0
	for k := range views {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	data := make([]*ViewData, len(keys))
	for i, k := range keys {
		view := views[k]
		title := view.Title
		if title == "" {
			title = view.Key
		}
		data[i] = &ViewData{
			Key:         view.Key,
			Title:       title,
			Description: view.Description,
		}
	}
	return data
}

const errindexT = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title>[ERROR]</title>
</head>
<body>
	<div class="error">
	<b>ERROR</b><br/><br/>
	{{ . }}
	</div>
	<script src="http://localhost:35729/livereload.js"></script>
</body>
</html>
`

const listindexT = `<!DOCTYPE html>
<html lang="end">
<head>
	<meta charset="utf-8">
	<title>Diagrams</title>
	<style>
		body {font-family: Arial;}
		h1 {font-weight: bold; font-size: 1.5rem;}
		ul {font-size: 1.2em;}
		li {line-height: 1.5;}
        a {color:#001f3f; text-decoration: none;}
        a:hover {color:#0074d9;}
	</style>
</head>
<body>
	<h1>Available views:</h1>
	<ul>
	{{- range . }}
		<li><a href="/{{ .Key }}">{{ .Title }}: <i>{{ .Description }}</i></a></li>
	{{ end }}
	</ul>
</body>
</html>
`
