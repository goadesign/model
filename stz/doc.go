/*
Package stz defines data structures that represent a software architecture
model and accompanying views that follow the C4 model (https://c4model.com).

The data structures can be serialized into JSON. The top level data structure
is Workspace which defines the model and views as well as a name, description
and version for the design. Model describes the people, software systems,
containers and components that make up the architecture as well as the
deployment nodes that represent runtime deployments. Views describes diagrams
that represent different level of details - from contextual views that
represent the overall system in context with other systems to component level
views that render software components and their relationships.

The JSON representation of the workspace data structure is compatible with
the Structurizr service API (https://structurizr.com). The package also
includes a Structurize API client that can be used to upload and download
workspaces to and from the service.
*/
package stz
