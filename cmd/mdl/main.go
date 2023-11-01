package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"goa.design/model/codegen"
	model "goa.design/model/pkg"
	"goa.design/model/svc/clients/repo"
)

func main() {
	var (
		gset           = flag.NewFlagSet("global", flag.ExitOnError)
		debug, help, h *bool

		genset = flag.NewFlagSet("gen", flag.ExitOnError)
		out    = genset.String("out", "design.json", "Set path to generated JSON representation")

		svrset = flag.NewFlagSet("serve", flag.ExitOnError)
		port   = svrset.Int("port", 8080, "set local HTTP port used to serve diagram editor")
		root   = svrset.String("root", "", "set root path to package files")

		devmode = os.Getenv("DEVMODE") == "1"

		showUsage = func() { printUsage(svrset, genset, gset) }
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
			showUsage()
		}
		idx++
	}

	switch cmd {
	case "gen":
		addGlobals(genset)
		if err := genset.Parse(os.Args[idx:]); err != nil {
			fail(err.Error())
		}
	case "serve":
		addGlobals(svrset)
		if err := svrset.Parse(os.Args[idx:]); err != nil {
			fail(err.Error())
		}
	default:
		addGlobals(gset)
		if err := gset.Parse(os.Args[idx:]); err != nil {
			fail(err.Error())
		}
	}

	if *h || *help {
		showUsage()
		os.Exit(0)
	}

	var err error
	switch cmd {
	case "gen":
		if pkg == "" {
			fail(`missing PACKAGE argument, use "--help" for usage`)
		}
		var b []byte
		b, err = codegen.JSON("", pkg, *debug)
		if err == nil {
			err = os.WriteFile(*out, b, 0644)
		}
	case "serve":
		store := repo.NewHandler(*root, devmode)
		err = serve(store, *port, devmode, *debug)
	case "version":
		fmt.Printf("%s %s\n", os.Args[0], model.Version())
	case "", "help":
		showUsage()
	default:
		fail(`unknown command %q, use "--help" for usage`, cmd)
	}
	if err != nil {
		fail(err.Error())
	}
}

func printUsage(fss ...*flag.FlagSet) {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintf(os.Stderr, "  %s serve MODEL_DIR [FLAGS].\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "    Start a HTTP server that serves a graphical editor for the models located in MODEL_DIR.\n")
	fmt.Fprintf(os.Stderr, "  %s gen PACKAGE [FLAGS].\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "    Generate a JSON representation of the design described in PACKAGE.\n")
	fmt.Fprintf(os.Stderr, "    PACKAGE must be the import path to a Go package containing Model DSL.\n\n")
	fmt.Fprintf(os.Stderr, "FLAGS:\n")
	for _, fs := range fss {
		fs.PrintDefaults()
	}
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
