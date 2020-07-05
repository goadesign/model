package expr

import (
	"bytes"
	"encoding/json"
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

	// View describes a view.
	View struct {
		// Title of the view.
		Title string `json:"title,omitempty"`
		// Description of view.
		Description string `json:"description,omitempty"`
		// Key used to refer to the view.
		Key string `json:"key"`
		// PaperSize is the paper size that should be used to render this view.
		PaperSize PaperSizeKind `json:"paperSize,omitempty"`
		// Layout describes the automatic layout mode for the diagram if
		// defined.
		Layout *Layout `json:"automaticLayout,omitempty"`
		// Elements list the elements included in the view.
		Elements []*ElementView `json:"elements,omitempty"`
		// Rels list the relationships included in the view.
		Relationships []*RelationshipView `json:"relationships,omitempty"`
		// AnimationSteps describes the animation steps if any.
		AnimationSteps []*AnimationStep `json:"animationSteps,omitempty"`
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

	// ElementView describes an instance of a model element (Person,
	// Software System, Container or Component) in a View.
	ElementView struct {
		// ID of element.
		ID string `json:"id"`
		// Horizontal position of element when rendered.
		X int `json:"x"`
		// Vertical position of element when rendered.
		Y int `json:"y"`
	}

	// RelationshipView describes an instance of a model relationship in a
	// view.
	RelationshipView struct {
		// ID of relationship.
		ID string `json:"id"`
		// Description of relationship used in dynamic views.
		Description string `json:"description,omitempty"`
		// Order of relationship in dynamic views.
		Order string `json:"order"`
		// Set of vertices used to render relationship
		Vertices []*Vertex `json:"vertices"`
		// Routing algorithm used to render relationship.
		Routing RoutingKind `json:"routing"`
		// Position of annotation along line; 0 (start) to 100 (end).
		Position int `json:"position"`
	}

	// Vertex describes the x and y coordinate of a bend in a line.
	Vertex struct {
		// Horizontal position of vertex when rendered.
		X int `json:"x"`
		// Vertical position of vertex when rendered.
		Y int `json:"y"`
	}

	// AnimationStep represents an animation step.
	AnimationStep struct {
		// Order of animation step.
		Order string `json:"order"`
		// Set of element IDs that should be included.
		Elements []string `json:"elements,omitempty"`
		// Set of relationship IDs tat should be included.
		Relationships []string `json:"relationships,omitempty"`
	}

	// Layout describes an automatic layout.
	Layout struct {
		// Algorithm rank direction.
		RankDirection RankDirectionKind `json:"rankDirection"`
		// RankSep defines the separation between ranks in pixels.
		RankSep int `json:"rankSeparation"`
		// NodeSep defines the separation between nodes in pixels.
		NodeSep int `json:"nodeSeparation"`
		// EdgeSep defines the separation between edges in pixels.
		EdgeSep int `json:"edgeSeparation"`
		// Render vertices if true.
		Vertices bool
	}

	// PaperSizeKind is the enum for possible paper kinds.
	PaperSizeKind int

	// RoutingKind is the enum for possible routing algorithms.
	RoutingKind int

	// RankDirectionKind is the enum for possible automatic layout rank
	// directions.
	RankDirectionKind int
)

const (
	SizeA6Portrait PaperSizeKind = iota + 1
	SizeA6Landscape
	SizeA5Portrait
	SizeA5Landscape
	SizeA4Portrait
	SizeA4Landscape
	SizeA3Portrait
	SizeA3Landscape
	SizeA2Portrait
	SizeA2Landscape
	SizeA1Portrait
	SizeA1Landscape
	SizeA0Portrait
	SizeA0Landscape
	SizeLetterPortrait
	SizeLetterLandscape
	SizeLegalPortrait
	SizeLegalLandscape
	SizeSlide4X3
	SizeSlide16X9
	SizeSlide16X10
)

const (
	RoutingDirect RoutingKind = iota + 1
	RoutingCurved
	RoutingOrthogonal
)

const (
	RankTopBottom RankDirectionKind = iota + 1
	RankBottomTop
	RankLeftRight
	RankRightLeft
)

// EvalName returns the generic expression name used in error messages.
func (v *Views) EvalName() string {
	return "views"
}

// Merge merges the given element views into the parent view.
func (v *View) Merge(evs []*ElementView) {
	for _, ev := range evs {
		var old *ElementView
		for _, e := range v.Elements {
			if e.ID == ev.ID {
				old = e
				break
			}
		}
		if old != nil {
			if ev.X > 0 {
				old.X = ev.X
			}
			if ev.Y > 0 {
				old.Y = ev.Y
			}
		} else {
			v.Elements = append(v.Elements, ev)
		}
	}
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

// MarshalJSON replaces the constant value with the proper structurizr schema
// string value.
func (p PaperSizeKind) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	switch p {
	case SizeA6Portrait:
		buf.WriteString("A6_Portrait")
	case SizeA6Landscape:
		buf.WriteString("A6_Landscape")
	case SizeA5Portrait:
		buf.WriteString("A5_Portrait")
	case SizeA5Landscape:
		buf.WriteString("A5_Landscape")
	case SizeA4Portrait:
		buf.WriteString("A4_Portrait")
	case SizeA4Landscape:
		buf.WriteString("A4_Landscape")
	case SizeA3Portrait:
		buf.WriteString("A3_Portrait")
	case SizeA3Landscape:
		buf.WriteString("A3_Landscape")
	case SizeA2Portrait:
		buf.WriteString("A2_Portrait")
	case SizeA2Landscape:
		buf.WriteString("A2_Landscape")
	case SizeA1Portrait:
		buf.WriteString("A1_Portrait")
	case SizeA1Landscape:
		buf.WriteString("A1_Landscape")
	case SizeA0Portrait:
		buf.WriteString("A0_Portrait")
	case SizeA0Landscape:
		buf.WriteString("A0_Landscape")
	case SizeLetterPortrait:
		buf.WriteString("Letter_Portrait")
	case SizeLetterLandscape:
		buf.WriteString("Letter_Landscape")
	case SizeLegalPortrait:
		buf.WriteString("Legal_Portrait")
	case SizeLegalLandscape:
		buf.WriteString("Legal_Landscape")
	case SizeSlide4X3:
		buf.WriteString("Slide_4_3")
	case SizeSlide16X9:
		buf.WriteString("Slide_16_9")
	case SizeSlide16X10:
		buf.WriteString("Slide_16_10")
	}
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

// UnmarshalJSON sets the constant from its JSON representation.
func (p *PaperSizeKind) UnmarshalJSON(data []byte) error {
	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	switch val {
	case "A6_Portrait":
		*p = SizeA6Portrait
	case "A6_Landscape":
		*p = SizeA6Landscape
	case "A5_Portrait":
		*p = SizeA5Portrait
	case "A5_Landscape":
		*p = SizeA5Landscape
	case "A4_Portrait":
		*p = SizeA4Portrait
	case "A4_Landscape":
		*p = SizeA4Landscape
	case "A3_Portrait":
		*p = SizeA3Portrait
	case "A3_Landscape":
		*p = SizeA3Landscape
	case "A2_Portrait":
		*p = SizeA2Portrait
	case "A2_Landscape":
		*p = SizeA2Landscape
	case "A1_Portrait":
		*p = SizeA1Portrait
	case "A1_Landscape":
		*p = SizeA1Landscape
	case "A0_Portrait":
		*p = SizeA0Portrait
	case "A0_Landscape":
		*p = SizeA0Landscape
	case "Letter_Portrait":
		*p = SizeLetterPortrait
	case "Letter_Landscape":
		*p = SizeLetterLandscape
	case "Legal_Portrait":
		*p = SizeLegalPortrait
	case "Legal_Landscape":
		*p = SizeLegalLandscape
	case "Slide_4_3":
		*p = SizeSlide4X3
	case "Slide_16_9":
		*p = SizeSlide16X9
	case "Slide_16_10":
		*p = SizeSlide16X10
	}
	return nil
}

// MarshalJSON replaces the constant value with the proper structurizr schema
// string value.
func (r RoutingKind) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	switch r {
	case RoutingDirect:
		buf.WriteString("Direct")
	case RoutingCurved:
		buf.WriteString("Curved")
	case RoutingOrthogonal:
		buf.WriteString("Orthogonal")
	}
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

// UnmarshalJSON sets the constant from its JSON representation.
func (r *RoutingKind) UnmarshalJSON(data []byte) error {
	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	switch val {
	case "Direct":
		*r = RoutingDirect
	case "Curved":
		*r = RoutingCurved
	case "Orthogonal":
		*r = RoutingOrthogonal
	}
	return nil
}

// MarshalJSON replaces the constant value with the proper structurizr schema
// string value.
func (r RankDirectionKind) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	switch r {
	case RankTopBottom:
		buf.WriteString("TopBottom")
	case RankBottomTop:
		buf.WriteString("BottomTop")
	case RankLeftRight:
		buf.WriteString("LeftRight")
	case RankRightLeft:
		buf.WriteString("RightLeft")
	}
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

// UnmarshalJSON sets the constant from its JSON representation.
func (r *RankDirectionKind) UnmarshalJSON(data []byte) error {
	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	switch val {
	case "TopBottom":
		*r = RankTopBottom
	case "BottomTop":
		*r = RankBottomTop
	case "LeftRight":
		*r = RankLeftRight
	case "RightLeft":
		*r = RankRightLeft
	}
	return nil
}
