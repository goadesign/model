package expr

import (
	"fmt"
)

// SoftwareSystem represents a software system.
type SoftwareSystem struct {
	*Element
	// Containers list the containers within the software system.
	Containers []*Container `json:"containers,omitempty"`
}

// EvalName returns the generic expression name used in error messages.
func (s *SoftwareSystem) EvalName() string {
	if s.Name == "" {
		return "unnamed software system"
	}
	return fmt.Sprintf("software system %q", s.Name)
}
