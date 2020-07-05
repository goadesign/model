package dsl

import (
	"goa.design/goa/v3/eval"
	"goa.design/structurizr/expr"
)

// Global is the keyword used to define dynamic views with global scope. See
// DynamicView.
const Global = 0

// Views defines one or more views.
//
// Views takes one argument: the function that defines the views.
//
// Views must appear in Workspace.
//
// Example:
//
//     var _ = Workspace(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         Views(func() {
//             SystemContext(MySystem, "SystemContext", "An example of a System Context diagram.", func() {
//                 AddAll()
//                 AutoLayout()
//             })
//         })
//     })
//
func Views(dsl func()) {
	w, ok := eval.Current().(*expr.Workspace)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	w.Views = &expr.Views{DSL: dsl}
}

// SystemLandscapeView defines a system landscape view.
//
// SystemLandscapeView must appear in Views.
//
// SystemLandscapeView accepts 1 to 3 arguments: the first argument is an optional
// key for the view which can be used to reference it when creating a fltered
// views. The second argument is an optional description, the key must be
// provided when giving a description. The last argument is a function
// describing the properties of the view.
//
// Valid usage of SystemLandscapeView are thus:
//
//    SystemLandscapeView(func())
//
//    SystemLandscapeView("[key]", func())
//
//    SystemLandscapeView("[key]", "[description]", func())
//
// Example:
//
//     var _ = Workspace(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         Views(func() {
//             SystemLandscapeView("landscape", "An overview diagram.", func() {
//                 Title("Overview of system")
//                 AddAll()
//                 Remove(Container3)
//                 AutoLayout()
//                 AnimationStep(Container1, Container2)
//                 PaperSize(SizeSlide4X3)
//                 EnterpriseBoundaryVisible()
//             })
//         })
//     })
//
func SystemLandscapeView(args ...interface{}) {
	vs, ok := eval.Current().(*expr.Views)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	key, description, dsl := parseView(args...)
	v := &expr.LandscapeView{
		View: expr.View{
			Key:         key,
			Description: description,
		},
	}
	if dsl != nil {
		eval.Execute(dsl, v)
	}
	vs.LandscapeViews = append(vs.LandscapeViews, v)
}

// SystemContextView defines a system context view.
//
// SystemContextView must appear in Views.
//
// SystemContextView accepts 2 to 4 arguments: the first argument is the system
// the view applies to. The second argument is an optional key for the view
// which can be used to reference it when creating a fltered views. The third
// argument is an optional description, the key must be provided when giving a
// description. The last argument is a function describing the properties of the
// view.
//
// Valid usage of SystemContextView are thus:
//
//    SystemContextView(SoftwareSystem, func())
//
//    SystemContextView(SoftwareSystem, "[key]", func())
//
//    SystemContextView(SoftwareSystem, "[key]", "[description]", func())
//
// Example:
//
//     var _ = Workspace(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         Views(func() {
//             SystemContextView(SoftwareSystem, "context", "An overview diagram.", func() {
//                 Title("Overview of system")
//                 AddAll()
//                 Remove(Container3)
//                 AutoLayout()
//                 AnimationStep(Container1, Container2)
//                 PaperSize(SizeSlide4X3)
//                 EnterpriseBoundaryVisible()
//             })
//         })
//     })
//
func SystemContextView(s *expr.SoftwareSystem, args ...interface{}) {
	vs, ok := eval.Current().(*expr.Views)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	key, description, dsl := parseView(args...)
	v := &expr.ContextView{
		View: expr.View{
			Key:         key,
			Description: description,
		},
		SoftwareSystemID: s.ID,
	}
	if dsl != nil {
		eval.Execute(dsl, v)
	}
	vs.ContextViews = append(vs.ContextViews, v)
}

// ContainerView defines a container view.
//
// ContainerView must appear in Views.
//
// ContainerView accepts 2 to 4 arguments: the first argument is the software
// system the container view applies to. The second argumetn is an optional key
// for the view which can be used to reference it when creating a fltered views.
// The third argument is an optional description, the key must be provided when
// giving a description. The last argument is a function describing the
// properties of the view.
//
// Valid usage of ContainerView are thus:
//
//    ContainerView(SoftwareSystem, func())
//
//    ContainerView(SoftwareSystem, "[key]", func())
//
//    ContainerView(SoftwareSystem, "[key]", "[description]", func())
//
// Example:
//
//     var _ = Workspace(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         Views(func() {
//             ContainerView(SoftwareSystem, "container", "An overview diagram.", func() {
//                 Title("System containers")
//                 AddAll()
//                 Remove(Container3)
//                 // Alternatively to AddAll + Remove: Add
//                 AutoLayout()
//                 AnimationStep(Container1, Container2)
//                 PaperSize(SizeSlide4X3)
//                 SystemBoundariesVisible()
//             })
//         })
//     })
//
func ContainerView(s *expr.SoftwareSystem, args ...interface{}) {
	vs, ok := eval.Current().(*expr.Views)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	key, description, dsl := parseView(args...)
	v := &expr.ContainerView{
		View: expr.View{
			Key:         key,
			Description: description,
		},
		SoftwareSystemID: s.ID,
	}
	if dsl != nil {
		eval.Execute(dsl, v)
	}
	vs.ContainerViews = append(vs.ContainerViews, v)
}

