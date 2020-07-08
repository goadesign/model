package dsl

import (
	"goa.design/goa/v3/eval"
	goaexpr "goa.design/goa/v3/expr"
	"goa.design/structurizr/expr"
)

// SoftwareSystem defines a software system.
//
// SoftwareSystem must appear in a Workspace expression.
//
// Software system takes 1 to 3 arguments. The first argument is the software
// system name and the last argument a function that contains the expressions
// that defines the content of the system. An optional description may be given
// after the name.
//
// The valid syntax for SoftwareSystem is thus:
//
//    SoftwareSystem("<name>")
//
//    SoftwareSystem("<name>", "[description]")
//
//    SoftwareSystem("<name>", func())
//
//    SoftwareSystem("<name>", "[description]", func())
//
// Example:
//
//    var _ = Workspace(func() {
//        SoftwareSystem("My system", "A system with a great architecture", func() {
//            Tag("bill processing")
//            URL("https://goa.design/mysystem")
//            External()
//            Uses(OtherSystem, "Uses", "gRPC", Synchronous)
//            Delivers(Customer, "Delivers emails to", "SMTP", Synchronous)
//        })
//    })
//
func SoftwareSystem(name string, args ...interface{}) {
	w, ok := eval.Current().(*expr.Workspace)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	description, _, dsl := parseElementArgs(args...)
	s := &expr.SoftwareSystem{
		Element: &expr.Element{
			Name:        name,
			Description: description,
		},
		Location: expr.LocationInternal,
	}
	if dsl != nil {
		eval.Execute(dsl, s)
	}
	expr.Identify(s)
	w.Model.Systems = append(w.Model.Systems, s)
}

// Container defines a container.
//
// Container must appear in a Workspace expression.
//
// Container takes 2 to 5 arguments. The first argument is the software system
// that the container belongs to. The second argument is the container name. The
// name may be optionally followed by a description. If a description is set
// then it may be followed by the technology details used by the container.
// Finally Container may take a func() as last argument to define additional
// properties of the container.
//
// The valid syntax for Container is thus:
//
//    Container(SoftwareSystem, "<name>")
//
//    Container(SoftwareSystem, "<name>", "[description]")
//
//    Container(SoftwareSystem, "<name>", "[description]", "[technology]")
//
//    Container(SoftwareSystem, "<name>", func())
//
//    Container(SoftwareSystem, "<name>", "[description]", func())
//
//    Container(SoftwareSystem, "<name>", "[description]", "[technology]", func())
//
// Container also accepts a Goa service as second argument in which case the
// name and description are taken from the service and the technology is set to
// "Go and Goa v3"
//
//    Container(SoftwareSystem, Service)
//
//    Container(SoftwareSystem, Service, func())
//
// Example:
//
//    var _ = Workspace(func() {
//        Container(SoftwareSystem, "My system", "A system with a great architecture", "Go and Goa", func() {
//            Tag("bill processing")
//            URL("https://goa.design/mysystem")
//            Uses(OtherSystem, "Uses", "gRPC", Synchronous)
//            Delivers(Customer, "Delivers emails to", "SMTP", Synchronous)
//        })
//
//        // Alternate syntax using a Goa service.
//        Container(SoftwareSystem, Service, func() {
//            // ...
//        })
//    })
//
func Container(system *expr.SoftwareSystem, args ...interface{}) {
	_, ok := eval.Current().(*expr.Workspace)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	if len(args) == 0 {
		eval.ReportError("missing argument")
		return
	}
	var (
		name        string
		description string
		technology  string
		dsl         func()
	)
	switch a := args[0].(type) {
	case string:
		name = a
		description, technology, dsl = parseElementArgs(args[1:]...)
	case *goaexpr.ServiceExpr:
		name = a.Name
		description = a.Description
		technology = "Go and Goa v3"
		if len(args) > 1 {
			if d, ok := args[1].(func()); ok {
				dsl = d
			} else {
				eval.InvalidArgError("DSL function", args[1])
			}
		}
		if len(args) > 2 {
			eval.ReportError("too many arguments")
		}
	default:
		eval.InvalidArgError("name or Goa service", args[0])
	}

	c := &expr.Container{
		Element: &expr.Element{
			Name:        name,
			Description: description,
			Technology:  technology,
		},
		System: system,
	}
	if dsl != nil {
		eval.Execute(dsl, c)
	}
	expr.Identify(c)
	system.Containers = append(system.Containers, c)
}

// Component defines a component.
//
// Component must appear in a Workspace expression.
//
// Component takes 2 to 5 arguments. The first argument is the container that
// the component belongs to. The second argument is the component name. The name
// may be optionally followed by a description. If a description is set then it
// may be followed by the technology details used by the component. Finally
// Component may take a func() as last argument to define additional properties
// of the component.
//
// The valid syntax for Component is thus:
//
//    Component(Container, "<name>")
//
//    Component(Container, "<name>", "[description]")
//
//    Component(Container, "<name>", "[description]", "[technology]")
//
//    Component(Container, "<name>", func())
//
//    Component(Container, "<name>", "[description]", func())
//
//    Component(Container, "<name>", "[description]", "[technology]", func())
//
// Example:
//
//    var _ = Workspace(func() {
//        Component(Container, "My component", "A component", "Go and Goa", func() {
//            Tag("bill processing")
//            URL("https://goa.design/mysystem")
//            Uses(OtherSystem, "Uses", "gRPC", Synchronous)
//            Delivers(Customer, "Delivers emails to", "SMTP", Synchronous)
//        })
//    })
//
func Component(container *expr.Container, name string, args ...interface{}) {
	_, ok := eval.Current().(*expr.Workspace)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	description, technology, dsl := parseElementArgs(args...)
	c := &expr.Component{
		Element: &expr.Element{
			Name:        name,
			Description: description,
			Technology:  technology,
		},
		Container: container,
	}
	if dsl != nil {
		eval.Execute(dsl, c)
	}
	expr.Identify(c)
	container.Components = append(container.Components, c)
}

// parseElement is a helper function that parses the given element DSL
// arguments. Accepted syntax are:
//
//     "[decription]"
//     "[description]", "[technology]"
//     func()
//     "[description]", func()
//     "[description]", "[technology]", func()
//
func parseElementArgs(args ...interface{}) (description, technology string, dsl func()) {
	if len(args) == 0 {
		return
	}
	switch a := args[0].(type) {
	case string:
		description = a
	case func():
		dsl = a
	default:
		eval.InvalidArgError("description or DSL function", args[0])
	}
	if len(args) > 1 {
		if dsl != nil {
			eval.ReportError("DSL function must be last argument")
		}
		switch a := args[1].(type) {
		case string:
			technology = a
		case func():
			dsl = a
		default:
			eval.InvalidArgError("technology or DSL function", args[1])
		}
		if len(args) > 2 {
			if dsl != nil {
				eval.ReportError("DSL function must be last argument")
			}
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
	return
}
