package main

//go:generate esc -o webapp.go -pkg main -prefix webapp/dist webapp/dist/

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"goa.design/goa/v3/codegen"
	"goa.design/model/mdl"
	model "goa.design/model/pkg"
)

func main() {
	var (
		gset           = flag.NewFlagSet("global", flag.ExitOnError)
		debug, help, h *bool

		genset = flag.NewFlagSet("gen", flag.ExitOnError)
		out    = genset.String("out", "design.json", "set path to generated JSON representation")

		svrset = flag.NewFlagSet("serve", flag.ExitOnError)
		dir    = svrset.String("dir", codegen.Gendir, "set output directory used by editor to save SVGs")
		port   = svrset.Int("port", 8080, "set local HTTP port used to serve diagram editor")

		devmode = os.Getenv("DEVMODE") == "1"
	)

	addGlobals := func(set *flag.FlagSet) {
		debug = set.Bool("debug", false, "print debug output")
		help = set.Bool("help", false, "print this information")
		h = set.Bool("h", false, "print this information")
	}

	var (
		cmd string
		pkg string
		idx = 1
	)
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
			break
		} else if cmd == "" {
			cmd = arg
		} else if pkg == "" {
			pkg = arg
		} else {
			addGlobals(gset)
			showUsage(cmd, genset, svrset, gset)
		}
		idx++
	}

	switch cmd {
	case "gen":
		addGlobals(genset)
		genset.Parse(os.Args[idx:])
		if *h || *help {
			showUsage(cmd, genset)
			os.Exit(0)
		}
	case "serve":
		addGlobals(svrset)
		svrset.Parse(os.Args[idx:])
		if *h || *help {
			showUsage(cmd, svrset)
			os.Exit(0)
		}
	default:
		addGlobals(gset)
		gset.Parse(os.Args[idx:])
		if *h || *help {
			showUsage(cmd, genset, svrset, gset)
			os.Exit(0)
		}
	}

	var err error
	switch cmd {
	case "gen":
		if pkg == "" {
			fail(`missing PACKAGE argument, use "--help" for usage`)
		}
		var b []byte
		b, err = gen(pkg, *debug)
		if err == nil {
			err = ioutil.WriteFile(*out, b, 0644)
		}
	case "serve":
		if pkg == "" {
			fail(`missing PACKAGE argument, use "--help" for usage`)
		}
		*dir, _ = filepath.Abs(*dir)
		if err := os.MkdirAll(*dir, 0777); err != nil {
			fail(err.Error())
		}
		err = serve(*dir, pkg, *port, devmode, *debug)
	case "version":
		fmt.Printf("%s %s\n", os.Args[0], model.Version())
	case "", "help":
		showUsage(cmd, genset, svrset, gset)
	default:
		fail(`unknown command %q, use "--help" for usage`, cmd)
	}
	if err != nil {
		fail(err.Error())
	}
}

func serve(out, pkg string, port int, devmode, debug bool) error {
	// Retrieve initial design and create server.
	b, err := gen(pkg, debug)
	if err != nil {
		return err
	}
	var design mdl.Design
	if err := json.Unmarshal(b, &design); err != nil {
		return fmt.Errorf("failed to load design: %s", err.Error())
	}
	s := NewServer(&design)

	// Update server whenever design changes on disk.
	err = watch(pkg, func() {
		b, err := gen(pkg, debug)
		if err != nil {
			fmt.Println("error parsing DSL:\n" + err.Error())
			return
		}
		var design mdl.Design
		if err := json.Unmarshal(b, &design); err != nil {
			fmt.Println("failed to load design: " + err.Error())
			return
		}
		s.SetDesign(&design)
	})
	if err != nil {
		return err
	}

	return s.Serve(out, devmode, port)
}

func showUsage(cmd string, fss ...*flag.FlagSet) {
	fmt.Fprintln(os.Stderr, "Usage:")
	if cmd != "gen" {
		fmt.Fprintf(os.Stderr, "  %s serve PACKAGE [FLAGS].\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "    Start a HTTP server that serves a graphical editor for the design described in PACKAGE.\n")
	}
	if cmd != "serve" {
		fmt.Fprintf(os.Stderr, "  %s gen PACKAGE [FLAGS].\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "    Generate a JSON representation of the design described in PACKAGE.\n")
	}
	fmt.Fprintf(os.Stderr, "\nPACKAGE must be the import path to a Go package containing Model DSL.\n\n")
	fmt.Fprintf(os.Stderr, "FLAGS:\n")
	for _, fs := range fss {
		fs.PrintDefaults()
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
