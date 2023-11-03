package editor

import (
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	gentypes "goa.design/model/svc/gen/types"
)

type (
	// Editor exposes method to modify the DSL in a model package.
	Editor struct {
		repo string
		dir  string
	}

	// ElementKind is the kind of a DSL element.
	ElementKind string
)

const (
	// DesignKind is the kind of a DSL design.
	DesignKind ElementKind = "Design"
	// PersonKind is the kind of a person.
	PersonKind = "Person"
	// SoftwareSystemKind is the kind of a software system.
	SoftwareSystemKind = "SoftwareSystem"
	// ContainerKind is the kind of a container.
	ContainerKind = "Container"
	// ComponentKind is the kind of a component.
	ComponentKind = "Component"
	// ViewsKind is the kind of a DSL view.
	ViewsKind = "Views"
	// LandscapeViewKind is the kind of a landscape view.
	LandscapeViewKind = "SystemLandscapeView"
	// SystemContextViewKind is the kind of a system context view.
	SystemContextViewKind = "SystemContextView"
	// ContainerViewKind is the kind of a container view.
	ContainerViewKind = "ContainerView"
	// ComponentViewKind is the kind of a component view.
	ComponentViewKind = "ComponentView"
	// ElementStyleKind is the kind of a DSL element style.
	ElementStyleKind = "ElementStyle"
	// RelationshipStyleKind is the kind of a DSL relationship style.
	RelationshipStyleKind = "RelationshipStyle"
)

// DefaultModelFilename is the name of the DSL file used by default (e.g. when
// adding new DSL).
const DefaultModelFilename = "model.go"

// NewEditor returns a new Editor.
func NewEditor(repo, dir string) *Editor {
	return &Editor{repo: repo, dir: dir}
}

