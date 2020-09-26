package mdl

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"goa.design/goa/v3/codegen"
	"goa.design/model/expr"
)

type (
	// legendData contains the data needed to render a diagram legend.
	legendData struct {
		LegendElements      []*legendElementData
		LegendRelationships []*legendRelationshipData
	}

	// legendElementData contains the data needed to render an element style legend.
	legendElementData struct {
		ID          string
		Name        string
		Description string
		Start, End  string
		Style       string
		IconURL     string
	}

	// legendRelationshipData contains the data needed to render a relationship
	// style legend.
	legendRelationshipData struct {
		SourceID, DestinationID string
		Start, End              string
		LinkIndex               int
		Style                   string
		Interpolate             string
		Description             string
	}
)

func legendDiagram(vp *expr.ViewProps) *codegen.File {
	// There is a many to many relationship between element types and styles. We
	// need to compute the set of unique combinations to produce the legend.
	styleMap := make(map[string]*legendElementData)
	seen := make(map[string]struct{})
	for i, ev := range vp.ElementViews {
		var id, name, tags string
		{
			id = strconv.Itoa(i)

			switch expr.Registry[ev.Element.ID].(type) {
			case *expr.Person:
				name = "Person"
			case *expr.SoftwareSystem:
				name = "Software System"
			case *expr.Container:
				name = "Container"
			case *expr.Component:
				name = "Component"
			case *expr.DeploymentNode:
				name = "Deployment Node"
			case *expr.InfrastructureNode:
				name = "Infrastructure Node"
			case *expr.ContainerInstance:
				name = "Container Instance"
			default:
				panic("unknown element type:" + fmt.Sprintf("%T", expr.Registry[ev.Element.ID]))
			}

			elems := strings.Split(ev.Element.Tags, ",")
			elems = remove(elems, "Element", "Person", "Software System", "Container", "Component", "Deployment Node", "Infrastructure Node", "Container Instance")
			sort.Strings(elems)
			tags = strings.Join(elems, "<br/>")
			key := tags + "--" + name
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
		}
		if _, ok := styleMap[id]; !ok {
			es := elemStyle(ev)
			start, end := nodeStartEnd(ev)
			styleMap[id] = &legendElementData{
				ID:          id,
				Name:        name,
				Description: tags,
				Start:       start,
				End:         end,
				Style:       styleDef(es.Background, es.Stroke, es.Color, es.Border, es.Opacity, nil),
				IconURL:     es.Icon,
			}
		}
	}
	legelems := make([]*legendElementData, len(styleMap))
	i := 0 // no need to sort, mermaid rendering reorders arbitrarily
	for _, le := range styleMap {
		legelems[i] = le
		i++
	}

	// Relationships are easier: there is only one type of relationship.
	legrels := make(map[string]*legendRelationshipData)
	for _, rv := range vp.RelationshipViews {
		var desc string
		{
			rel := expr.Registry[rv.RelationshipID].(*expr.Relationship)
			tags := strings.Split(rel.Tags, ",")
			if len(tags) > 1 {
				for i, tag := range tags {
					if tag == "Relationship" {
						tags = append(tags[:i], tags[i+1:]...)
						break
					}
				}
			}
			sort.Strings(tags)
			desc = strings.Join(tags, "<br/>")
		}
		if _, ok := legrels[desc]; ok {
			continue
		}
		rs := relStyle(rv)
		start, end := lineStartEnd(rs)
		border := expr.BorderUndefined
		if rs.Dashed != nil && *rs.Dashed {
			border = expr.BorderDashed
		}
		legrels[desc] = &legendRelationshipData{
			SourceID:      "src" + rs.Tag,
			DestinationID: "dest" + rs.Tag,
			Start:         start,
			End:           end,
			Style:         styleDef("", rs.Stroke, rs.Color, border, rs.Opacity, rs.Thickness),
			Interpolate:   interpolate(rs),
			Description:   desc,
		}
	}
	keys := make([]string, len(legrels))
	i = 0
	for d := range legrels {
		keys[i] = d
		i++
	}
	srels := make([]*legendRelationshipData, len(keys))
	for i, d := range keys {
		data := legrels[d]
		data.LinkIndex = i
		srels[i] = data
	}
	legend := &codegen.SectionTemplate{
		Name:   "legend",
		Source: legendT,
		Data: &legendData{
			LegendElements:      legelems,
			LegendRelationships: srels,
		},
		FuncMap: funcs,
	}
	path := filepath.Join(codegen.Gendir, filepath.Join(codegen.Gendir, "diagrams", vp.Key+".legend.mmd"))

	return &codegen.File{Path: path, SectionTemplates: []*codegen.SectionTemplate{legend}}
}

func remove(slice []string, vals ...string) []string {
	for _, val := range vals {
		idx := -1
		for i, elem := range slice {
			if elem == val {
				idx = i
				break
			}
		}
		if idx > -1 {
			slice = append(slice[:idx], slice[idx+1:]...)
		}
	}
	return slice
}

const legendT = `graph TD
{{ range .LegendElements }}
    {{ .ID }}{{ .Start }}"{{ if .IconURL }}<img src='{{ .IconURL }}'/>
	{{- end }}<div class='element'><div class='element-title'>{{ .Name }}</div class='element-title'><div class='element-technology'></div><div class='element-description'>{{ .Description }}</div class='element-description'></div>"{{ .End }}
{{ if .Style }}    style {{ .ID }} {{ .Style }}
{{ end }}{{ end }}

{{- range .LegendRelationships }}
    {{ .SourceID }}[A]{{ .Start }}"<div class='relationship'><div class='relationship-label'>{{ wrap .Description 30 }}</div></div>"{{ .End }}{{ .DestinationID }}[B]
{{ if .Style }}    linkStyle {{ .LinkIndex }} {{ .Style }}
{{ end }}{{ if .Interpolate }}        linkStyle {{ .LinkIndex }} interpolate {{ .Interpolate }}
{{ end }}
{{- end }}`
