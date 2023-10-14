package codegen

import (
	"bytes"
	"fmt"
	"strings"

	"goa.design/goa/v3/codegen"
	"golang.org/x/tools/imports"

	"goa.design/model/dsl"
	"goa.design/model/expr"
	"goa.design/model/mdl"
	model "goa.design/model/pkg"
)

// Tags that should not be generated in DSL.
var builtInTags []string

func init() {
	for _, tags := range [][]string{
		expr.PersonTags,
		expr.SoftwareSystemTags,
		expr.ContainerTags,
		expr.ComponentTags,
		expr.DeploymentNodeTags,
		expr.InfrastructureNodeTags,
		expr.ContainerInstanceTags,
		expr.RelationshipTags,
	} {
		builtInTags = append(builtInTags, tags...)
	}
}

// Model generates the model DSL from the given design package.
// pkg is the name of the generated package (e.g. "model").
func Model(d *mdl.Design, pkg string) ([]byte, error) {
	if d.Model == nil {
		return nil, fmt.Errorf("model is nil")
	}
	var template string
	for name, tmpl := range templates {
		template += fmt.Sprintf("{{define %q}}%s{{end}}", name, tmpl)
	}
	template += designT
	section := &codegen.SectionTemplate{
		Name:   pkg,
		Source: template,
		Data: map[string]any{
			"Design":      d,
			"Pkg":         pkg,
			"ToolVersion": model.Version(),
		},
		FuncMap: map[string]any{
			"elementPath":             elementPath,
			"filterTags":              filterTags,
			"softwareSystemData":      softwareSystemData,
			"containerData":           containerData,
			"componentData":           componentData,
			"personData":              personData,
			"systemLandscapeViewData": systemLandscapeViewData,
			"systemContextViewData":   systemContextViewData,
			"containerViewData":       containerViewData,
			"componentViewData":       componentViewData,
			"viewPropsData":           viewPropsData,
			"relData":                 relData,
			"relDSLFunc":              relDSLFunc,
			"systemHasFunc":           systemHasFunc,
			"containerHasFunc":        containerHasFunc,
			"componentHasFunc":        componentHasFunc,
			"personHasFunc":           personHasFunc,
			"autoLayoutHasFunc":       autoLayoutHasFunc,
			"hasViews":                hasViews,
			"findRelationship":        findRelationship,
			"deref":                   func(i *int) int { return *i },
		},
	}
	var buf bytes.Buffer
	if err := section.Write(&buf); err != nil {
		return nil, err
	}
	opt := imports.Options{Comments: true, FormatOnly: true}
	res, err := imports.Process("", buf.Bytes(), &opt)
	if err != nil {
		// Print content for troubleshooting
		fmt.Println(buf.String())
	}
	return res, err
}

// relDSLFunc is a function used by the DSL codegen to compute the name of the
// DSL function used to represent the corresponding relationship, one of "Uses",
// "Delivers" or "InteractsWith".
func relDSLFunc(mod *mdl.Model, rel *mdl.Relationship) string {
	var sourceIsPerson, destIsPerson bool
	for _, p := range mod.People {
		if p.ID == rel.SourceID {
			sourceIsPerson = true
		} else if p.ID == rel.DestinationID {
			destIsPerson = true
		}
	}
	if sourceIsPerson && destIsPerson {
		return "InteractsWith"
	}
	if destIsPerson {
		return "Delivers"
	}
	return "Uses"
}

// relData produces a data structure appropriate for running the useT template.
func relData(mod *mdl.Model, rel *mdl.Relationship, current string) map[string]any {
	return map[string]any{
		"Model":        mod,
		"Relationship": rel,
		"CurrentPath":  current,
	}
}

// softwareSystemData produces a data structure appropriate for running the systemT template.
func softwareSystemData(mod *mdl.Model, s *mdl.SoftwareSystem) map[string]any {
	return map[string]any{
		"Model":          mod,
		"SoftwareSystem": s,
	}
}

