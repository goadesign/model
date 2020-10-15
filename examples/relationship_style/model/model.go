package design

import (
	. "goa.design/model/dsl"
)

var _ = Design("Getting Started", "This is a model of my software system.", func() {
	SoftwareSystem("The System", "My Software System")
	Person("Person1", "A person using the system.", func() {
		Uses("The System", "Edge\nwith vertices", func() {
			Tag("pos75")
		})
		Tag("person")
	})
	Person("Person2", "A person using the system.", func() {
		Uses("The System", "Reads from")
		Uses("The System", "Writes to")
		Tag("person")
	})

	Views(func() {
		SystemContextView("The System", "SystemContext", "System Context diagram.", func() {
			AddAll()
			Link("Person1", "The System", "Edge\nwith vertices", func() {
				Vertices(300, 300, 500, 500)
			})
			AutoLayout(RankLeftRight)
		})
		Styles(func() {
			RelationshipStyle("pos75", func() {
				Position(75)
			})
		})
	})
})
