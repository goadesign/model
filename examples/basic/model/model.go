package model

import . "goa.design/model/dsl"

var _ = Workspace("Getting Started", "This is a model of my software system.", func() {
	var System = SoftwareSystem("Software System", "My software system.", func() {
		Tag("system")
	})

	Person("User", "A user of my software system.", func() {
		Uses(System, "Uses")
		Tag("person")
	})

	Views(func() {
		SystemContextView(System, "SystemContext", "An example of a System Context diagram.", func() {
			AddAll()
			AutoLayout(RankTopBottom)
		})
		Styles(func() {
			ElementStyle("system", func() {
				Background("#1168bd")
				Color("#ffffff")
			})
			ElementStyle("person", func() {
				Shape(ShapePerson)
				Background("#08427b")
				Color("#ffffff")
			})
		})
	})
})
