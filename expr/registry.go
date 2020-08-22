package expr

import (
	"fmt"
	"hash/fnv"
	"math/big"
	"sort"
)

// Registry captures all the elements, people and relationships.
var Registry = make(map[string]interface{})

// Iterate iterates through all elements, people and relationships in the
// registry in a consistent order.
func Iterate(visitor func(elem interface{})) {
	keys := make([]string, len(Registry))
	i := 0
	for k := range Registry {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		visitor(Registry[k])
	}
}

// IterateRelationships iterates through all relationships in the registry in a
// consistent order.
func IterateRelationships(visitor func(r *Relationship)) {
	Iterate(func(e interface{}) {
		if r, ok := e.(*Relationship); ok {
			visitor(r)
		}
	})
}

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
		prefix := "dn:"
		f := e.Parent
		for f != nil {
			prefix += f.ID + ":"
			f = f.Parent
		}
		e.ID = idify(prefix + e.Name)
		Registry[e.ID] = e
	case *InfrastructureNode:
		e.ID = idify(e.Parent.ID + ":" + e.Name)
		Registry[e.ID] = e
	case *ContainerInstance:
		e.ID = idify(e.Parent.ID + ":" + e.ContainerID)
		Registry[e.ID] = e
	case *Relationship:
		var dest string
		if e.Destination != nil {
			dest = e.Destination.ID
		} else {
			dest = e.DestinationPath
		}
		e.ID = idify(e.Source.ID + ":" + dest + ":" + e.Description)
		Registry[e.ID] = e
	default:
		panic(fmt.Sprintf("element of type %T does not have an ID", element)) // bug
	}
}

var h = fnv.New32a()

func idify(s string) string {
	h.Reset()
	h.Write([]byte(s))
	return encodeToBase36(h.Sum(nil))
}

var encodeStd = [36]byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e',
	'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
	'u', 'v', 'w', 'x', 'y', 'z',
}

var bigRadix = big.NewInt(36)
var bigZero = big.NewInt(0)

func encodeToBase36(b []byte) string {
	x := new(big.Int)
	x.SetBytes(b)
	res := make([]byte, 0, len(b)*136/100)
	for x.Cmp(bigZero) > 0 {
		mod := new(big.Int)
		x.DivMod(x, bigRadix, mod)
		res = append(res, encodeStd[mod.Int64()])
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
