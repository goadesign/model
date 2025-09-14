package main

import (
	"testing"

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
