package mdl

import (
	"goa.design/goa/v3/eval"
	"goa.design/model/expr"
)

type (
	// Design is the root node of the data model
	Design struct {
		// Name of the design.
		Name string `json:"name"`
		// Description of the design if any.
		Description string `json:"description,omitempty"`
		// Version number for the design.
		Version string `json:"version,omitempty"`
		// Model is the software architecture model.
		Model *Model `json:"model,omitempty"`
		// Views contains the views if any.
		Views *Views `json:"views,omitempty"`
	}
)

// RunDSL runs the DSL defined in a global variable and returns a JSON serializable version of it.
func RunDSL() (*Design, error) {
	if err := eval.RunDSL(); err != nil {
		return nil, err
	}
	return ModelizeDesign(expr.Root), nil
}

func ModelizeDesign(d *expr.Design) *Design {
	model := &Model{}
	m := d.Model
	if name := m.Enterprise; name != "" {
		model.Enterprise = &Enterprise{Name: name}
	}
	model.People = make([]*Person, len(m.People))
	for i, p := range m.People {
		model.People[i] = modelizePerson(p)
	}
	model.Systems = make([]*SoftwareSystem, len(m.Systems))
	for i, sys := range m.Systems {
		model.Systems[i] = modelizeSystem(sys)
	}
	model.DeploymentNodes = modelizeDeploymentNodes(m.DeploymentNodes)

	views := &Views{}
	v := d.Views
	views.LandscapeViews = make([]*LandscapeView, len(v.LandscapeViews))
	for i, lv := range v.LandscapeViews {
		views.LandscapeViews[i] = &LandscapeView{
			ViewProps:                 modelizeProps(lv.Props()),
			EnterpriseBoundaryVisible: lv.EnterpriseBoundaryVisible,
		}
	}
	views.ContextViews = make([]*ContextView, len(v.ContextViews))
	for i, lv := range v.ContextViews {
		views.ContextViews[i] = &ContextView{
			ViewProps:                 modelizeProps(lv.Props()),
			EnterpriseBoundaryVisible: lv.EnterpriseBoundaryVisible,
			SoftwareSystemID:          lv.SoftwareSystemID,
		}
	}
	views.ContainerViews = make([]*ContainerView, len(v.ContainerViews))
	for i, cv := range v.ContainerViews {
		views.ContainerViews[i] = &ContainerView{
			ViewProps:               modelizeProps(cv.Props()),
			SystemBoundariesVisible: cv.SystemBoundariesVisible,
			SoftwareSystemID:        cv.SoftwareSystemID,
		}
	}
	views.ComponentViews = make([]*ComponentView, len(v.ComponentViews))
	for i, cv := range v.ComponentViews {
		views.ComponentViews[i] = &ComponentView{
			ViewProps:                  modelizeProps(cv.Props()),
			ContainerBoundariesVisible: cv.ContainerBoundariesVisible,
			ContainerID:                cv.ContainerID,
		}
	}
	views.DynamicViews = make([]*DynamicView, len(v.DynamicViews))
	for i, dv := range v.DynamicViews {
		views.DynamicViews[i] = &DynamicView{
			ViewProps: modelizeProps(dv.Props()),
			ElementID: dv.ElementID,
		}
	}
	views.DeploymentViews = make([]*DeploymentView, len(v.DeploymentViews))
	for i, dv := range v.DeploymentViews {
		views.DeploymentViews[i] = &DeploymentView{
			ViewProps:        modelizeProps(dv.Props()),
			SoftwareSystemID: dv.SoftwareSystemID,
			Environment:      dv.Environment,
		}
	}
	views.FilteredViews = make([]*FilteredView, len(v.FilteredViews))
	for i, lv := range v.FilteredViews {
		mode := "Include"
		if lv.Exclude {
			mode = "Exclude"
		}
		views.FilteredViews[i] = &FilteredView{
			Title:       lv.Title,
			Description: lv.Description,
			Key:         lv.Key,
			BaseKey:     lv.BaseKey,
			Mode:        mode,
			Tags:        lv.FilterTags,
		}
	}
	views.Styles = modelizeStyles(v.Styles)

	return &Design{
		Name:        d.Name,
		Description: d.Description,
		Version:     d.Version,
		Model:       model,
		Views:       views,
	}
}

