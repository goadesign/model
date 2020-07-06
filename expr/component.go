package expr

import (
	"fmt"
)

// Component represents a component.
type Component struct {
	*Element
	// Container is the parent container.
	Container *Container `json:"-"`
}

// EvalName returns the generic expression name used in error messages.
func (c *Component) EvalName() string {
	if c.Name == "" {
		return "unnamed component"
	}
	return fmt.Sprintf("component %q", c.Name)
}
