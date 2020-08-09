/*
Package dsl implements a Go based DSL that makes it possible to describe
softare architecture models following the C4 model (https://c4model.com).

It is recommended to use "dot import" when using this package to write DSLs,
for example:

    package model

    import . "goa.design/model/dsl"

    var _ = Workspace("<name>", "[description]", func() {
        // ...
    })

Some DSL functions accept a anonymous function as last argument (such as
Workspace above) which makes it possible to define a nesting structure. The
general shape of the DSL is:

    Workspace                           Workspace
    ├── Version                         └── Views
    ├── Enterprise                          ├── SystemLandscapeView
    ├── Person                              │   ├── Title
    │   ├── Tag                             │   ├── AddDefault
    │   ├── URL                             │   ├── Add
    │   ├── External                        │   ├── AddAll
    │   ├── Prop                            │   ├── AddNeighbors
    │   ├── Uses                            │   ├── Link
    │   └── InteractsWith                   │   ├── Remove
    ├── SoftwareSystem                      │   ├── RemoveUnreachable
    │   ├── Tag                             │   ├── RemoveUnrelated
    │   ├── URL                             │   ├── AutoLayout
    │   ├── External                        │   ├── Animation
    │   ├── Prop                            │   ├── PaperSize
    │   ├── Uses                            │   └── EnterpriseBoundaryVisible
    │   └── Delivers                        ├── SystemContextView
    ├── Container                           │   └──  ... (same as SystemLandsapeView)
    │   ├── Tag                             ├── ContainerView
    │   ├── URL                             │   ├── AddContainers
    │   ├── Prop                            │   ├── AddInfluencers
    │   ├── Uses                            │   ├── SystemBoundariesVisible
    │   └── Delivers                        │   └── ... (same as SystemLandscapeView*)
    ├── Component                           ├── ComponentView
    │   ├── Tag                             │   ├── AddContainers
    │   ├── URL                             │   ├── AddComponents
    │   ├── Prop                            │   ├── ContainerBoundariesVisible
    │   ├── Uses                            │   └── ... (same as SystemLandscapeView*)
    │   └── Delivers                        ├── FilteredView
    └── DeploymentEnvironment               │   ├── FilterTag
        ├── DeploymentNode                  │   └── Exclude
        │   ├── Tag                         ├── DynamicView
        │   ├── Instances                   │   ├── Title
        │   ├── URL                         │   ├── AutoLayout
        │   └── Prop                        │   ├── PaperSize
        ├── InfrastructureNode              │   └── Add
        │   ├── Tag                         ├── DeploymentView
        │   ├── URL                         │   └── ... (same as SystemLandscapeView*)
        │   └── Prop                        ├── Style
        └── ContainerInstance               │   ├── ElementStyle
            ├── Tag                         │   └── RelationshipStyle
            ├── HealthCheck                 ├── Theme
            ├── Prop                        └── Branding
            └── RefName                     (* minus EnterpriseBoundaryVisible)
*/
package dsl
