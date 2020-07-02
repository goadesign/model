package expr

type (
	// EnterpriseExpr describes a named enterprise / organization.
	EnterpriseExpr struct {
		// Name of enterprise.
		Name string `json:"name"`
	}

	// ModelExpr describes a software architecture model.
	ModelExpr struct {
		// Enterprise associated with model if any.
		Enterprise *EnterpriseExpr `json:"enterprise"`
		// People lists Person elements.
		People []*PersonExpr `json:"people"`
		// Systems lists Software System elements.
		Systems []*SystemExpr `json:"softwareSystems"`
		// DeploymentNodes list the deployment nodes.
		DeploymentNodes []*DeploymentNodeExpr `json:"deploymentNodes"`
	}
)

// EvalName is the qualified name of the DSL expression e.g. "service
// bottle".
func (m *ModelExpr) EvalName() string {
	return "model"
}
