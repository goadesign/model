package mdlsvc

import (
	"context"
	"sync"

	genpackages "goa.design/model/mdlsvc/gen/packages"
)

type (
	// Service implements the model service business logic.
	// It exposes methods for manipulating a model DSL.
	Service struct {
		dir           string // SVG generation directory
		workspace     string // Go workspace
		debug         bool   // Whether to print debug output when generating
		lock          sync.Mutex
		subscriptions map[string]*subscription // subscriptions indexed by module
	}

	// subscription represents a client subscription to model updates.
	subscription struct {
		ch      chan []byte                         // Channel used to notify of model changes.
		streams []genpackages.SubscribeServerStream // List of active client streams.
	}
)

// New returns the Service implementation.
func New(ctx context.Context, workspace, dir string, debug bool) (*Service, error) {
	return &Service{
		dir:           dir,
		workspace:     workspace,
		debug:         debug,
		subscriptions: make(map[string]*subscription),
	}, nil
}
