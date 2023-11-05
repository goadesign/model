package repo

import (
	"context"
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"goa.design/clue/log"
	"goa.design/model/codegen"
	gentypes "goa.design/model/svc/gen/types"
	"gopkg.in/fsnotify.v1"
)

type (
	// RepoHandler is the repository client used by the model service.
	RepoHandler interface {
		// ListPackages returns the list of model packages in the given
		// repo.
		ListPackages(context.Context, string) ([]*gentypes.Package, error)
		// ReadPackage returns the list of files in the given model
		// package.
		ReadPackage(context.Context, *gentypes.PackageLocator) ([]*gentypes.PackageFile, error)
		// CreatePackage creates a new model package.
		CreatePackage(context.Context, *gentypes.PackageFile) error
		// DeletePackage deletes the given package.
		DeletePackage(context.Context, *gentypes.PackageLocator) error
		// GetImportPath returns the import path for the given package.
		GetImportPath(context.Context, *gentypes.PackageLocator) (string, error)
		// Save persists the content of the given file for future
		// retrieval with LoadPackageFiles.
		Save(context.Context, *gentypes.PackageFile) error
		// Subscribe sends notifications when files in the give package
		// change.  The subscription is closed when the context is
		// canceled.
		Subscribe(context.Context, *gentypes.PackageLocator) (<-chan struct{}, error)
	}

	// handler provides an implementation of the repo handler client that
	// uses the local filesystem as a backend.
	handler struct {
		root  string
		debug bool
	}
)

// ErrAlreadyExists is the error returned when a package already exists.
var ErrAlreadyExists = fmt.Errorf("already exists")

// ErrNotFound is the error returned when a package is not found.
var ErrNotFound = fmt.Errorf("not found")

// NewHandler returns a file system based repo handler rooted at the given
// directory. A file based package store stores packages under
// <root>/<repo>/<package dir>.
func NewHandler(root string, debug bool) RepoHandler {
	return &handler{root: root, debug: debug}
}

// CreatePackage creates a new model package.
func (h *handler) CreatePackage(ctx context.Context, pf *gentypes.PackageFile) error {
	dir := filepath.Join(h.root, pf.Locator.Repository, pf.Locator.Dir)
	// Check if directory already exists
	if _, err := os.Stat(dir); err == nil {
		return ErrAlreadyExists
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	file := filepath.Join(dir, "model.go")
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", file, err)
	}
	defer f.Close()
	if _, err := f.WriteString(pf.Content); err != nil {
		return fmt.Errorf("failed to write file %s: %w", file, err)
	}
	return nil
}

// DeletePackage deletes the given package.
func (h *handler) DeletePackage(ctx context.Context, p *gentypes.PackageLocator) error {
	dir := filepath.Join(p.Repository, p.Dir)
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to stat directory %s: %w", dir, err)
	}
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("failed to delete directory %s: %w", dir, err)
	}
	return nil
}

// ListPackages returns the list of model packages in the given repo.
func (h *handler) ListPackages(ctx context.Context, repo string) ([]*gentypes.Package, error) {
	var pkgs []*gentypes.Package
	err := iterateModelPackages(ctx, repo, func(fpath, gopath string) {
		pkgs = append(pkgs, &gentypes.Package{
			Dir:        fpath,
			ImportPath: gopath,
		})
	})
	if err != nil {
		if err == ErrNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("failed to list packages: %w", err)
	}
	return pkgs, nil
}

// ListPackageFiles returns the list of files in the given model package.
func (h *handler) ReadPackage(ctx context.Context, p *gentypes.PackageLocator) ([]*gentypes.PackageFile, error) {
	dir := filepath.Join(p.Repository, p.Dir)
	files, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to list files in %s", p.Dir)
	}
	var res []*gentypes.PackageFile
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if filepath.Ext(f.Name()) != ".go" {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s", f.Name())
		}
		res = append(res, &gentypes.PackageFile{
			Locator: &gentypes.FileLocator{
				Repository: p.Repository,
				Dir:        p.Dir,
				Filename:   f.Name(),
			},
			Content: string(b),
		})
	}
	return res, nil
}

