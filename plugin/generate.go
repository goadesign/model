package docs

import (
	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
)

// init registers the plugin generator function.
func init() {
	codegen.RegisterPlugin("model", "gen", nil, Generate)
}

// Generate produces the documentation JSON file.
func Generate(_ string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	return files, nil
}
