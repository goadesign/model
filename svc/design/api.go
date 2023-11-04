package design

import . "goa.design/goa/v3/dsl"

var _ = API("model", func() {
	Title("Model")
	Description("Model is a service for managing software models using diagram as code.")
	Server("model", func() {
		Host("localhost", func() {
			URI("http://localhost:8080")
		})
	})
	Docs(func() {
		Description("Model open-source project")
		URL("https://github.com/goadesign/model/")
	})
})
