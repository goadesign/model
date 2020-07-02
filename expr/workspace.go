package expr

import (
	"fmt"

	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
	goa "goa.design/goa/v3/pkg"
)

// WorkspaceExpr describes a workspace and is the root expression of the plugin.
type WorkspaceExpr struct {
	// ID of workspace.
	ID string `json:"id"`
	// Name of workspace.
	Name string `json:"name"`
	// Description of workspace if any.
	Description string `json:"description"`
	// Model is the software architecture model.
	Model *ModelExpr `json:"model"`
	// Views contains the views if any.
	Views *ViewsExpr `json:"views"`
}

// Root is the design root expression.
var Root = &WorkspaceExpr{}

// Register design root with eval engine.
func init() {
	eval.Register(Root)
}

// WalkSets iterates over the model and then the views.
func (w *WorkspaceExpr) WalkSets(walk eval.SetWalker) {
	if w.Model == nil {
		return
	}
	walk(eval.ToExpressionSet(w.Model.Systems))
	walk(eval.ToExpressionSet(w.Model.DeploymentNodes))
	if w.Views == nil {
		return
	}
	walk(eval.ToExpressionSet(w.Views.LandscapeViews))
	walk(eval.ToExpressionSet(w.Views.ContextViews))
	walk(eval.ToExpressionSet(w.Views.ContainerViews))
	walk(eval.ToExpressionSet(w.Views.ComponentViews))
	walk(eval.ToExpressionSet(w.Views.DynamicViews))
	walk(eval.ToExpressionSet(w.Views.FilteredViews))
	walk(eval.ToExpressionSet(w.Views.DeploymentViews))
}

// DependsOn tells the eval engine to run the goa DSL first.
func (w *WorkspaceExpr) DependsOn() []eval.Root { return []eval.Root{expr.Root} }

// Packages returns the import path to the Go packages that make
// up the DSL. This is used to skip frames that point to files
// in these packages when computing the location of errors.
func (w *WorkspaceExpr) Packages() []string {
	return []string{
		"goa.design/plugins/v3/structurizr/expr",
		"goa.design/plugins/v3/structurizr/dsl",
		fmt.Sprintf("goa.design/plugins/v3@%s/structurizr/expr", goa.Version()),
		fmt.Sprintf("goa.design/plugins/v3@%s/structurizr/dsl", goa.Version()),
	}
}