// containerData produces a data structure appropriate for running the containerT template.
func containerData(mod *mdl.Model, c *mdl.Container, current string) map[string]any {
	return map[string]any{
		"Model":       mod,
		"Container":   c,
		"CurrentPath": current,
	}
}

// componentData produces a data structure appropriate for running the componentT template.
func componentData(mod *mdl.Model, cmp *mdl.Component, current string) map[string]any {
	return map[string]any{
		"Model":       mod,
		"Component":   cmp,
		"CurrentPath": current,
	}
}

// personData produces a data structure appropriate for running the personT template.
func personData(mod *mdl.Model, p *mdl.Person) map[string]any {
	return map[string]any{
		"Model":  mod,
		"Person": p,
	}
}

// systemLandscapeViewData produces a data structure appropriate for running the systemLandscapeViewT template.
func systemLandscapeViewData(mod *mdl.Model, v *mdl.LandscapeView) map[string]any {
	return map[string]any{
		"Model": mod,
		"View":  v,
	}
}

// systemContextViewData produces a data structure appropriate for running the systemContextViewT template.
func systemContextViewData(mod *mdl.Model, v *mdl.ContextView) map[string]any {
	return map[string]any{
		"Model": mod,
		"View":  v,
	}
}

// containerViewData produces a data structure appropriate for running the containerViewT template.
func containerViewData(mod *mdl.Model, v *mdl.ContainerView) map[string]any {
	return map[string]any{
		"Model": mod,
		"View":  v,
	}
}

// componentViewData produces a data structure appropriate for running the componentViewT template.
func componentViewData(mod *mdl.Model, v *mdl.ComponentView) map[string]any {
	return map[string]any{
		"Model": mod,
		"View":  v,
	}
}

// viewPropsData produces a data structure appropriate for running the viewPropsT template.
func viewPropsData(mod *mdl.Model, v *mdl.ViewProps) map[string]any {
	return map[string]any{
		"Model":                 mod,
		"Props":                 v,
		"DefaultRankSeparation": dsl.DefaultRankSeparation,
		"DefaultNodeSeparation": dsl.DefaultNodeSeparation,
		"DefaultEdgeSeparation": dsl.DefaultEdgeSeparation,
	}
}

// elementPath is used by templates codegen to compute the path to the element with the given ID.
func elementPath(mod *mdl.Model, id string, roots ...string) string {
	var root string
	if len(roots) > 0 {
		root = roots[0]
	}
	for _, p := range mod.People {
		if p.ID == id {
			return p.Name
		}
	}
	for _, s := range mod.Systems {
		if s.ID == id {
			return s.Name
		}
		for _, c := range s.Containers {
			if c.ID == id {
				if root == s.Name {
					return c.Name
				}
				return fmt.Sprintf("%s/%s", s.Name, c.Name)
			}
			for _, cmp := range c.Components {
				if cmp.ID == id {
					if root == s.Name {
						return fmt.Sprintf("%s/%s", c.Name, cmp.Name)
					}
					if root == fmt.Sprintf("%s/%s", s.Name, c.Name) {
						return cmp.Name
					}
					return fmt.Sprintf("%s/%s/%s", s.Name, c.Name, cmp.Name)
				}
			}
		}
	}
	return ""
}

// findRelatioship returns the relationship with the given id.
func findRelationship(mod *mdl.Model, id string) *mdl.Relationship {
	for _, p := range mod.People {
		for _, rel := range p.Relationships {
			if rel.ID == id {
				return rel
			}
		}
	}
	for _, s := range mod.Systems {
		for _, rel := range s.Relationships {
			if rel.ID == id {
				return rel
			}
			for _, c := range s.Containers {
				for _, rel := range c.Relationships {
					if rel.ID == id {
						return rel
					}
				}
				for _, cmp := range c.Components {
					for _, rel := range cmp.Relationships {
						if rel.ID == id {
							return rel
						}
					}
				}
			}
		}
	}
	return nil
}

