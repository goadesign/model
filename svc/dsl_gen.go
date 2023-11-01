package svc

import (
	"fmt"
	"strings"
	"text/template"

	geneditor "goa.design/model/svc/gen/dsl_editor"
)

type (
	// RelationshipData contains the data needed to render a relationship.
	RelationshipData struct {
		*geneditor.Relationship
		RelationName string
	}
)

var (
	systemT            = template.Must(template.New("system").Parse(systemTemplate))
	personT            = template.Must(template.New("person").Parse(personTemplate))
	containerT         = template.Must(template.New("container").Parse(containerTemplate))
	componentT         = template.Must(template.New("component").Parse(componentTemplate))
	relationshipT      = template.Must(template.New("relationship").Parse(relationshipTemplate))
	landscapeViewT     = template.Must(template.New("landscapeView").Parse(landscapeViewTemplate))
	systemContextViewT = template.Must(template.New("systemContextView").Parse(systemContextViewTemplate))
	containerViewT     = template.Must(template.New("containerView").Parse(containerViewTemplate))
	componentViewT     = template.Must(template.New("componentView").Parse(componentViewTemplate))
	elementStyleT      = template.Must(template.New("elementStyle").Parse(elementStyleTemplate))
	relationshipStyleT = template.Must(template.New("relationshipStyle").Parse(relationshipStyleTemplate))
)

func systemDSL(s *geneditor.System) string                       { return exec(systemT, s) }
func personDSL(p *geneditor.Person) string                       { return exec(personT, p) }
func containerDSL(c *geneditor.Container) string                 { return exec(containerT, c) }
func componentDSL(c *geneditor.Component) string                 { return exec(componentT, c) }
func relationshipDSL(r *RelationshipData) string                 { return exec(relationshipT, r) }
func landscapeViewDSL(v *geneditor.LandscapeView) string         { return exec(landscapeViewT, v) }
func systemContextViewDSL(v *geneditor.SystemContextView) string { return exec(systemContextViewT, v) }
func containerViewDSL(v *geneditor.ContainerView) string         { return exec(containerViewT, v) }
func componentViewDSL(v *geneditor.ComponentView) string         { return exec(componentViewT, v) }
func elementStyleDSL(v *geneditor.ElementStyle) string           { return exec(elementStyleT, v) }
func relationshipStyleDSL(v *geneditor.RelationshipStyle) string { return exec(relationshipStyleT, v) }

// exec executes the given template with the given data and returns the
// result.
func exec(t *template.Template, data interface{}) string {
	var b strings.Builder
	if err := t.Execute(&b, data); err != nil {
		panic(fmt.Sprintf("failed to execute template: %s", err)) // should never happen
	}
	return b.String()
}

const systemTemplate = `SoftwareSystem({{ printf "%q" .Name }}{{ if .Description }}, {{ printf "%q" .Description }}{{ end }}, func() {
	{{- range .Tags }}
        Tag({{ printf "%q" . }})
	{{- end }}
	{{- if .URL }}
	URL({{ printf "%q" .URL }})
	{{- end }}
	{{- if eq .Location "External" }}
	External()
	{{- end }}
	{{- range .Properties }}
	Prop({{ printf "%q" .Key }}, {{ printf "%q" .Value }})
	{{- end }}
})`

const personTemplate = `Person({{ printf "%q" .Name }}{{ if .Description }}, {{ printf "%q" .Description }}{{ end }}, func() {
	{{- range .Tags }}
	Tag({{ printf "%q" . }})
	{{- end }}
	{{- if .URL }}
	URL({{ printf "%q" .URL }})
	{{- end }}
	{{- if eq .Location "External" }}
	External()
	{{- end }}
	{{- range .Properties }}
	Prop({{ printf "%q" .Key }}, {{ printf "%q" .Value }})
	{{- end }}
})`

const containerTemplate = `Container({{ printf "%q" .Name }}{{ if or .Description .Technology }}, {{ printf "%q" .Description }}{{ end }}{{ if .Technology }}, {{ printf "%q" .Technology }}{{ end }}, func() {
	{{- range .Tags }}
	Tag({{ printf "%q" . }})
	{{- end }}
	{{- if .URL }}
	URL({{ printf "%q" .URL }})
	{{- end }}
	{{- range .Properties }}
	Prop({{ printf "%q" .Key }}, {{ printf "%q" .Value }})
	{{- end }}
})`

const componentTemplate = `Component({{ printf "%q" .Name }}{{ if or .Description .Technology }}, {{ printf "%q" .Description }}{{ end }}{{ if .Technology }}, {{ printf "%q" .Technology }}{{ end }}, func() {
	{{- range .Tags }}
	Tag({{ printf "%q" . }})
	{{- end }}
	{{- if .URL }}
	URL({{ printf "%q" .URL }})
	{{- end }}
	{{- range .Properties }}
	Prop({{ printf "%q" .Key }}, {{ printf "%q" .Value }})
	{{- end }}
})`

const relationshipTemplate = `{{ .RelationName }}({{ printf "%q" .DestinationPath }}{{ if .Description }}, {{ printf "%q" .Description }}{{ end }}{{ if .Technology }}, {{ printf "%q" .Technology }}{{ end }}{{ if .InteractionStyle }}, {{ .InteractionStyle }}{{ end }}{{ if or .Tags .URL }}, func() {
	{{- range .Tags }}
	Tag({{ printf "%q" . }})
	{{- end }}
	{{- if .URL }}
	URL({{ printf "%q" .URL }})
	{{- end }}
  }{{ end }})`

