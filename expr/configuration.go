package expr

import (
	"bytes"
	"encoding/json"
)

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
	SymbolSquareBrackets SymbolKind = iota + 1
	SymbolRoundBrackets
	SymbolCurlyBrackets
	SymbolAngleBrackets
	SymbolDoubleAngleBrackets
	SymbolNone
)

const (
	ShapeBox ShapeKind = iota + 1
	ShapeRoundedBox
	ShapeComponent
	ShapeCircle
	ShapeEllipse
	ShapeHexagon
	ShapeFolder
	ShapeCylinder
	ShapePipe
	ShapeWebBrowser
	ShapeMobileDevicePortrait
	ShapeMobileDeviceLandscape
	ShapePerson
	ShapeRobot
)

const (
	BorderSolid BorderKind = iota + 1
	BorderDashed
	BorderDotted
)

// MarshalJSON replaces the constant value with the proper structurizr schema
// string value.
func (s SymbolKind) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	switch s {
	case SymbolSquareBrackets:
		buf.WriteString("SquareBrackets")
	case SymbolRoundBrackets:
		buf.WriteString("RoundBrackets")
	case SymbolCurlyBrackets:
		buf.WriteString("CurlyBrackets")
	case SymbolAngleBrackets:
		buf.WriteString("AngleBrackets")
	case SymbolDoubleAngleBrackets:
		buf.WriteString("DoubleAngleBrackets")
	case SymbolNone:
		buf.WriteString("None")
	}
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

// UnmarshalJSON sets the constant from its JSON representation.
func (s *SymbolKind) UnmarshalJSON(data []byte) error {
	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	switch val {
	case "SquareBrackets":
		*s = SymbolSquareBrackets
	case "RoundBrackets":
		*s = SymbolRoundBrackets
	case "CurlyBrackets":
		*s = SymbolCurlyBrackets
	case "AngleBrackets":
		*s = SymbolAngleBrackets
	case "DoubleAngleBrackets":
		*s = SymbolDoubleAngleBrackets
	case "None":
		*s = SymbolNone
	}
	return nil
}

// MarshalJSON replaces the constant value with the proper structurizr schema
// string value.
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

// MarshalJSON replaces the constant value with the proper structurizr schema
// string value.
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
