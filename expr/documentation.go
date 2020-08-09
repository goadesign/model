package expr

import (
	"bytes"
	"encoding/json"
)

type (
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
