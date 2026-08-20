package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"goa.design/model/mdl"
)

func TestCollectViewKeys(t *testing.T) {
	d := &mdl.Design{
		Views: &mdl.Views{
			LandscapeViews: []*mdl.LandscapeView{
				{ViewProps: &mdl.ViewProps{Key: "L1"}},
			},
			ContextViews: []*mdl.ContextView{
				{ViewProps: &mdl.ViewProps{Key: "C1"}},
			},
			ContainerViews: []*mdl.ContainerView{
				{ViewProps: &mdl.ViewProps{Key: "Ct1"}},
			},
			ComponentViews: []*mdl.ComponentView{
				{ViewProps: &mdl.ViewProps{Key: "Cm1"}},
			},
			DynamicViews: []*mdl.DynamicView{
				{ViewProps: &mdl.ViewProps{Key: "D1"}},
			},
			DeploymentViews: []*mdl.DeploymentView{
				{ViewProps: &mdl.ViewProps{Key: "Dp1"}},
			},
			FilteredViews: []*mdl.FilteredView{
				{Key: "F1"},
			},
		},
	}

	keys := collectViewKeys(d)
	if len(keys) != 7 {
		t.Fatalf("expected 7 keys, got %d: %v", len(keys), keys)
	}

	want := map[string]bool{"L1": true, "C1": true, "Ct1": true, "Cm1": true, "D1": true, "Dp1": true, "F1": true}
	for _, k := range keys {
		if !want[k] {
			t.Fatalf("unexpected key %q in %v", k, keys)
		}
		delete(want, k)
	}
	if len(want) != 0 {
		t.Fatalf("missing keys: %v", want)
	}
}

func TestNormalizeLayoutDirection(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  string
		shouldErr bool
	}{
		{name: "view default", input: "", expected: ""},
		{name: "explicit direction", input: "RIGHT", expected: "RIGHT"},
		{name: "case normalization", input: "left", expected: "LEFT"},
		{name: "invalid direction", input: "diagonal", shouldErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := normalizeLayoutDirection(test.input)
			if test.shouldErr {
				if err == nil {
					t.Fatal("expected an error")
				}
				return
			}
			if err != nil {
				t.Fatalf("normalize direction: %v", err)
			}
			if actual != test.expected {
				t.Fatalf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}

func TestHeadlessRenderURL(t *testing.T) {
	actual := headlessRenderURL(
		"http://127.0.0.1:1234",
		"digest",
		"AURA Services & Runtime",
		"RIGHT",
		true,
	)
	expected := "http://127.0.0.1:1234/headless.html?" +
		"compact=true&digest=digest&direction=RIGHT&view=AURA+Services+%26+Runtime"
	if actual != expected {
		t.Fatalf("expected %q, got %q", expected, actual)
	}
}

func TestChromedpExecNavigatesDirectPage(t *testing.T) {
	if !hasChrome() {
		t.Skip("skipping: Chrome/Chromium not available in PATH")
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, err := fmt.Fprint(w, "<!doctype html><html><body>ready</body></html>")
		if err != nil {
			t.Errorf("write direct page: %v", err)
		}
	}))
	defer server.Close()

	timeout := 30 * time.Second
	err := withChromedp(timeout, false, func(exec navigateExec) error {
		return exec(server.URL, timeout)
	})
	if err != nil {
		t.Fatalf("navigate direct page: %v", err)
	}
}
