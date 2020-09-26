package dsl

import (
	"fmt"
	"strconv"
	"strings"

	"goa.design/goa/v3/eval"
	"goa.design/model/expr"
)

type (
	// PaperSizeKind is the enum for possible paper kinds.
	PaperSizeKind int

	// RoutingKind is the enum for possible routing algorithms.
	RoutingKind int

	// RankDirectionKind is the enum for possible automatic layout rank
	// directions.
	RankDirectionKind int
)

// Global is the keyword used to define dynamic views with global scope. See
// DynamicView.
const Global = 0

const (
	// RankTopBottom indicates a layout that uses top to bottom rank.
	RankTopBottom RankDirectionKind = iota + 1
	// RankBottomTop indicates a layout that uses bottom to top rank.
	RankBottomTop
	// RankLeftRight indicates a layout that uses left to right rank.
	RankLeftRight
	// RankRightLeft indicates a layout that uses right to left rank.
	RankRightLeft
)

const (
	// SizeA0Landscape defines a render page size of A0 in landscape mode (46-13/16 x 33-1/8).
	SizeA0Landscape PaperSizeKind = iota + 1
	// SizeA0Portrait defines a render page size of A0 in portrait mode (33-1/8 x 46-13/16).
	SizeA0Portrait
	// SizeA1Landscape defines a render page size of A1 in landscape mode (33-1/8 x 23-3/8).
	SizeA1Landscape
	// SizeA1Portrait defines a render page size of A1 in portrait mode (23-3/8 x 33-1/8).
	SizeA1Portrait
	// SizeA2Landscape defines a render page size of A2 in landscape mode (23-3/8 x 16-1/2).
	SizeA2Landscape
	// SizeA2Portrait defines a render page size of A2 in portrait mode (16-1/2 x 23-3/8).
	SizeA2Portrait
	// SizeA3Landscape defines a render page size of A3 in landscape mode (16-1/2 x 11-3/4).
	SizeA3Landscape
	// SizeA3Portrait defines a render page size of A3 in portrait mode (11-3/4 x 16-1/2).
	SizeA3Portrait
	// SizeA4Landscape defines a render page size of A4 in landscape mode (11-3/4 x 8-1/4).
	SizeA4Landscape
	// SizeA4Portrait defines a render page size of A4 in portrait mode (8-1/4 x 11-3/4).
	SizeA4Portrait
	// SizeA5Landscape defines a render page size of A5 in landscape mode (8-1/4  x 5-7/8).
	SizeA5Landscape
	// SizeA5Portrait defines a render page size of A5 in portrait mode (5-7/8 x 8-1/4).
	SizeA5Portrait
	// SizeA6Landscape defines a render page size of A6 in landscape mode (4-1/8 x 5-7/8).
	SizeA6Landscape
	// SizeA6Portrait defines a render page size of A6 in portrait mode (5-7/8 x 4-1/8).
	SizeA6Portrait
	// SizeLegalLandscape defines a render page size of Legal in landscape mode (14 x 8-1/2).
	SizeLegalLandscape
	// SizeLegalPortrait defines a render page size of Legal in portrait mode (8-1/2 x 14).
	SizeLegalPortrait
	// SizeLetterLandscape defines a render page size of Letter in landscape mode (11 x 8-1/2).
	SizeLetterLandscape
	// SizeLetterPortrait defines a render page size of Letter in portrait mode (8-1/2 x 11).
	SizeLetterPortrait
	// SizeSlide16X10 defines a render page size ratio of 16 x 10.
	SizeSlide16X10
	// SizeSlide16X9 defines a render page size ratio of 16 x 9.
	SizeSlide16X9
	// SizeSlide4X3 defines a render page size ratio of 4 x 3.
	SizeSlide4X3
)

const (
	// RoutingDirect draws straight lines between ends of relationships.
	RoutingDirect RoutingKind = iota + 1
	// RoutingOrthogonal draws lines with right angles between ends of relationships.
	RoutingOrthogonal
	// RoutingCurved draws curved lines between ends of relationships.
	RoutingCurved
)

// Views defines one or more views.
//
// Views takes one argument: the function that defines the views.
//
// Views must appear in Design.
//
// Example:
//
//     var _ = Design(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         Views(func() {
//             SystemContext(System, "SystemContext", "An example of a System Context diagram.", func() {
//                 AddAll()
//                 AutoLayout()
//             })
//         })
//     })
//
func Views(dsl func()) {
	w, ok := eval.Current().(*expr.Design)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	if existing := w.Views.DSLFunc; existing != nil {
		newdsl := dsl
		dsl = func() { existing(); newdsl() }
	}
	w.Views.DSLFunc = dsl
}

