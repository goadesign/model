package expr

import (
	"fmt"
	"testing"
)

func TestPersonEvalName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, want string
	}{
		{name: "", want: "unnamed person"},
		{name: "foo", want: `person "foo"`},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			person := Person{
				Element: &Element{
					Name: tt.name,
				},
			}
			if got := person.EvalName(); got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestPersonFinalize(t *testing.T) {
	t.Parallel()
	person := Person{
		Element: &Element{
			Name: "foo",
		},
	}
	tests := []struct {
		pre  func()
		want string
	}{
		{want: ""},
		{pre: func() { person.Tags = "foo" }, want: "foo"},
		{pre: func() { person.Finalize() }, want: "Element,Person,foo"},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if tt.pre != nil {
				tt.pre()
			}
			if got := person.Tags; got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestPeopleElements(t *testing.T) {
	t.Parallel()
	people := People{
		{Element: &Element{Name: "foo"}},
		{Element: &Element{Name: "bar"}},
		{Element: &Element{Name: "baz"}},
	}
	if got := people.Elements(); len(got) != len(people) {
		t.Errorf("got %d, want %d", len(got), len(people))
	}
}
