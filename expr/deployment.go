package expr

import (
	"fmt"
)

type (
	// DeploymentEnvironment provides context to the other deployment expressions.
	DeploymentEnvironment struct {
		// Name of environment.
		Name string
	}

	// DeploymentNode describes a single deployment node.
	DeploymentNode struct {
		*Element
		Parent              *DeploymentNode
		Children            []*DeploymentNode
		InfrastructureNodes []*InfrastructureNode
		ContainerInstances  []*ContainerInstance
		Instances           *int
		Environment         string
	}

	// InfrastructureNode describes an infrastructure node.
	InfrastructureNode struct {
		*Element
		Parent      *DeploymentNode
		Environment string
	}

	// ContainerInstance describes an instance of a container.
	ContainerInstance struct {
		// cheating a bit: a ContainerInstance does not have a name,
		// description, technology or URL.
		*Element
		Parent       *DeploymentNode
		Container    *Container
		HealthChecks []*HealthCheck
		ContainerID  string
		InstanceID   int
		Environment  string
	}

	// InfrastructureNodes is a slice of infrastructure nodes that can be
	// converted into a slice of ElementHolder.
	InfrastructureNodes []*InfrastructureNode

	// ContainerInstances is a slice of container instances that can be
	// converted into a slice of ElementHolder.
	ContainerInstances []*ContainerInstance

	// HealthCheck is a HTTP-based health check.
	HealthCheck struct {
		Name     string
		URL      string
		Interval int
		Timeout  int
		Headers  map[string]string
	}
)

// EvalName returns the generic expression name used in error messages.
func (d *DeploymentEnvironment) EvalName() string {
	return fmt.Sprintf("deployment environment %q", d.Name)
}

// EvalName returns the generic expression name used in error messages.
func (d *DeploymentNode) EvalName() string { return fmt.Sprintf("deployment node %q", d.Name) }

// Finalize adds the 'Deployment Node' tag ands finalizes relationships.
func (d *DeploymentNode) Finalize() {
	d.PrefixTags("Element", "Deployment Node")
	d.Element.Finalize()
}

// Child returns the child deployment node with the given name if any,
// nil otherwise.
func (d *DeploymentNode) Child(name string) *DeploymentNode {
	for _, dd := range d.Children {
		if dd.Name == name {
			return dd
		}
	}
	return nil
}

// InfrastructureNode returns the infrastructure node with the given name if
// any, nil otherwise.
func (d *DeploymentNode) InfrastructureNode(name string) *InfrastructureNode {
	for _, i := range d.InfrastructureNodes {
		if i.Name == name {
			return i
		}
	}
	return nil
}

// ContainerInstance returns the container instance for the given container with
// the given instance ID if any, nil otherwise. container must be an instance of
// Container or the name of a container.
func (d *DeploymentNode) ContainerInstance(containerID string, instanceID int) *ContainerInstance {
	for _, ci := range d.ContainerInstances {
		if ci.ContainerID == containerID && ci.InstanceID == instanceID {
			return ci
		}
	}
	return nil
}

// AddChild adds the given child deployment node to the parent. If
// there is already a deployment node with the given name then AddChild
// merges both definitions. The merge algorithm:
//
//    * overrides the description, technology and URL if provided,
//    * merges any new tag or propery into the existing tags and properties,
//    * merges any new child deployment node into the existing children,
//    * merges any new container instance or infrastructure nodes into existing
//      ones.
//
// AddChild returns the new or merged deployment node.
func (d *DeploymentNode) AddChild(n *DeploymentNode) *DeploymentNode {
	existing := d.Child(n.Name)
	if existing == nil {
		Identify(n)
		d.Children = append(d.Children, n)
		return n
	}
	if n.Description != "" {
		existing.Description = n.Description
	}
	if n.Technology != "" {
		existing.Technology = n.Technology
	}
	if olddsl := existing.DSLFunc; olddsl != nil {
		existing.DSLFunc = func() { olddsl(); n.DSLFunc() }
	}
	return existing
}

