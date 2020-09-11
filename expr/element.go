package expr

import (
	"strings"
)

type (
	// Element describes an element.
	Element struct {
		ID            string
		Name          string
		Description   string
		Technology    string
		Tags          string
		URL           string
		Properties    map[string]string
		Relationships []*Relationship
		DSLFunc       func()
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

// DSL returns the attached DSL.
func (e *Element) DSL() func() { return e.DSLFunc }

// Finalize finalizes the relationships.
func (e *Element) Finalize() {
	for _, rel := range e.Relationships {
		rel.Finalize()
	}
}

// GetElement returns the underlying element.
func (e *Element) GetElement() *Element { return e }

// MergeTags adds the given tags. It skips tags already present in e.Tags.
func (e *Element) MergeTags(tags ...string) {
	e.Tags = mergeTags(e.Tags, tags)
}

// PrefixTags adds the given tags to the beginning of the comma separated list.
func (e *Element) PrefixTags(tags ...string) {
	prefix := strings.Join(tags, ",")
	if e.Tags == "" {
		e.Tags = prefix
		return
	}
	e.Tags = mergeTags(prefix, strings.Split(e.Tags, ","))
}

// mergeTags merges the comma separated tags in old with the ones in tags and
// returns a comma separated string with the results.
func mergeTags(existing string, tags []string) string {
	if existing == "" {
		return strings.Join(tags, ",")
	}
	old := strings.Split(existing, ",")
	var merged []string
	for _, o := range old {
		found := false
		for _, tag := range tags {
			if tag == o {
				found = true
				break
			}
		}
		if !found {
			merged = append(merged, o)
		}
	}
	for _, tag := range tags {
		found := false
		for _, o := range merged {
			if tag == o {
				found = true
				break
			}
		}
		if !found {
			merged = append(merged, tag)
		}
	}
	return strings.Join(merged, ",")
}
