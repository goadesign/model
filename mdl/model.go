package mdl

import (
	"encoding/json"
	"sort"
)

type (
	// Model describes a software architecture model.
	Model struct {
		// Enterprise associated with model if any.
		Enterprise *Enterprise `json:"enterprise,omitempty"`
		// People lists Person elements.
		People []*Person `json:"people,omitempty"`
		// Systems lists Software System elements.
		Systems []*SoftwareSystem `json:"softwareSystems,omitempty"`
		// DeploymentNodes list the deployment nodes.
		DeploymentNodes []*DeploymentNode `json:"deploymentNodes,omitempty"`
	}

	// Enterprise describes a named enterprise / organization.
	Enterprise struct {
		// Name of enterprise.
		Name string `json:"name"`
	}

	// alias to call original json.Unmarshal
	_model Model
)

// MarshalJSON guarantees the order of elements in generated JSON arrays that
// correspond to sets.
func (m *Model) MarshalJSON() ([]byte, error) {
	sort.Slice(m.People, func(i, j int) bool { return m.People[i].Name < m.People[j].Name })
	for _, p := range m.People {
		sort.Slice(p.Relationships, func(i, j int) bool { return p.Relationships[i].ID < p.Relationships[j].ID })
	}
	sort.Slice(m.Systems, func(i, j int) bool { return m.Systems[i].Name < m.Systems[j].Name })
	for _, sys := range m.Systems {
		sort.Slice(sys.Relationships, func(i, j int) bool { return sys.Relationships[i].ID < sys.Relationships[j].ID })
		sort.Slice(sys.Containers, func(i, j int) bool { return sys.Containers[i].Name < sys.Containers[j].Name })
		for _, c := range sys.Containers {
			sort.Slice(c.Relationships, func(i, j int) bool { return c.Relationships[i].ID < c.Relationships[j].ID })
			sort.Slice(c.Components, func(i, j int) bool { return c.Components[i].Name < c.Components[j].Name })
			for _, cmp := range c.Components {
				sort.Slice(cmp.Relationships, func(i, j int) bool { return cmp.Relationships[i].ID < cmp.Relationships[j].ID })
			}
		}
	}
	sortDeploymentNodes(m.DeploymentNodes)
	var mm _model = _model(*m)
	return json.Marshal(&mm)
}

func sortDeploymentNodes(nodes []*DeploymentNode) {
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].Name < nodes[j].Name })
	for _, node := range nodes {
		sortDeploymentNodes(node.Children)
		sort.Slice(node.InfrastructureNodes, func(i, j int) bool { return node.InfrastructureNodes[i].Name < node.InfrastructureNodes[j].Name })
		sort.Slice(node.ContainerInstances, func(i, j int) bool { return node.ContainerInstances[i].ID < node.ContainerInstances[j].ID })
	}
}
