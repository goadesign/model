package mdl

import (
	"bytes"
	"strings"
	"unicode"

	"goa.design/goa/v3/codegen"
	"goa.design/model/expr"
)

type (
	// elementsData is the data structure used to render the elements template.
	elementsData struct {
		// BoundaryName is the name of the subgraph rendered around the elements
		// if any.
		BoundaryName string
		// Elements to render
		Elements []*elementData
	}

	elementData struct {
		// Indent of line in rendered mermaid source
		Indent int
		// ID of element
		ID string
		// Start and End node mermaid symbols (e.g. "[", "]")
		Start, End string
		// Name of element
		Name string
		// Description of element
		Description string
		// Technology used by element if any
		Technology string
		// URL to redirect to when element is clicked if any
		URL string
		// IconURL is the URL to an icon if any
		IconURL string
		// Background is the background color defined in the design if any
		Background string
		// Stroke is the stroke color defined in the design if any
		Stroke string
	}
)

func elements(evs []*expr.ElementView, boundary string, ind int) *codegen.SectionTemplate {
	elems := make([]*elementData, len(evs))
	if boundary != "" {
		ind++
	}
	join := func(name, tech string) string {
		if tech != "" {
			return name + ": " + tech
		}
		return name
	}
	for i, ev := range evs {
		var tech string
		switch e := expr.Registry[ev.Element.ID].(type) {
		case *expr.Person:
			tech = "Person"
		case *expr.SoftwareSystem:
			tech = "Software System"
		case *expr.Container:
			tech = join("Container", e.Technology)
		case *expr.Component:
			tech = join("Component", e.Technology)
		case *expr.ContainerInstance:
			tech = join("Container", e.Technology)
		case *expr.InfrastructureNode:
			tech = join("Infrastructure Node", e.Technology)
		}
		es := elemStyle(ev)
		start, end := nodeStartEnd(ev)
		elems[i] = &elementData{
			ID:          ev.Element.ID,
			Indent:      ind,
			Start:       start,
			End:         end,
			Name:        ev.Element.Name,
			Description: ev.Element.Description,
			Technology:  tech,
			URL:         ev.Element.URL,
			IconURL:     es.Icon,
			Background:  es.Background,
			Stroke:      es.Stroke,
		}
	}
	data := &elementsData{
		BoundaryName: boundary,
		Elements:     elems,
	}
	funcs := map[string]interface{}{"wrap": wrap, "stroke": stroke, "indent": indent}
	return &codegen.SectionTemplate{Name: "elements", Source: elementT, Data: data, FuncMap: funcs}
}

// wrap wraps the given string to n charaters per line, encodes the
// results so mermaid is happy to use them as element description and separates
// each line with <br/> .
func wrap(s string, n uint) string {
	init := make([]byte, 0, len(s))
	buf := bytes.NewBuffer(init)
	var current uint
	var wordBuf, spaceBuf bytes.Buffer
	for _, char := range s {
		if char == '\n' {
			if wordBuf.Len() == 0 {
				if current+uint(spaceBuf.Len()) > n {
					current = 0
				} else {
					current += uint(spaceBuf.Len())
					spaceBuf.WriteTo(buf)
				}
				spaceBuf.Reset()
			} else {
				current += uint(spaceBuf.Len() + wordBuf.Len())
				spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				wordBuf.WriteTo(buf)
				wordBuf.Reset()
			}
			buf.WriteRune(char)
			current = 0
		} else if unicode.IsSpace(char) {
			if spaceBuf.Len() == 0 || wordBuf.Len() > 0 {
				current += uint(spaceBuf.Len() + wordBuf.Len())
				spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				wordBuf.WriteTo(buf)
				wordBuf.Reset()
			}
			spaceBuf.WriteRune(char)
		} else {
			wordBuf.WriteRune(char)
			if current+uint(spaceBuf.Len()+wordBuf.Len()) > n && uint(wordBuf.Len()) < n {
				buf.WriteString("<br/>")
				current = 0
				spaceBuf.Reset()
			}
		}
	}
	if wordBuf.Len() == 0 {
		if current+uint(spaceBuf.Len()) <= n {
			spaceBuf.WriteTo(buf)
		}
	} else {
		spaceBuf.WriteTo(buf)
		wordBuf.WriteTo(buf)
	}
	return strings.ReplaceAll(buf.String(), "\n", "<br/>")
}

func indent(n int) string {
	return strings.Repeat(" ", n*4)
}

func nodeStartEnd(ev *expr.ElementView) (string, string) {
	// Look for explicit shape first
	es := elemStyle(ev)
	switch es.Shape {
	case expr.ShapeBox:
		return `[`, `]`
	case expr.ShapeRoundedBox:
		return `(`, `)`
	case expr.ShapeCircle:
		return `((`, `))`
	case expr.ShapeEllipse:
		return `([`, `])` // Approximation - this is actually a stadium shape
	case expr.ShapeHexagon:
		return `{{`, `}}`
	case expr.ShapeCylinder:
		return `[(`, `)]`
	}

	// Compute default shape for given element type.
	switch expr.Registry[ev.Element.ID].(type) {
	case *expr.Person:
		return `([`, `])`
	case *expr.SoftwareSystem:
		return "[", "]"
	case *expr.Container, *expr.ContainerInstance:
		return "(", ")"
	default:
		return "[", "]"
	}
}

// input: ElementsData
const elementT = `{{ if .BoundaryName }}{{ indent 1 }}subgraph boundary [{{ .BoundaryName }}]
{{ end }}
{{- range .Elements }}{{ indent .Indent }}{{ .ID }}{{ .Start }}"
{{- if .IconURL }}<img src='{{ .IconURL }}'/>
{{ end -}}
<div class='element'><div class='element-title'>{{ wrap .Name 25 }}</div><div class='element-technology'>{{ if .Technology }}[{{ wrap .Technology 30 }}]{{ end }}</div><div class='element-description'>{{ wrap .Description 30 }}</div></div>"{{ .End }}
{{- if .URL }}
{{ indent .Indent }}click {{ .ID }} "{{ .URL }}"{{ if .URLTooltip }} "{{ .URLTooltip }}"{{ end }}
{{ end }}
{{- if not .Stroke }}
{{ indent .Indent }}style {{ .ID }} stroke:{{ stroke . }};
{{- end }}
{{ end }}
{{- if .BoundaryName }}{{ indent 1 }}end
{{ indent 1 }}style boundary fill:#ffffff,stroke:#909090,color:#000000,stroke-dasharray: 15 5;
{{ end }}`
