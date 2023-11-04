package expr

import (
	"fmt"
	"testing"
)

func TestComponentEvalName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, want string
	}{
		{name: "", want: "unnamed component"},
		{name: "foo", want: `component "foo"`},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			component := Component{
				Element: &Element{
					Name: tt.name,
				},
			}
			if got := component.EvalName(); got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestComponentFinalize(t *testing.T) {
	t.Parallel()
	component := Component{
		Element: &Element{
			Name: "foo",
		},
	}
	tests := []struct {
		pre  func()
		want string
	}{
		{want: ""},
		{pre: func() { component.Tags = "foo" }, want: "foo"},
		{pre: func() { component.Finalize() }, want: "Element,Component,foo"},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()
			if tt.pre != nil {
				tt.pre()
			}
			if got := component.Tags; got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestComponentsElements(t *testing.T) {
	t.Parallel()
	components := Components{
		{Element: &Element{Name: "foo"}},
		{Element: &Element{Name: "bar"}},
		{Element: &Element{Name: "baz"}},
	}
	if got := components.Elements(); len(got) != len(components) {
		t.Errorf("got %d, want %d", len(got), len(components))
	}
}
