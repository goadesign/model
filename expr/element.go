package expr

import (
	"bytes"
	"encoding/json"
)

type (
	// Element describes an element.
	Element struct {
		// ID of element.
		ID string `json:"id"`
		// Name of element.
		Name string `json:"name"`
		// Description of element if any.
		Description string `json:"description,omitempty"`
		// Technology used by element if any - not applicable to Person.
		Technology string `json:"technology,omitempty"`
		// Tags attached to element as comma separated list if any.
		Tags string `json:"tags,omitempty"`
		// URL where more information about this element can be found.
		URL string `json:"url,omitempty"`
		// Location of element.
		Location LocationKind `json:"location"`
		// Set of arbitrary name-value properties (shown in diagram tooltips).
		Properties map[string]string `json:"properties"`
		// Rels is the set of relationships from this element to other elements.
		Rels []*Relationship `json:"relationships,omitempty"`
	}

	// ElementHolder provides access to the underlying element.
	ElementHolder interface {
		GetElement() *Element
	}

	// LocationKind is the enum for possible locations.
	LocationKind int
)

const (
	// LocationUndefined means no location specified in design.
	LocationUndefined LocationKind = iota
	// LocationInternal defines an element internal to the enterprise.
	LocationInternal
	// LocationExternal defines an element external to the enterprise.
	LocationExternal
)

// GetElement returns the underlying element.
func (e *Element) GetElement() *Element { return e }

// RelatedPeople returns all people the element has a relationship with
// (either as source or as destination).
func (e *Element) RelatedPeople() (res People) {
	add := func(p *Person) {
		for _, ep := range res {
			if ep.ID == p.ID {
				return
			}
		}
		res = append(res, p)
	}
	for _, r := range AllRelationships() {
		if r.SourceID == e.ID {
			if p := GetPerson(r.DestinationID); p != nil {
				add(p)
			}
		}
		if r.DestinationID == e.ID {
			if p := GetPerson(r.SourceID); p != nil {
				add(p)
			}
		}
	}
	return
}

// RelatedSoftwareSystems returns all software systems the element has a
// relationship with (either as source or as destination).
func (e *Element) RelatedSoftwareSystems() (res SoftwareSystems) {
	add := func(s *SoftwareSystem) {
		for _, es := range res {
			if es.ID == s.ID {
				return
			}
		}
		res = append(res, s)
	}
	for _, r := range AllRelationships() {
		if r.SourceID == e.ID {
			if s := GetSoftwareSystem(r.DestinationID); s != nil {
				add(s)
			}
		}
		if r.DestinationID == e.ID {
			if s := GetSoftwareSystem(r.SourceID); s != nil {
				add(s)
			}
		}
	}
	return
}

// RelatedContainers returns all containers the element has a relationship with
// (either as source or as destination).
func (e *Element) RelatedContainers() (res Containers) {
	add := func(cc *Container) {
		for _, es := range res {
			if es.ID == cc.ID {
				return
			}
		}
		res = append(res, cc)
	}
	for _, r := range AllRelationships() {
		if r.SourceID == e.ID {
			if s := GetContainer(r.DestinationID); s != nil {
				add(s)
			}
		}
		if r.DestinationID == e.ID {
			if s := GetContainer(r.SourceID); s != nil {
				add(s)
			}
		}
	}
	return
}

// RelatedComponents returns all components the element has a relationship with
// (either as source or as destination).
func (e *Element) RelatedComponents() (res Components) {
	add := func(c *Component) {
		for _, es := range res {
			if es.ID == c.ID {
				return
			}
		}
		res = append(res, c)
	}
	for _, r := range AllRelationships() {
		if r.SourceID == e.ID {
			if s := GetComponent(r.DestinationID); s != nil {
				add(s)
			}
		}
		if r.DestinationID == e.ID {
			if s := GetComponent(r.SourceID); s != nil {
				add(s)
			}
		}
	}
	return
}

// Reachable returns the IDs of all elements that can be reached by traversing
// the relationships from the given root.
func (e *Element) Reachable() (res []string) {
	seen := make(map[string]struct{})
	traverse(e, seen)
	res = make([]string, len(seen))
	for k := range seen {
		res = append(res, k)
	}
	return
}
func traverse(e *Element, seen map[string]struct{}) {
	add := func(nid string) bool {
		for id := range seen {
			if id == nid {
				return false
			}
		}
		seen[nid] = struct{}{}
		return true
	}
	for _, r := range AllRelationships() {
		if r.SourceID == e.ID {
			if add(r.DestinationID) {
				traverse(r.Destination, seen)
			}
		}
		if r.DestinationID == e.ID {
			if add(r.SourceID) {
				traverse(r.Source, seen)
			}
		}
	}
	return
}

// MarshalJSON replaces the constant value with the proper structurizr schema
// string value.
func (l LocationKind) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	switch l {
	case LocationInternal:
		buf.WriteString("Internal")
	case LocationExternal:
		buf.WriteString("External")
	case LocationUndefined:
		buf.WriteString("Undefined")
	}
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

// UnmarshalJSON sets the constant from its JSON representation.
func (l *LocationKind) UnmarshalJSON(data []byte) error {
	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	switch val {
	case "Internal":
		*l = LocationInternal
	case "External":
		*l = LocationExternal
	case "Undefined":
		*l = LocationUndefined
	}
	return nil
}
