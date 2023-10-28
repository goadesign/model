package design

import (
	. "goa.design/goa/v3/dsl"

	"goa.design/model/expr"
)

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
			PUT("/model/system")
			Response(StatusNoContent)
		})
	})
	Method("UpsertPerson", func() {
		Description("Create or update a person in the model")
		Payload(Person)
		HTTP(func() {
			PUT("/model/person")
			Response(StatusNoContent)
		})
	})
	Method("UpsertContainer", func() {
		Description("Create or update a container in the model")
		Payload(Container)
		HTTP(func() {
			PUT("/model/container")
			Response(StatusNoContent)
		})
	})
	Method("UpsertComponent", func() {
		Description("Create or update a component in the model")
		Payload(Component)
		HTTP(func() {
			PUT("/model/component")
			Response(StatusNoContent)
		})
	})
	Method("UpsertRelationship", func() {
		Description("Create or update a relationship in the model")
		Payload(Relationship)
		HTTP(func() {
			PUT("/model/relationship")
			Response(StatusNoContent)
		})
	})
	Method("UpsertLandscapeView", func() {
		Description("Create or update a landscape view in the model")
		Payload(LandscapeView)
		HTTP(func() {
			PUT("/views/landscape")
			Response(StatusNoContent)
		})
	})
	Method("UpsertSystemContextView", func() {
		Description("Create or update a system context view in the model")
		Payload(SystemContextView)
		HTTP(func() {
			PUT("/views/systemcontext")
			Response(StatusNoContent)
		})
	})
	Method("UpsertContainerView", func() {
		Description("Create or update a container view in the model")
		Payload(ContainerView)
		HTTP(func() {
			PUT("/views/container")
			Response(StatusNoContent)
		})
	})
	Method("UpsertComponentView", func() {
		Description("Create or update a component view in the model")
		Payload(ComponentView)
		HTTP(func() {
			PUT("/views/component")
			Response(StatusNoContent)
		})
	})
	Method("UpserElementStyle", func() {
		Description("Create or update an element style in the model")
		Payload(ElementStyle)
		HTTP(func() {
			PUT("/views/elementstyle")
			Response(StatusNoContent)
		})
	})
	Method("UpsertRelationshipStyle", func() {
		Description("Create or update a relationship style in the model")
		Payload(RelationshipStyle)
		HTTP(func() {
			PUT("/views/relationshipstyle")
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
			DELETE("/model/system/{SystemName}")
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
			DELETE("/model/person/{PersonName}")
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
			DELETE("/model/system/{SystemName}/container/{ContainerName}")
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
			DELETE("/model/system/{SystemName}/container/{ContainerName}/component/{ComponentName}")
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
			DELETE("/model/relationship")
			Response(StatusNoContent)
			Response("NotFound", StatusNotFound)
		})
	})
	Method("DeleteLandscapeView", func() {
		Description("Delete an existing landscape view from the model")
		Payload(func() {
			Extend(FileLocator)
			Attribute("Key", String, "Key of landscape view to delete")
			Required("Key")
		})
		Error("NotFound", ErrorResult, "Landscape view not found")
		HTTP(func() {
			DELETE("/views/landscape/{Key}")
			Response(StatusNoContent)
			Response("NotFound", StatusNotFound)
		})
	})
	Method("DeleteSystemContextView", func() {
		Description("Delete an existing system context view from the model")
		Payload(func() {
			Extend(FileLocator)
			Attribute("Key", String, "Key of system context view to delete")
			Required("Key")
		})
		Error("NotFound", ErrorResult, "System context view not found")
		HTTP(func() {
			DELETE("/views/systemcontext/{Key}")
			Response(StatusNoContent)
			Response("NotFound", StatusNotFound)
		})
	})
	Method("DeleteContainerView", func() {
		Description("Delete an existing container view from the model")
		Payload(func() {
			Extend(FileLocator)
			Attribute("Key", String, "Key of container view to delete")
			Required("Key")
		})
		Error("NotFound", ErrorResult, "Container view not found")
		HTTP(func() {
			DELETE("/views/container/{Key}")
			Response(StatusNoContent)
			Response("NotFound", StatusNotFound)
		})
	})
	Method("DeleteComponentView", func() {
		Description("Delete an existing component view from the model")
		Payload(func() {
			Extend(FileLocator)
			Attribute("Key", String, "Key of component view to delete")
			Required("Key")
		})
		Error("NotFound", ErrorResult, "Component view not found")
		HTTP(func() {
			DELETE("/views/component/{Key}")
			Response(StatusNoContent)
			Response("NotFound", StatusNotFound)
		})
	})
	Method("DeleteElementStyle", func() {
		Description("Delete an existing element style from the model")
		Payload(func() {
			Extend(FileLocator)
			Attribute("Tag", String, "Tag of element style to delete")
			Required("Tag")
		})
		Error("NotFound", ErrorResult, "Element style not found")
		HTTP(func() {
			DELETE("/views/elementstyle/{Tag}")
			Response(StatusNoContent)
			Response("NotFound", StatusNotFound)
		})
	})
	Method("DeleteRelationshipStyle", func() {
		Description("Delete an existing relationship style from the model")
		Payload(func() {
			Extend(FileLocator)
			Attribute("Tag", String, "Tag of relationship style to delete")
			Required("Tag")
		})
		Error("NotFound", ErrorResult, "Relationship style not found")
		HTTP(func() {
			DELETE("/views/relationshipstyle/{Tag}")
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

var ViewBase = Type("ViewBase", func() {
	Description("Base type for all views")
	Attribute("Locator", FileLocator, "Path to file containing view DSL")
	Attribute("Key", String, "Key of view", func() {
		Example("key")
	})
	Attribute("Title", String, "Title of view", func() {
		Example("title")
	})
	Attribute("Description", String, "Description of view", func() {
		Example("description")
	})
	Attribute("PaperSize", String, "Paper size of view", func() {
		Enum(
			expr.SizeA0Landscape.Name(),
			expr.SizeA0Portrait.Name(),
			expr.SizeA1Landscape.Name(),
			expr.SizeA1Portrait.Name(),
			expr.SizeA2Landscape.Name(),
			expr.SizeA2Portrait.Name(),
			expr.SizeA3Landscape.Name(),
			expr.SizeA3Portrait.Name(),
			expr.SizeA4Landscape.Name(),
			expr.SizeA4Portrait.Name(),
			expr.SizeA5Landscape.Name(),
			expr.SizeA5Portrait.Name(),
			expr.SizeA6Landscape.Name(),
			expr.SizeA6Portrait.Name(),
			expr.SizeLegalLandscape.Name(),
			expr.SizeLegalPortrait.Name(),
			expr.SizeLetterLandscape.Name(),
			expr.SizeLetterPortrait.Name(),
			expr.SizeSlide16X10.Name(),
			expr.SizeSlide16X9.Name(),
			expr.SizeSlide4X3.Name(),
		)
	})
	Attribute("ElementViews", ArrayOf(ElementView), "Elements in view")
	Attribute("RelationshipViews", ArrayOf(RelationshipView), "Relationships in view")
	Required("Key", "Title")
})

var ElementView = Type("ElementView", func() {
	Description("ElementView defines an element in a view")
	Attribute("Element", String, "Path to element consisting of <software system name>[/<container name>[/<component name>]]", func() {
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
})

var RelationshipView = Type("RelationshipView", func() {
	Description("RelationshipView defines a relationship in a view")
	Attribute("Source", String, "Path to source element consisting of <software system name>[/<container name>[/<component name>]]", func() {
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
	Attribute("Destination", String, "Path to destination element, see SourcePath for details.", func() {
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
})

var LandscapeView = Type("LandscapeView", func() {
	Extend(ViewBase)
	Attribute("EnterpriseBoundaryVisible", Boolean, "Indicates whether the enterprise boundary is visible on the resulting diagram", func() {
		Default(true)
	})
})

var SystemContextView = Type("SystemContextView", func() {
	Extend(ViewBase)
	Attribute("SoftwareSystemName", String, "Name of software system to create view for", func() {
		Example("Software System")
	})
	Attribute("EnterpriseBoundaryVisible", Boolean, "Indicates whether the enterprise boundary is visible on the resulting diagram", func() {
		Default(true)
	})
	Required("SoftwareSystemName")
})

var ContainerView = Type("ContainerView", func() {
	Extend(ViewBase)
	Attribute("SoftwareSystemName", String, "Name of software system to create view for", func() {
		Example("Software System")
	})
	Attribute("SystemBoundariesVisible", Boolean, "Indicates whether the system boundaries are visible on the resulting diagram", func() {
		Default(true)
	})
})

var ComponentView = Type("ComponentView", func() {
	Extend(ViewBase)
	Attribute("SoftwareSystemName", String, "Name of software system to create view for", func() {
		Example("Software System")
	})
	Attribute("ContainerName", String, "Name of container to create view for", func() {
		Example("Container")
	})
	Attribute("ContainerBoundariesVisible", Boolean, "Indicates whether the container boundaries are visible on the resulting diagram", func() {
		Default(true)
	})
	Required("SoftwareSystemName", "ContainerName")
})

var ElementStyle = Type("ElementStyle", func() {
	Attribute("Tag", String, "Tag of elements to apply style onto", func() {
		Example("tag")
	})
	Attribute("Shape", String, "Shape of element", func() {
		Enum(
			expr.ShapeBox.Name(),
			expr.ShapeCircle.Name(),
			expr.ShapeCylinder.Name(),
			expr.ShapeEllipse.Name(),
			expr.ShapeHexagon.Name(),
			expr.ShapeRoundedBox.Name(),
			expr.ShapeComponent.Name(),
			expr.ShapeFolder.Name(),
			expr.ShapeMobileDeviceLandscape.Name(),
			expr.ShapeMobileDevicePortrait.Name(),
			expr.ShapePerson.Name(),
			expr.ShapePipe.Name(),
			expr.ShapeRobot.Name(),
			expr.ShapeWebBrowser.Name(),
		)
		Default(expr.ShapeBox.Name())
	})
	Attribute("Icon", String, "URL to icon of element", func() {
		Example("https://static.structurizr.com/images/icons/Person.png")
		Format(FormatURI)
	})
	Attribute("Background", String, "Background color of element", func() {
		Pattern("^#[0-9a-fA-F]{6}$")
	})
	Attribute("Color", String, "Text color of element", func() {
		Pattern("^#[0-9a-fA-F]{6}$")
	})
	Attribute("Stroke", String, "Stroke color of element", func() {
		Pattern("^#[0-9a-fA-F]{6}$")
	})
	Attribute("Width", Int, "Width of element", func() {
		Example(100)
	})
	Attribute("Height", Int, "Height of element", func() {
		Example(100)
	})
	Attribute("FontSize", Int, "Font size of element", func() {
		Example(20)
	})
	Attribute("Metadata", Boolean, "Indicates whether the element metadata should be visible on the resulting diagram")
	Attribute("Description", Boolean, "Indicates whether the element description should be visible on the resulting diagram", func() {
		Default(true)
	})
	Attribute("Opacity", Int, "Opacity of element as a percentage", func() {
		Minimum(0)
		Maximum(100)
	})
	Attribute("Border", String, "Type of border to apply to elements", func() {
		Enum(
			expr.BorderSolid.Name(),
			expr.BorderDashed.Name(),
			expr.BorderDotted.Name(),
		)
		Default(expr.BorderSolid.Name())
	})
	Required("Tag")
})

var RelationshipStyle = Type("RelationshipStyle", func() {
	Attribute("Tag", String, "Tag of relationships to apply style onto", func() {
		Example("tag")
	})
	Attribute("Thickness", Int, "Thickness of relationship in pixels", func() {
		Example(2)
		Minimum(0)
		Maximum(1000)
	})
	Attribute("FontSize", Int, "Font size of label on relationship", func() {
		Minimum(1)
		Maximum(100)
	})
	Attribute("Width", Int, "Width of label on relationship", func() {
		Minimum(1)
		Maximum(10000)
	})
	Attribute("Position", Int, "Position of label on relationship as a percentage (0 is next to source, 100 next to destination)", func() {
		Minimum(0)
		Maximum(100)
	})
	Attribute("Color", String, "Color of label", func() {
		Pattern("^#[0-9a-fA-F]{6}$")
	})
	Attribute("Stroke", String, "Stroke color of relationship", func() {
		Pattern("^#[0-9a-fA-F]{6}$")
	})
	Attribute("Dashed", Boolean, "Indicates whether the relationship is dashed", func() {
		Default(true)
	})
	Attribute("Routing", String, "Routing of relationship", func() {
		Enum(
			expr.RoutingDirect.Name(),
			expr.RoutingOrthogonal.Name(),
			expr.RoutingCurved.Name(),
		)
		Default(expr.RoutingDirect.Name())
	})
	Attribute("Opacity", Int, "Opacity of relationship as a percentage", func() {
		Minimum(0)
		Maximum(100)
	})
	Required("Tag")
})
