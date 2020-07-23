package model

import (
	. "goa.design/structurizr/dsl"
	"goa.design/structurizr/examples/nested/styles"
)

var Subsystem1 = Workspace("Subsystem 1", "This is a model of subsystem 1.", func() {
	var System = SoftwareSystem("Subsystem 1", "A software system that belongs to subsystem 1.", func() {
		Tag("system")
	})

	Person("User", "A user of Subsystem 1.", func() {
		Uses(System, "Uses")
		Tag("person")
	})

	Views(func() {
		styles.DefineAll() // Use shared styles

		SystemContextView(System, "Subsystem 1 context", "System context diagram for Subsystem 1.", func() {
			AddAll()
			AutoLayout(RankTopBottom)
		})
	})
})
