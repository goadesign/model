/*
Package design defines data structures that represent a software architecture
model and accompanying views that follow the C4 model (https://c4model.com).

The data structures can be serialized into JSON. The top level data structure
is Design which defines the model and views as well as a name, description
and version for the design. Model describes the people, software systems,
containers and components that make up the architecture as well as the
deployment nodes that represent runtime deployments. Views describes diagrams
that represent different level of details - from contextual views that
represent the overall system in context with other systems to component level
views that render software components and their relationships.
*/
package design
