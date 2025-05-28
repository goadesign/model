
# DSL Syntax

The code snippet below describes the entire syntax of the DSL. The complete
reference can be found in the `dsl`
[package documentation](https://pkg.go.dev/goa.design/model@v1.11.1/dsl?tab=doc)

```Go
// Design defines the architecture design containing the models and views.
// Design must appear exactly once.
var _ = Design("[name]", "[description]", func() {

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
        Prop("<name>", "<value>")

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
        Prop("<name>", "<value>")

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

            // RemoveTagged removes elements and relationships with the given tag.
            RemoveTagged("<tag>")

            // Remove given relationship from view.
            Unlink(Source, Destination)

            // Remove all elements and people that cannot be reached by
            // traversing the relashionships starting with given element or
            // person.
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

            // ElementStyle defines an element style.
            ElementStyle("<tag>", func() {
                Shape(ShapeBox) // ShapeBox, ShapeRoundedBox, ShapeCircle, ShapeEllipse,
                                // ShapeHexagon, ShapeCylinder, ShapePipe, ShapePerson
                                // ShapeRobot, ShapeFolder, ShapeWebBrowser,
                                // ShapeMobileDevicePortrait, ShapeMobileDeviceLandscape,
                                // ShapeComponent.
                Width(42)
                Height(42)
                FontSize(42)
                Border(BorderSolid) // BorderSolid, BorderDashed, BorderDotted
                Opacity(42)         // Between 0 and 100
                Icon("<url>")
                Background("#<rrggbb>")
                Color("#<rrggbb>")
                Stroke("#<rrggbb>")
                ShowMetadata()
                ShowDescription()
            })

            // RelationshipStyle defines a relationship style. All nested
            // properties (thickness, color, etc) are optional.
            RelationshipStyle("<tag>", func() {
                Thickness(42)
                FontSize(42)
                Width(42)
                Position(42) // Between 0 and 100
                Opacity(42)  // Between 0 and 100
                Color("#<rrggbb>")
                Solid()
                Dashed()
                Routing(RoutingDirect) // RoutingDirect, RoutingOrthogonal, RoutingCurved
            })
        })
    })
})
```
