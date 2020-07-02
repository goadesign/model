package dsl

import (
	"goa.design/goa/v3/eval"
	"goa.design/plugins/v3/structurizr/expr"

	// Register code generators for the structurizr plugin
	_ "goa.design/plugins/v3/structurizr"
)

// Workspace defines the workspace containing the models and views. Workspace
// must appear exactly once in a given design. A name must be provided if a
// description is.
//
// Workspace is a top-level DSL function.
//
// Workspace takes one to three arguments. The first argument is either a string
// or a function. If the first argument is a string then an optional description
// may be passed as second argument. The last argument must be a function that
// defines the models and views.
//
// The valid syntax for Workspace is thus:
//
//    Workspace(func())
//
//    Workspace("name", func())
//
//    Workspace("name", "description", func())
//
// Examples:
//
//    // Default workspace, no description
//    var _ = Workspace(func() {
//	      Model(func() {
//            SoftwareSystem("My Software System")
//        })
//    })
//
//    // Workspace with given name, no description
//    var _ = Workspace("name", func() {
//	      Model(func() {
//            SoftwareSystem("My Software System")
//        })
//    })
//
//    // Workspace with given name and description
//    var _ = Workspace("My Workspace", "A great architecture.", func() {
//	      Model(func() {
//            SoftwareSystem("My Software System")
//        })
//    })
//
func Workspace(args ...interface{}) {
	_, ok := eval.Current().(eval.TopExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	nargs := len(args)
	if nargs == 0 {
		eval.ReportError("missing child DSL")
		return
	}
	dsl, ok := args[nargs-1].(func())
	if !ok {
		eval.ReportError("missing child DSL (last argument must be a func)")
		return
	}
	var name, desc string
	if nargs > 1 {
		name, ok = args[0].(string)
		if !ok {
			eval.InvalidArgError("string", args[0])
		}
	}
	if nargs > 2 {
		desc, ok = args[1].(string)
		if !ok {
			eval.InvalidArgError("string", args[1])
		}
	}
	if nargs > 3 {
		eval.ReportError("too many arguments")
		return
	}
	w := &expr.WorkspaceExpr{Name: name, Description: desc, Model: &expr.ModelExpr{}}
	if !eval.Execute(dsl, w) {
		return
	}
	expr.Root = w
}
