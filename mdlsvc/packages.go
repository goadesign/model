package mdlsvc

import (
	"bytes"
	"context"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"

	"goa.design/clue/log"
	"goa.design/model/codegen"
	genpackages "goa.design/model/mdlsvc/gen/packages"
	gentypes "goa.design/model/mdlsvc/gen/types"
)

// List the model packages located below the service packages directory.
func (svc *Service) ListPackages(ctx context.Context) ([]*gentypes.Package, error) {
	var pkgs []*gentypes.Package
	err := filepath.Walk(svc.dir, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			log.Errorf(ctx, err, "failed to stat %s", fpath)
			return nil
		}
		if !info.IsDir() {
			return nil
		}
		// Parse directory and look for goa.design/model/dsl imports
		parsed, err := parser.ParseDir(token.NewFileSet(), fpath, nil, parser.ImportsOnly)
		if err != nil {
			log.Errorf(ctx, err, "failed to parse %s", fpath)
			return nil
		}
		if len(parsed) == 0 {
			return nil
		}
		fmpath, mpath, err := modulePath(fpath)
		if err != nil {
			log.Errorf(ctx, err, "failed to find module for %s", fpath)
			return nil
		}
		for _, pkg := range parsed {
			for _, f := range pkg.Files {
				for _, imp := range f.Imports {
					if imp.Path.Value == `"goa.design/model/dsl"` {
						pkgPath, err := packagePath(pkg.Name, fpath, mpath, fmpath)
						if err != nil {
							log.Error(ctx, err)
							continue
						}
						pkgs = append(pkgs, &gentypes.Package{
							Dir:        fpath,
							ImportPath: pkgPath,
						})
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Errorf(ctx, err, "failed to list packages")
		return nil, err
	}
	return pkgs, nil
}

// WebSocket endpoint for subscribing to updates to the model
func (svc *Service) Subscribe(ctx context.Context, dir *gentypes.PackageLocator, stream genpackages.SubscribeServerStream) error {
	pkgPath, err := findPackagePath(dir.Dir)
	if err != nil {
		return logAndReturn(ctx, err, "failed to find package path for %s", dir.Dir)
	}
	js, err := codegen.JSON(dir.Dir, pkgPath, svc.debug)
	if err != nil {
		return logAndReturn(ctx, err, "failed to generate JSON for %s, package %s", dir.Dir, pkgPath)
	}
	if err := stream.Send(gentypes.ModelJSON(js)); err != nil {
		return logAndReturn(ctx, err, "failed to send update for %s, package %s", dir.Dir, pkgPath)
	}
	svc.lock.Lock()
	sub, ok := svc.subscriptions[dir.Dir]
	if ok {
		sub.streams = append(sub.streams, stream)
		svc.lock.Unlock()
	} else {
		sub, err = svc.subscribe(ctx, dir.Dir, pkgPath, stream)
		svc.lock.Unlock()
		if err != nil {
			return logAndReturn(ctx, err, "failed to subscribe to %s", dir.Dir)
		}
	}
	<-ctx.Done() // wait for client to disconnect
	svc.lock.Lock()
	for i, s := range sub.streams {
		if s == stream {
			sub.streams = append(sub.streams[:i], sub.streams[i+1:]...)
			break
		}
	}
	if len(sub.streams) == 0 {
		delete(svc.subscriptions, dir.Dir)
	}
	svc.lock.Unlock()
	return nil
}

// Get the files and their content for the given package
func (svc *Service) ListPackageFiles(ctx context.Context, dir *gentypes.PackageLocator) ([]*gentypes.PackageFile, error) {
	files, err := os.ReadDir(dir.Dir)
	if err != nil {
		return nil, logAndReturn(ctx, err, "failed to read directory %s", dir.Dir)
	}
	res := make([]*gentypes.PackageFile, 0, len(files))
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir.Dir, f.Name()))
		if err != nil {
			return nil, logAndReturn(ctx, err, "failed to read file %s", f.Name())
		}
		res = append(res, &gentypes.PackageFile{
			Locator: &gentypes.FileLocator{
				Workspace: dir.Workspace,
				Dir:       dir.Dir,
				Filename:  f.Name(),
			},
			Content: string(b),
		})
	}
	return res, nil
}

// Stream the model JSON, see https://pkg.go.dev/goa.design/model/model#model
func (svc *Service) GetModelJSON(ctx context.Context, dir *gentypes.PackageLocator) (io.ReadCloser, error) {
	pkgPath, err := findPackagePath(dir.Dir)
	if err != nil {
		return nil, logAndReturn(ctx, err, "failed to find package path for %s", dir.Dir)
	}
	js, err := codegen.JSON(dir.Dir, pkgPath, svc.debug)
	if err != nil {
		return nil, logAndReturn(ctx, err, "failed to generate JSON for %s", dir.Dir)
	}
	return io.NopCloser(bytes.NewBuffer(js)), nil
}

// Subscribe to updates to the model, svc lock must be held.
func (svc *Service) subscribe(ctx context.Context, dir, pkgPath string, stream genpackages.SubscribeServerStream) (*subscription, error) {
	sub := &subscription{streams: []genpackages.SubscribeServerStream{stream}}
	svc.subscriptions[dir] = sub
	err := watch(ctx, svc.dir, dir, func() {
		js, err := codegen.JSON(dir, pkgPath, svc.debug)
		if err != nil {
			log.Errorf(ctx, err, "failed to generate JSON for %s", dir)
			return
		}
		svc.lock.Lock()
		streams := svc.subscriptions[dir].streams
		svc.lock.Unlock()
		for _, s := range streams {
			if err := s.Send(gentypes.ModelJSON(js)); err != nil {
				log.Errorf(ctx, err, "failed to send update for %s", dir)
			}
		}
	})
	if err != nil {
		delete(svc.subscriptions, dir)
		return nil, err
	}
	return sub, nil
}

// logAndReturn logs the given error and returns it.
func logAndReturn(ctx context.Context, err error, format string, args ...interface{}) error {
	log.Errorf(ctx, err, format, args...)
	return fmt.Errorf(format, args...)
}
