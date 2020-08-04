package expr

// Merge merges other into this workspace. The merge algorithm recursively
// overrides all fields of w with fields from other that do not have the zero
// value. Merge does not change the ID of w.
func (w *Workspace) Merge(other *Workspace) {
	w.Name = other.Name
	if other.Description != "" {
		w.Description = other.Description
	}
	if other.Version != "" {
		w.Version = other.Version
	}
	w.Revision = other.Revision
	if other.Thumbnail != "" {
		w.Thumbnail = other.Thumbnail
	}

	// idmap maps IDs of elements of w to IDs of elements of other.
	idmap := buildIDMap(w, other)

	if w.Model == nil {
		w.Model = other.Model
	} else {
		mergeModels(w.Model, other.Model, idmap)
	}
	if w.Views == nil {
		w.Views = other.Views
	} else {
		mergeViews(w.Views, other.Views, idmap)
	}
	if w.Documentation == nil {
		w.Documentation = other.Documentation
	} else {
		mergeDocumentation(w.Documentation, other.Documentation)
	}
	if w.Configuration == nil {
		w.Configuration = other.Configuration
	} else {
		mergeWorkspaceConfiguration(w.Configuration, other.Configuration)
	}
}

func buildIDMap(w, other *Workspace) map[string]string {
	if w.Model == nil || other.Model == nil {
		return nil
	}
	idmap := make(map[string]string)
	m, om := w.Model, other.Model
	for _, p := range m.People {
		for _, p2 := range om.People {
			if p.Name == p2.Name {
				idmap[p.ID] = p2.ID
				break
			}
		}
	}
	for _, s := range m.Systems {
		for _, s2 := range om.Systems {
			if s.Name == s2.Name {
				idmap[s.ID] = s2.ID
				for _, c := range s.Containers {
					for _, c2 := range s2.Containers {
						if c.Name == c2.Name {
							idmap[c.ID] = c2.ID
							for _, cmp := range c.Components {
								for _, cmp2 := range c2.Components {
									if cmp.Name == cmp2.Name {
										idmap[cmp.ID] = cmp2.ID
										break
									}
								}
							}
							break
						}
					}
				}
				break
			}
		}
	}
	for _, n := range m.DeploymentNodes {
		for _, n2 := range om.DeploymentNodes {
			if n.Name == n2.Name {
				buildDeploymentNodeIDMap(n, n2, idmap)
			}
		}
	}
	return idmap
}

func buildDeploymentNodeIDMap(n, other *DeploymentNode, idmap map[string]string) {
	idmap[n.ID] = other.ID
	for _, c := range n.Children {
		for _, c2 := range other.Children {
			if c.Name == c2.Name {
				buildDeploymentNodeIDMap(c, c2, idmap)
				break
			}
		}
	}
	for _, i := range n.InfrastructureNodes {
		for _, i2 := range other.InfrastructureNodes {
			if i.Name == i2.Name {
				idmap[i.ID] = i2.ID
				break
			}
		}
	}
	for _, ci := range n.ContainerInstances {
		for _, ci2 := range other.ContainerInstances {
			if ci.ContainerID == ci2.ContainerID && ci.InstanceID == ci2.InstanceID {
				idmap[ci.ID] = ci2.ID
				break
			}
		}
	}
}

func mergeModels(m, other *Model, idmap map[string]string) {
	if other == nil {
		return
	}
	if other.Enterprise != nil && other.Enterprise.Name != "" {
		m.Enterprise = other.Enterprise
	}
loopPeople:
	for _, p := range other.People {
		for _, p2 := range m.People {
			if p.Name == p2.Name {
				mergeElements(p2.Element, p.Element, idmap)
				p2.Location = p.Location
				continue loopPeople
			}
		}
		m.People = append(m.People, p)
	}
loopSystems:
	for _, s := range other.Systems {
		for _, s2 := range m.Systems {
			if s.Name == s2.Name {
				mergeElements(s2.Element, s.Element, idmap)
				s2.Location = s.Location
			loopContainers:
				for _, c := range s.Containers {
					for _, c2 := range s2.Containers {
						if c.Name == c2.Name {
							mergeContainers(c2, c, idmap)
						}
						continue loopContainers
					}
					s2.Containers = append(s2.Containers, c)
				}
				continue loopSystems
			}
		}
		m.Systems = append(m.Systems, s)
	}
loopDeploymentNodes:
	for _, n := range other.DeploymentNodes {
		for _, n2 := range m.DeploymentNodes {
			if n.Name == n2.Name {
				mergeDeploymentNodes(n2, n, idmap)
				continue loopDeploymentNodes
			}
		}
		m.DeploymentNodes = append(m.DeploymentNodes, n)
	}
}

