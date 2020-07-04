/*
Package dsl implements a Go based DSL that makes it possible to describe
softare architecture models following the C4 model (https://c4model.com).

It is recommended to use "dot import" when using this package to write DSLs,
for example:

	package model

	import . "goa.design/structurizr/dsl"

	var _ = Workspace("<name>", "[description]", func() {
		// ...
	})

The DSL can be executed via the eval package. The resulting data structure
JSON representation is suitable for uploading to the Structurizr service
(https://structurizr.com). It can also be used to render diagrams in
different formats (https://structurizr.com/help/code).

Some DSL functions accept a anonymous function as last argument (such as
Workspace above) which makes it possible to define a nesting structure. The
general shape of the DSL is:

    Workspace                           Workspace
    ├── Version                         └── Views
    ├── Enterprise                          ├── SystemLandscape
    ├── Person                              │   ├── Title
    │   ├── Tag                             │   ├── IncludeAll
    │   ├── URL                             │   ├── Include
    │   ├── External                        │   ├── Exclude
    │   ├── Properties                      │   ├── AutoLayout
    │   ├── Uses                            │   │   ├── RankSeparation
    │   │   └── Tag                         │   │   ├── NodeSeparation
    │   └── InteractsWith                   │   │   ├── EdgeSeparation
    │       └── Tag                         │   │   └── Vertices
    ├── SoftwareSystem                      │   ├── AnimationStep
    │   ├── Tag                             │   ├── PaperSize
    │   ├── URL                             │   └── EnterpriseBoundaryVisible
    │   ├── External                        ├── SystemContext
    │   ├── Properties                      │   └──  ... same as SystemLandsape
    │   ├── Uses                            ├── Container
    │   │   └── Tag                         │   ├── SystemBoundariesVisible
    │   └── Delivers                        │   └── ... same as SystemLandscape
    │       └── Tag                         ├── Component
    ├── Container                           │   ├── ContainerBoundariesVisible
    │   ├── Tag                             │   └── ... same as SystemLandscape
    │   ├── URL                             ├── Filtered
    │   ├── Properties                      │   ├── FilterTag
    │   ├── Uses                            │   └── Exclude
    │   │   └── Tag                         ├── Dynamic
    │   └── Delivers                        │   ├── Title
    │       └── Tag                         │   ├── AutoLayout
    ├── Component                           │   │   ├── RankSeparation
    │   ├── Tag                             │   │   ├── NodeSeparation
    │   ├── URL                             │   │   ├── EdgeSeparation
    │   ├── Properties                      │   │   └── Vertices
    │   ├── Uses                            │   ├── PaperSize
    │   │   └── Tag                         │   └── Relationship
    │   └── Delivers                        ├── Deployment
    │       └── Tag                         │   └── ... same as SystemLandscape
    └── DeploymentEnvironment               ├── Element
        ├── DeploymentNode                  ├── Relationship
        │   ├── Tag                         │   ├── Description
        │   ├── Instances                   │   ├── Order
        │   ├── URL                         │   ├── Vertices
        │   ├── Properties                  │   ├── Routing
        │   └── Uses                        │   └── Position
        │       └── Tag                     ├── Styles
        ├── InfrastructureNode              │   ├── Element
	    │   ├── Tag                         │   │   ├── Shape
	    │   ├── URL                         │   │   ├── Icon
	    │   ├── Uses                        │   │   ├── Width
	    │   └── Uses                        │   │   ├── Height
	    │       └── Tag                     │   │   ├── Background
	    └── ContainerInstance               │   │   ├── Color
            ├── Tag                         │   │   ├── Stroke
            ├── InstanceID                  │   │   ├── FontSize
            ├── HealthCheck                 │   │   ├── Border
            │   ├── URL                     │   │   ├── Opacity
            │   ├── Interval                │   │   ├── Metadata
            │   ├── Timeout                 │   │   └── Description
            │   └── Header                  │   └── Relationship
            ├── Properties                  │       ├── Thickness
            └── Uses                        │       ├── Color
                └── Tag                     │       ├── Dashed
                                            │       ├── Routing
                                            │       ├── FontSize
                                            │       ├── Width
                                            │       ├── Position
                                            │       └── Opacity
                                            ├── Theme
                                            └── Branding
*/
package dsl
