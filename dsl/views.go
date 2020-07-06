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

// Add adds a person, an element or a relationship to a view.
//
// Add must appear in SystemLandscapeView, SystemContextView, ContainerView,
// ComponentView or DynamicView (only relationships can be added to dynamic
// views).
//
// Add takes the person, element or relationship (as defined by the source and
// destination) as first argument and an optional function as last argument.
//
//      Add(PersonOrElement)
//
//      Add(PersonOrElement, func())
//
//      Add(Source, Destination)
//
//      Add(Source, Destination, func())
//
// Example:
//
//     var _ = Workspace(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         var Person = Person("Customer", func() {
//             External()
//             Uses(System, "Sends emails", "SMTP")
//         })
//         Views(func() {
//             SystemContextView(SoftwareSystem, "context", "An overview diagram.", func() {
//                 Add(System, func() {
//                     Coord(10, 10)
//                     NoRelationships()
//                 })
//                 Add(Person, System, func() {
//                     Vertices(10, 20, 10, 40)
//                     Routing(RoutingOrthogonal)
//                     Position(45)
//                 })
//             })
//             DynamicView(SoftwareSystem, "dynamic", func() {
//                 Title("Customer flow")
//                 Add(Person, System, func() {
//                     Vertices(10, 20, 10, 40)
//                     Routing(RoutingOrthogonal)
//                     Position(45)
//                     Description("Customer sends email to support")
//                     Order("1")
//                 })
//             })
//         })
//     })
//
func Add(first interface{}, rest ...interface{}) {
	var (
		per  *expr.Person
		sys  *expr.SoftwareSystem
		cont *expr.Container
		comp *expr.Component
		rel  *expr.Relationship
		dsl  func()
	)
	if len(rest) > 0 {
		switch a := rest[0].(type) {
		case expr.ElementHolder:
			destID := a.GetElement().ID
			srcID := first.(expr.ElementHolder).GetElement().ID
			rel = expr.FindRelationship(srcID, destID)
			if rel == nil {
				eval.ReportError("no existing relationship between %s and %s.", first.(eval.Expression).EvalName(), rest[0].(eval.Expression).EvalName())
				return
			}
		case func():
			dsl = a
		default:
			eval.InvalidArgError("person, software system, container, component or function DSL", a)
			return
		}
		if len(rest) > 1 {
			if d, ok := rest[1].(func()); ok {
				dsl = d
			} else {
				eval.InvalidArgError("function", rest[1])
				return
			}
			if len(rest) > 2 {
				eval.ReportError("too many arguments")
			}
		}
	}
	if _, ok := eval.Current().(*expr.DynamicView); ok && rel == nil {
		eval.ReportError("only relationships may be added to dynamic views")
		return
	}
	switch a := first.(type) {
	case *expr.Person:
		per = a
	case *expr.SoftwareSystem:
		sys = a
	case *expr.Container:
		cont = a
	case *expr.Component:
		comp = a
	default:
		eval.InvalidArgError("person, software system, container or component", a)
		return
	}

	vh, ok := eval.Current().(expr.ViewHolder)
	if !ok {
		eval.IncompatibleDSL()
	}
	v := vh.GetView()
	switch {
	case per != nil:
		v.AddPeople(per)
	case sys != nil:
		v.AddSoftwareSystems(sys)
	case cont != nil:
		v.AddContainers(cont)
	case comp != nil:
		v.AddComponents(comp)
	case rel != nil:
		v.AddRelationships(rel)
	}

	// Execute DSL last so that relationships that may be removed are not added
	// back.
	if dsl != nil {
		switch {
		case per != nil:
			eval.Execute(dsl, v.GetElementView(per.ID))
		case sys != nil:
			eval.Execute(dsl, v.GetElementView(sys.ID))
		case cont != nil:
			eval.Execute(dsl, v.GetElementView(cont.ID))
		case comp != nil:
			eval.Execute(dsl, v.GetElementView(comp.ID))
		case rel != nil:
			eval.Execute(dsl, v.GetRelationshipView(rel.ID))
		}
	}
}

