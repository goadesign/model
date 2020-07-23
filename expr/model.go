package expr

import (
	"goa.design/goa/v3/eval"
)

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
		// AddImpliedRelationships adds implied relationships automatically.
		AddImpliedRelationships bool `json:"-"`
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
		switch s := Registry[r.SourceID].(type) {
		case *Container:
			m.addMissingRelationships(s.System.Element, r.FindDestination().ID, r)
		case *Component:
			m.addMissingRelationships(s.Container.Element, r.FindDestination().ID, r)
			m.addMissingRelationships(s.Container.System.Element, r.FindDestination().ID, r)
		}
	}
}

// Person returns the person with the given name if any, nil otherwise.
func (m *Model) Person(name string) *Person {
	for _, pp := range m.People {
		if pp.Name == name {
			return pp
		}
	}
	return nil
}

// SoftwareSystem returns the software system with the given name if any, nil
// otherwise.
func (m *Model) SoftwareSystem(name string) *SoftwareSystem {
	for _, s := range m.Systems {
		if s.Name == name {
			return s
		}
	}
	return nil
}

// DeploymentNode returns the deployment node with the given name if any, nil
// otherwise.
func (m *Model) DeploymentNode(name string) *DeploymentNode {
	for _, d := range m.DeploymentNodes {
		if d.Name == name {
			return d
		}
	}
	return nil
}

// AddPerson adds the given person to the model. If there is already a person
// with the given name then AddPerson merges both definitions. The merge
// algorithm:
//
//    * overrides the description, technology and URL if provided,
//    * merges any new tag or propery into the existing tags and properties,
//    * merges any new relationship into the existing relationships.
//
// AddPerson returns the new or merged person.
func (m *Model) AddPerson(p *Person) *Person {
	existing := m.Person(p.Name)
	if existing == nil {
		Identify(p)
		m.People = append(m.People, p)
		return p
	}
	if p.Description != "" {
		existing.Description = p.Description
	}
	if olddsl := existing.DSLFunc; olddsl != nil {
		existing.DSLFunc = func() { olddsl(); p.DSLFunc() }
	}
	return existing
}

// AddSystem adds the given software system to the model. If there is already a
// software system with the given name then AddSystem merges both definitions.
// The merge algorithm:
//
//    * overrides the description, technology and URL if provided,
//    * merges any new tag or propery into the existing tags and properties,
//    * merges any new relationship into the existing relationships,
//    * merges any new container into the existing containers.
//
// AddSystem returns the new or merged software system.
func (m *Model) AddSystem(s *SoftwareSystem) *SoftwareSystem {
	existing := m.SoftwareSystem(s.Name)
	if existing == nil {
		Identify(s)
		m.Systems = append(m.Systems, s)
		return s
	}
	if s.Description != "" {
		existing.Description = s.Description
	}
	if olddsl := existing.DSLFunc; olddsl != nil {
		existing.DSLFunc = func() { olddsl(); s.DSLFunc() }
	}
	return existing
}

// AddDeploymentNode adds the given deployment node to the model. If there is
// already a deployment node with the given name then AddDeploymentNode merges
// both definitions. The merge algorithm:
//
//    * overrides the description, technology and URL if provided,
//    * merges any new tag or propery into the existing tags and properties,
//    * merges any new relationship into the existing relationships,
//    * merges any new child deployment node into the existing children,
//    * merges any new container instance or infrastructure nodes into existing
//      ones.
//
// AddDeploymentNode returns the new or merged deployment node.
func (m *Model) AddDeploymentNode(d *DeploymentNode) *DeploymentNode {
	existing := m.DeploymentNode(d.Name)
	if existing == nil {
		Identify(d)
		m.DeploymentNodes = append(m.DeploymentNodes, d)
		return d
	}
	if d.Description != "" {
		existing.Description = d.Description
	}
	if d.Technology != "" {
		existing.Technology = d.Technology
	}
	if olddsl := existing.DSLFunc; olddsl != nil {
		existing.DSLFunc = func() { olddsl(); d.DSLFunc() }
	}
	return existing
}

// defined in the new person. It merges relationships.
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
	switch e := Registry[destID].(type) {
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