// UpsertElement updates the code for the DSL element with the given kind, path
// and arguments if there's a match otherwise it adds the element. The path is
// the path to the element in the DSL file, e.g. "Person1" or
// "SoftwareSystem1/Container1/Component1". Only non nil and non-empty arguments
// are used to match the element.
func (e *Editor) UpsertElement(code string, kind ElementKind, path string, args ...string) (*gentypes.PackageFile, error) {
	res, parsed, fset, err := e.parseDir()
	if err != nil {
		return nil, err
	}

	if parsed == nil {
		// New repository, add the DSL file.
		res.Locator.Filename = DefaultModelFilename
		modelFile := filepath.Join(e.repo, e.dir, DefaultModelFilename)
		return setContent(modelFile, res, newDesign(code))
	}

	// Find node and/or parent.
	var node, parent ast.Node
	var file *ast.File
	for _, f := range parsed.Files {
		node, parent = findDSL(f, kind, path, args...)
		if node != nil {
			file = f
			break
		} else if parent != nil {
			file = f
		}
	}

	if file != nil {
		// Load node or parent file.
		tokFile := fset.File(file.Pos())
		filename := tokFile.Name()
		res.Locator.Filename = filepath.Base(filename)
		srcBytes, err := os.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to read DSL file %s: %w", tokFile.Name(), err)
		}
		src := string(srcBytes)

		if node != nil {
			// Found the element, so replace it.
			code, err := copyRelationships(tokFile, node, src, code)
			if err != nil {
				return nil, err
			}
			newContent := concat(
				src[:tokFile.Offset(node.Pos())],
				code,
				src[tokFile.Offset(node.End()):],
			)
			return setContent(filename, res, newContent)
		}

		if parent != nil {
			// We found the parent, so insert the child.
			endOffset := tokFile.Offset(parent.End()) - 2 // -2 to account for `})`
			newContent := concat(
				src[:endOffset],
				code,
				src[endOffset:],
			)
			return setContent(filename, res, newContent)
		}
	}

	// We didn't find the element, so add it
	modelFile := filepath.Join(e.repo, e.dir, DefaultModelFilename)
	res.Locator.Filename = DefaultModelFilename
	if _, err := os.Stat(modelFile); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("failed to stat DSL file %s: %w", modelFile, err)
		}
		// New model, add the DSL file.
		return setContent(modelFile, res, newDesign(code))
	}
	file, err = parser.ParseFile(fset, modelFile, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSL file %s: %w", modelFile, err)
	}
	tokFile := fset.File(file.Pos())
	filename := tokFile.Name()
	res.Locator.Filename = filepath.Base(filename)
	srcBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read DSL file %s: %w", tokFile.Name(), err)
	}
	src := string(srcBytes)

	// Make sure the DSL file contains a Design() call.
	design, _ := findDSL(file, DesignKind, "")
	if design == nil {
		return nil, fmt.Errorf("failed to find Design() call in DSL file %s", modelFile)
	}

	// Where to insert the code depends on the kind:
	//  - Person, SoftwareSystem, Container, Component: insert last before Views DSL
	//  - Views: insert last in the Views DSL
	//  - Styles: insert last in the Styles DSL
	var insertOffset int
	switch kind {
	case PersonKind, SoftwareSystemKind, ContainerKind, ComponentKind:
		n, _ := findDSL(file, ViewsKind, "")
		if n != nil {
			insertOffset = fset.File(file.Pos()).Offset(n.Pos()) - 2 // -2 to account for `})`
		} else {
			insertOffset = fset.File(file.Pos()).Offset(design.End()) - 2
		}
	case LandscapeViewKind, SystemContextViewKind, ContainerViewKind, ComponentViewKind:
		n, _ := findDSL(file, ViewsKind, "")
		if n != nil {
			insertOffset = fset.File(file.Pos()).Offset(n.Pos()) - 2 // -2 to account for `})`
		} else {
			insertOffset = fset.File(file.Pos()).Offset(design.End()) - 2
			code = concat("Views(func() {", code, "})")
		}
	case ElementStyleKind, RelationshipStyleKind:
		n, _ := findDSL(file, "Styles", "")
		if n != nil {
			insertOffset = fset.File(file.Pos()).Offset(n.End()) - 2 // -2 to account for `})`
		} else {
			vn, _ := findDSL(file, ViewsKind, "")
			if vn != nil {
				insertOffset = fset.File(file.Pos()).Offset(vn.Pos()) - 2
				code = concat("Styles(func() {", code, "})")
			} else {
				insertOffset = fset.File(file.Pos()).Offset(design.End()) - 2
				code = concat("Views(func() {", "Styles(func() {", code, "})", "})")
			}
		}
	default:
		return nil, fmt.Errorf("invalid element kind %s", kind)
	}

	newContent := concat(src[:insertOffset], code, src[insertOffset:])
	return setContent(modelFile, res, newContent)
}

// UpsertRelationship updates the code for the given relationship if there's a
// match otherwise it adds a new relationship. The sourcePath is the path to
// the source element, e.g. "Person1" or
// "SoftwareSystem1/Container1/Component1". The destinationPath is the path to
// the destination element.
func (e *Editor) UpsertRelationship(sourceKind ElementKind, sourcePath, destinationPath, code string) (*gentypes.PackageFile, error) {
	res, parsed, fset, err := e.parseDir()
	if err != nil {
		return nil, err
	}
	var sourceNode ast.Node
	var sourceFile *ast.File
	var existingRel ast.Node
	for _, file := range parsed.Files {
		if sourceNode == nil {
			fmt.Println("findDSL", sourceKind, sourcePath)
			sourceNode, _ = findDSL(file, sourceKind, sourcePath)
			if sourceNode != nil {
				sourceFile = file
				existingRel = findRelationship(sourceNode, destinationPath)
			}
		}
		if sourceNode != nil {
			break
		}
	}
	if sourceNode == nil {
		return nil, fmt.Errorf("failed to find source element %s", sourcePath)
	}

	srcFilename := fset.File(sourceNode.Pos()).Name()
	srcBytes, err := os.ReadFile(srcFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to read DSL file %s: %w", srcFilename, err)
	}
	src := string(srcBytes)
	res.Locator.Filename = filepath.Base(srcFilename)
	tokFile := fset.File(sourceFile.Pos())

	// Case 1: relationship already exists, replace it
	if existingRel != nil {
		newContent := concat(
			src[:tokFile.Offset(existingRel.Pos())],
			code,
			src[tokFile.Offset(existingRel.End()):],
		)
		return setContent(srcFilename, res, newContent)
	}

	// Case 2: relationship does not exist, add it
	relnode, err := parser.ParseExpr(code)
	if err != nil {
		return nil, fmt.Errorf("failed to parse relationship DSL: %w", err)
	}
	sourceSrc := src[tokFile.Offset(sourceNode.Pos()):tokFile.Offset(sourceNode.End())]
	updated, err := copyRelationships(tokFile, relnode, code, sourceSrc)
	if err != nil {
		return nil, err
	}
	newContent := concat(
		src[:tokFile.Offset(sourceNode.Pos())],
		updated,
		src[tokFile.Offset(sourceNode.End()):],
	)
	return setContent(srcFilename, res, newContent)
}

