package model

import (
	"bytes"
	"encoding/json"
	"sort"
)

type (
	// Views is the container for all views.
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
		// Styles associated with views.
		Styles *Styles `json:"styles,omitempty"`
	}

	// LandscapeView describes a system landscape view.
	LandscapeView struct {
		*ViewProps
		// EnterpriseBoundaryVisible specifies whether the enterprise boundary
		// (to differentiate internal elements from external elements) should be
		// visible on the resulting diagram.
		EnterpriseBoundaryVisible *bool `json:"enterpriseBoundaryVisible,omitempty"`
	}

	// ContextView describes a system context view.
	ContextView struct {
		*ViewProps
		// EnterpriseBoundaryVisible specifies whether the enterprise boundary
		// (to differentiate internal elements from external elements) should be
		// visible on the resulting diagram.
		EnterpriseBoundaryVisible *bool `json:"enterpriseBoundaryVisible,omitempty"`
		// SoftwareSystemID is the ID of the software system this view with is
		// associated with.
		SoftwareSystemID string `json:"softwareSystemId"`
	}

	// ContainerView describes a container view for a specific software
	// system.
	ContainerView struct {
		*ViewProps
		// Specifies whether software system boundaries should be visible for
		// "external" containers (those outside the software system in scope).
		SystemBoundariesVisible *bool `json:"externalSoftwareSystemBoundariesVisible,omitempty"`
		// SoftwareSystemID is the ID of the software system this view with is
		// associated with.
		SoftwareSystemID string `json:"softwareSystemId"`
	}

	// ComponentView describes a component view for a specific container.
	ComponentView struct {
		*ViewProps
		// Specifies whether container boundaries should be visible for
		// "external" containers (those outside the container in scope).
		ContainerBoundariesVisible *bool `json:"externalContainersBoundariesVisible,omitempty"`
		// The ID of the container this view is associated with.
		ContainerID string `json:"containerId"`
	}

	// DynamicView describes a dynamic view for a specified scope.
	DynamicView struct {
		*ViewProps
		// ElementID is the identifier of the element this view is associated with.
		ElementID string `json:"elementId"`
	}

	// DeploymentView describes a deployment view.
	DeploymentView struct {
		*ViewProps
		// SoftwareSystemID is the ID of the software system this view with is
		// associated with if any.
		SoftwareSystemID string `json:"softwareSystemId,omitempty"`
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
		Tags []string `json:"tags,omitempty"`
	}

	// ViewProps contains common properties for all views.
	ViewProps struct {
		// Title of the view
		Title string `json:"title,omitempty"`
		// Description of view
		Description string `json:"description,omitempty"`
		// Key used to identify the view
		Key string `json:"key"`
		// PaperSize is the paper size that should be used to render this view.
		PaperSize PaperSizeKind `json:"paperSize,omitempty"`
		// AutoLayout describes the automatic layout mode for the diagram if
		// defined.
		AutoLayout *AutoLayout `json:"automaticLayout,omitempty"`
		// ElementViews lists the elements included in the view.
		ElementViews []*ElementView `json:"elements,omitempty"`
		// RelationshipViews lists the relationships included in the view.
		RelationshipViews []*RelationshipView `json:"relationships,omitempty"`
		// Animations describes the animation steps if any.
		Animations []*AnimationStep `json:"animations,omitempty"`
	}

	// ElementView describes an instance of a model element (Person,
	// Software System, Container or Component) in a View.
	ElementView struct {
		// ID of element.
		ID string `json:"id"`
		// Horizontal position of element when rendered
		X *int `json:"x,omitempty"`
		// Vertical position of element when rendered.
		Y *int `json:"y,omitempty"`
	}

	// RelationshipView describes an instance of a model relationship in a
	// view.
	RelationshipView struct {
		// ID of relationship.
		ID string `json:"id"`
		// Description of relationship used in dynamic views.
		Description string `json:"description,omitempty"`
		// Order of relationship in dynamic views.
		Order string `json:"order,omitempty"`
		// Set of vertices used to render relationship
		Vertices []*Vertex `json:"vertices,omitempty"`
		// Routing algorithm used to render relationship.
		Routing RoutingKind `json:"routing,omitempty"`
		// Position of annotation along line; 0 (start) to 100 (end).
		Position *int `json:"position,omitempty"`
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
		Order int `json:"order"`
		// Set of element IDs that should be included.
		Elements []string `json:"elements,omitempty"`
		// Set of relationship IDs tat should be included.
		Relationships []string `json:"relationships,omitempty"`
	}

	// AutoLayout describes an automatic layout.
	AutoLayout struct {
		// Implementation used for automatic layouting of the view
		Implementation ImplementationKind `json:"implementation,omitempty"`
		// Algorithm rank direction.
		RankDirection RankDirectionKind `json:"rankDirection,omitempty"`
		// RankSep defines the separation between ranks in pixels.
		RankSep *int `json:"rankSeparation,omitempty"`
		// NodeSep defines the separation between nodes in pixels.
		NodeSep *int `json:"nodeSeparation,omitempty"`
		// EdgeSep defines the separation between edges in pixels.
		EdgeSep *int `json:"edgeSeparation,omitempty"`
		// Render vertices if true.
		Vertices *bool `json:"vertices,omitempty"`
	}

	// Styles describe styles associated with set of views.
	Styles struct {
		// Elements is the set of element styles.
		Elements []*ElementStyle `json:"elements,omitempty"`
		// Relationships is the set of relationship styles.
		Relationships []*RelationshipStyle `json:"relationships,omitempty"`
	}

	// ElementStyle defines an element style.
	ElementStyle struct {
		// Tag to which this style applies.
		Tag string `json:"tag,omitempty"`
		// Width of element, in pixels.
		Width *int `json:"width,omitempty"`
		// Height of element, in pixels.
		Height *int `json:"height,omitempty"`
		// Background color of element as HTML RGB hex string (e.g. "#ffffff")
		Background string `json:"background,omitempty"`
		// Stroke color of element as HTML RGB hex string (e.g. "#000000")
		Stroke string `json:"stroke,omitempty"`
		// Foreground (text) color of element as HTML RGB hex string (e.g. "#ffffff")
		Color string `json:"color,omitempty"`
		// Standard font size used to render text, in pixels.
		FontSize *int `json:"fontSize,omitempty"`
		// Shape used to render element.
		Shape ShapeKind `json:"shape,omitempty"`
		// URL of PNG/JPG/GIF file or Base64 data URI representation.
		Icon string `json:"icon,omitempty"`
		// Type of border used to render element.
		Border BorderKind `json:"border,omitempty"`
		// Opacity used to render element; 0-100.
		Opacity *int `json:"opacity,omitempty"`
		// Whether element metadata should be shown.
		Metadata *bool `json:"metadata,omitempty"`
		// Whether element description should be shown.
		Description *bool `json:"description,omitempty"`
	}

	// RelationshipStyle defines a relationship style.
	RelationshipStyle struct {
		// Tag to which this style applies.
		Tag string `json:"tag,omitempty"`
		// Thickness of line, in pixels.
		Thickness *int `json:"thickness,omitempty"`
		// Color of line as HTML RGB hex string (e.g. "#ffffff").
		Color string `json:"color,omitempty"`
		// Standard font size used to render relationship annotation, in pixels.
		FontSize *int `json:"fontSize,omitempty"`
		// Width of relationship annotation, in pixels.
		Width *int `json:"width,omitempty"`
		// Whether line is rendered dashed or not.
		Dashed *bool `json:"dashed,omitempty"`
		// Routing algorithm used to render lines.
		Routing RoutingKind `json:"routing,omitempty"`
		// Position of annotation along the line; 0 (start) to 100 (end).
		Position *int `json:"position,omitempty"`
		// Opacity used to render line; 0-100.
		Opacity *int `json:"opacity,omitempty"`
	}

	// PaperSizeKind is the enum for possible paper kinds.
	PaperSizeKind int

	// RoutingKind is the enum for possible routing algorithms.
	RoutingKind int

	// ImplementationKind is the enum for possible automatic layout implementations
	ImplementationKind int

	// RankDirectionKind is the enum for possible automatic layout rank
	// directions.
	RankDirectionKind int

	// ShapeKind is the enum used to represent shapes used to render elements.
	ShapeKind int

	// BorderKind is the enum used to represent element border styles.
	BorderKind int

	// for calling json.Marshal.
	_views          Views
	_landscapeView  LandscapeView
	_contextView    ContextView
	_containerView  ContainerView
	_componentView  ComponentView
	_dynamicView    DynamicView
	_deploymentView DeploymentView
	_filteredView   FilteredView
	_styles         Styles
)

