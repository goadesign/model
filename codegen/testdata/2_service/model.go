package model

import . "goa.design/model/dsl"

var _ = Design("Service", "This is a test model.", func() {
	SoftwareSystem("System", "Description", func() {
		Container("1_Empty", "Empty description")
		Container("2_WithTag", "WithTag description", func() {
			Tag("Foo")
		})
		Container("3_WithURL", "WithURL description", func() {
			URL("https://goa.design/docs/mysystem")
		})
		Container("4_WithProperties", "WithProperties description", func() {
			Prop("1_foo", "bar")
			Prop("2_baz", "qux")
		})
		Container("5_WithAll", "WithAll description", func() {
			Tag("Foo")
			URL("https://goa.design/docs/mysystem")
			Prop("1_foo", "bar")
			Prop("2_baz", "qux")
		})
	})
})
