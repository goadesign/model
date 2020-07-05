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
    ├── Enterprise                          ├── SystemLandscapeView
    ├── Person                              │   ├── Title
    │   ├── Tag                             │   ├── Add
    │   ├── URL                             │   ├── AddAll
    │   ├── External                        │   ├── AddNeighbors
    │   ├── Properties                      │   ├── AddContainers
    │   ├── Uses                            │   ├── Remove
    │   │   └── Tag                         │   ├── RemoveUnreachable
    │   └── InteractsWith                   │   ├── RemoveUnrelated
    │       └── Tag                         │   ├── AutoLayout
    ├── SoftwareSystem                      │   │   ├── RankSeparation
    │   ├── Tag                             │   │   ├── NodeSeparation
    │   ├── URL                             │   │   ├── EdgeSeparation
    │   ├── External                        │   │   └── Vertices
    │   ├── Properties                      │   ├── AnimationStep
    │   ├── Uses                            │   ├── PaperSize
    │   │   └── Tag                         │   └── EnterpriseBoundaryVisible
    │   └── Delivers                        ├── SystemContextView
    │       └── Tag                         │   └──  ... same as SystemLandsapeView*
    ├── Container                           ├── ContainerView
    │   ├── Tag                             │   ├── SystemBoundariesVisible
    │   ├── URL                             │   └── ... same as SystemLandscapeView*
    │   ├── Properties                      ├── ComponentView
    │   ├── Uses                            │   ├── ContainerBoundariesVisible
    │   │   └── Tag                         │   └── ... same as SystemLandscapeView*
    │   └── Delivers                        ├── FilteredView
    │       └── Tag                         │   ├── FilterTag
    ├── Component                           │   └── Exclude
    │   ├── Tag                             ├── DynamicView
    │   ├── URL                             │   ├── Title
    │   ├── Properties                      │   ├── AutoLayout
    │   ├── Uses                            │   │   ├── RankSeparation
    │   │   └── Tag                         │   │   ├── NodeSeparation
    │   └── Delivers                        │   │   ├── EdgeSeparation
    │       └── Tag                         │   │   └── Vertices
    └── DeploymentEnvironment               │   ├── PaperSize
        ├── DeploymentNode                  │   └── Relationship
        │   ├── Tag                         ├── DeploymentView
        │   ├── Instances                   │   └── ... same as SystemLandscapeView*
        │   ├── URL                         ├── Styles
        │   ├── Properties                  │   ├── ElementStyle
        │   └── Uses                        │   │   ├── Shape
        │       └── Tag                     │   │   ├── Icon
        ├── InfrastructureNode              │   │   ├── Width
        │   ├── Tag                         │   │   ├── Height
	    │   ├── URL                         │   │   ├── Background
	    │   ├── Uses                        │   │   ├── Color
	    │   └── Uses                        │   │   ├── Stroke
	    │       └── Tag                     │   │   ├── FontSize
	    └── ContainerInstance               │   │   ├── Border
            ├── Tag                         │   │   ├── Opacity
            ├── InstanceID                  │   │   ├── Metadata
            ├── HealthCheck                 │   │   └── Description
            │   ├── URL                     │   └── RelationshipStyle
            │   ├── Interval                │       ├── Thickness
            │   ├── Timeout                 │       ├── Color
            │   └── Header                  │       ├── Dashed
            ├── Properties                  │       ├── Routing
            └── Uses                        │       ├── FontSize
                └── Tag                     │       ├── Width
                                            │       └── Opacity
                                            ├── Theme
                                            └── Branding

                                            * minus EnterpriseBoundaryVisible
*/
package dsl
