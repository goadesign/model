/*
Package eval allows evaluating Go code that describes a software architecture
model using the DSL defined in Model: https://github.com/goadesign/model.
*/
package eval

import (
	"goa.design/goa/v3/eval"
	"goa.design/model/design"
	"goa.design/model/expr"
)

// RunDSL runs the DSL stored in a global variable and returns the corresponding
// Design expression.
func RunDSL() (*design.Design, error) {
	if err := eval.RunDSL(); err != nil {
		return nil, err
	}
	return modelize(expr.Root), nil
}

// modelize computes the software architecture design from the given DSL
// expressions.
func modelize(d *expr.Design) *design.Design {
	res := &design.Model{}
	m := d.Model
	if name := m.Enterprise; name != "" {
		res.Enterprise = &design.Enterprise{Name: name}
	}
	res.People = make([]*design.Person, len(m.People))
	for i, p := range m.People {
		res.People[i] = modelizePerson(p)
	}
	res.Systems = make([]*design.SoftwareSystem, len(m.Systems))
	for i, sys := range m.Systems {
		res.Systems[i] = modelizeSystem(sys)
	}
	res.DeploymentNodes = modelizeDeploymentNodes(m.DeploymentNodes)

	views := &design.Views{}
	v := d.Views
	views.LandscapeViews = make([]*design.LandscapeView, len(v.LandscapeViews))
	for i, lv := range v.LandscapeViews {
		views.LandscapeViews[i] = &design.LandscapeView{
			ViewProps:                 modelizeProps(lv.Props()),
			EnterpriseBoundaryVisible: lv.EnterpriseBoundaryVisible,
		}
	}
	views.ContextViews = make([]*design.ContextView, len(v.ContextViews))
	for i, lv := range v.ContextViews {
		views.ContextViews[i] = &design.ContextView{
			ViewProps:                 modelizeProps(lv.Props()),
			EnterpriseBoundaryVisible: lv.EnterpriseBoundaryVisible,
			SoftwareSystemID:          lv.SoftwareSystemID,
		}
	}
	views.ContainerViews = make([]*design.ContainerView, len(v.ContainerViews))
	for i, cv := range v.ContainerViews {
		views.ContainerViews[i] = &design.ContainerView{
			ViewProps:               modelizeProps(cv.Props()),
			SystemBoundariesVisible: cv.SystemBoundariesVisible,
			SoftwareSystemID:        cv.SoftwareSystemID,
		}
	}
	views.ComponentViews = make([]*design.ComponentView, len(v.ComponentViews))
	for i, cv := range v.ComponentViews {
		views.ComponentViews[i] = &design.ComponentView{
			ViewProps:                  modelizeProps(cv.Props()),
			ContainerBoundariesVisible: cv.ContainerBoundariesVisible,
			ContainerID:                cv.ContainerID,
		}
	}
	views.DynamicViews = make([]*design.DynamicView, len(v.DynamicViews))
	for i, dv := range v.DynamicViews {
		views.DynamicViews[i] = &design.DynamicView{
			ViewProps: modelizeProps(dv.Props()),
			ElementID: dv.ElementID,
		}
	}
	views.FilteredViews = make([]*design.FilteredView, len(v.FilteredViews))
	for i, lv := range v.FilteredViews {
		mode := "Include"
		if lv.Exclude {
			mode = "Exclude"
		}
		views.FilteredViews[i] = &design.FilteredView{
			Title:       lv.Title,
			Description: lv.Description,
			Key:         lv.Key,
			BaseKey:     lv.BaseKey,
			Mode:        mode,
			Tags:        lv.FilterTags,
		}
	}
	views.Styles = modelizeStyles(v.Styles)

	return &design.Design{
		Name:        d.Name,
		Description: d.Description,
		Version:     d.Version,
		Model:       res,
		Views:       views,
	}
}

func modelizePerson(p *expr.Person) *design.Person {
	return &design.Person{
		ID:            p.Element.ID,
		Name:          p.Element.Name,
		Description:   p.Element.Description,
		Technology:    p.Element.Technology,
		Tags:          p.Element.Tags,
		URL:           p.Element.URL,
		Properties:    p.Element.Properties,
		Relationships: modelizeRelationships(p.Relationships),
		Location:      p.Location,
	}
}

func modelizeRelationships(rels []*expr.Relationship) []*design.Relationship {
	res := make([]*design.Relationship, len(rels))
	for i, r := range rels {
		res[i] = &design.Relationship{
			ID:                   r.ID,
			Description:          r.Description,
			Tags:                 r.Tags,
			URL:                  r.URL,
			SourceID:             r.Source.ID,
			DestinationID:        r.Destination.ID,
			Technology:           r.Technology,
			InteractionStyle:     r.InteractionStyle,
			LinkedRelationshipID: r.LinkedRelationshipID,
		}
	}
	return res
}

