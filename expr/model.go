package expr

import "goa.design/goa/v3/eval"

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

// Validate makes sure all element names are unique.
func (m *Model) Validate() error {
	var verr *eval.ValidationErrors
	known := make(map[string]struct{})
	for _, p := range m.People {
		if _, ok := known[p.Name]; ok {
			verr.Add(p, "name already in use")
		}
		known[p.Name] = struct{}{}
	}
	for _, s := range m.Systems {
		if _, ok := known[s.Name]; ok {
			verr.Add(s, "name already in use")
		}
		known[s.Name] = struct{}{}
		for _, c := range s.Containers {
			if _, ok := known[c.Name]; ok {
				verr.Add(s, "name already in use")
			}
			known[c.Name] = struct{}{}
			for _, cm := range c.Components {
				if _, ok := known[cm.Name]; ok {
					verr.Add(s, "name already in use")
				}
				known[cm.Name] = struct{}{}
			}
		}
	}
	return verr
}

// FindElement retrieves the element with the given name or nil if there isn't
// one.
func (m *Model) FindElement(n string) ElementHolder {
	for _, s := range m.Systems {
		if s.Name == n {
			return s
		}
		for _, c := range s.Containers {
			if c.Name == n {
				return c
			}
			for _, cm := range c.Components {
				if cm.Name == n {
					return cm
				}
			}
		}
	}
	return nil
}
