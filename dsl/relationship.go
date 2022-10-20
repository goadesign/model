package dsl

import (
	"fmt"

	"goa.design/goa/v3/eval"
	"goa.design/model/expr"
)

// InteractionStyleKind is the enum for possible interaction styles.
type InteractionStyleKind int

const (
	// Synchronous describes a synchronous interaction.
	Synchronous InteractionStyleKind = iota + 1
	// Asynchronous describes an asynchronous interaction.
	Asynchronous
)

// Uses adds a uni-directional relationship between two elements.
//
// Uses may appear in Person, SoftwareSystem, Container or Component.
//
// Uses takes 2 to 5 arguments. The first argument identifies the target of the
// relationship. The following argument is a short description for the
// relationship. The description may optionally be followed by the technology
// used by the relationship and/or the type of relationship: Synchronous or
// Asynchronous. Finally Uses accepts an optional func() as last argument to
// define additional properties on the relationship.
//
// The target of the relationship is identified by providing an element (person,
// software system, container or component) or the path of an element. The path
// consists of the element name if a top level element (person or software
// system) or if the element is in scope (container in the same software system
// as the source or component in the same container as the source). When the
// element is not in scope the path specifies the parent element name followed
// by a slash and the element name. If the parent itself is not in scope (i.e. a
// component that is a child of a different software system than the source)
// then the path specifies the top-level software system followed by a slash,
// the container name, another slash and the component name.
//
// Usage:
//
//	Uses(Element, "<description>")
//
//	Uses(Element, "<description>", "[technology]")
//
//	Uses(Element, "<description>", Synchronous|Asynchronous)
//
//	Uses(Element, "<description>", "[technology]", Synchronous|Asynchronous)
//
//	Uses(Element, "<description>", func())
//
//	Uses(Element, "<description>", "[technology]", func())
//
//	Uses(Element, "<description>", Synchronous|Asynchronous, func())
//
//	Uses(Element, "<description>", "[technology]", Synchronous|Asynchronous, func())
//
// Where Element is one of:
//
//   - Person, SoftwareSystem, Container or Component
//   - "<Person>", "<SoftwareSystem>", "<SoftwareSystem>/<Container>" or "<SoftwareSystem>/<Container>/<Component>"
//   - "<Container>" (if container is a sibling of the source)
//   - "<Component>" (if component is a sibling of the source)
//   - "<Container>/<Component>" (if container is a sibling of the source)
//
// Example:
//
//	var _ = Design("my workspace", "a great architecture model", func() {
//	    SoftwareSystem("SystemA", func() {
//	        Container("ContainerA")
//	        Container("ContainerB", func() {
//	            Uses("ContainerA") // sibling, not need to specify system
//	        })
//	    })
//	    SoftwareSystem("SystemB", func() {
//	        Uses("SystemA/ContainerA") // not a sibling, need full path
//	    })
//	    Person("Customer", "Customers of enterprise", func () {
//	       Uses(SystemA, "Access", "HTTP", Synchronous)
//	    })
//	    Person("Staff", "Back office staff", func() {
//	       InteractsWith("Customer", "Sends invoices to", Synchronous)
//	    })
//	})
func Uses(element interface{}, description string, args ...interface{}) {
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
	uses(src, element, description, args...)
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
// Usage:
//
//	InteractsWith(Person|"Person", "<description>")
//
//	InteractsWith(Person|"Person", "<description>", "[technology]")
//
//	InteractsWith(Person|"Person", "<description>", Synchronous|Asynchronous)
//
//	InteractsWith(Person|"Person", "<description>", "[technology]", Synchronous|Asynchronous)
//
//	InteractsWith(Person|"Person", "<description>", func())
//
//	InteractsWith(Person|"Person", "<description>", "[technology]", func())
//
//	InteractsWith(Person|"Person", "<description>", Synchronous|Asynchronous, func())
//
//	InteractsWith(Person|"Person", "<description>", "[technology]", Synchronous|Asynchronous, func())
//
// Example:
//
//	var _ = Design("my workspace", "a great architecture model", func() {
//	    var Employee = Person("Employee")
//	    Person("Customer", "Customers of enterprise", func () {
//	       InteractsWith(Employee, "Sends requests to", "email")
//	    })
//	})
func InteractsWith(person interface{}, description string, args ...interface{}) {
	src, ok := eval.Current().(*expr.Person)
	if !ok {
		eval.IncompatibleDSL()
	}
	switch p := person.(type) {
	case *expr.Person:
		if err := uses(src.Element, p, description, args...); err != nil {
			eval.ReportError("InteractsWith: %s", err.Error())
		}
	case string:
		e := expr.Root.Model.Person(p)
		if e == nil {
			eval.ReportError("InteractsWith: unknown person %q", p)
			return
		}
		if err := uses(src.Element, e, description, args...); err != nil {
			eval.ReportError("InteractsWith: %s", err.Error())
		}
	default:
		eval.InvalidArgError("person or name of person", person)
		return
	}
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
// Usage:
//
//	Delivers(Person|"Person", "<description>")
//
//	Delivers(Person|"Person", "<description>", "[technology]")
//
//	Delivers(Person|"Person", "<description>", Synchronous|Asynchronous)
//
//	Delivers(Person|"Person", "<description>", "[technology]", Synchronous|Asynchronous)
//
//	Delivers(Person|"Person", "<description>", func())
//
//	Delivers(Person|"Person", "<description>", "[technology]", func())
//
//	Delivers(Person|"Person", "<description>", Synchronous|Asynchronous, func())
//
//	Delivers(Person|"Person", "<description>", "[technology]", Synchronous|Asynchronous, func())
//
// Example:
//
//	var _ = Design("my workspace", "a great architecture model", func() {
//	    var Customer = Person("Customer")
//	    SoftwareSystem("MySystem", func () {
//	       Delivers(Customer, "Sends requests to", "email")
//	    })
//	})
func Delivers(person interface{}, description string, args ...interface{}) {
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

	switch p := person.(type) {
	case *expr.Person:
		if err := uses(src, p, description, args...); err != nil {
			eval.ReportError("Delivers: %s", err.Error())
		}
	case string:
		e := expr.Root.Model.Person(p)
		if e == nil {
			eval.ReportError("Delivers: unknown person %q", p)
			return
		}
		if err := uses(src, e, description, args...); err != nil {
			eval.ReportError("Delivers: %s", err.Error())
		}
	default:
		eval.InvalidArgError("person or name of person", person)
		return
	}

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
func uses(src *expr.Element, dest interface{}, desc string, args ...interface{}) error {
	var (
		technology string
		style      InteractionStyleKind
		dsl        func()
	)
	if len(args) > 0 {
		switch a := args[0].(type) {
		case string:
			technology = a
		case InteractionStyleKind:
			style = a
		case func():
			dsl = a
		default:
			return fmt.Errorf("expected description, Synchronous or Asynchronous, got %T", args[0])
		}
		if len(args) > 1 {
			if dsl != nil {
				return fmt.Errorf("function DSL must be last argument")
			}
			switch a := args[1].(type) {
			case InteractionStyleKind:
				style = a
			case func():
				dsl = a
			default:
				return fmt.Errorf("expected Synchronous or Asynchronous, got %T", args[1])
			}
			if len(args) > 2 {
				if d, ok := args[2].(func()); ok {
					dsl = d
				} else {
					return fmt.Errorf("expected DSL function, got %T", args[2])
				}
				if len(args) > 3 {
					return fmt.Errorf("too many arguments")
				}
			}
		}
	}
	rel := &expr.Relationship{
		Description:      desc,
		Source:           src,
		Technology:       technology,
		InteractionStyle: expr.InteractionStyleKind(style),
	}
	// Note: we need to check the types explicitly below because
	// (*expr.Person)(nil) != (expr.ElementHolder)(nil) for example.
	switch d := dest.(type) {
	case *expr.Person:
		if d == nil {
			return fmt.Errorf("Person reference is nil")
		}
		rel.Destination = d.Element
	case *expr.SoftwareSystem:
		if d == nil {
			return fmt.Errorf("SoftwareSystem reference is nil")
		}
		rel.Destination = d.Element
	case *expr.Container:
		if d == nil {
			return fmt.Errorf("Container reference is nil")
		}
		rel.Destination = d.Element
	case *expr.Component:
		if d == nil {
			return fmt.Errorf("Component reference is nil")
		}
		rel.Destination = d.Element
	case string:
		rel.DestinationPath = d
	default:
		return fmt.Errorf("invalid argument type for destination expected element or string (element path) got %T", d)
	}
	if dsl != nil {
		eval.Execute(dsl, rel)
	}
	expr.Identify(rel)
	src.Relationships = append(src.Relationships, rel)

	return nil
}
