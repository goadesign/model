package expr

import (
	"bytes"
	"encoding/json"
)

type (
	// Relationship describes a uni-directional relationship between two elements.
	Relationship struct {
		// Description of relationship if any.
		Description string `json:"description"`
		// SourceID is the ID of the source element.
		SourceID string `json:"sourceId"`
		// Target is the relationship target.
		Target *Element
		// Tags attached to relationship as comma separated list if any.
		Tags string `json:"tags"`
		// InteractionStyle describes whether the interaction is synchronous or asynchronous
		InteractionStyle InteractionStyleKind
	}

	// InteractionStyleKind is the enum for possible interaction styles.
	InteractionStyleKind int
)

const (
	// InteractionUndefined means no interaction style specified in design.
	InteractionUndefined InteractionStyleKind = iota
	// InteractionSynchronous describes a synchronous interaction.
	InteractionSynchronous
	// InteractionAsynchronous describes an asynchronous interaction.
	InteractionAsynchronous
)

// MarshalJSON replaces the constant value with the proper structurizr schema
// string value.
func (i InteractionStyleKind) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	switch i {
	case InteractionSynchronous:
		buf.WriteString("Synchronous")
	case InteractionAsynchronous:
		buf.WriteString("Asynchronous")
	case InteractionUndefined:
		buf.WriteString("Undefined")
	}
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

// UnmarshalJSON sets the constant from its JSON representation.
func (i *InteractionStyleKind) UnmarshalJSON(data []byte) error {
	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	switch val {
	case "Synchronous":
		*i = InteractionSynchronous
	case "Asynchronous":
		*i = InteractionAsynchronous
	case "Undefined":
		*i = InteractionUndefined
	}
	return nil
}
