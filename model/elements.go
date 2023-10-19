package model

import (
	"bytes"
	"encoding/json"
)

type (
	// Person represents a person.
	Person struct {
		// ID of element.
		ID string `json:"id"`
		// Name of element - not applicable to ContainerInstance.
		Name string `json:"name,omitempty"`
		// Description of element if any.
		Description string `json:"description,omitempty"`
		// Tags attached to element as comma separated list if any.
		Tags string `json:"tags,omitempty"`
		// URL where more information about this element can be found.
		URL string `json:"url,omitempty"`
		// Set of arbitrary name-value properties (shown in diagram tooltips).
		Properties map[string]string `json:"properties,omitempty"`
		// Relationships is the set of relationships from this element to other
		// elements.
		Relationships []*Relationship `json:"relationships,omitempty"`
		// Location of person.
		Location LocationKind `json:"location,omitempty"`
	}

	// SoftwareSystem represents a software system.
	SoftwareSystem struct {
		// ID of element.
		ID string `json:"id"`
		// Name of element - not applicable to ContainerInstance.
		Name string `json:"name,omitempty"`
		// Description of element if any.
		Description string `json:"description,omitempty"`
		// Tags attached to element as comma separated list if any.
		Tags string `json:"tags,omitempty"`
		// URL where more information about this element can be found.
		URL string `json:"url,omitempty"`
		// Set of arbitrary name-value properties (shown in diagram tooltips).
		Properties map[string]string `json:"properties,omitempty"`
		// Relationships is the set of relationships from this element to other
		// elements.
		Relationships []*Relationship `json:"relationships,omitempty"`
		// Location of element.
		Location LocationKind `json:"location,omitempty"`
		// Containers list the containers within the software system.
		Containers []*Container `json:"containers,omitempty"`
	}

	// Container represents a container.
	Container struct {
		// ID of element.
		ID string `json:"id"`
		// Name of element - not applicable to ContainerInstance.
		Name string `json:"name,omitempty"`
		// Description of element if any.
		Description string `json:"description,omitempty"`
		// Technology used by element if any - not applicable to Person.
		Technology string `json:"technology,omitempty"`
		// Tags attached to element as comma separated list if any.
		Tags string `json:"tags,omitempty"`
		// URL where more information about this element can be found.
		URL string `json:"url,omitempty"`
		// Set of arbitrary name-value properties (shown in diagram tooltips).
		Properties map[string]string `json:"properties,omitempty"`
		// Relationships is the set of relationships from this element to other
		// elements.
		Relationships []*Relationship `json:"relationships,omitempty"`
		// Components list the components within the container.
		Components []*Component `json:"components,omitempty"`
	}

	// Component represents a component.
	Component struct {
		// ID of element.
		ID string `json:"id"`
		// Name of element - not applicable to ContainerInstance.
		Name string `json:"name,omitempty"`
		// Description of element if any.
		Description string `json:"description,omitempty"`
		// Technology used by element if any - not applicable to Person.
		Technology string `json:"technology,omitempty"`
		// Tags attached to element as comma separated list if any.
		Tags string `json:"tags,omitempty"`
		// URL where more information about this element can be found.
		URL string `json:"url,omitempty"`
		// Set of arbitrary name-value properties (shown in diagram tooltips).
		Properties map[string]string `json:"properties,omitempty"`
		// Relationships is the set of relationships from this element to other
		// elements.
		Relationships []*Relationship `json:"relationships,omitempty"`
	}

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

// MarshalJSON replaces the constant value with the proper string value.
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
