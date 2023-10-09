package codegen

import (
	"fmt"
	"path/filepath"
	"strings"

	"goa.design/goa/v3/codegen"
	"golang.org/x/tools/go/packages"

	"goa.design/model/mdl"
)

// Model generates the model DSL from the given design.
func Model(d *mdl.Design, pkg string) error {
	file, err := findModelFile(pkg)
	if err != nil {
		return fmt.Errorf("failed to find model file: %s", err)
	}
	var template string
	for name, tmpl := range templates {
		template += fmt.Sprintf("{{define %q}}%s{{end}}", name, tmpl)
	}
	template += designT
	sections := []*codegen.SectionTemplate{
		codegen.Header(d.Name, "model", []*codegen.ImportSpec{codegen.NewImport(".", "goa.design/model/dsl")}),
		{
			Name:   "design",
			Source: template,
			Data:   d,
			FuncMap: map[string]any{
				"relDSL":           relDSL,
				"elementPath":      elementPath,
				"elementPaths":     elementPaths,
				"findRelationship": findRelationship,
			},
		},
	}
	cf := &codegen.File{Path: filepath.Base(file), SectionTemplates: sections}
	if _, err := cf.Render(filepath.Dir(file)); err != nil {
		return fmt.Errorf("failed to render model file: %s", err)
	}
	return nil
}

// findModelFile finds the package directory for the given package name.
func findModelFile(name string) (string, error) {
	cfg := &packages.Config{
		Mode: packages.NeedFiles,
	}
	pkgs, err := packages.Load(cfg, name)
	if err != nil {
		return "", err
	}
	if len(pkgs) == 0 {
		return "", fmt.Errorf("package %q not found", name)
	}
	if len(pkgs[0].GoFiles) == 0 {
		return "", fmt.Errorf("package %q does not contain any Go file", name)
	}
	return pkgs[0].GoFiles[0], nil
}

