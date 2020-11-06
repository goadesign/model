package expr

import (
	"fmt"

	"goa.design/goa/v3/eval"
)

type (
	// Views is the container for all views.
	Views struct {
		LandscapeViews  []*LandscapeView
		ContextViews    []*ContextView
		ContainerViews  []*ContainerView
		ComponentViews  []*ComponentView
		DynamicViews    []*DynamicView
		DeploymentViews []*DeploymentView
		FilteredViews   []*FilteredView
		Styles          *Styles
		DSLFunc         func()
	}

	// LandscapeView describes a system landscape view.
	LandscapeView struct {
		*ViewProps
		EnterpriseBoundaryVisible *bool
	}

	// ContextView describes a system context view.
	ContextView struct {
		*ViewProps
		EnterpriseBoundaryVisible *bool
		SoftwareSystemID          string
	}

	// ContainerView describes a container view for a specific software
	// system.
	ContainerView struct {
		*ViewProps
		SystemBoundariesVisible *bool
		SoftwareSystemID        string
		AddInfluencers          bool
	}

	// ComponentView describes a component view for a specific container.
	ComponentView struct {
		*ViewProps
		ContainerBoundariesVisible *bool
		ContainerID                string
	}

	// DynamicView describes a dynamic view for a specified scope.
	DynamicView struct {
		*ViewProps
		ElementID string
	}

	// DeploymentView describes a deployment view.
	DeploymentView struct {
		*ViewProps
		SoftwareSystemID string
		Environment      string
	}

	// Styles describes the styles for a view.
	Styles struct {
		Elements      []*ElementStyle
		Relationships []*RelationshipStyle
	}

	// ElementStyle defines an element style.
	ElementStyle struct {
		Tag         string
		Shape       ShapeKind
		Icon        string
		Background  string
		Color       string
		Stroke      string
		Width       *int
		Height      *int
		FontSize    *int
		Metadata    *bool
		Description *bool
		Opacity     *int
		Border      BorderKind
	}

	// RelationshipStyle defines a relationship style.
	RelationshipStyle struct {
		Tag       string
		Thickness *int
		FontSize  *int
		Width     *int
		Position  *int
		Color     string
		Stroke    string
		Dashed    *bool
		Routing   RoutingKind
		Opacity   *int
	}

	// View is the common interface for all views.
	View interface {
		Props() *ViewProps
	}

	// ViewAdder is the interface implemented by views that allow adding
	// elements and animations explicitly.
	ViewAdder interface {
		AddElements(...ElementHolder) error
		AddAnimationStep(*AnimationStep) error
	}

	// ShapeKind is the enum used to represent shapes used to render elements.
	ShapeKind int

	// BorderKind is the enum used to represent element border styles.
	BorderKind int
)

const (
	ShapeUndefined ShapeKind = iota
	ShapeBox
	ShapeCircle
	ShapeCylinder
	ShapeEllipse
	ShapeHexagon
	ShapeRoundedBox
	ShapeComponent
	ShapeFolder
	ShapeMobileDeviceLandscape
	ShapeMobileDevicePortrait
	ShapePerson
	ShapePipe
	ShapeRobot
	ShapeWebBrowser
)

const (
	BorderUndefined BorderKind = iota
	BorderSolid
	BorderDashed
	BorderDotted
)

var (
	// Make sure views implement View.
	_ View = &LandscapeView{}
	_ View = &ContextView{}
	_ View = &ContainerView{}
	_ View = &ComponentView{}
	_ View = &DynamicView{}
	_ View = &DeploymentView{}

	// Make sure static views implement ViewAdder.
	_ ViewAdder = &LandscapeView{}
	_ ViewAdder = &ContextView{}
	_ ViewAdder = &ContainerView{}
	_ ViewAdder = &ComponentView{}
	_ ViewAdder = &DeploymentView{}
)

// DSL returns the DSL to execute.
func (vs *Views) DSL() func() {
	return vs.DSLFunc
}

// EvalName returns the generic expression name used in error messages.
func (vs *Views) EvalName() string {
	return "views"
}

