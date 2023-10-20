package design

import . "goa.design/model/dsl"

var _ = Design("Getting Started", "This is a model of my software system.", func() {
	var System = SoftwareSystem("Software System", "My software system.", func() {
		Container("Application Database", "Stores application data.", func() {
			Tag("database")
		})
		Container("Web Application", "Delivers content to users.", func() {
			Component("Dashboard Endpoint", "Serve dashboard content.", func() {
				Tag("endpoint")
			})
			Uses("Application Database", "Reads from and writes to", "MySQL", Synchronous)
		})
		Container("Load Balancer", "Distributes requests across the Web Application instances.", func() {
			Uses("Web Application/Dashboard Endpoint", "Routes requests to", "HTTPS", Synchronous)
		})
		Tag("system")
	})

	Person("User", "A user of my software system.", func() {
		Uses(System, "Uses", Synchronous)
		Tag("person")
	})

	Views(func() {
		SystemContextView(System, "SystemContext", "An example of a System Context diagram.", func() {
			AddAll()
			AutoLayout(ImplementationDagre, RankLeftRight)
		})
		Styles(func() {
			ElementStyle("system", func() {
				Background("#1168bd")
				Color("#ffffff")
			})
			ElementStyle("person", func() {
				Background("#08427b")
				Color("#ffffff")
				Shape(ShapePerson)
			})
			ElementStyle("database", func() {
				Shape(ShapeCylinder)
			})
		})
	})
})
