package dsl

import (
	"strings"

	"goa.design/goa/v3/eval"
	"goa.design/model/expr"
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
//	Person("name")
//
//	Person("name", "description")
//
//	Person("name", func())
//
//	Person("name", "description", func())
//
// Example:
//
//	var _ = Design(func() {
//	    Person("Employee")
//	    Person("Customer", "A customer", func() {
//	        Tag("system")
//	        External()
//	        URL("https://acme.com/docs/customer.html")
//	        Uses(System)
//	        InteractsWith(Employee)
//	    })
//	})
func Person(name string, args ...any) *expr.Person {
	w, ok := eval.Current().(*expr.Design)
	if !ok {
		eval.IncompatibleDSL()
		return nil
	}
	if strings.Contains(name, "/") {
		eval.ReportError("Person: name cannot include slashes")
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
				eval.ReportError("Person: DSL function must be last argument")
			}
			dsl, ok = args[1].(func())
			if !ok {
				eval.InvalidArgError("DSL function", args[1])
			}
			if len(args) > 2 {
				eval.ReportError("Person: too many arguments")
			}
		}
	}
	p := &expr.Person{
		Element: &expr.Element{
			Name:        name,
			Description: desc,
			DSLFunc:     dsl,
		},
	}
	return w.Model.AddPerson(p)
}
