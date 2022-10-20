package dsl

import (
	"strings"

	"goa.design/goa/v3/eval"

	"goa.design/model/expr"
)

// DeploymentEnvironment defines a deployment environment (e.g. development,
// production).
//
// DeploymentEnvironment must appear in a Design expression.
//
// DeploymentEnvironment accepts two arguments: the environment name and a DSL
// function used to describe the nodes within the environment.
//
// Example:
//
//	var _ = Design(func() {
//	     DeploymentEnvironment("production", func() {
//	         DeploymentNode("AppServer", "Application server", "Go and Goa v3")
//	         InfrastructureNote("Router", "External traffic router", "AWS Route 53")
//	         ContainerInstance(Container)
//	     })
//	 })
func DeploymentEnvironment(name string, dsl func()) {
	_, ok := eval.Current().(*expr.Design)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	env := &expr.DeploymentEnvironment{Name: name}
	eval.Execute(dsl, env)
}

// DeploymentNode defines a deployment node. Deployment nodes can be
// nested, so a deployment node can contain other deployment nodes.
// A deployment node can also contain InfrastructureNode and
// ContainerInstance elements.
//
// DeploymentNode must appear in a DeploymentEnvironment or DeploymentNode
// expression.
//
// DeploymentNode takes 1 to 4 arguments. The first argument is the node name.
// The name may be optionally followed by a description. If a description is set
// then it may be followed by the technology details used by the component.
// Finally DeploymentNode may take a func() as last argument to define
// additional properties of the component.
//
// Usage:
//
//	DeploymentNode("<name>")
//
//	DeploymentNode("<name>", "[description]")
//
//	DeploymentNode("<name>", "[description]", "[technology]")
//
//	DeploymentNode("<name>", func())
//
//	DeploymentNode("<name>", "[description]", func())
//
//	DeploymentNode("<name>", "[description]", "[technology]", func())
//
// Example:
//
//	var _ = Design(func() {
//	    DeploymentEnvironment("Production", func() {
//	        DeploymentNode("US", "US shard", func() {
//	            Tag("shard")
//	            Instances(3)
//	            URL("https://goa.design/docs/shard")
//	            InfrastructureNode("Gateway", "US gateway", func() {
//	                Tag("gateway")
//	                URL("https://goa.design/docs/shards/us")
//	            })
//	            ContainerInstance("Container", "System", func() {
//	                Tag("service")
//	                InstanceID(1)
//	                HealthCheck("check", func() {
//	                    URL("https://goa.design/health")
//	                    Interval(10)
//	                    Timeout(1000)
//	                })
//	            })
//	            DeploymentNode("Cluster", "K8 cluster", func() {
//	                // ...
//	            })
//	        })
//	    })
//	})
func DeploymentNode(name string, args ...interface{}) *expr.DeploymentNode {
	if strings.Contains(name, "/") {
		eval.ReportError("DeploymentNode: name cannot include slashes")
	}
	var (
		parent *expr.DeploymentNode
		env    string
	)
	switch d := eval.Current().(type) {
	case *expr.DeploymentEnvironment:
		env = d.Name
	case *expr.DeploymentNode:
		env = d.Environment
		parent = d
	default:
		eval.IncompatibleDSL()
		return nil
	}
	description, technology, dsl, err := parseElementArgs(args...)
	if err != nil {
		eval.ReportError("DeploymentNode: " + err.Error())
		return nil
	}
	one := 1
	node := &expr.DeploymentNode{
		Element: &expr.Element{
			Name:        name,
			Description: description,
			Technology:  technology,
			DSLFunc:     dsl,
		},
		Instances:   &one,
		Environment: env,
		Parent:      parent,
	}
	if parent != nil {
		return parent.AddChild(node)
	}
	return expr.Root.Model.AddDeploymentNode(node)
}

