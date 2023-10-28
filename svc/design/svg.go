package design

import . "goa.design/goa/v3/dsl"

var _ = Service("SVG", func() {
	HTTP(func() {
		Path("/api/diagrams")
	})
	Method("Load", func() {
		Description("Stream the SVG")
		Payload(Filename)
		Result(SVG)
		HTTP(func() {
			GET("/")
			Param("Filename:file")
		})
	})
	Method("Save", func() {
		Description("Save the SVG streamed in the request body")
		Payload(func() {
			Extend(Filename)
			Attribute("SVG", SVG, "Diagram SVG")
			Required("Filename", "SVG")
		})
		HTTP(func() {
			POST("/")
			Param("Filename:file")
			Response(StatusNoContent)
		})
	})
})

var Filename = Type("Filename", func() {
	Attribute("Filename", String, "Diagram SVG filename", func() {
		Pattern(`.*\.svg`)
		Example("diagram.svg")
	})
	Required("Filename")
})

var SVG = Type("SVG", String, func() {
	Description("Scalable Vector Graphics XML")
	Pattern(`<svg.*</svg>$`)
})
