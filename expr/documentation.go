package expr

import "bytes"

type (
	// Documentation associated with software architecture model.
	Documentation struct {
		// Documentation sections.
		Sections []*DocumentationSection `json:"sections"`
		// ADR decisions.
		Decisions []*Decision `json:"decisions"`
		// Images used in documentation.
		Images []*Image `json:"images"`
		// Information about template used to render documentation.
		Template *DocumentationTemplateMetadata `json:"template"`
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
		ElementID string `json:"elementId"`
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
		ElementID string `json:"elementId"`
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
		Author string `json:"author"`
		// URL that points to more information about template.
		URL string `json:"url"`
	}

	// DocFormatKind is the enum used to represent documentation format.
	DocFormatKind int

	// DecisionStatusKind is the enum used to represent status of decision.
	DecisionStatusKind int
)

const (
	Markdown DocFormatKind = iota + 1
	ASCIIDoc
)

const (
	Proposed DecisionStatusKind = iota + 1
	Accepted
	Superseded
	Deprecated
	Rejected
)

// MarshalJSON replaces the constant value with the proper structurizr schema
// string value.
func (d DocFormatKind) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	switch d {
	case Markdown:
		buf.WriteString("Markown")
	case ASCIIDoc:
		buf.WriteString("AsciiDoc")
	}
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

// MarshalJSON replaces the constant value with the proper structurizr schema
// string value.
func (d DecisionStatusKind) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	switch d {
	case Proposed:
		buf.WriteString("Proposed")
	case Accepted:
		buf.WriteString("Accepted")
	case Superseded:
		buf.WriteString("Superseded")
	case Deprecated:
		buf.WriteString("Deprecated")
	case Rejected:
		buf.WriteString("Rejected")
	}
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}
