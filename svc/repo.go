package svc

import (
	"context"
	"fmt"
	"path/filepath"

	"goa.design/clue/log"
	"goa.design/model/codegen"
	"goa.design/model/svc/clients/repo"
	genrepo "goa.design/model/svc/gen/repo"
	gentypes "goa.design/model/svc/gen/types"
)

// defaultPackageContent is the default content for a new package.
const defaultPackageContent = `package model

import . "goa.design/model/dsl"

var _ = Design("model", "System architecture model", func() {
})
`

// ListPackages lists the packages in the given workspace.
func (svc *Service) ListPackages(ctx context.Context, w *gentypes.Repository) ([]*gentypes.Package, error) {
	ctx = log.With(ctx, log.KV{K: "repo", V: w.Repository})
	ps, err := svc.handler.ListPackages(ctx, w.Repository)
	if err != nil {

		return nil, logAndReturn(ctx, err)
	}
	return ps, nil
}

// ReadPackage lists the DSL files and their content for the given package.
func (svc *Service) ReadPackage(ctx context.Context, p *gentypes.PackageLocator) ([]*gentypes.PackageFile, error) {
	ctx = log.With(ctx, log.KV{K: "repo", V: p.Repository}, log.KV{K: "dir", V: p.Dir})
	fs, err := svc.handler.ReadPackage(ctx, p)
	if err != nil {
		return nil, logAndReturn(ctx, err)
	}
	return fs, nil
}

// CreatePackage creates a new package in the given workspace.
func (svc *Service) CreatePackage(ctx context.Context, pf *gentypes.PackageFile) error {
	ctx = log.With(ctx, log.KV{K: "repo", V: pf.Locator.Repository}, log.KV{K: "dir", V: pf.Locator.Dir})
	if err := svc.handler.CreatePackage(ctx, pf); err != nil {
		if err == repo.ErrAlreadyExists {
			return genrepo.MakeAlreadyExists(err)
		}
		return logAndReturn(ctx, err)
	}
	return nil
}

// CreateDefaultPackage creates a new package with default content in the given workspace.
func (svc *Service) CreateDefaultPackage(ctx context.Context, p *gentypes.FileLocator) error {
	ctx = log.With(ctx, log.KV{K: "repo", V: p.Repository}, log.KV{K: "dir", V: p.Dir})
	pf := &gentypes.PackageFile{
		Locator: p,
		Content: defaultPackageContent,
	}
	if err := svc.handler.CreatePackage(ctx, pf); err != nil {
		if err == repo.ErrAlreadyExists {
			return genrepo.MakeAlreadyExists(err)
		}
		return logAndReturn(ctx, err)
	}
	return nil
}

// DeletePackage deletes the given package.
func (svc *Service) DeletePackage(ctx context.Context, p *gentypes.PackageLocator) error {
	ctx = log.With(ctx, log.KV{K: "repo", V: p.Repository}, log.KV{K: "dir", V: p.Dir})
	if err := svc.handler.DeletePackage(ctx, p); err != nil {
		if err == repo.ErrNotFound {
			return genrepo.MakeNotFound(err)
		}
		return logAndReturn(ctx, err)
	}
	return nil
}

// GetModelJSON returns the model JSON for the given package.
func (svc *Service) GetModelJSON(ctx context.Context, p *gentypes.PackageLocator) (gentypes.ModelJSON, error) {
	ctx = log.With(ctx, log.KV{K: "repo", V: p.Repository}, log.KV{K: "dir", V: p.Dir})
	pkgPath, err := svc.handler.GetImportPath(ctx, p)
	if err != nil {
		if err == repo.ErrNotFound {
			return "", genrepo.MakeNotFound(err)
		}
		return "", logAndReturn(ctx, err, "failed to get import path")
	}
	js, err := codegen.JSON(filepath.Join(p.Repository, p.Dir), pkgPath, svc.debug)
	if err != nil {
		return "", genrepo.MakeCompilationError(err)
	}
	return gentypes.ModelJSON(js), nil
}

// Subscribe streams the model JSON for the given package.
func (svc *Service) Subscribe(ctx context.Context, p *gentypes.PackageLocator, stream genrepo.SubscribeServerStream) (err error) {
	ctx = log.With(ctx, log.KV{K: "repo", V: p.Repository}, log.KV{K: "dir", V: p.Dir})
	pkgPath, err := svc.handler.GetImportPath(ctx, p)
	if err != nil {
		if err == repo.ErrNotFound {
			return genrepo.MakeNotFound(err)
		}
		return logAndReturn(ctx, err, "failed to get import path")
	}
	ch, err := svc.handler.Subscribe(ctx, p)
	if err != nil {
		if err == repo.ErrNotFound {
			return genrepo.MakeNotFound(err)
		}
		return logAndReturn(ctx, err, "failed to subscribe")
	}
	for range ch {
		js, err := codegen.JSON(filepath.Join(p.Repository, p.Dir), pkgPath, svc.debug)
		if err != nil {
			msg := err.Error()
			log.Infof(ctx, "compilation error: %s", err)
			err = stream.Send(&gentypes.CompilationResults{Error: &msg})
		} else {
			model := gentypes.ModelJSON(js)
			log.Info(ctx, log.KV{K: "msg", V: "compilation success"}, log.KV{K: "bytes", V: len(model)})
			err = stream.Send(&gentypes.CompilationResults{Model: &model})
		}
		if err != nil {
			return logAndReturn(ctx, err, "failed to stream compilation results")
		}
	}
	return nil
}

// logAndReturn logs the given error and returns it.
func logAndReturn(ctx context.Context, err error, msg ...string) error {
	if len(msg) > 0 {
		err = fmt.Errorf(msg[0]+": %w", err)
	}
	log.Error(ctx, err)
	return err
}

// fpath returns the file path for the given package file.
func fpath(p *gentypes.FileLocator) string {
	return filepath.Join(p.Repository, p.Dir, p.Filename)
}
