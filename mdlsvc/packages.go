package mdlsvc

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"goa.design/clue/log"
	"goa.design/model/codegen"
	genpackages "goa.design/model/mdlsvc/gen/packages"
)

// List the model pkgules in the current Go workspace
func (svc *Service) ListPackages(ctx context.Context) ([]*genpackages.Package, error) {
	_, err := os.Stat(filepath.Join(svc.workspace, "go.work"))
	if err != nil {
		_, err := os.Stat(filepath.Join(svc.workspace, "go.pkg"))
		if err != nil {
			return nil, nil
		}
		// return []*genpackages.Package{{
	}
	return nil, nil
}

// WebSocket endpoint for subscribing to updates to the model
func (svc *Service) Subscribe(ctx context.Context, pkg *genpackages.Package, stream genpackages.SubscribeServerStream) error {
	if js, err := codegen.JSON(svc.workspace, pkg.PackagePath, svc.debug); err == nil {
		if err := stream.Send(genpackages.Model(js)); err != nil {
			return err
		}
	} else {
		log.Errorf(ctx, err, "failed to generate JSON for %s", pkg.PackagePath)
	}
	svc.lock.Lock()
	sub, ok := svc.subscriptions[pkg.PackagePath]
	if ok {
		sub.streams = append(sub.streams, stream)
		svc.lock.Unlock()
	} else {
		ch := make(chan []byte)
		sub = &subscription{ch: ch, streams: []genpackages.SubscribeServerStream{stream}}
		svc.subscriptions[pkg.PackagePath] = sub
		svc.lock.Unlock()
		if err := watch(ctx, svc.workspace, pkg.PackagePath, func() {
			js, err := codegen.JSON(svc.workspace, pkg.PackagePath, svc.debug)
			if err != nil {
				log.Errorf(ctx, err, "failed to generate JSON for %s", pkg.PackagePath)
				return
			}
			ch <- js
		}); err != nil {
			return err
		}
		go func() {
			defer func() {
				close(ch)
			}()
			for {
				select {
				case <-ctx.Done():
					return
				case update := <-ch:
					svc.lock.Lock()
					streams := svc.subscriptions[pkg.PackagePath].streams
					svc.lock.Unlock()
					for _, stream := range streams {
						if err := stream.Send(genpackages.Model(update)); err != nil {
							log.Errorf(ctx, err, "failed to send update for %s", pkg.PackagePath)
						}
					}
				}
			}
		}()
	}
	defer func() {
		svc.lock.Lock()
		streams := svc.subscriptions[pkg.PackagePath].streams
		for i, s := range streams {
			if s == stream {
				streams = append(streams[:i], streams[i+1:]...)
				break
			}
		}
		if len(streams) == 0 {
			delete(svc.subscriptions, pkg.PackagePath)
		} else {
			svc.subscriptions[pkg.PackagePath].streams = streams
		}
		svc.lock.Unlock()
	}()
	for {
		select {
		case <-ctx.Done():
			return nil
		case update := <-sub.ch:
			if err := stream.Send(genpackages.Model(update)); err != nil {
				return err
			}
		}
	}
}

// Upload the package content, compile it and return the corresponding JSON
func (svc *Service) Upload(ctx context.Context, pkg *genpackages.Package, rc io.ReadCloser) (res genpackages.Model, err error) {
	panic("not implemented")
}

// Stream the model JSON, see https://pkg.go.dev/goa.design/model/model#model
func (svc *Service) GetModel(ctx context.Context, pkg *genpackages.Package) (body io.ReadCloser, err error) {
	panic("not implemented")
}
