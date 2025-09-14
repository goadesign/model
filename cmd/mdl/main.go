package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	goacodegen "goa.design/goa/v3/codegen"

	"goa.design/model/codegen"
	"goa.design/model/mdl"
	model "goa.design/model/pkg"

	"context"

	"github.com/chromedp/chromedp"
)

type config struct {
	debug   bool
	help    bool
	out     string
	dir     string
	port    int
	devmode bool
	devdist string
	// svg command options
	views     SliceFlag
	all       bool
	direction string
	compact   bool
	timeout   time.Duration
}

// SliceFlag implements flag.Value for repeated string flags
type SliceFlag []string

func (s *SliceFlag) String() string { return strings.Join(*s, ",") }
func (s *SliceFlag) Set(v string) error {
	*s = append(*s, v)
	return nil
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
	case "svg":
		err = runSVG(pkg, cfg)
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
		// defaults for svg command
		direction: "DOWN",
		timeout:   20 * time.Second,
	}

	flag.BoolVar(&cfg.debug, "debug", false, "print debug output")
	flag.BoolVar(&cfg.help, "help", false, "print this information")
	flag.BoolVar(&cfg.help, "h", false, "print this information")
	flag.StringVar(&cfg.out, "out", cfg.out, "set path to generated JSON representation")
	flag.StringVar(&cfg.dir, "dir", cfg.dir, "set output directory used by editor to save SVG files")
	flag.IntVar(&cfg.port, "port", cfg.port, "set local HTTP port used to serve diagram editor")
	// svg command flags (safe to always register)
	flag.Var(&cfg.views, "view", "view key to render (repeatable)")
	flag.BoolVar(&cfg.all, "all", false, "render all views")
	flag.StringVar(&cfg.direction, "direction", cfg.direction, "auto-layout direction: DOWN|UP|LEFT|RIGHT")
	flag.BoolVar(&cfg.compact, "compact", false, "enable compact auto-layout")
	flag.DurationVar(&cfg.timeout, "timeout", cfg.timeout, "timeout per view (e.g. 15s)")

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

// runSVG runs a local server without watch, opens a headless browser per view to trigger
// auto layout and save, and waits for the SVG files to be written.
func runSVG(pkg string, cfg config) error {
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

	// Load design to enumerate views
	design, err := loadDesign(pkg, cfg.debug)
	if err != nil {
		return err
	}

	// Collect view keys based on flags
	viewKeys := collectViewKeys(design)
	selected := make([]string, 0)
	if cfg.all || len(cfg.views) == 0 && !cfg.all {
		// default to all if nothing specified
		selected = viewKeys
	}
	if len(cfg.views) > 0 {
		m := make(map[string]bool)
		for _, k := range viewKeys {
			m[k] = true
		}
		for _, v := range cfg.views {
			if !m[v] {
				return fmt.Errorf("unknown view %q; known views: %s", v, strings.Join(viewKeys, ", "))
			}
			selected = append(selected, v)
		}
	}
	if len(selected) == 0 {
		return fmt.Errorf("no views to render; use --all or --view")
	}

	// Prepare server with custom mux (no watch)
	server := NewServer(design)
	mux := http.NewServeMux()

	// Pick port: if cfg.port == 0, choose a random free port
	port := cfg.port
	if port == 0 {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return err
		}
		port = ln.Addr().(*net.TCPAddr).Port
		_ = ln.Close()
	}

	httpServer := &http.Server{
		Addr:              fmt.Sprintf("127.0.0.1:%d", port),
		ReadHeaderTimeout: 3 * time.Second,
	}

	// Start server in background
	done := make(chan error, 1)
	go func() {
		done <- server.ServeOnMux(absDir, cfg.devdist, httpServer, mux)
	}()

	// Wait for server to accept connections by polling
	baseURL := fmt.Sprintf("http://127.0.0.1:%d", port)
	if err := waitForHTTP(baseURL+"/data/model.json", 5*time.Second); err != nil {
		_ = httpServer.Close()
		return fmt.Errorf("server did not start: %w", err)
	}

	// Drive headless browser to render and save
	if err := renderViewsHeadless(baseURL, absDir, selected, cfg); err != nil {
		_ = httpServer.Close()
		return err
	}

	// Shutdown server
	_ = httpServer.Close()
	// Ensure background server goroutine exits to avoid leaks
	select {
	case <-done:
		// server exited
	case <-time.After(2 * time.Second):
		// timeout waiting for server to exit; continue
	}
	return nil
}

