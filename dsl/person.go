package dsl

import (
	"goa.design/goa/v3/eval"
	"goa.design/structurizr/expr"
)

// Person defines a person (user, actor, role or persona).
//
// Person must appear in a Model expression.
//
// Person takes one to three arguments. The first argument is the name of the
// person. An optional description may be passed as second argument. The last
// argument may be a function that defines tags associated with the Person.
//
// The valid syntax for Person is thus:
//
//    Person("name")
//
//    Person("name", "description")
//
//    Person("name", func())
//
//    Person("name", "description", func())
//
// Example:
//
//    var _ = Workspace(func() {
//        Person("Employee")
//        Person("Customer", "A customer", func() {
//            Tag("system")
//            External()
//            URL("https://acme.com/docs/customer.html")
//            Uses(System)
//            InteractsWith(Employee)
//        })
//    })
//
func Person(name string, args ...interface{}) {
	w, ok := eval.Current().(*expr.Workspace)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	var (
		desc string
		dsl  func()
	)
	if len(args) > 0 {
		switch a := args[0].(type) {
		case string:
			desc = a
		case func():
			dsl = a
		default:
			eval.InvalidArgError("description or DSL function", args[0])
		}
		if len(args) > 1 {
			if dsl != nil {
				eval.ReportError("DSL function must be last argument")
			}
			dsl, ok = args[1].(func())
			if !ok {
				eval.InvalidArgError("DSL function", args[1])
			}
			if len(args) > 2 {
				eval.ReportError("too many arguments")
			}
		}
	}
	p := &expr.Person{
		Element: &expr.Element{
			Name:        name,
			Description: desc,
		},
		Location: expr.LocationInternal,
	}
	if dsl != nil {
		eval.Execute(dsl, p)
	}
	expr.Identify(p)
	w.Model.People = append(w.Model.People, p)
}
