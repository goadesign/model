package model

import . "goa.design/model/dsl"

var _ = Design("Service", "This is a test model.", func() {
	SoftwareSystem("System", "Description", func() {
		Container("Service", "Service description")
	})
})
