package svc

import (
	"context"

	"goa.design/model/svc/clients/repo"
)

type (
	// Service implements the model service business logic.
	// It exposes methods for manipulating a model DSL.
	Service struct {
		dir     string // Models root directory
		debug   bool   // Whether to print debug output when generating
		handler repo.RepoHandler
	}
)

// New returns the Service implementation.
func New(ctx context.Context, dir string, handler repo.RepoHandler, debug bool) (*Service, error) {
	return &Service{
		dir:     dir,
		debug:   debug,
		handler: handler,
	}, nil
}
