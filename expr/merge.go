package expr

// MergeLayout merges the layout of elements and relationships in the views of
// remote into the views of w. Two elements correspond to each other if they
// have the same name in their scope (model for software systems, software
// system for containers and container for components). Two relationships
// correspond to each other if they have the same source, destination and
// description.
func (w *Workspace) MergeLayout(remote *Workspace) {
	if remote.Views == nil {
		return
	}
	if w.Views == nil {
		w.Views = remote.Views
		return
	}

	// Map w's element and relationship IDs to remote's.
	idmap := buildIDMap(remote, w)

	mergeViews(w.Views, remote.Views, idmap)
}

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

func mergeViews(local, remote *Views, idmap map[string]string) {
loop:
	for _, lv := range local.all() {
		for _, rv := range remote.all() {
			if lv.Key == rv.Key {
			loopElementViews:
				for _, lv := range lv.ElementViews {
					id := lv.ID
					if mapped, ok := idmap[id]; ok {
						id = mapped
					}
					for _, rv := range rv.ElementViews {
						if rv.ID == id {
							if lv.X == nil {
								lv.X = rv.X
							}
							if lv.Y == nil {
								lv.Y = rv.Y
							}
							continue loopElementViews
						}
					}
				}
			loopRelationshipViews:
				for _, lv := range lv.RelationshipViews {
					id := lv.ID
					if mapped, ok := idmap[id]; ok {
						id = mapped
					}
					for _, rv := range rv.RelationshipViews {
						if rv.ID == id {
							if lv.Vertices == nil {
								lv.Vertices = rv.Vertices
							}
							if lv.Routing == RoutingUndefined {
								lv.Routing = rv.Routing
							}
							if lv.Position == nil {
								lv.Position = rv.Position
							}
							continue loopRelationshipViews
						}
					}
				}
				continue loop
			}
		}
	}
}
