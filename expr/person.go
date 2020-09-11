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

// EvalName returns the generic expression name used in error messages.
func (p *Person) EvalName() string {
	if p.Name == "" {
		return "unnamed person"
	}
	return fmt.Sprintf("person %q", p.Name)
}

// Finalize adds the 'Person' tag ands finalizes relationships.
func (p *Person) Finalize() {
	p.PrefixTags("Element", "Person")
}

// Elements returns a slice of ElementHolder that contains the people.
func (p People) Elements() []ElementHolder {
	res := make([]ElementHolder, len(p))
	for i, pp := range p {
		res[i] = pp
	}
	return res
}
