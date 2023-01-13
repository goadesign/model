package mdl

type (
	// DeploymentNode describes a single deployment node.
	DeploymentNode struct {
		// ID of element.
		ID string `json:"id"`
		// Name of element - not applicable to ContainerInstance.
		Name string `json:"name,omitempty"`
		// Description of element if any.
		Description string `json:"description,omitempty"`
		// Technology used by element if any - not applicable to Person.
		Technology string `json:"technology,omitempty"`
		// Environment is the deployment environment in which this deployment
		// node resides (e.g. "Development", "Live", etc).
		Environment string `json:"environment"`
		// Instances is the number of instances.
		Instances *string `json:"instances,omitempty"`
		// Tags attached to element as comma separated list if any.
		Tags string `json:"tags,omitempty"`
		// URL where more information about this element can be found.
		URL string `json:"url,omitempty"`
		// Children describe the child deployment nodes if any.
		Children []*DeploymentNode `json:"children,omitempty"`
		// InfrastructureNodes describe the infrastructure nodes (load
		// balancers, firewall etc.)
		InfrastructureNodes []*InfrastructureNode `json:"infrastructureNodes,omitempty"`
		// ContainerInstances describe instances of containers deployed in
		// deployment node.
		ContainerInstances []*ContainerInstance `json:"containerInstances,omitempty"`
		// Set of arbitrary name-value properties (shown in diagram tooltips).
		Properties map[string]string `json:"properties,omitempty"`
		// Relationships is the set of relationships from this element to other
		// elements.
		Relationships []*Relationship `json:"relationships,omitempty"`
	}

	// InfrastructureNode describes an infrastructure node.
	InfrastructureNode struct {
		// ID of element.
		ID string `json:"id"`
		// Name of element - not applicable to ContainerInstance.
		Name string `json:"name,omitempty"`
		// Description of element if any.
		Description string `json:"description,omitempty"`
		// Technology used by element if any - not applicable to Person.
		Technology string `json:"technology,omitempty"`
		// Tags attached to element as comma separated list if any.
		Tags string `json:"tags,omitempty"`
		// URL where more information about this element can be found.
		URL string `json:"url,omitempty"`
		// Set of arbitrary name-value properties (shown in diagram tooltips).
		Properties map[string]string `json:"properties,omitempty"`
		// Relationships is the set of relationships from this element to other
		// elements.
		Relationships []*Relationship `json:"relationships,omitempty"`
		// Environment is the deployment environment in which this
		// infrastructure node resides (e.g. "Development", "Live",
		// etc).
		Environment string `json:"environment"`
	}

	// ContainerInstance describes an instance of a container.
	ContainerInstance struct {
		// ID of element.
		ID string `json:"id"`
		// Tags attached to element as comma separated list if any.
		Tags string `json:"tags,omitempty"`
		// URL where more information about this element can be found.
		URL string `json:"url,omitempty"`
		// Set of arbitrary name-value properties (shown in diagram tooltips).
		Properties map[string]string `json:"properties,omitempty"`
		// Relationships is the set of relationships from this element to other
		// elements.
		Relationships []*Relationship `json:"relationships,omitempty"`
		// ID of container that is instantiated.
		ContainerID string `json:"containerId"`
		// InstanceID is the number/index of this instance.
		InstanceID int `json:"instanceId"`
		// Environment is the deployment environment of this instance.
		Environment string `json:"environment"`
		// HealthChecks is the set of HTTP-based health checks for this
		// container instance.
		HealthChecks []*HealthCheck `json:"healthChecks,omitempty"`
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
