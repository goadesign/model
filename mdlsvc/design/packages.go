package design

import . "goa.design/goa/v3/dsl"

var _ = Service("Packages", func() {
	HTTP(func() {
		Path("/api/packages")
	})
	Method("ListPackages", func() {
		Description("List the model packages in the given workspace")
		Payload(Workspace)
		Result(ArrayOf(GoPackage))
		HTTP(func() {
			GET("/")
			Param("Workspace:work")
			Response(StatusOK)
		})
	})
	Method("ListPackageFiles", func() {
		Description("Get the DSL files and their content for the given model package")
		Payload(PackageLocator)
		Result(ArrayOf(PackageFile))
		HTTP(func() {
			GET("/files")
			Param("Workspace:work")
			Param("Dir:dir")
			Response(StatusOK)
		})
	})
	Method("Subscribe", func() {
		Description("Send model JSON on initial subscription and when the model package changes")
		Payload(PackageLocator)
		StreamingResult(ModelJSON)
		HTTP(func() {
			GET("/subscribe")
			Param("Workspace:work")
			Param("Dir:dir")
			Response(StatusSwitchingProtocols)
		})
	})
	Method("GetModelJSON", func() {
		Description("Streams the model JSON for the given package, see https://pkg.go.dev/goa.design/model/model#Model")
		Payload(PackageLocator)
		HTTP(func() {
			GET("/model")
			Param("Workspace:work")
			Param("Dir:dir")
			SkipResponseBodyEncodeDecode()
		})
	})
	Method("GetLayout", func() {
		Description("Streams the model layout JSON for the given package")
		Payload(PackageLocator)
		HTTP(func() {
			GET("/layout")
			Param("Workspace:work")
			Param("Dir:dir")
			SkipResponseBodyEncodeDecode()
		})
	})
})

var Workspace = Type("Workspace", func() {
	Attribute("Workspace", String, "Workspace identifier", func() {
		Example("my-workspace")
		MinLength(1)
	})
	Meta("struct:pkg:path", "types")
	Required("Workspace")
})

var PackageDir = Type("PackageDir", func() {
	Attribute("Dir", String, "Path to directory containing a model package", func() {
		Example("src/repo/model")
		MinLength(1)
	})
	Meta("struct:pkg:path", "types")
	Required("Dir")
})

var PackageLocator = Type("PackageLocator", func() {
	Description("PackageLocator is the location of a model package in a workspace")
	Extend(Workspace)
	Extend(PackageDir)
	Meta("struct:pkg:path", "types")
})

var FileLocator = Type("FileLocator", func() {
	Description("FileLocator is the location of a DSL file in a model package")
	Extend(PackageLocator)
	Attribute("Filename", String, "Name of DSL file", func() {
		Pattern(`\.go$`)
		Example("model.go")
	})
	Meta("struct:pkg:path", "types")
	Required("Filename")
})

var GoPackage = Type("Package", func() {
	Extend(PackageDir)
	Attribute("ImportPath", String, "Design Go package import path", func() {
		domainRegex := `^([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*\.)+[a-zA-Z]{2,}`
		orgRegex := `[a-zA-Z0-9_\-]+`
		pathRegex := `(/([a-zA-Z0-9_\-]+))*$`
		Pattern(domainRegex + "/" + orgRegex + "/" + pathRegex)
		Example("goa.design/model/examples/basic/model")
	})
	Meta("struct:pkg:path", "types")
	Required("ImportPath")
})

var PackageFile = Type("PackageFile", func() {
	Attribute("Locator", FileLocator, "Path to file containing DSL code")
	Attribute("Content", String, "DSL code", func() {
		Example(`import . "goa.design/model/dsl"

var _ = Design(func() {})`)
		MinLength(58)
		Pattern(`import . "goa.design/model/dsl"`)
	})
	Meta("struct:pkg:path", "types")
	Required("Locator", "Content")
})

var ModelJSON = Type("ModelJSON", String, func() {
	Format(FormatJSON)
	Docs(func() {
		Description("A serialized representation of a model")
		URL("https://pkg.go.dev/goa.design/model/model#Model")
	})
	Meta("struct:pkg:path", "types")
})