func mergeContainers(c, other *Container, idmap map[string]string) {
	mergeElements(c.Element, other.Element, idmap)
loop:
	for _, cp := range other.Components {
		for _, cp2 := range c.Components {
			if cp.Name == cp2.Name {
				mergeElements(cp2.Element, cp.Element, idmap)
				continue loop
			}
		}
		c.Components = append(c.Components, cp)
	}
}

func mergeElements(e, other *Element, idmap map[string]string) {
	if other.Description != "" {
		e.Description = other.Description
	}
	if other.Technology != "" {
		e.Technology = other.Technology
	}
	e.Tags = other.Tags // do not merge to allow removing tags
	if other.URL != "" {
		e.URL = other.URL
	}
	if len(other.Properties) > 0 && e.Properties == nil {
		e.Properties = make(map[string]string)
	}
	for k, v := range other.Properties {
		e.Properties[k] = v
	}
loop:
	for _, rel := range other.Rels {
		for _, rel2 := range e.Rels {
			srcID := rel2.SourceID
			if mapped, ok := idmap[srcID]; ok {
				srcID = mapped
			}
			destID := rel2.DestinationID
			if mapped, ok := idmap[destID]; ok {
				destID = mapped
			}
			if rel.SourceID == srcID && rel.DestinationID == destID && rel.Description == rel2.Description {
				rel2.Tags = rel.Tags
				if rel.URL != "" {
					rel2.URL = rel.URL
				}
				if rel.Technology != "" {
					rel2.Technology = rel.Technology
				}
				rel2.InteractionStyle = rel.InteractionStyle
				if rel.LinkedRelationshipID != "" {
					rel2.LinkedRelationshipID = rel.LinkedRelationshipID
				}
				continue loop
			}
		}
		e.Rels = append(other.Rels, rel)
	}
}

func mergeDeploymentNodes(n, other *DeploymentNode, idmap map[string]string) {
	mergeElements(n.Element, other.Element, idmap)
	n.Environment = other.Environment
	if other.Instances != nil {
		n.Instances = other.Instances
	}
loopChildren:
	for _, c := range other.Children {
		for _, c2 := range n.Children {
			if c.Name == c2.Name {
				mergeDeploymentNodes(c2, c, idmap)
				continue loopChildren
			}
		}
		n.Children = append(n.Children, c)
	}
loopInfrastructureNodes:
	for _, i := range other.InfrastructureNodes {
		for _, i2 := range n.InfrastructureNodes {
			if i.Name == i2.Name {
				mergeElements(i2.Element, i.Element, idmap)
				i2.Environment = i.Environment
				continue loopInfrastructureNodes
			}
		}
		n.InfrastructureNodes = append(n.InfrastructureNodes, i)
	}
loopContainerInstances:
	for _, i := range other.ContainerInstances {
		for _, i2 := range n.ContainerInstances {
			if i.ContainerID == i2.ContainerID && i.InstanceID == i2.InstanceID {
				mergeElements(i2.Element, i.Element, idmap)
				i2.Environment = i.Environment
			loopHealthChecks:
				for _, hc := range i.HealthChecks {
					for _, hc2 := range i2.HealthChecks {
						if hc.Name == hc2.Name {
							hc2.URL = hc.URL
							hc2.Interval = hc.Interval
							hc2.Timeout = hc.Timeout
							hc2.Headers = hc.Headers
							continue loopHealthChecks
						}
					}
					i2.HealthChecks = append(i2.HealthChecks, hc)
				}
				continue loopContainerInstances
			}
		}
		n.ContainerInstances = append(n.ContainerInstances, i)
	}
}

