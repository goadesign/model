package expr

import (
	"fmt"
)

// Person represents a person.
type Person struct {
	*Element
}

// EvalName returns the generic expression name used in error messages.
func (p *Person) EvalName() string {
	if p.Name == "" {
		return "unnamed person"
	}
	return fmt.Sprintf("person %q", p.Name)
}
