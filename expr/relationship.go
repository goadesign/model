package expr

import "bytes"

type (
	// RelationshipExpr describes a uni-directional relationship between two elements.
	RelationshipExpr struct {
		// Description of relationship if any.
		Description string `json:"description"`
		// SourceID is the ID of the source element.
		SourceID string `json:"sourceId"`
		// Target is the relationship target.
		Target *ElementExpr
		// Tags attached to relationship if any.
		Tags []string `json:"tags"`
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
	}
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}
