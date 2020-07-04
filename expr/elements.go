package expr

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type (
	// Element describes an element.
	Element struct {
		// ID of element.
		ID string `json:"id"`
		// Name of element.
		Name string `json:"name"`
		// Description of element if any.
		Description string `json:"description,omitempty"`
		// Technology used by element if any.
		Technology string `json:"technology,omitempty"`
		// Tags attached to element as comma separated list if any.
		Tags string `json:"tags,omitempty"`
		// URL where more information about this element can be found.
		URL string `json:"url,omitempty"`
		// Location of element.
		Location LocationKind `json:"location"`
		// Set of arbitrary name-value properties (shown in diagram tooltips).
		Properties map[string]string `json:"properties"`
		// Rels is the set of relationships from this element to other elements.
		Rels []*Relationship `json:"relationships,omitempty"`
	}

	// System represents a software system.
	SoftwareSystem struct {
		Element
		// Containers list the containers within the software system.
		Containers []*Container `json:"containers,omitempty"`
	}

	// Container represents a container.
	Container struct {
		Element
		// Components list the components within the container.
		Components []*Component `json:"components,omitempty"`
	}

	// Component represents a component.
	Component Element

	// LocationKind is the enum for possible locations.
	LocationKind int
)

const (
	// LocationUndefined means no location specified in design.
	LocationUndefined LocationKind = iota
	// LocationInternal defines an element internal to the enterprise.
	LocationInternal
	// LocationExternal defines an element external to the enterprise.
	LocationExternal
)

// EvalName returns the generic expression name used in error messages.
func (w *Workspace) EvalName() string { return "Structurizr workspace" }

// EvalName returns the generic expression name used in error messages.
func (s *SoftwareSystem) EvalName() string { return fmt.Sprintf("system %q", s.Name) }

// EvalName returns the generic expression name used in error messages.
func (s *Container) EvalName() string { return fmt.Sprintf("container %q", s.Name) }

// EvalName returns the generic expression name used in error messages.
func (s *Component) EvalName() string { return fmt.Sprintf("component %q", s.Name) }

// MarshalJSON replaces the constant value with the proper structurizr schema
// string value.
func (l LocationKind) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	switch l {
	case LocationInternal:
		buf.WriteString("Internal")
	case LocationExternal:
		buf.WriteString("External")
	case LocationUndefined:
		buf.WriteString("Undefined")
	}
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

// UnmarshalJSON sets the constant from its JSON representation.
func (l *LocationKind) UnmarshalJSON(data []byte) error {
	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	switch val {
	case "Internal":
		*l = LocationInternal
	case "External":
		*l = LocationExternal
	case "Undefined":
		*l = LocationUndefined
	}
	return nil
}
