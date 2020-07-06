package expr

import (
	"fmt"
)

// Container represents a container.
type Container struct {
	*Element
	// Components list the components within the container.
	Components []*Component `json:"components,omitempty"`
	// System is the parent software system.
	System *SoftwareSystem `json:"-"`
}

// EvalName returns the generic expression name used in error messages.
func (c *Container) EvalName() string {
	if c.Name == "" {
		return "unnamed container"
	}
	return fmt.Sprintf("container %q", c.Name)
}
