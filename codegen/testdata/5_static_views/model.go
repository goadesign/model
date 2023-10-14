package model

import . "goa.design/model/dsl"

var _ = Design("Service", "This is a test model.", func() {
	SoftwareSystem("System", "Description", func() {
		Container("Dependency", "Dependency description", func() {
			Component("Endpoint", "Endpoint description", func() {
				Tag("Endpoint")
			})
		})
		Container("Service", "Service description", func() {
			Component("Endpoint", "Endpoint description", func() {
				Tag("Endpoint")
			})
			Uses("Dependency/Endpoint", "Uses endpoint description", "Technology", Synchronous)
			Uses("Dependency", "Uses description", "Technology", Asynchronous)
		})
	})
	Views(func() {
		SystemLandscapeView("01_SystemLandscape", "System landscape description", func() {
			Title("SystemLandscape")
		})
		SystemLandscapeView("02_SystemLandscape_WithAutoLayout", "System context description", func() {
			Title("SystemLandscape_WithAutoLayout")
			AutoLayout(RankLeftRight)
		})
		SystemLandscapeView("03_SystemLandscape_WithAutoLayoutExtended", "System context description", func() {
			Title("SystemLandscape_WithAutoLayoutExtended")
			AutoLayout(RankTopBottom, func() {
				RankSeparation(200)
				NodeSeparation(300)
				EdgeSeparation(400)
				RenderVertices()
			})
		})
		SystemLandscapeView("04_SystemLandscape_WithEnterpriseBoundary", "System context description", func() {
			Title("SystemLandscape_WithEnterpriseBoundary")
			EnterpriseBoundaryVisible()
		})
		SystemLandscapeView("05_SystemLandscape_WithPaperSize", "System context description", func() {
			Title("SystemLandscape_WithPaperSize")
			PaperSize(SizeSlide4X3)
		})
		SystemLandscapeView("06_SystemLandscape_WithAddAll", "System context description", func() {
			Title("SystemLandscape_WithAddAll")
			AddAll()
		})
	})
})
