package dsl

import (
	"goa.design/goa/v3/eval"
	"goa.design/structurizr/expr"
)

// DeploymentEnvironment defines a deployment environment (e.g. development,
// production).
//
// DeploymentEnvironment must appear in a Workspace expression.
//
// DeploymentEnvironment accepts two arguments: the environment name and a DSL
// function used to describe the nodes within the environment.
//
// Example:
//
//    var _ = Workspace(func() {
//         DeploymentEnvironment("production", func() {
//             DeploymentNode("AppServer", "Application server", "Go and Goa v3")
//             InfrastructureNote("Router", "External traffic router", "AWS Route 53")
//             ContainerInstance(Container)
//         })
//     })
//
func DeploymentEnvironment(name string, dsl func()) {
	_, ok := eval.Current().(*expr.Workspace)
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
// DeploymentNode must appear in a DeploymentEnvironment expression.
//
// DeploymentNode takes 1 to 4 arguments. The first argument is the node name.
// The name may be optionally followed by a description. If a description is set
// then it may be followed by the technology details used by the component.
// Finally DeploymentNode may take a func() as last argument to define
// additional properties of the component.
//
// The valid syntax for DeploymentNode is thus:
//
//    DeploymentNode("<name>")
//
//    DeploymentNode("<name>", "[description]")
//
//    DeploymentNode("<name>", "[description]", "[technology]")
//
//    DeploymentNode("<name>", func())
//
//    DeploymentNode("<name>", "[description]", func())
//
//    DeploymentNode("<name>", "[description]", "[technology]", func())
//
// Example:
//
//    var _ = Workspace(func() {
//        DeploymentEnvironment("Production", func() {
//            DeploymentNode("US", "US shard", func() {
//                Tag("shard")
//                Instances(3)
//                URL("https://goa.design/docs/shard")
//                Uses(OtherDeploymentNode, "Uses", "gRPC", Asynchronous)
//            })
//        })
//    })
//
func DeploymentNode(name string, args ...interface{}) {
	env, ok := eval.Current().(*expr.DeploymentEnvironment)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	description, technology, dsl := parseElementArgs(args...)
	node := &expr.DeploymentNode{
		Element: expr.Element{
			Name:        name,
			Description: description,
			Technology:  technology,
			DSLFunc:     dsl,
		},
		Environment: env.Name,
	}
	expr.Identify(node)
	expr.Root.Model.DeploymentNodes = append(expr.Root.Model.DeploymentNodes, node)
}

// InfrastructureNode defines an infrastructure node, typically something like a
// load balancer, firewall, DNS service, etc.
//
// InfrastructureNode must appear in a DeploymentEnvironment expression.
//
// InfrastructureNode takes 2 to 5 arguments. The first argument is the parent
// deployment node. The second argument is the infrastructure node name. The
// name may be optionally followed by a description. If a description is set
// then it may be followed by the technology details used by the component.
// Finally InfrastructureNode may take a func() as last argument to define
// additional properties of the component.
//
// The valid syntax for InfrastructureNode is thus:
//
//    InfrastructureNode(DeploymentNode, "<name>")
//
//    InfrastructureNode(DeploymentNode, "<name>", "[description]")
//
//    InfrastructureNode(DeploymentNode, "<name>", "[description]", "[technology]")
//
//    InfrastructureNode(DeploymentNode, "<name>", func())
//
//    InfrastructureNode(DeploymentNode, "<name>", "[description]", func())
//
//    InfrastructureNode(DeploymentNode, "<name>", "[description]", "[technology]", func())
//
// Example:
//
//    var _ = Workspace(func() {
//        DeploymentEnvironment("Production", func() {
//            InfrastructureNode(DeploymentNode, "US", "US shard", func() {
//                Tag("shard")
//                Instances(3)
//                URL("https://goa.design/docs/shards/us")
//                Uses(OtherInfrastructureNode, "Uses", "gRPC", Asynchronous)
//            })
//        })
//    })
//
func InfrastructureNode(d *expr.DeploymentNode, name string, args ...interface{}) {
	env, ok := eval.Current().(*expr.DeploymentEnvironment)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	description, technology, dsl := parseElementArgs(args...)
	node := &expr.InfrastructureNode{
		Element: expr.Element{
			Name:        name,
			Description: description,
			Technology:  technology,
			DSLFunc:     dsl,
		},
		Environment: env.Name,
	}
	expr.Identify(node)
	d.InfrastructureNodes = append(d.InfrastructureNodes, node)
}

// ContainerInstance defines an instance of the specified container that is
// deployed on the parent deployment node.
//
// ContainerInstance must appear in a DeploymentEnvironment expression.
//
// ContainerInstance takes 1 or 2 arguments. The first argument is the parent
// deployment node. The second argument is an optional func() that defines
// additional properties on the container instance.
//
// Example:
//
//    var _ = Workspace(func() {
//        DeploymentEnvironment("Production", func() {
//            ContainerInstance(DeploymentNode, func() {
//                Tag("shard")
//                InstanceID(1)
//                Uses(OtherContainerInstance, "Uses", "gRPC", Asynchronous)
//                HealthCheck("check", func() {
//                    URL("https://goa.design/health")
//                    Interval(10)
//                    Timeout(1000)
//                })
//            })
//        })
//    })
//
func ContainerInstance(d *expr.DeploymentNode, args ...func()) {
	env, ok := eval.Current().(*expr.DeploymentEnvironment)
	if !ok {
		eval.IncompatibleDSL()
	}
	var dsl func()
	if len(args) > 0 {
		dsl = args[0]
		if len(args) > 1 {
			eval.ReportError("too many arguments")
		}
	}
	ci := &expr.ContainerInstance{
		Element:     expr.Element{DSLFunc: dsl},
		ContainerID: d.ID,
		Environment: env.Name,
	}
	expr.Identify(ci)
	d.ContainerInstances = append(d.ContainerInstances, ci)
}

// Instances sets the number of instances of the deployment node.
//
// Instances must appear in a DeploymentNode expression.
//
// Instances accepts a single parameter which is the number.
//
// Example:
//
//    var _ = Workspace(func() {
//        DeploymentEnvironment("Production", func() {
//            DeploymentNode("Web app", func() {
//                Instances(3)
//            })
//        })
//    })
//
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
// InstanceID accepts a single parameter which is the number.
//
// Example:
//
//    var _ = Workspace(func() {
//        DeploymentEnvironment("Production", func() {
//            ContainerInstance(Container, func() {
//                InstanceID(3)
//            })
//        })
//    })
//
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
//    var _ = Workspace(func() {
//        DeploymentEnvironment("Production", func() {
//            ContainerInstance(Container, func() {
//                HealthCheck("check", func() {
//                    URL("https://goa.design/health")
//                    Interval(10)
//                    Timeout(1000)
//                    Header("X-Foo", "bar")
//                })
//            })
//        })
//    })
//
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
//    var _ = Workspace(func() {
//        DeploymentEnvironment("Production", func() {
//            ContainerInstance(Container, func() {
//                HealthCheck("check", func() {
//                    Interval(10)
//                })
//            })
//        })
//    })
//
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
//    var _ = Workspace(func() {
//        DeploymentEnvironment("Production", func() {
//            ContainerInstance(Container, func() {
//                HealthCheck("check", func() {
//                    Timeout(1000)
//                })
//            })
//        })
//    })
//
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
//    var _ = Workspace(func() {
//        DeploymentEnvironment("Production", func() {
//            ContainerInstance(Container, func() {
//                HealthCheck("check", func() {
//                    Header("X-Foo", "bar")
//                })
//            })
//        })
//    })
//
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
