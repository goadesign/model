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

func TestContainerComponent(t *testing.T) {
	t.Parallel()
	container := Container{
		Components: Components{
			{Element: &Element{Name: "foo"}},
			{Element: &Element{Name: "bar"}},
			{Element: &Element{Name: "baz"}},
		},
	}
	tests := []struct {
		name string
		want *Component
	}{
		{name: "thud", want: nil},
		{name: "bar", want: container.Components[1]},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if got := container.Component(tt.name); got != tt.want {
				t.Errorf("got %#v, want %#v", got.Element, tt.want.Element)
			}
		})
	}
}

func TestAddComponent(t *testing.T) {
	t.Parallel()
	mElementFoo := Element{ID: "1", Name: "foo", Description: ""}
	mElementBar := Element{ID: "2", Name: "bar", Description: ""}
	componentFoo := Component{
		Element: &Element{ID: "1", Name: "foo", Description: "", DSLFunc: func() {
			fmt.Println("1")
			return
		}},
	}
	components := make([]*Component, 1)
	components[0] = &componentFoo
	OuterContainer := Container{
		Element:    &mElementFoo,
		Components: components,
		System:     nil,
	}
	InnerContainer := Container{
		Element:    &mElementBar,
		Components: components,
		System:     nil,
	}
	componentBar := Component{
		Element:   &Element{ID: "2", Name: "bar", Description: ""},
		Container: &InnerContainer,
	}
	componentFooPlus := Component{
		Element: &Element{ID: "3", Name: "foo", Description: "Description", Technology: "Golang", URL: "https://github.com/goadesign/model/", DSLFunc: func() {
			fmt.Printf("hello")
			return
		}},
	}

	tests := []struct {
		name          string
		component2Add *Component
		want          *Component
	}{
		{name: "foo", component2Add: &componentFoo, want: &componentFoo}, // already in container
		{name: "bar", component2Add: &componentBar, want: &componentBar}, // now appended
		{name: "foo", component2Add: &componentFooPlus, want: &componentFoo},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := OuterContainer.AddComponent(tt.component2Add)
			//if got := container.AddComponent(tt.component2Add); got != tt.want {
			if got != tt.want {
				t.Errorf("got %#v, want %#v", got.Element, tt.want.Element)
			}
		})
	}
}