func mergeViews(v, other *Views, idmap map[string]string) {
	if other == nil {
		return
	}
loopLandscapeViews:
	for _, lv := range other.LandscapeViews {
		for _, lv2 := range v.LandscapeViews {
			if lv.Key == lv2.Key {
				mergeViewProps(lv2.ViewProps, lv.ViewProps, idmap)
				lv2.EnterpriseBoundaryVisible = lv.EnterpriseBoundaryVisible
				continue loopLandscapeViews
			}
		}
		v.LandscapeViews = append(v.LandscapeViews, lv)
	}
loopContextViews:
	for _, cv := range other.ContextViews {
		for _, cv2 := range v.ContextViews {
			if cv.Key == cv2.Key {
				mergeViewProps(cv2.ViewProps, cv.ViewProps, idmap)
				cv2.EnterpriseBoundaryVisible = cv.EnterpriseBoundaryVisible
				cv2.SoftwareSystemID = cv.SoftwareSystemID
				continue loopContextViews
			}
		}
		v.ContextViews = append(v.ContextViews, cv)
	}
loopContainerViews:
	for _, cv := range other.ContainerViews {
		for _, cv2 := range v.ContainerViews {
			if cv.Key == cv2.Key {
				mergeViewProps(cv2.ViewProps, cv.ViewProps, idmap)
				cv2.SystemBoundariesVisible = cv.SystemBoundariesVisible
				cv2.SoftwareSystemID = cv.SoftwareSystemID
				continue loopContainerViews
			}
		}
		v.ContainerViews = append(v.ContainerViews, cv)
	}
loopComponentViews:
	for _, cv := range other.ComponentViews {
		for _, cv2 := range v.ComponentViews {
			if cv.Key == cv2.Key {
				mergeViewProps(cv2.ViewProps, cv.ViewProps, idmap)
				cv2.ContainerBoundariesVisible = cv.ContainerBoundariesVisible
				cv2.ContainerID = cv.ContainerID
				continue loopComponentViews
			}
		}
		v.ComponentViews = append(v.ComponentViews, cv)
	}
loopDynamicViews:
	for _, dv := range other.DynamicViews {
		for _, dv2 := range v.DynamicViews {
			if dv.Key == dv2.Key {
				mergeViewProps(dv2.ViewProps, dv.ViewProps, idmap)
				dv2.ElementID = dv.ElementID
				continue loopDynamicViews
			}
		}
		v.DynamicViews = append(v.DynamicViews, dv)
	}
loopDeploymentViews:
	for _, dv := range other.DeploymentViews {
		for _, dv2 := range v.DeploymentViews {
			if dv.Key == dv2.Key {
				mergeViewProps(dv2.ViewProps, dv.ViewProps, idmap)
				dv2.SoftwareSystemID = dv.SoftwareSystemID
				dv2.Environment = dv.Environment
				continue loopDeploymentViews
			}
		}
		v.DeploymentViews = append(v.DeploymentViews, dv)
	}
loopFilteredViews:
	for _, lv := range other.FilteredViews {
		for _, lv2 := range v.FilteredViews {
			if lv.Key == lv2.Key {
				if lv.Title != "" {
					lv2.Title = lv.Title
				}
				if lv.Description != "" {
					lv2.Description = lv.Description
				}
				lv2.BaseKey = lv.BaseKey
				lv2.Mode = lv.Mode
				lv2.Tags = lv.Tags // do not merge to allow removing tags
				continue loopFilteredViews
			}
		}
		v.FilteredViews = append(v.FilteredViews, lv)
	}
	if other.Configuration != nil {
		if v.Configuration != nil {
			mergeStyles(v.Configuration.Styles, other.Configuration.Styles)
			if other.Configuration.DefaultView != "" {
				v.Configuration.DefaultView = other.Configuration.DefaultView
			}
			if other.Configuration.Branding != nil {
				v.Configuration.Branding = other.Configuration.Branding
			}
			if term := other.Configuration.Terminology; term != nil {
				if v.Configuration.Terminology == nil {
					v.Configuration.Terminology = term
				} else {
					term2 := v.Configuration.Terminology
					if term.Enterprise != "" {
						term2.Enterprise = term.Enterprise
					}
					if term.Person != "" {
						term2.Person = term.Person
					}
					if term.SoftwareSystem != "" {
						term2.SoftwareSystem = term.SoftwareSystem
					}
					if term.Container != "" {
						term2.Container = term.Container
					}
					if term.Component != "" {
						term2.Component = term.Component
					}
					if term.Code != "" {
						term2.Code = term.Code
					}
					if term.DeploymentNode != "" {
						term2.DeploymentNode = term.DeploymentNode
					}
					if term.Relationship != "" {
						term2.Relationship = term.Relationship
					}
				}
			}
			v.Configuration.MetadataSymbols = other.Configuration.MetadataSymbols
		loopThemes:
			for _, t := range other.Configuration.Themes {
				for _, t2 := range v.Configuration.Themes {
					if t == t2 {
						continue loopThemes
					}
				}
				v.Configuration.Themes = append(v.Configuration.Themes, t)
			}
		} else {
			v.Configuration = other.Configuration
		}
	}
}

