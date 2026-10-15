package main

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	cdruntime "github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

type componentRendering struct {
	Groups []string `json:"groups"`
	Nodes  []string `json:"nodes"`
	Edges  []string `json:"edges"`
}

// TestComponentBoundariesRendered proves the editor and CLI use the same real
// component owners and preserve neighboring containers as separate elements.
func TestComponentBoundariesRendered(t *testing.T) {
	if !hasChrome() {
		t.Skip("skipping: Chrome/Chromium not available in PATH")
	}
	const pkg = "goa.design/model/examples/component_boundaries/model"
	outDir := t.TempDir()
	if err := runSVG(pkg, config{dir: outDir, port: 0, timeout: 30 * time.Second, all: true}); err != nil {
		t.Fatalf("render component views: %v", err)
	}
	design, err := loadDesign(pkg, false)
	if err != nil {
		t.Fatalf("load component views: %v", err)
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	server := NewServer(design)
	httpServer := &http.Server{ReadHeaderTimeout: 3 * time.Second}
	mux := http.NewServeMux()
	mux.Handle("/exports/", http.StripPrefix("/exports/", http.FileServer(http.Dir(outDir))))
	serverDone := make(chan error, 1)
	go func() {
		serverDone <- server.ServeOnListener(outDir, "", httpServer, mux, listener)
	}()
	t.Cleanup(func() {
		if err := httpServer.Close(); err != nil {
			t.Errorf("close editor server: %v", err)
		}
		if err := <-serverDone; err != nil && err != http.ErrServerClosed {
			t.Errorf("editor server: %v", err)
		}
	})
	for _, test := range []struct {
		view   string
		groups []string
		nodes  []string
		edges  []string
	}{
		{"Owner Boundary", []string{"Application: Handler,Store"},
			[]string{"Database", "Handler", "Store", "Web Client", "Worker"},
			[]string{"Handler -> Worker", "Store -> Database"}},
		{"External Boundaries", []string{"Application: Handler,Store", "Jobs: Worker"},
			[]string{"Database", "Handler", "Store", "Web Client", "Worker"},
			[]string{"Handler -> Worker", "Store -> Database"}},
		{"Container Node", []string{"Application: Handler,Store"},
			[]string{"Database", "Handler", "Jobs", "Store", "Worker"},
			[]string{"Handler -> Jobs", "Handler -> Worker", "Store -> Database"}},
	} {
		t.Run(test.view, func(t *testing.T) {
			ctx, cleanup := newChromeContext(t)
			defer cleanup()
			pageURL := fmt.Sprintf("http://%s/?id=%s", listener.Addr().String(), url.QueryEscape(test.view))
			var exported string
			const exportScript = `(async () => {
				await window.gdata.autoLayout();
				return window.gdata.exportSVG();
			})()`
			if err := chromedp.Run(ctx,
				chromedp.Navigate(pageURL),
				chromedp.Poll(`window.gdata && window.gdata.nodes().every(node => node.ref?.isConnected)`, nil),
				chromedp.Evaluate(exportScript, &exported, func(params *cdruntime.EvaluateParams) *cdruntime.EvaluateParams {
					return params.WithAwaitPromise(true)
				}),
			); err != nil {
				t.Fatalf("render editor component view: %v", err)
			}
			editorPath := filepath.Join(outDir, "editor-"+test.view+".svg")
			if err := os.WriteFile(editorPath, []byte(exported), 0600); err != nil {
				t.Fatalf("save editor export: %v", err)
			}
			// Inspect real SVG rectangles after export. A missing boundary or a
			// neighboring element inside the wrong boundary changes these groups.
			const inspectScript = `(() => {
				const nodes = [...document.querySelectorAll("g.node")].map(node => ({
					id: node.id,
					name: node.querySelector('[data-field="name"]').textContent.trim(),
					box: node.querySelector(".nodeBorder").getBoundingClientRect(),
				}));
				const names = new Map(nodes.map(node => [node.id, node.name]));
				const groups = [...document.querySelectorAll("g.group")].map(group => {
					const box = group.querySelector("rect").getBoundingClientRect();
					const members = nodes.filter(node => node.box.left >= box.left - 0.5 &&
						node.box.right <= box.right + 0.5 && node.box.top >= box.top - 0.5 &&
						node.box.bottom <= box.bottom + 0.5).map(node => node.name).sort();
					return group.querySelector("text").textContent.trim() + ": " + members.join(",");
				}).sort();
				const edges = [...document.querySelectorAll("g.edge")].map(edge =>
					names.get(edge.dataset.from) + " -> " + names.get(edge.dataset.to)).sort();
				return {groups, nodes: nodes.map(node => node.name).sort(), edges};
			})()`
			want := componentRendering{Groups: test.groups, Nodes: test.nodes, Edges: test.edges}
			// Export inspection uses another browser so the editor can retain its
			// unsaved-layout navigation protection.
			exportContext, cleanupExports := newChromeContext(t)
			defer cleanupExports()
			for _, svgPath := range []string{filepath.Join(outDir, test.view+".svg"), editorPath} {
				var got componentRendering
				fileURL := fmt.Sprintf("http://%s/exports/%s", listener.Addr().String(), url.PathEscape(filepath.Base(svgPath)))
				if err := chromedp.Run(exportContext,
					chromedp.Navigate(fileURL),

					chromedp.Evaluate(inspectScript, &got),
				); err != nil {
					t.Fatalf("inspect component SVG: %v", err)
				}
				if !reflect.DeepEqual(got, want) {
					t.Errorf("%s: got %+v, want %+v", filepath.Base(svgPath), got, want)
				}
			}
		})
	}
}
