package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"goa.design/clue/debug"
	"goa.design/clue/log"
	goahttp "goa.design/goa/v3/http"

	svc "goa.design/model/svc"
	"goa.design/model/svc/clients/repo"
	geneditorsvc "goa.design/model/svc/gen/dsl_editor"
	genassetshttp "goa.design/model/svc/gen/http/assets/server"
	geneditorhttp "goa.design/model/svc/gen/http/dsl_editor/server"
	genrepohttp "goa.design/model/svc/gen/http/repo/server"
	gensvghttp "goa.design/model/svc/gen/http/svg/server"
	genrepovc "goa.design/model/svc/gen/repo"
	gensvgvc "goa.design/model/svc/gen/svg"
)

//go:embed webapp/dist
var webapp embed.FS

func serve(store repo.RepoHandler, port int, devmode, debugf bool) error {
	format := log.FormatJSON
	if log.IsTerminal() {
		format = log.FormatTerminal
	}
	ctx := log.Context(context.Background(), log.WithFormat(format))
	ctx = log.With(ctx, log.KV{K: "svc", V: "mdl"})
	if debugf {
		ctx = log.Context(ctx, log.WithDebug())
		log.Debugf(ctx, "debug logs enabled")
	}

	svc, err := svc.New(ctx, store, debugf)
	if err != nil {
		return err
	}
	var fs http.FileSystem
	if devmode {
		// in devmode (go run), serve the webapp from filesystem
		fs = http.FileSystem(http.Dir("./webapp/dist"))
		http.Handle("/", http.FileServer(fs))
	} else {
		// the TS/React webapp is embeded in the go executable using embed
		fs = http.FS(webapp)
	}

	mux := goahttp.NewMuxer()

	editorEndpoints := geneditorsvc.NewEndpoints(svc)
	editorEndpoints.Use(debug.LogPayloads())
	editorEndpoints.Use(log.Endpoint)
	editorsvr := geneditorhttp.New(editorEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)
	geneditorhttp.Mount(mux, editorsvr)
	for _, m := range editorsvr.Mounts {
		log.Print(ctx, log.KV{K: "method", V: m.Method}, log.KV{K: "endpoint", V: m.Verb + " " + m.Pattern})
	}

	packagesEndpoints := genrepovc.NewEndpoints(svc)
	packagesEndpoints.Use(debug.LogPayloads())
	packagesEndpoints.Use(log.Endpoint)
	modulesvr := genrepohttp.New(packagesEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil, &websocket.Upgrader{}, nil)
	genrepohttp.Mount(mux, modulesvr)
	for _, m := range modulesvr.Mounts {
		log.Print(ctx, log.KV{K: "method", V: m.Method}, log.KV{K: "endpoint", V: m.Verb + " " + m.Pattern})
	}

	svgEndpoints := gensvgvc.NewEndpoints(svc)
	svgEndpoints.Use(debug.LogPayloads())
	svgEndpoints.Use(log.Endpoint)
	svgsvr := gensvghttp.New(svgEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)
	gensvghttp.Mount(mux, svgsvr)
	for _, m := range svgsvr.Mounts {
		log.Print(ctx, log.KV{K: "method", V: m.Method}, log.KV{K: "endpoint", V: m.Verb + " " + m.Pattern})
	}

	assetssvr := genassetshttp.New(nil, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil, fs)
	genassetshttp.Mount(mux, assetssvr)
	for _, m := range assetssvr.Mounts {
		log.Print(ctx, log.KV{K: "method", V: m.Method}, log.KV{K: "endpoint", V: m.Verb + " " + m.Pattern})
	}

	debug.MountDebugLogEnabler(debug.Adapt(mux))
	handler := log.HTTP(ctx)(mux)
	addr := fmt.Sprintf(":%d", port)
	httpServer := &http.Server{Addr: addr, Handler: handler}

	errc := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		go func() {
			log.Printf(ctx, "HTTP server listening on %s", addr)
			errc <- httpServer.ListenAndServe()
		}()

		<-ctx.Done()
		log.Printf(ctx, "shutting down HTTP server")

		// Shutdown gracefully with a 30s timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			log.Errorf(ctx, err, "failed to shutdown HTTP server")
		}
	}()

	// Cleanup
	if err := <-errc; err != nil {
		log.Errorf(ctx, err, "exiting")
	}
	cancel()
	wg.Wait()
	log.Printf(ctx, "exited")
	return nil
}
