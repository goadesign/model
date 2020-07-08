package expr

import (
	"fmt"
)

type (
	// SoftwareSystem represents a software system.
	SoftwareSystem struct {
		*Element
		// Location of element.
		Location LocationKind `json:"location"`
		// Containers list the containers within the software system.
		Containers Containers `json:"containers,omitempty"`
	}

	// SoftwareSystems is a slice of software system that can be easily
	// converted into a slice of ElementHolder.
	SoftwareSystems []*SoftwareSystem
)

// EvalName returns the generic expression name used in error messages.
func (s *SoftwareSystem) EvalName() string {
	if s.Name == "" {
		return "unnamed software system"
	}
	return fmt.Sprintf("software system %q", s.Name)
}

// Elements returns a slice of ElementHolder that contains the elements of s.
func (s SoftwareSystems) Elements() []ElementHolder {
	res := make([]ElementHolder, len(s))
	for i, ss := range s {
		res[i] = ss
	}
	return res
}
