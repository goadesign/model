package expr

import (
	"fmt"
	"strings"
)

// addAllElements adds all top level elements (people and software systems) as
// well as all elements in scope to the view.
func addAllElements(view View) {
	m := Root.Model
	switch v := view.(type) {
	case *LandscapeView:
		v.AddElements(m.People.Elements()...)
		v.AddElements(m.Systems.Elements()...)
	case *ContextView:
		v.AddElements(m.People.Elements()...)
		v.AddElements(m.Systems.Elements()...)
	case *ContainerView:
		v.AddElements(m.People.Elements()...)
		v.AddElements(m.Systems.Elements()...)
		v.AddElements(Registry[v.SoftwareSystemID].(*SoftwareSystem).Containers.Elements()...)
		removeElements(v.Props(), Registry[v.SoftwareSystemID].(*SoftwareSystem).Element)
	case *ComponentView:
		v.AddElements(m.People.Elements()...)
		v.AddElements(m.Systems.Elements()...)
		c := Registry[v.ContainerID].(*Container)
		v.AddElements(c.System.Containers.Elements()...)
		v.AddElements(c.Components.Elements()...)
	case *DeploymentView:
		for _, n := range m.DeploymentNodes {
			if n.Environment == "" || n.Environment == v.Environment {
				v.AddElements(n)
			}
		}
	default:
		panic(fmt.Sprintf("AddAllElements called on %T", view))
	}
}

// addDefaultElements adds a default set of elements and relationships for the
// given view.
func addDefaultElements(view View) {
	// Add default elements if needed.
	switch v := view.(type) {
	case *LandscapeView:
		addAllElements(v)
	case *ContextView:
		s := Registry[v.SoftwareSystemID].(*SoftwareSystem)
		v.AddElements(s)
		addNeighbors(s.Element, v)
	case *ContainerView:
		s := Registry[v.SoftwareSystemID].(*SoftwareSystem)
		v.AddElements(s.Containers.Elements()...)
		for _, c := range s.Containers {
			v.AddElements(relatedSoftwareSystems(c.Element).Elements()...)
			v.AddElements(relatedPeople(c.Element).Elements()...)
		}
	case *ComponentView:
		c := Registry[v.ContainerID].(*Container)
		v.AddElements(c.Components.Elements()...)
		for _, c := range c.Components {
			v.AddElements(relatedContainers(c.Element).Elements()...)
			v.AddElements(relatedSoftwareSystems(c.Element).Elements()...)
			v.AddElements(relatedPeople(c.Element).Elements()...)
		}
	case *DeploymentView:
		addAllElements(v)
	}
}

// addMissingElementsAndRelationships adds all elements that form edges of
// relationships in the view and adds all relationships between elements that
// are in the view.
func addMissingElementsAndRelationships(vp *ViewProps) {
	for _, rv := range vp.RelationshipViews {
		addElements(vp, rv.Source)
		addElements(vp, rv.Destination)
	}
	for _, ev := range vp.ElementViews {
		var rels []*Relationship
		IterateRelationships(func(r *Relationship) {
			if r.Source.ID == ev.Element.ID {
				for _, ev2 := range vp.ElementViews {
					if r.Destination.ID == ev2.Element.ID {
						rels = append(rels, r)
						break
					}
				}
			}
		})
	loop:
		for _, r := range rels {
			// Do not add previously added relationship views
			for _, existing := range vp.RelationshipViews {
				if r.ID == existing.RelationshipID {
					continue loop
				}
			}

			// Do not automatically add relationship views across different
			// top-level deployment nodes.
			//
			// Note: this rule is a little bit arbitrary however it is possible
			// to override the behavior using `Link` and `Unlink` explicitly in
			// the design. We'll see how that works out over time.
			top := func(d *DeploymentNode) string {
				id := d.ID
				for p := d.Parent; p != nil; p = p.Parent {
					id = p.ID
				}
				return id
			}
			var srcTop, destTop string
			switch s := Registry[r.Source.ID].(type) {
			case *DeploymentNode:
				srcTop = top(s)
			case *InfrastructureNode:
				srcTop = top(s.Parent)
			case *ContainerInstance:
				srcTop = top(s.Parent)
			}
			switch d := Registry[r.Destination.ID].(type) {
			case *DeploymentNode:
				destTop = top(d)
			case *InfrastructureNode:
				destTop = top(d.Parent)
			case *ContainerInstance:
				destTop = top(d.Parent)
			}
			if srcTop != destTop {
				continue loop
			}

			vp.RelationshipViews = append(vp.RelationshipViews,
				&RelationshipView{
					Source:         r.Source,
					Destination:    r.Destination,
					Description:    r.Description,
					RelationshipID: r.ID,
				})
		}
	}
}

