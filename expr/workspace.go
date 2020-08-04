package expr

import (
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
)

type (
	// Workspace describes a workspace and is the root expression of the plugin.
	Workspace struct {
		// ID of workspace.
		ID int `json:"id,omitempty"`
		// Name of workspace.
		Name string `json:"name"`
		// Description of workspace if any.
		Description string `json:"description,omitempty"`
		// Version number for the workspace.
		Version string `json:"version,omitempty"`
		// Revision number, automatically generated.
		Revision int `json:"revision,omitempty"`
		// Thumbnail associated with the workspace; a Base64 encoded PNG file as a
		// data URI (data:image/png;base64).
		Thumbnail string `json:"thumbnail,omitempty"`
		// The last modified date, in ISO 8601 format (e.g. "2018-09-08T12:40:03Z").
		LastModifiedDate string `json:"lastModifiedDate,omitempty"`
		// A string identifying the user who last modified the workspace (e.g. an
		// e-mail address or username).
		LastModifiedUser string `json:"lastModifiedUser,omitempty"`
		//  A string identifying the agent that was last used to modify the workspace
		//  (e.g. "structurizr-java/1.2.0").
		LastModifiedAgent string `json:"lastModifiedAgent,omitempty"`
		// Model is the software architecture model.
		Model *Model `json:"model,omitempty"`
		// Views contains the views if any.
		Views *Views `json:"views,omitempty"`
		// Documentation associated with software architecture model.
		Documentation *Documentation `json:"documentation,omitempty"`
		// Configuration of workspace.
		Configuration *WorkspaceConfiguration `json:"configuration,omitempty"`
	}

	// WorkspaceConfiguration describes the workspace configuration.
	WorkspaceConfiguration struct {
		// Users that have access to the workspace.
		Users []*User `json:"users"`
	}

	// User of Structurizr service.
	User struct {
		Username string `json:"username"`
		// Role of user, one of "ReadWrite" or "ReadOnly".
		Role string `json:"role"`
	}
)

// Root is the design root expression.
var Root = &Workspace{Model: &Model{}, Views: &Views{}}

// Register design root with eval engine.
func init() {
	eval.Register(Root)
	eval.Register(Root.Views)
}

// WalkSets iterates over the elements and views.
// Elements DSL cannot be executed on init because all elements must first be
// loaded and their IDs captured in the registry before relationships can be
// built with DSL.
func (w *Workspace) WalkSets(walk eval.SetWalker) {
	// 1. Model
	walk([]eval.Expression{w.Model})
	// 2. People
	walk(eval.ToExpressionSet(w.Model.People))
	// 3. Systems
	walk(eval.ToExpressionSet(w.Model.Systems))
	// 4. Containers
	for _, s := range w.Model.Systems {
		walk(eval.ToExpressionSet(s.Containers))
	}
	// 5. Components
	for _, s := range w.Model.Systems {
		for _, c := range s.Containers {
			walk(eval.ToExpressionSet(c.Components))
		}
	}
	// 6. Deployment environments
	walkDeploymentNodes(w.Model.DeploymentNodes, walk)
}

func walkDeploymentNodes(n []*DeploymentNode, walk eval.SetWalker) {
	if n == nil {
		return
	}
	walk(eval.ToExpressionSet(n))
	for _, d := range n {
		walk(eval.ToExpressionSet(d.InfrastructureNodes))
		walk(eval.ToExpressionSet(d.ContainerInstances))
		walkDeploymentNodes(d.Children, walk)
	}
}

// DependsOn tells the eval engine to run the goa DSL first.
func (w *Workspace) DependsOn() []eval.Root { return []eval.Root{expr.Root} }

// Packages returns the import path to the Go packages that make
// up the DSL. This is used to skip frames that point to files
// in these packages when computing the location of errors.
func (w *Workspace) Packages() []string { return w.Views.Packages() }

// EvalName returns the generic expression name used in error messages.
func (w *Workspace) EvalName() string {
	return "Structurizr workspace"
}

// Person returns the person with the given name if any, nil otherwise.
func (w *Workspace) Person(name string) *Person {
	return w.Model.Person(name)
}

// SoftwareSystem returns the software system with the given name if any, nil
// otherwise.
func (w *Workspace) SoftwareSystem(name string) *SoftwareSystem {
	return w.Model.SoftwareSystem(name)
}

// DeploymentNode returns the deployment node with the given name if any, nil
// otherwise.
func (w *Workspace) DeploymentNode(name string) *DeploymentNode {
	return w.Model.DeploymentNode(name)
}
