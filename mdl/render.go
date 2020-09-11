package mdl

import (
	"bytes"

	"goa.design/goa/v3/codegen"
	"goa.design/model/expr"
)

type (
	// RenderedView contains the information needed to create a static web page containing a view diagram.
	RenderedView struct {
		// Key is the key of the view.
		Key string
		// Title of view
		Title string
		// Version of design
		Version string
		// Description of view if any
		Description string
		// Mermaid contains the Mermaid source for the diagram.
		Mermaid string
		// Nodes contains additional information for each node rendered in the
		// diagram and is indexed by node ID (which corresponds to the ID of the
		// underlying element).
		Nodes map[string]*Node
	}

	// Node contains information about a single diagram node (element).
	Node struct {
		// ID is the node ID.
		ID string
		// URL is the URL specified in the element design.
		URL string
		// Properties are the properties specified in the element design.
		Properties map[string]string
		// ElementViewKey is the key of the first view in the design that is
		// scoped to the node element if any. This is useful to create links
		// to child element diagrams for example.
		ElementViewKey string
	}
)

// MermaidFiles returns codegen files that render Mermaid diagrams describing
// the views described in the given design. There is one file generated per view
// in the design.
func MermaidFiles(d *expr.Design) (files []*codegen.File) {
	views := d.Views
	if views == nil {
		return nil
	}
	for _, lv := range views.LandscapeViews {
		files = append(files, landscapeDiagram(lv))
	}
	for _, cv := range views.ContextViews {
		files = append(files, contextDiagram(cv))
	}
	for _, cv := range views.ContainerViews {
		files = append(files, containerDiagram(cv))
	}
	for _, cv := range views.ComponentViews {
		files = append(files, componentDiagram(cv))
	}
	for _, dv := range views.DynamicViews {
		files = append(files, dynamicDiagram(dv))
	}
	for _, dv := range views.DeploymentViews {
		files = append(files, deploymentDiagram(dv))
	}
	return
}

// Render renders the views of the given design.
func Render(d *expr.Design) []*RenderedView {
	views := d.Views
	if views == nil {
		return nil
	}
	var rvs []*RenderedView
	for _, lv := range views.LandscapeViews {
		rvs = append(rvs, render(landscapeDiagram(lv), lv, d))
	}
	for _, cv := range views.ContextViews {
		rvs = append(rvs, render(contextDiagram(cv), cv, d))
	}
	for _, cv := range views.ContainerViews {
		rvs = append(rvs, render(containerDiagram(cv), cv, d))
	}
	for _, cv := range views.ComponentViews {
		rvs = append(rvs, render(componentDiagram(cv), cv, d))
	}
	for _, dv := range views.DynamicViews {
		rvs = append(rvs, render(dynamicDiagram(dv), dv, d))
	}
	for _, dv := range views.DeploymentViews {
		rvs = append(rvs, render(deploymentDiagram(dv), dv, d))
	}
	return rvs
}

func render(f *codegen.File, view expr.View, d *expr.Design) *RenderedView {
	var buf bytes.Buffer
	for _, s := range f.SectionTemplates {
		if err := s.Write(&buf); err != nil {
			panic("render: " + err.Error()) // bug
		}
	}
	vp := view.Props()
	nodes := make(map[string]*Node, len(vp.ElementViews))
	for _, ev := range vp.ElementViews {
		var evk string
		switch e := expr.Registry[ev.Element.ID].(type) {
		case *expr.SoftwareSystem:
			for _, vv := range d.Views.ContainerViews {
				if vv.SoftwareSystemID == e.ID {
					evk = vv.Key
					break
				}
			}
		case *expr.Container:
			for _, vv := range d.Views.ComponentViews {
				if vv.ContainerID == e.ID {
					evk = vv.Key
					break
				}
			}
		}
		nodes[ev.Element.ID] = &Node{
			ID:             ev.Element.ID,
			URL:            ev.Element.URL,
			Properties:     ev.Element.Properties,
			ElementViewKey: evk,
		}
	}
	title := vp.Title
	if title == "" {
		title = "Diagram for " + vp.Key
	}
	return &RenderedView{
		Key:         vp.Key,
		Title:       title,
		Version:     d.Version,
		Description: vp.Description,
		Mermaid:     buf.String(),
		Nodes:       nodes,
	}
}