func mergeViewProps(p, other *ViewProps, idmap map[string]string) {
	if other.Title != "" {
		p.Title = other.Title
	}
	if other.Description != "" {
		p.Description = other.Description
	}
	p.PaperSize = other.PaperSize
	if other.Layout != nil {
		p.Layout = other.Layout
	}
loopElementViews:
	for _, ev := range other.ElementViews {
		for _, ev2 := range p.ElementViews {
			id := ev2.ID
			if mapped, ok := idmap[id]; ok {
				id = mapped
			}
			if ev.ID == id {
				if ev.X != nil {
					ev2.X = ev.X
				}
				if ev.Y != nil {
					ev2.Y = ev.Y
				}
				continue loopElementViews
			}
		}
		p.ElementViews = append(p.ElementViews, ev)
	}
loopRelationshipViews:
	for _, rv := range other.RelationshipViews {
		for _, rv2 := range p.RelationshipViews {
			id := rv2.ID
			if mapped, ok := idmap[id]; ok {
				id = mapped
			}
			if rv.ID == id {
				if rv.Description != "" {
					rv2.Description = rv.Description
				}
				if rv.Order != "" {
					rv2.Order = rv.Order
				}
				if rv.Vertices != nil {
					rv2.Vertices = rv.Vertices
				}
				rv2.Routing = rv.Routing
				if rv.Position != nil {
					rv2.Position = rv.Position
				}
				continue loopRelationshipViews
			}
		}
		p.RelationshipViews = append(p.RelationshipViews, rv)
	}
loopAnimations:
	for _, a := range other.Animations {
		for _, a2 := range p.Animations {
			if a.Order == a2.Order {
			loopElements:
				for _, e := range a.ElementIDs {
					for _, e2 := range a2.ElementIDs {
						if mapped, ok := idmap[e2]; ok {
							e2 = mapped
						}
						if e == e2 {
							continue loopElements
						}
					}
					a2.ElementIDs = append(a2.ElementIDs, e)
				}
			loopRelationships:
				for _, r := range a.Relationships {
					for _, r2 := range a2.Relationships {
						if mapped, ok := idmap[r2]; ok {
							r2 = mapped
						}
						if r == r2 {
							continue loopRelationships
						}
					}
					a2.Relationships = append(a2.Relationships, r)
				}
				continue loopAnimations
			}
		}
		p.Animations = append(p.Animations, a)
	}
}

