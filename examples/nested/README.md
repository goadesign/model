# Nested Example

This examples illustrates how multiple packages can be contribute to the same
model. This makes it possible for multiple teams to collaborate on the same
overall software architecture model while maintaining only the section
relevant to them.

The example also illustrates how styles may be shared between multiple models
by defining them in a shared Go package.

The nesting is done simply by leveraging the Go `import` facility: the model
being generated imports the packages that defines the nested models. The
parent model can refer to elements defined in the child models by using the
functions exposed on the corresponding element structs:

* The `Person` and `SoftwareSystem` workspace methods return a person or
  software system element given their name.
* The `Container` software system method returns a container given its name.
* The `Component` container method returns a component given its name.
* The `DeploymentNode` workspace method returns a deployment node given its
  name, the `Child` deployment node method returns a child deployment node
  given its name.
* The `InfrastructureNode` deployment node method returns an infrastructure
  node given its name.
* The `ContainerInstance` deployment node method returns a container instance
  given a reference to the container being instantiated (or its name) and the
  corresponding instance ID (1 if there is only one instance of the container
  in the deployment node).

When used that way the `Workspace` expression defined in imported models gets
overridden by the one defined in the package being generated.

## Running

Using the `stz` command line:

```bash
stz gen goa.design/model/examples/nested/model
```

This generates the file `model.json` that contains the JSON representation of
the workspace.  This file can be used to upload the workspace, again using stz:

```bash
stz put model.json -workspace XXX -key YYY -secret ZZZ
```