// ComponentView defines a component view.
//
// ComponentView must appear in Views.
//
// ComponentView accepts 2 to 4 arguments: the first argument is the container
// being described by the component view. The second argument is an optional key
// for the view which can be used to reference it when creating a fltered views.
// The third argument is an optional description, the key must be provided when
// giving a description. The last argument is a function describing the
// properties of the view.
//
// Valid usage of ComponentView are thus:
//
//    ComponentView(Container, func())
//
//    ComponentView(Container, "[key]", func())
//
//    ComponentView(Container, "[key]", "[description]", func())
//
// Example:
//
//     var _ = Workspace(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         Views(func() {
//             ComponentView(Container, "component", "An overview diagram.", func() {
//                 Title("Overview of container")
//                 AddAll()
//                 Remove(Component3)
//                 AutoLayout()
//                 AnimationStep(Component1, Component2)
//                 PaperSize(SizeSlide4X3)
//                 ContainerBoundariesVisible()
//             })
//         })
//     })
//
func ComponentView(c *expr.Container, args ...interface{}) {
	vs, ok := eval.Current().(*expr.Views)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	key, description, dsl := parseView(args...)
	v := &expr.ComponentView{
		View: expr.View{
			Key:         key,
			Description: description,
		},
		ContainerID: c.ID,
	}
	if dsl != nil {
		eval.Execute(dsl, v)
	}
	vs.ComponentViews = append(vs.ComponentViews, v)
}

// FilteredView defines a filtered view on top of the specified view.
// The base key specifies the key of the System Landscape, System
// Context, Container, or Component view on which this filtered view
// should be based.
//
// FilteredView must appear in Views.
//
// FilteredView accepts 2 arguments: the view being filtered and a function
// describing additional properties.
//
// Example:
//
//     var _ = Workspace(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         Views(func() {
//             SystemContextView(SoftwareSystem, "context", "An overview diagram.", func() {
//                 AddAll()
//                 AutoLayout()
//             })
//             FilteredView(SystemContextView, func() {
//                 FilterTag("infra")
//                 Exclude()
//             })
//         })
//     })
//
func FilteredView(view interface{}, dsl func()) {
	vs, ok := eval.Current().(*expr.Views)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	var key string
	switch v := view.(type) {
	case *expr.LandscapeView:
		key = v.Key
	case *expr.ContextView:
		key = v.Key
	case *expr.ContainerView:
		key = v.Key
	case *expr.ComponentView:
		key = v.Key
	default:
		eval.IncompatibleDSL()
		return
	}
	if key == "" {
		eval.ReportError("Filtered view applied on a view with no key. Make sure the view given as argument defines a key.")
		return
	}
	fv := &expr.FilteredView{BaseKey: key}
	eval.Execute(dsl, fv)
	vs.FilteredViews = append(vs.FilteredViews, fv)
}

// DynamicView defines a Dynamic view for the specified scope. The
// first argument defines the scope of the view, and therefore what can
// be added to the view, as follows:
//
//   * Global scope: People and software systems.
//   * Software system scope: People, other software systems, and
//     containers belonging to the software system.
//   * Container scope: People, other software systems, other
//     containers, and components belonging to the container.
//
// DynamicView must appear in Views.
//
// DynamicView accepts 2 to 4 arguments: the first argument is the scope: either
// the keyword 'Global' or a software system or container identifier. The second
// argument is an optional key for the view. The third argument is an optional
// description, the key must be provided when giving a description. The last
// argument is a function describing the properties of the view.
//
// A dynamic view is created by specifying relationships that should be
// rendered. See Relationship for additional information.
//
// Valid usage of DynamicView are thus:
//
//    DynamicView(Scope, func())
//
//    DynamicView(Scope, "[key]", func())
//
//    DynamicView(Scope, "[key]", "[description]", func())
//
// Example:
//
//     var _ = Workspace(func() {
//         var FirstSystem = SoftwareSystem("First system")
//         var SecondSystem = SoftwareSystem("Second system", func() {
//             Uses(FirstSystem, "Uses")
//         })
//         Views(func() {
//             DynamicView(Global, "dynamic", "A dynamic diagram.", func() {
//                 Title("Overview of system")
//                 AutoLayout()
//                 PaperSize(SizeSlide4X3)
//                 Relationship(SecondSystem, FirstSystem)
//             })
//         })
//     })
//
func DynamicView(scope interface{}, args ...interface{}) {
	vs, ok := eval.Current().(*expr.Views)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	var id string
	switch s := scope.(type) {
	case int:
		id = "" // Global scope
	case *expr.SoftwareSystem:
		id = s.ID
	case *expr.Container:
		id = s.ID
	default:
		eval.IncompatibleDSL()
		return
	}
	key, description, dsl := parseView(args...)
	v := &expr.DynamicView{
		View: expr.View{
			Key:         key,
			Description: description,
		},
		ElementID: id,
	}
	if dsl != nil {
		eval.Execute(dsl, v)
	}
	vs.DynamicViews = append(vs.DynamicViews, v)
}