// Validate makes sure the right element are in the right views, it also makes
// sure all animation steps have elements.
func (vs *Views) Validate() error {
	verr := new(eval.ValidationErrors)

	// Make sure views don't include elements that are not allowed for that type
	// of view.
	checkElements := func(title string, evs []*ElementView, allowContainers bool) {
		for _, ev := range evs {
			switch Registry[ev.Element.ID].(type) {
			case *SoftwareSystem, *Person:
				// all good
			case *Container:
				if !allowContainers {
					verr.Add(vs, fmt.Sprintf("%s can only contain software systems and people", title))
				}
			default:
				var suffix = " and people"
				if allowContainers {
					suffix = ", people and containers"
				}
				verr.Add(vs, fmt.Sprintf("%s can only contain software systems%s", title, suffix))
			}
		}
	}
	for _, lv := range vs.LandscapeViews {
		checkElements("software landscape views", lv.ElementViews, false)
	}
	for _, cv := range vs.ContextViews {
		checkElements("software context views", cv.ElementViews, false)
	}
	for _, cv := range vs.ContainerViews {
		checkElements("container views", cv.ElementViews, true)
	}

	for _, view := range vs.All() {
		v := view.Props()

		// Map relationship views created explicitly to model relationships.
		for _, rv := range v.RelationshipViews {
			srcID := rv.Source.ID
			destID := rv.Destination.ID
			desc := rv.Description

			// The relationships between container instances is implicitly
			// derived from the relationships between the corresponding
			// containers so make sure there is one for all relationships added
			// explicitly to the deployment view and if so create the
			// relationship between the container instances.
			sci, srcIsCI := Registry[rv.Source.ID].(*ContainerInstance)
			dci, destIsCI := Registry[rv.Destination.ID].(*ContainerInstance)
			if srcIsCI && destIsCI {
				srcID = sci.ContainerID
				destID = dci.ContainerID
			}

			IterateRelationships(func(r *Relationship) {
				if r.Destination == nil {
					return // a validation error was already created in model.Validate
				}
				if r.Source.ID == srcID && r.Destination.ID == destID && r.Description == desc {
					if srcIsCI && destIsCI {
						rci := r.Dup(sci.Element, dci.Element)
						rci.LinkedRelationshipID = r.ID
						sci.Relationships = append(sci.Relationships, rci)
						r = rci
					}
					rv.RelationshipID = r.ID
				}
			})
			if rv.RelationshipID == "" {
				verr.Add(rv, "could not find relationship %q [%s -> %s] to add to view %q", desc, rv.Source.Name, rv.Destination.Name, v.Key)
			}
		}

		// Make sure all elements used to remove unreachable are in scope.
		for _, e := range v.RemoveUnreachable {
			validateElementInView(v, e, "RemoveUnreachable", verr)
		}

		for i, s := range v.AnimationSteps {
			// Make sure all animation steps define at least one element.
			if len(s.Elements) == 0 {
				verr.AddError(v, fmt.Errorf("animation step %d in view %q introduces no new elements", i, v.Key))
			}
			// Make sure all animation step elements are in scope.
			for _, eh := range s.Elements {
				validateElementInView(v, eh.GetElement(), fmt.Sprintf("animation step %d", i), verr)
			}
		}
	}

	return verr
}

// Finalize relationships.
func (vs *Views) Finalize() {
	// Add influencers to container views.
	for _, view := range vs.ContainerViews {
		if view.AddInfluencers {
			addInfluencers(view)
		}
	}

	for _, view := range vs.All() {
		vp := view.Props()

		if vp.AddAll {
			addAllElements(view)
		} else if vp.AddDefault {
			addDefaultElements(view)
		}
		for _, e := range vp.AddNeighbors {
			addNeighbors(e, vp)
		}
		addMissingElementsAndRelationships(vp)
		addAnimationStepRelationships(vp)

		// Then remove elements and relationships that need to be removed
		// explicitly.
		for _, e := range vp.RemoveElements {
			removeElements(vp, e)
		}
		for _, r := range vp.RemoveRelationships {
			removeRelationship(vp, r)
		}
		for _, tag := range vp.RemoveTags {
			removeElements(vp, tagged(vp, tag)...)
		}
		for _, e := range vp.RemoveUnreachable {
			removeElements(vp, unreachable(vp, e)...)
		}
		if vp.RemoveUnrelated {
			removeElements(vp, unrelated(vp)...)
		}
		for _, ev := range vp.ElementViews {
			if ev.NoRelationship {
				i := 0
				for _, rv := range vp.RelationshipViews {
					if rv.Source.ID != ev.Element.ID && rv.Destination.ID != ev.Element.ID {
						vp.RelationshipViews[i] = rv
						i++
					}
				}
				vp.RelationshipViews = vp.RelationshipViews[:i]
			}
		}
	}
}

// All returns all the views in a single slice.
func (vs Views) All() (vps []View) {
	for _, lv := range vs.LandscapeViews {
		vps = append(vps, lv)
	}
	for _, cv := range vs.ContextViews {
		vps = append(vps, cv)
	}
	for _, cv := range vs.ContainerViews {
		vps = append(vps, cv)
	}
	for _, cv := range vs.ComponentViews {
		vps = append(vps, cv)
	}
	for _, dv := range vs.DynamicViews {
		vps = append(vps, dv)
	}
	for _, dv := range vs.DeploymentViews {
		vps = append(vps, dv)
	}
	return
}