// DeleteElement deletes the element with the given kind and path. The
// elementPath is the path to the element in the DSL file, e.g. "Person1" or
// "SoftwareSystem1/Container1/Component1".
func (e *Editor) DeleteElement(kind ElementKind, path string) (*gentypes.PackageFile, error) {
	res, parsed, fset, err := e.parseDir()
	if err != nil {
		return nil, err
	}
	for _, file := range parsed.Files {
		node, _ := findDSL(file, kind, path)
		if node != nil {
			// Found the element, so delete it.
			tokFile := fset.File(file.Pos())
			filename := tokFile.Name()
			res.Locator.Filename = filepath.Base(filename)
			srcBytes, err := os.ReadFile(filename)
			if err != nil {
				return nil, fmt.Errorf("failed to read DSL file %s: %w", tokFile.Name(), err)
			}
			src := string(srcBytes)
			newContent := concat(
				src[:tokFile.Offset(node.Pos())],
				src[tokFile.Offset(node.End()):],
			)
			return setContent(filename, res, newContent)
		}
	}
	return nil, fmt.Errorf("failed to find %s %s", kind, path)
}

// DeleteRelationship deletes the relationship between the given source and
// destination elements.
func (e *Editor) DeleteRelationship(sourceKind ElementKind, sourcePath, destinationPath string) (*gentypes.PackageFile, error) {
	res, parsed, fset, err := e.parseDir()
	if err != nil {
		return nil, err
	}
	var sourceNode ast.Node
	var sourceFile *ast.File
	var existingRel ast.Node
	for _, file := range parsed.Files {
		sourceNode, _ = findDSL(file, sourceKind, sourcePath)
		if sourceNode != nil {
			sourceFile = file
			existingRel = findRelationship(sourceNode, destinationPath)
			break
		}
	}
	if existingRel == nil {
		return nil, fmt.Errorf("failed to find relationship between %s and %s", sourcePath, destinationPath)
	}

	srcFilename := fset.File(sourceNode.Pos()).Name()
	srcBytes, err := os.ReadFile(srcFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to read DSL file %s: %w", srcFilename, err)
	}
	src := string(srcBytes)
	res.Locator.Filename = filepath.Base(srcFilename)
	tokFile := fset.File(sourceFile.Pos())

	newContent := concat(
		src[:tokFile.Offset(existingRel.Pos())],
		src[tokFile.Offset(existingRel.End()):],
	)
	return setContent(srcFilename, res, newContent)
}

