package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	cdruntime "github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"

	"goa.design/model/mdl"
)

type (
	// twoGraphsResult reports one page after two graphs were laid out together.
	twoGraphsResult struct {
		// Duplicates lists ids that appear more than once in the page.
		Duplicates []string `json:"duplicates"`
		// Graphs describes the first and the second view, in that order.
		Graphs []graphIsolation `json:"graphs"`
	}

	// graphIsolation describes one graph's live svg and its exported file.
	graphIsolation struct {
		// Nodes lists the model element ids of the graph's nodes.
		Nodes []string `json:"nodes"`
		// OwnNodes is true when the svg holds exactly this graph's node elements.
		OwnNodes bool `json:"ownNodes"`
		// Zoom is the transform that fitting the layout set on this svg.
		Zoom string `json:"zoom"`
		// Grid is true when this svg shows its own grid.
		Grid bool `json:"grid"`
		// Unresolved lists url(#…) and #… references with no target in this svg.
		Unresolved []string `json:"unresolved"`
		// Export is the graph's exported SVG file.
		Export string `json:"export"`
		// ExportRootID is the id of the exported root element.
		ExportRootID string `json:"exportRootId"`
		// ExportNodes lists the ids of the exported node elements.
		ExportNodes []string `json:"exportNodes"`
		// ExportLeaks lists exported values that still name the live id prefix.
		ExportLeaks []string `json:"exportLeaks"`
		// ExportDefs is true when the export has the saved definition ids.
		ExportDefs bool `json:"exportDefs"`
		// ExportGrid is true when the exported grid uses the saved grid ids.
		ExportGrid bool `json:"exportGrid"`
	}
)

