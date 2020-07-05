package dsl

import (
	"fmt"

	"goa.design/goa/v3/eval"
	"goa.design/structurizr/expr"
)

const (
	// Synchronous describes a synchronous interaction.
	Synchronous = expr.InteractionSynchronous
	// Asynchronous describes an asynchronous interaction.
	Asynchronous = expr.InteractionAsynchronous
)

// Uses adds a uni-directional relationship between two elements.
//
// Uses may appear in Person, SoftwareSystem, Container, Component,
// DeploymentNode, InfrastructureNode or ContainerInstance.
//
// Uses tags 2 to 5 arguments. The first argument is the target of the
// relationship, it must be a software system, container or component if the
// scope is Person, SoftwareSystem, Container or Component. It must be a
// DeploymentNode, InfrastructureNode or ContainerInstance if the scope is one
// of these. The second argument is a short description for the relationship.
// The description may optionally be followed by the technology used by the
// relationship. If technology is set then Uses accepts an additional argument
// to indicate the type of relationship: Synchronous or Asynchronous. Finally
// Uses accepts an optional func() as last argument to define additional
// properties on the relationship.
//
// Usage is thus:
//
//    Uses(Element, "<description>")
//
//    Uses(Element, "<description>", "[technology]")
//
//    Uses(Element, "<description>", "[technology]", Synchronous|Asynchronous)
//
//    Uses(Element, "<description>", func())
//
//    Uses(Element, "<description>", "[technology]", func())
//
//    Uses(Element, "<description>", "[technology]", Synchronous|Asynchronous, func())
//
// Example:
//
//     var _ = Workspace("my workspace", "a great architecture model", func() {
//         var MySystem = SoftwareSystem("My system")
//         Person("Customer", "Customers of enterprise", func () {
//            Uses(MySystem, "Access", "HTTP", InteractionSynchronous)
//         })
//     })
//
func Uses(element interface{}, desc string, args ...interface{}) {
	// 1. Relationships between elements (software systems, containers and
	// components)
	var srcID string
	switch e := eval.Current().(type) {
	case *expr.SoftwareSystem:
		srcID = e.ID
	case *expr.Container:
		srcID = e.ID
	case *expr.Component:
		srcID = e.ID
	default:
		eval.IncompatibleDSL()
		return
	}

	var destID string
	switch e := eval.Current().(type) {
	case *expr.SoftwareSystem:
		destID = e.ID
	case *expr.Container:
		destID = e.ID
	case *expr.Component:
		destID = e.ID
	default:
		eval.IncompatibleDSL()
		return
	}

	if srcID != "" && destID == "" || srcID == "" && destID != "" {
		eval.ReportError("Uses used in an element (SoftareSystem, Container or Component) must target another element.")
	}

	if srcID != "" {
		uses(srcID, destID, desc, args...)
		return
	}

	// 2. Relationships between deployment nodes.
	if d, ok := eval.Current().(*expr.DeploymentNode); ok {
		if dd, ok := element.(*expr.DeploymentNode); ok {
			uses(d.ID, dd.ID, desc, args...)
		} else {
			eval.InvalidArgError("deployment node", fmt.Sprintf("%T", element))
		}
		return
	}

	// 3. Relationships between infrastructure node and another deployment
	// element.
	if i, ok := eval.Current().(*expr.InfrastructureNode); ok {
		srcID := i.ID
		var destID string
		switch e := element.(type) {
		case *expr.DeploymentNode:
			destID = e.ID
		case *expr.InfrastructureNode:
			destID = e.ID
		case *expr.ContainerInstance:
			destID = e.ID
		default:
			eval.InvalidArgError("deployment node, infrastructure node or container instance", fmt.Sprintf("%T", element))
			return
		}
		uses(srcID, destID, desc, args...)
	}

	// 4. Relationships between container instances.
	if c, ok := eval.Current().(*expr.ContainerInstance); ok {
		if cc, ok := element.(*expr.ContainerInstance); ok {
			uses(c.ID, cc.ID, desc, args...)
		} else {
			eval.InvalidArgError("container instance", fmt.Sprintf("%T", element))
		}
		return
	}

}

