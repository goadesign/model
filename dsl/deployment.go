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
//        DeploymentEnvironment("production", func() {
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
		ID:          expr.NewID(),
		Name:        name,
		Description: description,
		Technology:  technology,
		Environment: env.Name,
	}
	if dsl != nil {
		eval.Execute(dsl, node)
	}
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
//        DeploymentEnvironment("production", func() {
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
		ID:          expr.NewID(),
		Name:        name,
		Description: description,
		Technology:  technology,
		Environment: env.Name,
	}
	if dsl != nil {
		eval.Execute(dsl, node)
	}
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
//        DeploymentEnvironment("production", func() {
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
		ID:          expr.NewID(),
		ContainerID: d.ID,
		Environment: env.Name,
	}
	if dsl != nil {
		eval.Execute(dsl, ci)
	}
	d.ContainerInstances = append(d.ContainerInstances, ci)
}
