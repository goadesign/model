package model

import . "goa.design/model/dsl"

var _ = Design("Service", "This is a test model.", func() {
	SoftwareSystem("System", "Description", func() {
		Container("Service", "Service description", func() {
			Component("01_Endpoint", "Endpoint description", func() {
				Tag("Endpoint")
			})
			Component("02_WithURL", "WithURL description", func() {
				URL("https://goa.design/docs/mysystem")
			})
			Component("03_WithProperties", "WithProperties description", func() {
				Prop("1_foo", "bar")
				Prop("2_baz", "qux")
			})
			Component("04_WithAll", "WithAll description", func() {
				Tag("Endpoint")
				URL("https://goa.design/docs/mysystem")
				Prop("1_foo", "bar")
				Prop("2_baz", "qux")
			})
		})
	})
})
