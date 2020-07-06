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
    │   ├── Tag                             │   ├── AddDefault
    │   ├── URL                             │   ├── Add
    │   ├── External                        │   ├── AddAll
    │   ├── Properties                      │   ├── AddNeighbors
    │   ├── Uses                            │   ├── Remove
    │   └── InteractsWith                   │   ├── RemoveUnreachable
    ├── SoftwareSystem                      │   ├── RemoveUnrelated
    │   ├── Tag                             │   ├── AutoLayout
    │   ├── URL                             │   ├── AnimationStep
    │   ├── External                        │   ├── PaperSize
    │   ├── Properties                      │   └── EnterpriseBoundaryVisible
    │   ├── Uses                            ├── SystemContextView
    │   └── Delivers                        │   └──  ... (same as SystemLandsapeView)
    ├── Container                           ├── ContainerView
    │   ├── Tag                             │   ├── AddContainers
    │   ├── URL                             │   ├── AddInfluencers
    │   ├── Properties                      │   ├── SystemBoundariesVisible
    │   ├── Uses                            │   └── ... (same as SystemLandscapeView*)
    │   └── Delivers                        ├── ComponentView
    ├── Component                           │   ├── AddContainers
    │   ├── Tag                             │   ├── AddComponents
    │   ├── URL                             │   ├── ContainerBoundariesVisible
    │   ├── Properties                      │   └── ... (same as SystemLandscapeView*)
    │   ├── Uses                            ├── FilteredViee
    │   └── Delivers                        │   ├── FilterTag
    └── DeploymentEnvironment               │   └── Exclude
        ├── DeploymentNode                  ├── DynamicView
        │   ├── Tag                         │   ├── Title
        │   ├── Instances                   │   ├── AutoLayout
        │   ├── URL                         │   ├── PaperSize
        │   └── Properties                  │   └── Add
        ├── InfrastructureNode              ├── DeploymentView
        │   ├── Tag                         │   └── ... (same as SystemLandscapeView*)
        │   ├── URL                         ├── Style
        │   └── Properties                  │   ├── ElementStyle
        └── ContainerInstance               │   └── RelationshipStyle
            ├── Tag                         ├── Theme
            ├── InstanceID                  └── Branding
            ├── HealthCheck
            └── Properties                  (* minus EnterpriseBoundaryVisible)
*/
package dsl
