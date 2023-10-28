package design

import . "goa.design/goa/v3/dsl"

var _ = Service("Repo", func() {
	HTTP(func() {
		Path("/api/repo")
	})
	Method("CreatePackage", func() {
		Description("Create a new model package")
		Payload(PackageFile)
		Error("already_exists", ErrorResult, "Package or directory already exists")
		HTTP(func() {
			POST("/")
			Response(StatusCreated)
			Response("already_exists", StatusConflict)
		})
	})
	Method("DeletePackage", func() {
		Description("Delete the given model package")
		Payload(PackageLocator)
		Error("not_found", ErrorResult, "Package not found")
		HTTP(func() {
			DELETE("/")
			Param("Repository:repo")
			Param("Dir:dir")
			Response(StatusNoContent)
			Response("not_found", StatusNotFound)
		})
	})
	Method("ListPackages", func() {
		Description("List the model packages in the given workspace")
		Payload(Repository)
		Result(ArrayOf(GoPackage))
		Error("not_found", ErrorResult, "Package not found")
		HTTP(func() {
			GET("/")
			Param("Repository:repo")
			Response(StatusOK)
			Response("not_found", StatusNotFound)
		})
	})
	Method("ReadPackage", func() {
		Description("Get the DSL files and their content for the given model package")
		Payload(PackageLocator)
		Result(ArrayOf(PackageFile))
		Error("not_found", ErrorResult, "Package not found")
		HTTP(func() {
			GET("/files")
			Param("Repository:repo")
			Param("Dir:dir")
			Response(StatusOK)
			Response("not_found", StatusNotFound)
		})
	})
	Method("GetModelJSON", func() {
		Description("Compile the given model package and return the model JSON")
		Payload(PackageLocator)
		Result(ModelJSON)
		Error("not_found", ErrorResult, "Package not found")
		Error("compilation_error", ErrorResult, "Compilation error")
		HTTP(func() {
			GET("/json")
			Param("Repository:repo")
			Param("Dir:dir")
			Response(StatusOK)
			Response("not_found", StatusNotFound)
			Response("compilation_error", StatusUnprocessableEntity)
		})
	})
	Method("Subscribe", func() {
		Description("Send model JSON on initial subscription and when the model package changes")
		Payload(PackageLocator)
		StreamingResult(CompilationResults)
		Error("not_found", ErrorResult, "Package not found")
		HTTP(func() {
			GET("/subscribe")
			Param("Repository:repo")
			Param("Dir:dir")
			Response(StatusSwitchingProtocols)
			Response("not_found", StatusNotFound)
		})
	})
})

var Repository = Type("Repository", func() {
	Attribute("Repository", String, "Path to repository root", func() {
		Example("my-repo")
		MinLength(1)
	})
	Meta("struct:pkg:path", "types")
	Required("Repository")
})

var PackageDir = Type("PackageDir", func() {
	Attribute("Dir", String, "Path to directory containing a model package", func() {
		Example("services/my-service/diagram")
		MinLength(1)
	})
	Meta("struct:pkg:path", "types")
	Required("Dir")
})

var PackageLocator = Type("PackageLocator", func() {
	Description("PackageLocator is the location of a model package")
	Extend(Repository)
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
		Example(`package model

import . "goa.design/model/dsl"

var _ = Design(func() {})`)
		MinLength(58)
		Pattern(`import . "goa.design/model/dsl"`)
	})
	Meta("struct:pkg:path", "types")
	Required("Locator", "Content")
})

var CompilationResults = Type("CompilationResults", func() {
	Attribute("Model", ModelJSON, "Model JSON if compilation succeeded")
	Attribute("Error", String, "Compilation error if any")
	Meta("struct:pkg:path", "types")
})

var ModelJSON = Type("ModelJSON", String, func() {
	Format(FormatJSON)
	Docs(func() {
		Description("A serialized representation of a model")
		URL("https://pkg.go.dev/goa.design/model/model#Model")
	})
	Meta("struct:pkg:path", "types")
})
