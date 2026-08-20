package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	goacodegen "goa.design/goa/v3/codegen"

	"goa.design/model/codegen"
	"goa.design/model/mdl"
	model "goa.design/model/pkg"

	cdnetwork "github.com/chromedp/cdproto/network"
	cdruntime "github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

type (
	config struct {
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
		force     bool
	}

	// SliceFlag implements flag.Value for repeated string flags.
	SliceFlag []string
)

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
	case "skill":
		if pkg != "install" {
			err = fmt.Errorf(`unknown skill command %q, expected "install"`, pkg)
		} else {
			err = runSkillInstall(cfg.force)
		}
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
		port:    0,
		devmode: os.Getenv("DEVMODE") == "1",
		devdist: os.Getenv("DEVDIST"),
		// defaults for svg command
		timeout: 20 * time.Second,
	}

	flag.BoolVar(&cfg.debug, "debug", false, "print debug output")
	flag.BoolVar(&cfg.help, "help", false, "print this information")
	flag.BoolVar(&cfg.help, "h", false, "print this information")
	flag.StringVar(&cfg.out, "out", cfg.out, "set path to generated JSON representation")
	flag.StringVar(&cfg.dir, "dir", cfg.dir, "set output directory used by editor to save SVG files")
	flag.IntVar(
		&cfg.port,
		"port",
		cfg.port,
		"set local HTTP port; serve defaults to 8080 and svg defaults to a private free port",
	)
	// svg command flags (safe to always register)
	flag.Var(&cfg.views, "view", "view key to render (repeatable)")
	flag.BoolVar(&cfg.all, "all", false, "render all views")
	flag.StringVar(
		&cfg.direction,
		"direction",
		cfg.direction,
		"override the view auto-layout direction: DOWN|UP|LEFT|RIGHT",
	)
	flag.BoolVar(&cfg.compact, "compact", false, "enable compact auto-layout")
	flag.DurationVar(&cfg.timeout, "timeout", cfg.timeout, "timeout per view (e.g. 15s)")
	flag.BoolVar(&cfg.force, "force", false, "replace a locally modified installed skill")

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
	if cfg.port == 0 {
		cfg.port = 8080
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

// runSVG serves one fixed model, renders selected views in one browser process,
// and saves only matched, validated browser results.
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

	design, err := loadDesign(pkg, cfg.debug)
	if err != nil {
		return err
	}

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

	server := NewServer(design)
	digest, err := designDigest(design)
	if err != nil {
		return err
	}
	broker := newRenderBroker()
	mux := http.NewServeMux()
	mux.HandleFunc("/headless/result", broker.handleResult)

	listener, err := net.Listen("tcp", listenAddress(cfg.port))
	if err != nil {
		return fmt.Errorf("bind headless server: %w", err)
	}

	httpServer := &http.Server{
		ReadHeaderTimeout: 3 * time.Second,
	}

	done := make(chan error, 1)
	go func() {
		done <- server.ServeOnListener(absDir, cfg.devdist, httpServer, mux, listener)
	}()

	baseURL := "http://" + listener.Addr().String()
	if err := renderViewsHeadless(baseURL, digest, selected, cfg, broker, server); err != nil {
		if closeErr := httpServer.Close(); closeErr != nil {
			return fmt.Errorf("%w; close headless server: %v", err, closeErr)
		}
		return err
	}

	if err := httpServer.Close(); err != nil {
		return fmt.Errorf("close headless server: %w", err)
	}
	select {
	case err := <-done:
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("headless server: %w", err)
		}
	case <-time.After(2 * time.Second):
		return fmt.Errorf("timeout stopping headless server")
	}
	return nil
}

// designDigest gives every browser result the identity of the exact model JSON.
func designDigest(design *mdl.Design) (string, error) {
	data, err := json.Marshal(design)
	if err != nil {
		return "", fmt.Errorf("serialize design digest: %w", err)
	}
	digest := sha256.Sum256(data)
	return fmt.Sprintf("%x", digest), nil
}

