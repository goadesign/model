package expr

import (
	"fmt"
)

type (
	// Component represents a component.
	Component struct {
		*Element
		Container *Container
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

// Finalize adds the 'Component' tag ands finalizes relationships.
func (c *Component) Finalize() {
	c.PrefixTags("Element", "Component")
}

// Elements returns a slice of ElementHolder that contains the elements of c.
func (cs Components) Elements() []ElementHolder {
	res := make([]ElementHolder, len(cs))
	for i, cc := range cs {
		res[i] = cc
	}
	return res
}
