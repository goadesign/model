package design

import . "goa.design/goa/v3/dsl"

var _ = Service("DSLEditor", func() {
	Error("compilation_failed", ErrorResult, "Compilation failed")
	HTTP(func() {
		Path("/api/dsl")
		Response("compilation_failed", StatusBadRequest)
	})
	Method("UpdateDSL", func() {
		Description("Update the DSL for the given package, compile it and return the corresponding JSON if successful")
		Payload(PackageFile)
		HTTP(func() {
			POST("/")
			Response(StatusNoContent)
		})
	})
	Method("UpsertSystem", func() {
		Description("Create or update a software system in the model")
		Payload(System)
		HTTP(func() {
			PUT("/system")
			Response(StatusNoContent)
		})
	})
	Method("UpsertPerson", func() {
		Description("Create or update a person in the model")
		Payload(Person)
		HTTP(func() {
			PUT("/person")
			Response(StatusNoContent)
		})
	})
	Method("UpsertContainer", func() {
		Description("Create or update a container in the model")
		Payload(Container)
		HTTP(func() {
			PUT("/container")
			Response(StatusNoContent)
		})
	})
	Method("UpsertComponent", func() {
		Description("Create or update a component in the model")
		Payload(Component)
		HTTP(func() {
			PUT("/component")
			Response(StatusNoContent)
		})
	})
	Method("UpsertRelationship", func() {
		Description("Create or update a relationship in the model")
		Payload(Relationship)
		HTTP(func() {
			PUT("/relationship")
			Response(StatusNoContent)
		})
	})
	Method("DeleteSystem", func() {
		Description("Delete an existing software system from the model")
		Payload(func() {
			Extend(FileLocator)
			Attribute("SystemName", String, "Name of software system to delete")
			Required("SystemName")
		})
		Error("NotFound", ErrorResult, "Software system not found")
		HTTP(func() {
			DELETE("/system/{SystemName}")
			Response(StatusNoContent)
			Response("NotFound", StatusNotFound)
		})
	})
	Method("DeletePerson", func() {
		Description("Delete an existing person from the model")
		Payload(func() {
			Extend(FileLocator)
			Attribute("PersonName", String, "Name of person to delete")
			Required("PersonName")
		})
		Error("NotFound", ErrorResult, "Person not found")
		HTTP(func() {
			DELETE("/person/{PersonName}")
			Response(StatusNoContent)
			Response("NotFound", StatusNotFound)
		})
	})
	Method("DeleteContainer", func() {
		Description("Delete an existing container from the model")
		Payload(func() {
			Extend(FileLocator)
			Attribute("SystemName", String, "Name of container software system")
			Attribute("ContainerName", String, "Name of container to delete")
			Required("ContainerName")
		})
		Error("NotFound", ErrorResult, "Container not found")
		HTTP(func() {
			DELETE("/system/{SystemName}/container/{ContainerName}")
			Response(StatusNoContent)
			Response("NotFound", StatusNotFound)
		})
	})
	Method("DeleteComponent", func() {
		Description("Delete an existing component from the model")
		Payload(func() {
			Extend(FileLocator)
			Attribute("SystemName", String, "Name of component software system", func() {
				Example("My System")
			})
			Attribute("ContainerName", String, "Name of component software system", func() {
				Example("My Container")
			})
			Attribute("ComponentName", String, "Name of component to delete", func() {
				Example("My Component")
			})
			Required("SystemName", "ContainerName", "ComponentName")
		})
		Error("NotFound", ErrorResult, "Component not found")
		HTTP(func() {
			DELETE("/system/{SystemName}/container/{ContainerName}/component/{ComponentName}")
			Response(StatusNoContent)
			Response("NotFound", StatusNotFound)
		})
	})
	Method("DeleteRelationship", func() {
		Description("Delete an existing relationship from the model")
		Payload(func() {
			Extend(FileLocator)
			Attribute("SourcePath", String, "Path to source element consisting of <software system name>[/<container name>[/<component name>]]", func() {
				Example("Software System", func() {
					Value("Software System")
				})
				Example("Container", func() {
					Value("Software System/Container")
				})
				Example("Component", func() {
					Value("Software System/Container/Component")
				})
			})
			Attribute("DestinationPath", String, "Path to destination element, see SourcePath for details.", func() {
				Example("Software System", func() {
					Value("Software System")
				})
				Example("Container", func() {
					Value("Software System/Container")
				})
				Example("Component", func() {
					Value("Software System/Container/Component")
				})
			})
			Required("SourcePath", "DestinationPath")
		})
		Error("NotFound", ErrorResult, "Relationship not found")
		HTTP(func() {
			DELETE("/relationship")
			Response(StatusNoContent)
			Response("NotFound", StatusNotFound)
		})
	})
})

