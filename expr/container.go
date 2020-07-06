package expr

import (
	"fmt"
)

type (
	// Container represents a container.
	Container struct {
		*Element
		// Components list the components within the container.
		Components Components `json:"components,omitempty"`
		// System is the parent software system.
		System *SoftwareSystem `json:"-"`
	}

	// Containers is a slice of containers that can be easily
	// converted into a slice of ElementHolder.
	Containers []*Container
)

// EvalName returns the generic expression name used in error messages.
func (c *Container) EvalName() string {
	if c.Name == "" {
		return "unnamed container"
	}
	return fmt.Sprintf("container %q", c.Name)
}

// Elements returns a slice of ElementHolder that contains the elements of c.
func (c Containers) Elements() []ElementHolder {
	res := make([]ElementHolder, len(c))
	for i, cc := range c {
		res[i] = cc
	}
	return res
}
