package expr

import (
	"fmt"
	"strings"
)

type (
	// Container represents a container.
	Container struct {
		*Element
		Components Components
		System     *SoftwareSystem
	}

	// Containers is a slice of containers that can be easily
	// converted into a slice of ElementHolder.
	Containers []*Container
)

// ContainerTags lists the tags that are added to all containers.
var ContainerTags = []string{"Element", "Container"}

// EvalName returns the generic expression name used in error messages.
func (c *Container) EvalName() string {
	if c.Name == "" {
		return "unnamed container"
	}
	return fmt.Sprintf("container %q", c.Name)
}

// Finalize adds the 'Container' tag ands finalizes relationships.
func (c *Container) Finalize() {
	c.PrefixTags(ContainerTags...)
	c.Element.Finalize()
}

// Elements returns a slice of ElementHolder that contains the elements of c.
func (c Containers) Elements() []ElementHolder {
	res := make([]ElementHolder, len(c))
	for i, cc := range c {
		res[i] = cc
	}
	return res
}

// Component returns the component with the given name if any, nil otherwise.
func (c *Container) Component(name string) *Component {
	for _, cc := range c.Components {
		if cc.Name == name {
			return cc
		}
	}
	return nil
}

// AddComponent adds the given component to the container. If there is already a
// component with the given name then AddComponent merges both definitions. The
// merge algorithm:
//
//   - overrides the description, technology and URL if provided,
//   - merges any new tag or propery into the existing tags and properties,
//
// AddComponent returns the new or merged component.
func (c *Container) AddComponent(cmp *Component) *Component {
	existing := c.Component(cmp.Name)
	if existing == nil {
		Identify(cmp)
		c.Components = append(c.Components, cmp)
		return cmp
	}
	if cmp.Description != "" {
		existing.Description = cmp.Description
	}
	if cmp.Technology != "" {
		existing.Technology = cmp.Technology
	}
	if cmp.URL != "" {
		existing.URL = cmp.URL
	}
	existing.MergeTags(strings.Split(cmp.Tags, ",")...)
	if olddsl := existing.DSLFunc; olddsl != nil {
		existing.DSLFunc = func() { olddsl(); cmp.DSLFunc() }
	}
	return existing
}