const (
	SizeUndefined PaperSizeKind = iota
	SizeA0Landscape
	SizeA0Portrait
	SizeA1Landscape
	SizeA1Portrait
	SizeA2Landscape
	SizeA2Portrait
	SizeA3Landscape
	SizeA3Portrait
	SizeA4Landscape
	SizeA4Portrait
	SizeA5Landscape
	SizeA5Portrait
	SizeA6Landscape
	SizeA6Portrait
	SizeLegalLandscape
	SizeLegalPortrait
	SizeLetterLandscape
	SizeLetterPortrait
	SizeSlide16X10
	SizeSlide16X9
	SizeSlide4X3
)

const (
	RoutingUndefined RoutingKind = iota
	RoutingDirect
	RoutingOrthogonal
	RoutingCurved
)

const (
	ImplementationUndefined ImplementationKind = iota
	ImplementationGraphviz
	ImplementationDagre
)

const (
	RankUndefined RankDirectionKind = iota
	RankTopBottom
	RankBottomTop
	RankLeftRight
	RankRightLeft
)

const (
	ShapeUndefined ShapeKind = iota
	ShapeBox
	ShapeCircle
	ShapeCylinder
	ShapeEllipse
	ShapeHexagon
	ShapeRoundedBox
	ShapeComponent
	ShapeFolder
	ShapeMobileDeviceLandscape
	ShapeMobileDevicePortrait
	ShapePerson
	ShapePipe
	ShapeRobot
	ShapeWebBrowser
)

