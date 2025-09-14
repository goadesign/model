package main

import (
	"os"
	"path/filepath"
	"testing"
)

// This test runs the full svg command against the basic example.
// Requires headless Chrome available in environment.
func TestSVGEndToEnd(t *testing.T) {
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
	// Cleanup generated file explicitly (t.TempDir will be removed automatically)
	if err := os.Remove(p); err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}
}