func filterTags(s string) []string {
	parts := strings.Split(s, ",")
	var res []string
loop:
	for _, p := range parts {
		p = strings.TrimSpace(p)
		for _, builtIn := range builtInTags {
			if p == builtIn {
				continue loop
			}
		}
		res = append(res, p)
	}
	return res
}

// personHasFunc returns true if an anonymous DSL function must be generated for p.
func personHasFunc(p *mdl.Person) bool {
	return len(filterTags(p.Tags)) > 0 ||
		p.URL != "" ||
		len(p.Properties) > 0 ||
		len(p.Relationships) > 0
}

// systemHasFunc returns true if an anonymous DSL function must be generated for s.
func systemHasFunc(s *mdl.SoftwareSystem) bool {
	return len(filterTags(s.Tags)) > 0 ||
		s.URL != "" ||
		s.Location == mdl.LocationExternal ||
		len(s.Properties) > 0 ||
		len(s.Relationships) > 0 ||
		len(s.Containers) > 0
}

// containerHasFunc returns true if an anonymous DSL function must be generated for c.
func containerHasFunc(c *mdl.Container) bool {
	return len(filterTags(c.Tags)) > 0 ||
		c.URL != "" ||
		len(c.Properties) > 0 ||
		len(c.Relationships) > 0 ||
		len(c.Components) > 0
}

// componentHasFunc returns true if an anonymous DSL function must be generated for cmp.
func componentHasFunc(cmp *mdl.Component) bool {
	return len(filterTags(cmp.Tags)) > 0 ||
		cmp.URL != "" ||
		len(cmp.Properties) > 0 ||
		len(cmp.Relationships) > 0
}

// autoLayoutHasFunc returns true if an anonymous DSL function must be generated for v.
func autoLayoutHasFunc(l *mdl.AutoLayout) bool {
	return l.RankSep != nil && *l.RankSep != dsl.DefaultRankSeparation ||
		l.NodeSep != nil && *l.NodeSep != dsl.DefaultNodeSeparation ||
		l.EdgeSep != nil && *l.EdgeSep != dsl.DefaultEdgeSeparation ||
		l.Vertices != nil && *l.Vertices
}

// hasViews returns true if the given views is not empty.
func hasViews(v *mdl.Views) bool {
	return len(v.LandscapeViews) > 0 ||
		len(v.ContextViews) > 0 ||
		len(v.ContainerViews) > 0 ||
		len(v.ComponentViews) > 0 ||
		len(v.DeploymentViews) > 0 ||
		len(v.FilteredViews) > 0
}

var templates = map[string]string{
	"systemT":                systemT,
	"containerT":             containerT,
	"componentT":             componentT,
	"personT":                personT,
	"useT":                   useT,
	"deploymentEnvironmentT": deploymentEnvironmentT,
	"deploymentNodeT":        deploymentNodeT,
	"infrastructureNodeT":    infrastructureNodeT,
	"containerInstanceT":     containerInstanceT,
	"healthCheckT":           healthCheckT,
	"viewPropsT":             viewPropsT,
	"systemLandscapeViewT":   systemLandscapeViewT,
	"systemContextViewT":     systemContextViewT,
	"containerViewT":         containerViewT,
	"componentViewT":         componentViewT,
	"filteredViewT":          filteredViewT,
	"deploymentView":         deploymentViewT,
	"dynamicView":            dynamicViewT,
	"styleT":                 styleT,
}

