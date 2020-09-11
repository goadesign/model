/*
Package mdl generates all the information required to render diagrams from a
model design into static web pages. The information includes the Mermaid
source code for the diagram as well as information for each node that can be
used to further style the diagram, add links etc.

The package can be used in two different ways:

* As a library: the function Render produces data
  structures for each view in a design that include the source code for a
  Mermaid diagram and additional information useful for visualizing the
  diagrams (e.g. links, style classes and element properties).

* As a Goa plugin: by adding 'import _ "goa.design/model/plugin"' to the Goa
  design. Running 'goa gen' produces Mermaid diagrams in the 'gen/diagrams'
  directory and a JSON representation of a structurizr workspace in
  'gen/structurizr'.

The tool located in goa.design/model/cmd/mdl makes use of the library to
generate the information or pages directly given the import path to the Go
package containing the model design.
*/
package mdl
