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

// UpsertElementByPath updates the code for the person, software system ,
// container or component with the given kind and path, if there's a match
// otherwise it adds a new element. The elementPath is the path to the element
// in the DSL file, e.g. "Person1" or "SoftwareSystem1/Container1/Component1".
func (e *Editor) UpsertElementByPath(kind ElementKind, elementPath, code string) (*gentypes.PackageFile, error) {
	res, parsed, fset, err := e.parseDir()
	if err != nil {
		return nil, err
	}
	if parsed == nil {
		res.Locator.Filename = DefaultModelFilename
		modelFile := filepath.Join(e.repo, e.dir, DefaultModelFilename)
		return setContent(modelFile, res, newDesign(code))
	}
	for _, file := range parsed.Files {
		tokFile := fset.File(file.Pos())
		filename := tokFile.Name()
		res.Locator.Filename = filepath.Base(filename)
		srcBytes, err := os.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to read DSL file %s: %w", tokFile.Name(), err)
		}
		src := string(srcBytes)
		elemNode, pos, end := findDSL(file, kind, elementPath)
		if elemNode != nil {
			// Found the element, so replace it.
			code, err := copyRelationships(tokFile, elemNode, src, code)
			if err != nil {
				return nil, err
			}
			newContent := concat(
				src[:tokFile.Offset(pos)],
				code,
				src[tokFile.Offset(end):],
			)
			return setContent(filename, res, newContent)
		}
		if end != token.NoPos {
			// We found the parent, but not the child, so add the child in the parent.
			endOffset := tokFile.Offset(end)
			newContent := concat(
				src[:endOffset-2], // -2 to account for `})`
				code,
				src[endOffset-2:],
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
		return setContent(modelFile, res, newDesign(code))
	}
	// Find the position of the Views() call if any.
	f, err := parser.ParseFile(fset, modelFile, nil, parser.AllErrors)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSL file %s: %w", modelFile, err)
	}
	tokFile := fset.File(f.Pos())
	var viewsCallOffset int
	n, _, _ := findDSL(f, ViewsKind, "")
	if n != nil {
		viewsCallOffset = tokFile.Offset(n.Pos())
	}
	if viewsCallOffset == 0 {
		_, _, end := findDSL(f, DesignKind, "")
		if end == token.NoPos {
			return nil, fmt.Errorf("failed to find Design() call in DSL file %s", modelFile)
		}
		viewsCallOffset = tokFile.Offset(end) - 2
	}
	data, err := os.ReadFile(modelFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read DSL file %s: %w", modelFile, err)
	}
	content := string(data)
	newContent := concat(content[:viewsCallOffset], code, content[viewsCallOffset:])
	return setContent(modelFile, res, newContent)
}

