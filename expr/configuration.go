package expr

import "bytes"

type (
	// Configuration associated with a set of views.
	Configuration struct {
		// Styles associated with views.
		Styles *Styles `json:"styles"`
		// Key of view that was saved most recently.
		LastSavedView string `json:"lastSavedView"`
		// Key of view shown by default.
		DefaultView string `json:"defaultView"`
		// URL(s) of theme(s) used when rendering diagram.
		Themes []string `json:"themes"`
		// Branding used in views.
		Branding *Branding `json:"branding"`
		// Terminology used in workspace.
		Terminology *Terminology `json:"terminology"`
		// Type of symbols used when rendering metadata.
		MetadataSymbols SymbolKind `json:"metadataSymbols"`
	}

	// Styles describe styles associated with set of views.
	Styles struct {
		// Set of element styles.
		elements []*ElementStyle `json:"elements"`
		// Set of relationship styles.
		Relationships []*RelationshipStyle `json:"relationships`
	}

	// Branding is a wrapper for font and logo for diagram/documentation
	// branding purposes.
	Branding struct {
		// URL of PNG/JPG/GIF file, or Base64 data URI representation.
		Logo string `json:"logo"`
		// Font details.
		Font *Font `json:"font"`
	}

	// Terminology used on diagrams.
	Terminology struct {
		// Terminology used when rendering enterprise boundaries.
		Enterprise string `json:"enterprise"`
		// Terminology used when rendering people.
		Person string `json:"person"`
		// Terminology used when rendering software systems.
		SoftwareSystem string `json:"softwareSystem"`
		// Terminology used when rendering containers.
		Container string `json:"container"`
		// Terminology used when rendering components.
		Component string `json:"component"`
		// Terminology used when rendering code elements.
		Code string `json:"code"`
		// Terminology used when rendering deployment nodes.
		DeploymentNode string `json:"deploymentNode"`
		// Terminology used when rendering relationships.
		Relationship string `json:"relationship"`
	}

	// Font details including name and optional URL for web fonts.
	Font struct {
		// Name of font.
		Name string `json:"name"`
		// Web font URL.
		URL string `json:"url"`
	}

	// ElementStyle defines an element style.
	ElementStyle struct {
		// Tag to which this style applies.
		Tag string `json:"string"`
		// Width of element, in pixels.
		Width int `json:"width"`
		// Height of element, in pixels.
		Height int `json:"height"`
		// Background color of element as HTML RGB hex string (e.g. "#ffffff")
		Background string `json:"background"`
		// Stroke color of element as HTML RGB hex string (e.g. "#000000")
		Stroke string `json:"stroke"`
		// Foreground (text) color of element as HTML RGB hex string (e.g. "#ffffff")
		Color string `json:"color"`
		// Standard font size used to render text, in pixels.
		FontSize int `json:"fontSize"`
		// Shape used to render element.
		Shape ShapeKind `json:"shape"`
		// URL of PNG/JPG/GIF file or Base64 data URI representation.
		Icon string `json:"icon"`
		// Type of border used to render element.
		Border BorderKind `json:"border"`
		// Opacity used to render element; 0-100.
		Opacity int `json:"opacity"`
		// Whether element metadata should be shown.
		Metadata bool `json:"metadata"`
		// Whether element description should be shown.
		Description bool `json:"description"`
	}

	// RelationshipStyle defines a relationship style.
	RelationshipStyle struct {
		// Tag to which this style applies.
		Tag string `json:"tag"`
		// Thickness of line, in pixels.
		Thickness int `json:"thickness"`
		// Color of line as HTML RGB hex string (e.g. "#ffffff").
		Color string `json:"color"`
		// Standard font size used to render relationship annotation, in pixels.
		FontSize int `json:"fontSize"`
		// Width of relationship annotation, in pixels.
		Width int `json:"width"`
		// Whether line is rendered dashed or not.
		Dashed bool `json:"dashed"`
		// Routing algorithm used to render lines.
		Routing RoutingKind `json:"routing"`
		// Position of annotation along the line; 0 (start) to 100 (end).
		Position int `json:"position"`
		// Opacity used to render line; 0-100.
		Opacity int `json:"opacity"`
	}

	// SymbolKind is the enum used to represent symbols used to render metadata.
	SymbolKind int

	// ShapeKind is the enum used to represent shapes used to render elements.
	ShapeKind int

	// BorderKind is the enum used to represent element border styles.
	BorderKind int
)

const (
	SquareBrackets SymbolKind = iota + 1
	RoundBrackets
	CurlyBrackets
	AngleBrackets
	DoubleAngleBrackets
	None
)

const (
	Box ShapeKind = iota + 1
	RoundedBox
	Component
	Circle
	Ellipse
	Hexagon
	Folder
	Cylinder
	Pipe
	WebBrowser
	MobileDevicePortrait
	MobileDeviceLandscape
	PersonShape
	Robot
)

const (
	Solid BorderKind = iota + 1
	Dashed
	Dotted
)

// MarshalJSON replaces the constant value with the proper structurizr schema
// string value.
func (s SymbolKind) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	switch s {
	case SquareBrackets:
		buf.WriteString("SquareBrackets")
	case RoundBrackets:
		buf.WriteString("RoundBrackets")
	case CurlyBrackets:
		buf.WriteString("CurlyBrackets")
	case AngleBrackets:
		buf.WriteString("AngleBrackets")
	case DoubleAngleBrackets:
		buf.WriteString("DoubleAngleBrackets")
	case None:
		buf.WriteString("None")
	}
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

// MarshalJSON replaces the constant value with the proper structurizr schema
// string value.
func (s ShapeKind) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	switch s {
	case Box:
		buf.WriteString("Box")
	case RoundedBox:
		buf.WriteString("RoundedBox")
	case Component:
		buf.WriteString("Component")
	case Circle:
		buf.WriteString("Circle")
	case Ellipse:
		buf.WriteString("Ellipse")
	case Hexagon:
		buf.WriteString("Hexagon")
	case Folder:
		buf.WriteString("Folder")
	case Cylinder:
		buf.WriteString("Cylinder")
	case Pipe:
		buf.WriteString("Pipe")
	case WebBrowser:
		buf.WriteString("WebBrowser")
	case MobileDevicePortrait:
		buf.WriteString("MobileDevicePortrait")
	case MobileDeviceLandscape:
		buf.WriteString("MobileDeviceLandscape")
	case PersonShape:
		buf.WriteString("Person")
	case Robot:
		buf.WriteString("Robot")
	}
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

// MarshalJSON replaces the constant value with the proper structurizr schema
// string value.
func (b BorderKind) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	switch b {
	case Solid:
		buf.WriteString("Solid")
	case Dashed:
		buf.WriteString("Dashed")
	case Dotted:
		buf.WriteString("Dotted")
	}
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}
