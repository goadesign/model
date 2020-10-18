package main

//go:generate esc -o webapp.go -pkg main -prefix webapp/dist webapp/dist/

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"goa.design/goa/v3/codegen"
	"goa.design/model/mdl"
	model "goa.design/model/pkg"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func main() {
	out := flag.String("out", codegen.Gendir, "set diagram persistence directory")
	port := flag.Int("port", 8080, "set local HTTP port")
	debug := flag.Bool("debug", false, "pring debug output")
	help := flag.Bool("help", false, "print this help and exit")
	version := flag.Bool("version", false, "print version information and exit")
	devmode := os.Getenv("DEVMODE") == "1"
	flag.Parse()

	if *help {
		showUsage()
		os.Exit(0)
	}

	if *version {
		fmt.Printf("%s %s\n", os.Args[0], model.Version())
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		fail(`missing PACKAGE argument, use "--help" for usage.`)
	}

	pkg := os.Args[1]
	if strings.HasPrefix(pkg, "-") {
		fail(`missing PACKAGE argument, use "--help" for usage.`)
	}

	// Initialize saved views directory if needed.
	outDir, _ := filepath.Abs(*out)
	err := os.MkdirAll(outDir, 0777)
	if err != nil {
		fail(err.Error())
	}

	// Retrieve initial design and create server.
	b, err := gen(pkg, *debug)
	if err != nil {
		fail(err.Error())
	}
	var design mdl.Design
	if err := json.Unmarshal(b, &design); err != nil {
		fail("failed to load design: " + err.Error())
	}
	s := NewServer(&design)

	// Update server whenever design changes on disk.
	err = watch(pkg, func() {
		b, err := gen(pkg, *debug)
		if err != nil {
			fmt.Println("error parsing DSL:\n" + err.Error())
			return
		}
		if err := json.Unmarshal(b, &design); err != nil {
			fmt.Println("failed to load design: " + err.Error())
			return
		}
		s.SetDesign(&design)
	})
	if err != nil {
		fail(err.Error())
	}

	err = s.Serve(outDir, devmode, *port)
	fail(err.Error())
}

func showUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s PACKAGE [FLAGS].\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "Start a HTTP server that serves a graphical editor for the design described in PACKAGE.")
	fmt.Fprintln(os.Stderr, "PACKAGE must be the import path to a Go package containing Model DSL.\n")
	flag.PrintDefaults()
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