const (
	BorderUndefined BorderKind = iota
	BorderSolid
	BorderDashed
	BorderDotted
)

// MarshalJSON guarantees the order of elements in generated JSON arrays that
// correspond to sets.
func (v *Views) MarshalJSON() ([]byte, error) {
	sort.Slice(v.LandscapeViews, func(i, j int) bool { return v.LandscapeViews[i].Key < v.LandscapeViews[j].Key })
	sort.Slice(v.ContextViews, func(i, j int) bool { return v.ContextViews[i].Key < v.ContextViews[j].Key })
	sort.Slice(v.ContainerViews, func(i, j int) bool { return v.ContainerViews[i].Key < v.ContainerViews[j].Key })
	sort.Slice(v.ComponentViews, func(i, j int) bool { return v.ComponentViews[i].Key < v.ComponentViews[j].Key })
	sort.Slice(v.DynamicViews, func(i, j int) bool { return v.DynamicViews[i].Key < v.DynamicViews[j].Key })
	sort.Slice(v.DeploymentViews, func(i, j int) bool { return v.DeploymentViews[i].Key < v.DeploymentViews[j].Key })
	sort.Slice(v.FilteredViews, func(i, j int) bool { return v.FilteredViews[i].Key < v.FilteredViews[j].Key })
	vv := _views(*v)
	return json.Marshal(&vv)
}

// MarshalJSON guarantees the order of elements in generated JSON arrays that
// correspond to sets.
func (v *LandscapeView) MarshalJSON() ([]byte, error) {
	sortViews(v.ViewProps)
	vv := _landscapeView(*v)
	return json.Marshal(&vv)
}

// MarshalJSON guarantees the order of elements in generated JSON arrays that
// correspond to sets.
func (v *ContextView) MarshalJSON() ([]byte, error) {
	sortViews(v.ViewProps)
	vv := _contextView(*v)
	return json.Marshal(&vv)
}

