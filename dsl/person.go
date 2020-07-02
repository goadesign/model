package dsl

import (
	"goa.design/goa/v3/eval"
	"goa.design/plugins/v3/structurizr/expr"
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
func Person(args ...interface{}) {
	w, ok := eval.Current().(*expr.WorkspaceExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	if len(args) == 0 {
		eval.ReportError("missing argument")
		return
	}
	var (
		name, desc string
		dsl        func()
	)
	{
		name, ok = args[0].(string)
		if !ok {
			eval.InvalidArgError("name", args[0])
			return
		}
		if len(args) > 1 {
			desc, ok = args[1].(string)
		}
		dsl, ok = args[len(args)-1].(func())
	}
	p := &expr.PersonExpr{Name: name, Description: desc}
	if dsl != nil && !eval.Execute(dsl, p) {
		return
	}
	w.Model.People = append(w.Model.People, p)
}
