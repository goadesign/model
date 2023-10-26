package mdlsvc

import (
	"context"

	pstore "goa.design/model/mdlsvc/clients/package_store"
)

type (
	// Service implements the model service business logic.
	// It exposes methods for manipulating a model DSL.
	Service struct {
		dir   string // Models root directory
		debug bool   // Whether to print debug output when generating
		store pstore.PackageStore
	}
)

// New returns the Service implementation.
func New(ctx context.Context, dir string, store pstore.PackageStore, debug bool) (*Service, error) {
	return &Service{
		dir:   dir,
		debug: debug,
		store: store,
	}, nil
}
