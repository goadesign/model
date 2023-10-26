package svc

import (
	"context"
	"fmt"

	"goa.design/clue/log"
	genpackages "goa.design/model/svc/gen/packages"
	"goa.design/model/svc/gen/types"
)

// List the known workspaces
func (s *Service) ListWorkspaces(ctx context.Context) (res []*types.Workspace, err error) {
	ws, err := s.store.ListWorkspaces(ctx)
	if err != nil {
		return nil, logAndReturn(ctx, err, "failed to list workspaces")
	}
	return ws, nil
}

// CreatePackage creates a new package in the given workspace.
func (s *Service) CreatePackage(ctx context.Context, p *genpackages.CreatePackagePayload) (err error) {
	pf := &types.PackageFile{
		Locator: &types.FileLocator{
			Workspace: p.Workspace,
			Dir:       p.Dir,
			Filename:  "design.go",
		},
		Content: p.Content,
	}
	if err := s.store.CreatePackage(ctx, pf); err != nil {
		return logAndReturn(ctx, err, "failed to create package %s", p.Dir)
	}
	return nil
}

// DeletePackage deletes the given package.
func (s *Service) DeletePackage(ctx context.Context, p *types.PackageLocator) (err error) {
	if err := s.store.DeletePackage(ctx, p); err != nil {
		return logAndReturn(ctx, err, "failed to delete package %s", p.Dir)
	}
	return nil
}

// ListPackages lists the packages in the given workspace.
func (s *Service) ListPackages(ctx context.Context, w *types.Workspace) (res []*types.Package, err error) {
	ps, err := s.store.ListPackages(ctx, w.Workspace)
	if err != nil {
		return nil, logAndReturn(ctx, err, "failed to list packages in workspace %s", w.Workspace)
	}
	return ps, nil
}

// ReadPackageFiles lists the DSL files and their content for the given package.
func (s *Service) ReadPackageFiles(ctx context.Context, p *types.PackageLocator) (res []*types.PackageFile, err error) {
	fs, err := s.store.ReadPackageFiles(ctx, p)
	if err != nil {
		return nil, logAndReturn(ctx, err, "failed to read files in package %s", p.Dir)
	}
	return fs, nil
}

// Subscribe streams the model JSON for the given package.
func (s *Service) Subscribe(ctx context.Context, p *types.PackageLocator, stream genpackages.SubscribeServerStream) (err error) {
	ch, err := s.store.Subscribe(ctx, p)
	if err != nil {
		return logAndReturn(ctx, err, "failed to subscribe to package %s", p.Dir)
	}
	for m := range ch {
		if err := stream.Send(types.ModelJSON(m)); err != nil {
			return logAndReturn(ctx, err, "failed to send model JSON for package %s", p.Dir)
		}
	}
	return nil
}

// logAndReturn logs the given error and returns it.
func logAndReturn(ctx context.Context, err error, format string, args ...interface{}) error {
	log.Errorf(ctx, err, format, args...)
	return fmt.Errorf(format, args...)
}
