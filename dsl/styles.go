package dsl

import (
	"goa.design/goa/v3/eval"
	"goa.design/structurizr/expr"
)

// Styles is a wrapper for one or more element/relationship styles,
// which are used when rendering diagrams.
//
// Styles must appear in Views.
//
// Styles accepts a single argument: a function that defines the styles.
func Styles(dsl func()) {
	vs, ok := eval.Current().(*expr.Views)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	cfg := &expr.Configuration{}
	eval.Execute(dsl, cfg)
	vs.Configuration = cfg
}
