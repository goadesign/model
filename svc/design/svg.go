package design

import . "goa.design/goa/v3/dsl"

var _ = Service("SVG", func() {
	HTTP(func() {
		Path("/api/diagrams")
	})
	Method("Load", func() {
		Description("Stream the SVG")
		Payload(FileLocator)
		Result(SVG)
		Error("NotFound", ErrorResult, "File not found")
		HTTP(func() {
			GET("/")
			Param("Filename:filename")
			Param("Repository:repo")
			Param("Dir:dir")
			Response(StatusOK)
			Response("NotFound", StatusNotFound)
		})
	})
	Method("Save", func() {
		Description("Save the SVG streamed in the request body")
		Payload(func() {
			Extend(FileLocator)
			Attribute("SVG", SVG, "Diagram SVG")
			Required("Filename", "SVG")
		})
		HTTP(func() {
			POST("/")
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