const designT = `// Code generated by mdl {{.ToolVersion}}.

package {{ .Pkg }}

import . "goa.design/model/dsl"

{{- with .Design }}
var _ = Design("{{.Name}}", "{{.Description}}", func() {
	{{- if .Model.Enterprise }}
	Enterprise({{ printf "%q" .Model.Enterprise }})
	{{- end }}
	{{- range .Model.People }}
	{{ template "personT" (personData $.Model .) }}
	{{- end }}
	{{- range .Model.Systems }}
	{{ template "systemT" (softwareSystemData $.Design.Model .) }}
	{{- end }}
	{{- range .Model.DeploymentNodes }}
	{{ template "deploymentEnvironmentT" . }}
	{{- end }}
	{{- if hasViews .Views }}
	Views(func() {
	{{- range .Views.LandscapeViews }}
	{{ template "systemLandscapeViewT" (systemLandscapeViewData $.Design.Model .) }}
	{{- end }}
	{{- range .Views.ContextViews }}
	{{ template "systemContextViewT" (systemContextViewData $.Design.Model .) }}
	{{- end }}
	{{- range .Views.ContainerViews }}
	{{ template "containerViewT" (containerViewData $.Design.Model .) }}
	{{- end }}
	{{- range .Views.ComponentViews }}
	{{ template "componentViewT" (componentViewData $.Design.Model .) }}
	{{- end }}
	{{- range .Views.DeploymentViews }}
	{{ template "deploymentViewT" . }}
	{{- end }}
	{{- range .Views.FilteredViews }}
	{{ template "filteredViewT" . }}
	{{- end }}
	{{- if .Views.Styles }}
	{{ template "styleT" .Views.Styles }}
	{{- end }}
	})
	{{- end }}
})
{{- end }}`

var systemT = `SoftwareSystem("{{ .SoftwareSystem.Name}}", "{{ .SoftwareSystem.Description }}"{{ if systemHasFunc .SoftwareSystem }}, func() {
	{{- range filterTags .SoftwareSystem.Tags }}
	Tag({{ printf "%q" . }})
	{{- end }}
	{{- if .SoftwareSystem.URL }}
	URL({{ printf "%q" .SoftwareSystem.URL }})
	{{- end }}
	{{- if eq .SoftwareSystem.Location 2 }}
	External()
	{{- end }}
	{{- range $k, $v := .SoftwareSystem.Properties }}
	Prop({{ printf "%q" $k }}, {{ printf "%q" $v }})
	{{- end }}
	{{- range .SoftwareSystem.Containers }}
	{{ template "containerT" (containerData $.Model . $.SoftwareSystem.Name) }}
	{{- end }}
	{{- range .SoftwareSystem.Relationships }}
	{{ template "useT" (relData $.Model . .SoftwareSystem.Name) }}
	{{- end }}
}{{ end }})`

var containerT = `Container("{{ .Container.Name }}"{{ if .Container.Description }}, {{ printf "%q" .Container.Description }}{{ end }}{{ if .Container.Technology }}, {{ printf "%q" .Container.Technology }}{{ end }}{{ if containerHasFunc .Container }}, func() {
	{{- range filterTags .Container.Tags }}
	Tag({{ printf "%q" . }})
	{{- end }}
	{{- if .Container.URL }}
	URL({{ printf "%q" .Container.URL }})
	{{- end }}
	{{- range $k, $v := .Container.Properties }}
	Prop({{ printf "%q" $k }}, {{ printf "%q" $v }})
	{{- end }}
	{{- range .Container.Components }}
	{{ template "componentT" (componentData $.Model . (printf "%s/%s" $.CurrentPath  $.Container.Name)) }}
	{{- end }}
	{{- range .Container.Relationships }}
	{{ template "useT" (relData $.Model . $.CurrentPath) }}
	{{- end }}
}{{ end }})`

var componentT = `Component("{{ .Component.Name }}"{{ if .Component.Description }}, {{ printf "%q" .Component.Description }}{{ end }}{{ if .Component.Technology }}, {{ printf "%q" .Component.Technology }}{{ end }}{{ if componentHasFunc .Component }}, func() {
	{{- range filterTags .Component.Tags }}
	Tag({{ printf "%q" . }})
	{{- end }}
	{{- if .Component.URL }}
	URL({{ printf "%q" .Component.URL }})
	{{- end }}
	{{- range $k, $v := .Component.Properties }}
	Prop({{ printf "%q" $k }}, {{ printf "%q" $v }})
	{{- end }}
	{{- range .Component.Relationships }}
	{{ template "useT" (relData $.Model . $.CurrentPath) }}
	{{- end }}
}{{ end }})`

