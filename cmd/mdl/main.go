package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	goacodegen "goa.design/goa/v3/codegen"

	"goa.design/model/codegen"
	"goa.design/model/mdl"
	model "goa.design/model/pkg"
)

type config struct {
	debug   bool
	help    bool
	out     string
	dir     string
	port    int
	devmode bool
	devdist string
}

func main() {
	cfg := parseArgs()

	if cfg.help {
		printUsage()
		os.Exit(0)
	}

	cmd, pkg := parseCommand()

	var err error
	switch cmd {
	case "gen":
		err = generateJSON(pkg, cfg)
	case "serve":
		err = startServer(pkg, cfg)
	case "version":
		fmt.Printf("%s %s\n", os.Args[0], model.Version())
	case "", "help":
		printUsage()
	default:
		fail(`unknown command %q, use "--help" for usage`, cmd)
	}

	if err != nil {
		fail(err.Error())
	}
}

func parseArgs() config {
	cfg := config{
		out:     "design.json",
		dir:     goacodegen.Gendir,
		port:    8080,
		devmode: os.Getenv("DEVMODE") == "1",
		devdist: os.Getenv("DEVDIST"),
	}

	flag.BoolVar(&cfg.debug, "debug", false, "print debug output")
	flag.BoolVar(&cfg.help, "help", false, "print this information")
	flag.BoolVar(&cfg.help, "h", false, "print this information")
	flag.StringVar(&cfg.out, "out", cfg.out, "set path to generated JSON representation")
	flag.StringVar(&cfg.dir, "dir", cfg.dir, "set output directory used by editor to save SVG files")
	flag.IntVar(&cfg.port, "port", cfg.port, "set local HTTP port used to serve diagram editor")

	// Parse only the flags, not the command and package
	args := os.Args[1:]
	flagStart := findFlagStart(args)
	if flagStart > 0 {
		if err := flag.CommandLine.Parse(args[flagStart:]); err != nil {
			fail("failed to parse flags: %s", err.Error())
		}
	}

	return cfg
}

func parseCommand() (string, string) {
	args := os.Args[1:]
	var cmd, pkg string

	for i, arg := range args {
		if strings.HasPrefix(arg, "-") {
			break
		}
		switch i {
		case 0:
			cmd = arg
		case 1:
			pkg = arg
		default:
			printUsage()
			os.Exit(1)
		}
	}

	return cmd, pkg
}

func findFlagStart(args []string) int {
	for i, arg := range args {
		if strings.HasPrefix(arg, "-") {
			return i
		}
	}
	return len(args)
}

func generateJSON(pkg string, cfg config) error {
	if pkg == "" {
		return fmt.Errorf(`missing PACKAGE argument, use "--help" for usage`)
	}

	b, err := codegen.JSON(pkg, cfg.debug)
	if err != nil {
		return err
	}

	return os.WriteFile(cfg.out, b, 0600)
}

func startServer(pkg string, cfg config) error {
	if pkg == "" {
		return fmt.Errorf(`missing PACKAGE argument, use "--help" for usage`)
	}

	absDir, err := filepath.Abs(cfg.dir)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(absDir, 0700); err != nil {
		return err
	}

	if cfg.devmode && cfg.devdist == "" {
		cfg.devdist = "./cmd/mdl/webapp/dist"
	}

	return serve(absDir, pkg, cfg.port, cfg.devdist, cfg.debug)
}

func serve(out, pkg string, port int, devdist string, debug bool) error {
	// Load initial design
	design, err := loadDesign(pkg, debug)
	if err != nil {
		return err
	}

	server := NewServer(design)

	// Watch for changes and update server
	if err := watch(pkg, func() {
		if newDesign, err := loadDesign(pkg, debug); err != nil {
			fmt.Println("error parsing DSL:\n" + err.Error())
		} else {
			server.SetDesign(newDesign)
		}
	}); err != nil {
		return err
	}

	return server.Serve(out, devdist, port)
}

func loadDesign(pkg string, debug bool) (*mdl.Design, error) {
	b, err := codegen.JSON(pkg, debug)
	if err != nil {
		return nil, err
	}

	var design mdl.Design
	if err := json.Unmarshal(b, &design); err != nil {
		return nil, fmt.Errorf("failed to load design: %s", err.Error())
	}

	return &design, nil
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintf(os.Stderr, "  %s serve PACKAGE [FLAGS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "    Start a HTTP server that serves a graphical editor for the design described in PACKAGE.\n")
	fmt.Fprintf(os.Stderr, "  %s gen PACKAGE [FLAGS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "    Generate a JSON representation of the design described in PACKAGE.\n")
	fmt.Fprintf(os.Stderr, "\nPACKAGE must be the import path to a Go package containing Model DSL.\n\n")
	fmt.Fprintf(os.Stderr, "FLAGS:\n")
	flag.PrintDefaults()
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
