/*
Package main implements a simple AST visualizer for Go programs.

Usage:

	astviz [targetDir]

If targetDir is not specified, the current directory is used.
The output is in DOT format and can be piped to the dot command to generate a graph.

Example:

	astviz | dot -Tpng -o ast.png
*/
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strings"
)

type (
	stack struct {
		nodes []ast.Node
	}

	visitor struct {
		stack stack
		src   string
		fset  *token.FileSet
	}
)

func main() {
	var targetDir string
	if len(os.Args) > 1 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			fmt.Println("Usage: astviz [targetDir]")
			os.Exit(0)
		}
		targetDir = os.Args[1]
	} else {
		targetDir = "./"
	}
	// Assuming you have a Go file in the current directory
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, targetDir, nil, parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Print the header
	fmt.Println("digraph G {")

	// Print the nodes and edges in DOT format
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			pos := file.Pos()
			tokFile := fset.File(pos)
			srcBytes, err := os.ReadFile(tokFile.Name())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			src := string(srcBytes)
			v := &visitor{src: src, fset: fset}
			ast.Walk(v, file)
		}
	}

	// Print the footer
	fmt.Println("}")
}

func (s *stack) push(n ast.Node) {
	s.nodes = append(s.nodes, n)
}

func (s *stack) pop() {
	s.nodes = s.nodes[:len(s.nodes)-1]
}

func (s *stack) peek() ast.Node {
	if len(s.nodes) == 0 {
		return nil
	}
	return s.nodes[len(s.nodes)-1]
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		v.stack.pop()
		return v
	}

	fmt.Printf("\"%p\" [label=\"%s\"];\n", node, nodeDescription(node, v.src, v.fset))

	// If there's a parent node, create an edge from the parent to the current node
	if parent := v.stack.peek(); parent != nil {
		fmt.Printf("\"%p\" -> \"%p\";\n", parent, node)
	}

	// Push the current node onto the stack as it's now the current parent
	v.stack.push(node)
	return v
}

func nodeDescription(node ast.Node, src string, fset *token.FileSet) string {
	if node == nil {
		return "nil"
	}

	pos := fset.Position(node.Pos())
	end := fset.Position(node.End())

	// Extract text from the source
	text := string([]rune(src)[pos.Offset:end.Offset])
	// Trim to the first newline
	if idx := strings.Index(text, "\n"); idx != -1 {
		text = text[:idx]
	}
	text = escapeForDotLabel(text)

	// Return a combination of the node type and the corresponding text
	nodeType := strings.TrimPrefix(reflect.TypeOf(node).String(), "*ast.")
	return fmt.Sprintf("%s: %s", nodeType, text)
}

func escapeForDotLabel(s string) string {
	// Escape backslashes first to avoid double escaping
	s = strings.Replace(s, "\\", "\\\\", -1)
	// Escape double quotes
	s = strings.Replace(s, "\"", "\\\"", -1)
	// Escape newlines
	s = strings.Replace(s, "\n", "\\n", -1)
	return s
}