// UpsertRelationship updates the code for the given relationship if there's a
// match otherwise it adds a new relationship. The sourcePath is the path to
// the source element, e.g. "Person1" or
// "SoftwareSystem1/Container1/Component1". The destinationPath is the path to
// the destination element.
func (e *Editor) UpsertRelationship(sourcePath, destinationPath, code string) (*gentypes.PackageFile, error) {
	res, parsed, fset, err := e.parseDir()
	if err != nil {
		return nil, err
	}
	var sourceNode, destinationNode ast.Node
	var sourceFile *ast.File
	var existingRel ast.Node
	for _, file := range parsed.Files {
		if sourceNode == nil {
			sourceNode = findElement(file, sourcePath)
			if sourceNode != nil {
				sourceFile = file
				existingRel = findRelationship(sourceNode, destinationPath)
			}
		}
		if destinationNode == nil {
			destinationNode = findElement(file, destinationPath)
		}
		if sourceNode != nil && destinationNode != nil {
			break
		}
	}
	if sourceNode == nil {
		return nil, fmt.Errorf("failed to find source element %s", sourcePath)
	}
	if destinationNode == nil {
		return nil, fmt.Errorf("failed to find destination element %s", destinationPath)
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

// UpsertElementByID updates the code for the view, element style or
// relationship style with the given kind and identified with the given ID if
// there's a match otherwise it adds new DSL.  The ID corresponds to the first
// argument of the DSL function, e.g. the ID of a view is the key.
func (e *Editor) UpsertElementByID(kind ElementKind, id, code string) (*gentypes.PackageFile, error) {
	res, parsed, fset, err := e.parseDir()
	if err != nil {
		return nil, err
	}
	var node ast.Node
	var file *ast.File
	for _, f := range parsed.Files {
		node, _, _ = findDSL(file, kind, id)
		if node != nil {
			file = f
			break
		}
	}
	if node != nil {
		// Found the element, so replace it.
		srcFilename := fset.File(node.Pos()).Name()
		srcBytes, err := os.ReadFile(srcFilename)
		if err != nil {
			return nil, fmt.Errorf("failed to read DSL file %s: %w", srcFilename, err)
		}
		src := string(srcBytes)
		res.Locator.Filename = filepath.Base(srcFilename)
		tokFile := fset.File(file.Pos())

		newContent := concat(
			src[:tokFile.Offset(node.Pos())],
			code,
			src[tokFile.Offset(node.End()):],
		)
		return setContent(srcFilename, res, newContent)
	}
	// We didn't find the element, so add it
	modelFile := filepath.Join(e.repo, e.dir, DefaultModelFilename)
	res.Locator.Filename = DefaultModelFilename
	if _, err := os.Stat(modelFile); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("failed to stat DSL file %s: %w", modelFile, err)
		}
		return setContent(modelFile, res, newDesign(code))
	}
	// Where to insert the code depends on the kind:
	//  - Views: insert last in the Views DSL
	//  - Styles: insert last in the Styles DSL
	var insertOffset int
	switch kind {
	case LandscapeViewKind, SystemContextViewKind, ContainerViewKind, ComponentViewKind:
		n, _, _ := findDSL(file, ViewsKind, "")
		if n != nil {
			insertOffset = fset.File(file.Pos()).Offset(n.Pos()) - 2 // -2 to account for `})`
		} else {
			dn, _, _ := findDSL(file, DesignKind, "")
			if dn == nil {
				return nil, fmt.Errorf("failed to find Design() call in DSL file %s", modelFile)
			}
			insertOffset = fset.File(file.Pos()).Offset(dn.End()) - 2
			code = concat("Views(func() {", code, "})")
		}
	case ElementStyleKind, RelationshipStyleKind:
		n, _, _ := findDSL(file, "Styles", "")
		if n != nil {
			insertOffset = fset.File(file.Pos()).Offset(n.End()) - 2 // -2 to account for `})`
		} else {
			vn, _, _ := findDSL(file, ViewsKind, "")
			if vn != nil {
				insertOffset = fset.File(file.Pos()).Offset(vn.Pos()) - 2
				code = concat("Styles(func() {", code, "})")
			} else {
				dn, _, _ := findDSL(file, DesignKind, "")
				if dn == nil {
					return nil, fmt.Errorf("failed to find Design() call in DSL file %s", modelFile)
				}
				insertOffset = fset.File(file.Pos()).Offset(dn.End()) - 2
				code = concat("Views(func() {", "Styles(func() {", code, "})", "})")
			}
		}
	default:
		return nil, fmt.Errorf("invalid element kind %s", kind)
	}
	data, err := os.ReadFile(modelFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read DSL file %s: %w", modelFile, err)
	}
	content := string(data)
	newContent := concat(content[:insertOffset], code, content[insertOffset:])
	return setContent(modelFile, res, newContent)
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

// findElement looks for a software system, container or component with the
// given path.
func findElement(file *ast.File, path string) ast.Node {
	stack := strings.Split(path, "/")
	var elemKind ElementKind
	switch len(stack) {
	case 1:
		elemKind = SoftwareSystemKind
	case 2:
		elemKind = ContainerKind
	case 3:
		elemKind = ComponentKind
	}
	n, _, _ := findDSL(file, elemKind, path)
	if n == nil && elemKind == SoftwareSystemKind {
		// Try person
		n, _, _ = findDSL(file, PersonKind, path)
	}
	return n
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

// findDSL looks for a DSL node (ast.CallExpr) with the given kind, path and
// first argument. If found it returns the node and its position in file. If
// it finds a parent DSL function but not the actual DSL it returns the position
// of the parent DSL function and a nil node.
func findDSL(file *ast.File, kind ElementKind, path string) (node ast.Node, pos, end token.Pos) {
	pathStack := strings.Split(path, "/")
	var curPos, curEnd, prevPos, prevEnd token.Pos
	recurse := true
	ast.Inspect(file, func(n ast.Node) bool {
		if n != nil {
			return false
		}
		curPos, curEnd, pathStack, recurse = parseNode(n, kind, pathStack)
		if curPos != token.NoPos {
			prevPos = curPos
			prevEnd = curEnd
		}
		if !recurse {
			node = n
		}
		return recurse
	})
	pos = prevPos
	end = prevEnd
	return
}

// parseNode checks if node is a call expression with a first argument that
// matches the first element of pathStack. If so there are two cases:
//
//  1. The function matches ElementKind and the pathStack has only one element:
//     return the position of the call expression and an empty pathStack.
//  2. The function matches SoftwareSystemKind or ContainerKind and the
//     pathStack has more than one element: return the position of the call
//     expression and the remaining pathStack so the caller can recurse.
//
// The last returned value is true if the search should continue.
func parseNode(node ast.Node, kind ElementKind, pathStack []string) (pos, end token.Pos, stack []string, recurse bool) {
	stack = pathStack
	recurse = true
	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return
	}
	if callExpr.Args == nil || len(callExpr.Args) == 0 {
		return
	}
	arg, ok := callExpr.Args[0].(*ast.BasicLit)
	if !ok {
		return
	}
	if arg.Value != fmt.Sprintf("%q", pathStack[0]) {
		return
	}
	switch len(pathStack) {
	case 1:
		if ident, ok := callExpr.Fun.(*ast.Ident); ok {
			if ident.Name == string(kind) {
				return callExpr.Pos(), callExpr.End(), nil, false
			}
		}
	case 2:
		// We must be in a system or a container for a match
		if ident, ok := callExpr.Fun.(*ast.Ident); ok {
			if ident.Name == string(SoftwareSystemKind) && kind == ContainerKind || ident.Name == string(ContainerKind) && kind == ComponentKind {
				return callExpr.Pos(), callExpr.End(), pathStack[1:], true
			}
		}
	case 3:
		// We must be in a system for a match
		if ident, ok := callExpr.Fun.(*ast.Ident); ok {
			if ident.Name == string(SoftwareSystemKind) {
				return callExpr.Pos(), callExpr.End(), pathStack[1:], true
			}
		}
	}
	return
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
