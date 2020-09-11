package mdl

import (
	"path/filepath"
	"strings"

	"goa.design/goa/v3/codegen"
	"goa.design/model/expr"
	model "goa.design/model/pkg"
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

var funcs = map[string]interface{}{
	"direction":               direction,
	"className":               className,
	"background":              func(ev *expr.ElementView) string { return elemStyle(ev).Background },
	"stroke":                  stroke,
	"color":                   func(ev *expr.ElementView) string { return elemStyle(ev).Color },
	"idsByTags":               idsByTags,
	"join":                    strings.Join,
	"indexOfRelationshipView": indexOfRelationshipView,
	"relColor":                func(rv *expr.RelationshipView) string { return relStyle(rv).Color },
	"relStroke":               func(rv *expr.RelationshipView) string { return relStyle(rv).Stroke },
	"relOpacity":              func(rv *expr.RelationshipView) *int { return relStyle(rv).Opacity },
	"interpolate":             interpolate,
	"indent":                  indent,
	"fromPercent":             func(p *int) float64 { return float64(*p) / 100.0 },
}

func viewFile(vp *expr.ViewProps, sections []*codegen.SectionTemplate) *codegen.File {
	var styles []*expr.ElementStyle
	if expr.Root.Views.Styles != nil {
		styles = expr.Root.Views.Styles.Elements
	}
	data := map[string]interface{}{"ViewProps": vp, "Styles": styles, "Version": model.Version()}
	header := &codegen.SectionTemplate{Name: "header", Source: headerT, Data: data, FuncMap: funcs}
	footer := &codegen.SectionTemplate{Name: "footer", Source: footerT, Data: data, FuncMap: funcs}
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

func className(tag string) string {
	res := strings.ReplaceAll(tag, "-", "--")
	return strings.ReplaceAll(res, " ", "-")
}

// idsByTags maps the ids of the elements in evs to the corresponding element
// tags. It returns a map of ids indexed by tag names. Entries are only added if
// they have a style defined in styles.
func idsByTags(evs []*expr.ElementView, styles []*expr.ElementStyle) map[string][]string {
	res := make(map[string][]string)
	for _, ev := range evs {
		tags := strings.Split(ev.Element.Tags, ",")
	loop:
		for _, tag := range tags {
			for _, es := range styles {
				if es.Tag == tag {
					res[tag] = append(res[tag], ev.Element.ID)
					continue loop
				}
			}
		}
	}
	return res
}
func indexOfRelationshipView(rv *expr.RelationshipView, set []*expr.RelationshipView) (index int) {
	for _, rrv := range set {
		if rrv.RelationshipID == rv.RelationshipID {
			return
		}
		index++
	}
	panic("relationship missing") // bug
}

func indent(i int) string {
	return strings.Repeat(" ", 4*i)
}

const headerT = `%% Graph generated by mdl {{ .Version }} - DO NOT EDIT
graph {{ direction .ViewProps }}
`

const footerT = `{{- range .Styles }}
	{{- if or (.Background) (.Stroke) (.Color) }}{{ indent 1 }}classDef {{ className .Tag }}
        {{- if .Background }} fill:{{ .Background }}{{ if or .Stroke .Color }},{{ end }}{{ end }}
		{{- if .Stroke }} stroke:{{ .Stroke }}{{ if .Color }},{{ end }}{{ end }}
		{{- if .Color }} color:{{ .Color }}{{ end }}
		{{- if .Opacity }} opacity:{{ fromPercent .Opacity }}{{ end }};
{{ end }}
{{- end }}
{{- range $tag, $ids := (idsByTags .ViewProps.ElementViews .Styles) }}{{ indent 1 }}class {{ join $ids "," }} {{ className $tag }};
{{ end }}
{{- range .ViewProps.RelationshipViews }}
	{{- if or (relStroke .) (relColor .) (relOpacity .) }}{{ indent 1 }}linkStyle {{ indexOfRelationshipView . $.ViewProps.RelationshipViews }}
		{{- if relStroke . }} stroke:{{ relStroke . }}{{- if relColor . }},{{ end }}{{ end }}
		{{- if relColor . }} color:{{ relColor . }}{{ end }}
		{{- if relOpacity . }} opacity:{{ fromPercent (relOpacity .) }}{{ end }};
{{ end }}
	{{- if not (eq (interpolate .) "linear") }}{{ indent 1 }}linkStyle {{ indexOfRelationshipView . $.ViewProps.RelationshipViews }} interpolate {{ interpolate . }};
{{ end }}
{{- end }}`
