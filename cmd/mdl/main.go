package main

//go:generate esc -o webapp.go -pkg main -prefix webapp/dist webapp/dist/

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"goa.design/goa/v3/codegen"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func main() {
	pkg := flag.String("pkg", "", "Model package to read")
	out := flag.String("out", codegen.Gendir, "Write diagrams to given directory")
	port := flag.Int("port", 8080, "Local HTTP port.")
	debug := flag.Bool("debug", false, "Pring debug output")
	devmode := os.Getenv("DEVMODE") == "1"
	flag.Parse()

	_, _, _, _ = out, port, debug, devmode

	outDir, _ := filepath.Abs(*out)
	err := os.MkdirAll(outDir, 0777)
	if err != nil {
		fail(err.Error())
	}

	model, err := gen(*pkg, *debug)
	if err != nil {
		fail(err.Error())
	}

	srv := Server{model: model}
	err = srv.Serve(outDir, devmode, *port)
	fail(err.Error())
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
