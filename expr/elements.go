package expr

import (
	"bytes"
	"encoding/json"
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
		// Rels is the set of relationships from this element to other elements.
		Rels []*Relationship `json:"relationships,omitempty"`
	}

	// System represents a software system.
	System struct {
		Element
		// Containers list the containers within the software system.
		Containers []*Element `json:"containers,omitempty"`
	}

	// Container represents a container.
	Container struct {
		Element
		// Components list the components within the container.
		Components []*Element `json:"components,omitempty"`
	}

	// DeploymentNode describes a single deployment node.
	DeploymentNode struct {
		// ID of deployment node.
		ID string `json:"id"`
		// Name of deployment node.
		Name string `json:"name"`
		// Description of deployment node if any.
		Description string `json:"description"`
		// Technology used by deployment node if any.
		Technology string `json:"technology"`
		// Environment is the deployment environment in which this deployment node resides (e.g. "Development", "Live", etc).
		Environment string `json:"environment"`
		// Instances is the number of instances.
		Instances int `json:"instances"`
		// Tags attached to deployment node as comma separated list if any.
		Tags string `json:"tags"`
		// URL where more information about this deployment node can be found.
		URL string `json:"url"`
		// Children describe the child deployment nodes if any.
		Children []*DeploymentNode `json:"children"`
		// InfrastructureNodes describe the infrastructure nodes (load
		// balancers, firewall etc.)
		InfrastructureNodes []*Element `json:"infrastrctureNodes"`
		// ContainerInstances describe instances of containers deployed in
		// deployment node.
		ContainerInstances []*ContainerInstance `json:"containerInstances"`
		// Rels is the set of relationships from this deployment node to other elements.
		Rels []*Relationship `json:"relationships"`
	}

	// ContainerInstance describes an instance of a container.
	ContainerInstance struct {
		// Container that is instantiated.
		Container *Element
		// Tags attached to container instance as comma separated list if any.
		Tags string `json:"tags"`
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

// EvalName returns the generic expression name used in error messages.
func (w *Workspace) EvalName() string { return "Structurizr workspace" }

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