func modelizePerson(p *expr.Person) *Person {
	return &Person{
		ID:            p.Element.ID,
		Name:          p.Element.Name,
		Description:   p.Element.Description,
		Tags:          p.Element.Tags,
		URL:           p.Element.URL,
		Properties:    p.Element.Properties,
		Relationships: modelizeRelationships(p.Relationships),
		Location:      LocationKind(p.Location),
	}
}

func modelizeRelationships(rels []*expr.Relationship) []*Relationship {
	res := make([]*Relationship, len(rels))
	for i, r := range rels {
		res[i] = &Relationship{
			ID:                   r.ID,
			Description:          r.Description,
			Tags:                 r.Tags,
			URL:                  r.URL,
			SourceID:             r.Source.ID,
			DestinationID:        r.Destination.ID,
			Technology:           r.Technology,
			InteractionStyle:     InteractionStyleKind(r.InteractionStyle),
			LinkedRelationshipID: r.LinkedRelationshipID,
		}
	}
	return res
}

func modelizeSystem(sys *expr.SoftwareSystem) *SoftwareSystem {
	return &SoftwareSystem{
		ID:            sys.ID,
		Name:          sys.Name,
		Description:   sys.Description,
		Tags:          sys.Tags,
		URL:           sys.URL,
		Properties:    sys.Properties,
		Relationships: modelizeRelationships(sys.Relationships),
		Location:      LocationKind(sys.Location),
		Containers:    modelizeContainers(sys.Containers),
	}
}

