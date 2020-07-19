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
		// AddImpliedRelationships adds implied relationships automatically.
		AddImpliedRelationships bool
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
	verr := new(eval.ValidationErrors)
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

// Finalize add all implied relationships if needed.
func (m *Model) Finalize() {
	if !m.AddImpliedRelationships {
		return
	}
	for _, elem := range Registry {
		r, ok := elem.(*Relationship)
		if !ok {
			continue
		}
		switch s := m.FindElement(r.SourceID).(type) {
		case *Container:
			m.addMissingRelationships(s.System.Element, r.FindDestination().ID, r)
		case *Component:
			m.addMissingRelationships(s.Container.Element, r.FindDestination().ID, r)
			m.addMissingRelationships(s.Container.System.Element, r.FindDestination().ID, r)
		}
	}
}

// addRelationshioIfNotExsists adds relationships from src to element with ID
// destID and its parents (container system software and component container)
// based on the properties of existing. It only adds a relationship if one
// doesn't already exist with the same description.
func (m *Model) addMissingRelationships(src *Element, destID string, existing *Relationship) {
	for _, r := range src.Rels {
		if r.FindDestination().ID == destID && r.Description == existing.Description {
			return
		}
	}
	r := existing.Dup()
	r.SourceID = src.ID
	r.DestinationID = destID
	src.Rels = append(src.Rels, r)

	// Add relationships to destination parents as well.
	switch e := m.FindElement(destID).(type) {
	case *Container:
		m.addMissingRelationships(src, e.System.ID, existing)
	case *Component:
		m.addMissingRelationships(src, e.Container.ID, existing)
		m.addMissingRelationships(src, e.Container.System.ID, existing)
	}
}

// FindElement retrieves the element with the given name or nil if there isn't
// one.
func (m *Model) FindElement(n string) ElementHolder {
	for _, x := range Registry {
		eh, ok := x.(ElementHolder)
		if !ok {
			continue
		}
		if eh.GetElement().Name == n {
			return eh
		}
	}
	return nil
}
