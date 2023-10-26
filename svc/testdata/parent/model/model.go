package model

import . "goa.design/model/dsl"

var _ = Design("test", func() {
	SoftwareSystem("parent", func() {
		Container("test")
	})
})