var personT = `Person("{{ .Person.Name }}"{{ if .Person.Description }}{{ printf "%q" .Person.Description}}{{ end }}{{ if personHasFunc .Person }}, func() {
	{{- range filterTags .Person.Tags }}
	Tag({{ printf "%q" . }})
	{{- end }}
	{{- if .Person.URL }}
	URL({{ printf "%q" .Person.URL }})
	{{- end }}
	{{- range $k, $v := .Person.Properties }}
	Prop({{ printf "%q" $k }}, {{ printf "%q" $v }})
	{{- end }}
	{{- range .Person.Relationships }}
	{{ template "useT" (relData $.Model . .Person.Name) }}
	{{- end }}
}{{ end }})`

var useT = `{{ relDSLFunc .Model .Relationship }}{{ with .Relationship }}("{{ elementPath $.Model .DestinationID $.CurrentPath }}"{{ if .Description}}, {{ printf "%q" .Description }}{{ end }}{{ if .Technology}}, {{ printf "%q" .Technology }}{{ end }}{{ if eq .InteractionStyle 1 }}, Synchronous{{ end }}{{ if eq .InteractionStyle 2 }}, Asynchronous{{ end }}{{ if or (filterTags .Tags) .URL }}, func() {
	{{- range filterTags .Tags }}
	Tag({{ printf "%q" . }})
	{{- end }}
	{{- if .URL }}
	URL({{ printf "%q" .URL }})
	{{- end }}
}{{ end }}){{ end }}`

var deploymentEnvironmentT = `DeploymentEnvironment({{ printf "%q" .Environment }}, func() {
	{{ template "deploymentNodeT" . }}
})`

var deploymentNodeT = `DeploymentNode("{{.Name}}"{{ if .Description}}, {{ printf "%q" .Description }}{{ end }}{{ if .Technology}}, {{ printf "%q" .Technology }}{{ end }}, func() {
	{{- range filterTags .Tags }}
	Tag({{ printf "%q" . }})
	{{- end }}
	{{- if .URL }}
	URL({{ printf "%q" .URL }})
	{{- end }}
	{{- if .Instances }}
	Instances({{ .Instances }})
	{{- end }}
	{{- range $k, $v := .Properties }}
	Prop({{ printf "%q" $k }}, {{ printf "%q" $v }})
	{{- end }}
	{{- range .Relationships }}
	{{ template "useT" (relData $.Model .  "") }}
	{{- end }}
	{{- range .Children }}
	{{ template "deploymentNodeT" . }}
	{{- end }}
	{{- range .InfrastructureNodes }}
	{{ template "infrastructureNodeT" . }}
	{{- end }}
	{{- range .ContainerInstances }}
	{{ template "containerInstanceT" . }}
	{{- end }}
})`

var infrastructureNodeT = `InfrastructureNode("{{.Name}}"{{ if .Description}}, {{ printf "%q" .Description }}{{ end }}{{ if .Technology}}, {{ printf "%q" .Technology }}{{ end }}{{ if or .Tags .URL .Properties .Relationships }}, func() {
	{{- range filterTags .Tags }}
	Tag({{ printf "%q" . }})
	{{- end }}
	{{- if .URL }}
	URL({{ printf "%q" .URL }})
	{{- end }}
	{{- range $k, $v := .Properties }}
	Prop({{ printf "%q" $k }}, {{ printf "%q" $v }})
	{{- end }}
	{{- range .Relationships }}
	{{ template "useT" (relData $.Model . "") }}
	{{- end }}
}{{ end }})`

var containerInstanceT = `ContainerInstance("{{ elementPath $.Model .ContainerID .CurrentPath }}", func() {
	InstanceID({{ .InstanceID }})
	{{- range filterTags .Tags "ContainerInstance" }}
	Tag({{ printf "%q" . }})
	{{- end }}
	{{- if .URL }}
	URL({{ printf "%q" .URL }})
	{{- end }}
	{{- range $k, $v := .Properties }}
	Prop({{ printf "%q" $k }}, {{ printf "%q" $v }})
	{{- end }}
	{{- range .HealthChecks }}
	{{ template "healthCheckT" . }}
	{{- end }}
})`

