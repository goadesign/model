package expr

import (
	"fmt"
	"testing"
)

func TestContainerEvalName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, want string
	}{
		{name: "", want: "unnamed container"},
		{name: "foo", want: `container "foo"`},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			container := Container{
				Element: &Element{
					Name: tt.name,
				},
			}
			if got := container.EvalName(); got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestContainerFinalize(t *testing.T) {
	t.Parallel()
	container := Container{
		Element: &Element{
			Name: "foo",
		},
	}
	tests := []struct {
		pre  func()
		want string
	}{
		{want: ""},
		{pre: func() { container.Tags = "foo" }, want: "foo"},
		{pre: func() { container.Finalize() }, want: "Element,Container,foo"},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if tt.pre != nil {
				tt.pre()
			}
			if got := container.Tags; got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestContainersElements(t *testing.T) {
	t.Parallel()
	containers := Containers{
		{Element: &Element{Name: "foo"}},
		{Element: &Element{Name: "bar"}},
		{Element: &Element{Name: "baz"}},
	}
	if got := containers.Elements(); len(got) != len(containers) {
		t.Errorf("got %d, want %d", len(got), len(containers))
	}
}
