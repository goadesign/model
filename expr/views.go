package expr

import (
	"fmt"
)

type (
	// Views defines one or more views.
	Views struct {
		// LandscapeViewss describe the system landscape views.
		LandscapeViews []*LandscapeView `json:"systemLandscapeViews,omitempty"`
		// ContextViews lists the system context views.
		ContextViews []*ContextView `json:"systemContextViews,omitempty"`
		// ContainerViews lists the container views.
		ContainerViews []*ContainerView `json:"containerViews,omitempty"`
		// ComponentViews lists the component views.
		ComponentViews []*ComponentView `json:"componentViews,omitempty"`
		// DynamicViews lists the dynamic views.
		DynamicViews []*DynamicView `json:"dynamicViews,omitempty"`
		// DeploymentViews lists the deployment views.
		DeploymentViews []*DeploymentView `json:"deploymentViews,omitempty"`
		// FilteredViews lists the filtered views.
		FilteredViews []*FilteredView `json:"filteredViews,omitempty"`
		// DSL to be run once all elements have been evaluated.
		DSL func() `json:"-"`
	}

	// LandscapeView describes a system landscape view.
	LandscapeView struct {
		View
		// EnterpriseBoundaryVisible specifies whether the enterprise boundary
		// (to differentiate internal elements from external elements) should be
		// visible on the resulting diagram.
		EnterpriseBoundaryVisible bool `json:"enterpriseBoundaryVisible"`
	}

	// ContextView describes a system context view.
	ContextView struct {
		View
		// EnterpriseBoundaryVisible specifies whether the enterprise boundary
		// (to differentiate internal elements from external elements) should be
		// visible on the resulting diagram.
		EnterpriseBoundaryVisible bool `json:"enterpriseBoundaryVisible"`
		// SoftwareSystemID is the ID of the software system this view with is
		// associated with.
		SoftwareSystemID string `json:"softwareSystemId"`
	}

	// ContainerView describes a container view for a specific software
	// system.
	ContainerView struct {
		View
		// Specifies whether software system boundaries should be visible for
		// "external" containers (those outside the software system in scope).
		ExternalSoftwareSystemBoundariesVisible bool `json:"externalSoftwareSystemBoundariesVisible"`
		// SoftwareSystemID is the ID of the software system this view with is
		// associated with.
		SoftwareSystemID string `json:"softwareSystemId"`
	}

	// ComponentView describes a component view for a specific container.
	ComponentView struct {
		View
		// Specifies whether container boundaries should be visible for
		// "external" containers (those outside the container in scope).
		ExternalContainerBoundariesVisible bool `json:"externalContainersBoundariesVisible"`
		// The ID of the container this view is associated with.
		ContainerID string `json:"containerID"`
	}

	// DynamicView describes a dynamic view for a specified scope.
	DynamicView struct {
		View
		// ElementID is the identifier of the element this view is associated with.
		ElementID string
	}

	// DeploymentView describes a deployment view.
	DeploymentView struct {
		View
		// SoftwareSystemID is the ID of the software system this view with is
		// associated with.
		SoftwareSystemID string `json:"softwareSystemId"`
		// The name of the environment that this deployment view is for (e.g.
		// "Development", "Live", etc).
		Environment string `json:"environment"`
	}

	// FilteredView describes a filtered view on top of a specified view.
	FilteredView struct {
		// Title of the view.
		Title string `json:"title,omitempty"`
		// Description of view.
		Description string `json:"description,omitempty"`
		// Key used to refer to the view.
		Key string `json:"key"`
		// BaseKey is the key of the view on which this filtered view is based.
		BaseKey string `json:"baseViewKey"`
		// Whether elements/relationships are being included ("Include") or
		// excluded ("Exclude") based upon the set of tags.
		Mode string `json:"mode"`
		// The set of tags to include/exclude elements/relationships when
		// rendering this filtered view.
		Tags []string `json:"tags"`
	}
)

// EvalName returns the generic expression name used in error messages.
func (v *Views) EvalName() string {
	return "views"
}

// EvalName returns the generic expression name used in error messages.
func (v *LandscapeView) EvalName() string {
	var suffix string
	if v.Key != "" {
		suffix = fmt.Sprintf(" with key %q", v.Key)
	}
	return fmt.Sprintf("system landscape view%s", suffix)
}

// EvalName returns the generic expression name used in error messages.
func (v *ContextView) EvalName() string {
	var suffix string
	if v.Key != "" {
		suffix = fmt.Sprintf(" with key %q", v.Key)
	}
	return fmt.Sprintf("system context view%s", suffix)
}

// EvalName returns the generic expression name used in error messages.
func (v *ContainerView) EvalName() string {
	var suffix string
	if v.Key != "" {
		suffix = fmt.Sprintf(" with key %q", v.Key)
	}
	return fmt.Sprintf("container view%s", suffix)
}

// EvalName returns the generic expression name used in error messages.
func (v *ComponentView) EvalName() string {
	var suffix string
	if v.Key != "" {
		suffix = fmt.Sprintf(" with key %q", v.Key)
	}
	return fmt.Sprintf("component view%s", suffix)
}

// EvalName returns the generic expression name used in error messages.
func (v *FilteredView) EvalName() string {
	var suffix string
	if v.Key != "" {
		suffix = fmt.Sprintf(" with base key %q", v.Key)
	}
	return fmt.Sprintf("filtered view%s", suffix)
}

// EvalName returns the generic expression name used in error messages.
func (v *DynamicView) EvalName() string {
	var suffix string
	if v.Key != "" {
		suffix = fmt.Sprintf(" with key %q", v.Key)
	}
	return fmt.Sprintf("dynamic view%s", suffix)
}

// EvalName returns the generic expression name used in error messages.
func (v *DeploymentView) EvalName() string {
	var suffix string
	if v.Key != "" {
		suffix = fmt.Sprintf(" with key %q", v.Key)
	}
	return fmt.Sprintf("deployment view%s", suffix)
}