var healthCheckT = `HealthCheck({{ printf "%q" .Name }}, func() {
	URL({{ printf "%q" .URL }})
	Interval({{ .Interval }})
	Timeout({{ .Timeout }})
	{{- range $k, $v := .Headers }}
	Header({{ printf "%q" $k }}, {{ printf "%q" $v }})
	{{- end }}
})`

var viewPropsT = `{{ with .Props }}Title({{ printf "%q" .Title }})
{{- if .PaperSize }}
PaperSize({{ .PaperSize.Name }})
{{- end }}
{{- if .AutoLayout }}{{ with .AutoLayout }}
AutoLayout({{ .RankDirection.Name }}{{ if autoLayoutHasFunc . }}, func () {
	{{- if and .RankSep  (ne (deref .RankSep) $.DefaultRankSeparation) }}
	RankSeparation({{ .RankSep }})
	{{- end }}
	{{- if and .NodeSep (ne (deref .NodeSep) $.DefaultNodeSeparation) }}
	NodeSeparation({{ .NodeSep }})
	{{- end }}
	{{- if and .EdgeSep (ne (deref .EdgeSep) $.DefaultEdgeSeparation) }}
	EdgeSeparation({{ .EdgeSep }})
	{{- end }}
	{{- if .Vertices }}
	RenderVertices()
	{{- end }}
}{{ end }})
{{- end }}{{- end }}
{{- with .Settings }}
{{- if .AddAll }}
	AddAll()
{{- end }}
{{- if .AddDefault }}
	AddDefault()
{{- end }}
{{- range .AddNeighborIDs }}
	AddNeighbors("{{ elementPath $.Model . }}")
{{- end }}
{{- range .RemoveElementIDs }}
	Remove("{{ elementPath $.Model . }}")
{{- end }}
{{- range .RemoveTags }}
	RemoveTagged({{ printf "%q" . }})
{{- end }}
{{- range .RemoveRelationshipIDs }}
	{{- $rel := findRelationship $ . }}
	{{- if $rel }}
	Unlink("{{ elementPath $.Model $rel.SourceID }}", "{{ elementPath $.Model $rel.DestinationID }}"{{ if $rel.Description }}, {{ printf "%q" $rel.Description }}{{ end }})
	{{- end }}
{{- end }}
{{- range .RemoveUnreachableIDs }}
	RemoveUnreachable("{{ elementPath $.Model .ID }}")
{{- end }}
{{- if .RemoveUnrelated }}
	RemoveUnrelated()
{{- end }}
{{- end }}
{{- range .ElementViews }}
	Add("{{ elementPath $.Model .ID }}"{{ if .X }}, func() {
		Coord({{ .X }}, {{ .Y }})
	}{{ end }})
{{- end }}
{{- range .RelationshipViews }}
	{{- if .Source }}
		Link("{{ elementPath $.Model .Source.ID }}", "{{ elementPath $.Model .Destination.ID }}"{{ if .Description }}, {{ printf "%q" .Description }}{{ end }}{{ if .Order }}, {{ printf "%q" .Order }}{{ end }}{{ if .Vertices }}, func() {
			{{- range .Vertices }}
			Vertex({{ .X }}, {{ .Y }})
			{{- end }}
		}{{ end }}{{ if .Routing }}, {{ .Routing.Name }}{{ end }}{{ if .Position }}, {{ .Position }}{{ end }})
	{{- else }}
		Link("{{ elementPath $.Model .Destination.ID }}"{{ if .Description }}, {{ printf "%q" .Description }}{{ end }}{{ if .Order }}, {{ printf "%q" .Order }}{{ end }}{{ if .Vertices }}, func() {
			{{- range .Vertices }}
			Vertex({{ .X }}, {{ .Y }})
			{{- end }}
		}{{ end }}{{ if .Routing }}, {{ .Routing.Name }}{{ end }}{{ if .Position }}, {{ .Position }}{{ end }})
	{{- end }}
{{- end }}
{{- range .Animations }}
	AnimationStep({{ range .Elements }}"{{ elementPath .GetElement.ID }}", {{ end }})
{{- end }}
{{- end }}`

