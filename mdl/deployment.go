package mdl

import (
	"goa.design/goa/v3/codegen"
	"goa.design/model/expr"
)

func deploymentNodeSections(dv *expr.DeploymentView, dn *expr.DeploymentNode, indent int) []*codegen.SectionTemplate {
	var inst int
	if dn.Instances != nil {
		inst = *dn.Instances
	}
	sections := []*codegen.SectionTemplate{{
		Name:    "deploymentNodeStart",
		Source:  deploymentNodeStartT,
		FuncMap: funcs,
		Data: struct {
			ID           string
			Indent       int
			BoundaryName string
			Instances    int
		}{dn.ID, indent, dn.Name, inst},
	}}
	var evs []*expr.ElementView
	for _, inf := range dn.InfrastructureNodes {
		if infv := findElement(dv, inf.Element); infv != nil {
			evs = append(evs, infv)
		}
	}
	for _, ci := range dn.ContainerInstances {
		if civ := findElement(dv, ci.Element); civ != nil {
			evs = append(evs, civ)
		}
	}
	sections = append(sections, elements(evs, "", indent+1))
	for _, c := range dn.Children {
		sections = append(sections, deploymentNodeSections(dv, c, indent+1)...)
	}
	sections = append(sections, &codegen.SectionTemplate{
		Name:    "deploymentNodeEnd",
		Source:  deploymentNodeEndT,
		FuncMap: funcs,
		Data: struct {
			ID     string
			Indent int
		}{dn.ID, indent},
	})
	return sections
}

func findElement(dv *expr.DeploymentView, elem *expr.Element) *expr.ElementView {
	id := elem.ID
	for _, ev := range dv.ElementViews {
		if ev.Element.ID == id {
			return ev
		}
	}
	return nil
}

const deploymentNodeStartT = `{{ indent .Indent }}subgraph {{ .ID }} [{{ .BoundaryName }}{{ if gt .Instances 1 }} x{{ .Instances }}{{ end }}]
`

const deploymentNodeEndT = `{{ indent .Indent }}end
%%{{ indent .Indent }}style {{ .ID }} fill:#ffffff,stroke:#606060,color:#000000;
`
