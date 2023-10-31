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

	// ViewKind is the kind of a DSL view.
	ViewKind string
)

const (
	// PersonKind is the kind of a person.
	PersonKind ElementKind = "Person"
	// SoftwareSystemKind is the kind of a software system.
	SoftwareSystemKind = "SoftwareSystem"
	// ContainerKind is the kind of a container.
	ContainerKind = "Container"
	// ComponentKind is the kind of a component.
	ComponentKind = "Component"
)

const (
	// LandscapeViewKind is the kind of a landscape view.
	LandscapeViewKind ViewKind = "SystemLandscapeView"
	// SystemContextViewKind is the kind of a system context view.
	SystemContextViewKind = "SystemContextView"
	// ContainerViewKind is the kind of a container view.
	ContainerViewKind = "ContainerView"
	// ComponentViewKind is the kind of a component view.
	ComponentViewKind = "ComponentView"
)

// DefaultModelFilename is the name of the DSL file used by default (e.g. when
// adding new DSL).
const DefaultModelFilename = "model.go"

// NewEditor returns a new Editor.
func NewEditor(repo, dir string) *Editor {
	return &Editor{repo: repo, dir: dir}
}

// UpsertElement updates the code for the given element if there's a match
// otherwise it adds a new element. The elementPath is the path to the element
// in the DSL file, e.g. "Person1" or "SoftwareSystem1/Container1/Component1".
func (e *Editor) UpsertElement(kind ElementKind, elementPath, code string) (*gentypes.PackageFile, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, filepath.Join(e.repo, e.dir), nil, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSL in %s: %w", e.dir, err)
	}
	if len(pkgs) > 1 {
		return nil, fmt.Errorf("found %d packages in %s, expected 1", len(pkgs), e.dir)
	}
	res := gentypes.PackageFile{
		Locator: &gentypes.FileLocator{
			Repository: e.repo,
			Dir:        e.dir,
		},
	}
	if len(pkgs) == 0 {
		res.Locator.Filename = DefaultModelFilename
		modelFile := filepath.Join(e.repo, e.dir, DefaultModelFilename)
		return setContent(modelFile, res, newDesign(code))
	}
	var parsed *ast.Package
	for _, pkg := range pkgs {
		parsed = pkg
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
		pathStack := strings.Split(elementPath, "/")
		var pos, end, prevPos, prevEnd token.Pos
		recurse := true
		var elemNode ast.Node
		ast.Inspect(file, func(node ast.Node) bool {
			if elemNode != nil {
				return false
			}
			pos, end, pathStack, recurse = findElement(node, kind, pathStack)
			if pos != token.NoPos {
				prevPos = pos
				prevEnd = end
			}
			if !recurse {
				elemNode = node
			}
			return recurse
		})
		if elemNode != nil {
			// Found the element, so replace it.
			code, err := copyRelationships(tokFile, elemNode, src, code)
			if err != nil {
				return nil, err
			}
			newContent := concat(src[:tokFile.Offset(prevPos)], code, src[tokFile.Offset(prevEnd):])
			return setContent(filename, res, newContent)
		}
		if prevEnd != token.NoPos && pos == token.NoPos {
			// We found the parent, but not the child, so add the child in the parent.
			endOffset := tokFile.Offset(prevEnd)
			newContent := concat(src[:endOffset-2], code, src[endOffset-2:])
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
	ast.Inspect(f, func(node ast.Node) bool {
		if callExpr, ok := node.(*ast.CallExpr); ok {
			if ident, ok := callExpr.Fun.(*ast.Ident); ok {
				if ident.Name == "Views" {
					viewsCallOffset = tokFile.Offset(callExpr.Pos())
					return false
				}
			}
		}
		return true
	})
	if viewsCallOffset == 0 {
		// Find end of Design DSL and insert before ending brackets.
		ast.Inspect(f, func(node ast.Node) bool {
			if callExpr, ok := node.(*ast.CallExpr); ok {
				if ident, ok := callExpr.Fun.(*ast.Ident); ok {
					if ident.Name == "Design" {
						viewsCallOffset = tokFile.Offset(callExpr.End()) - 2
						return false
					}
				}
			}
			return true
		})
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
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, filepath.Join(e.repo, e.dir), nil, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSL in %s: %w", e.dir, err)
	}
	if len(pkgs) != 1 {
		return nil, fmt.Errorf("found %d packages in %s, expected 1", len(pkgs), e.dir)
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
	var sourceNode, destinationNode ast.Node
	var sourceFile *ast.File
	var existingRel ast.Node
	for _, file := range parsed.Files {
		if sourceNode == nil {
			sourceNode = findNode(file, sourcePath)
			if sourceNode != nil {
				sourceFile = file
				existingRel = findRelationship(sourceNode, destinationPath)
			}
		}
		if destinationNode == nil {
			destinationNode = findNode(file, destinationPath)
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
		start := tokFile.Offset(existingRel.Pos())
		end := tokFile.Offset(existingRel.End())
		newContent := concat(src[:start], code, src[end:])
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
	fmt.Println("sourceSrc:", sourceSrc)
	fmt.Println("code:", code)
	fmt.Println("updated:", updated)
	newContent := concat(src[:tokFile.Offset(sourceNode.Pos())], updated, src[tokFile.Offset(sourceNode.End()):])
	fmt.Println("newContent:", newContent)
	return setContent(srcFilename, res, newContent)
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

// findNode looks for a ast Node with the given path.
func findNode(file *ast.File, path string) ast.Node {
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
	recurse := true
	var n ast.Node
	ast.Inspect(file, func(node ast.Node) bool {
		if n != nil {
			return false
		}
		_, _, stack, recurse = findElement(node, elemKind, stack)
		if !recurse {
			n = node
		}
		return recurse
	})
	if n == nil && elemKind == SoftwareSystemKind {
		// Try person
		ast.Inspect(file, func(node ast.Node) bool {
			if n != nil {
				return false
			}
			_, _, stack, recurse = findElement(node, PersonKind, stack)
			if !recurse {
				n = node
			}
			return recurse
		})
	}
	return n
}

// setContent replaces the content of the DSL file with newContent and returns
// the updated file. It first formats the new content using the Go formatter.
func setContent(filename string, res gentypes.PackageFile, newContent string) (*gentypes.PackageFile, error) {
	contentBytes, err := format.Source([]byte(newContent))
	if err != nil {
		return nil, fmt.Errorf("failed to format DSL file %s: %w, original content:\n%s", filename, err, newContent)
	}
	if err := os.WriteFile(filename, contentBytes, 0644); err != nil {
		return nil, fmt.Errorf("failed to write DSL file %s: %w", filename, err)
	}
	res.Content = string(contentBytes)
	return &res, nil
}

// findElement looks for a call expression with a first argument that matches the first element of pathStack.
// If found there are two cases:
//
//  1. The function matches ElementKind and the pathStack has only one element: return the position of the call
//     expression and an empty pathStack.
//  2. The function matches ContainerViewKind or ComponentViewKind and the pathStack has more than one element:
//     return the position of the call expression and the remaining pathStack.
//
// The last returned value is true if the search should continue.
func findElement(node ast.Node, kind ElementKind, pathStack []string) (pos, end token.Pos, stack []string, recurse bool) {
	stack = pathStack
	recurse = true
	if callExpr, ok := node.(*ast.CallExpr); ok {
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
