package model

import . "goa.design/model/dsl"

var _ = Design("Empty", "This is a test model.", func() {
	SoftwareSystem("System", "Description")
})