const landscapeViewTemplate = `SystemLandscapeView({{ printf "%q" .Key }}{{ if .Description }}, {{ printf "%q" .Description }}{{ end }}, func() {
	Title({{ printf "%q" .Title }})
	{{- if .EnterpriseBoundaryVisible }}
	EnterpriseBoundaryVisible()
	{{- end }}
	{{- if .PaperSize }}
	PaperSize({{ printf "%q" .PaperSize }})
	{{- end }}
	{{- range .ElementViews }}
	Add({{ printf "%q" .Element }})
	{{- end }}
	{{- range .RelationshipViews }}
	Link({{ printf "%q" .Source }}, {{ printf "%q" .Destination }})
	{{- end }}
})`

const systemContextViewTemplate = `SystemContextView({{ printf "%q" .SoftwareSystemName }}, {{ printf "%q" .Key }}{{ if .Description }}, {{ printf "%q" .Description }}{{ end }}, func() {
	Title({{ printf "%q" .Title }})
	{{- if .EnterpriseBoundaryVisible }}
	EnterpriseBoundaryVisible()
	{{- end }}
	{{- if .PaperSize }}
	PaperSize({{ printf "%q" .PaperSize }})
	{{- end }}
	{{- range .ElementViews }}
	Add({{ printf "%q" .Element }})
	{{- end }}
	{{- range .RelationshipViews }}
	Link({{ printf "%q" .Source }}, {{ printf "%q" .Destination }})
	{{- end }}
})`

const containerViewTemplate = `ContainerView({{ printf "%q" .SoftwareSystemName }}, {{ printf "%q" .Key }}{{ if .Description }}, {{ printf "%q" .Description }}{{ end }}, func() {
	Title({{ printf "%q" .Title }})
	{{- if .SystemBoundaryVisible }}
	SystemBoundaryVisible()
	{{- end }}
	{{- if .PaperSize }}
	PaperSize({{ printf "%q" .PaperSize }})
	{{- end }}
	{{- range .ElementViews }}
	Add({{ printf "%q" .Element }})
	{{- end }}
	{{- range .RelationshipViews }}
	Link({{ printf "%q" .Source }}, {{ printf "%q" .Destination }})
	{{- end }}
})`

const componentViewTemplate = `ComponentView({{ printf "%q" (printf "%s/%s" .SoftwareSystemName .ContainerName) }}, {{ printf "%q" .Key }}{{ if .Description }}, {{ printf "%q" .Description }}{{ end }}, func() {
	Title({{ printf "%q" .Title }})
	{{- if .ContainerBoundaryVisible }}
	ContainerBoundaryVisible()
	{{- end }}
	{{- if .PaperSize }}
	PaperSize({{ printf "%q" .PaperSize }})
	{{- end }}
	{{- range .ElementViews }}
	Add({{ printf "%q" .Element }})
	{{- end }}
	{{- range .RelationshipViews }}
	Link({{ printf "%q" .Source }}, {{ printf "%q" .Destination }})
	{{- end }}
})`

const elementStyleTemplate = `ElementStyle({{ printf "%q" .Tag }}, func() {
	Shape({{ printf "%q" .Shape }})
	{{- if .Icon }}
	Icon({{ printf "%q" .Icon }})
	{{- end }}
	{{- if .Background }}
	Background({{ printf "%q" .Background }})
	{{- end }}
	{{- if .Color }}
	Color({{ printf "%q" .Color }})
	{{- end }}
	{{- if .Stroke }}
	Stroke({{ printf "%q" .Stroke }})
	{{- end }}
	{{- if .Width }}
	Width({{ .Width }})
	{{- end }}
	{{- if .Height }}
	Height({{ .Height }})
	{{- end }}
	{{- if .FontSize }}
	FontSize({{ .FontSize }})
	{{- end }}
	{{- if .Metadata }}
	ShowMetadata()
	{{- end }}
	{{- if .Description }}
	ShowDescription()
	{{- end }}
	{{- if .Opacity }}
	Opacity({{ .Opacity }})
	{{- end }}
	{{- if .Border }}
	Border({{ printf "%q" .Border }})
	{{- end }}
})`

const relationshipStyleTemplate = `RelationshipStyle({{ printf "%q" .Tag }}, func() {
	{{- if .Thickness }}
	Thickness({{ .Thickness }})
	{{- end }}
	{{- if .FontSize }}
	FontSize({{ .FontSize }})
	{{- end }}
	{{- if .Width }}
	Width({{ .Width }})
	{{- end }}
	{{- if .Position }}
	Position({{ .Position }})
	{{- end }}
	{{- if .Color }}
	Color({{ printf "%q" .Color }})
	{{- end }}
	{{- if .Stroke }}
	Stroke({{ printf "%q" .Stroke }})
	{{- end }}
	{{- if .Dashed }}
	Dashed()
	{{- end }}
	{{- if .Routing }}
	Routing({{ printf "%q" .Routing }})
	{{- end }}
	{{- if .Opacity }}
	Opacity({{ .Opacity }})
	{{- end }}
})`
