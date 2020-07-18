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
		// Environment is the deployment environment in which this deployment
		// node resides (e.g. "Development", "Live", etc).
		Environment string `json:"environment"`
		// Instances is the number of instances.
		Instances *int `json:"instances,omitempty"`
		// Children describe the child deployment nodes if any.
		Children []*DeploymentNode `json:"children,omitempty"`
		// Parent is the parent deployment node if any.
		Parent *DeploymentNode `json:"-"`
		// InfrastructureNodes describe the infrastructure nodes (load
		// balancers, firewall etc.)
		InfrastructureNodes []*InfrastructureNode `json:"infrastrctureNodes,omitempty"`
		// ContainerInstances describe instances of containers deployed in
		// deployment node.
		ContainerInstances []*ContainerInstance `json:"containerInstances,omitempty"`
	}

	// InfrastructureNode describes an infrastructure node.
	InfrastructureNode struct {
		*Element
		// Parent deployment node.
		Parent *DeploymentNode `json:"-"`
		// Environment is the deployment environment in which this
		// infrastructure node resides (e.g. "Development", "Live",
		// etc).
		Environment string `json:"environment"`
	}

	// ContainerInstance describes an instance of a container.
	ContainerInstance struct {
		// cheating a bit: a ContainerInstance does not have a name,
		// description, technology or URL.
		*Element
		// Parent deployment node.
		Parent *DeploymentNode `json:"-"`
		// ID of container that is instantiated.
		ContainerID string `json:"containerId"`
		// InstanceID is the number/index of this instance.
		InstanceID int `json:"instanceId"`
		// Environment is the deployment environment of this instance.
		Environment string `json:"environment"`
		// HealthChecks is the set of HTTP-based health checks for this
		// container instance.
		HealthChecks []*HealthCheck `json:"healthChecks"`
	}

	// HealthCheck is a HTTP-based health check.
	HealthCheck struct {
		// Name of health check.
		Name string `json:"name"`
		// Health check URL/endpoint.
		URL string `json:"url"`
		// Polling interval, in seconds.
		Interval int `json:"interval"`
		// Timeout after which health check is deemed as failed, in milliseconds.
		Timeout int `json:"timeout"`
		// Set of name-value pairs corresponding to HTTP headers to be sent with request.
		Headers map[string]string `json:"headers"`
	}
)

// EvalName returns the generic expression name used in error messages.
func (d *DeploymentEnvironment) EvalName() string {
	return fmt.Sprintf("deployment environment %q", d.Name)
}

// EvalName returns the generic expression name used in error messages.
func (d *DeploymentNode) EvalName() string { return fmt.Sprintf("deployment node %q", d.Name) }

// EvalName returns the generic expression name used in error messages.
func (i *InfrastructureNode) EvalName() string { return fmt.Sprintf("infrastructure node %q", i.Name) }

// EvalName returns the generic expression name used in error messages.
func (ci *ContainerInstance) EvalName() string {
	n := "unknown container"
	if cn, ok := Registry[ci.ContainerID]; ok {
		n = fmt.Sprintf("container %q", cn.(*Container).Name)
	}
	return fmt.Sprintf("instance %d of %s", ci.InstanceID, n)
}

// Finalize removes the name value as it should not appear in the final JSON. It
// also adds all the implied relationships.
func (ci *ContainerInstance) Finalize() {
	ci.Name = ""
	c := Root.Model.FindElement(ci.ContainerID).(*Container)
	for _, r := range c.Rels {
		dc, ok := Root.Model.FindElement(r.DestinationID).(*Container)
		if !ok {
			continue
		}
		for _, e := range Registry {
			eci, ok := e.(*ContainerInstance)
			if !ok {
				continue
			}
			if eci.ContainerID == dc.ID {
				rc := &Relationship{
					Description:          r.Description,
					Tags:                 r.Tags,
					URL:                  r.URL,
					SourceID:             ci.ID,
					DestinationID:        eci.ID,
					Technology:           r.Technology,
					InteractionStyle:     r.InteractionStyle,
					LinkedRelationshipID: r.ID,
				}
				Identify(rc)
				ci.Rels = append(c.Rels, rc)
			}
		}
	}
}

// EvalName returns the generic expression name used in error messages.
func (hc *HealthCheck) EvalName() string {
	return fmt.Sprintf("health check %q", hc.Name)
}
