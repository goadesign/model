package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"goa.design/goa/v3/codegen"
	"goa.design/structurizr/expr"
	structurizr "goa.design/structurizr/pkg"
	"goa.design/structurizr/service"
)

func main() {
	var (
		out    = flag.String("out", "", "Write structurizr JSON to given file path instead of uploading.")
		wid    = flag.String("wid", "", "Structurizr workspace ID [ignored if -out is used]")
		key    = flag.String("key", "", "Structurizr API key [ignored if -out is used]")
		secret = flag.String("secret", "", "Structurizr API secret [ignored if -out is used]")
	)
	flag.Parse()

	var path string
	inFlag := false
	for _, arg := range os.Args[1:] {
		if !strings.HasPrefix(arg, "-") && !inFlag {
			if path != "" {
				fail("only one argument can be provided.")
			}
			path = arg
		} else {
			inFlag = strings.HasPrefix(arg, "-")
		}
	}

	switch path {
	case "":
		showUsage()
		os.Exit(1)
	case "version":
		fmt.Printf("%s version %s\n", os.Args[0], structurizr.Version())
		os.Exit(0)
	case "help":
		showUsage()
		os.Exit(0)
	default:
		if isFilePath(path) {
			if err := upload(path, *wid, *key, *secret); err != nil {
				fail("upload failed: %s", err.Error())
			}
			os.Exit(0)
		}
		var (
			up     = *out != ""
			output = *out
		)
		if output == "" {
			tdir, err := ioutil.TempDir(".", "model")
			if err != nil {
				fail("failed to create temp dir: %s", err.Error())
			}
			defer os.RemoveAll(tdir)
			output = filepath.Join(tdir, "model.json")
		}
		if err := generate(path, output); err != nil {
			fail(err.Error())
		}
		if up {
			if err := upload(output, *wid, *key, *secret); err != nil {
				fail("upload failed: %s", err.Error())
			}
		}
	}
}

func isFilePath(f string) bool {
	s, err := os.Stat(f)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

func upload(path, wid, key, secret string) error {
	c := service.NewClient(key, secret)
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	var w expr.Workspace
	if err := json.NewDecoder(f).Decode(&w); err != nil {
		return fmt.Errorf("failed to read %q: %s", path, err.Error())
	}
	return c.Put(wid, &w)
}

func generate(pkg, out string) error {
	tmpDir, err := ioutil.TempDir("", "stz")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	// Write program that generates JSON
	var sections []*codegen.SectionTemplate
	{
		imports := []*codegen.ImportSpec{
			codegen.SimpleImport("flag"),
			codegen.SimpleImport("fmt"),
			codegen.SimpleImport("io/ioutil"),
			codegen.SimpleImport("encoding/json"),
			codegen.SimpleImport("os"),
			codegen.SimpleImport("goa.design/structurizr/eval"),
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
	if err := runCmd(gobin, tmpDir, gobin, "get", "-v", "goa.design/structurizr/eval"); err != nil {
		return err
	}
	if err := runCmd(gobin, tmpDir, gobin, "build", "-o", "stz"); err != nil {
		return err
	}

	// Run program
	return runCmd(filepath.Join(tmpDir, "stz"), "-o", out)
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func showUsage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintf(os.Stderr, "%s PACKAGE [FLAGS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s FILE FLAGS\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s help\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s version\n\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "Where:")
	fmt.Fprintln(os.Stderr, "PACKAGE is the import path to a Go package containing the DSL describing a Structurizr workspace.")
	fmt.Fprintln(os.Stderr, "FILE is the path to a file containing a valid JSON representation of a Structurizr workspace")
	fmt.Fprintln(os.Stderr, "FLAGS is a sequence of:")
	flag.PrintDefaults()
}

func runCmd(path, dir string, args ...string) error {
	os.Setenv("GO111MODULE", "on")
	c := exec.Cmd{Path: path, Args: args, Dir: dir}
	b, err := c.CombinedOutput()
	if err != nil {
		if len(b) > 0 {
			return fmt.Errorf(string(b))
		}
		return fmt.Errorf("failed to run command %q in directory %q: %s", path, dir, err)
	}
	return nil
}

// mainT is the template for the generator main.
const mainT = `func main() {
	// Retrieve output path
	out = flag.String("out", "", "")
	flag.Parse()
		
    // Run the model DSL
    w, err := eval.RunDSL()
    if err != nil {
        fmt.Fprintf(os.Stderr, "invalid model: %s", err.String())
        os.Exit(1)
    }
	b, err := json.MarshalIndent("", "    ", w)
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to encode JSON: %s", err.String())
        os.Exit(1)
	}
	if err := ioutil.WriteFile(*out, b, 0644); err != nil {
        fmt.Fprintf(os.Stderr, "failed to write file: %s", err.String())
        os.Exit(1)
	}
}
`
