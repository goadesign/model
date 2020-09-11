package mdl

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"goa.design/goa/v3/codegen"
	"goa.design/model/expr"
	model "goa.design/model/pkg"
)

type (
	// headerData is the data structure used to render the header template.
	headerData struct {
		Version   string
		Direction string
	}

	// footerData is the data structure used to render the footer template.
	footerData struct {
		Classes         []*elementClassData
		IDsByClassNames map[string][]string
		Links           []*linkStyleData
	}

	// elementClassData contains the data needed to render element classes.
	elementClassData struct {
		Style     string
		ClassName string
	}

	// linkStyleData contains the data needed to render link styles.
	linkStyleData struct {
		LinkIndex   int
		Style       string
		Interpolate string
	}

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

func landscapeDiagram(lv *expr.LandscapeView) *codegen.File {
	ebv := false
	if lv.EnterpriseBoundaryVisible != nil {
		ebv = *lv.EnterpriseBoundaryVisible
	}
	return landscapeOrContextDiagram(lv.ViewProps, ebv)
}

func contextDiagram(cv *expr.ContextView) *codegen.File {
	ebv := false
	if cv.EnterpriseBoundaryVisible != nil {
		ebv = *cv.EnterpriseBoundaryVisible
	}
	return landscapeOrContextDiagram(cv.ViewProps, ebv)
}

func landscapeOrContextDiagram(vp *expr.ViewProps, ebv bool) *codegen.File {
	var internal, external []*expr.ElementView
	for _, ev := range vp.ElementViews {
		switch a := expr.Registry[ev.Element.ID].(type) {
		case *expr.Person:
			if a.Location == expr.LocationUndefined || a.Location == expr.LocationInternal {
				internal = append(internal, ev)
			} else {
				external = append(external, ev)
			}
		case *expr.SoftwareSystem:
			if a.Location == expr.LocationUndefined || a.Location == expr.LocationInternal {
				internal = append(internal, ev)
			} else {
				external = append(external, ev)
			}
		}
	}
	ebv = ebv && len(internal) > 0
	var boundaryName string
	if ebv {
		boundaryName = expr.Root.Model.Enterprise
		if boundaryName == "" {
			boundaryName = "Internal"
		}
	}
	var sections []*codegen.SectionTemplate
	if len(external) > 0 {
		sections = append(sections, elements(external, ""))
	}
	if len(internal) > 0 {
		sections = append(sections, elements(internal, boundaryName))
	}
	if len(vp.RelationshipViews) > 0 {
		sections = append(sections, relationships(vp.RelationshipViews))
	}

	return viewFile(vp, sections)
}

func containerDiagram(cv *expr.ContainerView) *codegen.File {
	var others []*expr.ElementView
	bySystem := make(map[string][]*expr.ElementView)
	for _, ev := range cv.ElementViews {
		switch c := expr.Registry[ev.Element.ID].(type) {
		case *expr.Container:
			bySystem[c.System.Name] = append(bySystem[c.System.Name], ev)
		default:
			others = append(others, ev)
		}
	}
	var sections []*codegen.SectionTemplate
	if len(others) > 0 {
		sections = append(sections, elements(others, ""))
	}
	for name, elems := range bySystem {
		sections = append(sections, elements(elems, name))
	}
	if len(cv.RelationshipViews) > 0 {
		sections = append(sections, relationships(cv.RelationshipViews))
	}

	return viewFile(cv.ViewProps, sections)
}

func componentDiagram(cv *expr.ComponentView) *codegen.File {
	var others []*expr.ElementView
	byContainer := make(map[string][]*expr.ElementView)
	for _, ev := range cv.ElementViews {
		switch c := expr.Registry[ev.Element.ID].(type) {
		case *expr.Component:
			byContainer[c.Container.Name] = append(byContainer[c.Container.Name], ev)
		default:
			others = append(others, ev)
		}
	}
	var sections []*codegen.SectionTemplate
	if len(others) > 0 {
		sections = append(sections, elements(others, ""))
	}
	for name, elems := range byContainer {
		sections = append(sections, elements(elems, name))
	}
	if len(cv.RelationshipViews) > 0 {
		sections = append(sections, relationships(cv.RelationshipViews))
	}
	return viewFile(cv.ViewProps, sections)
}

func dynamicDiagram(dv *expr.DynamicView) *codegen.File {
	return viewFile(dv.ViewProps, nil)
}

func deploymentDiagram(dv *expr.DeploymentView) *codegen.File {
	return viewFile(dv.ViewProps, nil)
}

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
			for i, tag := range tags {
				if tag == "Relationship" {
					tags = append(tags[:i], tags[i+1:]...)
					break
				}
			}
			if len(tags) == 0 {
				continue
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
			Style:         styleDef("", rs.Stroke, rs.Color, border, rs.Opacity, rs.Thick),
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

var funcs = map[string]interface{}{"indent": indent, "join": strings.Join}

func viewFile(vp *expr.ViewProps, sections []*codegen.SectionTemplate) *codegen.File {
	var classes []*elementClassData
	var styles []*expr.ElementStyle
	if expr.Root.Views.Styles != nil {
		styles = expr.Root.Views.Styles.Elements
		classes = make([]*elementClassData, len(styles))
		for i, es := range styles {
			style := styleDef(es.Background, es.Stroke, es.Color, es.Border, es.Opacity, nil)
			classes[i] = &elementClassData{
				Style:     style,
				ClassName: className(es.Tag),
			}
		}
	}

	links := make([]*linkStyleData, len(vp.RelationshipViews))
	for i, rv := range vp.RelationshipViews {
		rs := relStyle(rv)
		border := expr.BorderUndefined
		if rs.Dashed != nil && *rs.Dashed {
			border = expr.BorderDashed
		}
		links[i] = &linkStyleData{
			LinkIndex:   i,
			Style:       styleDef("", rs.Stroke, rs.Color, border, rs.Opacity, rs.Thick),
			Interpolate: interpolate(relStyle(rv)),
		}
	}
	header := &codegen.SectionTemplate{
		Name:    "header",
		Source:  headerT,
		Data:    &headerData{Direction: direction(vp), Version: model.Version()},
		FuncMap: funcs,
	}
	footer := &codegen.SectionTemplate{
		Name:   "footer",
		Source: footerT,
		Data: &footerData{
			Classes:         classes,
			Links:           links,
			IDsByClassNames: idsByClassNames(vp.ElementViews, styles),
		},
		FuncMap: funcs,
	}
	sections = append([]*codegen.SectionTemplate{header}, append(sections, footer)...)
	path := filepath.Join(codegen.Gendir, filepath.Join(codegen.Gendir, "diagrams", vp.Key+".mmd"))

	return &codegen.File{Path: path, SectionTemplates: sections}
}

func direction(vp *expr.ViewProps) string {
	if vp.AutoLayout == nil {
		return "TB"
	}
	switch vp.AutoLayout.RankDirection {
	case expr.RankBottomTop:
		return "BT"
	case expr.RankLeftRight:
		return "LR"
	case expr.RankRightLeft:
		return "RL"
	default:
		return "TB"
	}
}

// className attempts to produce a mermaid safe styles class name from the given
// tag value.
func className(tag string) string {
	res := strings.ReplaceAll(tag, "-", "--")
	return strings.ReplaceAll(res, " ", "-")
}

// styleDef renders a valid mermaid/SVG style line from the given values.
func styleDef(fill, stroke, color string, border expr.BorderKind, opacity *int, thick *bool) string {
	var elems []string
	if fill != "" {
		elems = append(elems, "fill:"+fill)
	}
	if stroke != "" {
		elems = append(elems, "stroke:"+stroke)
	}
	if color != "" {
		elems = append(elems, "color:"+color)
	}
	switch border {
	case expr.BorderDashed:
		elems = append(elems, "stroke-dasharray: 15 5")
	case expr.BorderDotted:
		elems = append(elems, "stroke-dasharray: 3 3")
	}
	if opacity != nil {
		elems = append(elems, "opacity:"+strconv.FormatFloat(float64(*opacity)/100.0, 'f', 2, 64))
	}
	if thick != nil && *thick {
		elems = append(elems, "stroke-width:4")
	}
	if len(elems) == 0 {
		return ""
	}
	return strings.Join(elems, ",") + ",foo; %% yeah FOO! https://github.com/mermaid-js/mermaid/issues/1666"
}

// idsByTags maps the ids of the elements in evs to the corresponding element
// tags. It returns a map of ids indexed by tag names. Entries are only added if
// they have a style defined in styles.
func idsByClassNames(evs []*expr.ElementView, styles []*expr.ElementStyle) map[string][]string {
	res := make(map[string][]string)
	for _, ev := range evs {
		tags := strings.Split(ev.Element.Tags, ",")
	loop:
		for _, tag := range tags {
			for _, es := range styles {
				if es.Tag == tag {
					res[className(tag)] = append(res[className(tag)], ev.Element.ID)
					continue loop
				}
			}
		}
	}
	return res
}

func indent(i int) string {
	return strings.Repeat(" ", 4*i)
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

const headerT = `%% Graph generated by mdl {{ .Version }} - DO NOT EDIT
graph {{ .Direction }}
`

const footerT = `{{- range .Classes }}
{{- if .Style }}{{ indent 1 }}classDef {{ .ClassName }} {{ .Style }}
{{ end }}
{{- end }}
{{- range $className, $ids := .IDsByClassNames }}{{ indent 1 }}class {{ join $ids "," }} {{ $className }};
{{ end }}

{{- range .Links }}
	{{- if .Style }}{{ indent 1 }}linkStyle {{ .LinkIndex }} {{ .Style }}
{{ end }}{{ if .Interpolate }}{{ indent 1 }}linkStyle {{ .LinkIndex }} interpolate {{ .Interpolate }};
{{ end }}
{{- end }}`

const legendT = `graph TD
{{ range .LegendElements }}
{{- indent 1 }}{{ .ID }}{{ .Start }}"{{ if .IconURL }}<img src='{{ .IconURL }}'/>
	{{- end }}<div class='element'><div class='element-title'>{{ .Name }}</div class='element-title'><div class='element-technology'></div><div class='element-description'>{{ .Description }}</div class='element-description'></div>"{{ .End }}
{{ indent 1 }}style {{ .ID }} {{ .Style }}
{{ end }}

{{- range .LegendRelationships }}
{{- indent 1 }}{{ .SourceID }}[A]{{ .Start }}{{ .Description }}{{ .End }}{{ .DestinationID }}[B]
{{ if .Style }}{{ indent 1 }}linkStyle {{ .LinkIndex }} {{ .Style }}
{{ end }}{{ if .Interpolate }}{{ indent 2 }}linkStyle {{ .LinkIndex }} interpolate {{ .Interpolate }}
{{ end }}
{{- end }}`
