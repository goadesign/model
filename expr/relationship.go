package expr

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type (
	// Relationship describes a uni-directional relationship between two elements.
	Relationship struct {
		// ID of relationship.
		ID string `json:"id"`
		// Description of relationship if any.
		Description string `json:"description"`
		// Tags attached to relationship as comma separated list if any.
		Tags string `json:"tags,omitempty"`
		// URL where more information can be found.
		URL string `json:"url,omitempty"`
		// SourceID is the ID of the source element.
		SourceID string `json:"sourceId"`
		// DestinationID is ID the destination element.
		DestinationID string `json:"destinationId"`
		// Technology associated with relationship.
		Technology string `json:"technology,omitempty"`
		// InteractionStyle describes whether the interaction is synchronous or asynchronous
		InteractionStyle InteractionStyleKind `json:"interactionStyle"`
		// ID of container-container relationship upon which this container
		// instance-container instance relationship is based.
		LinkedRelationshipID string `json:"linkedRelationshipId,omitempty"`
		// Source element.
		Source *Element `json:"-"`
		// Destination element.
		Destination *Element `json:"-"`
		// DestinationName element name.
		DestinationName string `json:"-"`
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

// EvalName is the qualified name of the expression.
func (r *Relationship) EvalName() string {
	var src, dest = "unknown source", "unknown destination"
	if r.Source != nil {
		src = r.Source.Name
	}
	if r.FindDestination() != nil {
		dest = r.Destination.Name
	}
	return fmt.Sprintf("%s [%s -> %s]", r.Description, src, dest)
}

// Validate makes sure the named destination exists.
func (r *Relationship) Validate() error {
	if r.FindDestination() == nil {
		return fmt.Errorf("could not find relationship destination %q", r.DestinationName)
	}
	return nil
}

// Finalize computes the destinations when name is used to define relationship.
func (r *Relationship) Finalize() {
	r.MergeTags("Relationship")
	r.FindDestination()
}

// FindDestination computes the relationship destination.
func (r *Relationship) FindDestination() *Element {
	if r.Destination != nil {
		return r.Destination
	}
	srcDepl := false
	switch Registry[r.Source.ID].(type) {
	case *DeploymentNode, *InfrastructureNode, *ContainerInstance:
		srcDepl = true
	}
	for _, e := range Registry {
		eh, ok := e.(ElementHolder)
		if !ok {
			continue
		}
		destDepl := false
		switch e.(type) {
		case *DeploymentNode, *InfrastructureNode, *ContainerInstance:
			destDepl = true
		}
		if (srcDepl || destDepl) && (!srcDepl || !destDepl) {
			continue
		}
		ee := eh.GetElement()
		if ee.Name == r.DestinationName {
			r.Destination = ee
			r.DestinationID = ee.ID
			return ee
		}
	}
	return nil
}

// Dup creates a new relationship with identical description, tags, URL,
// technology and interaction style as r. Dup also creates a new ID for the
// result.
func (r *Relationship) Dup(newSrc, newDest string) *Relationship {
	dup := &Relationship{
		SourceID:         newSrc,
		DestinationID:    newDest,
		Description:      r.Description,
		Tags:             r.Tags,
		URL:              r.URL,
		Technology:       r.Technology,
		InteractionStyle: r.InteractionStyle,
	}
	Identify(dup)
	return dup
}

// MergeTags adds the given tags. It skips tags already present in e.Tags.
func (r *Relationship) MergeTags(tags ...string) {
	r.Tags = mergeTags(r.Tags, tags)
}

// MarshalJSON replaces the constant value with the proper string value.
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
