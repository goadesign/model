package expr

import (
	"fmt"

	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
	goa "goa.design/goa/v3/pkg"
)

// Workspace describes a workspace and is the root expression of the plugin.
type Workspace struct {
	// ID of workspace.
	ID string `json:"id"`
	// Name of workspace.
	Name string `json:"name"`
	// Description of workspace if any.
	Description string `json:"description"`
	// Version number for the workspace.
	Version string `json:"version"`
	// Thumbnail associated with the workspace; a Base64 encoded PNG file as a
	// data URI (data:image/png;base64).
	Thumbnail string `json:"thumbnail"`
	// The last modified date, in ISO 8601 format (e.g. "2018-09-08T12:40:03Z").
	LastModifiedDate string `json:"lastModifiedData"`
	// A string identifying the user who last modified the workspace (e.g. an
	// e-mail address or username).
	LastModifiedUser string `json:"lastModifiedUser"`
	//  A string identifying the agent that was last used to modify the workspace
	//  (e.g. "structurizr-java/1.2.0").
	LastModifiedAgent string `json:"lastModifiedAgent"`
	// Model is the software architecture model.
	Model *Model `json:"model"`
	// Views contains the views if any.
	Views *Views `json:"views"`
	// Documentation associated with software architecture model.
	Documentation *Documentation `json:"documentation"`
	// Configuration of workspace.
	Configuration *Configuration `json:"configuration"`
}

// Root is the design root expression.
var Root = &Workspace{}

// Register design root with eval engine.
func init() {
	eval.Register(Root)
}

// WalkSets iterates over the views, elements are completely evaluated during
// init.
func (w *Workspace) WalkSets(walk eval.SetWalker) {
	if w.Views == nil {
		return
	}
	walk([]eval.Expression{w.Views})
}

// DependsOn tells the eval engine to run the goa DSL first.
func (w *Workspace) DependsOn() []eval.Root { return []eval.Root{expr.Root} }

// Packages returns the import path to the Go packages that make
// up the DSL. This is used to skip frames that point to files
// in these packages when computing the location of errors.
func (w *Workspace) Packages() []string {
	return []string{
		"goa.design/plugins/v3/structurizr/expr",
		"goa.design/plugins/v3/structurizr/dsl",
		fmt.Sprintf("goa.design/plugins/v3@%s/structurizr/expr", goa.Version()),
		fmt.Sprintf("goa.design/plugins/v3@%s/structurizr/dsl", goa.Version()),
	}
}

// EvalName returns the generic expression name used in error messages.
func (w *Workspace) EvalName() string {
	return "Structurizr workspace"
}
