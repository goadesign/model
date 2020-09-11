package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"goa.design/goa/v3/codegen"
	model "goa.design/model/pkg"
)

func main() {
	var (
		fs     = flag.NewFlagSet("flags", flag.ExitOnError)
		out    = fs.String("out", codegen.Gendir, "Write diagrams to given directory")
		config = fs.String("config", "", "Path to Mermaid config JSON used to serve diagrams")
		port   = fs.Int("port", 6070, "Listen port used to serve diagrams")
		debug  = fs.Bool("debug", false, "Pring debug output")
	)

	var (
		cmd  string
		path string
		idx  int
	)
	for _, arg := range os.Args[1:] {
		switch cmd {
		case "":
			cmd = arg
			idx++
		case "gen", "render", "serve":
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
	fs.Parse(os.Args[idx+1:])

	var err error
	switch cmd {
	case "gen":
		if path == "" {
			fail("missing Go import package path")
		}
		err = gen(path, *out, *debug)
	case "render", "serve":
		if path == "" {
			fail("missing Go import package path")
		}
		var cfg string
		if *config != "" {
			b, err := ioutil.ReadFile(*config)
			if err != nil {
				fail("failed to load config: %s", err)
			}
			cfg = string(b)
		}
		if cmd == "render" {
			err = render(path, cfg, *out, *debug)
		} else {
			err = serve(path, cfg, *out, *port, *debug)
		}
	case "help":
		showUsage(fs)
	case "version":
		fmt.Printf("%s %s\n", os.Args[0], model.Version())
	default:
		showUsage(fs)
		os.Exit(1)
	}
	if err != nil {
		fail(err.Error())
	}
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func showUsage(fs *flag.FlagSet) {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintf(os.Stderr, "%s gen PACKAGE [FLAGS]\t# Generate Mermaid diagrams from DSL.\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s serve PACAKGE [FLAGS]\t# Serve diagrams to HTTP clients.\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s help\t\t\t# Print this help message.\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s version\t\t\t# Print the tool version.\n\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "Where:")
	fmt.Fprintln(os.Stderr, "PACKAGE is the import path to a Go package containing the DSL describing a Structurizr workspace.")
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