func addNeighbors(e *Element, view View) {
	switch v := view.(type) {
	case *LandscapeView:
		v.AddElements(relatedPeople(e).Elements()...)
		v.AddElements(relatedSoftwareSystems(e).Elements()...)
	case *ContextView:
		v.AddElements(relatedPeople(e).Elements()...)
		v.AddElements(relatedSoftwareSystems(e).Elements()...)
	case *ContainerView:
		v.AddElements(relatedPeople(e).Elements()...)
		v.AddElements(relatedSoftwareSystems(e).Elements()...)
		v.AddElements(relatedContainers(e).Elements()...)
	case *ComponentView:
		v.AddElements(relatedPeople(e).Elements()...)
		v.AddElements(relatedSoftwareSystems(e).Elements()...)
		v.AddElements(relatedContainers(e).Elements()...)
		v.AddElements(relatedComponents(e).Elements()...)
	case *DeploymentView:
		v.AddElements(relatedInfrastructureNodes(e).Elements()...)
		v.AddElements(relatedContainerInstances(e).Elements()...)
	}

}

func addInfluencers(cv *ContainerView) {
	system := Registry[cv.SoftwareSystemID].(*SoftwareSystem)
	m := Root.Model
	for _, s := range m.Systems {
		for _, r := range s.Relationships {
			if r.Destination.ID == cv.SoftwareSystemID {
				cv.AddElements(s)
			}
		}
		for _, r := range system.Relationships {
			if r.Destination.ID == s.ID {
				cv.AddElements(s)
			}
		}
	}

	for _, p := range m.People {
		for _, r := range p.Relationships {
			if r.Destination.ID == cv.SoftwareSystemID {
				cv.AddElements(p)
			}
		}
		for _, r := range system.Relationships {
			if r.Destination.ID == p.ID {
				cv.AddElements(p)
			}
		}
	}

	for i, rv := range cv.RelationshipViews {
		src := rv.Source
		var keep bool
		if keep = src.ID == cv.SoftwareSystemID; !keep {
			if c, ok := Registry[src.ID].(*Container); ok {
				keep = c.System.ID == cv.SoftwareSystemID
			} else if c, ok := Registry[src.ID].(*Component); ok {
				keep = c.Container.System.ID == cv.SoftwareSystemID
			}
		}
		if !keep {
			dest := rv.Destination
			if keep = dest.ID == cv.SoftwareSystemID; !keep {
				if c, ok := Registry[dest.ID].(*Container); ok {
					keep = c.System.ID == cv.SoftwareSystemID
				} else if c, ok := Registry[dest.ID].(*Component); ok {
					keep = c.Container.System.ID == cv.SoftwareSystemID
				}
			}
		}
		if !keep {
			cv.RelationshipViews = append(cv.RelationshipViews[:i], cv.RelationshipViews[i+1:]...)
		}
	}
}

// Add implied animation step relationships
func addAnimationStepRelationships(vp *ViewProps) {
	for _, s := range vp.AnimationSteps {
		var newSrc, newDest, oldSrc, oldDest bool
		for _, rv := range vp.RelationshipViews {
			for _, eh := range s.Elements {
				id := eh.GetElement().ID
				if id == rv.Source.ID {
					oldSrc = true
				} else if id == rv.Destination.ID {
					oldDest = true
				}
				if oldSrc && oldDest {
					break
				}
			}
			if oldSrc && oldDest {
				break
			}
			for _, ev := range vp.ElementViews {
				if ev.Element.ID == rv.Source.ID {
					newSrc = true
				} else if ev.Element.ID == rv.Destination.ID {
					newDest = true
				}
				if newSrc && newDest {
					break
				}
			}
			if newSrc && oldDest || oldSrc && newDest {
				s.RelationshipIDs = append(s.RelationshipIDs, rv.RelationshipID)
			}
		}
	}
}

// removeElements removes the given elements from the given view as well as any
// relationship that uses the element as source or destination.
func removeElements(vp *ViewProps, elems ...*Element) {
	for _, e := range elems {
		i := 0
		for _, ev := range vp.ElementViews {
			if ev.Element.ID != e.ID {
				vp.ElementViews[i] = ev
				i++
			} else {
				// Remove corresponding relationships.
				j := 0
				for _, rv := range vp.RelationshipViews {
					if rv.Source.ID != e.ID && rv.Destination.ID != e.ID {
						vp.RelationshipViews[j] = rv
						j++
					}
				}
				vp.RelationshipViews = vp.RelationshipViews[:j]
			}
		}
		vp.ElementViews = vp.ElementViews[:i]
	}
}

func removeRelationship(vp *ViewProps, r *Relationship) {
	i := 0
	for _, rv := range vp.RelationshipViews {
		if rv.Source.ID != r.Source.ID || rv.Destination.ID != r.Destination.ID || rv.Description != r.Description {
			vp.RelationshipViews[i] = rv
			i++
		}
	}
	vp.RelationshipViews = vp.RelationshipViews[:i]
}

// allUnrelated fetches all elements that have no relationship to other elements
// in the view.
func unrelated(v *ViewProps) (elems []*Element) {
loop:
	for _, ev := range v.ElementViews {
		for _, rv := range v.RelationshipViews {
			if rv.Source.ID == ev.Element.ID || rv.Destination.ID == ev.Element.ID {
				continue loop
			}
		}
		elems = append(elems, ev.Element)
	}
	return
}