// findDSL searches the AST for a CallExpr node representing a DSL with the
// specified kind, path, and arguments.  If the desired DSL node is found, it
// returns the node alongside its parent (a function containing an anonymous
// function argument that calls the found node).  If only the parent DSL
// function is located without the actual DSL node, it returns (nil,
// parentNode).
func findDSL(file *ast.File, kind ElementKind, path string, args ...string) (node, parent ast.Node) {
	var stack []string
	lookupKind := kind
	lookupArgs := args
	if path != "" {
		stack = strings.Split(path, "/")
		if len(stack) > 1 {
			lookupKind = SoftwareSystemKind
			if len(stack) == 2 {
				if kind == ComponentKind {
					lookupKind = ContainerKind
				}
			}
		}
		lookupArgs = append([]string{stack[0]}, args...)
	}
	ast.Inspect(file, func(n ast.Node) bool {
		if node != nil {
			return false
		}
		if checkNode(n, lookupKind, lookupArgs...) {
			if len(stack) <= 1 {
				node = n
			} else {
				parent = n
				stack = stack[1:]
				lookupArgs = append([]string{stack[0]}, args...)
				if len(stack) == 2 {
					lookupKind = ContainerKind
				} else {
					lookupKind = kind
				}
			}
		}
		return node == nil
	})
	return
}

// checkNode checks if the given node is a CallExpr with the specified
// ElementKind and arguments.  It returns true if the node matches the criteria,
// false otherwise. Only non nil and non-empty arguments are matched.
func checkNode(node ast.Node, kind ElementKind, args ...string) bool {
	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return false
	}
	ident, ok := callExpr.Fun.(*ast.Ident)
	if !ok {
		return false
	}
	if ident.Name != string(kind) {
		return false
	}
	if len(callExpr.Args) < len(args) {
		return false
	}
	for i, arg := range args {
		if arg == "" {
			continue
		}
		val, ok := callExpr.Args[i].(*ast.BasicLit)
		if !ok {
			return false
		}
		if fmt.Sprintf("%q", arg) != val.Value {
			return false
		}
	}
	return true
}

// parseDir parses the DSL in the editor directory and returns the parsed AST.
func (e *Editor) parseDir() (*gentypes.PackageFile, *ast.Package, *token.FileSet, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, filepath.Join(e.repo, e.dir), nil, 0)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to parse DSL in %s: %w", e.dir, err)
	}
	res := gentypes.PackageFile{
		Locator: &gentypes.FileLocator{
			Repository: e.repo,
			Dir:        e.dir,
		},
	}
	var parsed *ast.Package
	for _, pkg := range pkgs {
		parsed = pkg
	}
	return &res, parsed, fset, nil
}

// findRelationship looks for a relationship DSL with the given destination
// path.
func findRelationship(node ast.Node, destPath string) ast.Node {
	var relNode ast.Node
	var currentSystem, currentContainer string
	ast.Inspect(node, func(node ast.Node) bool {
		if relNode != nil {
			return false
		}
		if callExpr, ok := node.(*ast.CallExpr); ok {
			ident, ok := callExpr.Fun.(*ast.Ident)
			if !ok {
				return true
			}
			if len(callExpr.Args) == 0 {
				return true
			}
			arg, ok := callExpr.Args[0].(*ast.BasicLit)
			if !ok {
				return true
			}
			switch ident.Name {
			case "SoftwareSystem":
				currentSystem = arg.Value
				return true
			case "Container":
				currentContainer = arg.Value
				return true
			case "Delivers", "Uses", "InteractsWith":
				if pathMatches(arg.Value, destPath, currentSystem, currentContainer) {
					relNode = node
					return false
				}
				relNode = node
				return false
			}
		}
		return true
	})
	return relNode
}

// pathMatches returns true if the given paths match relative to the given
// current system and container (if any)
func pathMatches(path, other, currentSystem, currentContainer string) bool {
	if path == other {
		return true
	}
	if currentSystem != "" {
		if path == fmt.Sprintf("%s/%s", currentSystem, other) {
			return true
		}
		if fmt.Sprintf("%s/%s", currentSystem, path) == other {
			return true
		}
	}
	if currentContainer != "" {
		if path == fmt.Sprintf("%s/%s", currentContainer, other) {
			return true
		}
		if fmt.Sprintf("%s/%s", currentContainer, path) == other {
			return true
		}
		if path == fmt.Sprintf("%s/%s/%s", currentSystem, currentContainer, other) {
			return true
		}
		if fmt.Sprintf("%s/%s/%s", currentSystem, currentContainer, path) == other {
			return true
		}
	}
	return false
}

