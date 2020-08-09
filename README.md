# Model

![Build](https://github.com/goadesign/model/workflows/CI/badge.svg)
![Version](https://img.shields.io/badge/Version-v1.0.2-blue.svg)
![Go version](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)
[![DSL Reference](https://img.shields.io/badge/Doc-DSL-orange)](https://pkg.go.dev/goa.design/model@v1.0.2/dsl?tab=doc)
[![Go Packages](https://img.shields.io/badge/Doc-packages-orange)](https://pkg.go.dev/goa.design/model)

## Overview

Model provides a way to describe software architectures using
*diagrams as code*. This approach provides many benefit over the use of
 graphical tools, in particular:

* **Built-in history via source control versioning**: Each change to the software
  architecture model results in a commit making it natural to implement
  [Architectural Decision Records](https://adr.github.io/).
* **Simple and consistent reuse of shared architecture components**: The code
  of shared architecture components can be imported by other software
  systems. Apart from the obvious time saving a major advantage is ensuring
  that all references are automatically kept up-to-date as the shared
  architecture component evolves. The shared architecture component code can
  be versioned using traditional source code versioning techniques as well.
* **Consistent notation and visual styles**: The same mechanism used to import
  shared architecture diagram can also be used to import shared style
  definitions.
* **Ability to control drift between reality and diagrams**: Static code
  analysis allows writing tools that compare the software models with actual
  code to detect discrepencies (this repo does not provide such a tool at
  this time).

The Model DSL is implemented in [Go](https://golang.org/) and follows the
[C4 Model](https://c4model.com) to describe the software architecture. Using
Go to implement the DSL makes it possible to leverage Go packages to share
and version models. It also allows for extending or customizing the DSL by
writing simple Go functions.

The C4 (Context, Containers, Components and Code) model provides a clean and
simple set of constructs (software elements, deployment nodes and views) that
can be learned in minutes. The C4 model is very flexible and only focuses on
a few key concepts making it possible to express many different styles of
architectures while still adding value.

Model includes a tool that uploads the software architecture described via
its DSL to the [Structurizr](https://structurizr.com) service. This service
renders the model and includes a visual editor to rearrange the results
(layout modifications done visually are kept on the next uploads).

Other outputs (such as [Mermaid](https://mermaid-js.github.io/mermaid/#/))
could be added in the future (consider contributing!).

Model also provides a [Goa](https://github.com/goadesign/goa) plugin so that
the design of APIs and microservices written in Goa can be augmented with a
description of the corresponding software architecture.

## Example

Here is a complete and correct DSL for an architecture model:

```Go
package model

import . "goa.design/model/dsl"

var _ = Workspace("Getting Started", "This is a model of my software system.", func() {
    var System = SoftwareSystem("Software System", "My software system.", func() {
        Tag("system")
    })

    Person("User", "A user of my software system.", func() {
        Uses(System, "Uses")
        Tag("person")
    })

    Views(func() {
        SystemContextView(System, "SystemContext", "An example of a System Context diagram.", func() {
            AddAll()
            AutoLayout(RankLeftRight)
        })
        Styles(func() {
            ElementStyle("system", func() {
                Background("#1168bd")
                Color("#ffffff")
            })
            ElementStyle("person", func() {
                Shape(ShapePerson)
                Background("#08427b")
                Color("#ffffff")
            })
        })
    })
})
```

This code creates a model containing elements and relationships, creates a
single view and adds some styling.
![Getting Started Diagram](https://structurizr.com/static/img/getting-started.png)

Other examples can be found in the repo
[examples](https://github.com/goadesign/model/tree/master/examples)
directory.

## Library

The [eval](https://github.com/goadesign/model/tree/master/eval) package
makes it convenient to run the DSL above. The
[service](https://github.com/goadesign/model/tree/master/service)
package contains a client library for the
[Structurizr service APIs](https://structurizr.com/help/web-api).

Here is a complete example that takes advantage of both to upload the
workspace described in a DSL to the Structurizr service:

```Go
package main

import (
    "fmt"
    "os"

    . "goa.design/model/dsl"
    "goa.design/model/eval"
    "goa.design/model/service"
)

// DSL that describes software architecture model.
var _ = Workspace("Getting Started", "This is a model of my software system.", func() {
    var System = SoftwareSystem("Software System", "My software system.", func() {
        Tag("system")
    })

    Person("User", "A user of my software system.", func() {
        Uses(System, "Uses")
        Tag("person")
    })

    Views(func() {
        SystemContextView(System, "SystemContext", "An example of a System Context diagram.", func() {
            AddAll()
            AutoLayout(RankLeftRight)
        })
        Styles(func() {
            ElementStyle("system", func() {
                Background("#1168bd")
                Color("#ffffff")
            })
            ElementStyle("person", func() {
                Shape(ShapePerson)
                Background("#08427b")
                Color("#ffffff")
            })
        })
    })
})

// Executes the DSL and uploads the corresponding workspace to Structurizr.
func main() {
    // Run the model DSL
    w, err := eval.RunDSL()
    if err != nil {
        fmt.Fprintf(os.Stderr, "invalid model: %s", err.Error())
        os.Exit(1)
    }

    // Upload the model to the Structurizr service.
    // The API key and secret must be set in the STRUCTURIZR_KEY and
    // STRUCTURIZR_SECRET environment variables respectively. The
    // workspace ID must be set in STRUCTURIZR_WORKSPACE_ID.
    var (
        key    = os.Getenv("STRUCTURIZR_KEY")
        secret = os.Getenv("STRUCTURIZR_SECRET")
        wid    = os.Getenv("STRUCTURIZR_WORKSPACE_ID")
    )
    c := service.NewClient(key, secret)
    if err := c.Put(wid, w); err != nil {
        fmt.Fprintf(os.Stderr, "failed to store workspace: %s", err.Error())
        os.Exit(1)
    }
}
```

## Tool

Alternatively, the `stz` tool included in this repo can be used to generate a
file containing the JSON representation of the structurizr API
[Workspace object](https://github.com/structurizr/json) described via DSL.
The tool can can also retrieve or upload such files from and to the
Structurizr service. Finally the tool can also lock or unlock a workspace in
the service.

Upload DSL defined in package `goa.design/model/examples/basic`:

```bash
stz gen goa.design/model/examples/basic && stz put -id ID -key KEY -secret SECRET
```

Where `ID` is the Structurizr service workspace ID, `KEY` the
Structurizr service API key and `SECRET` the corresponding secret.

Retrieve the JSON representation of a workspace from the service:

```bash
stz get -id ID -key KEY -secret SECRET -out model.json
```

Upload an existing file to the Structurizr service:

```bash
stz put model.json -id ID -key KEY -secret SECRET
```

### Tool Setup

Assuming a working Go setup, run the following command in the root of the
repo:

```bash
go install cmd/stz
```

This will create a `stz` executable under `$GOPATH/bin` which should be in
your `PATH` environment variable.

## Goa Plugin

This package can also be used as a Goa plugin by including the DSL package in
the Goa design:

```Go
package design

import . "goa.design/goa/v3/dsl"
import "goa.design/model/dsl"

// ... DSL describing API, services and architecture model
```

Running `goa gen` creates a `model.json` file in the `gen` folder. This
file follows the
[structurizr JSON schema](https://github.com/structurizr/json) and can be
uploaded to the Structurizr service for example using the `stz` tool included
in this repo.

## DSL Syntax

### Rules

The following rules apply to all elements and views declared in a model:

* Software and people names must be unique.
* Container names must be unique within the context of a software system.
* Component names must be unique within the context of a container.
* Deployment node names must be unique with their parent context.
* Infrastructure node names must be unique with their parent context.
* All relationships from a given source element to a given destination element
  must have a unique description.
* View keys must be unique.

Note that uniqueness of names is enforced by combining the evaluated
definitions for a given element. For example if a model contained:

```Go
var Person1 = Person("User", "A user", func() {
    Uses("System 1", "Uses")
})
var Person2 = Person("User", "The same user again", func() {
    Uses("System 2", "Uses")
})
```

Then the final model would only define a single person named `"User"` with
the description `"The same user again"` and both relationships. This makes it
possible to import shared models and "edit" existing elements, for example to
add new relationships.

### References

The functions `Uses`, `Delivers`, `InteractsWith`, `Add` and `Link` accept
references to other elements as argument. The reference can be done either
through a variable (which holds the element being referred to) or the name of
the element. Note that names do not necessarily have to be globally unique
(see rules above) so it may sometimes be necessary to use a variable to
disambiguate. Also container instances do not have names per se however the
`RefName` function makes it possible to define a name that can be used to
refer to the container instance in deployment views (when using `Add` or
`Link`).

### Syntax

The code snippet below describes the entire syntax of the DSL. The complete
reference can be found in the `dsl`
[package documentation](https://pkg.go.dev/goa.design/model@v1.0.2/dsl?tab=doc)

```Go
// Workspace defines the workspace containing the models and views. Workspace
// must appear exactly once in a given design. A name must be provided if a
// description is.
var _ = Workspace("[name]", "[description]", func() {

    // Version number.
    Version("<version>")

    // Enterprise defines a named "enterprise" (e.g. an organisation). On System
    // Landscape and System Context diagrams, an enterprise is represented as a
    // dashed box. Only a single enterprise can be defined within a model.
    Enterprise("<name>")

    // Person defines a person (user, actor, role or persona).
    var Person = Person("<name>", "[description]", func() {
        Tag("<name>", "[name]") // as many tags as needed

        // URL where more information about this system can be found.
        URL("<url>")

        // External indicates the person is external to the enterprise.
        External()

        // Prop defines an arbitrary set of associated key-value pairs.
        Prop("<name>", "<value">)

        // Adds a uni-directional relationship between this person and the given element.
        Uses(Element, "<description>", "[technology]", Synchronous /* or Asynchronous */, func() {
            Tag("<name>", "[name]") // as many tags as needed
        })

        // Adds an interaction between this person and another.
        InteractsWith(Person, "<description>", "[technology]", Synchronous /* or Asynchronous */, func() {
            Tag("<name>", "[name]") // as many tags as needed
        })
    })

    // SoftwareSystem defines a software system.
    var SoftwareSystem = SoftwareSystem("<name>", "[description]", func() {
        Tag("<name>",  "[name]") // as many tags as needed

        // URL where more information about this software system can be
        // found.
        URL("<url>")

        // External indicates the software system is external to the enterprise.
        External()

        // Prop defines an arbitrary set of associated key-value pairs.
        Prop("<name>", "<value">)

        // Adds a uni-directional relationship between this software system and the given element.
        Uses(Element, "<description>", "[technology]", Synchronous /* or Asynchronous */, func() {
            Tag("<name>", "[name]") // as many tags as needed
        })

        // Adds an interaction between this software system and a person.
        Delivers(Person, "<description>", "[technology]", Synchronous /* or Asynchronous */, func() {
            Tag("<name>", "[name]") // as many tags as needed
        })

        // Container defines a container within a software system.
        var Container = Container("<name>",  "[description]",  "[technology]",  func() {
            Tag("<name>",  "[name]") // as many tags as neede

            // URL where more information about this container can be found.
            URL("<url>")

            // Prop defines an arbitrary set of associated key-value pairs.
            Prop("<name>", "<value">)

            // Adds a uni-directional relationship between this container and the given element.
            Uses(Element, "<description>", "[technology]", Synchronous /* or Asynchronous */, func () {
                Tag("<name>", "[name]") // as many tags as needed
            })

            // Adds an interaction between this container and a person.
            Delivers(Person, "<description>", "[technology]", Synchronous /* or Asynchronous */, func() {
                Tag("<name>", "[name]") // as many tags as needed
            })
        })

        // Container may also refer to a Goa service in which case the name
        // and description are taken from the given service definition and
        // the technology is set to "Go and Goa v3".
        var Container = Container(GoaService, func() {
            // ... see above

            // Component defines a component within a container.
            var Component = Component("<name>",  "[description]",  "[technology]",  func() {
                Tag("<name>",  "[name]") // as many tags as need
                // URL where more information about this container can be found.
                URL("<url>")
                // Prop defines an arbitrary set of associated key-value pairs.
                Prop("<name>", "<value">)
                // Adds a uni-directional relationship between this component and the given element.
                Uses(Element, "<description>", "[technology]", Synchronous /* or Asynchronous */, func() {
                    Tag("<name>", "[name]") // as many tags as needed
                })
                // Adds an interaction between this component and a person.
                Delivers(Person, "<description>", "[technology]", Synchronous /* or Asynchronous */, func() {
                    Tag("<name>", "[name]") // as many tags as needed
                })
            })
        })
    })

    // DeploymentEnvironment provides a way to define a deployment
    // environment (e.g. development, staging, production, etc).
    DeploymentEnvironment("<name>", func() {

        // DeploymentNode defines a deployment node. Deployment nodes can be
        // nested, so a deployment node can contain other deployment nodes.
        // A deployment node can also contain InfrastructureNode and
        // ContainerInstance elements.
        var DeploymentNode = DeploymentNode("<name>",  "[description]",  "[technology]",  func() {
            Tag("<name>",  "[name]") // as many tags as needed

            // Instances sets the number of instances, defaults to 1.
            Instances(2)

            // URL where more information about this deployment node can be
            // found.
            URL("<url>")

            // Prop defines an arbitrary set of associated key-value pairs.
            Prop("<name>", "<value">)

            // InfrastructureNode defines an infrastructure node, typically
            // something like a load balancer, firewall, DNS service, etc.
            var InfrastructureNode = InfrastructureNode("<name>", "[description]", "[technology]", func() {
                Tag("<name>",  "[name]") // as many tags as needed

                // URL where more information about this infrastructure node can be
                // found.
                URL("<url>")

                // Prop defines an arbitrary set of associated key-value pairs.
                Prop("<name>", "<value">)
            })

            // ContainerInstance defines an instance of the specified
            // container that is deployed on the parent deployment node.
            var ContainerInstance = ContainerInstance(Container, func() {
                Tag("<name>",  "[name]") // as many tags as needed

                // Sets a name that can be used in deployment views.
                RefName("<name>")

                // Sets instance number or index.
                InstanceID(1)

                // Prop defines an arbitrary set of associated key-value pairs.
                Prop("<name>", "<value">)

                // HealthCheck defines a HTTP-based health check for this
                // container instance.
                HealthCheck("<name>", func() {

                    // URL is the health check URL/endpoint.
                    URL("<url>")

                    // Interval is the polling interval, in seconds.
                    Interval(42)

                    // Timeout after which a health check is deemed as failed,
                    // in milliseconds.
                    Timeout(42)

                    // Header defines a header that should be sent with the
                    // request.
                    Header("<name>", "<value>")
                })
            })

            // DeploymentNode within a deployment node defines child nodes.
            var ChildNode = DeploymentNode("<name>", "[description]", "[technology]", func() {
                // ... see above
            })
        })
    })

    // Views is optional and defines one or more views.
    Views(func() {

        // SystemLandscapeView defines a System Landscape view.
        SystemLandscapeView("[key]", "[description]", func() {

            // Title of this view.
            Title("<title>")

            // AddDefault adds default elements that are relevant for the
            // specific view:
            //
            //    - System landscape view: adds all software systems and people
            //    - System context view: adds softare system and other related
            //      software systems and people.
            //    - Container view: adds all containers in software system as well
            //      as related software systems and people.
            //    - Component view: adds all components in container as well as
            //      related containers, software systems and people.
            AddDefault()

            // Add given person or element to view. If person or element was
            // already added implictely (e.g. via AddAll()) then overrides how
            // the person or element is rendered.
            Add(PersonOrElement, func() {

                // Set explicit coordinates for where to render person or
                // element.
                Coord(X, Y)

                // Do not render relationships when rendering person or element.
                NoRelationship()
            })

            // Add given relationship to view. If relationship was already added
            // implictely (e.g. via AddAll()) then overrides how the
            // relationship is rendered.
            Link(Source, Destination, func() {

                // Vertices lists the x and y coordinate of the vertices used to
                // render the relationship. The number of arguments must be even.
                Vertices(10, 20, 10, 40)

                // Routing algorithm used when rendering relationship, one of
                // RoutingDirect, RoutingCurved or RoutingOrthogonal.
                Routing(RoutingOrthogonal)

                // Position of annotation along line; 0 (start) to 100 (end).
                Position(50)
            })

            // Add all elements and people in scope.
            AddAll()

            // Add default set of elements depending on view type.
            AddDefault()

            // Add all elements that are directly connected to given person
            // or element.
            AddNeighbors(PersonOrElement)

            // Remove given element or person from view.
            Remove(ElementOrPerson)

            // Remove given relationship from view.
            Remove(Source, Destination)

            // Remove elements and relationships with given tag.
            Remove("<tag>")

            // Remove all elements and people that cannot be reached by
            // traversing the graph of relationships starting with given element
            // or person.
            RemoveUnreachable(ElementOrPerson)

            // Remove all elements that have no relationships to other elements.
            RemoveUnrelated()

            // AutoLayout enables automatic layout mode for the diagram. The
            // first argument indicates the rank direction, it must be one of
            // RankTopBottom, RankBottomTop, RankLeftRight or RankRightLeft.
            AutoLayout(RankTopBottom, func() {

                // Separation between ranks in pixels, defaults to 300.
                RankSeparation(300)

                // Separation between nodes in the same rank in pixels, defaults to 600.
                NodeSeparation(600)

                // Separation between edges in pixels, defaults to 200.
                EdgeSeparation(200)

                // Create vertices during automatic layout, false by default.
                RenderVertices()
            })

            // Animation defines an animation step consisting of the
            // specified elements.
            Animation(Element, Element/*, ...*/)

            // PaperSize defines the paper size that should be used to render
            // the view. The possible values for the argument follow the
            // patterns SizeA[0-6][Portrait|Landscape], SizeLetter[Portrait|Landscape]
            // or SizeLegal[Portrait_Landscape]. Alternatively the argument may be
            // one of SizeSlide4X3, SizeSlide16X9 or SizeSlide16X10.
            PaperSize(SizeSlide4X3)

            // Make enterprise boundary visible to differentiate internal
            // elements from external elements on the resulting diagram.
            EnterpriseBoundaryVisible()
        })

        SystemContextView(SoftwareSystem, "[key]", "[description]", func() {
            // ... same usage as SystemLandscapeView.
        })

        ContainerView(SoftwareSystem, "[key]", "[description]", func() {
            // ... same usage as SystemLandscapeView without EnterpriseBoundaryVisible.

            // All all containers in software system to view.
            AddContainers()

            // Make software system boundaries visible for "external" containers
            // (those outside the software system in scope).
            SystemBoundariesVisible()
        })

        ComponentView(Container, "[key]", "[description]", func() {
            // ... same usage as SystemLandscapeView without EnterpriseBoundaryVisible.

            // All all containers in software system to view.
            AddContainers()

            // All all components in container to view.
            AddComponents()

            // Make container boundaries visible for "external" components
            // (those outside the container in scope).
            ContainerBoundariesVisible()
        })

        // FilteredView defines a Filtered view on top of the specified view.
        // The given view must be a System Landscape, System Context, Container,
        // or Component view on which this filtered view should be based.
        FilteredView(View, func() {
            // Set of tags to include or exclude (if Exclude() is used)
            // elements/relationships when rendering this filtered view.
            FilterTag("<tag>", "[tag]") // as many as needed

            // Exclude elements and relationships with the given tags instead of
            // including.
            Exclude()
        })

        // DynamicView defines a Dynamic view for the specified scope. The
        // first argument defines the scope of the view, and therefore what can
        // be added to the view, as follows:
        //
        //   * Global scope: People and software systems.
        //   * Software system scope: People, other software systems, and
        //     containers belonging to the software system.
        //   * Container scope: People, other software systems, other
        //     containers, and components belonging to the container.
        DynamicView(Global, "[key]", "[description]", func() {

            // Title of this view.
            Title("<title>")

            // AutoLayout enables automatic layout mode for the diagram. The
            // first argument indicates the rank direction, it must be one of
            // RankTopBottom, RankBottomTop, RankLeftRight or RankRightLeft.
            AutoLayout(RankTopBottom, func() {

                // Separation between ranks in pixels
                RankSeparation(200)

                // Separation between nodes in the same rank in pixels
                NodeSeparation(200)

                // Separation between edges in pixels
                EdgeSeparation(10)

                // Create vertices during automatic layout.
                Vertices()
            })

            // PaperSize defines the paper size that should be used to render
            // the view. The possible values for the argument follow the
            // patterns SizeA[0-6][Portrait|Landscape], SizeLetter[Portrait|Landscape]
            // or SizeLegal[Portrait_Landscape]. Alternatively the argument may be
            // one of SizeSlide4X3, SizeSlide16X9 or SizeSlide16X10.
            PaperSize(SizeSlide4X3)

            // Set of relationships that make up dynamic diagram.
            Link(Source, Destination, func() {

                // Vertices lists the x and y coordinate of the vertices used to
                // render the relationship. The number of arguments must be even.
                Vertices(10, 20, 10, 40)

                // Routing algorithm used when rendering relationship, one of
                // RoutingDirect, RoutingCurved or RoutingOrthogonal.
                Routing(RoutingOrthogonal)

                // Position of annotation along line; 0 (start) to 100 (end).
                Position(50)

                // Description used in dynamic views.
                Description("<description>")

                // Order of relationship in dynamic views, e.g. 1.0, 1.1, 2.0
                Order("<order>")
            })
        })

        // DynamicView on software system or container uses the corresponding
        // identifier as first argument.
        DynamicView(SoftwareSystemOrContainer, "[key]", "[description]", func() {
            // see usage above
        })

        // DeploymentView defines a Deployment view for the specified scope and
        // deployment environment. The first argument defines the scope of the
        // view, and the second property defines the deployment environment. The
        // combination of these two arguments determines what can be added to
        // the view, as follows:
        //  * Global scope: All deployment nodes, infrastructure nodes, and
        //    container instances within the deployment environment.
        //  * Software system scope: All deployment nodes and infrastructure
        //    nodes within the deployment environment. Container instances within
        //    the deployment environment that belong to the software system.
        DeploymentView(Global, "<environment name>", "[key]", "[description]", func() {
            // ... same usage as SystemLandscape without EnterpriseBoundaryVisible.
        })

        // DeploymentView on a software system uses the software system as first
        // argument.
        DeploymentView(SoftwareSystem, "<environment name>", "[key]", "[description]", func() {
            // see usage above
        })

        // Styles is a wrapper for one or more element/relationship styles,
        // which are used when rendering diagrams.
        Styles(func() {

            // ElementStyle defines an element style. All nested properties
            // (shape, icon, etc) are optional.
            ElementStyle("<tag>", func() {
                Shape(ShapeBox) // ShapeBox, ShapeRoundedBox, ShapeCircle, ShapeEllipse,
                                // ShapeHexagon, ShapeCylinder, ShapePipe, ShapePerson
                                // ShapeRobot, ShapeFolder, ShapeWebBrowser,
                                // ShapeMobileDevicePortrait, ShapeMobileDeviceLandscape,
                                // ShapeComponent.
                Icon("<file>")
                Width(42)
                Height(42)
                Background("#<rrggbb>")
                Color("#<rrggbb>")
                Stroke("#<rrggbb>")
                FontSize(42)
                Border(BorderSolid) // BorderSolid, BorderDashed, BorderDotted
                Opacity(42) // Between 0 and 100
                ShowMetadata()
                ShowDescription()
            })

            // RelationshipStyle defines a relationship style. All nested
            // properties (thickness, color, etc) are optional.
            RelationshipStyle("<tag>", func() {
                Thickness(42)
                Color("#<rrggbb>")
                Solid()
                Routing(RoutingDirect) // RoutingDirect, RoutingOrthogonal, RoutingCurved
                FontSize(42)
                Width(42)
                Position(42) // Between 0 and 100
                Opacity(42)  // Between 0 and 100
            })
        })

        // Theme specifies one or more themes that should be used when
        // rendering diagrams.
        Theme("<theme URL>", "[theme URL]") // as many theme URLs as needed

        // Branding defines custom branding that should be used when rendering
        // diagrams and documentation.
        Branding(func() {
            Logo("<file>")
            Font("<name>", "[url]")
        })
    })
})
```