// AddAll includes all elements and relationships in scope to the view.
//
// AddAll may appear in SystemLandscapeView, SystemContextView, ContainerView,
// ComponentView, DynamicView or DeploymentView.
//
// AddAll takes no argument.
//
// Example:
//
//     var _ = Workspace(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         var Person = Person("Customer", func() {
//             External()
//             Uses(System, "Sends emails", "SMTP")
//         })
//         Views(func() {
//             SystemContextView(SoftwareSystem, "context", "An overview diagram.", func() {
//                 AddAll()
//             })
//         })
//     })
//
func AddAll() {
	model := expr.Root.Model
	switch v := eval.Current().(type) {
	case *expr.LandscapeView:
		v.AddPeople(model.People...)
		v.AddSoftwareSystems(model.Systems...)
	case *expr.ContextView:
		v.AddPeople(model.People...)
		v.AddSoftwareSystems(model.Systems...)
	case *expr.ContainerView:
		v.AddPeople(model.People...)
		v.AddSoftwareSystems(model.Systems...)
		AddContainers()
	case *expr.ComponentView:
		v.AddPeople(model.People...)
		v.AddSoftwareSystems(model.Systems...)
		AddContainers()
		AddComponents()
	default:
		eval.IncompatibleDSL()
	}
}

// AddNeighbors Adds all of the permitted elements which are directly connected
// to the specified element to the view. Permitted elements are software
// systems and people for system landscape and system context views, software
// systems, people and containers for container views and software system,
// people, containers and components for component views.
//
// AddNeighbors must appear in SystemLandscapeView, SystemContextView,
// ContainerView or ComponentView.
//
// AddNeighbors accept a single argument which is the element that should be
// added with its direct relationships. It must be a software system or person
// for system landscape and system context views, a software system, person or
// container for container views or a software system, person, container or
// component for component views.
//
// Example:
//
//     var _ = Workspace(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         var Customer = Person("Customer", func() {
//             External()
//             Uses(System, "Sends emails", "SMTP")
//         })
//         Views(func() {
//             SystemContextView(SoftwareSystem, "context", "An overview diagram.", func() {
//                 AddNeighbors(System)
//                 AddNeighbors(Customer)
//             })
//         })
//     })
//
func AddNeighbors(element interface{}) {
	var (
		elt        *expr.Element
		cont, comp bool
	)
	switch e := element.(type) {
	case *expr.Person:
		elt = e.Element
	case *expr.SoftwareSystem:
		elt = e.Element
	case *expr.Container:
		elt = e.Element
		cont = true
	case *expr.Component:
		elt = e.Element
		comp = true
	default:
		eval.InvalidArgError("person, software system, container or component", element)
		return
	}
	switch v := eval.Current().(type) {
	case *expr.LandscapeView:
		if cont || comp {
			eval.ReportError("AddNeighbors in a software landscape view must be given a software system or a person.")
			return
		}
		v.AddPeople(elt.RelatedPeople()...)
		v.AddSoftwareSystems(elt.RelatedSoftwareSystems()...)
	case *expr.ContextView:
		if cont || comp {
			eval.ReportError("AddNeighbors in a software context view must be given a software system or a person.")
			return
		}
		v.AddPeople(elt.RelatedPeople()...)
		v.AddSoftwareSystems(elt.RelatedSoftwareSystems()...)
	case *expr.ContainerView:
		if comp {
			eval.ReportError("AddNeighbors in a container view must be given a person, software system or a container.")
			return
		}
		v.AddPeople(elt.RelatedPeople()...)
		v.AddSoftwareSystems(elt.RelatedSoftwareSystems()...)
		v.AddContainers(elt.RelatedContainers()...)
	case *expr.ComponentView:
		v.AddPeople(elt.RelatedPeople()...)
		v.AddSoftwareSystems(elt.RelatedSoftwareSystems()...)
		v.AddContainers(elt.RelatedContainers()...)
		v.AddComponents(elt.RelatedComponents()...)
	default:
		eval.IncompatibleDSL()
	}
}

