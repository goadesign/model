package expr

import "fmt"

type (
	// DeploymentEnvironment provides context to the other deployment expressions.
	DeploymentEnvironment struct {
		// Name of environment.
		Name string
	}

	// DeploymentElement describes a deployment element.
	DeploymentElement struct {
		// ID of deployment element.
		ID string `json:"id"`
		// Name of deployment element.
		Name string `json:"name"`
		// Description of deployment element if any.
		Description string `json:"description"`
		// Technology used by deployment element if any.
		Technology string `json:"technology"`
		// Environment is the deployment environment in which this deployment
		// node resides (e.g. "Development", "Live", etc).
		Environment string `json:"environment"`
		// Tags attached to deployment node as comma separated list if any.
		Tags string `json:"tags"`
		// URL where more information about this deployment node can be found.
		URL string `json:"url"`
		// Set of arbitrary name-value properties (shown in diagram tooltips).
		Properties map[string]string `json:"properties"`
		// Rels is the set of relationships from this deployment node to other
		// elements.
		Rels []*Relationship `json:"relationships"`
	}

	// DeploymentNode describes a single deployment node.
	DeploymentNode struct {
		DeploymentElement
		// Instances is the number of instances.
		Instances int `json:"instances"`
		// Children describe the child deployment nodes if any.
		Children []*DeploymentNode `json:"children"`
		// InfrastructureNodes describe the infrastructure nodes (load
		// balancers, firewall etc.)
		InfrastructureNodes []*InfrastructureNode `json:"infrastrctureNodes"`
		// ContainerInstances describe instances of containers deployed in
		// deployment node.
		ContainerInstances []*ContainerInstance `json:"containerInstances"`
	}

	// InfrastructureNode describes an infrastructure node.
	InfrastructureNode DeploymentElement

	// ContainerInstance describes an instance of a container.
	ContainerInstance struct {
		// ID of container instance.
		ID string `json:"id"`
		// ID of container that is instantiated.
		ContainerID string `json:"containerId"`
		// InstanceID is the number/index of this instance.
		InstanceID int `json:"instanceId"`
		// Environment is the deployment environment of this instance.
		Environment string `json:"environment"`
		// Tags attached to container instance as comma separated list if any.
		Tags string `json:"tags"`
		// Set of arbitrary name-value properties (shown in diagram tooltips).
		Properties map[string]string `json:"properties"`
		// Rels is the set of relationships from this container to other elements.
		Rels []*Relationship `json:"relationships"`
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
func (c *ContainerInstance) EvalName() string {
	return fmt.Sprintf("instance %d of container with ID %s", c.InstanceID, c.ContainerID)
}