func mergeStyles(s, other *Styles) {
loopElements:
	for _, es := range other.Elements {
		for _, es2 := range s.Elements {
			if es.Tag == es2.Tag {
				if es.Width != nil {
					es2.Width = es.Width
				}
				if es.Height != nil {
					es2.Height = es.Height
				}
				if es.Background != "" {
					es2.Background = es.Background
				}
				if es.Stroke != "" {
					es2.Stroke = es.Stroke
				}
				if es.Color != "" {
					es2.Color = es.Color
				}
				if es.FontSize != nil {
					es2.FontSize = es.FontSize
				}
				if es.Shape != ShapeUndefined {
					es2.Shape = es.Shape
				}
				if es.Icon != "" {
					es2.Icon = es.Icon
				}
				if es.Border != BorderUndefined {
					es2.Border = es.Border
				}
				if es.Opacity != nil {
					es2.Opacity = es.Opacity
				}
				if es.Metadata != nil {
					es2.Metadata = es.Metadata
				}
				if es.Description != nil {
					es2.Description = es.Description
				}
				continue loopElements
			}
		}
		s.Elements = append(s.Elements, es)
	}
loopRelationships:
	for _, rs := range other.Relationships {
		for _, rs2 := range s.Relationships {
			if rs.Tag == rs2.Tag {
				if rs.Thickness != nil {
					rs2.Thickness = rs.Thickness
				}
				if rs.Color != "" {
					rs2.Color = rs.Color
				}
				if rs.FontSize != nil {
					rs2.FontSize = rs.FontSize
				}
				if rs.Width != nil {
					rs2.Width = rs.Width
				}
				if rs.Dashed != nil {
					rs2.Dashed = rs.Dashed
				}
				if rs.Routing != RoutingUndefined {
					rs2.Routing = rs.Routing
				}
				if rs.Position != nil {
					rs2.Position = rs.Position
				}
				if rs.Opacity != nil {
					rs2.Opacity = rs.Opacity
				}
				continue loopRelationships
			}
		}
		s.Relationships = append(s.Relationships, rs)
	}
}

func mergeDocumentation(doc, other *Documentation) {
	if other == nil {
		return
	}
loopSections:
	for _, s := range doc.Sections {
		for _, s2 := range other.Sections {
			if s.Title == s2.Title {
				s2.Content = s.Content
				if s.Format != FormatUndefined {
					s2.Format = s.Format
				}
				s2.Order = s.Order
				if s.ElementID != "" {
					s2.ElementID = s.ElementID
				}
				continue loopSections
			}
		}
		doc.Sections = append(doc.Sections, s)
	}
loopDecisions:
	for _, dec := range doc.Decisions {
		for _, dec2 := range other.Decisions {
			if dec.ID == dec2.ID {
				if dec.Date != "" {
					dec2.Date = dec.Date
				}
				if dec.Decision != DecisionUndefined {
					dec2.Decision = dec.Decision
				}
				if dec.Title != "" {
					dec2.Title = dec.Title
				}
				dec2.Content = dec.Content
				if dec.Format != FormatUndefined {
					dec2.Format = dec.Format
				}
				if dec.ElementID != "" {
					dec2.ElementID = dec.ElementID
				}
				continue loopDecisions
			}
		}
		doc.Decisions = append(doc.Decisions, dec)
	}
loopImages:
	for _, img := range doc.Images {
		for _, img2 := range other.Images {
			if img.Name == img2.Name {
				img2.Content = img.Content
				if img.Type != "" {
					img2.Type = img.Type
				}
				continue loopImages
			}
		}
		doc.Images = append(doc.Images, img)
	}
	if other.Template != nil {
		doc.Template = other.Template
	}
}

func mergeWorkspaceConfiguration(cfg, other *WorkspaceConfiguration) {
	if other == nil {
		return
	}
loop:
	for _, u := range other.Users {
		for _, u2 := range cfg.Users {
			if u.Username == u2.Username {
				u2.Role = u.Role
				continue loop
			}
		}
		cfg.Users = append(cfg.Users, u)
	}
}