var systemLandscapeViewT = `SystemLandscapeView("{{.View.Key}}"{{ if .View.Description}}, {{ printf "%q" .View.Description }}{{ end }}, func() {
	{{ template "viewPropsT" (viewPropsData $.Model .View.ViewProps) }}
	{{- if .View.EnterpriseBoundaryVisible }}
	EnterpriseBoundaryVisible()
	{{- end }}
})`

var systemContextViewT = `SystemContextView("{{ elementPath .View.SoftwareSystemID }}", "{{ .View.Key }}"{{ if .View.Description}}, {{ printf "%q" .View.Description }}{{ end }}, func() {
	{{ template "viewPropsT" (viewPropsData .Model .View.ViewProps) }}
	{{- if .View.EnterpriseBoundaryVisible }}
	EnterpriseBoundaryVisible()
	{{- end }}
})`

var containerViewT = `ContainerView("{{ elementPath .View.SoftwareSystemID }}", "{{ .View.Key }}"{{ if .View.Description}}, {{ printf "%q" .View.Description }}{{ end }}, func() {
	{{ template "viewPropsT" (viewPropsData .Model .View.ViewProps) }}
	{{- if .View.SystemBoundaryVisible }}
	SystemBoundaryVisible()
	{{- end }}
})`

var componentViewT = `ComponentView("{{ elementPath .View.SoftwareSystemID }}", "{{ .View.Key }}"{{ if .View.Description}}, {{ printf "%q" .View.Description }}{{ end }}, func() {
	{{ template "viewPropsT" (viewPropsData .Model .View.ViewProps) }}
	{{- if .View.ContainerBoundaryVisible }}
	ContainerBoundaryVisible()
	{{- end }}
})`

var filteredViewT = `FilteredView("{{ .Key }}", func() {
	{{- range .FilterTags }}
	FilterTag({{ printf "%q" . }})
	{{- end }}
	{{- if .Exclude }}
	Exclude()
	{{- end }}
})`

var deploymentViewT = `DeploymentView("{{ elementPath .SoftwareSystemID }}", {{ printf "%q" .Environment }}, {{ printf "%q" .Key }}{{ if .Description}}, {{ printf "%q" .Description }}{{ end }}, func() {
	{{ template "viewPropsT" . }}
})`

var dynamicViewT = `DynamicView("{{ elementPath .ElementID }}", {{ printf "%q" .Key }}, func() {
	{{ template "viewPropsT" . }}
})`

var styleT = `Style(func() {
	{{- range .Elements }}
	ElementStyle({{ printf "%q" .Tag }}, func() {
		{{- if gt .Shape 0 }}
		Shape("{{ .Shape.Name }}")
		{{- end }}
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
		{{- if gt .Border 0 }}
		Border("{{ .Border.Name }}")
		{{- end }}
	})
	{{- end }}
	{{- range .Relationships }}
	RelationshipStyle({{ printf "%q" .Tag }}, func() {
		{{- if .Thickness }}
		Thickness({{ .Thickness }})
		{{- end }}
		{{- if .FontSize }}
		FontSize({{ .FontSize }})
		{{- end }}
		{{- if .Width    }}
		Width({{ .Width }})
		{{- end }}
		{{- if .Position }}
		Position({{ .Position }})
		{{- end }}
		{{- if .Color    }}
		Color({{ .Color }})
		{{- end }}
		{{- if .Stroke   }}
		Stroke({{ .Stroke }})
		{{- end }}
		{{- if .Dashed   }}
		Dashed()
		{{- end }}
		{{- if ge .Routing 0 }}
		Routing({{ .Routing.Name }})
		{{- end }}
		{{- if .Opacity  }}
		Opacity({{ .Opacity }})
		{{- end }}
	})
	{{- end }}
})`
