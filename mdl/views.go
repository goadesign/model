package mdl

import (
	"path/filepath"
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
)

// funcs contains the functions used by templates.
var funcs = map[string]interface{}{"join": strings.Join, "wrap": wrap, "indent": indent}

// landscapeDiagram produces a file that contains Mermaid code representing the
// given view diagram.
func landscapeDiagram(lv *expr.LandscapeView) *codegen.File {
	ebv := false
	if lv.EnterpriseBoundaryVisible != nil {
		ebv = *lv.EnterpriseBoundaryVisible
	}
	return landscapeOrContextDiagram(lv.ViewProps, ebv)
}

// contextDiagram produces a file that contains Mermaid code representing the
// given view diagram.
func contextDiagram(cv *expr.ContextView) *codegen.File {
	ebv := false
	if cv.EnterpriseBoundaryVisible != nil {
		ebv = *cv.EnterpriseBoundaryVisible
	}
	return landscapeOrContextDiagram(cv.ViewProps, ebv)
}

// landscapeOrContextDiagram contains the shared logic between landscapeDiagram
// and contextDiagram.
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
		sections = append(sections, elements(external, "", 1))
	}
	if len(internal) > 0 {
		sections = append(sections, elements(internal, boundaryName, 1))
	}
	if len(vp.RelationshipViews) > 0 {
		sections = append(sections, relationships(vp.RelationshipViews))
	}

	return viewDiagram(vp, sections)
}

// containerDiagram produces a file that contains Mermaid code representing the
// given view diagram.
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
		sections = append(sections, elements(others, "", 1))
	}
	for name, elems := range bySystem {
		sections = append(sections, elements(elems, name, 1))
	}
	if len(cv.RelationshipViews) > 0 {
		sections = append(sections, relationships(cv.RelationshipViews))
	}

	return viewDiagram(cv.ViewProps, sections)
}

// componentDiagram produces a file that contains Mermaid code representing the
// given view diagram.
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
		sections = append(sections, elements(others, "", 1))
	}
	for name, elems := range byContainer {
		sections = append(sections, elements(elems, name, 1))
	}
	if len(cv.RelationshipViews) > 0 {
		sections = append(sections, relationships(cv.RelationshipViews))
	}
	return viewDiagram(cv.ViewProps, sections)
}

// deploymentDiagram produces a file that contains Mermaid code representing the
// given view diagram.
func deploymentDiagram(dv *expr.DeploymentView) *codegen.File {
	var sections []*codegen.SectionTemplate
	for _, ev := range dv.ElementViews {
		if dn, ok := expr.Registry[ev.Element.ID].(*expr.DeploymentNode); ok {
			sections = append(sections, deploymentNodeSections(dv, dn, 1)...)
		}
	}
	return viewDiagram(dv.ViewProps, sections)
}

// viewDiagram renders the common Mermaid code for all types of views. This includes the
// header and footer with class and style definitions.
func viewDiagram(vp *expr.ViewProps, sections []*codegen.SectionTemplate) *codegen.File {
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
			Style:       styleDef("", rs.Stroke, rs.Color, border, rs.Opacity, rs.Thickness),
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

// direction returns the Mermaid value for the AutoLayout direction defined in
// vp.
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
func styleDef(fill, stroke, color string, border expr.BorderKind, opacity, thickness *int) string {
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
	if thickness != nil {
		elems = append(elems, "stroke-width:"+strconv.Itoa(*thickness))
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
