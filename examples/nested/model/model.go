package design

import (
	_ "goa.design/model/examples/nested/model/subsystem1"
	s2 "goa.design/model/examples/nested/model/subsystem2"

	. "goa.design/model/dsl"
)

var _ = Design("Global workspace", "The model for all systems", func() {
	// Add a new dependency for the person "User" defined in subsystem 1 to the
	// software system defined in subsystem 2.
	Person("User", "A user of both Subsystems.", func() {
		Uses(s2.Subsystem2.SoftwareSystem("Subsystem 2"), "Uses")
		Tag("person")
	})
})