// InteractsWith adds an interaction between a person and another.
//
// InteractsWith must appear in Person.
//
// InteractsWith accepts 2 to 5 arguments. The first argument is the target of
// the relationship, it must be a person. The target may optionally be followed
// by a short description of the relationship. The description may optionally be
// followed by the technology used by the relationship. If technology is set
// then InteractsWith accepts an additional argument to indicate the type of
// relationship: Synchronous or Asynchronous. Finally InteractsWith accepts an
// optional func() as last argument to add further properties to the relationship.
//
// Usage is thus:
//
//    InteractsWith(Person, "<description>")
//
//    InteractsWith(Person, "<description>", "[technology]")
//
//    InteractsWith(Person, "<description>", "[technology]", Synchronous|Asynchronous)
//
//    InteractsWith(Person, "<description>", func())
//
//    InteractsWith(Person, "<description>", "[technology]", func())
//
//    InteractsWith(Person, "<description>", "[technology]", Synchronous|Asynchronous, func())
//
// Example:
//
//     var _ = Workspace("my workspace", "a great architecture model", func() {
//         var Employee = Person("Employee")
//         Person("Customer", "Customers of enterprise", func () {
//            InteractsWith(Employee, "Sends requests to", "email")
//         })
//     })
//
func InteractsWith(p *expr.Person, desc string, args ...interface{}) {
	if c, ok := eval.Current().(*expr.Person); ok {
		uses(c.ID, p.ID, desc, args...)
	}
}

// Delivers adds an interaction between an element and a person.
//
// Delivers must appear in SoftareSystem, Container or Component.
//
// Delivers accepts 2 to 5 arguments. The first argument is the target of
// the relationship, it must be a person. The target may optionally be followed
// by a short description of the relationship. The description may optionally be
// followed by the technology used by the relationship. If technology is set
// then Delivers accepts an additional argument to indicate the type of
// relationship: Synchronous or Asynchronous. Finally Delivers accepts an
// optional func() as last argument to add further properties to the relationship.
//
// Usage is thus:
//
//    Delivers(Person, "<description>")
//
//    Delivers(Person, "<description>", "[technology]")
//
//    Delivers(Person, "<description>", "[technology]", Synchronous|Asynchronous)
//
//    Delivers(Person, "<description>", func())
//
//    Delivers(Person, "<description>", "[technology]", func())
//
//    Delivers(Person, "<description>", "[technology]", Synchronous|Asynchronous, func())
//
// Example:
//
//     var _ = Workspace("my workspace", "a great architecture model", func() {
//         var Customer = Person("Customer")
//         SoftwareSystem("MySystem", func () {
//            Delivers(Customer, "Sends requests to", "email")
//         })
//     })
//
func Delivers(p *expr.Person, desc string, args ...interface{}) {
	var srcID string
	switch e := eval.Current().(type) {
	case *expr.SoftwareSystem:
		srcID = e.ID
	case *expr.Container:
		srcID = e.ID
	case *expr.Component:
		srcID = e.ID
	default:
		eval.IncompatibleDSL()
		return
	}
	uses(srcID, p.ID, desc, args...)
}

// uses adds a relationship between the given source and destination. The caller
// must make sure that the relationship is valid.
func uses(srcID, destID string, desc string, args ...interface{}) *expr.Relationship {
	var (
		technology string
		style      expr.InteractionStyleKind
		dsl        func()
	)
	if len(args) > 0 {
		switch a := args[0].(type) {
		case string:
			technology = a
		case expr.InteractionStyleKind:
			style = a
		case func():
			dsl = a
		default:
			eval.InvalidArgError("description or InteractionSynchronous or InteractionAsynchronous", args[0])
		}
		if len(args) > 1 {
			if dsl != nil {
				eval.ReportError("function DSL must be last argument")
			}
			switch a := args[1].(type) {
			case expr.InteractionStyleKind:
				style = a
			case func():
				dsl = a
			default:
				eval.InvalidArgError("InteractionSynchronous or InteractionAsynchronous", args[1])
			}
			if len(args) > 2 {
				if d, ok := args[2].(func()); ok {
					dsl = d
				} else {
					eval.InvalidArgError("DSL function", args[2])
				}
				if len(args) > 3 {
					eval.ReportError("too many arguments")
				}
			}
		}
	}
	rel := &expr.Relationship{
		Description:      desc,
		SourceID:         srcID,
		DestinationID:    destID,
		Technology:       technology,
		InteractionStyle: style,
	}
	if dsl != nil {
		eval.Execute(dsl, rel)
	}
	expr.Identify(rel)

	switch e := eval.Current().(type) {
	case *expr.Person:
		e.Rels = append(e.Rels, rel)
	case *expr.SoftwareSystem:
		e.Rels = append(e.Rels, rel)
	case *expr.Container:
		e.Rels = append(e.Rels, rel)
	case *expr.Component:
		e.Rels = append(e.Rels, rel)
	case *expr.DeploymentNode:
		e.Rels = append(e.Rels, rel)
	case *expr.InfrastructureNode:
		e.Rels = append(e.Rels, rel)
	case *expr.ContainerInstance:
		e.Rels = append(e.Rels, rel)
	default:
		panic(fmt.Sprintf("unexpected expression type %T", eval.Current())) // bug
	}
	return rel
}
