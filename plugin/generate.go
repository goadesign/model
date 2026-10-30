package plugin

import (
	"encoding/json"
	"path/filepath"
	"text/template"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
	"goa.design/model/expr"
	"goa.design/model/stz"
)

// init registers the plugin generator function.
func init() {
	codegen.RegisterPlugin("model", "gen", nil, Generate)
}

// Generate adds the design's Structurizr JSON representation, gen/structurizr/
// model.json, to files. Goa evaluates the design before it calls plugins, so
// Generate reads the evaluated model and neither runs the DSL again nor touches
// the file system: rendering the returned file replaces any previous one.
func Generate(_ string, _ []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	path := filepath.Join(codegen.Gendir, "structurizr", "model.json")
	section := &codegen.SectionTemplate{
		Name:    "model",
		FuncMap: template.FuncMap{"toJSON": toJSON},
		Source:  "{{ toJSON . }}",
		Data:    stz.WorkspaceFromDesign(expr.Root),
	}
	files = append(files, &codegen.File{
		Path:             path,
		SectionTemplates: []*codegen.SectionTemplate{section},
	})

	return files, nil
}

func toJSON(d any) string {
	b, err := json.MarshalIndent(d, "", "    ")
	if err != nil {
		panic("design: " + err.Error()) // bug
	}
	return string(b)
}
