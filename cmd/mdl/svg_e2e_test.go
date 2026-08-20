package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	cdruntime "github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/jaschaephraim/lrserver"
)

type nodeOverflow struct {
	Name   string  `json:"name"`
	Top    float64 `json:"top"`
	Right  float64 `json:"right"`
	Bottom float64 `json:"bottom"`
	Left   float64 `json:"left"`
}

type verticalEdgeLabelOverlap struct {
	Source      string  `json:"source"`
	Destination string  `json:"destination"`
	LineX       float64 `json:"lineX"`
	LabelLeft   float64 `json:"labelLeft"`
	LabelRight  float64 `json:"labelRight"`
}

func hasChrome() bool {
	if os.Getenv("CHROME_BIN") != "" {
		return true
	}
	names := []string{"google-chrome", "chromium", "chromium-browser"}
	for _, n := range names {
		if _, err := exec.LookPath(n); err == nil {
			return true
		}
	}
	return false
}

// This test runs the full svg command against the basic example.
// Requires headless Chrome available in environment.
func TestSVGEndToEnd(t *testing.T) {
	if !hasChrome() {
		t.Skip("skipping: Chrome/Chromium not available in PATH")
	}

	outDir := t.TempDir()
	cfg := config{
		dir:     outDir,
		port:    0,
		timeout: 30_000_000_000, // 30s
		all:     true,
	}
	if err := runSVG("goa.design/model/examples/basic/model", cfg); err != nil {
		t.Fatalf("runSVG failed: %v", err)
	}
	// basic example defines SystemContext view
	p := filepath.Join(outDir, "SystemContext.svg")
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("missing generated svg: %v", err)
	}
	target := filepath.Join(outDir, "Container View.svg")
	if _, err := os.Stat(target); err != nil {
		t.Fatalf("missing linked container svg: %v", err)
	}
	links := svgLinks(t, p)
	if len(links) != 1 || links[0] != "Container%20View.svg" {
		t.Fatalf("expected one container-view link, got %v", links)
	}
	if orientation := inspectNodeOrientation(t, p); orientation != "horizontal" {
		t.Fatalf("expected view-defined left-to-right layout, got %s", orientation)
	}
	if overlaps := inspectVerticalEdgeLabelOverlaps(t, p); len(overlaps) > 0 {
		t.Fatalf("vertical relationship labels overlap their lines: %+v", overlaps)
	}

	overrideDir := t.TempDir()
	overrideConfig := config{
		dir:       overrideDir,
		port:      0,
		views:     SliceFlag{"SystemContext"},
		direction: "DOWN",
		timeout:   30 * time.Second,
	}
	if err := runSVG("goa.design/model/examples/basic/model", overrideConfig); err != nil {
		t.Fatalf("runSVG with direction override failed: %v", err)
	}
	overridePath := filepath.Join(overrideDir, "SystemContext.svg")
	if orientation := inspectNodeOrientation(t, overridePath); orientation != "vertical" {
		t.Fatalf("expected explicit top-to-bottom override, got %s", orientation)
	}

	// Cleanup generated files explicitly (t.TempDir will be removed automatically).
	for _, path := range []string{p, target} {
		if err := os.Remove(path); err != nil {
			t.Fatalf("cleanup %s: %v", path, err)
		}
	}
}