// InfrastructureNode defines an infrastructure node, typically something like a
// load balancer, firewall, DNS service, etc.
//
// InfrastructureNode must appear in a DeploymentNode expression.
//
// InfrastructureNode takes 1 to 4 arguments. The first argument is the
// infrastructure node name. The name may be optionally followed by a
// description. If a description is set then it may be followed by the
// technology details used by the component. Finally InfrastructureNode may take
// a func() as last argument to define additional properties of the component.
//
// Usage:
//
//	InfrastructureNode("<name>")
//
//	InfrastructureNode("<name>", "[description]")
//
//	InfrastructureNode("<name>", "[description]", "[technology]")
//
//	InfrastructureNode("<name>", func())
//
//	InfrastructureNode("<name>", "[description]", func())
//
//	InfrastructureNode("<name>", "[description]", "[technology]", func())
//
// Example:
//
//	var _ = Design(func() {
//	    DeploymentEnvironment("Production", func() {
//	        DeploymentNode("US", "US shard", func() {
//	            InfrastructureNode("Gateway", "US gateway", func() {
//	                Tag("gateway")
//	                URL("https://goa.design/docs/shards/us")
//	            })
//	        })
//	    })
//	})
func InfrastructureNode(name string, args ...interface{}) *expr.InfrastructureNode {
	d, ok := eval.Current().(*expr.DeploymentNode)
	if !ok {
		eval.IncompatibleDSL()
		return nil
	}
	if strings.Contains(name, "/") {
		eval.ReportError("InfrastructureNode: name cannot include slashes")
	}
	description, technology, dsl, err := parseElementArgs(args...)
	if err != nil {
		eval.ReportError("InfrastructureNode: " + err.Error())
		return nil
	}
	node := &expr.InfrastructureNode{
		Element: &expr.Element{
			Name:        name,
			Description: description,
			Technology:  technology,
			DSLFunc:     dsl,
		},
		Environment: d.Environment,
		Parent:      d,
	}
	return d.AddInfrastructureNode(node)
}

// ContainerInstance defines an instance of the specified container that is
// deployed on the parent deployment node.
//
// ContainerInstance must appear in a DeploymentNode expression.
//
// ContainerInstance takes one or two arguments: the first argument identifies
// the container by reference or by path (software system name followed by colon
// and container name). The second optional argument is a func() that defines
// additional properties on the container instance including the container
// instance ID.
//
// Usage:
//
//	ContainerInstance(Container)
//
//	ContainerInstance(Container, func())
//
//	ContainerInstance("<Software System>/<Container>")
//
//	ContainerInstance("<Software System>/<Container>", func())
//
// Example:
//
//	var _ = Design(func() {
//	    var MyContainer *expr.Container
//	    SoftwareSystem("SoftwareSystem", "A software system", func() {
//	        MyContainer = Container("Container")
//	    })
//	    DeploymentEnvironment("Production", func() {
//	        DeploymentNode("US", "US shard", func() {
//	            ContainerInstance(MyContainer, func() {
//	                Tag("service")
//	                InstanceID(1)
//	                HealthCheck("check", func() {
//	                    URL("https://goa.design/health")
//	                    Interval(10)
//	                    Timeout(1000)
//	                })
//	            })
//	            // Using the name instead:
//	            ContainerInstance("SoftwareSystem/Container", func() {
//	                InstanceID(2
//	            })
//	        })
//	    })
//	})
func ContainerInstance(container interface{}, dsl ...func()) *expr.ContainerInstance {
	d, ok := eval.Current().(*expr.DeploymentNode)
	if !ok {
		eval.IncompatibleDSL()
	}
	var cont *expr.Container
	switch c := container.(type) {
	case *expr.Container:
		cont = c
	case string:
		eh, err := expr.Root.Model.FindElement(nil, c)
		if err != nil {
			eval.ReportError("ContainerInstance: " + err.Error())
			return nil
		}
		var ok bool
		cont, ok = eh.(*expr.Container)
		if !ok {
			eval.ReportError("ContainerInstance: expected path to container, got %T", eh)
			return nil
		}
	default:
		eval.ReportError("ContainerInstance: expected container or path to container, got %T", container)
		return nil
	}
	var f func()
	if len(dsl) > 0 {
		f = dsl[0]
		if len(dsl) > 1 {
			eval.ReportError("ContainerInstance: too many arguments")
		}
	}
	ci := &expr.ContainerInstance{
		Element: &expr.Element{
			Name:        cont.Name,
			Description: cont.Description,
			URL:         cont.URL,
			Technology:  cont.Technology,
			DSLFunc:     f,
		},
		Parent:      d,
		Environment: d.Environment,
		ContainerID: cont.ID,
		InstanceID:  1,
	}
	return d.AddContainerInstance(ci)
}

