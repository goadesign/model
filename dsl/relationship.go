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
// Uses may appear in Person, SoftwareSystem, Container or Component.
//
// Uses tags 2 to 5 arguments. The first argument is the target of the
// relationship, it must be a software system, container, component or the name
// of one of these. The second argument is a short description for the
// relationship. The description may optionally be followed by the technology
// used by the relationship and/or the type of relationship: Synchronous or
// Asynchronous. Finally Uses accepts an optional func() as last argument to
// define additional properties on the relationship.
//
// Usage is thus:
//
//    Uses(Element, "<description>")
//
//    Uses(Element, "<description>", "[technology]")
//
//    Uses(Element, "<description>", Synchronous|Asynchronous)
//
//    Uses(Element, "<description>", "[technology]", Synchronous|Asynchronous)
//
//    Uses(Element, "<description>", func())
//
//    Uses(Element, "<description>", "[technology]", func())
//
//    Uses(Element, "<description>", Synchronous|Asynchronous, func())
//
//    Uses(Element, "<description>", "[technology]", Synchronous|Asynchronous, func())
//
// Example:
//
//     var _ = Workspace("my workspace", "a great architecture model", func() {
//         var MySystem = SoftwareSystem("My system")
//         Person("Customer", "Customers of enterprise", func () {
//            Uses(MySystem, "Access", "HTTP", Synchronous)
//         })
//         Person("Staff", "Back office staff", func() {
//            InteractsWith("Customer", "Sends invoices to", Synchronous)
//         })
//     })
//
func Uses(element interface{}, desc string, args ...interface{}) {
	var src *expr.Element
	switch e := eval.Current().(type) {
	case *expr.Person:
		src = e.Element
	case *expr.SoftwareSystem:
		src = e.Element
	case *expr.Container:
		src = e.Element
	case *expr.Component:
		src = e.Element
	default:
		eval.IncompatibleDSL()
		return
	}
	var n string
	switch e := element.(type) {
	case *expr.SoftwareSystem:
		n = e.Name
	case *expr.Container:
		n = e.Name
	case *expr.Component:
		n = e.Name
	case string:
		n = e
	default:
		eval.IncompatibleDSL()
		return
	}
	uses(src, n, desc, args...)
}

// InteractsWith adds an interaction between a person and another.
//
// InteractsWith must appear in Person.
//
// InteractsWith accepts 2 to 5 arguments. The first argument is the target of
// the relationship, it must be a person or the name of a person. The target may
// optionally be followed by a short description of the relationship. The
// description may optionally be followed by the technology used by the
// relationship and/or the type of relationship: Synchronous or Asynchronous.
// Finally InteractsWith accepts an optional func() as last argument to add
// further properties to the relationship.
//
// Usage is thus:
//
//    InteractsWith(Person, "<description>")
//
//    InteractsWith(Person, "<description>", "[technology]")
//
//    InteractsWith(Person, "<description>", Synchronous|Asynchronous)
//
//    InteractsWith(Person, "<description>", "[technology]", Synchronous|Asynchronous)
//
//    InteractsWith(Person, "<description>", func())
//
//    InteractsWith(Person, "<description>", "[technology]", func())
//
//    InteractsWith(Person, "<description>", Synchronous|Asynchronous, func())
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
func InteractsWith(p interface{}, desc string, args ...interface{}) {
	var n string
	switch a := p.(type) {
	case *expr.Person:
		n = a.Name
	case string:
		n = a
	default:
		eval.InvalidArgError("person or name of person", p)
		return
	}
	if c, ok := eval.Current().(*expr.Person); ok {
		uses(c.Element, n, desc, args...)
		return
	}
	eval.IncompatibleDSL()
}

// Delivers adds an interaction between an element and a person.
//
// Delivers must appear in SoftareSystem, Container or Component.
//
// Delivers accepts 2 to 5 arguments. The first argument is the target of the
// relationship, it must be a person or the name of a person. The target may
// optionally be followed by a short description of the relationship. The
// description may optionally be followed by the technology used by the
// relationship and/or the type of relationship: Synchronous or Asynchronous.
// Finally Delivers accepts an optional func() as last argument to add further
// properties to the relationship.
//
// Usage is thus:
//
//    Delivers(Person, "<description>")
//
//    Delivers(Person, "<description>", "[technology]")
//
//    Delivers(Person, "<description>", Synchronous|Asynchronous)
//
//    Delivers(Person, "<description>", "[technology]", Synchronous|Asynchronous)
//
//    Delivers(Person, "<description>", func())
//
//    Delivers(Person, "<description>", "[technology]", func())
//
//    Delivers(Person, "<description>", Synchronous|Asynchronous, func())
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
func Delivers(p interface{}, desc string, args ...interface{}) {
	var src *expr.Element
	switch e := eval.Current().(type) {
	case *expr.SoftwareSystem:
		src = e.Element
	case *expr.Container:
		src = e.Element
	case *expr.Component:
		src = e.Element
	default:
		eval.IncompatibleDSL()
		return
	}

	var n string
	switch a := p.(type) {
	case *expr.Person:
		n = a.Name
	case string:
		n = a
	default:
		eval.InvalidArgError("person or name of person", p)
		return
	}

	uses(src, n, desc, args...)
}

// Description provides a short description for a relationship displayed in a
// dynamic view.
//
// Description must appear in Add.
//
// Description takes one argument: the relationship description used in the
// dynamic view.
func Description(desc string) {
	v, ok := eval.Current().(*expr.RelationshipView)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	v.Description = desc
}

// uses adds a relationship between the given source and destination. The caller
// must make sure that the relationship is valid.
func uses(src *expr.Element, dest, desc string, args ...interface{}) *expr.Relationship {
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
			eval.InvalidArgError("description, Synchronous or Asynchronous", args[0])
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
				eval.InvalidArgError("Synchronous or Asynchronous", args[1])
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
		SourceID:         src.ID,
		Technology:       technology,
		InteractionStyle: style,
		Source:           src,
		DestinationName:  dest,
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
