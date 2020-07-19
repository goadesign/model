package expr

import (
	"fmt"

	"github.com/rs/xid"
)

// Registry captures all the elements, people and relationships.
var Registry = make(map[string]interface{})

// Identify sets the ID field of the given element and registers it with the
// global registery.
func Identify(element interface{}) {
	id := xid.New().String()
	switch e := element.(type) {
	case *Person:
		e.ID = id
	case *SoftwareSystem:
		e.ID = id
	case *Container:
		e.ID = id
	case *Component:
		e.ID = id
	case *DeploymentNode:
		e.ID = id
	case *InfrastructureNode:
		e.ID = id
	case *ContainerInstance:
		e.ID = id
	case *Relationship:
		e.ID = id
	default:
		panic(fmt.Sprintf("element of type %T does not have an ID", element)) // bug
	}
	Registry[id] = element
}
