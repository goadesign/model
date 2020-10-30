package design

import (
	. "goa.design/model/dsl"
)

var _ = Design("Getting Started", "This is a model of my software system.", func() {
	SoftwareSystem("The System", "My Software System")
	Person("Person1", "A person using the system.", func() {
		Uses("The System", "Thick, red edge\nwith vertices", func() {
			Tag("labelPos")
		})
	})
	Person("Person2", "Two relationships\nautomatically spread", func() {
		Uses("The System", "Right")
		Uses("The System", "Left")
		Tag("Customer")
	})
	Person("Person3", "Another person", func() {
		Uses("The System", "Solid, Dashed\nOrthogonal", func() {
			Tag("knows")
		})
	})

	Views(func() {
		SystemContextView("The System", "SystemContext", "System Context diagram.", func() {
			AddAll()
			Link("Person1", "The System", "Thick, red edge\nwith vertices", func() {
				Vertices(300, 300, 300, 800)
			})
			AutoLayout(RankLeftRight)
		})
		Styles(func() {
			ElementStyle("Person", func() {
				Shape(ShapePerson)
				Background("#fffff0")
			})
			ElementStyle("Customer", func() {
				Background("#ffffa0")
			})
			// Defaults all relationships to solid line
			RelationshipStyle("Relationship", func() {
				Solid()
			})
			RelationshipStyle("labelPos", func() {
				Position(40)
				Color("#FF0000")
				Thickness(7)
			})
			RelationshipStyle("knows", func() {
				Routing(RoutingOrthogonal)
				// overwrite line style
				Dashed()
			})
		})
	})
})
