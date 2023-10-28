package svc

import (
	"context"

	gensvg "goa.design/model/svc/gen/svg"
)

// Returns the SVG for the given file.
func (svc *Service) Load(ctx context.Context, file *gensvg.Filename) (gensvg.SVG, error) {
	panic("not implemented")
}

// Save the SVG to the given file.
func (svc *Service) Save(ctx context.Context, payload *gensvg.SavePayload) error {
	panic("not implemented")
}
