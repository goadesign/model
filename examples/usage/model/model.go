package design

import . "goa.design/model/dsl"

var _ = Design("Model Usage", "Not a software architecture but a diagram illustrating how to use Model.", func() {
	SoftwareSystem("Model Usage", func() {
		Container("Model Usage", func() {
			Component("Design", "Go package containing Model DSL that describes the system architecture", func() {
				Uses("Visual Editor", "mdl serve")
				Uses("Design JSON", "mdl gen")
				// Uses("Static HTML", "mdl gen")
				Uses("Structurizr workspace JSON", "stz gen")
				Tag("Design")
			})
			Component("Design JSON", "JSON representation of the design")
			Component("Structurizr workspace JSON", "JSON representation of a Structurizr workspace corresponding to the software architecture model described in the design", func() {
				Uses("Structurizr Service", "stz put")
				Tag("Structurizr")
			})
			Component("Visual Editor", "Edit diagram element positions and save to SVG", func() {
				Uses("View Renderings", "Save")
				Tag("Editor")
			})
			// Component("Static HTML", "Static HTML rendering of each view")
			Component("View Renderings", "SVG files corresponding to design views.", func() {
				Tag("SVG")
			})
			Component("Structurizr Service", "Structurizr service hosted at https://structurizr.com", func() {
				Tag("Editor", "Structurizr")
			})
		})
	})

	Person("User", "Model user.", func() {
		Uses("Model Usage/Model Usage/Design", "Writes")
		Uses("Model Usage/Model Usage/Visual Editor", "Uses")
		Tag("Person")
	})

	Views(func() {
		ComponentView("Model Usage/Model Usage", "view", func() {
			AddDefault()
		})
		Styles(func() {
			ElementStyle("Person", func() {
				Shape(ShapePerson)
				Stroke("#55AAEE")
			})
			ElementStyle("Design", func() {
				Shape(ShapeFolder)
			})
			ElementStyle("Editor", func() {
				Shape(ShapeWebBrowser)
				Stroke("#EEAA55")
			})
			ElementStyle("SVG", func() {
				Shape(ShapeRoundedBox)
				Stroke("#EEAA55")
			})
			ElementStyle("Structurizr", func() {
				Stroke("#55AAEE")
			})
		})
	})
})
