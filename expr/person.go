package expr

import (
	"fmt"
)

type (
	// Person represents a person.
	Person struct {
		*Element
		Location LocationKind
	}

	// People is a slide of Person that can easily be converted into a slice of ElementHolder.
	People []*Person
)

// PersonTags list the tags that are automatically added to all people.
var PersonTags = []string{"Element", "Person"}

// EvalName returns the generic expression name used in error messages.
func (p *Person) EvalName() string {
	if p.Name == "" {
		return "unnamed person"
	}
	return fmt.Sprintf("person %q", p.Name)
}

// Finalize adds the 'Person' tag ands finalizes relationships.
func (p *Person) Finalize() {
	p.PrefixTags(PersonTags...)
	p.Element.Finalize()
}

// Elements returns a slice of ElementHolder that contains the people.
func (p People) Elements() []ElementHolder {
	res := make([]ElementHolder, len(p))
	for i, pp := range p {
		res[i] = pp
	}
	return res
}