// copyRelationships copies any relationship DSL that is defined in the given
// node to the body of the first anonymous function defined in code.  It copies
// only the top level relationships, i.e. those that are not defined inside a
// child element.
func copyRelationships(tokFile *token.File, node ast.Node, src, code string) (string, error) {
	// 1. Extract relationship DSL from src
	var rels []string
	var stack []*ast.CallExpr
	ast.Inspect(node, func(node ast.Node) bool {
		if callExpr, ok := node.(*ast.CallExpr); ok {
			for len(stack) > 0 && !isInside(node.Pos(), stack[len(stack)-1].Pos(), stack[len(stack)-1].End()) {
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, callExpr)
			if len(stack) > 2 {
				return false
			}
			if ident, ok := callExpr.Fun.(*ast.Ident); ok {
				if ident.Name == "Uses" || ident.Name == "Delivers" || ident.Name == "InteractsWith" {
					start := tokFile.Offset(callExpr.Pos())
					end := tokFile.Offset(callExpr.End())
					rels = append(rels, src[start:end])
				}
			}
		}
		return true
	})
	if len(rels) == 0 {
		return code, nil
	}
	relCode := strings.Join(rels, "\n")

	// 2. Find the first anonymous function in code and insert the relationship DSL
	node, err := parser.ParseExpr(code)
	if err != nil {
		return "", fmt.Errorf("failed to parse DSL code: %w", err)
	}
	done := false
	ast.Inspect(node, func(node ast.Node) bool {
		if done {
			return false
		}
		if callExpr, ok := node.(*ast.CallExpr); ok {
			if len(callExpr.Args) > 0 {
				lastArg := callExpr.Args[len(callExpr.Args)-1]
				if funcLit, ok := lastArg.(*ast.FuncLit); ok {
					end := tokFile.Offset(funcLit.Body.Rbrace)
					code = concat(code[:end], relCode, code[end:])
					done = true
					return false
				}
			}
			// Add anonymous function with relationship DSL
			start := tokFile.Offset(callExpr.End())
			code = concat(code[:start-1]+", func() {", relCode, "\n}"+code[start-1:])
			done = true
		}
		return done
	})
	return code, nil
}

// isInside returns true if target is inside the range [start, end].
func isInside(target, start, end token.Pos) bool {
	return target >= start && target <= end
}

// setContent replaces the content of the DSL file with newContent and returns
// the updated file. It first formats the new content using the Go formatter.
func setContent(filename string, res *gentypes.PackageFile, newContent string) (*gentypes.PackageFile, error) {
	contentBytes, err := format.Source([]byte(newContent))
	if err != nil {
		return nil, fmt.Errorf("failed to format DSL file %s: %w, original content:\n%s", filename, err, newContent)
	}
	if err := os.WriteFile(filename, contentBytes, 0644); err != nil {
		return nil, fmt.Errorf("failed to write DSL file %s: %w", filename, err)
	}
	res.Content = string(contentBytes)
	return res, nil
}

// newDesign returns the DSL code for a new design with the given code.
func newDesign(code string) string {
	return fmt.Sprintf(`package model

import . "goa.design/model/dsl"

var _ = Design(func() {
%s
})`, code)
}

// concat returns the concatenation of the given strings ensuring one and
// exactly one newline between each.
func concat(s ...string) string {
	var trimmed []string
	for i := 0; i < len(s); i++ {
		s[i] = strings.TrimSpace(s[i])
		if len(s[i]) > 0 {
			trimmed = append(trimmed, s[i])
		}
	}
	return strings.Join(trimmed, "\n")
}
