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
	"golang.org/x/tools/go/packages"
)

func main() {
	var (
		out    = flag.String("out", "model.json", "Write structurizr JSON to given file path [use with get or gen].")
		wid    = flag.String("workspace", "", "Structurizr workspace ID [ignored for gen]")
		key    = flag.String("key", "", "Structurizr API key [ignored for gen]")
		secret = flag.String("secret", "", "Structurizr API secret [ignored for gen]")
	)
	flag.Parse()

	var (
		cmd  string
		path string
	)
	inFlag := false
	for _, arg := range os.Args[1:] {
		if !strings.HasPrefix(arg, "-") && !inFlag {
			switch cmd {
			case "":
				cmd = arg
			case "gen", "put":
				path = arg
			default:
				fail("too many arguments, use 'help' for usage.")
			}
		} else {
			inFlag = strings.HasPrefix(arg, "-")
		}
	}

	pathOrDefault := func(p string) string {
		if p == "" {
			return "model.json"
		}
		return p
	}

	var err error
	switch cmd {
	case "gen":
		err = gen(pathOrDefault(path), *out)
	case "get":
		err = get(pathOrDefault(*out), *wid, *key, *secret)
	case "put":
		err = put(pathOrDefault(path), *wid, *key, *secret)
	case "lock":
		err = lock(*wid, *key, *secret)
	case "unlock":
		err = unlock(*wid, *key, *secret)
	case "version":
		fmt.Printf("%s version %s\n", os.Args[0], structurizr.Version())
	case "help":
		showUsage()
	default:
		showUsage()
		os.Exit(1)
	}
	if err != nil {
		fail(err.Error())
	}
}

func gen(pkg, out string) error {
	// Validate package import path
	if _, err := packages.Load(&packages.Config{Mode: packages.NeedName}, pkg); err != nil {
		return err
	}

	// Write program that generates JSON
	tmpDir, err := ioutil.TempDir("", "stz")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)
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
	if err := runCmd(gobin, tmpDir, gobin, "mod", "init", "stz"); err != nil {
		return err
	}
	if err := runCmd(gobin, tmpDir, gobin, "build", "-o", "stz"); err != nil {
		return err
	}

	// Run program
	out, _ = filepath.Abs(out)
	return runCmd(filepath.Join(tmpDir, "stz"), tmpDir, "-out", out)
}

func get(out, wid, key, secret string) error {
	c := service.NewClient(key, secret)
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

func put(path, wid, key, secret string) error {
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

func lock(wid, key, secret string) error {
	c := service.NewClient(key, secret)
	_, err := c.Lock(wid)
	return err
}

func unlock(wid, key, secret string) error {
	c := service.NewClient(key, secret)
	_, err := c.Unlock(wid)
	return err
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func showUsage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintf(os.Stderr, "\n%s gen PACKAGE [FLAGS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s get [FLAGS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s put FILE FLAGS\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s lock [FLAGS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s unlock [FLAGS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s help\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s version\n\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "Where:")
	fmt.Fprintln(os.Stderr, "PACKAGE is the import path to a Go package containing the DSL describing a Structurizr workspace.")
	fmt.Fprintln(os.Stderr, "FILE is the path to a file containing a valid JSON representation of a Structurizr workspace.")
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
	out := os.Args[1]
		
    // Run the model DSL
    w, err := eval.RunDSL()
    if err != nil {
        fmt.Fprintf(os.Stderr, "invalid model: %s", err.Error())
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
