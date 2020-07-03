package expr

import "bytes"

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
		Relationships []*Relationship `json:"relationships,omitempty"`
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
	A6_Portrait PaperSizeKind = iota + 1
	A6_Landscape
	A5_Portrait
	A5_Landscape
	A4_Portrait
	A4_Landscape
	A3_Portrait
	A3_Landscape
	A2_Portrait
	A2_Landscape
	A1_Portrait
	A1_Landscape
	A0_Portrait
	A0_Landscape
	Letter_Portrait
	Letter_Landscape
	Legal_Portrait
	Legal_Landscape
	Slide_4_3
	Slide_16_9
	Slide_16_10
)

const (
	Direct RoutingKind = iota + 1
	Curved
	Orthogonal
)

const (
	TopBottom RankDirectionKind = iota + 1
	BottomTop
	LeftRight
	RightLeft
)

// MarshalJSON replaces the constant value with the proper structurizr schema
// string value.
func (p PaperSizeKind) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	switch p {
	case A6_Portrait:
		buf.WriteString("A6_Portrait")
	case A6_Landscape:
		buf.WriteString("A6_Landscape")
	case A5_Portrait:
		buf.WriteString("A5_Portrait")
	case A5_Landscape:
		buf.WriteString("A5_Landscape")
	case A4_Portrait:
		buf.WriteString("A4_Portrait")
	case A4_Landscape:
		buf.WriteString("A4_Landscape")
	case A3_Portrait:
		buf.WriteString("A3_Portrait")
	case A3_Landscape:
		buf.WriteString("A3_Landscape")
	case A2_Portrait:
		buf.WriteString("A2_Portrait")
	case A2_Landscape:
		buf.WriteString("A2_Landscape")
	case A1_Portrait:
		buf.WriteString("A1_Portrait")
	case A1_Landscape:
		buf.WriteString("A1_Landscape")
	case A0_Portrait:
		buf.WriteString("A0_Portrait")
	case A0_Landscape:
		buf.WriteString("A0_Landscape")
	case Letter_Portrait:
		buf.WriteString("Letter_Portrait")
	case Letter_Landscape:
		buf.WriteString("Letter_Landscape")
	case Legal_Portrait:
		buf.WriteString("Legal_Portrait")
	case Legal_Landscape:
		buf.WriteString("Legal_Landscape")
	case Slide_4_3:
		buf.WriteString("Slide_4_3")
	case Slide_16_9:
		buf.WriteString("Slide_16_9")
	case Slide_16_10:
		buf.WriteString("Slide_16_10")
	}
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

// MarshalJSON replaces the constant value with the proper structurizr schema
// string value.
func (r RoutingKind) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	switch r {
	case Direct:
		buf.WriteString("Direct")
	case Curved:
		buf.WriteString("Curved")
	case Orthogonal:
		buf.WriteString("Orthogonal")
	}
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

// MarshalJSON replaces the constant value with the proper structurizr schema
// string value.
func (r RankDirectionKind) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	switch r {
	case TopBottom:
		buf.WriteString("TopBottom")
	case BottomTop:
		buf.WriteString("BottomTop")
	case LeftRight:
		buf.WriteString("LeftRight")
	case RightLeft:
		buf.WriteString("RightLeft")
	}
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}
