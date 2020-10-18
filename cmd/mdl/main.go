package main

//go:generate esc -o webapp.go -pkg main -prefix webapp/dist webapp/dist/

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
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
	var (
		fs      = flag.NewFlagSet("flags", flag.ContinueOnError)
		out     = fs.String("out", codegen.Gendir, "set output directory used to save model JSON and SVGs")
		port    = fs.Int("port", 8080, "set local HTTP port used to serve diagram editor")
		debug   = fs.Bool("debug", false, "pring debug output")
		devmode = os.Getenv("DEVMODE") == "1"
	)

	var (
		cmd   string
		fpath string
		idx   int
	)
	for _, arg := range os.Args[1:] {
		idx++
		switch cmd {
		case "":
			cmd = arg
		case "serve", "gen":
			if !strings.HasPrefix(arg, "-") {
				fpath = arg
				idx++
			}
			goto done
		default:
			goto done
		}
	}
done:
	fs.Parse(os.Args[idx:])

	if cmd == "gen" || cmd == "serve" {
		// Initialize saved views directory if needed.
		*out, _ = filepath.Abs(*out)
		err := os.MkdirAll(*out, 0777)
		if err != nil {
			fail(err.Error())
		}
	}

	var err error
	switch cmd {
	case "gen":
		if fpath == "" {
			fail(`missing PACKAGE argument, use "--help" for usage.`)
		}
		var b []byte
		b, err = gen(fpath, *debug)
		if err == nil {
			err = ioutil.WriteFile(path.Join(*out, "model.json"), b, 0644)
		}
	case "serve":
		err = serve(*out, fpath, *port, devmode, *debug)
	case "version":
		fmt.Printf("%s %s\n", os.Args[0], model.Version())
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

func showUsage(fs *flag.FlagSet) {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintf(os.Stderr, " - %s serve PACKAGE [FLAGS].\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "   Start a HTTP server that serves a graphical editor for the design described in PACKAGE.\n")
	fmt.Fprintf(os.Stderr, " - %s gen PACKAGE [FLAGS].\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "   Generate a JSON representation of the design described in PACKAGE.\n\n")
	fmt.Fprintf(os.Stderr, "   PACKAGE must be the import path to a Go package containing Model DSL.\n\n")
	fs.PrintDefaults()
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