// AddDefault adds default elements that are relevant for the specific view:
//
//    - System landscape view: adds all software systems and people
//    - System context view: adds softare system and other related software systems
//      and people.
//    - Container view: adds all containers in software system as well as related
//      software systems and people.
//    - Component view: adds all components in container as well as related
//      containers, software systems and people.
//
// AddDefault must appear in SystemLandscapeView, SystemContextView,
// ContainerView or ComponentView.
//
// AddDefault takes no argument.
//
// Example:
//
//     var _ = Workspace(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         var Customer = Person("Customer", func() {
//             External()
//             Uses(System, "Sends emails", "SMTP")
//         })
//         Views(func() {
//             SystemContextView(SoftwareSystem, "context", "An overview diagram.", func() {
//                 AddDefault()
//             })
//         })
//     })
//
func AddDefault() {
	switch v := eval.Current().(type) {
	case *expr.LandscapeView:
		AddAll()
	case *expr.ContextView:
		AddNeighbors(expr.GetSoftwareSystem(v.SoftwareSystemID))
	case *expr.ContainerView:
		s := expr.GetSoftwareSystem(v.SoftwareSystemID)
		v.AddContainers(s.Containers...)
		for _, c := range s.Containers {
			v.AddSoftwareSystems(c.RelatedSoftwareSystems()...)
			v.AddPeople(c.RelatedPeople()...)
		}
	case *expr.ComponentView:
		c := expr.GetContainer(v.ContainerID)
		v.AddComponents(c.Components...)
		for _, c := range c.Components {
			v.AddContainers(c.RelatedContainers()...)
			v.AddSoftwareSystems(c.RelatedSoftwareSystems()...)
			v.AddPeople(c.RelatedPeople()...)
		}
	default:
		eval.IncompatibleDSL()
	}
}

// Remove given person, element or relationship from view. Alternatively remove
// all persons, elements and relationships tagged with the given tag.
//
// Remove must appear in SystemLandscapeView, SystemContextView,
// ContainerView or ComponentView.
//
// Remove takes one or two argument: the first argument must be a person, an
// element or a tag value. The second argument is needed when removing
// relationships and indicates the destination of the relationship (the first
// argument is the source in this case).
//
// Usage:
//
//     Remove(PersonOrElement)
//
//     Remove(Source, Destination)
//
//     Remove("<tag>")
//
// Example:
//
//     var _ = Workspace(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         var Customer = Person("Customer", func() {
//             External()
//             Uses(System, "Sends emails", "SMTP")
//         })
//         Container(System, "Unwanted", func() {
//             Tag("irrelevant")
//         })
//         Views(func() {
//             SystemContextView(SoftwareSystem, "context", "An overview diagram.", func() {
//                 AddDefault()
//                 Remove(Customer)
//                 Remove(Customer, System)
//                 Remove("irrelevant")
//             })
//         })
//     })
//
func Remove(e interface{}, dest ...interface{}) {
	if len(dest) > 1 {
		eval.ReportError("too many arguments")
		return
	}

	var destID string
	if len(dest) > 0 {
		if eh, ok := dest[0].(expr.ElementHolder); ok {
			destID = eh.GetElement().ID
		} else {
			eval.InvalidArgError("person, software system, container or component", dest[0])
			return
		}
	}

	var id, tag string
	switch a := e.(type) {
	case expr.ElementHolder:
		id = a.GetElement().ID
	case string:
		tag = a
	default:
		eval.InvalidArgError("string, person, software system, container or component", e)
		return
	}
	if destID != "" {
		if tag != "" {
			eval.ReportError("only one argument allowed when using a tag as first argument")
			return
		}
		if r := expr.FindRelationship(id, destID); r != nil {
			id = r.ID
		} else {
			eval.ReportError("no existing relationship with source %s and destination %s", e.(eval.Expression).EvalName(), dest[0].(eval.Expression).EvalName())
			return
		}
	}

	if vh, ok := eval.Current().(expr.ViewHolder); ok {
		if id != "" {
			vh.GetView().Remove(id)
		} else {
			vh.GetView().RemoveTagged(tag)
		}
	} else {
		eval.IncompatibleDSL()
	}
}

