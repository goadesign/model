package main

import (
	"encoding/xml"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

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
		dir:       outDir,
		port:      0,
		direction: "DOWN",
		timeout:   30_000_000_000, // 30s
		all:       true,
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
	// Cleanup generated files explicitly (t.TempDir will be removed automatically).
	for _, path := range []string{p, target} {
		if err := os.Remove(path); err != nil {
			t.Fatalf("cleanup %s: %v", path, err)
		}
	}
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
