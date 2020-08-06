package expr

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestMerge(t *testing.T) {
	cases := []struct {
		Name   string
		Target *Workspace
		Merged *Workspace

		Expected *Workspace
	}{{
		Name:     "name only",
		Target:   &Workspace{},
		Merged:   &Workspace{Name: "foo"},
		Expected: &Workspace{Name: "foo"},
	}, {
		Name:     "override",
		Target:   &Workspace{Name: "old"},
		Merged:   &Workspace{Name: "foo"},
		Expected: &Workspace{Name: "foo"},
	}, {
		Name:     "existing fields",
		Target:   &Workspace{ID: 42, Description: "desc"},
		Merged:   &Workspace{ID: 43, Name: "foo"},
		Expected: &Workspace{ID: 42, Name: "foo", Description: "desc"},
	}, {
		Name:     "nested",
		Target:   &Workspace{ID: 42},
		Merged:   &Workspace{Model: &Model{Enterprise: &Enterprise{Name: "ent"}}},
		Expected: &Workspace{ID: 42, Model: &Model{Enterprise: &Enterprise{Name: "ent"}}},
	}, {
		Name:     "deep override",
		Target:   &Workspace{Model: &Model{Enterprise: &Enterprise{Name: "old"}}},
		Merged:   &Workspace{Model: &Model{Enterprise: &Enterprise{Name: "new"}}},
		Expected: &Workspace{Model: &Model{Enterprise: &Enterprise{Name: "new"}}},
	}, {
		Name:     "deep merge",
		Target:   &Workspace{Model: &Model{Enterprise: &Enterprise{Name: "old"}, People: People{{Element: &Element{Name: "jane"}}}}},
		Merged:   &Workspace{Model: &Model{Enterprise: &Enterprise{Name: "new"}, Systems: SoftwareSystems{{Element: &Element{Name: "sys"}}}}},
		Expected: &Workspace{Model: &Model{Enterprise: &Enterprise{Name: "new"}, People: People{{Element: &Element{Name: "jane"}}}, Systems: SoftwareSystems{{Element: &Element{Name: "sys"}}}}},
	}, {
		Name:     "deep merge 2",
		Target:   &Workspace{Model: &Model{Enterprise: &Enterprise{Name: "old"}, Systems: SoftwareSystems{{Element: &Element{Name: "sys1"}}}}},
		Merged:   &Workspace{Model: &Model{Enterprise: &Enterprise{Name: "new"}, Systems: SoftwareSystems{{Element: &Element{Name: "sys2"}, Containers: Containers{{Element: &Element{Name: "cont"}}}}}}},
		Expected: &Workspace{Model: &Model{Enterprise: &Enterprise{Name: "new"}, Systems: SoftwareSystems{{Element: &Element{Name: "sys1"}}, {Element: &Element{Name: "sys2"}, Containers: Containers{{Element: &Element{Name: "cont"}}}}}}},
	}}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			c.Target.Merge(c.Merged)
			if !reflect.DeepEqual(c.Target, c.Expected) {
				b, err := json.MarshalIndent(c.Target, "", "  ")
				if err != nil {
					t.Fatalf("failed to serialized merged: %s", err)
				}
				t.Errorf("merged and expected differ, merged:\n%s", string(b))
			}
		})
	}
}
