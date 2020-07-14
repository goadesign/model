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
		Enterprise *Enterprise `json:"enterprise,omitempty"`
		// People lists Person elements.
		People People `json:"people,omitempty"`
		// Systems lists Software System elements.
		Systems SoftwareSystems `json:"softwareSystems,omitempty"`
		// DeploymentNodes list the deployment nodes.
		DeploymentNodes []*DeploymentNode `json:"deploymentNodes,omitempty"`
	}
)

// EvalName is the qualified name of the DSL expression.
func (m *Model) EvalName() string { return "model" }
