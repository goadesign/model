# Basic Model Example

This example `model` package contains a valid DSL describing the example used
in the README of the repo.

## Usage

### Using the example command

The example contains a command that uploads the diagram to the Structurizr
service. The `main` function loads the Structurizr workspace ID, API key and
API secret from the environment:

* `$STRUCTURIZR_WORKSPACE_ID`: Workspace ID
* `$STRUCTURIZR_KEY`: API key
* `$STRUCTURIZR_SECRET`: API secret

Follow the steps below to run the command in `bash` (substitute the values
between brackets):

```bash
cd $GOPATH/src/goa.design/model/examples/basic
export STRUCTURIZR_WORKSPACE_ID="<Workspace ID>"
export STRUCTURIZR_KEY="<API key>"
export STRUCTURIZR_SECRET="<API secret>"
go run main.go
```

Open the diagram in a browser:

```bash
open https://structurizr.com/workspace/$STRUCTURUZR_WORKSPACE_ID/diagrams#SystemContext
```

## Using the `mdl` tool

Alternatively the `mdl` tool can be used to render the diagram locally. Make sure the
tool is installed:

```bash
mdl version
```

If the command above returns an error then try reinstalling the tool:

```bash
go install goa.design/model/cmd/mdl
```

Serve the static page using the tool:

```bash
mdl serve goa.design/model/examples/basic/model
```

Open the diagram in a browser:

```bash
open http://localhost:6070/SystemContext
```
