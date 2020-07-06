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

// GetPerson retrieves the person with the given ID from the registry. It
// returns nil if none is found.
func GetPerson(id string) *Person {
	if p, ok := Registry[id].(*Person); ok {
		return p
	}
	return nil
}

// GetSoftwareSystem retrieves the software system with the given ID from the
// registry. It returns nil if none is found.
func GetSoftwareSystem(id string) *SoftwareSystem {
	if s, ok := Registry[id].(*SoftwareSystem); ok {
		return s
	}
	return nil
}

// GetContainer retrieves the container with the given ID from the registry. It
// returns nil if none is found.
func GetContainer(id string) *Container {
	if c, ok := Registry[id].(*Container); ok {
		return c
	}
	return nil
}

// GetComponent retrieves the component with the given ID from the registry. It
// returns nil if none is found.
func GetComponent(id string) *Component {
	if c, ok := Registry[id].(*Component); ok {
		return c
	}
	return nil
}

// GetRelationship retrieves the relationship with the given ID from the registry. It
// returns nil if none is found.
func GetRelationship(id string) *Relationship {
	if r, ok := Registry[id].(*Relationship); ok {
		return r
	}
	return nil
}

// FindRelationship retrieves the relationship with the given source and
// destination ID from the registry. It returns nil if none is found.
func FindRelationship(srcID, destID string) *Relationship {
	for _, e := range Registry {
		if r, ok := e.(*Relationship); ok {
			if r.SourceID == srcID && r.DestinationID == destID {
				return r
			}
		}
	}
	return nil
}

// AllRelationships returns all the relationships in the registry.
func AllRelationships() []*Relationship {
	var res []*Relationship
	for _, e := range Registry {
		if r, ok := e.(*Relationship); ok {
			res = append(res, r)
		}
	}
	return res
}