// GetImportPath returns the import path for the given package.
func (h *handler) GetImportPath(ctx context.Context, p *gentypes.PackageLocator) (string, error) {
	dir := filepath.Join(h.root, p.Repository, p.Dir)
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("failed to stat directory %s: %w", dir, err)
	}
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.PackageClauseOnly)
	if err != nil {
		return "", fmt.Errorf("failed to parse %s", dir)
	}
	var pkgPath string
	for _, pkg := range pkgs {
		pkgPath = pkg.Name
		break
	}
	if pkgPath == "" {
		return "", fmt.Errorf("failed to find package in %s", dir)
	}
	gomod, modpath, err := modulePath(dir)
	if err != nil {
		return "", fmt.Errorf("failed to find module for %s", dir)
	}
	rel, err := filepath.Rel(filepath.Dir(gomod), dir)
	if err != nil {
		return "", fmt.Errorf("failed to compute relative path for %s relative to %s", dir, gomod)
	}
	rel = filepath.ToSlash(rel)
	return path.Join(modpath, rel, pkgPath), nil
}

// Save implements Save.
func (h *handler) Save(ctx context.Context, p *gentypes.PackageFile) error {
	dir := filepath.Join(h.root, p.Locator.Repository, p.Locator.Dir)
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}
		} else {
			return fmt.Errorf("failed to stat directory %s: %w", dir, err)
		}
	}
	err := os.WriteFile(filepath.Join(dir, p.Locator.Filename), []byte(p.Content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", filepath.Join(dir, p.Locator.Filename), err)
	}
	return nil
}

// Subscribe listens to changes in the given model package files. When
// notifications are received from the filesystem, it notifies the returned
// channel.
func (h *handler) Subscribe(ctx context.Context, p *gentypes.PackageLocator) (<-chan struct{}, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Errorf(ctx, err, "Error watching files")
		return nil, err
	}
	dir := filepath.Join(h.root, p.Repository, p.Dir)
	if err = watcher.Add(dir); err != nil {
		log.Errorf(ctx, err, "Error watching files in %s", dir)
		return nil, err
	}
	ch := make(chan struct{})
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
				fset := token.NewFileSet()
				pkgs, err := parser.ParseDir(fset, dir, nil, parser.PackageClauseOnly)
				if err != nil {
					log.Errorf(ctx, err, "failed to stat %s", dir)
					continue
				}
				var pkgPath string
				for _, pkg := range pkgs {
					pkgPath = pkg.Name
					break
				}
				if pkgPath == "" {
					log.Error(ctx, fmt.Errorf("failed to find package in %s", dir))
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
				ch <- struct{}{}

			case err := <-watcher.Errors:
				log.Errorf(ctx, err, "Error watching files")

			case <-ctx.Done():
				return
			}
		}
	}()

	return ch, nil
}

func iterateModelPackages(ctx context.Context, root string, fn func(dirPath string, importPath string)) error {
	return filepath.Walk(root, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return ErrNotFound
			}
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
		gomod, modpath, err := modulePath(fpath)
		if err != nil {
			log.Errorf(ctx, err, "failed to find module for %s", fpath)
			return nil
		}
		for _, pkg := range parsed {
			for _, f := range pkg.Files {
				for _, imp := range f.Imports {
					if imp.Path.Value == `"goa.design/model/dsl"` {
						rel, err := filepath.Rel(filepath.Dir(gomod), fpath)
						if err != nil {
							log.Errorf(ctx, err, "failed to compute relative path for %s relative to %s", fpath, gomod)
						}
						rel = filepath.ToSlash(rel)
						gopath := path.Join(modpath, rel, pkg.Name)
						fn(fpath, gopath)
					}
				}
			}
		}
		return nil
	})
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