func modelizeContainers(cs []*expr.Container) []*Container {
	res := make([]*Container, len(cs))
	for i, c := range cs {
		res[i] = &Container{
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

func modelizeComponents(cs []*expr.Component) []*Component {
	res := make([]*Component, len(cs))
	for i, c := range cs {
		res[i] = &Component{
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

func modelizeDeploymentNodes(dns []*expr.DeploymentNode) []*DeploymentNode {
	res := make([]*DeploymentNode, len(dns))
	for i, dn := range dns {
		children := modelizeDeploymentNodes(dn.Children)
		infs := make([]*InfrastructureNode, len(dn.InfrastructureNodes))
		for i, inf := range dn.InfrastructureNodes {
			infs[i] = &InfrastructureNode{
				ID:            inf.ID,
				Name:          inf.Name,
				Description:   inf.Description,
				Technology:    inf.Technology,
				Tags:          inf.Tags,
				URL:           inf.URL,
				Properties:    inf.Properties,
				Relationships: modelizeRelationships(inf.Relationships),
				Environment:   inf.Environment,
			}
		}
		cis := make([]*ContainerInstance, len(dn.ContainerInstances))
		for i, ci := range dn.ContainerInstances {
			cis[i] = &ContainerInstance{
				ID:            ci.ID,
				Tags:          ci.Tags,
				URL:           ci.URL,
				Properties:    ci.Properties,
				Relationships: modelizeRelationships(ci.Relationships),
				ContainerID:   ci.ContainerID,
				InstanceID:    ci.InstanceID,
				Environment:   ci.Environment,
				HealthChecks:  modelizeHealthChecks(ci.HealthChecks),
			}
		}
		res[i] = &DeploymentNode{
			ID:                  dn.ID,
			Name:                dn.Name,
			Description:         dn.Description,
			Technology:          dn.Technology,
			Environment:         dn.Environment,
			Children:            children,
			InfrastructureNodes: infs,
			ContainerInstances:  cis,
			Instances:           dn.Instances,
			Tags:                dn.Tags,
			URL:                 dn.URL,
		}
	}
	return res
}

func modelizeHealthChecks(hcs []*expr.HealthCheck) []*HealthCheck {
	res := make([]*HealthCheck, len(hcs))
	for i, hc := range hcs {
		res[i] = &HealthCheck{
			Name:     hc.Name,
			URL:      hc.URL,
			Interval: hc.Interval,
			Timeout:  hc.Timeout,
			Headers:  hc.Headers,
		}
	}
	return res
}

func modelizeProps(prop *expr.ViewProps) *ViewProps {
	props := &ViewProps{
		Title:             prop.Title,
		Description:       prop.Description,
		Key:               prop.Key,
		PaperSize:         PaperSizeKind(prop.PaperSize),
		ElementViews:      modelizeElementViews(prop.ElementViews),
		RelationshipViews: modelizeRelationshipViews(prop.RelationshipViews),
		Animations:        modelizeAnimationSteps(prop.AnimationSteps),
		Settings:          modelizeSettings(prop),
	}
	if layout := prop.AutoLayout; layout != nil {
		props.AutoLayout = &AutoLayout{
			RankDirection: RankDirectionKind(layout.RankDirection),
			RankSep:       layout.RankSep,
			NodeSep:       layout.NodeSep,
			EdgeSep:       layout.EdgeSep,
			Vertices:      layout.Vertices,
		}
	}
	return props
}

func modelizeElementViews(evs []*expr.ElementView) []*ElementView {
	res := make([]*ElementView, len(evs))
	for i, ev := range evs {
		res[i] = &ElementView{
			ID: ev.Element.ID,
			X:  ev.X,
			Y:  ev.Y,
		}
	}
	return res
}

func modelizeRelationshipViews(rvs []*expr.RelationshipView) []*RelationshipView {
	res := make([]*RelationshipView, len(rvs))
	for i, rv := range rvs {
		vertices := make([]*Vertex, len(rv.Vertices))
		for i, v := range rv.Vertices {
			vertices[i] = &Vertex{v.X, v.Y}
		}
		res[i] = &RelationshipView{
			ID:          rv.RelationshipID,
			Description: rv.Description,
			Order:       rv.Order,
			Vertices:    vertices,
			Routing:     RoutingKind(rv.Routing),
			Position:    rv.Position,
		}
	}
	return res
}

func modelizeAnimationSteps(as []*expr.AnimationStep) []*AnimationStep {
	res := make([]*AnimationStep, len(as))
	for i, s := range as {
		elems := make([]string, len(s.Elements))
		for i, e := range s.Elements {
			elems[i] = e.GetElement().ID
		}
		res[i] = &AnimationStep{
			Order:         s.Order,
			Elements:      elems,
			Relationships: s.RelationshipIDs,
		}
	}
	return res
}

func modelizeSettings(p *expr.ViewProps) *ViewSettings {
	var addNeighborsIDs []string
	for _, e := range p.AddNeighbors {
		addNeighborsIDs = append(addNeighborsIDs, e.GetElement().ID)
	}
	var removeElementsIDs []string
	for _, e := range p.RemoveElements {
		removeElementsIDs = append(removeElementsIDs, e.GetElement().ID)
	}
	var removeRelationshipsIDs []string
	for _, e := range p.RemoveRelationships {
		removeRelationshipsIDs = append(removeRelationshipsIDs, e.ID)
	}
	var removeUnreachableIDs []string
	for _, e := range p.RemoveUnreachable {
		removeUnreachableIDs = append(removeUnreachableIDs, e.GetElement().ID)
	}
	return &ViewSettings{
		AddAll:                p.AddAll,
		AddDefault:            p.AddDefault,
		AddNeighborIDs:        addNeighborsIDs,
		RemoveElementIDs:      removeElementsIDs,
		RemoveTags:            p.RemoveTags,
		RemoveRelationshipIDs: removeRelationshipsIDs,
		RemoveUnreachableIDs:  removeUnreachableIDs,
		RemoveUnrelated:       p.RemoveUnrelated,
	}
}

func modelizeStyles(s *expr.Styles) *Styles {
	if s == nil {
		return nil
	}
	elems := make([]*ElementStyle, len(s.Elements))
	for i, es := range s.Elements {
		elems[i] = &ElementStyle{
			Tag:         es.Tag,
			Background:  es.Background,
			Width:       es.Width,
			Height:      es.Height,
			FontSize:    es.FontSize,
			Icon:        es.Icon,
			Stroke:      es.Stroke,
			Color:       es.Color,
			Shape:       ShapeKind(es.Shape),
			Metadata:    es.Metadata,
			Description: es.Description,
			Opacity:     es.Opacity,
			Border:      BorderKind(es.Border),
		}
	}
	rels := make([]*RelationshipStyle, len(s.Relationships))
	for i, rs := range s.Relationships {
		rels[i] = &RelationshipStyle{
			Tag:       rs.Tag,
			Color:     rs.Color,
			Dashed:    rs.Dashed,
			Routing:   RoutingKind(rs.Routing),
			Opacity:   rs.Opacity,
			Thickness: rs.Thickness,
			Width:     rs.Width,
			FontSize:  rs.FontSize,
			Position:  rs.Position,
		}
	}
	return &Styles{
		Elements:      elems,
		Relationships: rels,
	}
}
