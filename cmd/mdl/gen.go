package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"goa.design/goa/v3/codegen"
	"golang.org/x/tools/go/packages"
)

// tmpDirPrefix is the prefix used to create the temporary directory where the
// code generation tool source code is generated and compiled.
const tmpDirPrefix = "mdl--"

func gen(pkg string, out string, debug bool) error {
	// Validate package import path
	if _, err := packages.Load(&packages.Config{Mode: packages.NeedName}, pkg); err != nil {
		return err
	}

	// Write program that generates Mermaid
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	tmpDir, err := ioutil.TempDir(cwd, tmpDirPrefix)
	if err != nil {
		return err
	}
	defer func() {
		os.RemoveAll(tmpDir)
	}()
	var sections []*codegen.SectionTemplate
	{
		imports := []*codegen.ImportSpec{
			codegen.SimpleImport("encoding/json"),
			codegen.SimpleImport("fmt"),
			codegen.SimpleImport("io/ioutil"),
			codegen.SimpleImport("os"),
			codegen.SimpleImport("path/filepath"),
			codegen.SimpleImport("goa.design/goa/v3/eval"),
			codegen.SimpleImport("goa.design/model/expr"),
			codegen.SimpleImport("goa.design/model/mdl"),
			codegen.NewImport("_", pkg),
		}
		sections = []*codegen.SectionTemplate{
			codegen.Header("Code Generator", "main", imports),
			{Name: "main", Source: mainT},
		}
	}
	cf := &codegen.File{Path: "main.go", SectionTemplates: sections}
	if _, err := cf.Render(tmpDir); err != nil {
		return err
	}

	// Compile program
	gobin, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf(`failed to find a go compiler, looked in "%s"`, os.Getenv("PATH"))
	}
	if _, err := runCmd(gobin, tmpDir, gobin, "build", "-o", "mdl"); err != nil {
		return err
	}

	// Run program
	out, _ = filepath.Abs(out)
	o, err := runCmd(filepath.Join(tmpDir, "mdl"), tmpDir, "-out", out)
	if err != nil && len(o) > 0 {
		err = fmt.Errorf("%s, output:\n%s", err.Error(), o)
	}
	if debug {
		fmt.Fprintln(os.Stderr, o)
	}

	return err
}

// mainT is the template for the generator main package.
const mainT = `func main() {
	// Retrieve output path
	out := os.Args[1]

	// (Re)Create directory
	os.RemoveAll(out)
	if err := os.MkdirAll(out, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create %q: %s", out, err.Error())
		os.Exit(1)
	}
		
	// Run the model DSL
	if err := eval.RunDSL(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	// Render the views and serialize them
	views := mdl.Render(expr.Root)
	for _, view := range views {
		path := filepath.Join(out, view.Key+".json")
		js, _ := json.MarshalIndent(view, "", "    ")
		if err := ioutil.WriteFile(path, js, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write %q: %s", view.Key+".json", err.Error())
			os.Exit(1)
		}
	}
}
`