var System = Type("System", func() {
	Attribute("Locator", FileLocator, "Path to file containing system DSL")
	Attribute("Name", String, "Name of software system", func() {
		Example("System")
	})
	Attribute("Description", String, "Description of system", func() {
		Example("System description")
	})
	Attribute("Tags", ArrayOf(String), "Attached tags", func() {
		Example([]string{"Tag1", "Tag2"})
	})
	Attribute("URL", String, "Documentation URL", func() {
		Example("https://system.com")
	})
	Attribute("Location", String, "Indicates whether the system is in-house (Internal) or hosted by a third party (External)", func() {
		Enum("Internal", "External")
		Default("Internal")
	})
	Attribute("Properties", MapOf(String, String), "Set of arbitrary name-value properties (shown in diagram tooltips)", func() {
		Example(map[string]string{"key1": "value1", "key2": "value2"})
	})
	Required("Name")
})

var Person = Type("Person", func() {
	Attribute("Locator", FileLocator, "Path to file containing person DSL")
	Attribute("Name", String, "Name of person", func() {
		Example("Person")
	})
	Attribute("Description", String, "Description of person", func() {
		Example("Person description")
	})
	Attribute("Tags", ArrayOf(String), "Attached tags", func() {
		Example([]string{"Tag1", "Tag2"})
	})
	Attribute("URL", String, "Documentation URL", func() {
		Example("https://person.com")
	})
	Attribute("Location", String, "Indicates whether the person is an employee (Internal) or a third party (External)", func() {
		Enum("Internal", "External")
		Default("Internal")
	})
	Attribute("Properties", MapOf(String, String), "Set of arbitrary name-value properties (shown in diagram tooltips)", func() {
		Example(map[string]string{"key1": "value1", "key2": "value2"})
	})
	Required("Name")
})

var Container = Type("Container", func() {
	Attribute("Locator", FileLocator, "Path to file containing container DSL")
	Attribute("SystemName", String, "Name of parent software system", func() {
		Example("My System")
	})
	Attribute("Name", String, "Name of container", func() {
		Example("Container")
	})
	Attribute("Description", String, "Description of container", func() {
		Example("Container description")
	})
	Attribute("Technology", String, "Technology used by container", func() {
		Example("Technology")
	})
	Attribute("Tags", ArrayOf(String), "Attached tags", func() {
		Example([]string{"Tag1", "Tag2"})
	})
	Attribute("URL", String, "Documentation URL", func() {
		Example("https://container.com")
	})
	Attribute("Properties", MapOf(String, String), "Set of arbitrary name-value properties (shown in diagram tooltips)", func() {
		Example(map[string]string{"key1": "value1", "key2": "value2"})
	})
	Required("SystemName", "Name")
})

var Component = Type("Component", func() {
	Attribute("Locator", FileLocator, "Path to file containing component DSL")
	Attribute("SystemName", String, "Name of parent software system", func() {
		Example("My System")
	})
	Attribute("ContainerName", String, "Name of parent container", func() {
		Example("My Container")
	})
	Attribute("Name", String, "Name of component", func() {
		Example("Component")
	})
	Attribute("Description", String, "Description of component", func() {
		Example("Component description")
	})
	Attribute("Technology", String, "Technology used by component", func() {
		Example("Technology")
	})
	Attribute("Tags", ArrayOf(String), "Attached tags", func() {
		Example([]string{"Tag1", "Tag2"})
	})
	Attribute("URL", String, "Documentation URL", func() {
		Example("https://component.com")
	})
	Attribute("Properties", MapOf(String, String), "Set of arbitrary name-value properties (shown in diagram tooltips)", func() {
		Example(map[string]string{"key1": "value1", "key2": "value2"})
	})
	Required("SystemName", "ContainerName", "Name")
})

var Relationship = Type("Relationship", func() {
	Attribute("Locator", FileLocator, "Path to file containing relationship DSL")
	Attribute("SourcePath", String, "Path to source element consisting of <software system name>[/<container name>[/<component name>]]", func() {
		Example("Software System", func() {
			Value("Software System")
		})
		Example("Container", func() {
			Value("Software System/Container")
		})
		Example("Component", func() {
			Value("Software System/Container/Component")
		})
	})
	Attribute("DestinationPath", String, "Path to destination element, see SourcePath for details.", func() {
		Example("Software System", func() {
			Value("Software System")
		})
		Example("Container", func() {
			Value("Software System/Container")
		})
		Example("Component", func() {
			Value("Software System/Container/Component")
		})
	})
	Attribute("Description", String, "Description of relationship", func() {
		Example("Relationship description")
	})
	Attribute("Technology", String, "Technology used by relationship", func() {
		Example("Technology")
	})
	Attribute("InteractionStyle", String, "Indicates whether the relationship is synchronous or asynchronous", func() {
		Enum("Synchronous", "Asynchronous")
		Default("Synchronous")
	})
	Attribute("Tags", ArrayOf(String), "Attached tags", func() {
		Example([]string{"Tag1", "Tag2"})
	})
	Attribute("URL", String, "Documentation URL", func() {
		Format(FormatURI)
		Example("https://relationship.com")
	})
	Required("SourcePath", "DestinationPath")
})