// MarshalJSON guarantees the order of elements in generated JSON arrays that
// correspond to sets.
func (v *ContainerView) MarshalJSON() ([]byte, error) {
	sortViews(v.ViewProps)
	vv := _containerView(*v)
	return json.Marshal(&vv)
}

// MarshalJSON guarantees the order of elements in generated JSON arrays that
// correspond to sets.
func (v *ComponentView) MarshalJSON() ([]byte, error) {
	sortViews(v.ViewProps)
	vv := _componentView(*v)
	return json.Marshal(&vv)
}

// MarshalJSON guarantees the order of elements in generated JSON arrays that
// correspond to sets.
func (v *DynamicView) MarshalJSON() ([]byte, error) {
	sortViews(v.ViewProps)
	vv := _dynamicView(*v)
	return json.Marshal(&vv)
}

// MarshalJSON guarantees the order of elements in generated JSON arrays that
// correspond to sets.
func (v *DeploymentView) MarshalJSON() ([]byte, error) {
	sortViews(v.ViewProps)
	vv := _deploymentView(*v)
	return json.Marshal(&vv)
}

// MarshalJSON guarantees the order of elements in generated JSON arrays that
// correspond to sets.
func (v *FilteredView) MarshalJSON() ([]byte, error) {
	sort.Strings(v.Tags)
	vv := _filteredView(*v)
	return json.Marshal(&vv)
}

// MarshalJSON guarantees the order of elements in generated JSON arrays that
// correspond to sets.
func (s *Styles) MarshalJSON() ([]byte, error) {
	sort.Slice(s.Elements, func(i, j int) bool { return s.Elements[i].Tag < s.Elements[j].Tag })
	sort.Slice(s.Relationships, func(i, j int) bool { return s.Relationships[i].Tag < s.Relationships[j].Tag })
	ss := _styles(*s)
	return json.Marshal(&ss)
}

// Name is the name of the paper size kind as it appears in the DSL.
func (p PaperSizeKind) Name() string {
	switch p {
	case SizeA0Landscape:
		return "SizeA0Landscape"
	case SizeA0Portrait:
		return "SizeA0Portrait"
	case SizeA1Landscape:
		return "SizeA1Landscape"
	case SizeA1Portrait:
		return "SizeA1Portrait"
	case SizeA2Landscape:
		return "SizeA2Landscape"
	case SizeA2Portrait:
		return "SizeA2Portrait"
	case SizeA3Landscape:
		return "SizeA3Landscape"
	case SizeA3Portrait:
		return "SizeA3Portrait"
	case SizeA4Landscape:
		return "SizeA4Landscape"
	case SizeA4Portrait:
		return "SizeA4Portrait"
	case SizeA5Landscape:
		return "SizeA5Landscape"
	case SizeA5Portrait:
		return "SizeA5Portrait"
	case SizeA6Landscape:
		return "SizeA6Landscape"
	case SizeA6Portrait:
		return "SizeA6Portrait"
	case SizeLegalLandscape:
		return "SizeLegalLandscape"
	case SizeLegalPortrait:
		return "SizeLegalPortrait"
	case SizeLetterLandscape:
		return "SizeLetterLandscape"
	case SizeLetterPortrait:
		return "SizeLetterPortrait"
	case SizeSlide16X10:
		return "SizeSlide16X10"
	case SizeSlide16X9:
		return "SizeSlide16X9"
	case SizeSlide4X3:
		return "SizeSlide4X3"
	default:
		return "SizeUndefined"
	}
}

// MarshalJSON replaces the constant value with the proper string value.
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
	default:
		*p = SizeUndefined
	}
	return nil
}

// Name returns the name of the implementation used in the automatic layout of the view.
func (i ImplementationKind) Name() string {
	switch i {
	case ImplementationGraphviz:
		return "Graphviz"
	case ImplementationDagre:
		return "Dagre"
	default:
		return "ImplementationUndefined"
	}
}

// Name returns the name of the rank direction is specified in the DSL.
func (r RankDirectionKind) Name() string {
	switch r {
	case RankTopBottom:
		return "RankTopBottom"
	case RankBottomTop:
		return "RankBottomTop"
	case RankLeftRight:
		return "RankLeftRight"
	case RankRightLeft:
		return "RankRightLeft"
	default:
		return "RankUndefined"
	}
}

