package codegen

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"goa.design/goa/v3/codegen"
	"golang.org/x/tools/go/packages"
)

// TmpDirPrefix is the prefix used to create temporary directories.
const TmpDirPrefix = "mdl--"

// JSON generates a JSON representation of the model described in pkg.  dir is
// the directory in which to run the build system's query tool.  pkg must be a
// valid Go package import path.
func JSON(dir, pkg string, debug bool) ([]byte, error) {
	// Validate package import path
	if _, err := packages.Load(&packages.Config{Dir: dir, Mode: packages.NeedName}, pkg); err != nil {
		return nil, err
	}

	// Write program that generates JSON
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	tmpDir, err := os.MkdirTemp(cwd, TmpDirPrefix)
	if err != nil {
		return nil, err
	}
	defer func() { os.RemoveAll(tmpDir) }()
	var sections []*codegen.SectionTemplate
	{
		imports := []*codegen.ImportSpec{
			codegen.SimpleImport("fmt"),
			codegen.SimpleImport("encoding/json"),
			codegen.SimpleImport("os"),
			codegen.SimpleImport("goa.design/model/model"),
			codegen.NewImport("_", pkg),
		}
		sections = []*codegen.SectionTemplate{
			codegen.Header("Code Generator", "main", imports),
			{Name: "main", Source: mainT},
		}
	}
	cf := &codegen.File{Path: "main.go", SectionTemplates: sections}
	if _, err := cf.Render(tmpDir); err != nil {
		return nil, err
	}

	// Compile program
	gobin, err := exec.LookPath("go")
	if err != nil {
		return nil, fmt.Errorf(`failed to find a go compiler, looked in "%s"`, os.Getenv("PATH"))
	}
	if _, err := runCmd(gobin, tmpDir, "build", "-o", "mdl"); err != nil {
		return nil, err
	}

	// Run program
	o, err := runCmd(path.Join(tmpDir, "mdl"), tmpDir, "model.json")
	if debug {
		fmt.Fprintln(os.Stderr, o)
	}
	if err != nil {
		return nil, err
	}
	return os.ReadFile(path.Join(tmpDir, "model.json"))
}

func runCmd(path, dir string, args ...string) (string, error) {
	args = append([]string{path}, args...) // args[0] becomes exec path
	c := exec.Cmd{Path: path, Args: args, Dir: dir}
	b, err := c.CombinedOutput()
	if err != nil {
		if len(b) > 0 {
			return "", fmt.Errorf(string(b))
		}
		return "", fmt.Errorf("failed to run command %q in directory %q: %s", path, dir, err)
	}
	return string(b), nil
}

// mainT is the template for the generator main.
const mainT = `func main() {
	// Retrieve output path
	out := os.Args[1]
		
	// Run the model DSL
	w, err := mdl.RunDSL()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
	b, err := json.MarshalIndent(w, "", "    ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode into JSON: %s", err.Error())
		os.Exit(1)
	}
	if err := os.WriteFile(out, b, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write file: %s", err.Error())
		os.Exit(1)
	}
}
`
