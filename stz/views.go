package stz

import (
	"goa.design/model/model"
)

type (
	// Views is the container for all views.
	Views struct {
		// LandscapeViewss describe the system landscape views.
		LandscapeViews []*model.LandscapeView `json:"systemLandscapeViews,omitempty"`
		// ContextViews lists the system context views.
		ContextViews []*model.ContextView `json:"systemContextViews,omitempty"`
		// ContainerViews lists the container views.
		ContainerViews []*model.ContainerView `json:"containerViews,omitempty"`
		// ComponentViews lists the component views.
		ComponentViews []*model.ComponentView `json:"componentViews,omitempty"`
		// DynamicViews lists the dynamic views.
		DynamicViews []*model.DynamicView `json:"dynamicViews,omitempty"`
		// DeploymentViews lists the deployment views.
		DeploymentViews []*model.DeploymentView `json:"deploymentViews,omitempty"`
		// FilteredViews lists the filtered views.
		FilteredViews []*model.FilteredView `json:"filteredViews,omitempty"`
		// Configuration contains view specific configuration information.
		Configuration *Configuration `json:"configuration,omitempty"`
	}
)
