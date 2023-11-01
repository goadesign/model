package svc

import (
	"context"

	"goa.design/model/svc/clients/repo"
)

type (
	// Service implements the model service business logic.
	// It exposes methods for manipulating a model DSL.
	Service struct {
		handler repo.RepoHandler
		debug   bool // Whether to print debug output when generating
	}
)

// New returns the Service implementation.
func New(ctx context.Context, handler repo.RepoHandler, debug bool) (*Service, error) {
	return &Service{
		handler: handler,
		debug:   debug,
	}, nil
}