// AddElements adds the given elements to the view if not already present.
func (cv *LandscapeView) AddElements(ehs ...ElementHolder) error {
	for _, eh := range ehs {
		if !isPS(eh) {
			return fmt.Errorf("elements of type %T cannot be added to landscape view", eh)
		}
	}
	addElements(cv.ViewProps, ehs...)
	return nil
}

// AddAnimationStep adds the given animation step to the view.
func (cv *LandscapeView) AddAnimationStep(s *AnimationStep) error {
	for _, eh := range s.Elements {
		if !isPS(eh) {
			return fmt.Errorf("elements of type %T cannot be added to an animation step in a landscape view", eh)
		}
	}
	return addAnimationStep(cv.ViewProps, s)
}

// AddElements adds the given elements to the view if not already present.
func (cv *ContextView) AddElements(ehs ...ElementHolder) error {
	for _, eh := range ehs {
		if !isPS(eh) {
			return fmt.Errorf("elements of type %T cannot be added to context view", eh)
		}
	}
	addElements(cv.ViewProps, ehs...)
	return nil
}

// AddAnimationStep adds the given animation step to the view.
func (cv *ContextView) AddAnimationStep(s *AnimationStep) error {
	for _, eh := range s.Elements {
		if !isPS(eh) {
			return fmt.Errorf("elements of type %T cannot be added to an animation step in a context view", eh)
		}
	}
	return addAnimationStep(cv.ViewProps, s)
}

// AddElements adds the given elements to the view if not already present.
func (cv *ContainerView) AddElements(ehs ...ElementHolder) error {
	for _, eh := range ehs {
		if !isPSC(eh) {
			return fmt.Errorf("elements of type %T cannot be added to container view", eh)
		}
	}
	addElements(cv.ViewProps, ehs...)
	return nil
}

// AddAnimationStep adds the given animation step to the view.
func (cv *ContainerView) AddAnimationStep(s *AnimationStep) error {
	for _, eh := range s.Elements {
		if !isPSC(eh) {
			return fmt.Errorf("elements of type %T cannot be added to an animation step in a container view", eh)
		}
	}
	return addAnimationStep(cv.ViewProps, s)
}

// AddElements adds the given elements to the view if not already present.
func (cv *ComponentView) AddElements(ehs ...ElementHolder) error {
	for _, eh := range ehs {
		if !isPSCC(eh) {
			return fmt.Errorf("elements of type %T cannot be added to component view", eh)
		}
	}
	addElements(cv.ViewProps, ehs...)
	return nil
}

// AddAnimationStep adds the given animation step to the view.
func (cv *ComponentView) AddAnimationStep(s *AnimationStep) error {
	for _, eh := range s.Elements {
		if !isPSCC(eh) {
			return fmt.Errorf("elements of type %T cannot be added to an animation step in a component view", eh)
		}
	}
	return addAnimationStep(cv.ViewProps, s)
}

// AddElements adds the given elements to the view if not already present.
func (dv *DeploymentView) AddElements(ehs ...ElementHolder) error {
	var nodes []*DeploymentNode
	for _, eh := range ehs {
		switch e := eh.(type) {
		case *DeploymentNode:
			if addDeploymentNodeChildren(dv, e) {
				nodes = append(nodes, e)
			}
		case *ContainerInstance:
			if dv.SoftwareSystemID == "" || dv.SoftwareSystemID == Registry[e.ContainerID].(*Container).System.ID {
				addElements(dv.ViewProps, e)
				nodes = append(nodes, e.Parent)
			}
		case *InfrastructureNode:
			addElements(dv.ViewProps, e)
			nodes = append(nodes, e.Parent)
		default:
			return fmt.Errorf("elements of type %T cannot be added to deployment views", eh)
		}
	}

	// Add deployment node hierarchy.
	for _, n := range nodes {
		addElements(dv.ViewProps, n)
		p := n.Parent
		for p != nil {
			addElements(dv.ViewProps, p)
			p = p.Parent
		}
	}

	return nil
}

// AddAnimationStep adds the given animation step to the view.
func (dv *DeploymentView) AddAnimationStep(s *AnimationStep) error {
	for _, eh := range s.Elements {
		if !isDCI(eh) {
			return fmt.Errorf("elements of type %T cannot be added to an animation step in a deployment view", eh)
		}
	}
	return addAnimationStep(dv.ViewProps, s)
}

