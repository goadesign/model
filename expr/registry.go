package expr

import (
	"fmt"
	"hash/fnv"
	"math/big"
	"sort"
)

// Registry captures all the elements, people and relationships.
var Registry = make(map[string]any)

// Iterate iterates through all elements, people and relationships in the
// registry in a consistent order.
func Iterate(visitor func(elem any)) {
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
	Iterate(func(e any) {
		if r, ok := e.(*Relationship); ok {
			visitor(r)
		}
	})
}

// Identify sets the ID field of the given element or relationship and registers
// it with the global registry. The algorithm first compute a unique moniker
// for the element or relatioship (based on names and parent scope ID) then
// hashes and base36 encodes the result.
func Identify(element any) {
	var id string
	switch e := element.(type) {
	case *Person:
		id = idify(e.Name)
		e.ID = id
	case *SoftwareSystem:
		id = idify(e.Name)
		e.ID = id
	case *Container:
		id = idify(e.System.ID + ":" + e.Name)
		e.ID = id
	case *Component:
		id = idify(e.Container.ID + ":" + e.Name)
		e.ID = id
	case *DeploymentNode:
		prefix := "dn:" + e.Environment + ":"
		for f := e.Parent; f != nil; f = f.Parent {
			prefix += f.ID + ":"
		}
		id = idify(prefix + e.Name)
		e.ID = id
	case *InfrastructureNode:
		id = idify(e.Environment + ":" + e.Parent.ID + ":" + e.Name)
		e.ID = id
	case *ContainerInstance:
		id = idify(e.Environment + ":" + e.Parent.ID + ":" + e.ContainerID)
		e.ID = id
	case *Relationship:
		var dest string
		if e.Destination != nil {
			dest = e.Destination.ID
		} else {
			dest = e.DestinationPath
		}
		id = idify(e.Source.ID + ":" + dest + ":" + e.Description)
		e.ID = id
	default:
		panic(fmt.Sprintf("element of type %T does not have an ID", element)) // bug
	}
	if _, ok := Registry[id]; ok {
		// Could have been imported from another model package
		return
	}
	Registry[id] = element
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
