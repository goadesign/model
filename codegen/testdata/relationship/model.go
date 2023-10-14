package model

import . "goa.design/model/dsl"

var _ = Design("Service", "This is a test model.", func() {
	SoftwareSystem("System", "Description", func() {
		Container("Service", "Service description", func() {
			Component("Endpoint", "Endpoint description", func() {
				Tag("Endpoint")
			})
		})
		Container("Service2", "Service2 description", func() {
			Uses("Service", "Uses description")
			Uses("Service/Endpoint", "Uses endpoint description")
		})
	})
})
