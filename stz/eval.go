package stz

import (
	"goa.design/goa/v3/eval"
	"goa.design/model/expr"
	"goa.design/model/model"
)

// RunDSL runs the DSL defined in a global variable and returns the corresponding
// Structurize workspace.
func RunDSL() (*Workspace, error) {
	if err := eval.RunDSL(); err != nil {
		return nil, err
	}
	return WorkspaceFromDesign(expr.Root), nil
}

// WorkspaceFromDesign returns a Structurizr workspace initialized from the
// given design.
func WorkspaceFromDesign(d *expr.Design) *Workspace {
	design := model.ModelizeDesign(d)
	v := design.Views

	return &Workspace{
		Name:        d.Name,
		Description: d.Description,
		Version:     d.Version,
		Model:       design.Model,
		Views: &Views{
			LandscapeViews:  v.LandscapeViews,
			ContextViews:    v.ContextViews,
			ContainerViews:  v.ContainerViews,
			ComponentViews:  v.ComponentViews,
			DynamicViews:    v.DynamicViews,
			DeploymentViews: v.DeploymentViews,
			FilteredViews:   v.FilteredViews,
			Configuration:   &Configuration{Styles: v.Styles},
		},
	}
}
