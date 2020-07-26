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

// Finalize adds the 'SoftwareSystem' tag ands finalizes relationships.
func (s *SoftwareSystem) Finalize() {
	s.MergeTags("Element", "SoftwareSystem")
	s.Element.Finalize()
}

// Elements returns a slice of ElementHolder that contains the elements of s.
func (s SoftwareSystems) Elements() []ElementHolder {
	res := make([]ElementHolder, len(s))
	for i, ss := range s {
		res[i] = ss
	}
	return res
}

// Container returns the container with the given name if any, nil otherwise.
func (s *SoftwareSystem) Container(name string) *Container {
	for _, c := range s.Containers {
		if c.Name == name {
			return c
		}
	}
	return nil
}

// AddContainer adds the given container to the software system. If there is
// already a container with the given name then AddContainer merges both
// definitions. The merge algorithm:
//
//    * overrides the description, technology and URL if provided,
//    * merges any new tag or propery into the existing tags and properties,
//    * merges any new relationship into the existing relationships,
//    * merges any new component into the existing components.
//
// AddContainer returns the new or merged person.
func (s *SoftwareSystem) AddContainer(c *Container) *Container {
	existing := s.Container(c.Name)
	if existing == nil {
		Identify(c)
		s.Containers = append(s.Containers, c)
		return c
	}
	if c.Description != "" {
		existing.Description = c.Description
	}
	if c.Technology != "" {
		existing.Technology = c.Technology
	}
	if olddsl := existing.DSLFunc; olddsl != nil {
		existing.DSLFunc = func() { olddsl(); s.DSLFunc() }
	}
	return existing
}
