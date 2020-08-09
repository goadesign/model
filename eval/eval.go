/*
Package eval allows evaluating Go code that describes a software architecture
model using the DSL defined in Model: https://github.com/goadesign/model.
*/
package eval

import (
	"goa.design/goa/v3/eval"
	"goa.design/model/expr"
)

// RunDSL runs the DSL stored in a global variable and returns the corresponding
// Workspace expression.
func RunDSL() (*expr.Workspace, error) {
	if err := eval.RunDSL(); err != nil {
		return nil, err
	}
	return expr.Root, nil
}
