/*
Package dsl implements a Go based DSL that makes it possible to describe
softare architecture models following the C4 model (https://c4model.com).

It is recommended to use "dot import" when using this package to write DSLs,
for example:

	package design

	import . "goa.design/model/dsl"

	var _ = Design("<name>", "[description]", func() {
	    // ...
	})

Some DSL functions accept a anonymous function as last argument (such as
Design above) which makes it possible to define a nesting structure. The
general shape of the DSL is:

	Design                              Design
	├── Version                         └── Views
	├── Enterprise                          ├── SystemLandscapeView
	├── Person                              │   ├── Title
	│   ├── Tag                             │   ├── AddDefault
	│   ├── URL                             │   ├── Add
	│   ├── External                        │   ├── AddAll
	│   ├── Prop                            │   ├── AddNeighbors
	│   ├── Uses                            │   ├── Link
	│   └── InteractsWith                   │   ├── Remove
	├── SoftwareSystem                      │   ├── RemoveTagged
	│   ├── Tag                             │   ├── RemoveUnreachable
	│   ├── URL                             │   ├── RemoveUnrelated
	│   ├── External                        │   ├── Unlink
	│   ├── Prop                            │   ├── AutoLayout
	│   ├── Uses                            │   ├── AnimationStep
	│   ├── Delivers                        │   ├── PaperSize
	│   └─── Container                      │   └── EnterpriseBoundaryVisible
	│       ├── Tag                         ├── SystemContextView
	│       ├── URL                         │   └──  ... (same as SystemLandsapeView)
	│       ├── Prop                        ├── ContainerView
	│       ├── Uses                        │   ├── AddContainers
	│       ├── Delivers                    │   ├── AddInfluencers
	│       └── Component                   │   ├── SystemBoundariesVisible
	│           ├── Tag                     │   └── ... (same as SystemLandscapeView*)
	│           ├── URL                     ├── ComponentView
	│           ├── Prop                    │   ├── AddContainers
	│           ├── Uses                    │   ├── AddComponents
	│           └── Delivers                │   ├── ContainerBoundariesVisible
	└── DeploymentEnvironment               │   └── ... (same as SystemLandscapeView*)
	    ├── DeploymentNode                  ├── FilteredView
	    │   ├── Tag                         │   ├── FilterTag
	    │   ├── Instances                   │   └── Exclude
	    │   ├── URL                         ├── DynamicView
	    │   ├── Prop                        │   ├── Title
	    │   └── DeploymentNode              │   ├── AutoLayout
	    │       └── ...                     │   ├── PaperSize
	    ├── InfrastructureNode              │   ├── Add
	    │   ├── Tag                         ├── DeploymentView
	    │   ├── URL                         │   └── ... (same as SystemLandscapeView*)
	    │   └── Prop                        └── Style
	    └── ContainerInstance                   ├── ElementStyle
	        ├── Tag                             └── RelationshipStyle
	        ├── HealthCheck
	        └── Prop                        (* minus EnterpriseBoundaryVisible)
*/
package dsl
