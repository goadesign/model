package svc

import (
	"fmt"
	"strings"
	"text/template"

	geneditor "goa.design/model/svc/gen/dsl_editor"
)

var systemT = template.Must(template.New("system").Parse(systemTemplate))
var personT = template.Must(template.New("person").Parse(personTemplate))
var containerT = template.Must(template.New("container").Parse(containerTemplate))
var componentT = template.Must(template.New("component").Parse(componentTemplate))

func systemDSL(s *geneditor.System) string       { return execute(systemT, s) }
func personDSL(p *geneditor.Person) string       { return execute(personT, p) }
func containerDSL(c *geneditor.Container) string { return execute(containerT, c) }
func componentDSL(c *geneditor.Component) string { return execute(componentT, c) }

// execute executes the given template with the given data and returns the
// result.
func execute(t *template.Template, data interface{}) string {
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
