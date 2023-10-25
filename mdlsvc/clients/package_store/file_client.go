package pstore

import (
	"context"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"goa.design/clue/log"

	"goa.design/model/codegen"
	gentypes "goa.design/model/mdlsvc/gen/types"
)

type (
	// fileClient provides an implementation of the package store client
	// that uses the local filesystem as a backend.
	fileClient struct {
		debug bool
	}
)

// NewFileClient returns a file client.
func NewFileClient(debug bool) Client {
	return &fileClient{debug: debug}
}

// ListPackages returns the list of model packages in the given workspace.
func (c *fileClient) ListPackages(ctx context.Context, workspace string) ([]*gentypes.Package, error) {
	// Walk through the workspace directory and look for directories that
	// contain a Go file that imports goa.design/model/dsl.
	var pkgs []*gentypes.Package
	err := filepath.Walk(workspace, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			log.Errorf(ctx, err, "failed to stat %s", fpath)
			return nil
		}
		if !info.IsDir() {
			return nil
		}
		// Parse directory using Go parser and look for goa.design/model/dsl imports
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
		return nil, fmt.Errorf("failed to list packages: %w", err)
	}
	return pkgs, nil
}

// ListPackageFiles returns the list of files in the given model package.
func (c *fileClient) ReadPackageFiles(ctx context.Context, p *gentypes.PackageLocator) ([]*gentypes.PackageFile, error) {
	dir := filepath.Join(p.Workspace, p.Dir)
	files, err := listGoFiles(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to list files in %s", p.Dir)
	}
	res := make([]*gentypes.PackageFile, 0, len(files))
	for _, f := range files {
		b, err := os.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s", f.Name())
		}
		res = append(res, &gentypes.PackageFile{
			Locator: &gentypes.FileLocator{
				Workspace: p.Workspace,
				Dir:       p.Dir,
				Filename:  f.Name(),
			},
			Content: string(b),
		})
	}
	return res, nil
}

// Save implements Save.
func (c *fileClient) Save(ctx context.Context, p *gentypes.PackageFile) error {
	err := os.WriteFile(filepath.Join(p.Locator.Workspace, p.Locator.Dir, p.Locator.Filename), []byte(p.Content), 0644)
	if err != nil {
		log.Errorf(ctx, err, "failed to save %s", filepath.Join(p.Locator.Dir, p.Locator.Filename))
		return err
	}
	return nil
}

// Subscribe listens to changes in the given model package files. When
// notifications are received from the filesystem, the model is rebuild and
// its JSON written to the returned channel. Call the returned stop function to
// stop listening.
func (c *fileClient) Subscribe(ctx context.Context, p *gentypes.PackageLocator) (chan []byte, func(), error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Errorf(ctx, err, "Error watching files")
		return nil, nil, err
	}
	dir := filepath.Join(p.Workspace, p.Dir)
	if err = watcher.Add(dir); err != nil {
		log.Errorf(ctx, err, "Error watching files in %s", dir)
		return nil, nil, err
	}
	stop := make(chan struct{})
	ch := make(chan []byte)
	go func() {
		defer func() {
			watcher.Close()
			close(ch)
		}()
		for {
			select {
			case ev := <-watcher.Events:
				if strings.HasPrefix(filepath.Base(ev.Name), codegen.TmpDirPrefix) {
					// ignore temporary (generated) files
					continue
				}
				pkgPath, err := findPackagePath(dir)
				if err != nil {
					log.Errorf(ctx, err, "failed to find package path for %s", dir)
					continue
				}

				// debounce, because some editors do several file operations when you save
				// we wait for the stream of events to become silent for `interval`
				interval := 100 * time.Millisecond
				timer := time.NewTimer(interval)
			outer:
				for {
					select {
					case ev = <-watcher.Events:
						timer.Reset(interval)
					case <-timer.C:
						break outer
					}
				}

				log.Infof(ctx, ev.String())
				js, err := codegen.JSON(dir, pkgPath, c.debug)
				if err != nil {
					log.Errorf(ctx, err, "failed to generate JSON for %s", dir)
					continue
				}
				ch <- js

			case err := <-watcher.Errors:
				log.Errorf(ctx, err, "Error watching files")

			case <-stop:
				return
			}
		}
	}()

	return ch, func() { close(stop) }, nil
}

// modulePath searches for the go.mod file recursively starting at the given
// directory and returns the corresponding file and module path.
func modulePath(dir string) (string, string, error) {
	fpath := filepath.Join(dir, "go.mod")
	if _, err := os.Stat(fpath); err == nil {
		f, err := os.Open(fpath)
		if err != nil {
			return "", "", err
		}
		defer f.Close()
		// read first line
		var line string
		if _, err := fmt.Fscanf(f, "%s", &line); err != nil {
			return "", "", err
		}
		if line != "module" {
			return "", "", fmt.Errorf("invalid go.mod file %s", fpath)
		}
		// read module path
		var mod string
		if _, err := fmt.Fscanf(f, "%s", &mod); err != nil {
			return "", "", err
		}
		return fpath, mod, nil
	}
	parent := filepath.Dir(dir)
	if parent == dir {
		return "", "", fmt.Errorf("failed to find go.mod file")
	}
	return modulePath(parent)
}

// findPackagePath returns the Go import path for the package with the given
// directory. For example:
//
//	findPackagePath("testdata/model")
//
// returns "goa.design/model/model".
func findPackagePath(dir string) (string, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.PackageClauseOnly)
	if err != nil {
		return "", err
	}
	for _, pkg := range pkgs {
		return pkg.Name, nil
	}
	return "", fmt.Errorf("failed to find package in %s", dir)
}

// packagePath returns the Go import path for the package with the given name in
// the given directory that lives under the Go module with the given Go and file
// paths. For example:
//
//	packagePath("model", "testdata/model", "goa.design/model", "testdata/go.mod")
//
// returns "goa.design/model/model".
func packagePath(pkgName, pkgDir, modGoPath, modFilePath string) (string, error) {
	rel, err := filepath.Rel(filepath.Dir(modFilePath), pkgDir)
	if err != nil {
		return "", fmt.Errorf("failed to compute relative path for %s: %w", pkgDir, err)
	}
	rel = filepath.ToSlash(rel)
	return path.Join(modGoPath, rel, pkgName), nil
}