func collectViewKeys(d *mdl.Design) []string {
	var keys []string
	if d.Views == nil {
		return keys
	}
	add := func(vs []*mdl.ViewProps) {
		for _, vp := range vs {
			if vp != nil && vp.Key != "" {
				keys = append(keys, vp.Key)
			}
		}
	}
	for _, v := range d.Views.LandscapeViews {
		add([]*mdl.ViewProps{v.ViewProps})
	}
	for _, v := range d.Views.ContextViews {
		add([]*mdl.ViewProps{v.ViewProps})
	}
	for _, v := range d.Views.ContainerViews {
		add([]*mdl.ViewProps{v.ViewProps})
	}
	for _, v := range d.Views.ComponentViews {
		add([]*mdl.ViewProps{v.ViewProps})
	}
	for _, v := range d.Views.DynamicViews {
		add([]*mdl.ViewProps{v.ViewProps})
	}
	for _, v := range d.Views.DeploymentViews {
		add([]*mdl.ViewProps{v.ViewProps})
	}
	// FilteredViews base on different struct, but still has Key
	for _, v := range d.Views.FilteredViews {
		if v != nil && v.Key != "" {
			keys = append(keys, v.Key)
		}
	}
	return keys
}

func waitForHTTP(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url) //nolint:gosec
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return nil
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for %s", url)
}

// renderViewsHeadless uses chromedp to visit each view and trigger auto-layout and save.
func renderViewsHeadless(baseURL, outDir string, views []string, cfg config) error {
	// Lazy import of chromedp via build tag is overkill; add module dependency and use directly
	// We inline a minimal wrapper to avoid leaking chromedp symbols elsewhere.

	// Prepare function that opens a headless tab to the URL and waits for SVG file
	run := func(url string, svgPath string, timeout time.Duration) error {
		// Defer importing chromedp to here for clarity
		return withChromedp(func(exec navigateExec) error {
			return exec(url, svgPath, timeout)
		})
	}

	for _, key := range views {
		// Build URL with automation params
		q := fmt.Sprintf("?id=%s&auto=1&save=1&direction=%s", key, strings.ToUpper(cfg.direction))
		if cfg.compact {
			q += "&compact=1"
		}
		url := baseURL + "/" + q

		// Remove any existing file to ensure fresh wait
		svgPath := filepath.Join(outDir, key+".svg")
		_ = os.Remove(svgPath)

		if err := run(url, svgPath, cfg.timeout); err != nil {
			return fmt.Errorf("render %s: %w", key, err)
		}
		fmt.Println("Saved:", svgPath)
	}
	return nil
}

func waitForFile(path string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if st, err := os.Stat(path); err == nil && st.Size() > 0 {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for %s", path)
}

// navigateExec abstracts chromedp.Navigate+Wait flow
type navigateExec func(url string, svgPath string, timeout time.Duration) error

// withChromedp wraps the chromedp session lifecycle
func withChromedp(fn func(exec navigateExec) error) error {
	return chromedpExec(fn)
}

// chromedpExec encapsulates direct chromedp usage.
func chromedpExec(fn func(exec navigateExec) error) error {
	// Use an explicit exec allocator with flags suitable for CI environments
	allocatorOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-setuid-sandbox", true),
	)
	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), allocatorOpts...)
	defer allocCancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	exec := func(url string, svgPath string, timeout time.Duration) error {
		// Use a tab context so the page stays open while we wait for the file
		tabCtx, tabCancel := chromedp.NewContext(ctx)
		defer tabCancel()

		navCtx, navCancel := context.WithTimeout(tabCtx, timeout)
		defer navCancel()
		if err := chromedp.Run(navCtx,
			chromedp.Navigate(url),
			// Wait for the graph svg to exist to ensure the app is ready
			chromedp.WaitVisible(`svg#graph`, chromedp.ByQuery),
		); err != nil {
			return err
		}

		// Keep tab alive while waiting for saved file
		if err := waitForFile(svgPath, timeout); err != nil {
			return err
		}
		return nil
	}

	return fn(exec)
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
	fmt.Fprintf(os.Stderr, "  %s svg PACKAGE [FLAGS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "    Auto-layout and export SVG diagram(s) for the design described in PACKAGE.\n")
	fmt.Fprintf(os.Stderr, "\nPACKAGE must be the import path to a Go package containing Model DSL.\n\n")
	fmt.Fprintf(os.Stderr, "FLAGS:\n")
	flag.PrintDefaults()
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