// relatedPeople returns all people the element has a relationship with
// (either as source or as destination).
func relatedPeople(elem *Element) (res People) {
	add := func(p *Person) {
		for _, ep := range res {
			if ep.ID == p.ID {
				return
			}
		}
		res = append(res, p)
	}
	IterateRelationships(func(r *Relationship) {
		if r.Source.ID == elem.ID {
			if p, ok := Registry[r.Destination.ID].(*Person); ok {
				add(p)
			}
		}
		if r.Destination.ID == elem.ID {
			if p, ok := Registry[r.Source.ID].(*Person); ok {
				add(p)
			}
		}
	})
	return
}

// relatedSoftwareSystems returns all software systems the element has a
// relationship with (either as source or as destination).
func relatedSoftwareSystems(elem *Element) (res SoftwareSystems) {
	add := func(s *SoftwareSystem) {
		for _, es := range res {
			if es.ID == s.ID {
				return
			}
		}
		res = append(res, s)
	}
	IterateRelationships(func(r *Relationship) {
		if r.Source.ID == elem.ID {
			if s, ok := Registry[r.Destination.ID].(*SoftwareSystem); ok {
				add(s)
			}
		}
		if r.Destination.ID == elem.ID {
			if s, ok := Registry[r.Source.ID].(*SoftwareSystem); ok {
				add(s)
			}
		}
	})
	return
}

// relatedContainers returns all containers the element has a relationship with
// (either as source or as destination).
func relatedContainers(elem *Element) (res Containers) {
	add := func(cc *Container) {
		for _, es := range res {
			if es.ID == cc.ID {
				return
			}
		}
		res = append(res, cc)
	}
	IterateRelationships(func(r *Relationship) {
		if r.Source.ID == elem.ID {
			if c, ok := Registry[r.Destination.ID].(*Container); ok {
				add(c)
			}
		}
		if r.Destination.ID == elem.ID {
			if c, ok := Registry[r.Source.ID].(*Container); ok {
				add(c)
			}
		}
	})
	return
}

// relatedComponents returns all components the element has a relationship with
// (either as source or as destination).
func relatedComponents(elem *Element) (res Components) {
	add := func(c *Component) {
		for _, es := range res {
			if es.ID == c.ID {
				return
			}
		}
		res = append(res, c)
	}
	IterateRelationships(func(r *Relationship) {
		if r.Source.ID == elem.ID {
			if c, ok := Registry[r.Destination.ID].(*Component); ok {
				add(c)
			}
		}
		if r.Destination.ID == elem.ID {
			if c, ok := Registry[r.Source.ID].(*Component); ok {
				add(c)
			}
		}
	})
	return
}

// relatedInfrastructureNodes returns all infrastructure nodes the element has a
// relationship with (either as source or as destination).
func relatedInfrastructureNodes(elem *Element) (res InfrastructureNodes) {
	add := func(i *InfrastructureNode) {
		for _, inf := range res {
			if inf.ID == i.ID {
				return
			}
		}
		res = append(res, i)
	}
	IterateRelationships(func(r *Relationship) {
		if r.Source.ID == elem.ID {
			if inf, ok := Registry[r.Destination.ID].(*InfrastructureNode); ok {
				add(inf)
			}
		}
		if r.Destination.ID == elem.ID {
			if inf, ok := Registry[r.Source.ID].(*InfrastructureNode); ok {
				add(inf)
			}
		}
	})
	return
}

// relatedContainerInstances returns all container instances the element has a
// relationship with (either as source or as destination).
func relatedContainerInstances(elem *Element) (res ContainerInstances) {
	add := func(ci *ContainerInstance) {
		for _, eci := range res {
			if eci.ID == ci.ID {
				return
			}
		}
		res = append(res, ci)
	}
	IterateRelationships(func(r *Relationship) {
		if r.Source.ID == elem.ID {
			if ci, ok := Registry[r.Destination.ID].(*ContainerInstance); ok {
				add(ci)
			}
		}
		if r.Destination.ID == elem.ID {
			if ci, ok := Registry[r.Source.ID].(*ContainerInstance); ok {
				add(ci)
			}
		}
	})
	return
}

// allUnreachable fetches all elements in view not reachable from eh (directory
// or not).
func unreachable(v *ViewProps, eh ElementHolder) (elems []*Element) {
	ids := reachable(eh.GetElement())
loop:
	for _, e := range v.ElementViews {
		for _, id := range ids {
			if id == e.Element.ID {
				continue loop
			}
		}
		elems = append(elems, e.Element)
	}
	return
}

// allTagged returns all elements with the given tag in the view.
func tagged(v *ViewProps, tag string) (elems []*Element) {
	for _, ev := range v.ElementViews {
		vals := strings.Split(ev.Element.Tags, ",")
		for _, val := range vals {
			if val == tag {
				elems = append(elems, ev.Element)
				break
			}
		}
	}
	return
}
