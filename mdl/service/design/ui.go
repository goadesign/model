package design

import . "goa.design/goa/v3/dsl"

var _ = Service("Assets", func() {
	Files("/assets", "assets", func() {
		Description("Serve the diagram editor UI")
	})
})
