package design

import . "goa.design/goa/v3/dsl"

var _ = Service("SVG", func() {
	HTTP(func() {
		Path("/api/diagrams")
	})
	Method("Load", func() {
		Description("Stream the model layout JSON saved in the SVG")
		Payload(Filename)
		HTTP(func() {
			GET("/layout")
			Param("Filename:file")
			SkipResponseBodyEncodeDecode()
		})
	})
	Method("Save", func() {
		Description("Save the SVG streamed in the request body")
		Payload(Filename)
		HTTP(func() {
			POST("/")
			Param("Filename:file")
			SkipRequestBodyEncodeDecode()
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