// RemoveUnreachable removes all elements and people that cannot be reached by
// traversing the graph of relationships starting with the given element or
// person.
//
// RemoveUnreachable must appear in SystemLandscapeView, SystemContextView,
// ContainerView or ComponentView.
//
// RemoveUnreachable takes one argument: the person or element from which the
// graph traversal should be initiated.
//
// Example:
//
//     var _ = Workspace(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         var OtherSystem = SoftwareSystem("Other software System")
//         var Customer = Person("Customer", func() {
//             External()
//             Uses(System, "Sends emails", "SMTP")
//         })
//         Views(func() {
//             SystemContextView(SoftwareSystem, "context", "An overview diagram.", func() {
//                 AddDefault()
//                 RemoveUnreachable(System) // Removes OtherSystem
//             })
//         })
//     })
//
func RemoveUnreachable(e interface{}) {
	var elt *expr.Element
	if eh, ok := e.(expr.ElementHolder); ok {
		elt = eh.GetElement()
	} else {
		eval.InvalidArgError("person, software system, container or component", e)
		return
	}
	if vh, ok := eval.Current().(expr.ViewHolder); ok {
		vh.GetView().RemoveUnreachable(elt)
	} else {
		eval.IncompatibleDSL()
	}
}

// RemoveUnrelated removes all elements that have no relationship to other
// elements in the view.
//
// RemoveUnrelated must appear in SystemLandscapeView, SystemContextView,
// ContainerView or ComponentView.
//
// RemoveUnrelated takes no argument.
//
// Example:
//
//     var _ = Workspace(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         var OtherSystem = SoftwareSystem("Other software System")
//         var Customer = Person("Customer", func() {
//             External()
//             Uses(System, "Sends emails", "SMTP")
//         })
//         Views(func() {
//             SystemContextView(SoftwareSystem, "context", "An overview diagram.", func() {
//                 AddDefault()
//                 RemoveUnrelated()) // Removes OtherSystem
//             })
//         })
//     })
//
func RemoveUnrelated() {
	if vh, ok := eval.Current().(expr.ViewHolder); ok {
		vh.GetView().RemoveUnrelated()
	} else {
		eval.IncompatibleDSL()
	}
}

// AutoLayout enables automatic layout mode for the diagram. The
// first argument indicates the rank direction, it must be one of
// RankTopBottom, RankBottomTop, RankLeftRight or RankRightLeft
//
// AutoLayout must appear in SystemLandscapeView, SystemContextView,

// ContainerView or ComponentView.

// AddContainers includes all containers in scope to the view.
//
// AddContainers may appear in ContainerView or ComponentView.
//
// AddContainers takes no argument.
//
func AddContainers() {
	switch v := eval.Current().(type) {
	case *expr.ContainerView:
		v.AddContainers(expr.GetSoftwareSystem(v.SoftwareSystemID).Containers...)
	case *expr.ComponentView:
		c := expr.GetContainer(v.ContainerID)
		v.AddContainers(c.System.Containers...)
	default:
		eval.IncompatibleDSL()
	}
}

// AddComponents includes all components in scope to the view.
//
// AddComponents must appear in ComponentView.
//
// AddComponents takes no argument
//
func AddComponents() {
	if cv, ok := eval.Current().(*expr.ComponentView); ok {
		cv.AddComponents(expr.GetContainer(cv.ContainerID).Components...)
		return
	}
	eval.IncompatibleDSL()
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
		eval.InvalidArgError("string or function", args[0])
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
			eval.InvalidArgError("string or function", args[1])
		}
		if len(args) > 2 {
			if dsl != nil {
				eval.ReportError("DSL function must be last argument")
			}
			if d, ok := args[2].(func()); ok {
				dsl = d
			} else {
				eval.InvalidArgError("function", args[2])
			}
			if len(args) > 3 {
				eval.ReportError("too many arguments")
			}
		}
	}
	return
}
