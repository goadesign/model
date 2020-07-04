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

// EvalName is the qualified name of the DSL expression e.g. "service
// bottle".
func (m *Model) EvalName() string {
	return "model"
}
