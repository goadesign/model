package stz

import (
	"bytes"
	"encoding/json"

	"goa.design/model/mdl"
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
		Model *mdl.Model `json:"model,omitempty"`
		// Views contains the views if any.
		Views *Views `json:"views,omitempty"`
		// Documentation associated with software architecture model.
		Documentation *Documentation `json:"documentation,omitempty"`
		// Configuration of workspace.
		Configuration *WorkspaceConfiguration `json:"configuration,omitempty"`
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

	// DocFormatKind is the enum used to represent documentation format.
	DocFormatKind int

	// DecisionStatusKind is the enum used to represent status of decision.
	DecisionStatusKind int
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
