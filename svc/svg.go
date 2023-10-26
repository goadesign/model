package svc

import (
	"context"

	gensvg "goa.design/model/svc/gen/svg"
)

// Stream the model layout JSON saved in the SVG
func (svc *Service) Load(ctx context.Context, file *gensvg.Filename) (gensvg.SVG, error) {
	panic("not implemented")
}

// Save the SVG streamed in the request body
func (svc *Service) Save(ctx context.Context, payload *gensvg.SavePayload) error {
	panic("not implemented")
}
