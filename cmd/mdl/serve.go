package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/jaschaephraim/lrserver"
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
				views, err = loadViews(pkg, out, debug)
				lr.Reload(ev.Name)
			case err := <-watcher.Errors:
				fmt.Fprintln(os.Stderr, err.Error())
			}
		}
	}()

	// Serve generated diagrams
	errindex := template.Must(template.New("errindex").Parse(errindexT))
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if err != nil {
			errindex.Execute(w, err.Error())
			return
		}
		view, ok := views[req.URL.Path[1:]]
		if !ok {
			links := make([]string, len(views))
			i := 0
			for v := range views {
				links[i] = fmt.Sprintf(`<a href="%s">%s</a>`, v, v)
				i++
			}
			data := fmt.Sprintf("no view with key %s, available views:<br/>%s", req.URL.Path[1:], strings.Join(links, "<br/>"))
			errindex.Execute(w, template.HTML(data))
			return
		}
		data := &ViewData{
			Title:         view.Title,
			Description:   view.Description,
			Version:       view.Version,
			MermaidSource: template.JS(view.Mermaid),
			MermaidConfig: template.JS(config),
			CSS:           template.CSS(DefaultCSS),
		}
		if err = indexTmpl.Execute(w, data); err != nil {
			errindex.Execute(w, err.Error())
			return
		}
	})

	fmt.Printf("[Model] listening on :%d\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
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
