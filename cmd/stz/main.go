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
	"goa.design/model/expr"
	model "goa.design/model/pkg"
	"goa.design/model/service"
	"golang.org/x/tools/go/packages"
)

func main() {
	var (
		fs     = flag.NewFlagSet("flags", flag.ContinueOnError)
		out    = fs.String("out", "model.json", "Write structurizr JSON to given file path [use with get or gen].")
		wid    = fs.String("id", "", "Structurizr workspace ID [ignored for gen]")
		key    = fs.String("key", "", "Structurizr API key [ignored for gen]")
		secret = fs.String("secret", "", "Structurizr API secret [ignored for gen]")
		debug  = fs.Bool("debug", false, "Print debug information to stderr.")
	)

	var (
		cmd  string
		path string
		idx  int
	)
	for _, arg := range os.Args[1:] {
		idx++
		switch cmd {
		case "":
			cmd = arg
		case "gen", "sync":
			if !strings.HasPrefix(arg, "-") {
				path = arg
				idx++
			}
			goto done
		default:
			goto done
		}
	}
done:
	fs.Parse(os.Args[idx:])

	pathOrDefault := func(p string) string {
		if p == "" {
			return "model.json"
		}
		return p
	}

	var err error
	switch cmd {
	case "gen":
		if path == "" {
			err = fmt.Errorf("missing Go import package path")
			break
		}
		err = gen(path, *out, *debug)
	case "get":
		err = get(pathOrDefault(*out), *wid, *key, *secret, *debug)
	case "sync":
		err = sync(pathOrDefault(path), *wid, *key, *secret, *debug)
	case "version":
		fmt.Printf("%s version %s\n", os.Args[0], model.Version())
	case "help":
		showUsage(fs)
	default:
		showUsage(fs)
		os.Exit(1)
	}
	if err != nil {
		fail(err.Error())
	}
}

func gen(pkg, out string, debug bool) error {
	// Validate package import path
	if _, err := packages.Load(&packages.Config{Mode: packages.NeedName}, pkg); err != nil {
		return err
	}

	// Write program that generates JSON
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	tmpDir, err := ioutil.TempDir(cwd, "stz")
	if err != nil {
		return err
	}
	defer func() {
		if debug {
			fmt.Printf("temp dir: %q\n", tmpDir)
		} else {
			os.RemoveAll(tmpDir)
		}
	}()
	var sections []*codegen.SectionTemplate
	{
		imports := []*codegen.ImportSpec{
			codegen.SimpleImport("flag"),
			codegen.SimpleImport("fmt"),
			codegen.SimpleImport("io/ioutil"),
			codegen.SimpleImport("encoding/json"),
			codegen.SimpleImport("os"),
			codegen.SimpleImport("goa.design/model/eval"),
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
	if _, err := runCmd(gobin, tmpDir, gobin, "build", "-o", "stz"); err != nil {
		return err
	}

	// Run program
	out, _ = filepath.Abs(out)
	o, err := runCmd(filepath.Join(tmpDir, "stz"), tmpDir, "-out", out)
	if debug {
		fmt.Fprintln(os.Stderr, o)
	}
	return err
}

func get(out, wid, key, secret string, debug bool) error {
	c := service.NewClient(key, secret)
	if debug {
		c.EnableDebug()
	}
	w, err := c.Get(wid)
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(w, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(out, b, 0644)
}

func sync(path, wid, key, secret string, debug bool) error {
	// Load local workspace
	var local expr.Workspace
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = json.NewDecoder(f).Decode(&local); err != nil {
		return err
	}

	// Apply local layout if any
	ext := filepath.Ext(path)
	layoutPath := strings.TrimSuffix(path, ext) + ".layout" + ext
	if _, err := os.Stat(layoutPath); err == nil {
		llf, err := os.Open(layoutPath)
		if err != nil {
			return err
		}
		defer llf.Close()
		layout := make(expr.WorkspaceLayout)
		if err := json.NewDecoder(llf).Decode(&layout); err != nil {
			return err
		}
		local.ApplyLayout(layout)
	}

	// Get remote workspace
	c := service.NewClient(key, secret)
	if debug {
		c.EnableDebug()
	}
	remote, err := c.Get(wid)
	if err != nil {
		return err
	}

	// Merge layouts and persist result
	local.MergeLayout(remote)
	b, err := json.MarshalIndent(local.Layout(), "", "   ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(layoutPath, b, 0777); err != nil {
		return err
	}

	// Upload result to Structurizr
	local.Revision = remote.Revision
	return c.Put(wid, &local)
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func showUsage(fs *flag.FlagSet) {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintf(os.Stderr, "\n%s gen PACKAGE [FLAGS]\t# Generate workspace JSON from DSL.\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s get [FLAGS]\t\t# Get remote workspace.\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s sync FILE FLAGS\t# Sync local workspace layout with remote workspace if any and upload result.\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s help\t\t# Print this help message.\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s version\t\t# Print the tool version.\n\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "Where:")
	fmt.Fprintln(os.Stderr, "\nPACKAGE is the import path to a Go package containing the DSL describing a Structurizr workspace.")
	fmt.Fprintln(os.Stderr, "FILE is the path to a file containing a valid JSON representation of a Structurizr workspace.")
	fmt.Fprintln(os.Stderr, "FLAGS is a sequence of:")
	fs.PrintDefaults()
}

func runCmd(path, dir string, args ...string) (string, error) {
	os.Setenv("GO111MODULE", "on")
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
    w, err := eval.RunDSL()
    if err != nil {
        fmt.Fprint(os.Stderr, err.Error())
        os.Exit(1)
    }
	b, err := json.MarshalIndent(w, "", "    ")
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to encode JSON: %s", err.Error())
        os.Exit(1)
	}
	if err := ioutil.WriteFile(out, b, 0644); err != nil {
        fmt.Fprintf(os.Stderr, "failed to write file: %s", err.Error())
        os.Exit(1)
	}
}
`