// relDSL is a function used by the DSL codegen to compute the name of the DSL function used to represent the corresponding relationship, one of "Uses", "Delivers" or "InteractsWith".
func relDSL(mod *mdl.Model, rel *mdl.Relationship) string {
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

// elementPath is used by templates codegen to compute the path to the element with the given ID.
func elementPath(mod *mdl.Model, id string) string {
	for _, p := range mod.People {
		if p.ID == id {
			return fmt.Sprintf("%q", p.Name)
		}
	}
	for _, s := range mod.Systems {
		if s.ID == id {
			return fmt.Sprintf("%q", s.Name)
		}
		for _, c := range s.Containers {
			if c.ID == id {
				return fmt.Sprintf("%q/%q", s.Name, c.Name)
			}
			for _, cmp := range c.Components {
				if cmp.ID == id {
					return fmt.Sprintf("%q/%q/%q", s.Name, c.Name, cmp.Name)
				}
			}
		}
	}
	return ""
}

// ElementPaths is used by templates to compute a comma separated list of element paths.
func elementPaths(mod *mdl.Model, ids []string) string {
	res := make([]string, len(ids))
	for _, id := range ids {
		res = append(res, elementPath(mod, id))
	}
	return strings.Join(res, ", ")
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

var designT = `package design

import . "goa.design/model/dsl"

var _ = Design("{{.Name}}", "{{.Description}}", func() {
	{{- if .Enterprise }}
	Enterprise({{ printf "%q" .Enterprise }})
	{{ end }}
	{{- range .Model.People }}
	{{ template "personT" . }}
	{{ end }}
	{{- range .Model.Systems }}
	{{ template "systemT" . }}
	{{ end }}
	{{- range .Model.DeploymentNodes }}
	{{ template "deploymentEnvironmentT" . }}
	{{ end }}
	Views(func() {
	{{- range .Views.LandscapeViews }}
	{{ template "systemLandscapeViewT" . }}
	{{ end }}
	{{- range .Views.ContextViews }}
	{{ template "systemContextViewT" . }}
	{{ end }}
	{{- range .Views.ContainerViews }}
	{{ template "containerViewT" . }}
	{{ end }}
	{{- range .Views.ComponentViews }}
	{{ template "componentViewT" . }}
	{{ end }}
	{{- range .Views.DeploymentViews }}
	{{ template "deploymentViewT" . }}
	{{ end }}
	{{- range .Views.FilteredViews }}
	{{ template "filteredViewT" . }}
	{{ end }}
	{{- if .Views.Styles }}
	{{ template "styleT" .Views.Styles }}
	{{ end }}
	})
})`

var systemT = `SoftwareSystem("{{.Name}}", "{{.Description}}", func() {
	{{- range .Tags }}
	Tag({{ printf "%q" . }})
	{{ end }}
	{{- if .URL }}
	URL({{ printf "%q" .URL }})
	{{ end }}
	{{- if eq .Location 2 }}
	External()
	{{ end }}
	{{- range $k, $v := .Properties }}
	Prop({{ printf "%q" $k }}, {{ printf "%q" $v }})
	{{ end }}
	{{- range .Relationships }}
	{{ template "useT" . }}
	{{ end }}
	{{- range .Containers }}
	{{ template "containerT" . }}
	{{ end }}
})`

var containerT = `Container("{{.Name}}"{{ if .Description }}, {{ printf "%q".Description }}{{ end }}{{ if .Technology }}, {{ printf "%q" .Technology }}{{ end }}", func() {
	{{- range .Tags }}
	Tag({{ printf "%q" . }})
	{{ end }}
	{{- if .URL }}
	URL({{ printf "%q" .URL }})
	{{ end }}
	{{- range $k, $v := .Properties }}
	Prop({{ printf "%q" $k }}, {{ printf "%q" $v }})
	{{ end }}
	{{- range .Relationships }}
	{{ template "useT" . }}
	{{ end }}
	{{- range .Components }}
	{{ template "componentT" . }}
	{{ end }}
})`

var componentT = `Component("{{.Name}}"{{ if .Description }}, {{ printf "%q".Description }}{{ end }}{{ if .Technology }}, {{ printf "%q" .Technology }}{{ end }}", func() {
	{{- range .Tags }}
	Tag({{ printf "%q" . }})
	{{ end }}
	{{- if .URL }}
	URL({{ printf "%q" .URL }})
	{{ end }}
	{{- if eq .Location 2 }}
	External()
	{{ end }}
	{{- range $k, $v := .Properties }}
	Prop({{ printf "%q" $k }}, {{ printf "%q" $v }})
	{{ end }}
	{{- range .Relationships }}
	{{ template "useT" . }}
	{{ end }}
})`

var personT = `Person("{{.Name}}"{{ if .Description }}{{ printf "%q" .Description}}{{ end }}, func() {
	{{- range .Tags }}
	Tag({{ printf "%q" . }})
	{{ end }}
	{{- if .URL }}
	URL({{ printf "%q" .URL }})
	{{ end }}
	{{- range $k, $v := .Properties }}
	Prop({{ printf "%q" $k }}, {{ printf "%q" $v }})
	{{ end }}
	{{- range .Relationships }}
	{{ template "useT" . }}
	{{ end }}
})`

var useT = `{{ relDSL $.Model . }}("{{.Name}}"{{ if .Description}}, {{ printf "%q" .Description }}{{ end }}{{ if .Technology}}, {{ printf "%q" .Technology }}{{ end }}{{ if eq .InteractionStyle 1 }}, Synchronous{{ end }}{{ if eq .InteractionStyle 2 }}, Asynchronous{{ end }}{{ if or .Tags .URL }}, func() {
	{{- range .Tags }}
	Tag({{ printf "%q" . }})
	{{ end }}
	{{- if .URL }}
	URL({{ printf "%q" .URL }})
	{{ end }}
}{{ end }})`

var deploymentEnvironmentT = `DeploymentEnvironment({{ printf "%q" .Environment }}, func() {
	{{ template "deploymentNodeT" . }}
})`

var deploymentNodeT = `DeploymentNode("{{.Name}}"{{ if .Description}}, {{ printf "%q" .Description }}{{ end }}{{ if .Technology}}, {{ printf "%q" .Technology }}{{ end }}, func() {
	{{- range .Tags }}
	Tag({{ printf "%q" . }})
	{{ end }}
	{{- if .URL }}
	URL({{ printf "%q" .URL }})
	{{ end }}
	{{- if .Instances }}
	Instances({{ .Instances }})
	{{ end }}
	{{- range $k, $v := .Properties }}
	Prop({{ printf "%q" $k }}, {{ printf "%q" $v }})
	{{ end }}
	{{- range .Relationships }}
	{{ template "useT" . }}
	{{ end }}
	{{- range .Children }}
	{{ template "deploymentNodeT" . }}
	{{ end }}
	{{- range .InfrastructureNodes }}
	{{ template "infrastructureNodeT" . }}
	{{ end }}
	{{- range .ContainerInstances }}
	{{ template "containerInstanceT" . }}
	{{ end }}
})`

var infrastructureNodeT = `InfrastructureNode("{{.Name}}"{{ if .Description}}, {{ printf "%q" .Description }}{{ end }}{{ if .Technology}}, {{ printf "%q" .Technology }}{{ end }}{{ if or .Tags .URL .Properties .Relationships }}, func() {
	{{- range .Tags }}
	Tag({{ printf "%q" . }})
	{{ end }}
	{{- if .URL }}
	URL({{ printf "%q" .URL }})
	{{ end }}
	{{- range $k, $v := .Properties }}
	Prop({{ printf "%q" $k }}, {{ printf "%q" $v }})
	{{ end }}
	{{- range .Relationships }}
	{{ template "useT" . }}
	{{ end }}
}{{ end }})`

var containerInstanceT = `ContainerInstance("{{ elementPath $.Model .ContainerID }}", func() {
	InstanceID({{ .InstanceID }})
	{{- range .Tags }}
	Tag({{ printf "%q" . }})
	{{ end }}
	{{- if .URL }}
	URL({{ printf "%q" .URL }})
	{{ end }}
	{{- range $k, $v := .Properties }}
	Prop({{ printf "%q" $k }}, {{ printf "%q" $v }})
	{{ end }}
	{{- range .HealthChecks }}
	{{ template "healthCheckT" . }}
	{{ end }}
})`

var healthCheckT = `HealthCheck({{ printf "%q" .Name }}, func() {
	URL({{ printf "%q" .URL }})
	Interval({{ .Interval }})
	Timeout({{ .Timeout }})
	{{- range $k, $v := .Headers }}
	Header({{ printf "%q" $k }}, {{ printf "%q" $v }})
	{{ end }}
})`

var viewPropsT = `Title({{ printf "%q" .Title }})
{{ if .PaperSize }}PaperSize({{ .PaperSize.Name }}){{ end }}
{{ if .AutoLayout }}AutoLayout({{ .AutoLayout.RankDirection.Name }}{{ if or .RankSep .NodeSep .EdgeSep .Vertices }}, func () {
	{{- if .RankSep }}RankSeparation({{ .RankSep }}){{ end }}
	{{- if .NodeSep }}NodeSeparation({{ .NodeSep }}){{ end }}
	{{- if .EdgeSep }}EdgeSeparation({{ .EdgeSep }}){{ end }}
	{{- if .Vertices }}RenderVertices(){{ end }}
}{{ end }}){{ end }}
{{- if .ViewSettings.AddAll }}
	AddAll()
{{ end }}
{{- if .ViewSettings.AddDefault }}
	AddDefault()
{{ end }}
{{- range .ViewSettings.AddNeighborIDs }}
	AddNeighbors({{ elementPath $.Model . }})
{{ end }}
{{- range .ViewSettings.RemoveElementIDs }}
	Remove({{ elementPath $.Model .ID }})
{{ end }}
{{- range .ViewSettings.RemoveTags }}
	RemoveTagged({{ printf "%q" . }})
{{ end }}
{{- range .ViewSettings.RemoveRelationshipIDs }}
	{{- $rel := findRelationship $ . }}
	{{- if $rel }}
	Unlink({{ elementPath $.Model $rel.SourceID }}, {{ elementPath $.Model $rel.DestinationID }}{{ if $rel.Description }}, {{ printf "%q" $rel.Description }}{{ end }})
	{{- end }}
{{ end }}
{{- range .ViewSettings.RemoveUnreachable }}
	RemoveUnreachable({{ elementPath $.Model .ID }})
{{ end }}
{{- if .ViewSettings.RemoveUnrelated }}
	RemoveUnrelated()
{{ end }}
{{- range .ElementViews }}
	Add({{ elementPath $.Model .ID }}{{ if .X }}, func() {
		Coord({{ .X }}, {{ .Y }})
	}{{ end }})
{{ end }}
{{- range .RelationshipViews }}
	{{- if .Source }}
		Link({{ elementPath $.Model .Source.ID }}, {{ elementPath $.Model .Destination.ID }}{{ if .Description }}, {{ printf "%q" .Description }}{{ end }}{{ if .Order }}, {{ printf "%q" .Order }}{{ end }}{{ if .Vertices }}, func() {
			{{- range .Vertices }}
			Vertex({{ .X }}, {{ .Y }})
			{{ end }}
		}{{ end }}{{ if .Routing }}, {{ .Routing.Name }}{{ end }}{{ if .Position }}, {{ .Position }}{{ end }})
	{{- else }}
		Link({{ elementPath $.Model .Destination.ID }}{{ if .Description }}, {{ printf "%q" .Description }}{{ end }}{{ if .Order }}, {{ printf "%q" .Order }}{{ end }}{{ if .Vertices }}, func() {
			{{- range .Vertices }}
			Vertex({{ .X }}, {{ .Y }})
			{{ end }}
		}{{ end }}{{ if .Routing }}, {{ .Routing.Name }}{{ end }}{{ if .Position }}, {{ .Position }}{{ end }})
	{{- end }}
{{ end }}
{{- range .AnimationSteps }}
	AnimationStep({{ range .Elements }}{{ elementPath .GetElement.ID }}, {{ end }})
{{- end }}`

var systemLandscapeViewT = `SystemLandscapeView("{{.Key}}"{{ if .Description}}, {{ printf "%q" .Description }}{{ end }}, func() {
	{{ template "viewPropT" . }}
	{{- if .EnterpriseBoundaryVisible }}
	EnterpriseBoundaryVisible()
	{{- end }}
})`

var systemContextViewT = `SystemContextView({{ elementPath .SoftwareSystemID }}, "{{.Key}}"{{ if .Description}}, {{ printf "%q" .Description }}{{ end }}, func() {
	{{ template "viewPropT" . }}
	{{- if .EnterpriseBoundaryVisible }}
	EnterpriseBoundaryVisible()
	{{- end }}
})`

var containerViewT = `ContainerView({{ elementPath .SoftwareSystemID }}, "{{.Key}}"{{ if .Description}}, {{ printf "%q" .Description }}{{ end }}, func() {
	{{ template "viewPropT" . }}
	{{- if .SystemBoundaryVisible }}
	SystemBoundaryVisible()
	{{- end }}
})`

var componentViewT = `ComponentView({{ elementPath .SoftwareSystemID }}, "{{.Key}}"{{ if .Description}}, {{ printf "%q" .Description }}{{ end }}, func() {
	{{ template "viewPropT" . }}
	{{- if .ContainerBoundaryVisible }}
	ContainerBoundaryVisible()
	{{- end }}
})`

var filteredViewT = `FilteredView("{{.Key}}", func() {
	{{- range .FilterTags }}
	FilterTag({{ printf "%q" . }})
	{{- end }}
	{{-if .Exclude }}
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
		{{- if Background }}
		Background({{ printf "%q" .Background }})
		{{- end }}
		{{- if Color }}
		Color({{ printf "%q" .Color }})
		{{- end }}
		{{- if Stroke }}
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