// listenAddress uses a caller port or asks the operating system for a free one.
func listenAddress(port int) string {
	if port == 0 {
		return "127.0.0.1:0"
	}
	return fmt.Sprintf("127.0.0.1:%d", port)
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

// renderViewsHeadless reuses one browser process and waits for typed HTTP results.
func renderViewsHeadless(
	baseURL string,
	modelDigest string,
	views []string,
	cfg config,
	broker *renderBroker,
	server *Server,
) error {
	direction, err := normalizeLayoutDirection(cfg.direction)
	if err != nil {
		return err
	}
	return withChromedp(cfg.timeout, cfg.debug, func(exec navigateExec) error {
		for _, viewID := range views {
			if err := renderViewHeadless(
				baseURL,
				modelDigest,
				viewID,
				direction,
				cfg,
				broker,
				server,
				exec,
			); err != nil {
				return fmt.Errorf("render %s: %w", viewID, err)
			}
		}
		return nil
	})
}

// renderViewHeadless waits for the exact view and model result before saving it.
func renderViewHeadless(
	baseURL string,
	modelDigest string,
	viewID string,
	direction string,
	cfg config,
	broker *renderBroker,
	server *Server,
	exec navigateExec,
) error {
	results, unregister, err := broker.register(viewID, modelDigest)
	if err != nil {
		return err
	}
	defer unregister()

	deadline := time.Now().Add(cfg.timeout)
	renderURL := headlessRenderURL(baseURL, modelDigest, viewID, direction, cfg.compact)
	if err := exec(renderURL, time.Until(deadline)); err != nil {
		return fmt.Errorf("open headless page: %w", err)
	}
	remaining := time.Until(deadline)
	if remaining <= 0 {
		return fmt.Errorf("timeout waiting for browser result")
	}
	timer := time.NewTimer(remaining)
	defer timer.Stop()

	var result headlessResult
	select {
	case result = <-results:
	case <-timer.C:
		return fmt.Errorf("timeout waiting for browser result")
	}
	if result.Status == "error" {
		return fmt.Errorf("browser failed: %s", result.Error)
	}

	err = server.storeSVG(viewID, strings.NewReader(result.SVG))
	if err != nil {
		return fmt.Errorf("save SVG: %w", err)
	}
	fmt.Println("Saved:", filepath.Join(server.outDir, viewID+".svg"))
	return nil
}

// headlessRenderURL encodes one view request without relying on router state.
func headlessRenderURL(
	baseURL string,
	modelDigest string,
	viewID string,
	direction string,
	compact bool,
) string {
	query := url.Values{
		"digest": {modelDigest},
		"view":   {viewID},
	}
	if direction != "" {
		query.Set("direction", direction)
	}
	if compact {
		query.Set("compact", "true")
	}
	return baseURL + "/headless.html?" + query.Encode()
}

func normalizeLayoutDirection(direction string) (string, error) {
	normalized := strings.ToUpper(direction)
	switch normalized {
	case "", "DOWN", "UP", "LEFT", "RIGHT":
		return normalized, nil
	default:
		return "", fmt.Errorf(
			"invalid auto-layout direction %q: use DOWN, UP, LEFT, or RIGHT",
			direction,
		)
	}
}

// navigateExec opens one isolated page in the shared browser process.
type navigateExec func(url string, timeout time.Duration) error

// withChromedp wraps the chromedp session lifecycle
func withChromedp(timeout time.Duration, debug bool, fn func(exec navigateExec) error) error {
	return chromedpExec(timeout, debug, fn)
}

// chromedpExec encapsulates direct chromedp usage.
func chromedpExec(timeout time.Duration, debug bool, fn func(exec navigateExec) error) error {
	// Use an explicit exec allocator with flags suitable for CI environments
	allocatorOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.WSURLReadTimeout(timeout),
		chromedp.DisableGPU,
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-setuid-sandbox", true),
	)
	if chromeBin := os.Getenv("CHROME_BIN"); chromeBin != "" {
		allocatorOpts = append(allocatorOpts, chromedp.ExecPath(chromeBin))
	}
	if debug {
		allocatorOpts = append(allocatorOpts, chromedp.CombinedOutput(os.Stderr))
	}
	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), allocatorOpts...)
	defer allocCancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	tabCtx, tabCancel := chromedp.NewContext(ctx)
	defer tabCancel()
	if err := chromedp.Run(tabCtx); err != nil {
		return fmt.Errorf("start headless browser tab: %w", err)
	}
	if debug {
		if err := chromedp.Run(tabCtx, cdnetwork.Enable()); err != nil {
			return fmt.Errorf("enable browser network diagnostics: %w", err)
		}
		chromedp.ListenTarget(tabCtx, func(event any) {
			switch typed := event.(type) {
			case *cdruntime.EventExceptionThrown:
				fmt.Fprintf(
					os.Stderr,
					"browser exception: %v\n",
					typed.ExceptionDetails,
				)
			case *cdruntime.EventConsoleAPICalled:
				for _, argument := range typed.Args {
					fmt.Fprintf(
						os.Stderr,
						"browser console: %s %s\n",
						argument.Description,
						string(argument.Value),
					)
				}
			case *cdnetwork.EventRequestWillBeSent:
				fmt.Fprintf(os.Stderr, "browser request: %s\n", typed.Request.URL)
			case *cdnetwork.EventResponseReceived:
				fmt.Fprintf(
					os.Stderr,
					"browser response: %d %s\n",
					typed.Response.Status,
					typed.Response.URL,
				)
			case *cdnetwork.EventLoadingFailed:
				fmt.Fprintf(os.Stderr, "browser request failed: %s\n", typed.ErrorText)
			}
		})
	}

	exec := func(url string, timeout time.Duration) error {
		if timeout <= 0 {
			return fmt.Errorf("browser navigation deadline exceeded")
		}
		navCtx, navCancel := context.WithTimeout(tabCtx, timeout)
		defer navCancel()
		if err := chromedp.Run(navCtx, chromedp.Navigate(url)); err != nil {
			return fmt.Errorf("navigate browser: %w", err)
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
	fmt.Fprintf(os.Stderr, "  %s skill install [-force]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "    Install the MDL diagram-editing skill in the current repository.\n")
	fmt.Fprintf(os.Stderr, "\nPACKAGE must be the import path to a Go package containing Model DSL.\n")
	fmt.Fprintf(os.Stderr, "PACKAGE is required by serve, gen, and svg.\n\n")
	fmt.Fprintf(os.Stderr, "FLAGS:\n")
	flag.PrintDefaults()
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