// SystemLandscapeView defines a system landscape view.
//
// SystemLandscapeView must appear in Views.
//
// SystemLandscapeView accepts 2 to 3 arguments: the first argument is a unique
// key for the view which can be used to reference it when creating a filtered
// views. The second argument is an optional description. The last argument is a
// function describing the properties of the view.
//
// Usage:
//
//    SystemLandscapeView("<key>", func())
//
//    SystemLandscapeView("<key>", "[description]", func())
//
// Example:
//
//     var _ = Design(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         var OtherSystem = SoftwareSystem("Other System")
//         Views(func() {
//             SystemLandscapeView("landscape", "An overview diagram.", func() {
//                 Title("Overview of system")
//                 AddAll()
//                 Remove(OtherSystem)
//                 AutoLayout()
//                 AnimationStep(System)
//                 PaperSize(SizeSlide4X3)
//                 EnterpriseBoundaryVisible()
//             })
//         })
//     })
//
func SystemLandscapeView(key string, args ...interface{}) {
	vs, ok := eval.Current().(*expr.Views)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	description, dsl, err := parseView(args...)
	if err != nil {
		eval.ReportError("SystemLandscapeView: " + err.Error())
		return
	}
	v := &expr.LandscapeView{
		ViewProps: &expr.ViewProps{
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
// SystemContextView accepts 3 to 4 arguments: the first argument is the system
// or the name of the system the view applies to. The second argument is a
// unique key for the view which can be used to reference it when creating a
// fltered views. The third argument is an optional description. The last
// argument is a function describing the properties of the view.
//
// Usage:
//
//    SystemContextView(SoftwareSystem, "<key>", func())
//
//    SystemContextView("<Software System>", "<key>", func())
//
//    SystemContextView(SoftwareSystem, "<key>", "[description]", func())
//
//    SystemContextView("<Software System>", "<key>", "[description]", func())
//
// Example:
//
//     var _ = Design(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         var OtherSystem = SoftwareSystem("Other System")
//         Views(func() {
//             SystemContextView(System, "context", "An overview diagram.", func() {
//                 Title("Overview of system")
//                 AddAll()
//                 Remove(OtherSystem)
//                 AutoLayout()
//                 AnimationStep(System)
//                 PaperSize(SizeSlide4X3)
//                 EnterpriseBoundaryVisible()
//             })
//         })
//     })
//
func SystemContextView(system interface{}, key string, args ...interface{}) {
	vs, ok := eval.Current().(*expr.Views)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	var sid string
	switch s := system.(type) {
	case *expr.SoftwareSystem:
		sid = s.ID
	case string:
		sys := expr.Root.Model.SoftwareSystem(s)
		if sys == nil {
			eval.ReportError("SystemContextView: no software system named %q", s)
			return
		}
		sid = sys.ID
	default:
		eval.InvalidArgError("software system or software system name", system)
		return
	}
	description, dsl, err := parseView(args...)
	if err != nil {
		eval.ReportError("SystemContextView: " + err.Error())
		return
	}
	v := &expr.ContextView{
		ViewProps: &expr.ViewProps{
			Key:         key,
			Description: description,
		},
		SoftwareSystemID: sid,
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
// ContainerView accepts 3 to 4 arguments: the first argument is the software
// system or the name of the software system the container view applies to. The
// second argument is a unique key which can be used to reference the view when
// creating a filtered views. The third argument is an optional description. The
// last argument is a function describing the properties of the view.
//
// Usage:
//
//    ContainerView(SoftwareSystem, "<key>", func())
//
//    ContainerView("<Software System>", "<key>", func())
//
//    ContainerView(SoftwareSystem, "<key>", "[description]", func())
//
//    ContainerView("<Software System>", "<key>", "[description]", func())
//
// Example:
//
//     var _ = Design(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         var OtherSystem = SoftwareSystem("Other System")
//         Views(func() {
//             ContainerView(SoftwareSystem, "container", "An overview diagram.", func() {
//                 Title("System containers")
//                 AddAll()
//                 Remove(OtherSystem)
//                 // Alternatively to AddAll + Remove: Add
//                 AutoLayout()
//                 AnimationStep(System)
//                 PaperSize(SizeSlide4X3)
//                 SystemBoundariesVisible()
//             })
//         })
//     })
//
func ContainerView(system interface{}, key string, args ...interface{}) {
	vs, ok := eval.Current().(*expr.Views)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	var sid string
	switch s := system.(type) {
	case *expr.SoftwareSystem:
		sid = s.ID
	case string:
		sys := expr.Root.Model.SoftwareSystem(s)
		if sys == nil {
			eval.ReportError("ContainerView: no software system named %q", s)
			return
		}
		sid = sys.ID
	default:
		eval.InvalidArgError("software system or software system name", system)
		return
	}
	description, dsl, err := parseView(args...)
	if err != nil {
		eval.ReportError("ContainerView: " + err.Error())
		return
	}
	v := &expr.ContainerView{
		ViewProps: &expr.ViewProps{
			Key:         key,
			Description: description,
		},
		SoftwareSystemID: sid,
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
// ComponentView accepts 3 to 4 arguments: the first argument is the container
// or the path to the container being described by the component view. The path
// consists of the name of the software system that contains the container followed
// by a slash followed by the name of the container. The following argument is a
// unique key which can be used to reference the view when creating a filtered
// views. Next is an optional description. The last argument must be a function
// describing the properties of the view.
//
// Usage:
//
//    ComponentView(Container, "<key>", func())
//
//    ComponentView("<Software System>/<Container>", "<key>", func())
//
//    ComponentView(Container, "<key>", "[description]", func())
//
//    ComponentView("<Software System>/<Container>", "<key>", "[description]", func())
//
// Example:
//
//     var _ = Design(func() {
//         SoftwareSystem("Software System", "My software system.", func() {
//             Container("Container", func() {
//                 Uses("Other System")
//             })
//         })
//         SoftwareSystem("Other System")
//         Views(func() {
//             ComponentView("Software System/Container", "component", "An overview diagram.", func() {
//                 Title("Overview of container")
//                 AddAll()
//                 Remove("Other System")
//                 AutoLayout()
//                 AnimationStep("Software System")
//                 PaperSize(SizeSlide4X3)
//                 ContainerBoundariesVisible()
//             })
//         })
//     })
//
func ComponentView(container interface{}, key string, args ...interface{}) {
	vs, ok := eval.Current().(*expr.Views)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	var c *expr.Container
	switch a := container.(type) {
	case *expr.Container:
		c = a
	case string:
		cont, err := expr.Root.Model.FindElement(nil, a)
		if err != nil {
			eval.ReportError("ComponentView: " + err.Error())
			return
		}
		c, ok = cont.(*expr.Container)
		if !ok {
			eval.ReportError("ComponentView: %q is not a container", a)
			return
		}
	}
	description, dsl, err := parseView(args...)
	if err != nil {
		eval.ReportError("ComponentView: " + err.Error())
		return
	}
	v := &expr.ComponentView{
		ViewProps: &expr.ViewProps{
			Key:         key,
			Description: description,
		},
		ContainerID: c.GetElement().ID,
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
//     var _ = Design(func() {
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
	if v, ok := view.(expr.View); ok {
		key = v.Props().Key
	} else {
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

// Exclude indicates that the filtered view should include all elements except
// the ones identified through the filter tag.
//
// Exclude must appear in FilteredView
//
// Exclude takes no argument.
func Exclude() {
	if v, ok := eval.Current().(*expr.FilteredView); ok {
		v.Exclude = true
		return
	}
	eval.IncompatibleDSL()
}

// FilterTag defines the set of tags to include or exclude (when Exclude() is
// used) elements and relationships when rendering the filtered view.
//
// FilterTag must appear in FilteredView
//
// FilterTag takes the list of tags as arguments. Multiple calls to FilterTag
// accumulate the tags.
func FilterTag(tag string, tags ...string) {
	if v, ok := eval.Current().(*expr.FilteredView); ok {
		v.FilterTags = append(v.FilterTags, tag)
		v.FilterTags = append(v.FilterTags, tags...)
		return
	}
	eval.IncompatibleDSL()
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
// DynamicView accepts 3 to 4 arguments. The first argument defines the scope:
// either the keyword 'Global', a software system, a software system name, a
// container or a container path. The path to a container is the name of the
// parent software system followed by a slash and the name of the container. The
// following argument is a unique key for the view. Next is an optional
// description. The last argument is a function describing the properties of the
// view.
//
// A dynamic view is created by specifying the relationships that should be
// rendered via Link.
//
// Usage:
//
//    DynamicView(Scope, "<key>", func())
//
//    DynamicView(Scope, "<key>", "[description]", func())
//
//    DynamicView("SoftwareSystem/Container", "<key>", func())
//
//    DynamicView("SoftwareSystem/Container", "<key>", "[description]", func())
//
// Where Scope is 'Global', a software system, a software system name or a
// container.
//
// Example:
//
//     var _ = Design(func() {
//         var FirstSystem = SoftwareSystem("First system")
//         var SecondSystem = SoftwareSystem("Second system", func() {
//             Uses(FirstSystem, "Uses")
//         })
//         Views(func() {
//             DynamicView(Global, "dynamic", "A dynamic diagram.", func() {
//                 Title("Overview of system")
//                 Link(FirstSystem, SecondSystem)
//                 AutoLayout()
//                 PaperSize(SizeSlide4X3)
//             })
//         })
//     })
//
func DynamicView(scope interface{}, key string, args ...interface{}) {
	vs, ok := eval.Current().(*expr.Views)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	var (
		id string
	)
	switch s := scope.(type) {
	case int:
		id = "" // Global scope
	case *expr.SoftwareSystem, *expr.Container:
		id = s.(expr.ElementHolder).GetElement().ID
	case string:
		elems := strings.Split(s, "/")
		switch len(elems) {
		case 1:
			if s := expr.Root.Model.SoftwareSystem(s); s != nil {
				id = s.ID
			} else {
				eval.ReportError("no software system named %q", s)
				return
			}
		case 2:
			if s := expr.Root.Model.SoftwareSystem(elems[0]); s != nil {
				if c := s.Container(elems[1]); c != nil {
					id = c.ID
				} else {
					eval.ReportError("no container named %q in software system %q", elems[1], elems[0])
					return
				}
			} else {
				eval.ReportError("no software system named %q", elems[0])
			}
		default:
			eval.ReportError("too many slashes in path (%d)", len(elems)-1)
			return
		}
	default:
		eval.ReportError("DynamicView: invalid scope, expected software system, container, software system name or container path, got %T", scope)
		return
	}
	description, dsl, err := parseView(args...)
	if err != nil {
		eval.ReportError("DynamicView: " + err.Error())
		return
	}
	v := &expr.DynamicView{
		ViewProps: &expr.ViewProps{
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
// DeploymentView accepts 4 to 5 arguments: the first argument is the scope:
// either the keyword 'Global', a software system or the name of a software
// system. The second argument is the name of the environment. The third
// argument is a unique key for the view. The fourth argument is an optional
// description. The last argument is a function describing the properties of the
// view.
//
// Usage:
//
//    DeploymentView(Scope, "<environment>", "<key>", func())
//
//    DeploymentView(Scope, "<environment>", "<key>", "[description]", func())
//
// Where Scope is 'Global', a software system or its name.
//
// Example:
//
//     var _ = Design(func() {
//         System("System", func() {
//              Container("Container")
//              Container("OtherContainer")
//         })
//         DeploymentEnvironment("Production", func() {
//             DeploymentNode("Cloud", func() {
//                 ContainerInstance("System/Container")
//                 ContainerInstance("System/OtherContainer")
//             })
//         })
//         Views(func() {
//             DeploymentView(Global, "Production", "deployment", "A deployment overview diagram.", func() {
//                 Title("Overview of deployment")
//                 AddAll()
//                 Remove("System/OtherContainer")
//                 AutoLayout()
//                 AnimationStep("System/Container")
//                 PaperSize(SizeSlide4X3)
//             })
//         })
//     })
//
func DeploymentView(scope interface{}, env, key string, args ...interface{}) {
	vs, ok := eval.Current().(*expr.Views)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	missing := true
	expr.Iterate(func(e interface{}) {
		if dn, ok := e.(*expr.DeploymentNode); ok {
			if dn.Environment == env {
				missing = false
			}
		}
	})
	if missing {
		eval.ReportError("DeploymentView: environment %q not defined", env)
		return
	}
	description, dsl, err := parseView(args...)
	if err != nil {
		eval.ReportError("DeploymentView: " + err.Error())
		return
	}
	var id string
	switch s := scope.(type) {
	case int:
		id = "" // Global scope
	case *expr.SoftwareSystem:
		id = s.ID
	case string:
		if se := expr.Root.Model.SoftwareSystem(s); se == nil {
			eval.InvalidArgError("'Global', a software system or a software system name", scope)
		} else {
			id = se.ID
		}
	default:
		eval.InvalidArgError("'Global', a software system or a software system name", scope)
		return
	}
	v := &expr.DeploymentView{
		ViewProps: &expr.ViewProps{
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
	if v, ok := eval.Current().(expr.View); ok {
		v.Props().Title = t
	} else {
		eval.IncompatibleDSL()
	}
}

// Add adds a person or an element to a view.
//
// Add must appear in SystemLandscapeView, SystemContextView, ContainerView,
// ComponentView or DeploymentView.
//
// Usage depends on the view Add is used in. In all cases Add supports an
// optional DSL function as last argument that can be used to specify rendering
// details.
//
//    - In SystemLandscapeView, SystemContextView, ContainerView and ComponentView
//      Add accepts a person, a software system or their names.
//
//    - In ContainerView Add also accepts a container or the path to a container.
//      The path to a container is either its name if it is a child of the the
//      software system the container view is for or the name of the parent software
//      system followed by a slash and the name of the container otherwise.
//
//    - In ComponentView Add also accepts a component or the path to a component.
//      The path of a component is either its name if it is a child of the container
//      the component view is for or the name of the parent software system followed
//      by a slash, the name of the parent container, another slash and the name of
//      the component.
//
//    - In DeploymentView, Add accepts a deployment node, a container instance, an
//      infrastructure node or their paths. The path is constructed by appending the
//      top level deployment node name with the child deployment name recursively
//      followed by the name of the inner deployment node, infrastructure node or
//      container instance. The names must be separated by slashes in the path. For
//      container instances the path may end with the container instance ID.
//
// Usage (SystemLandscapeView, SystemContextView, ContainerView and ComponentView):
//
//     Add(Person|"<Person>"[, func()])
//
//     Add(SoftwareSystem|"<Software System>"[, func()])
//
// Additionally for ContainerView:
//
//     Add(Container[, func()])
//
//     Add("<Container>"[, func()]) // If container is a child of the system software
//                                  // the container view is for.
//     Add("<Software System/Container>"[, func()])
//
// Additionally for ComponentView:
//
//     Add(Component[, func()])
//
//     Add("<Component>"[, func()]) // If component is a child of the container the
//                                  // component view is for.
//     Add("<Container/Component>"[, func()]) // if container is a child of the
//                                            // software system that contains the
//                                            // container the component view is for.
//     Add("<Software System/Container/Component>"[, func()])
//
// Usage (DeploymentView):
//
//     Add(DeploymentNode[, func()])
//
//     Add(InfrastructureNode[, func()])
//
//     Add(ContainerInstance[, func()])
//
//     Add("<Deployment Node>"[, func()]) // top level deployment node
//
//     Add("<Parent Deployment Node>/.../<Child Deployment Node>"[, func()]) // child deployment node
//
//     Add("<Parent Deployment Node>/.../<Infrastructure Node>"[, func()])
//
//     Add("<Parent Deployment Node>/.../<Container Instance>:<Container Instance ID>"[, func()])
//
// Where "<Parent Deployment Node>/..." describes a deployment node hierarchy
// starting with the top level deployment node name followed by its child
// deployment node name etc.
//
// Example:
//
//     var _ = Design(func() {
//         var Customer = Person("Customer", "A customer", func() {
//             Uses("Software System", "Sends emails", "SMTP")
//         })
//         var MyContainer *expr.Container
//         var System = SoftwareSystem("Software System", "My software system.", func() {
//             MyContainer = Container("Container", "A container")
//         })
//         var OtherSystem = SoftwareSystem("Other System", "My other software system.", func() {
//             Container("Other Container", "Another container")
//         })
//         var Kubernetes *expr.DeploymentNode
//         DeploymentEnvironment("Production", func() {
//             DeploymentNode("Cloud", func() {
//                 Kubernetes = DeploymentNode("Kubernetes", func() {
//                     ContainerInstance(MyContainer) // Same as ContainerInstance("Software System/Container")
//                     ContainerInstance(MyContainer, func() {
//                         InstanceID(2)
//                     })
//                 })
//                 InfrastructureNode("API Gateway")
//             })
//         })
//         Views(func() {
//             SystemContextView(SoftwareSystem, "context", "An overview diagram.", func() {
//                 Add(Customer, func() { // Same as Add("Customer", func() {
//                     Coord(10, 10)
//                     NoRelationships()
//                 })
//             })
//             ContainerView(SoftwareSystem, "container", "A diagram of Software System", func() {
//                 Add(Customer)
//                 Add(MyContainer) // Same as Add("Container")
//                 Add("Other System/Other Container")
//             })
//             DeploymentView(Global, "Production", "deployment", "A deployment overview diagram.", func() {
//                 Add("Cloud/Kubernetes/Container")
//                 Add("Cloud/Kubernetes/Container:2")
//                 Add("Cloud:API Gateway")
//             })
//         })
//     })
//
func Add(element interface{}, dsl ...func()) {
	var (
		eh   expr.ElementHolder
		err  error
		view expr.View
	)
	switch v := eval.Current().(type) {
	case *expr.LandscapeView, *expr.ContextView, *expr.ContainerView, *expr.ComponentView:
		view = v.(expr.View)
		eh, err = findViewElement(view, element)
	case *expr.DeploymentView:
		view = v
		eh, err = findViewElement(v, element)
	case *expr.DynamicView:
		eval.ReportError("Add: only relationships may be added explicitly to dynamic views using Link")
		return
	default:
		eval.IncompatibleDSL()
		return
	}
	if err != nil {
		eval.ReportError("Add: " + err.Error())
		return
	}
	if err := view.(expr.ViewAdder).AddElements(eh); err != nil {
		eval.ReportError("Add: " + err.Error()) // Element type not supported in view
		return
	}
	if len(dsl) > 0 {
		eval.Execute(dsl[0], view.Props().ElementView(eh.GetElement().ID))
		if len(dsl) > 1 {
			eval.ReportError("Add: too many arguments")
		}
	}
}

// Link adds a relationship to a view.
//
// Link must appear in SystemLandscapeView, SystemContextView, ContainerView,
// ComponentView, DynamicView or DeploymentView.
//
// Link takes the relationship as defined by its source, destination and when
// needed to distinguish its description as first arguments and an optional
// function as last argument.
//
// The source and destination are identified by reference or by path. The path
// consists of the element name if a top level element (person or software system)
// or if the element is in scope (container in container view software system,
// container in component view software system or component in component view
// container). When the element is not in scope the path specifies the parent
// element name followed by a slash and the element name. If the parent itself
// is not in scope (i.e. a component that is a child of a different software
// system in a ComponentView) then the path specifies the top-level software
// system followed by a slash, the container name, another slash and the
// component name.
//
// Usage:
//
//      Link(Source, Destination) // If only one relationship exists between Source
//                                // and Destination.
//      Link(Source, Destination, func()) // If only one relationship exists between
//                                        // Source and Destination.
//
//      Link(Source, Destination, Description)
//
//      Link(Source, Destination, Description, func())
//
//
// Where Source and Destination are one of:
//
//    - Person|"<Person>"
//    - SoftwareSystem|"<Software System>"
//    - Container
//    - "<Container>" (if container is in the container or component view software system)
//    - "<Software System/Container>"
//    - Component
//    - "<Component>" (if component is in the component view container)
//    - "<Container>/<Component>" (if container is in the component view software system)
//    - "<Software System>/<Container>/<Component>"
//
// Example:
//
//     var _ = Design(func() {
//         var System = SoftwareSystem("Software System", "My software system.", func() {
//             Container("Container")
//             Container("Other Container", func() {
//                 Uses("Container", "Makes requests to")
//             })
//         })
//         var Person = Person("Customer", func() {
//             External()
//             Uses(System, "Sends emails", "SMTP")
//         })
//         DeploymentEnvironment("Production", func() {
//             DeploymentNode("Cloud", func() {
//                 ContainerInstance("Container")
//                 ContainerInstance("Container", func() {
//                     InstanceID(2)
//                 })
//                 ContainerInstance("Other Container")
//             })
//         })
//         Views(func() {
//             SystemContextView(SoftwareSystem, "context", "An overview diagram", func() {
//                 Add(System, func() {
//                     Coord(10, 10)
//                     NoRelationships()
//                 })
//                 Link(Person, System, "Sends emails", func() {
//                     Vertices(10, 20, 10, 40)
//                     Routing(RoutingOrthogonal)
//                     Position(45)
//                 })
//             })
//             DeploymentView(Global, "Production", "deployment", "A deployment view", func() {
//                 Add("Cloud/Container")
//                 Add("Cloud/Container/2")
//                 Add("Cloud/Other Container")
//                 Link("Cloud/Container", "Cloud/Other Container")
//                 Link("Cloud/Container/2", "Cloud/Other Container")
//             })
//             DynamicView(SoftwareSystem, "dynamic", func() {
//                 Title("Customer flow")
//                 Link(Person, System, "Sends emails", func() {
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
func Link(source, destination interface{}, args ...interface{}) {
	v, ok := eval.Current().(expr.View)
	if !ok {
		eval.IncompatibleDSL()
	}
	src, dest, desc, dsl, err := parseLinkArgs(v, source, destination, args)
	if err != nil {
		eval.ReportError("Link: " + err.Error())
		return
	}
	rel := &expr.RelationshipView{
		Source:      src.GetElement(),
		Destination: dest.GetElement(),
		Description: desc,
	}
	if dsl != nil {
		eval.Execute(dsl, rel)
	}
	v.Props().RelationshipViews = append(v.Props().RelationshipViews, rel)
}

// AddAll includes all elements and relationships in the view scope.
//
// AddAll may appear in SystemLandscapeView, SystemContextView, ContainerView,
// ComponentView or DeploymentView.
//
// AddAll takes no argument.
//
// Example:
//
//     var _ = Design(func() {
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
	switch v := eval.Current().(type) {
	case *expr.DynamicView:
		eval.IncompatibleDSL()
	case expr.View:
		v.Props().AddAll = true
	default:
		eval.IncompatibleDSL()
	}
}

// AddNeighbors Adds all of the permitted elements which are directly connected
// to the specified element. Permitted elements are software systems and people
// for system landscape and system context views, software systems, people and
// containers for container views and software system, people, containers and
// components for component views.
//
// AddNeighbors must appear in SystemLandscapeView, SystemContextView,
// ContainerView, ComponentView or DeploymentView.
//
// AddNeighbors accept a single argument which is the element that should be
// added with its direct relationships. The element is identified by reference
// or by path. The path consists of the element name if a top level element
// (person or software system) or if the element is in scope (container in
// container view software system, container in component view software system
// or component in component view container). When the element is not in scope
// the path specifies the parent element name followed by a slash and the
// element name. If the parent itself is not in scope (i.e. a component that is
// a child of a different software system in a ComponentView) then the path
// specifies the top-level software system followed by a slash, the container
// name, another slash and the component name.
//
// Example:
//
//     var _ = Design(func() {
//         var System = SoftwareSystem("Software System", "My software system.", func() {
//             Container("Container", func() {
//                 Component("Component")
//             })
//         })
//         SoftwareSystem("Other System", func() {
//             Container("Container", func() {
//                 Component("Component")
//             })
//         })
//         var Customer = Person("Customer", func() {
//             External()
//             Uses(System, "Sends emails", "SMTP")
//         })
//         Views(func() {
//             SystemContextView(SoftwareSystem, "context", "An overview diagram.", func() {
//                 AddNeighbors(System)
//                 AddNeighbors(Customer)
//             })
//             ComponentView("Software System/Container", func() {
//                  Add("Component")
//                  AddNeighbors("Other System/Container/Component")
//             })
//         })
//     })
//
func AddNeighbors(element interface{}) {
	v, ok := eval.Current().(expr.View)
	if !ok {
		eval.IncompatibleDSL()
	}
	eh, err := findViewElement(v, element)
	if err != nil {
		eval.ReportError("AddNeighbors: " + err.Error())
		return
	}
	v.Props().AddNeighbors = append(v.Props().AddNeighbors, eh.GetElement())
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
//    - Deployment view: adds all deployment nodes.
//
// AddDefault must appear in SystemLandscapeView, SystemContextView,
// ContainerView or ComponentView.
//
// AddDefault takes no argument.
//
// Example:
//
//     var _ = Design(func() {
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
	case expr.View:
		v.Props().AddDefault = true
	default:
		eval.IncompatibleDSL()
	}
}

// AddContainers includes all containers in scope to the view.
//
// AddContainers may appear in ContainerView or ComponentView.
//
// AddContainers takes no argument.
func AddContainers() {
	switch v := eval.Current().(type) {
	case *expr.ContainerView:
		v.AddElements(expr.Registry[v.SoftwareSystemID].(*expr.SoftwareSystem).Containers.Elements()...)
	case *expr.ComponentView:
		c := expr.Registry[v.ContainerID].(*expr.Container)
		v.AddElements(c.System.Containers.Elements()...)
	default:
		eval.IncompatibleDSL()
	}
}

// AddInfluencers adds all containers of the ContainerView as well as all
// external influencers, that is all persons and all other software systems with
// incoming or outgoing dependencies. Additionally, all relationships of
// external dependencies are omitted to keep the diagram clean.
//
// AddInfluencers must appear in ContainerView.
//
// AddInfluencers takes no argument.
func AddInfluencers() {
	cv, ok := eval.Current().(*expr.ContainerView)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	cv.AddInfluencers = true
}

// AddComponents includes all components in scope to the view.
//
// AddComponents must appear in ComponentView.
//
// AddComponents takes no argument
func AddComponents() {
	if cv, ok := eval.Current().(*expr.ComponentView); ok {
		cv.AddElements(expr.Registry[cv.ContainerID].(*expr.Container).Components.Elements()...)
		return
	}
	eval.IncompatibleDSL()
}

// Remove given person, software system, container or component from the view.
//
// Remove must appear in SystemLandscapeView, SystemContextView, ContainerView
// or ComponentView.
//
// Remove takes one argument: the element or the path to the element to be
// removed. The path consists of the element name if a top level element (person
// or software system) or if the element is in scope (container in container
// view software system, container in component view software system or
// component in component view container). When the element is not in scope the
// path specifies the parent element name followed by a slash and the element
// name. If the parent itself is not in scope (i.e. a component that is a child
// of a different software system in a ComponentView) then the path specifies
// the top-level software system followed by a slash, the container name,
// another slash and the component name.
//
// Usage:
//
//     Remove(Person|"<Person>")
//
//     Remove(SoftwareSystem|"<Software System>")
//
//     Remove(Container)
//
//     Remove("<Container>") // if container is in scope
//
//     Remove("<Software System>/<Container>")
//
//     Remove(Component)
//
//     Remove("<Component>") // if component is in scope
//
//     Remove("<Container>/<Component>") // if container is in scope
//
//     Remove("<Software System>/<Container>/<Component>")
//
// Example:
//
//     var _ = Design(func() {
//         var System = SoftwareSystem("Software System", "My software system.", func() {
//             Container("Unwanted")
//         })
//         var Customer = Person("Customer", func() {
//             External()
//             Uses(System, "Sends emails", "SMTP")
//         })
//         Views(func() {
//             SystemContextView(System, "context", "An overview diagram.", func() {
//                 AddDefault()
//                 Remove(Customer)
//                 Remove("Unwanted")
//             })
//         })
//     })
//
func Remove(element interface{}) {
	v, ok := eval.Current().(expr.View)
	if !ok {
		eval.IncompatibleDSL()
	}
	eh, err := findViewElement(v, element)
	if err != nil {
		eval.ReportError("Remove: " + err.Error())
		return
	}
	v.Props().RemoveElements = append(v.Props().RemoveElements, eh.GetElement())
}

// RemoveTagged removes all elements and relationships with the given tag from
// the view.
//
// RemoveTagged must appear in SystemLandscapeView, SystemContextView,
// ContainerView or ComponentView.
//
// Remove takes one argument: the tag identifying the elements and relationships
// to be removed.
//
// Usage:
//
//     RemoveTagged("<tag>")
//
// Example:
//
//     var _ = Design(func() {
//         var System = SoftwareSystem("Software System", "My software system.", func() {
//             Container(System, "Unwanted", func() {
//                 Tag("irrelevant")
//             })
//         })
//         var Customer = Person("Customer", func() {
//             External()
//             Uses(System, "Sends emails", "SMTP", func() {
//                 Tag("irrelevant")
//             })
//         })
//         Views(func() {
//             SystemContextView(SoftwareSystem, "context", "An overview diagram.", func() {
//                 AddDefault()
//                 RemoveTagged("irrelevant")
//             })
//         })
//     })
//
func RemoveTagged(tag string) {
	v, ok := eval.Current().(expr.View)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	v.Props().RemoveTags = append(v.Props().RemoveTags, tag)
}

// Unlink removes a relationship from a view.
//
// Unlink must appear in SystemLandscapeView, SystemContextView, ContainerView
// or ComponentView.
//
// Unlink takes the relationship as defined by its source, destination and when
// needed to distinguish its description.
//
// The source and destination are identified by reference or by path. The path
// consists of the element name if a top level element (person or software
// system) or if the element is in scope (container in container view software
// system, container in component view software system or component in component
// view container). When the element is not in scope the path specifies the
// parent element name followed by a slash and the element name. If the parent
// itself is not in scope (i.e. a component that is a child of a different
// software system in a ComponentView) then the path specifies the top-level
// software system followed by a slash, the container name, another slash and
// the component name.
//
// Usage:
//
//      Unlink(Source, Destination) // If only one relationship exists between
//                                  // Source and Destination
//      Unlink(Source, Destination, Description)
//
// Where Source and Destination are one of:
//    - Person|"<Person>"
//    - SoftwareSystem|"<Software System>"
//    - Container
//    - "<Container>" (if container is in the container or component view software system)
//    - "<Software System/Container>"
//    - Component
//    - "<Component>" (if component is in the component view container)
//    - "<Container>/<Component>" (if container is in the component view software system)
//    - "<Software System>/<Container>/<Component>"
//
// Example:
//
//     var _ = Design(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         var Person = Person("Customer", func() {
//             External()
//             Uses(System, "Sends emails", "SMTP")
//         })
//         Views(func() {
//             SystemContextView(SoftwareSystem, "context", "An overview diagram.", func() {
//                 AddDefault()
//                 Unlink(Person, System, "Sends emails")
//             })
//         })
//     })
//
func Unlink(source, destination interface{}, description ...string) {
	v, ok := eval.Current().(expr.View)
	if !ok {
		eval.IncompatibleDSL()
	}
	var args []interface{}
	if len(description) > 0 {
		args = []interface{}{description[0]}
		if len(description) > 1 {
			eval.ReportError("Unlink: too many arguments")
		}
	}
	src, dest, desc, _, err := parseLinkArgs(v, source, destination, args)
	if err != nil {
		eval.ReportError("Unlink: " + err.Error())
		return
	}
	v.Props().RemoveRelationships = append(v.Props().RemoveRelationships,
		&expr.Relationship{
			Source:      src.GetElement(),
			Destination: dest.GetElement(),
			Description: desc,
		})
}

// RemoveUnreachable removes all elements and people that cannot be reached by
// traversing the graph of relationships starting with the given element or
// person.
//
// RemoveUnreachable must appear in SystemLandscapeView, SystemContextView,
// ContainerView or ComponentView.
//
// RemoveUnreachable accept a single argument which is the element used to start
// the graph traversal. The element is identified by reference or by path. The
// path consists of the element name if a top level element (person or software
// system) or if the element is in scope (container in container view software
// system, container in component view software system or component in component
// view container). When the element is not in scope the path specifies the
// parent element name followed by a slash and the element name. If the parent
// itself is not in scope (i.e. a component that is a child of a different
// software system in a ComponentView) then the path specifies the top-level
// software system followed by a slash, the container name, another slash and
// the component name.
//
// Example:
//
//     var _ = Design(func() {
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
func RemoveUnreachable(element interface{}) {
	v, ok := eval.Current().(expr.View)
	if !ok {
		eval.IncompatibleDSL()
	}
	eh, err := findViewElement(v, element)
	if err != nil {
		eval.ReportError("RemoveUnreachable: " + err.Error())
		return
	}
	v.Props().RemoveUnreachable = append(v.Props().RemoveUnreachable, eh.GetElement())
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
//     var _ = Design(func() {
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
	if v, ok := eval.Current().(expr.View); ok {
		v.Props().RemoveUnrelated = true
		return
	}
	eval.IncompatibleDSL()
}

// AutoLayout enables automatic layout mode for the diagram. The
// first argument indicates the rank direction, it must be one of
// RankTopBottom, RankBottomTop, RankLeftRight or RankRightLeft
//
// AutoLayout must appear in SystemLandscapeView, SystemContextView,
// ContainerView, ComponentView, DynamicView or DeploymentView.
//
// AutoLayout accepts one or two arguments: the layout rank direction and
// an optional function DSL that describes the layout properties.
//
// Example:
//
//     var _ = Design(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         var OtherSystem = SoftwareSystem("Other software System")
//         var Customer = Person("Customer", func() {
//             External()
//             Uses(System, "Sends emails", "SMTP")
//         })
//         Views(func() {
//             SystemContextView(SoftwareSystem, "context", "An overview diagram.", func() {
//                 AddDefault()
//                 AutoLayout(RankTopBottom, func() {
//                     RankSeparation(200)
//                     NodeSeparation(100)
//                     EdgeSeparation(10)
//                     Vertices()
//                 })
//             })
//         })
//     })
//
func AutoLayout(rank RankDirectionKind, args ...func()) {
	v, ok := eval.Current().(expr.View)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	var dsl func()
	if len(args) > 0 {
		dsl = args[0]
		if len(args) > 1 {
			eval.ReportError("AutoLayout: too many arguments")
		}
	}
	r, n, e := 300, 600, 200
	layout := &expr.AutoLayout{
		RankDirection: expr.RankDirectionKind(rank),
		RankSep:       &r,
		NodeSep:       &n,
		EdgeSep:       &e,
	}
	if dsl != nil {
		eval.Execute(dsl, layout)
	}
	v.Props().AutoLayout = layout
}

// AnimationStep defines an animation step consisting of the specified elements.
//
// AnimationStep must appear in SystemLandscapeView, SystemContextView,
// ContainerView, ComponentView or DeploymentView.
//
// AnimationStep accepts the list of elements that should be rendered in the
// animation step as argument. Each element is identified by reference or by
// path. The path consists of the element name if a top level element (person or
// software system) or if the element is in scope (container in container view
// software system, container in component view software system or component in
// component view container). When the element is not in scope the path
// specifies the parent element name followed by a slash and the element name.
// If the parent itself is not in scope (i.e. a component that is a child of a
// different software system in a ComponentView) then the path specifies the
// top-level software system followed by a slash, the container name, another
// slash and the component name.
//
// Example
//
//     var _ = Design(func() {
//         SoftwareSystem("Software System", "My software system.")
//         var OtherSystem = SoftwareSystem("Other software System")
//         var Customer = Person("Customer", func() {
//             External()
//             Uses("Software System", "Sends emails", "SMTP")
//         })
//         Views(func() {
//             SystemContextView(SoftwareSystem, "context", "An overview diagram.", func() {
//                 AddDefault()
//                 AnimationStep(OtherSystem, Customer)
//                 AnimationStep("Software System")
//             })
//         })
//     })
//
func AnimationStep(elements ...interface{}) {
	v, ok := eval.Current().(expr.ViewAdder)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	step := &expr.AnimationStep{View: v.(expr.View)}
	for _, elem := range elements {
		eh, err := findViewElement(v.(expr.View), elem)
		if err != nil {
			eval.ReportError("AnimationStep: " + err.Error())
			continue
		}
		step.Add(eh)
	}
	if err := v.AddAnimationStep(step); err != nil {
		eval.ReportError("AnimationStep: " + err.Error())
	}
}

// PaperSize defines the paper size that should be used to render
// the view in the Structurizr service.
//
// PaperSize must appear in SystemLandscapeView, SystemContextView,
// ContainerView, ComponentView, DynamicView or DeploymentView.
//
// PaperSize accepts a single argument: the paper size. The possible values for
// the argument follow the patterns SizeA[0-6][Portrait|Landscape],
// SizeLetter[Portrait|Landscape] or SizeLegal[Portrait_Landscape].
// Alternatively the argument may be one of SizeSlide4X3, SizeSlide16X9 or
// SizeSlide16X10.
//
// Example
//
//     var _ = Design(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         var OtherSystem = SoftwareSystem("Other software System")
//         var Customer = Person("Customer", func() {
//             External()
//             Uses(System, "Sends emails", "SMTP")
//         })
//         Views(func() {
//             SystemContextView(SoftwareSystem, "context", "An overview diagram.", func() {
//                 AddDefault()
//                 PaperSize(SizeSlide4X3)
//             })
//         })
//     })
//
func PaperSize(size PaperSizeKind) {
	v, ok := eval.Current().(expr.View)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	v.Props().PaperSize = expr.PaperSizeKind(size)
}

// EnterpriseBoundaryVisible makes the enterprise boundary visible to differentiate internal
// elements from external elements on the resulting diagram.
//
// EnterpriseBoundaryVisible must appear in SystemLandscapeView or SystemContextView.
//
// EnterpriseBoundaryVisible takes no argument
func EnterpriseBoundaryVisible() {
	t := true
	switch v := eval.Current().(type) {
	case *expr.LandscapeView:
		v.EnterpriseBoundaryVisible = &t
	case *expr.ContextView:
		v.EnterpriseBoundaryVisible = &t
	default:
		eval.IncompatibleDSL()
	}
}

// SystemBoundariesVisible makes the system boundaries visible for "external" containers
// (those outside the software system in scope)
//
// SystemBoundariesVisible must appear in ContainerView.
//
// SystemBoundariesVisible takes no argument
func SystemBoundariesVisible() {
	if v, ok := eval.Current().(*expr.ContainerView); ok {
		t := true
		v.SystemBoundariesVisible = &t
		return
	}
	eval.IncompatibleDSL()
}

// ContainerBoundariesVisible makes the enterprise boundary visible to differentiate internal
// elements from external elements on the resulting diagram.
//
// ContainerBoundariesVisible must appear in ComponentView.
//
// ContainerBoundariesVisible takes no argument
func ContainerBoundariesVisible() {
	if v, ok := eval.Current().(*expr.ComponentView); ok {
		t := true
		v.ContainerBoundariesVisible = &t
		return
	}
	eval.IncompatibleDSL()
}

// Coord defines explicit coordinates for where to render a person or element.
//
// Coord must appear in Add.
//
// Coord takes two arguments: the X and Y where the person or element is rendered.
//
// Example:
//
//     var _ = Design(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         var Customer = Person("Customer", func() {
//             External()
//             Uses(System, "Sends emails", "SMTP")
//         })
//         Views(func() {
//             SystemContextView(System, "context", "An overview diagram.", func() {
//                 Add(Customer, func() {
//                     Coord(200,200)
//                 })
//             })
//         })
//     })
//
func Coord(x, y int) {
	if ev, ok := eval.Current().(*expr.ElementView); ok {
		ev.X = &x
		ev.Y = &y
		return
	}
	eval.IncompatibleDSL()
}

// NoRelationship indicates that no relationship should be rendered to and from the person or element.
//
// NoRelationship must appear in Add.
//
// NoRelationship takes no argument.
//
// Example:
//
//     var _ = Design(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         var OtherSystem = SoftwareSystem("Other software System")
//         var Customer = Person("Customer", func() {
//             External()
//             Uses(System, "Sends emails", "SMTP")
//         })
//         Views(func() {
//             SystemContextView(SoftwareSystem, "context", "An overview diagram.", func() {
//                 Add(Customer, func() {
//                     NoRelationship()
//                 })
//             })
//         })
//     })
//
func NoRelationship() {
	if ev, ok := eval.Current().(*expr.ElementView); ok {
		ev.NoRelationship = true
		return
	}
	eval.IncompatibleDSL()
}

// Vertices lists the x and y coordinate of the vertices used to render the
// relationship.
//
// Vertices must appear in Add when adding relationships.
//
// Vertices takes the x and y coordinates of the vertices as argument. The
// number of arguments must be even.
//
// Example:
//
//     var _ = Design(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         var Customer = Person("Customer", func() {
//             External()
//             Uses(System, "Sends emails", "SMTP")
//         })
//         Views(func() {
//             SystemContextView(SoftwareSystem, "context", "An overview diagram.", func() {
//                 Add(Customer, System, func() {
//                     Vertices(300, 100, 400, 200)
//                 })
//             })
//         })
//     })
//
func Vertices(args ...int) {
	if len(args)%2 != 0 {
		eval.ReportError("Vertices: must be given an even number of arguments")
	}
	rv, ok := eval.Current().(*expr.RelationshipView)
	if !ok {
		eval.IncompatibleDSL()
	}
	for i := 0; i < len(args); i += 2 {
		rv.Vertices = append(rv.Vertices, &expr.Vertex{args[i], args[i+1]})
	}
}

// Routing algorithm used when rendering relationship, defaults to
// RoutingDirect.
//
// Routing must appear in a Add expr.ssion that adds a relationship or in a
// RelationshipStyle expr.ssion.
//
// Routing takes one argument: one of RoutingDirect, RoutingCurved or
// RoutingOrthogonal.
//
// Example:
//
//     var _ = Design(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         var Customer = Person("Customer", func() {
//             External()
//             Uses(System, "Sends emails", "SMTP")
//         })
//         Views(func() {
//             SystemContextView(SoftwareSystem, "context", "An overview diagram.", func() {
//                 Add(Customer, System, func() {
//                     Routing(RoutingDirect)
//                 })
//             })
//         })
//     })
//
func Routing(kind RoutingKind) {
	switch a := eval.Current().(type) {
	case *expr.RelationshipView:
		a.Routing = expr.RoutingKind(kind)
	case *expr.RelationshipStyle:
		a.Routing = expr.RoutingKind(kind)
	default:
		eval.IncompatibleDSL()
	}
}

// Position sets the position of a relationship annotation along the line.
//
// Position must appear in a Add expression that adds a relationship or in
// RelationshipStyle.
//
// Position takes one argument: the position value between 0 (start of line) and
// 100 (end of line).
//
// Example:
//
//     var _ = Design(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         var Customer = Person("Customer", func() {
//             External()
//             Uses(System, "Sends emails", "SMTP")
//         })
//         Views(func() {
//             SystemContextView(SoftwareSystem, "context", "An overview diagram.", func() {
//                 Add(Customer, System, func() {
//                     Position(40)
//                 })
//             })
//         })
//     })
//
func Position(pos int) {
	if pos < 0 || pos > 100 {
		eval.InvalidArgError("integer between 0 and 100", pos)
		return
	}
	switch a := eval.Current().(type) {
	case *expr.RelationshipView:
		a.Position = &pos
	case *expr.RelationshipStyle:
		a.Position = &pos
	default:
		eval.IncompatibleDSL()
	}
}

// RankSeparation sets the separation between ranks in pixels, defaults to 300.
//
// RankSeparation must appear in AutoLayout.
//
// RankSeparation takes one argument: the rank separation in pixels.
//
// Example:
//
//     var _ = Design(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         Views(func() {
//             SystemContextView(SoftwareSystem, "context", "An overview diagram.", func() {
//                 AutoLayout(func() {
//                     RankSeparation(500)
//                 })
//             })
//         })
//     })
//
func RankSeparation(sep int) {
	if sep < 0 {
		eval.ReportError("RankSeparation: value must be positive")
		return
	}
	if a, ok := eval.Current().(*expr.AutoLayout); ok {
		a.RankSep = &sep
		return
	}
	eval.IncompatibleDSL()
}

// NodeSeparation sets the separation between nodes in pixels, defaults to 600.
//
// NodeSeparation must appear in AutoLayout.
//
// NodeSeparation takes one argument: the node separation in pixels.
//
// Example:
//
//     var _ = Design(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         Views(func() {
//             SystemContextView(SoftwareSystem, "context", "An overview diagram.", func() {
//                 AutoLayout(func() {
//                     NodeSeparation(500)
//                 })
//             })
//         })
//     })
//
func NodeSeparation(sep int) {
	if sep < 0 {
		eval.ReportError("NodeSeparation: value must be positive")
		return
	}
	if a, ok := eval.Current().(*expr.AutoLayout); ok {
		a.NodeSep = &sep
		return
	}
	eval.IncompatibleDSL()
}

// EdgeSeparation sets the separation between edges in pixels, defaults to 200.
//
// EdgeSeparation must appear in AutoLayout.
//
// EdgeSeparation takes one argument: the edge separation in pixels.
//
// Example:
//
//     var _ = Design(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         Views(func() {
//             SystemContextView(SoftwareSystem, "context", "An overview diagram.", func() {
//                 AutoLayout(func() {
//                     EdgeSeparation(500)
//                 })
//             })
//         })
//     })
//
func EdgeSeparation(sep int) {
	if sep < 0 {
		eval.ReportError("EdgeSeparation: value must be positive")
		return
	}
	if a, ok := eval.Current().(*expr.AutoLayout); ok {
		a.EdgeSep = &sep
		return
	}
	eval.IncompatibleDSL()
}

// RenderVertices indicates that vertices should be created during automatic
// layout, false by default.
//
// RenderVertices must appear in AutoLayout.
//
// RenderVertices takes no argument.
//
// Example:
//
//     var _ = Design(func() {
//         var System = SoftwareSystem("Software System", "My software system.")
//         Views(func() {
//             SystemContextView(SoftwareSystem, "context", "An overview diagram.", func() {
//                 AutoLayout(func() {
//                     RenderVertices()
//                 })
//             })
//         })
//     })
//
func RenderVertices() {
	if a, ok := eval.Current().(*expr.AutoLayout); ok {
		t := true
		a.Vertices = &t
		return
	}
	eval.IncompatibleDSL()
}

// parseView is a helper function that parses the given view DSL
// arguments. Accepted syntax are:
//
//     func()
//     "[description]", func()
//
func parseView(args ...interface{}) (description string, dsl func(), err error) {
	if len(args) == 0 {
		err = fmt.Errorf("missing argument")
		return
	}
	switch a := args[0].(type) {
	case string:
		description = a
	case func():
		dsl = a
	default:
		err = fmt.Errorf("expected string or function, got %T", args[0])
		return
	}
	if len(args) > 1 {
		if dsl != nil {
			err = fmt.Errorf("DSL function must be last argument")
			return
		}
		switch a := args[1].(type) {
		case func():
			dsl = a
		default:
			err = fmt.Errorf("expected function, got %T", args[1])
			return
		}
		if len(args) > 2 {
			err = fmt.Errorf("too many arguments")
		}
	}
	return
}

// findViewElement returns the element identifed by element that
// is in scope for the given view. See model.FindElement for details.
func findViewElement(view expr.View, element interface{}) (expr.ElementHolder, error) {
	if eh, ok := element.(expr.ElementHolder); ok {
		return eh, nil
	}
	name, ok := element.(string)
	if !ok {
		return nil, fmt.Errorf("expected element or element name, got %T", element)
	}
	switch v := view.(type) {
	case *expr.LandscapeView, *expr.ContextView:
		return expr.Root.Model.FindElement(nil, name)
	case *expr.ContainerView:
		scope := expr.Registry[v.SoftwareSystemID].(expr.ElementHolder)
		return expr.Root.Model.FindElement(scope, name)
	case *expr.ComponentView:
		scope := expr.Registry[v.ContainerID].(expr.ElementHolder)
		res, err := expr.Root.Model.FindElement(scope, name)
		return res, err
	case *expr.DeploymentView:
		return findDeploymentViewElement(name)
	case *expr.DynamicView:
		var scope expr.ElementHolder
		if v.ElementID != "" {
			scope, _ = expr.Registry[v.ElementID].(expr.ElementHolder)
		}
		return expr.Root.Model.FindElement(scope, name)
	default:
		panic("view does not support adding elements") // bug
	}
}

// findDeploymentViewElement returns the element identified as follows:
//
//    - DeploymentNode: returns the given deployment node.
//    - InfrastructureNode: returns the given infrastructure node.
//    - ContainerInstance: returns the given container instance.
//    - "DeploymentNode/.../Child DeploymentNode": returns the deployment node with
//      the given path (top level deployemnt node name to child deployment node name
//      all separated with slashes).
//    - "DeploymentNode/.../InfrastructureNode": returns the infrastructure node
//      with the given name in the given deployment node path.
//    - "DeploymentNode/.../Container:InstanceID": returns the container instance
//      in the given deployment node path and with the given container name and
//      instance ID.
//
func findDeploymentViewElement(e interface{}, cid ...int) (expr.ElementHolder, error) {
	switch s := e.(type) {
	case *expr.DeploymentNode, *expr.InfrastructureNode, *expr.ContainerInstance:
		return s.(expr.ElementHolder), nil
	case string:
		elems := strings.Split(s, "/")
		parent := expr.Root.Model.DeploymentNode(elems[0])
		if parent == nil {
			return nil, fmt.Errorf("no top level deployment node named %q", s)
		}
		cid := 1
		if len(elems) > 2 {
			last := elems[len(elems)-1]
			if id, err := strconv.Atoi(last); err == nil {
				cid = id
				elems = elems[:len(elems)-1]
			}
		}
		for i := 1; i < len(elems)-1; i++ {
			parent = parent.Child(elems[i])
			if parent == nil {
				return nil, fmt.Errorf("no deployment node named %q in path %q", elems[i], s)
			}
		}
		name := elems[len(elems)-1]
		if dn := parent.Child(name); dn != nil {
			return dn, nil
		}
		if in := parent.InfrastructureNode(name); in != nil {
			return in, nil
		}
		if ci := parent.ContainerInstance(name, cid); ci != nil {
			return ci, nil
		}
		return nil, fmt.Errorf("could not find %q in path %q", name, s)
	default:
		return nil, fmt.Errorf("expected deployment node, infrastructure node, container instance or the path to one of these but bot %T", e)
	}
}

func parseLinkArgs(v expr.View, source interface{}, destination interface{}, args []interface{}) (src, dest expr.ElementHolder, desc string, dsl func(), err error) {
	var ok bool
	if dsl, ok = args[len(args)-1].(func()); ok {
		args = args[:len(args)-1]
	}
	if len(args) > 0 {
		desc, ok = args[0].(string)
		if !ok {
			err = fmt.Errorf("expected string (description), got %T", args[0])
			return
		}
		if len(args) > 1 {
			err = fmt.Errorf("too many arguments")
			return
		}
	}
	if src, err = findViewElement(v, source); err != nil {
		return
	}
	dest, err = findViewElement(v, destination)
	return
}