func modelizeSystem(sys *expr.SoftwareSystem) *design.SoftwareSystem {
	return &design.SoftwareSystem{
		ID:            sys.ID,
		Name:          sys.Name,
		Description:   sys.Description,
		Technology:    sys.Technology,
		Tags:          sys.Tags,
		URL:           sys.URL,
		Properties:    sys.Properties,
		Relationships: modelizeRelationships(sys.Relationships),
		Location:      sys.Location,
		Containers:    modelizeContainers(sys.Containers),
	}
}

func modelizeContainers(cs []*expr.Container) []*design.Container {
	res := make([]*design.Container, len(cs))
	for i, c := range cs {
		res[i] = &design.Container{
			ID:            c.ID,
			Name:          c.Name,
			Description:   c.Description,
			Technology:    c.Technology,
			Tags:          c.Tags,
			URL:           c.URL,
			Properties:    c.Properties,
			Relationships: modelizeRelationships(c.Relationships),
			Components:    modelizeComponents(c.Components),
		}
	}
	return res
}

func modelizeComponents(cs []*expr.Component) []*design.Component {
	res := make([]*design.Component, len(cs))
	for i, c := range cs {
		res[i] = &design.Component{
			ID:            c.ID,
			Name:          c.Name,
			Description:   c.Description,
			Technology:    c.Technology,
			Tags:          c.Tags,
			URL:           c.URL,
			Properties:    c.Properties,
			Relationships: modelizeRelationships(c.Relationships),
		}
	}
	return res
}

func modelizeDeploymentNodes(dns []*expr.DeploymentNode) []*design.DeploymentNode {
	res := make([]*design.DeploymentNode, len(dns))
	for i, dn := range dns {
		res[i] = &design.DeploymentNode{
			ID:          dn.ID,
			Name:        dn.Name,
			Description: dn.Description,
			Technology:  dn.Technology,
			Environment: dn.Environment,
			Instances:   dn.Instances,
			Tags:        dn.Tags,
			URL:         dn.URL,
		}
	}
	return res
}

func modelizeProps(prop *expr.ViewProps) *design.ViewProps {
	props := &design.ViewProps{
		Title:             prop.Title,
		Description:       prop.Description,
		Key:               prop.Key,
		PaperSize:         prop.PaperSize,
		ElementViews:      modelizeElementViews(prop.ElementViews),
		RelationshipViews: modelizeRelationshipViews(prop.RelationshipViews),
		Animations:        modelizeAnimationSteps(prop.AnimationSteps),
	}
	if layout := prop.AutoLayout; layout != nil {
		props.AutoLayout = &design.AutoLayout{
			RankDirection: layout.RankDirection,
			RankSep:       layout.RankSep,
			NodeSep:       layout.NodeSep,
			EdgeSep:       layout.EdgeSep,
			Vertices:      layout.Vertices,
		}
	}
	return props
}

func modelizeElementViews(evs []*expr.ElementView) []*design.ElementView {
	res := make([]*design.ElementView, len(evs))
	for i, ev := range evs {
		res[i] = &design.ElementView{
			ID: ev.Element.ID,
			X:  ev.X,
			Y:  ev.Y,
		}
	}
	return res
}

func modelizeRelationshipViews(rvs []*expr.RelationshipView) []*design.RelationshipView {
	res := make([]*design.RelationshipView, len(rvs))
	for i, rv := range rvs {
		res[i] = &design.RelationshipView{
			ID:          rv.RelationshipID,
			Description: rv.Description,
			Order:       rv.Order,
			Vertices:    rv.Vertices,
			Routing:     rv.Routing,
			Position:    rv.Position,
		}
	}
	return res
}

func modelizeAnimationSteps(as []*expr.AnimationStep) []*design.AnimationStep {
	res := make([]*design.AnimationStep, len(as))
	for i, s := range as {
		elems := make([]string, len(s.Elements))
		for i, e := range s.Elements {
			elems[i] = e.GetElement().ID
		}
		res[i] = &design.AnimationStep{
			Order:         s.Order,
			Elements:      elems,
			Relationships: s.RelationshipIDs,
		}
	}
	return res
}

func modelizeStyles(s *expr.Styles) *design.Styles {
	elems := make([]*design.ElementStyle, len(s.Elements))
	for i, elem := range s.Elements {
		tmp := design.ElementStyle(*elem)
		elems[i] = &tmp
	}
	rels := make([]*design.RelationshipStyle, len(s.Relationships))
	for i, rel := range s.Relationships {
		tmp := design.RelationshipStyle(*rel)
		rels[i] = &tmp
	}
	return &design.Styles{
		Elements:      elems,
		Relationships: rels,
	}
}
