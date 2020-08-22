package stz

import (
	"bytes"
	"encoding/json"
	"sort"

	"goa.design/model/design"
)

type (
	// Workspace describes a Structurizr service workspace.
	Workspace struct {
		// ID of workspace.
		ID int `json:"id,omitempty"`
		// Name of workspace.
		Name string `json:"name"`
		// Description of workspace if any.
		Description string `json:"description,omitempty"`
		// Version number for the workspace.
		Version string `json:"version,omitempty"`
		// Revision number, automatically generated.
		Revision int `json:"revision,omitempty"`
		// Thumbnail associated with the workspace; a Base64 encoded PNG file as a
		// data URI (data:image/png;base64).
		Thumbnail string `json:"thumbnail,omitempty"`
		// The last modified date, in ISO 8601 format (e.g. "2018-09-08T12:40:03Z").
		LastModifiedDate string `json:"lastModifiedDate,omitempty"`
		// A string identifying the user who last modified the workspace (e.g. an
		// e-mail address or username).
		LastModifiedUser string `json:"lastModifiedUser,omitempty"`
		//  A string identifying the agent that was last used to modify the workspace
		//  (e.g. "model-go/1.2.0").
		LastModifiedAgent string `json:"lastModifiedAgent,omitempty"`
		// Model is the software architecture model.
		Model *design.Model `json:"model,omitempty"`
		// Views contains the views if any.
		Views *Views `json:"views,omitempty"`
		// Documentation associated with software architecture model.
		Documentation *Documentation `json:"documentation,omitempty"`
		// Configuration of workspace.
		Configuration *WorkspaceConfiguration `json:"configuration,omitempty"`
	}

	// Views adds Structurizr specific configuration to the model views.
	Views struct {
		// LandscapeViewss describe the system landscape views.
		LandscapeViews []*design.LandscapeView `json:"systemLandscapeViews,omitempty"`
		// ContextViews lists the system context views.
		ContextViews []*design.ContextView `json:"systemContextViews,omitempty"`
		// ContainerViews lists the container views.
		ContainerViews []*design.ContainerView `json:"containerViews,omitempty"`
		// ComponentViews lists the component views.
		ComponentViews []*design.ComponentView `json:"componentViews,omitempty"`
		// DynamicViews lists the dynamic views.
		DynamicViews []*design.DynamicView `json:"dynamicViews,omitempty"`
		// DeploymentViews lists the deployment views.
		DeploymentViews []*design.DeploymentView `json:"deploymentViews,omitempty"`
		// FilteredViews lists the filtered views.
		FilteredViews []*design.FilteredView `json:"filteredViews,omitempty"`
		// Configuration contains view specific configuration information.
		Configuration *Configuration `json:"configuration,omitempty"`
	}

	// Configuration encapsulate Structurizr service specific view configuration information.
	Configuration struct {
		// Styles associated with views.
		Styles *design.Styles `json:"styles,omitempty"`
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

	// WorkspaceConfiguration describes the workspace configuration.
	WorkspaceConfiguration struct {
		// Users that have access to the workspace.
		Users []*User `json:"users"`
	}

	// User of Structurizr service.
	User struct {
		Username string `json:"username"`
		// Role of user, one of "ReadWrite" or "ReadOnly".
		Role string `json:"role"`
	}

	// Documentation associated with software architecture model.
	Documentation struct {
		// Documentation sections.
		Sections []*DocumentationSection `json:"sections,omitempty"`
		// ADR decisions.
		Decisions []*Decision `json:"decisions,omitempty"`
		// Images used in documentation.
		Images []*Image `json:"images,omitempty"`
		// Information about template used to render documentation.
		Template *DocumentationTemplateMetadata `json:"template,omitempty"`
	}

	// DocumentationSection corresponds to a documentation section.
	DocumentationSection struct {
		// Title (name/section heading) of section.
		Title string `json:"title"`
		// Markdown or AsciiDoc content of section.
		Content string `json:"string"`
		// Content format.
		Format DocFormatKind `json:"format"`
		// Order (index) of section in document.
		Order int `json:"order"`
		// ID of element (in model) that section applies to (optional).
		ElementID string `json:"elementId,omitempty"`
	}

	// Decision record (e.g. architecture decision record).
	Decision struct {
		// ID of decision.
		ID string `json:"id"`
		// Date of decision in ISO 8601 format.
		Date string `json:"date"`
		// Status of decision.
		Decision DecisionStatusKind `json:"decision"`
		// Title of decision
		Title string `json:"title"`
		// Markdown or AsciiDoc content of decision.
		Content string `json:"content"`
		// Content format.
		Format DocFormatKind `json:"format"`
		// ID of element (in model) that decision applies to (optional).
		ElementID string `json:"elementId,omitempty"`
	}

	// Image represents a Base64 encoded image (PNG/JPG/GIF).
	Image struct {
		// Name of image.
		Name string `json:"image"`
		// Base64 encoded content.
		Content string `json:"content"`
		// Image MIME type (e.g. "image/png")
		Type string `json:"type"`
	}

	// DocumentationTemplateMetadata provides information about a documentation
	// template used to create documentation.
	DocumentationTemplateMetadata struct {
		// Name of documentation template.
		Name string `json:"name"`
		// Name of author of documentation template.
		Author string `json:"author,omitempty"`
		// URL that points to more information about template.
		URL string `json:"url,omitempty"`
	}

	// SymbolKind is the enum used to represent symbols used to render metadata.
	SymbolKind int

	// DocFormatKind is the enum used to represent documentation format.
	DocFormatKind int

	// DecisionStatusKind is the enum used to represent status of decision.
	DecisionStatusKind int

	// for calling json.Marshal.
	_views Views
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

const (
	FormatUndefined DocFormatKind = iota
	FormatMarkdown
	FormatASCIIDoc
)

const (
	DecisionUndefined DecisionStatusKind = iota
	DecisionProposed
	DecisionAccepted
	DecisionSuperseded
	DecisionDeprecated
	DecisionRejected
)

// WorkspaceFromDesign creates a Workspace data structure compatible with the
// Structurizr service API from a design.
func WorkspaceFromDesign(d *design.Design) *Workspace {
	return &Workspace{
		Name:        d.Name,
		Description: d.Description,
		Version:     d.Version,
		Model:       d.Model,
		Views: &Views{
			LandscapeViews:  d.Views.LandscapeViews,
			ContextViews:    d.Views.ContextViews,
			ContainerViews:  d.Views.ContainerViews,
			ComponentViews:  d.Views.ComponentViews,
			DynamicViews:    d.Views.DynamicViews,
			DeploymentViews: d.Views.DeploymentViews,
			FilteredViews:   d.Views.FilteredViews,
			Configuration: &Configuration{
				Styles: d.Views.Styles,
			},
		},
	}
}

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

// MarshalJSON replaces the constant value with the proper string value.
func (d DocFormatKind) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	switch d {
	case FormatMarkdown:
		buf.WriteString("Markdown")
	case FormatASCIIDoc:
		buf.WriteString("AsciiDoc")
	}
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

// UnmarshalJSON sets the constant from its JSON representation.
func (d *DocFormatKind) UnmarshalJSON(data []byte) error {
	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	switch val {
	case "Markdown":
		*d = FormatMarkdown
	case "AsciiDoc":
		*d = FormatASCIIDoc
	}
	return nil
}

// MarshalJSON replaces the constant value with the proper string value.
func (d DecisionStatusKind) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	switch d {
	case DecisionProposed:
		buf.WriteString("Proposed")
	case DecisionAccepted:
		buf.WriteString("Accepted")
	case DecisionSuperseded:
		buf.WriteString("Superseded")
	case DecisionDeprecated:
		buf.WriteString("Deprecated")
	case DecisionRejected:
		buf.WriteString("Rejected")
	}
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

// UnmarshalJSON sets the constant from its JSON representation.
func (d *DecisionStatusKind) UnmarshalJSON(data []byte) error {
	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	switch val {
	case "Proposed":
		*d = DecisionProposed
	case "Accepted":
		*d = DecisionAccepted
	case "Superseded":
		*d = DecisionSuperseded
	case "Deprecated":
		*d = DecisionDeprecated
	case "Rejected":
		*d = DecisionRejected
	}
	return nil
}
