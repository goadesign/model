package model

import . "goa.design/model/dsl"

var _ = Design("Dependency", "This is a test model.", func() {
	SoftwareSystem("1_External", "External system", func() {
		External()
		Container("Dependency", "External dependency", func() {
			Component("Endpoint", "Endpoint description", func() {
				Tag("Endpoint")
			})
		})
	})
	SoftwareSystem("2_System", "Description", func() {
		Container("01_Dependency", "Dependency description", func() {
			Component("Endpoint", "Endpoint description", func() {
				Tag("Endpoint")
			})
		})
		Container("02_Basic", "Basic description", func() {
			Uses("01_Dependency/Endpoint", "Uses endpoint description")
			Uses("01_Dependency", "Uses description")
		})
		Container("03_WithTag", "WithTag description", func() {
			Uses("01_Dependency", "Uses description", func() {
				Tag("Foo")
			})
			Uses("01_Dependency/Endpoint", "Uses endpoint description", func() {
				Tag("Foo2")
			})
		})
		Container("04_WithTech", "WithTech description", func() {
			Uses("01_Dependency", "Uses description", "Technology")
			Uses("01_Dependency/Endpoint", "Uses endpoint description", "Technology")
		})
		Container("05_WithTechAndTag", "WithTechAndTag description", func() {
			Uses("01_Dependency", "Uses description", "Technology", func() {
				Tag("Foo")
			})
			Uses("01_Dependency/Endpoint", "Uses endpoint description", "Technology", func() {
				Tag("Foo2")
			})
		})
		Container("06_WithSynchronous", "WithSynchronous description", func() {
			Uses("01_Dependency", "Uses description", "Technology", Synchronous)
			Uses("01_Dependency/Endpoint", "Uses endpoint description", "Technology", Synchronous)
		})
		Container("07_WithSynchronousAndTag", "WithSynchronousAndTag description", func() {
			Uses("01_Dependency", "Uses description", "Technology", Synchronous, func() {
				Tag("Foo")
			})
			Uses("01_Dependency/Endpoint", "Uses endpoint description", "Technology", Synchronous, func() {
				Tag("Foo2")
			})
		})
		Container("08_WithAsync", "WithAsync description", func() {
			Uses("01_Dependency/Endpoint", "Uses endpoint description", "Technology", Asynchronous)
			Uses("01_Dependency", "Uses description", "Technology", Asynchronous)
		})
		Container("09_WithAsyncAndTag", "WithAsyncAndTag description", func() {
			Uses("01_Dependency/Endpoint", "Uses endpoint description", "Technology", Asynchronous, func() {
				Tag("Foo2")
			})
			Uses("01_Dependency", "Uses description", "Technology", Asynchronous, func() {
				Tag("Foo")
			})
		})
		Container("10_External", "External description", func() {
			Uses("1_External/Dependency/Endpoint", "Uses endpoint description")
			Uses("1_External", "Uses description")
			Uses("1_External/Dependency", "Uses dependency description")
		})
		Container("11_ExternalWithTag", "ExternalWithTag description", func() {
			Uses("1_External", "Uses description", func() {
				Tag("Foo")
			})
			Uses("1_External/Dependency", "Uses dependency description", func() {
				Tag("Foo2")
			})
			Uses("1_External/Dependency/Endpoint", "Uses endpoint description", func() {
				Tag("Foo3")
			})
		})
	})
})