// AddInfrastructureNode adds the given infrastructure node to the deployment
// node. If there is already an infrastructure node with the given name then
// AddInfrastructureNode merges both definitions. The merge algorithm:
//
//    * overrides the description, technology and URL if provided,
//    * merges any new tag or propery into the existing tags and properties.
//
// AddInfrastructureNode returns the new or merged infrastructure node.
func (d *DeploymentNode) AddInfrastructureNode(n *InfrastructureNode) *InfrastructureNode {
	existing := d.InfrastructureNode(n.Name)
	if existing == nil {
		Identify(n)
		d.InfrastructureNodes = append(d.InfrastructureNodes, n)
		return n
	}
	if n.Description != "" {
		existing.Description = n.Description
	}
	if n.Technology != "" {
		existing.Technology = n.Technology
	}
	if olddsl := existing.DSLFunc; olddsl != nil {
		existing.DSLFunc = func() { olddsl(); n.DSLFunc() }
	}
	return existing
}

// AddContainerInstance adds the given container instance to the deployment
// node. If there is already a container instance with the given container and
// instance ID then AddContainerInstance merges both definitions. The merge
// algorithm:
//
//    * overrides the description, technology and URL if provided,
//    * merges any new tag or propery into the existing tags and properties,
//    * merges any new health check into the existing health checks.
//
// AddContainerInstance returns the new or merged container instance.
func (d *DeploymentNode) AddContainerInstance(ci *ContainerInstance) *ContainerInstance {
	c := Registry[ci.ContainerID].(*Container)
	existing := d.ContainerInstance(c.ID, ci.InstanceID)
	if existing == nil {
		Identify(ci)
		d.ContainerInstances = append(d.ContainerInstances, ci)
		return ci
	}
	if ci.Description != "" {
		existing.Description = ci.Description
	}
	if ci.Technology != "" {
		existing.Technology = ci.Technology
	}
	existing.HealthChecks = append(existing.HealthChecks, ci.HealthChecks...)
	if olddsl := existing.DSLFunc; olddsl != nil {
		existing.DSLFunc = func() { olddsl(); ci.DSLFunc() }
	}
	return existing
}

// EvalName returns the generic expression name used in error messages.
func (i *InfrastructureNode) EvalName() string {
	return fmt.Sprintf("infrastructure node %q", i.Name)
}

// Finalize adds the 'Infrastructure Node' tag ands finalizes relationships.
func (i *InfrastructureNode) Finalize() {
	i.PrefixTags("Element", "Infrastructure Node")
	i.Element.Finalize()
}

// EvalName returns the generic expression name used in error messages.
func (ci *ContainerInstance) EvalName() string {
	n := "unknown container"
	if cn, ok := Registry[ci.ContainerID]; ok {
		n = fmt.Sprintf("container %q", cn.(*Container).Name)
	}
	return fmt.Sprintf("instance %d of %s", ci.InstanceID, n)
}

// Finalize adds the "Container Instance" tag if not present.
func (ci *ContainerInstance) Finalize() {
	ci.PrefixTags("Container Instance")
	ci.Element.Finalize()
}

// EvalName returns the generic expression name used in error messages.
func (hc *HealthCheck) EvalName() string {
	return fmt.Sprintf("health check %q", hc.Name)
}

// Elements returns a slice of ElementHolder that contains the elements of inf.
func (inf InfrastructureNodes) Elements() []ElementHolder {
	res := make([]ElementHolder, len(inf))
	for i, cc := range inf {
		res[i] = cc
	}
	return res
}

// Elements returns a slice of ElementHolder that contains the elements of ci.
func (ci ContainerInstances) Elements() []ElementHolder {
	res := make([]ElementHolder, len(ci))
	for i, cc := range ci {
		res[i] = cc
	}
	return res
}