// TestManualEditAfterAutoLayout proves editor changes enter the complete manual
// contract instead of mixing changed nodes with stale automatic geometry.
func TestManualEditAfterAutoLayout(t *testing.T) {
	if !hasChrome() {
		t.Skip("skipping: Chrome/Chromium not available in PATH")
	}
	design, err := loadDesign("goa.design/model/examples/basic/model", false)
	if err != nil {
		t.Fatalf("load design: %v", err)
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	server := NewServer(design)
	httpServer := &http.Server{ReadHeaderTimeout: 3 * time.Second}
	outDir := t.TempDir()
	serverDone := make(chan error, 1)
	go func() {
		serverDone <- server.ServeOnListener(
			outDir,
			"",
			httpServer,
			http.NewServeMux(),
			listener,
		)
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
	pageURL := fmt.Sprintf(
		"http://%s/?id=%s",
		listener.Addr().String(),
		url.QueryEscape("SystemContext"),
	)
	const script = `(async () => {
		const graph = window.gdata;
		await graph.autoLayout();
		const node = graph.nodes()[0];
		graph.moveNode(node, node.x + 20, node.y, true);
		return graph.exportSVG().length;
	})()`
	var svgLength int
	if err := chromedp.Run(
		testContext,
		chromedp.Navigate(pageURL),
		chromedp.Poll(
			`window.gdata && window.gdata.nodes().every(node => node.ref?.isConnected)`,
			nil,
		),
		chromedp.Evaluate(
			script,
			&svgLength,
			func(params *cdruntime.EvaluateParams) *cdruntime.EvaluateParams {
				return params.WithAwaitPromise(true)
			},
		),
	); err != nil {
		t.Fatalf("auto layout, edit, and export: %v", err)
	}
	if svgLength == 0 {
		t.Fatal("manual export is empty")
	}
}

// TestSVGHeadlessDoesNotConnectToLiveReload verifies automated rendering is
// isolated from interactive editor reloads sent on the process-wide port.
func TestSVGHeadlessDoesNotConnectToLiveReload(t *testing.T) {
	if !hasChrome() {
		t.Skip("skipping: Chrome/Chromium not available in PATH")
	}

	reloadServer := lrserver.New(lrserver.DefaultName, lrserver.DefaultPort)
	serverDone := make(chan error, 1)
	go func() {
		serverDone <- reloadServer.ListenAndServe()
	}()
	if err := waitForTCP("127.0.0.1:35729", 5*time.Second); err != nil {
		select {
		case serverErr := <-serverDone:
			t.Skipf("live-reload port unavailable: %v", serverErr)
		default:
			t.Fatalf("start live-reload server: %v", err)
		}
	}
	t.Cleanup(func() {
		if err := reloadServer.Close(); err != nil {
			t.Errorf("close live-reload server: %v", err)
		}
	})

	var connected atomic.Bool
	stopObserving := make(chan struct{})
	observerDone := make(chan struct{})
	go func() {
		defer close(observerDone)
		ticker := time.NewTicker(time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if reloadServer.IsConnected() {
					connected.Store(true)
				}
			case <-stopObserving:
				return
			}
		}
	}()

	outDir := t.TempDir()
	cfg := config{
		dir:     outDir,
		port:    0,
		timeout: 30 * time.Second,
		views:   SliceFlag{"SystemContext"},
	}
	err := runSVG("goa.design/model/examples/basic/model", cfg)
	close(stopObserving)
	<-observerDone
	if err != nil {
		t.Fatalf("runSVG failed: %v", err)
	}
	if connected.Load() {
		t.Fatal("headless renderer connected to the interactive LiveReload server")
	}
}

func TestSVGNodeTextFits(t *testing.T) {
	if !hasChrome() {
		t.Skip("skipping: Chrome/Chromium not available in PATH")
	}

	outDir := t.TempDir()
	cfg := config{
		dir:       outDir,
		port:      0,
		direction: "DOWN",
		timeout:   30 * time.Second,
		all:       true,
	}
	if err := runSVG("goa.design/model/examples/text_fit/model", cfg); err != nil {
		t.Fatalf("runSVG failed: %v", err)
	}

	path := filepath.Join(outDir, "Text Fit.svg")
	overflows, height, boundaryIntersections, groupCount := inspectNodeTextFit(t, path)
	if len(overflows) > 0 {
		t.Fatalf("node text exceeds its border: %+v", overflows)
	}
	if height <= 180 {
		t.Fatalf("expected long-content node to grow beyond 180px, got %.1f", height)
	}
	if groupCount < 2 {
		t.Fatalf("expected at least two system boundaries, got %d", groupCount)
	}
	if len(boundaryIntersections) > 0 {
		t.Fatalf("system boundaries intersect: %v", boundaryIntersections)
	}
}

func TestSVGEdgeLabelsAvoidElements(t *testing.T) {
	if !hasChrome() {
		t.Skip("skipping: Chrome/Chromium not available in PATH")
	}

	outDir := t.TempDir()
	cfg := config{
		dir:       outDir,
		port:      0,
		direction: "DOWN",
		timeout:   30 * time.Second,
		all:       true,
	}
	if err := runSVG("goa.design/model/examples/label_collision/model", cfg); err != nil {
		t.Fatalf("runSVG failed: %v", err)
	}

	path := filepath.Join(outDir, "Label Collision.svg")
	if collisions := inspectEdgeLabelCollisions(t, path); len(collisions) > 0 {
		t.Fatalf("relationship labels overlap diagram elements: %v", collisions)
	}
}

func inspectNodeTextFit(t *testing.T, path string) ([]nodeOverflow, float64, []string, int) {
	t.Helper()

	testContext, cleanup := newChromeContext(t)
	defer cleanup()

	var overflows []nodeOverflow
	var height float64
	var boundaryIntersections []string
	var groupCount int
	fileURL := (&url.URL{Scheme: "file", Path: path}).String()
	const overflowScript = `(() => {
		const tolerance = 0.5;
		return [...document.querySelectorAll("g.node")].flatMap((node) => {
			const border = node.querySelector(".nodeBorder");
			const texts = [...node.querySelectorAll("text")];
			if (!border || texts.length === 0) return [];
			const rect = border.getBoundingClientRect();
			const boxes = texts.map((text) => text.getBoundingClientRect());
			const left = Math.min(...boxes.map((box) => box.left));
			const top = Math.min(...boxes.map((box) => box.top));
			const right = Math.max(...boxes.map((box) => box.right));
			const bottom = Math.max(...boxes.map((box) => box.bottom));
			const overflow = {
				name: texts[0].textContent.trim(),
				top: Math.max(0, rect.top - top),
				right: Math.max(0, right - rect.right),
				bottom: Math.max(0, bottom - rect.bottom),
				left: Math.max(0, rect.left - left),
			};
			return overflow.top > tolerance || overflow.right > tolerance ||
				overflow.bottom > tolerance || overflow.left > tolerance
				? [overflow]
				: [];
		});
	})()`
	const boundaryScript = `(() => {
		const groups = [...document.querySelectorAll("g.group")].map((group) => ({
			name: group.querySelector("text")?.textContent.trim() || group.id || "group",
			rect: group.querySelector("rect").getBoundingClientRect(),
		}));
		const intersections = [];
		for (let i = 0; i < groups.length; i++) {
			for (let j = i + 1; j < groups.length; j++) {
				const width = Math.max(0,
					Math.min(groups[i].rect.right, groups[j].rect.right) -
					Math.max(groups[i].rect.left, groups[j].rect.left));
				const height = Math.max(0,
					Math.min(groups[i].rect.bottom, groups[j].rect.bottom) -
					Math.max(groups[i].rect.top, groups[j].rect.top));
				if (width > 0.5 && height > 0.5) {
					intersections.push(groups[i].name + " / " + groups[j].name);
				}
			}
		}
		return intersections;
	})()`

	if err := chromedp.Run(testContext,
		chromedp.Navigate(fileURL),
		chromedp.WaitVisible("g.node", chromedp.ByQuery),
		chromedp.Evaluate(overflowScript, &overflows),
		chromedp.Evaluate(`document.querySelector("g.node .nodeBorder").getBBox().height`, &height),
		chromedp.Evaluate(boundaryScript, &boundaryIntersections),
		chromedp.Evaluate(`document.querySelectorAll("g.group").length`, &groupCount),
	); err != nil {
		t.Fatalf("inspect generated SVG: %v", err)
	}

	return overflows, height, boundaryIntersections, groupCount
}

func inspectEdgeLabelCollisions(t *testing.T, path string) []string {
	t.Helper()
	testContext, cleanup := newChromeContext(t)
	defer cleanup()

	var collisions []string
	fileURL := (&url.URL{Scheme: "file", Path: path}).String()
	const collisionScript = `(() => {
		const tolerance = 1;
		const overlaps = (first, second) =>
			Math.min(first.right, second.right) - Math.max(first.left, second.left) > tolerance &&
			Math.min(first.bottom, second.bottom) - Math.max(first.top, second.top) > tolerance;
		const nodes = [...document.querySelectorAll("g.node")].map((node) => ({
			name: node.querySelector("[data-field=name]")?.textContent.trim() || node.id,
			rect: node.querySelector(".nodeBorder")?.getBoundingClientRect(),
		})).filter((node) => node.rect);
		const labels = [...document.querySelectorAll("g.edge")].map((edge) => ({
			name: edge.querySelector(":scope > text")?.textContent.trim() || edge.id,
			rect: edge.querySelector(":scope > rect")?.getBoundingClientRect(),
		})).filter((label) => label.rect);
		const collisions = labels.flatMap((label) =>
			nodes.filter((node) => overlaps(label.rect, node.rect))
				.map((node) => label.name + " / node: " + node.name));
		for (let first = 0; first < labels.length; first++) {
			for (let second = first + 1; second < labels.length; second++) {
				if (overlaps(labels[first].rect, labels[second].rect)) {
					collisions.push(labels[first].name + " / label: " + labels[second].name);
				}
			}
		}
		return collisions;
	})()`
	if err := chromedp.Run(testContext,
		chromedp.Navigate(fileURL),
		chromedp.WaitVisible("g.edge", chromedp.ByQuery),
		chromedp.Evaluate(collisionScript, &collisions),
	); err != nil {
		t.Fatalf("inspect relationship label collisions: %v", err)
	}
	return collisions
}

func inspectNodeOrientation(t *testing.T, path string) string {
	t.Helper()
	testContext, cleanup := newChromeContext(t)
	defer cleanup()

	var orientation string
	fileURL := (&url.URL{Scheme: "file", Path: path}).String()
	const orientationScript = `(() => {
		const nodes = [...document.querySelectorAll("g.node")];
		if (nodes.length < 2) return "unknown";
		const first = nodes[0].getBoundingClientRect();
		const second = nodes[1].getBoundingClientRect();
		const deltaX = Math.abs(
			(first.left + first.right) / 2 - (second.left + second.right) / 2);
		const deltaY = Math.abs(
			(first.top + first.bottom) / 2 - (second.top + second.bottom) / 2);
		return deltaX > deltaY ? "horizontal" : "vertical";
	})()`
	if err := chromedp.Run(testContext,
		chromedp.Navigate(fileURL),
		chromedp.WaitVisible("g.node", chromedp.ByQuery),
		chromedp.Evaluate(orientationScript, &orientation),
	); err != nil {
		t.Fatalf("inspect node orientation: %v", err)
	}
	return orientation
}

func inspectVerticalEdgeLabelOverlaps(t *testing.T, path string) []verticalEdgeLabelOverlap {
	t.Helper()

	testContext, cleanup := newChromeContext(t)
	defer cleanup()

	var overlaps []verticalEdgeLabelOverlap
	fileURL := (&url.URL{Scheme: "file", Path: path}).String()
	const overlapScript = `(() => {
		const minimumGap = 4;
		return [...document.querySelectorAll("g.edge[data-from][data-to]")].flatMap((edge) => {
			const source = document.getElementById(edge.dataset.from);
			const destination = document.getElementById(edge.dataset.to);
			const label = edge.querySelector(":scope > rect");
			if (!source || !destination || !label) return [];

			const sourceRect = source.getBoundingClientRect();
			const destinationRect = destination.getBoundingClientRect();
			const sourceX = (sourceRect.left + sourceRect.right) / 2;
			const destinationX = (destinationRect.left + destinationRect.right) / 2;
			const sourceY = (sourceRect.top + sourceRect.bottom) / 2;
			const destinationY = (destinationRect.top + destinationRect.bottom) / 2;
			if (Math.abs(sourceX - destinationX) > 0.5 ||
				Math.abs(sourceY - destinationY) <= 0.5) {
				return [];
			}

			const labelRect = label.getBoundingClientRect();
			if (labelRect.left - sourceX >= minimumGap ||
				sourceX - labelRect.right >= minimumGap) {
				return [];
			}
			return [{
				source: edge.dataset.from,
				destination: edge.dataset.to,
				lineX: sourceX,
				labelLeft: labelRect.left,
				labelRight: labelRect.right,
			}];
		});
	})()`
	if err := chromedp.Run(testContext,
		chromedp.Navigate(fileURL),
		chromedp.WaitVisible("g.edge", chromedp.ByQuery),
		chromedp.Evaluate(overlapScript, &overlaps),
	); err != nil {
		t.Fatalf("inspect vertical relationship labels: %v", err)
	}
	return overlaps
}

// waitForTCP waits until a local test server accepts connections.
func waitForTCP(address string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		connection, err := net.DialTimeout("tcp", address, 100*time.Millisecond)
		if err == nil {
			if closeErr := connection.Close(); closeErr != nil {
				return closeErr
			}
			return nil
		}
		time.Sleep(10 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for %s", address)
}

func newChromeContext(t *testing.T) (context.Context, func()) {
	t.Helper()

	options := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-setuid-sandbox", true),
	)
	if chromeBin := os.Getenv("CHROME_BIN"); chromeBin != "" {
		options = append(options, chromedp.ExecPath(chromeBin))
	}
	allocatorContext, cancelAllocator := chromedp.NewExecAllocator(context.Background(), options...)
	browserContext, cancelBrowser := chromedp.NewContext(allocatorContext)
	testContext, cancelTest := context.WithTimeout(browserContext, 30*time.Second)
	cleanup := func() {
		cancelTest()
		cancelBrowser()
		cancelAllocator()
	}
	return testContext, cleanup
}

func svgLinks(t *testing.T, path string) []string {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open svg: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			t.Errorf("close svg: %v", err)
		}
	}()

	var links []string
	decoder := xml.NewDecoder(file)
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			return links
		}
		if err != nil {
			t.Fatalf("decode svg: %v", err)
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "a" {
			continue
		}
		for _, attr := range start.Attr {
			if attr.Name.Local == "href" {
				links = append(links, attr.Value)
			}
		}
	}
}
