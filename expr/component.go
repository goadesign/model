package expr

import (
	"fmt"
)

type (
	// Component represents a component.
	Component struct {
		*Element
		// Container is the parent container.
		Container *Container `json:"-"`
	}

	// Components is a slice of components that can be easily converted into
	// a slice of ElementHolder.
	Components []*Component
)

// EvalName returns the generic expression name used in error messages.
func (c *Component) EvalName() string {
	if c.Name == "" {
		return "unnamed component"
	}
	return fmt.Sprintf("component %q", c.Name)
}

// Elements returns a slice of ElementHolder that contains the elements of c.
func (c Components) Elements() []ElementHolder {
	res := make([]ElementHolder, len(c))
	for i, cc := range c {
		res[i] = cc
	}
	return res
}