// Instances sets the number of instances of the deployment node.
//
// Instances must appear in a DeploymentNode expression.
//
// Instances accepts a single argument which is the number.
//
// Example:
//
//	var _ = Design(func() {
//	    DeploymentEnvironment("Production", func() {
//	        DeploymentNode("Web app", func() {
//	            Instances(3)
//	        })
//	    })
//	})
func Instances(n int) {
	node, ok := eval.Current().(*expr.DeploymentNode)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	node.Instances = &n
}

// InstanceID sets the instance number or index of a container instance.
//
// InstanceID must appear in a ContainerInstance expression.
//
// InstanceID accepts a single argument which is the number.
//
// Example:
//
//	var _ = Design(func() {
//	    DeploymentEnvironment("Production", func() {
//	        ContainerInstance(Container, func() {
//	            InstanceID(3)
//	        })
//	    })
//	})
func InstanceID(n int) {
	node, ok := eval.Current().(*expr.ContainerInstance)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	node.InstanceID = n
}

// HealthCheck defines a HTTP-based health check for a container instance.
//
// HealthCheck must appear in a ContainerInstance expression.
//
// HealthCheck accepts two arguments: the health check name and a function used
// to define additional required properties.
//
// Example:
//
//	var _ = Design(func() {
//	    DeploymentEnvironment("Production", func() {
//	        ContainerInstance(Container, func() {
//	            HealthCheck("check", func() {
//	                URL("https://goa.design/health")
//	                Interval(10)
//	                Timeout(1000)
//	                Header("X-Foo", "bar")
//	            })
//	        })
//	    })
//	})
func HealthCheck(name string, dsl func()) {
	c, ok := eval.Current().(*expr.ContainerInstance)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	hc := &expr.HealthCheck{Name: name}
	eval.Execute(dsl, hc)
	c.HealthChecks = append(c.HealthChecks, hc)
}

// Interval defines a health check polling interval in seconds.
//
// Interval must appear in a HealthCheck expression.
//
// Interval takes one argument: the number of seconds.
//
// Example:
//
//	var _ = Design(func() {
//	    DeploymentEnvironment("Production", func() {
//	        ContainerInstance(Container, func() {
//	            HealthCheck("check", func() {
//	                Interval(10)
//	            })
//	        })
//	    })
//	})
func Interval(n int) {
	hc, ok := eval.Current().(*expr.HealthCheck)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	hc.Interval = n
}

// Timeout defines a health check timeout in milliseconds.
//
// Timeout must appear in a HealthCheck expression.
//
// Timeout takes one argument: the number of milliseconds.
//
// Example:
//
//	var _ = Design(func() {
//	    DeploymentEnvironment("Production", func() {
//	        ContainerInstance(Container, func() {
//	            HealthCheck("check", func() {
//	                Timeout(1000)
//	            })
//	        })
//	    })
//	})
func Timeout(n int) {
	hc, ok := eval.Current().(*expr.HealthCheck)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	hc.Timeout = n
}

// Header defines a header name and value to be set in requests sent for health
// checks.
//
// Header must appear in a HealthCheck expression.
//
// Header takes two arguments: the header name and value.
//
// Example:
//
//	var _ = Design(func() {
//	    DeploymentEnvironment("Production", func() {
//	        ContainerInstance(Container, func() {
//	            HealthCheck("check", func() {
//	                Header("X-Foo", "bar")
//	            })
//	        })
//	    })
//	})
func Header(n, v string) {
	hc, ok := eval.Current().(*expr.HealthCheck)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	if hc.Headers == nil {
		hc.Headers = make(map[string]string)
	}
	hc.Headers[n] = v
}
