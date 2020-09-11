package stz

import (
	"bytes"
	"encoding/json"
)

type (
	// Configuration encapsulate Structurizr service specific view configuration information.
	Configuration struct {
		// Styles associated with views.
		Styles *Styles `json:"styles,omitempty"`
		// Key of view that was saved most recently.
		LastSavedView string `json:"lastSavedView,omitempty"`
		// Key of view shown by default.
		DefaultView string `json:"defaultView,omitempty"`
		// URL(s) of theme(s) used when rendering diagram.
		Themes []string `json:"themes,omitempty"`
		// Branding used in views.
		Branding *Branding `json:"branding,omitempty"`
		// Terminology used in workspace.
		Terminology *Terminology `json:"terminology,omitempty"`
		// Type of symbols used when rendering metadata.
		MetadataSymbols SymbolKind `json:"metadataSymbols,omitempty"`
	}

	// Branding is a wrapper for font and logo for diagram/documentation
	// branding purposes.
	Branding struct {
		// URL of PNG/JPG/GIF file, or Base64 data URI representation.
		Logo string `json:"logo,omitempty"`
		// Font details.
		Font *Font `json:"font,omitempty"`
	}

	// Terminology used on diagrams.
	Terminology struct {
		// Terminology used when rendering enterprise boundaries.
		Enterprise string `json:"enterprise,omitempty"`
		// Terminology used when rendering people.
		Person string `json:"person,omitempty"`
		// Terminology used when rendering software systems.
		SoftwareSystem string `json:"softwareSystem,omitempty"`
		// Terminology used when rendering containers.
		Container string `json:"container,omitempty"`
		// Terminology used when rendering components.
		Component string `json:"component,omitempty"`
		// Terminology used when rendering code elements.
		Code string `json:"code,omitempty"`
		// Terminology used when rendering deployment nodes.
		DeploymentNode string `json:"deploymentNode,omitempty"`
		// Terminology used when rendering relationships.
		Relationship string `json:"relationship,omitempty"`
	}

	// Font details including name and optional URL for web fonts.
	Font struct {
		// Name of font.
		Name string `json:"name,omitempty"`
		// Web font URL.
		URL string `json:"url,omitempty"`
	}

	// SymbolKind is the enum used to represent symbols used to render metadata.
	SymbolKind int
)

const (
	SymbolUndefined SymbolKind = iota
	SymbolSquareBrackets
	SymbolRoundBrackets
	SymbolCurlyBrackets
	SymbolAngleBrackets
	SymbolDoubleAngleBrackets
	SymbolNone
)

// MarshalJSON replaces the constant value with the proper string value.
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