// TestTwoGraphsInOneDocument proves that two graphs drawn in one page keep
// their own nodes, ids, fit, grid, and export. The editor shows one view at a
// time, so the test opens the first view, switches to the second view with the
// view selector, puts the first view's svg back in the page, and lays out both
// graphs at the same time. Both graphs use the editor's fixed 25px grid.
func TestTwoGraphsInOneDocument(t *testing.T) {
	if !hasChrome() {
		t.Skip("skipping: Chrome/Chromium not available in PATH")
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	server := NewServer(twoViewDesign())
	httpServer := &http.Server{ReadHeaderTimeout: 3 * time.Second}
	serverDone := make(chan error, 1)
	outDir := t.TempDir()
	go func() {
		serverDone <- server.ServeOnListener(outDir, "", httpServer, http.NewServeMux(), listener)
	}()
	t.Cleanup(func() {
		if err := httpServer.Close(); err != nil {
			t.Errorf("close editor server: %v", err)
		}
		if err := <-serverDone; err != nil && err != http.ErrServerClosed {
			t.Errorf("editor server: %v", err)
		}
	})
	testContext, cleanup := newChromeContext(t)
	defer cleanup()

	const switchScript = `(() => {
		window.firstGraph = ` + editorGraphScript + `;
		const views = document.querySelector(".toolbar select");
		views.value = "Containers";
		views.dispatchEvent(new Event("change", {bubbles: true}));
		return true;
	})()`
	const secondDrawnScript = `(() => {
		const graph = ` + editorGraphScript + `;
		return graph !== null && graph !== window.firstGraph &&
			graph.nodes().every(node => node.ref?.isConnected);
	})()`
	const layoutScript = `(async () => {
		const graphs = [window.firstGraph, ` + editorGraphScript + `];
		const holder = document.createElement("div");
		document.body.append(holder);
		holder.append(graphs[0].svg);
		graphs.forEach(graph => graph.toggleGrid());
		await Promise.all(graphs.map(graph => graph.autoLayout()));
		const ids = [...document.querySelectorAll("[id]")].map(el => el.id);
		return {
			duplicates: ids.filter((id, index) => ids.indexOf(id) !== index),
			graphs: graphs.map(graph => {
				const svg = graph.svg;
				const refs = graph.nodes().map(node => node.ref);
				const drawn = [...svg.querySelectorAll("g.nodes > g.node")];
				const unresolved = [];
				for (const el of svg.querySelectorAll("*")) {
					for (const attr of el.attributes) {
						const target = attr.value.match(/^url\(#(.+)\)$/)?.[1] ??
							(attr.localName === "href" && attr.value.startsWith("#") ? attr.value.slice(1) : null);
						if (target && !svg.querySelector('[id="' + target + '"]')) unresolved.push(attr.value);
					}
				}
				const markup = graph.exportSVG();
				const exported = new DOMParser().parseFromString(markup, "image/svg+xml").documentElement;
				const leaks = [...exported.querySelectorAll("*")].flatMap(el => [...el.attributes]
					.filter(attr => (attr.name === "id" && attr.value.startsWith(graph.idPrefix)) ||
						attr.value.includes("#" + graph.idPrefix))
					.map(attr => attr.value));
				return {
					nodes: graph.nodes().map(node => node.id).sort(),
					ownNodes: drawn.length === refs.length && drawn.every(el => refs.includes(el)),
					zoom: svg.querySelector("g.zoom").getAttribute("transform") || "",
					grid: svg.querySelector("pattern") !== null,
					unresolved,
					export: markup,
					exportRootId: exported.id,
					exportNodes: [...exported.querySelectorAll("g.node")].map(node => node.id).sort(),
					exportLeaks: leaks,
					exportDefs: ["#inter-font-license", "marker#arrow", "#icon-circle", "#icon-cube"]
						.every(selector => exported.querySelector(selector) !== null),
					exportGrid: exported.querySelector("pattern#grid-pattern") !== null &&
						exported.querySelector("rect#grid-background")?.getAttribute("fill") === "url(#grid-pattern)",
				};
			}),
		};
	})()`
	var switched bool
	var result twoGraphsResult
	if err := chromedp.Run(testContext,
		chromedp.Navigate(fmt.Sprintf("http://%s/?id=Context", listener.Addr().String())),
		chromedp.Poll(editorDrawnScript, nil),
		chromedp.Evaluate(switchScript, &switched),
		chromedp.Poll(secondDrawnScript, nil),
		chromedp.Evaluate(layoutScript, &result, func(params *cdruntime.EvaluateParams) *cdruntime.EvaluateParams {
			return params.WithAwaitPromise(true)
		}),
	); err != nil {
		t.Fatalf("lay out two graphs in one page: %v", err)
	}

	if len(result.Duplicates) > 0 {
		t.Errorf("ids shared by the two graphs: %v", result.Duplicates)
	}
	if len(result.Graphs) != 2 {
		t.Fatalf("expected two graphs, got %d", len(result.Graphs))
	}
	for i, want := range []struct {
		view        string
		nodes       []string
		orientation string
	}{
		{"Context", []string{"shop", "user"}, "horizontal"},
		{"Containers", []string{"db", "user", "web"}, "vertical"},
	} {
		got := result.Graphs[i]
		if !reflect.DeepEqual(got.Nodes, want.nodes) || !got.OwnNodes {
			t.Errorf("%s: svg nodes %v (only its own: %t), want %v", want.view, got.Nodes, got.OwnNodes, want.nodes)
		}
		if got.Zoom == "" {
			t.Errorf("%s: layout fit was not applied to this graph's svg", want.view)
		}
		if !got.Grid {
			t.Errorf("%s: svg has no grid", want.view)
		}
		if len(got.Unresolved) > 0 {
			t.Errorf("%s: references to ids outside this svg: %v", want.view, got.Unresolved)
		}
		if got.ExportRootID != "graph" || !reflect.DeepEqual(got.ExportNodes, want.nodes) {
			t.Errorf("%s: export root %q nodes %v, want graph %v", want.view, got.ExportRootID, got.ExportNodes, want.nodes)
		}
		if len(got.ExportLeaks) > 0 || !got.ExportDefs || !got.ExportGrid {
			t.Errorf("%s: export ids not saved names: leaks %v, defs %t, grid %t",
				want.view, got.ExportLeaks, got.ExportDefs, got.ExportGrid)
		}
		path := filepath.Join(t.TempDir(), want.view+".svg")
		if err := os.WriteFile(path, []byte(got.Export), 0600); err != nil {
			t.Fatalf("save %s export: %v", want.view, err)
		}
		if issues := inspectArrowMarkerIssues(t, path); len(issues) > 0 {
			t.Errorf("%s: exported arrow markers: %v", want.view, issues)
		}
		if orientation := inspectNodeOrientation(t, path); orientation != want.orientation {
			t.Errorf("%s: exported layout is %s, want %s", want.view, orientation, want.orientation)
		}
	}
}

// twoViewDesign returns one model with two views that share the element
// "user" and differ in shapes, colors, and layout direction.
func twoViewDesign() *mdl.Design {
	return &mdl.Design{
		Name: "Two views",
		Model: &mdl.Model{
			People: []*mdl.Person{{
				ID: "user", Name: "Customer", Tags: "Element,Person,person",
				Relationships: []*mdl.Relationship{
					{ID: "browse", SourceID: "user", DestinationID: "shop", Description: "Browses", Tags: "Relationship"},
					{ID: "order", SourceID: "user", DestinationID: "web", Description: "Orders through", Tags: "Relationship"},
				},
			}},
			Systems: []*mdl.SoftwareSystem{{
				ID: "shop", Name: "Shop", Tags: "Element,Software System",
				Containers: []*mdl.Container{
					{ID: "web", Name: "Web", Tags: "Element,Container,web", Relationships: []*mdl.Relationship{
						{ID: "save", SourceID: "web", DestinationID: "db", Description: "Saves orders", Tags: "Relationship"},
					}},
					{ID: "db", Name: "Orders", Tags: "Element,Container,store"},
				},
			}},
		},
		Views: &mdl.Views{
			ContextViews: []*mdl.ContextView{{SoftwareSystemID: "shop", ViewProps: &mdl.ViewProps{
				Key:               "Context",
				AutoLayout:        &mdl.AutoLayout{RankDirection: mdl.RankLeftRight},
				ElementViews:      []*mdl.ElementView{{ID: "user"}, {ID: "shop"}},
				RelationshipViews: []*mdl.RelationshipView{{ID: "browse"}},
			}}},
			ContainerViews: []*mdl.ContainerView{{SoftwareSystemID: "shop", ViewProps: &mdl.ViewProps{
				Key:               "Containers",
				AutoLayout:        &mdl.AutoLayout{RankDirection: mdl.RankTopBottom},
				ElementViews:      []*mdl.ElementView{{ID: "user"}, {ID: "web"}, {ID: "db"}},
				RelationshipViews: []*mdl.RelationshipView{{ID: "order"}, {ID: "save"}},
			}}},
			Styles: &mdl.Styles{Elements: []*mdl.ElementStyle{
				{Tag: "person", Shape: mdl.ShapePerson},
				{Tag: "web", Background: "#e0f2fe", Stroke: "#0369a1"},
				{Tag: "store", Shape: mdl.ShapeCylinder, Background: "#fef3c7"},
			}},
		},
	}
}
