package expr

type (
	// WorkspaceLayout describes the view layouts of a workspace. The layout
	// information includes element positions and relationship styles and
	// vertices and is indexed by view keys.
	WorkspaceLayout map[string]*ViewLayout

	// ViewLayout contains the layout information for a given view.
	ViewLayout struct {
		Elements      []*ElementView      `json:"elements,omitempty"`
		Relationships []*RelationshipView `json:"relationships,omitempty"`
	}
)

// Layout returns the workspace layout. It makes sure to only return relevant
// data. That is the entries in the layout all have at least one non-default
// field value (X or Y not 0 for elements, position not 0 or routing not
// undefined or vertices exist for relationships).
func (w *Workspace) Layout() WorkspaceLayout {
	views := w.Views
	if views == nil {
		return nil
	}
	layout := make(map[string]*ViewLayout)
	for _, v := range views.all() {
		var evs []*ElementView
		for _, ev := range v.ElementViews {
			if ev.X != nil && *ev.X != 0 || ev.Y != nil && *ev.Y != 0 {
				evs = append(evs, ev)
			}
		}
		var rvs []*RelationshipView
		for _, rv := range v.RelationshipViews {
			if rv.Position != nil || rv.Routing != RoutingUndefined || len(rv.Vertices) > 0 {
				rvs = append(rvs, rv)
			}
		}
		if len(evs) > 0 || len(rvs) > 0 {
			layout[v.Key] = &ViewLayout{
				Elements:      evs,
				Relationships: rvs,
			}
		}
	}
	return layout
}

// ApplyLayout merges the layout into the views of w.
func (w *Workspace) ApplyLayout(layout WorkspaceLayout) {
	for _, v := range w.Views.all() {
		if vl, ok := layout[v.Key]; ok {
			for _, el := range v.ElementViews {
				for _, vle := range vl.Elements {
					if el.ID == vle.ID {
						el.X = vle.X
						el.Y = vle.Y
						break
					}
				}
			}
			for _, rl := range v.RelationshipViews {
				for _, vlr := range vl.Relationships {
					if rl.ID == vlr.ID {
						rl.Vertices = vlr.Vertices
						rl.Routing = vlr.Routing
						rl.Position = vlr.Position
						break
					}
				}
			}
		}
	}
}

// MergeLayout merges the layout of elements and relationships in the views of
// remote into the views of w. The merge algorithm matches elements by name and
// relationships by matching source, destination and description (i.e. IDs don't
// have to be identical).
func (w *Workspace) MergeLayout(remote *Workspace) {
	if remote.Views == nil {
		return
	}
	if w.Views == nil {
		w.Views = remote.Views
		return
	}
	wl := remote.Layout()
	idmap := buildIDMap(remote, w)
	for _, m := range wl {
		for _, el := range m.Elements {
			if mapped, ok := idmap[el.ID]; ok {
				el.ID = mapped
			}
		}
		for _, rl := range m.Relationships {
			if mapped, ok := idmap[rl.ID]; ok {
				rl.ID = mapped
			}
		}
	}
	w.ApplyLayout(wl)
}

// buildIDMap returns a map that indexes the IDs of elements and relationships
// of remote with the IDs of matching elements and relationships of local. Two
// elements match if they have the same name in their scope (model for software
// systems, software system for containers and container for components). Two
// relationships match if they have matching source and destination and
// identical description.
func buildIDMap(remote, local *Workspace) map[string]string {
	if remote.Model == nil || local.Model == nil {
		return nil
	}
	idmap := make(map[string]string)
	rm, lm := remote.Model, local.Model

	// Compute element ID mappings.
	for _, rp := range rm.People {
		for _, lp := range lm.People {
			if rp.Name == lp.Name {
				idmap[lp.ID] = rp.ID
				break
			}
		}
	}
	for _, rs := range rm.Systems {
		for _, ls := range lm.Systems {
			if rs.Name == ls.Name {
				idmap[ls.ID] = rs.ID
				for _, rc := range rs.Containers {
					for _, lc := range ls.Containers {
						if rc.Name == lc.Name {
							idmap[lc.ID] = rc.ID
							for _, rcmp := range rc.Components {
								for _, lcmp := range lc.Components {
									if rcmp.Name == lcmp.Name {
										idmap[lcmp.ID] = rcmp.ID
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

	// Now compute relationship ID mappings using element mappings.
	for _, rp := range rm.People {
		for _, lp := range lm.People {
			if rp.Name == lp.Name {
				buildRelationshipIDMap(rp.Rels, lp.Rels, idmap)
				break
			}
		}
	}
	for _, rs := range rm.Systems {
		for _, ls := range lm.Systems {
			if rs.Name == ls.Name {
				buildRelationshipIDMap(rs.Rels, ls.Rels, idmap)
				for _, rc := range rs.Containers {
					for _, lc := range ls.Containers {
						if rc.Name == lc.Name {
							buildRelationshipIDMap(rc.Rels, lc.Rels, idmap)
							for _, rcmp := range rc.Components {
								for _, lcmp := range lc.Components {
									if rcmp.Name == lcmp.Name {
										buildRelationshipIDMap(rcmp.Rels, lcmp.Rels, idmap)
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
	return idmap
}

func buildRelationshipIDMap(remote, local []*Relationship, idmap map[string]string) {
	for _, lrel := range local {
		srcID := lrel.SourceID
		if mapped, ok := idmap[srcID]; ok {
			srcID = mapped
		}
		destID := lrel.DestinationID
		if mapped, ok := idmap[destID]; ok {
			destID = mapped
		}
		for _, rrel := range remote {
			if rrel.SourceID == srcID && rrel.DestinationID == destID && rrel.Description == lrel.Description {
				idmap[lrel.ID] = rrel.ID
				break
			}
		}
	}
}
