# Model

![Build](https://github.com/goadesign/model/workflows/CI/badge.svg)
![Version](https://img.shields.io/badge/Version-v1.6.2-blue.svg)
![Go version](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)
[![DSL Reference](https://img.shields.io/badge/Doc-DSL-orange)](https://pkg.go.dev/goa.design/model@v1.6.2/dsl?tab=doc)
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

Model includes a couple of code generation tools that generate diagrams from
the DSL:

* The `mdl` tool generates static web pages for all the views defined in the
  DSL. The tool can also serve the pages and includes live reload so that
  edits to the DSL are reflected in real-time making it convenient to author
  diagrams.

* The `stz` tool uploads the software architecture described via
  the DSL to the [Structurizr](https://structurizr.com) service. This service
  renders the model and includes a visual editor to rearrange the results
  (the tool takes care of keeping any change made graphically on the next
  upload).

Model also provides a [Goa](https://github.com/goadesign/goa) plugin so that
the design of APIs and microservices written in Goa can be augmented with a
description of the corresponding software architecture.

## Example

Here is a complete and correct DSL for an architecture model:

```Go
package design

import . "goa.design/model/dsl"

var _ = Design("Getting Started", "This is a model of my software system.", func() {
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

Static page rendering:

![Getting Started Diagram](https://raw.githubusercontent.com/goadesign/model/master/model.png)

Structurizr rendering:

![Getting Started Diagram](https://structurizr.com/static/img/getting-started.png)

Additional examples can be found in the
[examples](https://github.com/goadesign/model/tree/master/examples)
directory.

## Library

The [mdl](https://pkg.go.dev/goa.design/model@v1.6.2/mdl?tab=doc) package
[RunDSL](https://pkg.go.dev/goa.design/model@v1.6.2/mdl?tab=doc#RunDSL)
method runs the DSL and produces data structures that contain all the
information needed to render the views it defines including
[mermaid](https://mermaid-js.github.io) definitions for all the diagrams.

The [stz](https://pkg.go.dev/goa.design/model@v1.6.2/stz?tab=doc) package
[RunDSL](https://pkg.go.dev/goa.design/model@v1.6.2/stz?tab=doc#RunDSL)
method runs the DSL and produces a data structure that can be serialized into
JSON and uploaded to the [Structurizr service](https://structurizr.com).

The [stz](https://github.com/goadesign/model/tree/master/stz)
package also contains a client library for the
[Structurizr service APIs](https://structurizr.com/help/web-api).

Here is a complete example that takes advantage of both to upload the
workspace described in a DSL to the Structurizr service:

```Go
package main

import (
    "fmt"
    "os"

    . "goa.design/model/dsl"
    "goa.design/model/stz"
)

// DSL that describes software architecture model.
var _ = Design("Getting Started", "This is a model of my software system.", func() {
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
    w, err := stz.RunDSL()
    if err != nil {
        fmt.Fprintf(os.Stderr, "invalid design: %s", err.Error())
        os.Exit(1)
    }

    // Upload the design to the Structurizr service.
    // The API key and secret must be set in the STRUCTURIZR_KEY and
    // STRUCTURIZR_SECRET environment variables respectively. The
    // workspace ID must be set in STRUCTURIZR_WORKSPACE_ID.
    var (
        key    = os.Getenv("STRUCTURIZR_KEY")
        secret = os.Getenv("STRUCTURIZR_SECRET")
        wid    = os.Getenv("STRUCTURIZR_WORKSPACE_ID")
    )
    if key == "" || secret == "" || wid == "" {
        fmt.Fprintln(os.Stderr, "missing STRUCTURIZR_KEY, STRUCTURIZR_SECRET or STRUCTURIZR_WORKSPACE_ID environment variable.")
        os.Exit(1)
    }
    c := stz.NewClient(key, secret)
    if err := c.Put(wid, w); err != nil {
        fmt.Fprintf(os.Stderr, "failed to store workspace: %s\n", err.Error())
        os.Exit(1)
    }
}
```

## Tools

### mdl

The `mdl` tool can be used to render JSON files that contain all the
information needed to render the diagrams. The tool can also serve static
pages that reload when the DSL changes.

Generate the JSON files for the diagrams described in package
`goa.design/model/examples/basic`:

```bash
mdl gen goa.design/model/examples/basic
```

The command above created a `gen` folder containing one JSON file per view
defined in the DSL. In the case of the `basic` example there is a single
view. The data structure serialized into the JSON is
[RenderedView](https://pkg.go.dev/goa.design/model@v1.6.2/mdl?tab=doc#RenderedView).

Serve static pages produced from the same package:

```bash
mdl serve goa.design/model/examples/basic
[Model] listening on :6070
```

The pages can be browsed by visiting localhost:6070. Modifying and saving the
DSL will cause the page to reload and reflect the changes.

### stz

Alternatively, the `stz` tool included in this repo can be used to generate a
file containing the JSON representation of the design described via DSL. The
tool can then upload the generated file to the Structurizr service. The tool
can also retrieve the JSON representation of
[Workspace objects](https://github.com/structurizr/json) from the service.

Upload DSL defined in package `goa.design/model/examples/basic`:

```bash
stz gen goa.design/model/examples/basic && stz put model.json -id ID -key KEY -secret SECRET
```

Where `ID` is the Structurizr service workspace ID, `KEY` the
Structurizr service API key and `SECRET` the corresponding secret.

Retrieve the JSON representation of a workspace from the service:

```bash
stz get -id ID -key KEY -secret SECRET -out workspace.json
```

### Tools Setup

Assuming a working Go setup, run the following commands in the root of the
repo:

```bash
go install cmd/mdl
go install cmd/stz
```

This will create both a `mdl` and `stz` executables under `$GOPATH/bin` which
should be in your `PATH` environment variable.

## Goa Plugin

This package can also be used as a Goa plugin by including the DSL package in
the Goa design:

```Go
package design

import . "goa.design/goa/v3/dsl"
import "goa.design/model/dsl"

// ... DSL describing API, services and architecture model
```

Running `goa gen` creates a `model.json` file in the `gen` folder as well as
a `model` subdirectory. `model.json` follows the [structurizr JSON
schema](https://github.com/structurizr/json) and can be uploaded to the
Structurizr service for example using the `stz` tool included in this repo.
The `model` subdirectory contains the static rendered views JSON.

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

Some DSL functions such as `Uses`, `Delivers`, `InteractsWith`, `Add` and
`Link` accept references to elements as argument. The references can be done
either through a variable (which holds the element being referred to) or by
the path of the element. The path of an element is constructured by appending
the parent element names separated by slashes (/) and the name of the
element. For example the path of the component 'C' in the container 'CO' and
software system 'S' is 'S/CO/C'. The path can be relative when the reference
is made within a scoped function. For example when adding an element to a
view that is scoped to a parent element.

### Syntax

The DSL package
[documentation](https://pkg.go.dev/goa.design/model@v1.6.2/dsl?tab=doc) lists
all the DSL keywords and their usage.

The file [DSL.md](https://github.com/goadesign/model/blob/master/DSL.md)
illustrates the complete syntax in one design.

### Examples

Refer to the
[examples](https://github.com/goadesign/model/blob/master/examples) directory
for working examples.