// MarshalJSON replaces the constant value with the proper string value.
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

// MarshalJSON replaces the constant value with the proper string value.
func (i ImplementationKind) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	switch i {
	case ImplementationGraphviz:
		buf.WriteString("Graphviz")
	case ImplementationDagre:
		buf.WriteString("Dagre")
	}
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

// UnmarshalJSON sets the constant from its JSON representation.
func (i *ImplementationKind) UnmarshalJSON(data []byte) error {
	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	switch val {
	case "Graphviz":
		*i = ImplementationGraphviz
	case "Dagre":
		*i = ImplementationDagre
	}
	return nil
}

// MarshalJSON replaces the constant value with the proper string value.
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

// MarshalJSON replaces the constant value with the proper string value.
func (s ShapeKind) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	switch s {
	case ShapeBox:
		buf.WriteString("Box")
	case ShapeRoundedBox:
		buf.WriteString("RoundedBox")
	case ShapeComponent:
		buf.WriteString("Component")
	case ShapeCircle:
		buf.WriteString("Circle")
	case ShapeEllipse:
		buf.WriteString("Ellipse")
	case ShapeHexagon:
		buf.WriteString("Hexagon")
	case ShapeFolder:
		buf.WriteString("Folder")
	case ShapeCylinder:
		buf.WriteString("Cylinder")
	case ShapePipe:
		buf.WriteString("Pipe")
	case ShapeWebBrowser:
		buf.WriteString("WebBrowser")
	case ShapeMobileDevicePortrait:
		buf.WriteString("MobileDevicePortrait")
	case ShapeMobileDeviceLandscape:
		buf.WriteString("MobileDeviceLandscape")
	case ShapePerson:
		buf.WriteString("Person")
	case ShapeRobot:
		buf.WriteString("Robot")
	}
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

// UnmarshalJSON sets the constant from its JSON representation.
func (s *ShapeKind) UnmarshalJSON(data []byte) error {
	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	switch val {
	case "Box":
		*s = ShapeBox
	case "RoundedBox":
		*s = ShapeRoundedBox
	case "Component":
		*s = ShapeComponent
	case "Circle":
		*s = ShapeCircle
	case "Ellipse":
		*s = ShapeEllipse
	case "Hexagon":
		*s = ShapeHexagon
	case "Folder":
		*s = ShapeFolder
	case "Cylinder":
		*s = ShapeCylinder
	case "Pipe":
		*s = ShapePipe
	case "WebBrowser":
		*s = ShapeWebBrowser
	case "MobileDevicePortrait":
		*s = ShapeMobileDevicePortrait
	case "MobileDeviceLandscape":
		*s = ShapeMobileDeviceLandscape
	case "Person":
		*s = ShapePerson
	case "Robot":
		*s = ShapeRobot
	}
	return nil
}

// MarshalJSON replaces the constant value with the proper string value.
func (b BorderKind) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	switch b {
	case BorderSolid:
		buf.WriteString("Solid")
	case BorderDashed:
		buf.WriteString("Dashed")
	case BorderDotted:
		buf.WriteString("Dotted")
	}
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

// UnmarshalJSON sets the constant from its JSON representation.
func (b *BorderKind) UnmarshalJSON(data []byte) error {
	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	switch val {
	case "Solid":
		*b = BorderSolid
	case "Dashed":
		*b = BorderDashed
	case "Dotted":
		*b = BorderDotted
	}
	return nil
}

// Sort guarantees the order of elements in generated JSON arrays that
// correspond to sets.
func sortViews(v *ViewProps) {
	sort.Slice(v.ElementViews, func(i, j int) bool { return v.ElementViews[i].ID < v.ElementViews[j].ID })
	sort.Slice(v.RelationshipViews, func(i, j int) bool { return v.RelationshipViews[i].ID < v.RelationshipViews[j].ID })
	for _, s := range v.Animations {
		sort.Strings(s.Elements)
		sort.Strings(s.Relationships)
	}
}
