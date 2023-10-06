package plugin

import (
	"encoding/json"
	"os"
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

// Generate produces the design JSON representation inside the top level gen
// directory.
func Generate(_ string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	err := eval.RunDSL()
	if err != nil {
		return nil, err
	}

	// Generate structurizr JSON
	path := filepath.Join(codegen.Gendir, "structurizr", "model.json")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		if err := os.Remove(path); err != nil {
			return nil, err
		}
	}
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
