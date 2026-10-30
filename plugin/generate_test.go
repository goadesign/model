package plugin_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
	_ "goa.design/model/examples/big_bank_plc/model"
	"goa.design/model/plugin"
)

// TestGenerateDescribesTheEvaluatedDesignOnce runs the design once, as goa gen
// does before calling plugins, and checks that the generated Structurizr
// workspace lists every view and every identified element once.
func TestGenerateDescribesTheEvaluatedDesignOnce(t *testing.T) {
	if err := eval.RunDSL(); err != nil {
		t.Fatal(err)
	}
	files, err := plugin.Generate("", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("got %d files, want 1", len(files))
	}
	if want := codegen.Gendir + "/structurizr/model.json"; files[0].Path != want {
		t.Errorf("got path %q, want %q", files[0].Path, want)
	}
	var buf bytes.Buffer
	for _, s := range files[0].SectionTemplates {
		if err := s.Write(&buf); err != nil {
			t.Fatal(err)
		}
	}
	var workspace struct {
		Model json.RawMessage
		Views map[string]json.RawMessage
	}
	if err := json.Unmarshal(buf.Bytes(), &workspace); err != nil {
		t.Fatal(err)
	}
	for kind, raw := range workspace.Views {
		var views []struct{ Key string }
		if err := json.Unmarshal(raw, &views); err != nil {
			continue // not a list of views, such as the view configuration
		}
		seen := map[string]bool{}
		for _, v := range views {
			if seen[v.Key] {
				t.Errorf("%s lists view %q more than once", kind, v.Key)
			}
			seen[v.Key] = true
		}
	}
	var model any
	if err := json.Unmarshal(workspace.Model, &model); err != nil {
		t.Fatal(err)
	}
	checkUniqueIDs(t, "model", model)
}

// checkUniqueIDs reports any list under v holding two objects with the same
// "id".
func checkUniqueIDs(t *testing.T, path string, v any) {
	t.Helper()
	switch x := v.(type) {
	case map[string]any:
		for k, child := range x {
			checkUniqueIDs(t, path+"."+k, child)
		}
	case []any:
		seen := map[any]bool{}
		for _, e := range x {
			if obj, ok := e.(map[string]any); ok {
				if id, ok := obj["id"]; ok {
					if seen[id] {
						t.Errorf("%s lists id %v more than once", path, id)
					}
					seen[id] = true
				}
			}
			checkUniqueIDs(t, path+"[]", e)
		}
	}
}