// EvalName returns the generic expression name used in error messages.
func (c *Styles) EvalName() string {
	return "styles"
}

// EvalName returns the generic expression name used in error messages.
func (es *ElementStyle) EvalName() string {
	return fmt.Sprintf("element style for tag %q", es.Tag)
}

// EvalName returns the generic expression name used in error messages.
func (rs *RelationshipStyle) EvalName() string {
	return fmt.Sprintf("relationship style for tag %q", rs.Tag)
}

// isPS returns true if element is a person or software system, false otherwise.
func isPS(eh ElementHolder) bool {
	switch eh.(type) {
	case *Person, *SoftwareSystem:
		return true
	}
	return false
}

// isPSC returns true if element is a person, a software system or a container,
// false otherwise.
func isPSC(eh ElementHolder) bool {
	if isPS(eh) {
		return true
	}
	_, ok := eh.(*Container)
	return ok
}

// isPSCC returns true if element is a person, a software system, a container or
// a component, false otherwise.
func isPSCC(eh ElementHolder) bool {
	if isPSC(eh) {
		return true
	}
	_, ok := eh.(*Component)
	return ok
}

// isDCI returns true if element is a deployment node, a container instance or
// an infrastructure node, false otherwise.
func isDCI(eh ElementHolder) bool {
	switch eh.(type) {
	case *DeploymentNode, *ContainerInstance, *InfrastructureNode:
		return true
	}
	return false
}

// addElements adds the given elements to the view if not already present.
func addElements(v *ViewProps, ehs ...ElementHolder) {
loop:
	for _, eh := range ehs {
		e := eh.GetElement()
		for _, ev := range v.ElementViews {
			if ev.Element.ID == e.ID {
				continue loop
			}
		}
		v.ElementViews = append(v.ElementViews, &ElementView{Element: e})
	}
}

// addDeploymentNodeChildren adds the children, infrastructure nodes and container
// instances of n to dv and returns true if anything was added, false otherwise.
func addDeploymentNodeChildren(dv *DeploymentView, n *DeploymentNode) bool {
	var nested bool
	for _, inst := range n.ContainerInstances {
		if dv.SoftwareSystemID == "" || Registry[inst.ContainerID].(*Container).System.ID == dv.SoftwareSystemID {
			addElements(dv.ViewProps, inst)
			nested = true
		}
	}
	for _, inf := range n.InfrastructureNodes {
		addElements(dv.ViewProps, inf)
		nested = true
	}
	for _, c := range n.Children {
		if nest := addDeploymentNodeChildren(dv, c); nest {
			addElements(dv.ViewProps, c)
			nested = true
		}
	}
	return nested
}

// addAnimation adds the animations to the view after normalizing their content.
// It makes sure that all the elements are in the view, that each element is
// only included in one animation step, that the relationships are initialized
// and that any dependent deployment nodes is added.
func addAnimationStep(v *ViewProps, s *AnimationStep) error {
	var known []ElementHolder
	for _, as := range v.AnimationSteps {
		known = append(known, as.Elements...)
	}
	s.Order = len(v.AnimationSteps) + 1
	var filtered []ElementHolder
loop:
	for _, e := range s.Elements {
		if e == nil {
			return fmt.Errorf("element not initialized")
		}
		id := e.GetElement().ID
		for _, k := range known {
			if k.GetElement().ID == id {
				continue loop // item already in a step, skip
			}
		}
		known = append(known, e)
		filtered = append(filtered, e)

		// Add parent deployment nodes for infrastructure nodes
		// and container instances.
		var node *DeploymentNode
		if inf, ok := e.(*InfrastructureNode); ok {
			node = inf.Parent
		} else if ci, ok := e.(*ContainerInstance); ok {
			node = ci.Parent
		}
		for node != nil {
			known = append(known, node)
			filtered = append(filtered, node)
			node = node.Parent
		}
	}
	if len(filtered) == 0 {
		return fmt.Errorf("none of the specified elements exist in this view or do not already appear in previous animation steps")
	}
	s.Elements = filtered

	v.AnimationSteps = append(v.AnimationSteps, s)
	return nil
}

// validateElementInView makes sure there is an ElementView corresponding to e
// in v. It adds an error to verr using title if that's not the case.
func validateElementInView(v *ViewProps, e *Element, title string, verr *eval.ValidationErrors) {
	for _, ev := range v.ElementViews {
		if ev.Element.ID == e.ID {
			return
		}
	}
	verr.Add(v, "%T %q used in %s not added to the view %q", e, e.Name, title, v.Key)
}
