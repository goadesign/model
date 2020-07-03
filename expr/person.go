package expr

import "fmt"

type (
	// Person represents a person who uses a software system.
	Person struct {
		// ID of person.
		ID string `json:"id"`
		// Name of person.
		Name string `json:"name"`
		// Description of person if any.
		Description string `json:"description,omitempty"`
		// Tags attached to person as comma separated list if any.
		Tags string `json:"tags,omitempty"`
		// URL where more information about this person can be found.
		URL string `json:"url,omitempty"`
		// Location of element.
		Location LocationKind `json:"location"`
		// Rels is the set of relationships from this element to other elements.
		Rels []*Relationship `json:"relationships,omitempty"`
	}
)

// EvalName is the qualified name of the DSL expression e.g. "service
// bottle".
func (p *Person) EvalName() string {
	if p.Name == "" {
		return "unnamed person"
	}
	return fmt.Sprintf("person %#v", p.Name)
}
