package stz

import (
	"goa.design/model/mdl"
)

type (
	// Views is the container for all views.
	Views struct {
		// LandscapeViewss describe the system landscape views.
		LandscapeViews []*mdl.LandscapeView `json:"systemLandscapeViews,omitempty"`
		// ContextViews lists the system context views.
		ContextViews []*mdl.ContextView `json:"systemContextViews,omitempty"`
		// ContainerViews lists the container views.
		ContainerViews []*mdl.ContainerView `json:"containerViews,omitempty"`
		// ComponentViews lists the component views.
		ComponentViews []*mdl.ComponentView `json:"componentViews,omitempty"`
		// DynamicViews lists the dynamic views.
		DynamicViews []*mdl.DynamicView `json:"dynamicViews,omitempty"`
		// DeploymentViews lists the deployment views.
		DeploymentViews []*mdl.DeploymentView `json:"deploymentViews,omitempty"`
		// FilteredViews lists the filtered views.
		FilteredViews []*mdl.FilteredView `json:"filteredViews,omitempty"`
		// Configuration contains view specific configuration information.
		Configuration *Configuration `json:"configuration,omitempty"`
	}
)
