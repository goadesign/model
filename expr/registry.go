package expr

import (
	"fmt"
	"hash/fnv"
	"math/big"
)

// Registry captures all the elements, people and relationships.
var Registry = make(map[string]interface{})

// Identify sets the ID field of the given element or relationship and registers
// it with the global registery. The algorithm first compute a unique moniker
// for the element or relatioship (based on names and parent scope ID) then
// hashes and base36 encodes the result.
func Identify(element interface{}) {
	switch e := element.(type) {
	case *Person:
		e.ID = idify(e.Name)
		Registry[e.ID] = e
	case *SoftwareSystem:
		e.ID = idify(e.Name)
		Registry[e.ID] = e
	case *Container:
		e.ID = idify(e.System.ID + ":" + e.Name)
		Registry[e.ID] = e
	case *Component:
		e.ID = idify(e.Container.ID + ":" + e.Name)
		Registry[e.ID] = e
	case *DeploymentNode:
		var prefix string
		f := e.Parent
		for f != nil {
			prefix += f.ID
			f = f.Parent
		}
		e.ID = idify(prefix + e.Name)
		Registry[e.ID] = e
	case *InfrastructureNode:
		e.ID = idify(e.Parent.ID + e.Name)
		Registry[e.ID] = e
	case *ContainerInstance:
		e.ID = idify(e.Parent.ID + e.Name)
		Registry[e.ID] = e
	case *Relationship:
		e.ID = idify(e.SourceID + ":" + e.DestinationID + ":" + e.Description)
		Registry[e.ID] = e
	default:
		panic(fmt.Sprintf("element of type %T does not have an ID", element)) // bug
	}
}

var bigRadix = big.NewInt(36)
var bigZero = big.NewInt(0)
var h = fnv.New32a()

func idify(s string) string {
	h.Reset()
	h.Write([]byte(s))
	b := h.Sum(nil)
	x := new(big.Int)
	x.SetBytes(b)
	res := make([]byte, 0, len(b)*136/100)
	for x.Cmp(bigZero) > 0 {
		mod := new(big.Int)
		x.DivMod(x, bigRadix, mod)
		offset := mod.Int64()
		if offset < 10 {
			res = append(res, byte(48+offset))
		} else {
			res = append(res, byte(87+offset))
		}
	}
	for _, i := range b {
		if i != 0 {
			break
		}
		res = append(res, byte(48))
	}
	alen := len(res)
	for i := 0; i < alen/2; i++ {
		res[i], res[alen-1-i] = res[alen-1-i], res[i]
	}
	return string(res)
}