// DeploymentView defines a Deployment view for the specified scope and
// deployment environment. The first argument defines the scope of the
// view, and the second argument defines the deployment environment. The
// combination of these two arguments determines what can be added to
// the view, as follows:
//
//   * Global scope: All deployment nodes, infrastructure nodes, and
//     container instances within the deployment environment.
//   * Software system scope: All deployment nodes and infrastructure
//     nodes within the deployment environment. Container instances within
//     the deployment environment that belong to the software system.
//
// DeploymentView must appear in Views.
//
// DeploymentView accepts 3 to 5 arguments: the first argument is the scope:
// either the keyword 'Global' or a software system. The second argument is the
// name of the environment. The third argument is an optional key for the view.
// The fourth argument is an optional description, the key must be provided when
// giving a description. The last argument is a function describing the
// properties of the view.
//
// Valid usage of DeploymentView are thus:
//
//    DeploymentView(Scope, "<environment>", func())
//
//    DeploymentView(Scope, "<environment>", "[key]", func())
//
//    DeploymentView(Scope, "<environment>", "[key]", "[description]", func())
//
// Example:
//
//     var _ = Workspace(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         Views(func() {
//             DeploymentView(Global, "Production", "deployment", "A deployment overview diagram.", func() {
//                 Title("Overview of deployment")
//                 AddAll()
//                 Remove(Container3)
//                 AutoLayout()
//                 AnimationStep(Container1, Container2)
//                 PaperSize(SizeSlide4X3)
//             })
//         })
//     })
//
func DeploymentView(scope interface{}, env string, args ...interface{}) {
	vs, ok := eval.Current().(*expr.Views)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	key, description, dsl := parseView(args...)
	var id string
	switch s := scope.(type) {
	case int:
		id = "" // Global scope
	case *expr.SoftwareSystem:
		id = s.ID
	default:
		eval.IncompatibleDSL()
		return
	}
	v := &expr.DeploymentView{
		View: expr.View{
			Key:         key,
			Description: description,
		},
		SoftwareSystemID: id,
		Environment:      env,
	}
	if dsl != nil {
		eval.Execute(dsl, v)
	}
	vs.DeploymentViews = append(vs.DeploymentViews, v)
}

// Title sets the view diagram title.
//
// Title may appear in SystemLandscapeView, SystemContextView, ContainerView,
// ComponentView, DynamicView or DeploymentView.
//
// Title accepts one argument: the view title.
func Title(t string) {
	switch v := eval.Current().(type) {
	case *expr.LandscapeView:
		v.Title = t
	case *expr.ContextView:
		v.Title = t
	case *expr.ContainerView:
		v.Title = t
	case *expr.ComponentView:
		v.Title = t
	case *expr.DynamicView:
		v.Title = t
	case *expr.DeploymentView:
		v.Title = t
	default:
		eval.IncompatibleDSL()
	}
}

// AddAll includes all elements and relationships in scope to the view.
//
// AddAll may appear in SystemLandscapeView, SystemContextView, ContainerView,
// ComponentView, DynamicView or DeploymentView.
//
// AddAll takes no argument.
//
func AddAll() {
	model := expr.Root.Model
	switch v := eval.Current().(type) {
	case *expr.LandscapeView:
		v.Merge(elementViews(model.PeopleElements()))
		v.Merge(elementViews(model.SystemElements()))
	case *expr.ContextView:
		v.Merge(elementViews(model.PeopleElements()))
		v.Merge(elementViews(model.SystemElements()))
	case *expr.ContainerView:
		v.Merge(elementViews(model.PeopleElements()))
		v.Merge(elementViews(model.SystemElements()))
		cs := expr.Registry[v.SoftwareSystemID].(*expr.SoftwareSystem).ContainerElements()
		v.Merge(elementViews(cs))
	case *expr.ComponentView:
	case *expr.DynamicView:
	case *expr.DeploymentView:
	default:
		eval.IncompatibleDSL()
	}
}

// elementViews is a helper method that converts the given elements to views.
func elementViews(elements []*expr.Element) []*expr.ElementView {
	res := make([]*expr.ElementView, len(elements))
	for i, e := range elements {
		res[i] = &expr.ElementView{ID: e.ID}
	}
	return res
}

// parseView is a helper function that parses the given view DSL
// arguments. Accepted syntax are:
//
//     func()
//     "[key]", func()
//     "[key]", "[description]", func()
//
func parseView(args ...interface{}) (key, description string, dsl func()) {
	if len(args) == 0 {
		eval.ReportError("missing argument")
		return
	}
	switch a := args[0].(type) {
	case string:
		key = a
	case func():
		dsl = a
	default:
		eval.InvalidArgError("key or DSL function", args[0])
	}
	if len(args) > 1 {
		if dsl != nil {
			eval.ReportError("DSL function must be last argument")
		}
		switch a := args[1].(type) {
		case string:
			description = a
		case func():
			dsl = a
		default:
			eval.InvalidArgError("desciption or DSL function", args[1])
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
