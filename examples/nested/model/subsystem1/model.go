package design

import (
	. "goa.design/model/dsl"
	"goa.design/model/examples/nested/styles"
)

// Subsystem1 defines the design for subsystem 1.
var Subsystem1 = Design("Subsystem 1", "This is a model of subsystem 1.", func() {
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
