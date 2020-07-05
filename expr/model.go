package expr

type (
	// Enterprise describes a named enterprise / organization.
	Enterprise struct {
		// Name of enterprise.
		Name string `json:"name"`
	}

	// Model describes a software architecture model.
	Model struct {
		// Enterprise associated with model if any.
		Enterprise *Enterprise `json:"enterprise"`
		// People lists Person elements.
		People []*Person `json:"people"`
		// Systems lists Software System elements.
		Systems []*SoftwareSystem `json:"softwareSystems"`
		// DeploymentNodes list the deployment nodes.
		DeploymentNodes []*DeploymentNode `json:"deploymentNodes"`
	}
)

// EvalName is the qualified name of the DSL expression.
func (m *Model) EvalName() string { return "model" }

// PeopleElements returns all the model people as a slice of *Element.
func (m *Model) PeopleElements() []*Element {
	res := make([]*Element, len(m.People))
	for i, p := range m.People {
		e := Element(*p)
		res[i] = &e
	}
	return res
}

// SystemElements returns all the model software systems as a slice of *Element.
func (m *Model) SystemElements() []*Element {
	res := make([]*Element, len(m.Systems))
	for i, s := range m.Systems {
		res[i] = &s.Element
	}
	return res
}
