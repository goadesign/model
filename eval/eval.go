/*
Package eval allows evaluating Go code that describes a software architecture
model using the DSL defined in Structurizr for Go:
https://github.com/goadesign/structurizr.
*/
package eval

import (
	"goa.design/goa/v3/eval"
	"goa.design/structurizr/expr"
)

// RunDSL runs the DSL stored in a global variable and returns the corresponding
// Workspace expression. The expression can be serialized to JSON to obtain a
// representation that is compatible with the Structurizr JSON schema.
func RunDSL() (*expr.Workspace, error) {
	if err := eval.RunDSL(); err != nil {
		return nil, err
	}
	return expr.Root, nil
}
